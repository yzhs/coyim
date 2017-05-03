package gui

import (
	"fmt"
	"log"

	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/gotk3adapter/gtki"
)

type verifier struct {
	state               verifierState
	peer                *rosters.Peer
	parent              gtki.Window
	currentResource     string
	session             access.Session
	notifier            *notifier
	newPinDialog        gtki.Dialog
	waitingForSMP       gtki.InfoBar
	peerRequestsSMP     gtki.InfoBar
	verificationWarning gtki.InfoBar
}

type notifier struct {
	notificationArea gtki.Box
}

func (n *notifier) notify(i gtki.InfoBar) {
	n.notificationArea.Add(i)
}

type verifierState int

const (
	unverified verifierState = iota
	peerRequestsSMP
	waitingForAnswerFromPeer
	success
	failure
	smpErr
)

func newVerifier(conv *conversationPane) *verifier {
	v := &verifier{session: conv.account.session}
	v.notifier = &notifier{conv.notificationArea}
	peer, ok := conv.currentPeer()
	if !ok {
		// ???
	}
	v.peer = peer
	v.parent = conv.transientParent
	v.currentResource = conv.currentResource()
	v.buildStartVerificationNotification()
	v.notifier.notify(v.verificationWarning)
	return v
}

func (v *verifier) buildStartVerificationNotification() {
	builder := newBuilder("StartVerificationNotification")
	v.verificationWarning = builder.getObj("infobar").(gtki.InfoBar)
	message := builder.getObj("message").(gtki.Label)
	message.SetText(i18n.Local("Make sure no one else is reading your messages."))
	button := builder.getObj("button_verify").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			d := v.showNewPinDialog()
			d.Run()
			d.Destroy()
		})
	})
	v.verificationWarning.ShowAll()
}

func (v *verifier) smpError(err error) {
	v.state = smpErr
	v.disableNotifications()
	v.showNotificationWhenWeCannotGeneratePINForSMP(err)
}

func (v *verifier) showNewPinDialog() gtki.Dialog {
	pinBuilder := newBuilder("GeneratePIN")
	sharePINDialog := pinBuilder.getObj("dialog").(gtki.Dialog)
	msg := pinBuilder.getObj("SharePinLabel").(gtki.Label)
	msg.SetText(fmt.Sprintf(i18n.Local("Share the one-time PIN below with %s"), v.peer.NameForPresentation()))
	var pinLabel gtki.Label
	pinBuilder.getItems(
		"PinLabel", &pinLabel,
	)
	pin, err := createPIN()
	if err != nil {
		sharePINDialog.Destroy()
		v.smpError(err)
		return sharePINDialog
	}
	pinBuilder.ConnectSignals(map[string]interface{}{
		"on_gen_pin": func() {
			pin, err = createPIN()
			if err != nil {
				sharePINDialog.Destroy()
				v.smpError(err)
			}
			pinLabel.SetText(pin)
		},
		"close_share_pin": func() {
			if v.peerRequestsSMP != nil {
				showSMPHasAlreadyStarted(v.peer.NameForPresentation(), sharePINDialog)
				return
			}
			v.showWaitingForPeerToCompleteSMPDialog(sharePINDialog)
			v.session.StartSMP(v.peer.Jid, v.currentResource, i18n.Local("Please enter the PIN that your contact shared with you."), pin)
		},
	})
	pinLabel.SetText(pin)
	sharePINDialog.SetTransientFor(v.parent)
	sharePINDialog.ShowAll()
	return sharePINDialog
}

func (v *verifier) showWaitingForPeerToCompleteSMPDialog(sharePINDialog gtki.Dialog) {
	v.state = waitingForAnswerFromPeer
	v.disableNotifications()
	builderWaitingSMP := newBuilder("WaitingSMPComplete")
	waitingInfoBar := builderWaitingSMP.getObj("smp_waiting_infobar").(gtki.InfoBar)
	waitingSMPMessage := builderWaitingSMP.getObj("message").(gtki.Label)
	waitingSMPMessage.SetText(fmt.Sprintf(i18n.Local("Waiting for %s to finish securing the channel..."), v.peer.NameForPresentation()))
	waitingInfoBar.ShowAll()
	v.waitingForSMP = waitingInfoBar
	v.notifier.notify(waitingInfoBar)
	sharePINDialog.Destroy()
}

