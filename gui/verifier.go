package gui

import (
	"fmt"
	"log"

	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/gotk3adapter/gtki"
)

type verifier struct {
	state        verificationState
	newPinDialog gtki.Dialog
}

type verificationState int

const (
	unverified verificationState = iota
	success
	waitingForAnswerFromPeer
)

func newVerifier(conv *conversationPane) *verifier {
	v := &verifier{}
	conv.verificationWarning = v.buildStartVerificationNotification(conv)
	conv.addNotification(conv.verificationWarning)
	return v
}

func (v *verifier) buildStartVerificationNotification(convPane *conversationPane) gtki.InfoBar {
	builder := newBuilder("StartVerificationNotification")
	infoBar := builder.getObj("infobar").(gtki.InfoBar)
	message := builder.getObj("message").(gtki.Label)
	message.SetText(i18n.Local("Make sure no one else is reading your messages."))
	button := builder.getObj("button_verify").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			d := v.showNewPinDialog(convPane.transientParent, convPane, infoBar)
			d.Run()
			d.Destroy()
		})
	})
	infoBar.ShowAll()
	return infoBar
}

func (v *verifier) showNewPinDialog(parent gtki.Window, conv *conversationPane, infoBar gtki.InfoBar) gtki.Dialog {
	pinBuilder := newBuilder("GeneratePIN")
	sharePINDialog := pinBuilder.getObj("dialog").(gtki.Dialog)
	peer, ok := conv.currentPeer()
	if !ok {
		// ???
	}
	msg := pinBuilder.getObj("SharePinLabel").(gtki.Label)
	msg.SetText(fmt.Sprintf(i18n.Local("Share the one-time PIN below with %s"), peer.NameForPresentation()))

	var pinLabel gtki.Label
	pinBuilder.getItems(
		"PinLabel", &pinLabel,
	)
	pin, err := createPIN()
	if err != nil {
		if conv.verificationWarning != nil {
			conv.verificationWarning.Hide()
		}
		v.showNotificationWhenWeCannotGeneratePINForSMP(err, sharePINDialog, conv)
		return sharePINDialog
	}
	pinBuilder.ConnectSignals(map[string]interface{}{
		"on_gen_pin": func() {
			pin, err = createPIN()
			if err != nil {
				if conv.verificationWarning != nil {
					conv.verificationWarning.Hide()
				}
				v.showNotificationWhenWeCannotGeneratePINForSMP(err, sharePINDialog, conv)
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
			if conv.peerRequestsSMP != nil {
				peer, ok := conv.currentPeer()
				if !ok {
					// ????
				}
				showSMPHasAlreadyStarted(peer.NameForPresentation(), sharePINDialog)
				return
			}
			v.showWaitingForPeerToCompleteSMPDialog(peer, infoBar, sharePINDialog, conv)
			conv.account.session.StartSMP(peer.Jid, conv.currentResource(), i18n.Local("Please enter the PIN that your contact shared with you."), pin)
		},
	})
	pinLabel.SetText(pin)
	sharePINDialog.SetTransientFor(parent)
	sharePINDialog.ShowAll()
	return sharePINDialog
}

func (v *verifier) showWaitingForPeerToCompleteSMPDialog(peer *rosters.Peer, infoBar gtki.InfoBar, sharePINDialog gtki.Dialog, conv *conversationPane) {
	builderWaitingSMP := newBuilder("WaitingSMPComplete")
	waitingInfoBar := builderWaitingSMP.getObj("smp_waiting_infobar").(gtki.InfoBar)
	waitingSMPMessage := builderWaitingSMP.getObj("message").(gtki.Label)
	waitingSMPMessage.SetText(fmt.Sprintf(i18n.Local("Waiting for %s to finish securing the channel..."), peer.NameForPresentation()))
	infoBar.Hide()
	waitingInfoBar.ShowAll()
	conv.waitingForSMP = waitingInfoBar
	conv.addNotification(waitingInfoBar)
	sharePINDialog.Destroy()
}

func (v *verifier) showNotificationWhenWeCannotGeneratePINForSMP(err error, pinDialog gtki.Dialog, conv *conversationPane) {
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
