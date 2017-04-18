package gui

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/gotk3adapter/gtki"
)

func (u *gtkUI) showConnectAccountNotification(account *account) func() {
	var notification gtki.InfoBar

	doInUIThread(func() {
		notification = account.buildConnectionNotification(u)
		account.setCurrentNotification(notification, u.notificationArea)
	})

	return func() {
		doInUIThread(func() {
			account.removeCurrentNotificationIf(notification)
		})
	}
}

func (u *gtkUI) notifyTorIsNotRunning(account *account) {
	doInUIThread(func() {
		notification := account.buildTorNotRunningNotification(u)
		account.setCurrentNotification(notification, u.notificationArea)
	})
}

func (u *gtkUI) notifyConnectionFailure(account *account, moreInfo func()) {
	doInUIThread(func() {
		notification := account.buildConnectionFailureNotification(u, moreInfo)
		account.setCurrentNotification(notification, u.notificationArea)
	})
}

func buildStartVerificationNotification(convPane *conversationPane) gtki.InfoBar {
	builder := newBuilder("StartVerificationNotification")
	infoBar := builder.getObj("infobar").(gtki.InfoBar)
	message := builder.getObj("message").(gtki.Label)
	message.SetText(i18n.Local("Make sure no one else is reading your messages."))
	button := builder.getObj("button_verify").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			showNewPinDialog(convPane.transientParent, convPane, infoBar)
		})
	})
	infoBar.ShowAll()
	return infoBar
}

func showNotificationWhenWeCannotGeneratePINForSMP(err error, pinDialog gtki.Dialog, conv *conversationPane) {
	log.Printf("Cannot recover from error: %v. Quitting verification using SMP.", err)
	pinDialog.Destroy()
	errBuilder := newBuilder("CannotVerifyWithSMP")
	errInfoBar := errBuilder.getObj("error_verifying_smp").(gtki.InfoBar)
	message := errBuilder.getObj("message").(gtki.Label)
	message.SetText(i18n.Local("Unable to verify the channel at this time."))
	button := errBuilder.getObj("try_later_button").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			errInfoBar.Destroy()
		})
	})
	errInfoBar.ShowAll()
	conv.addNotification(errInfoBar)
}

func showWaitingForPeerToCompleteSMPDialog(peer *rosters.Peer, infoBar gtki.InfoBar, sharePinDialog gtki.Dialog, conv *conversationPane) {
	builderWaitingSMP := newBuilder("WaitingSMPComplete")
	waitingInfoBar := builderWaitingSMP.getObj("smp_waiting_infobar").(gtki.InfoBar)
	waitingSMPMessage := builderWaitingSMP.getObj("message").(gtki.Label)
	waitingSMPMessage.SetText(i18n.Local(fmt.Sprintf("Waiting for %s to finish securing the channel...", peer.NameForPresentation())))
	infoBar.Hide()
	waitingInfoBar.ShowAll()
	conv.waitingForSMP = waitingInfoBar
	conv.addNotification(waitingInfoBar)
	sharePinDialog.Destroy()
}

func showNewPinDialog(parent gtki.Window, conv *conversationPane, infoBar gtki.InfoBar) {
	pinBuilder := newBuilder("GeneratePIN")
	sharePinDialog := pinBuilder.getObj("dialog").(gtki.Dialog)
	var pinLabel gtki.Label
	pinBuilder.getItems(
		"PinLabel", &pinLabel,
	)
	pin, err := createPIN()
	if err != nil {
		if conv.verificationWarning != nil {
			conv.verificationWarning.Hide()
		}
		showNotificationWhenWeCannotGeneratePINForSMP(err, sharePinDialog, conv)
		return
	}
	pinBuilder.ConnectSignals(map[string]interface{}{
		"on_gen_pin": func() {
			pin, err = createPIN()
			if err != nil {
				if conv.verificationWarning != nil {
					conv.verificationWarning.Hide()
				}
				showNotificationWhenWeCannotGeneratePINForSMP(err, sharePinDialog, conv)
				return
			}
			pinLabel.SetText(pin)
		},
		"close_share_pin": func() {
			peer, ok := conv.currentPeer()
			if !ok {
				// print that contact does not exist? this is impossible situation
				return
			}
			showWaitingForPeerToCompleteSMPDialog(peer, infoBar, sharePinDialog, conv)
			conv.account.session.StartSMP(peer.Jid, conv.currentResource(), "Please enter the PIN that your contact shared with you.", pin)
		},
	})
	pinLabel.SetText(pin)
	sharePinDialog.SetTransientFor(parent)
	sharePinDialog.ShowAll()
}

func createPIN() (string, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(1000000)))
	if err != nil {
		log.Printf("Error encountered when creating a new PIN: %v", err)
		return "", err
	}
	return fmt.Sprintf("%06d", val), nil
}

func verifyChannelDialog(conv *conversationPane, infoBar gtki.InfoBar) {
	builder := newBuilder("VerifyChannel")
	d := builder.getObj("dialog").(gtki.Dialog)
	submit := builder.getObj("button_submit").(gtki.Button)
	submit.Connect("clicked", func() {
		doInUIThread(func() {
			e := builder.getObj("pin").(gtki.Entry)
			pin, err := e.GetText()
			if err != nil {
				// TODO: Don't allow blank submission - require a PIN. This is a WIP.

				// notificationBuilder := newBuilder("BadPINNotification")
				// notification := notificationBuilder.getObj("infoBar").(gtki.InfoBar)
				// msg := notificationBuilder.getObj("message").(gtki.Label)
				// msg.SetText("A PIN is required")
				// area := builder.getObj("notification-area").(gtki.Box)
				// area.Add(notification)
			}
			peer, ok := conv.currentPeer()
			if !ok {
				// TODO: handle when getting the current peer fails
			}
			conv.account.session.FinishSMP(peer.Jid, conv.currentResource(), pin)
			d.Destroy()
		})
	})
	d.SetTransientFor(conv.transientParent)
	d.ShowAll()
}

func showPINWasIncorrect(parent gtki.Window) {
	builder := newBuilder("OopsWrongPIN")
	d := builder.getObj("dialog").(gtki.Dialog)
	continueButton := builder.getObj("continue_anyway").(gtki.Button)
	continueButton.Connect("clicked", func() {
		doInUIThread(func() {
			d.Destroy()
		})
	})
	getOutButton := builder.getObj("get_out").(gtki.Button)
	getOutButton.Connect("clicked", func() {
		doInUIThread(func() {
			d.Destroy()
			parent.Hide()
		})
	})
	d.SetTransientFor(parent)
	d.ShowAll()
}

func showThatChannelIsVerified(peerName string, parent gtki.Window) {
	builder := newBuilder("SecureChannelVerified")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("Horray! No one is listening in on your conversations with %s", peerName)))
	button := builder.getObj("button_ok").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			d.Destroy()
		})
	})
	d.SetTransientFor(parent)
	d.ShowAll()
}

func (u *gtkUI) notify(title, message string) {
	builder := newBuilder("SimpleNotification")
	obj := builder.getObj("dialog")
	dlg := obj.(gtki.MessageDialog)

	dlg.SetProperty("title", title)
	dlg.SetProperty("text", message)
	dlg.SetTransientFor(u.window)

	doInUIThread(func() {
		dlg.Run()
		dlg.Destroy()
	})
}