func (v *verifier) showNotificationWhenWeCannotGeneratePINForSMP(err error) {
	log.Printf("Cannot recover from error: %v. Quitting verification using SMP.", err)
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
	v.notifier.notify(errInfoBar)
}

func (v *verifier) displayRequestForSecret() {
	v.disableNotifications()
	b := newBuilder("PeerRequestsSMP")
	infobar := b.getObj("peer_requests_smp").(gtki.InfoBar)
	infobarMsg := b.getObj("message").(gtki.Label)
	verificationButton := b.getObj("verification_button").(gtki.Button)
	verificationButton.Connect("clicked", func() {
		builder := newBuilder("EnterPIN")
		d := builder.getObj("dialog").(gtki.Dialog)
		msg := builder.getObj("verification_message").(gtki.Label)
		msg.SetText(i18n.Local(fmt.Sprintf("Type the PIN that %s sent you", v.peer.NameForPresentation())))
		button := builder.getObj("button_submit").(gtki.Button)
		button.SetSensitive(false)
		builder.ConnectSignals(map[string]interface{}{
			"on_text_entry": func() {
				e := builder.getObj("pin").(gtki.Entry)
				pin, _ := e.GetText()
				button.SetSensitive(len(pin) > 0)
			},
			"close_share_pin": func() {
				e := builder.getObj("pin").(gtki.Entry)
				pin, _ := e.GetText()
				v.session.FinishSMP(v.peer.Jid, v.currentResource, pin)
				d.Destroy()
			},
		})
		d.SetTransientFor(v.parent)
		d.ShowAll()
		d.Run()
		d.Destroy()
	})
	message := fmt.Sprintf("%s is waiting for you to finish verifying the security of this channel...", v.peer.NameForPresentation())
	infobarMsg.SetText(i18n.Local(message))
	infobar.ShowAll()
	v.peerRequestsSMP = infobar
	v.notifier.notify(infobar)
}

func (v *verifier) displayVerificationSuccess() {
	builder := newBuilder("VerificationSucceeded")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("Horray! No one is listening in on your conversations with %s", v.peer.NameForPresentation())))
	button := builder.getObj("button_ok").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			d.Destroy()
		})
	})
	d.SetTransientFor(v.parent)
	d.ShowAll()
	v.disableNotifications()
}

func (v *verifier) displayVerificationFailure() {
	builder := newBuilder("VerificationFailed")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("We failed to verify this channel with %s\n\n Maybe:", v.peer.NameForPresentation())))
	tryLaterButton := builder.getObj("try_later").(gtki.Button)
	tryLaterButton.Connect("clicked", func() {
		doInUIThread(func() {
			v.disableNotifications()
			d.Destroy()
		})
	})
	d.SetTransientFor(v.parent)
	d.ShowAll()
}

func (v *verifier) disableNotifications() {
	switch v.state {
	case success:
		v.removeInProgressNotifications()
		v.verificationWarning.Destroy()
	case failure:
		v.removeInProgressNotifications()
		v.verificationWarning.Show()
	case waitingForAnswerFromPeer, peerRequestsSMP, smpErr:
		v.verificationWarning.Hide()
	}
}

func (v *verifier) removeInProgressNotifications() {
	if v.peerRequestsSMP != nil {
		v.peerRequestsSMP.Destroy()
		v.peerRequestsSMP = nil
	}
	if v.waitingForSMP != nil {
		v.waitingForSMP.Destroy()
		v.waitingForSMP = nil
	}
}

func (v *verifier) handle(ev events.SMP) {
	switch ev.Type {
	case events.SecretNeeded:
		v.state = peerRequestsSMP
		v.displayRequestForSecret()
	case events.Success:
		v.state = success
		v.displayVerificationSuccess()
	case events.Failure:
		v.state = failure
		v.displayVerificationFailure()
	}
}
