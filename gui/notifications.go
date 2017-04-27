package gui

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/session/access"
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

func showSMPHasAlreadyStarted(peerName string, parent gtki.Window) {
	b := newBuilder("SMPHasAlreadyStarted")
	d := b.getObj("dialog").(gtki.Dialog)
	msg := b.getObj("smp_has_already_started").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("%s has already started verification and generated a PIN.\nPlease ask them for it.", peerName)))
	button := b.getObj("button_ok").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			d.Destroy()
		})
	})
	d.SetTransientFor(parent)
	d.ShowAll()
}

func createPIN() (string, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(1000000)))
	if err != nil {
		log.Printf("Error encountered when creating a new PIN: %v", err)
		return "", err
	}
	return fmt.Sprintf("%06d", val), nil
}

type genPinDialog struct {
	noPINNotification gtki.InfoBar
}

func pinInputDialog(peer *rosters.Peer, session access.Session, parent gtki.Window, currentResource string) {
	gpDialog := &genPinDialog{}
	builder := newBuilder("EnterPIN")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("Type the PIN that %s sent you", peer.NameForPresentation())))
	builder.ConnectSignals(map[string]interface{}{
		"close_share_pin": func() {
			e := builder.getObj("pin").(gtki.Entry)
			pin, _ := e.GetText()
			if pin == "" {
				area := builder.getObj("notification-area").(gtki.Box)
				if gpDialog.noPINNotification != nil {
					area.Remove(gpDialog.noPINNotification)
				}
				notificationBuilder := newBuilder("NoPINNotification")
				gpDialog.noPINNotification = notificationBuilder.getObj("infobar").(gtki.InfoBar)
				msg := notificationBuilder.getObj("message").(gtki.Label)
				msg.SetText(i18n.Local("PIN is required"))
				area.Add(gpDialog.noPINNotification)
				area.ShowAll()
				d.Run()
				return
			}
			session.FinishSMP(peer.Jid, currentResource, pin)
			d.Destroy()
		},
	})

	d.SetTransientFor(parent)
	d.ShowAll()
	d.Run()
	d.Destroy()
	return
}

func verifyChannelDialog(conv *conversationPane, infoBar gtki.InfoBar) {
	peer, ok := conv.currentPeer()
	if !ok {
		// ????
	}
	pinInputDialog(peer, conv.account.session, conv.transientParent, conv.currentResource())
}

func showThatVerificationFailed(peer string, conv *conversationPane, parent gtki.Window) {
	builder := newBuilder("VerificationFailed")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("We failed to verify this channel with %s\n\n Maybe:", peer)))
	tryLaterButton := builder.getObj("try_later").(gtki.Button)
	tryLaterButton.Connect("clicked", func() {
		doInUIThread(func() {
			// TODO: This is hacky and the checks will only apply to one of the peers at a time. We should do something better.
			if conv.peerRequestsSMP != nil {
				conv.peerRequestsSMP.Destroy()
				conv.peerRequestsSMP = nil
			}
			if conv.waitingForSMP != nil {
				conv.waitingForSMP.Destroy()
			}
			conv.verificationWarning.Show()
			d.Destroy()
		})
	})
	d.SetTransientFor(parent)
	d.ShowAll()
}

func showThatVerificationSucceeded(peer string, parent gtki.Window) {
	builder := newBuilder("VerificationSucceeded")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("Horray! No one is listening in on your conversations with %s", peer)))
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
