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
	state               verificationState
	parent              gtki.Window
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

func (n *notifier) stopNotifying(i gtki.InfoBar) {

}

type verificationState int

const (
	unverified verificationState = iota
	success
	waitingForAnswerFromPeer
)

func newVerifier(conv *conversationPane) *verifier {
	v := &verifier{session: conv.account.session}
	v.notifier = &notifier{conv.notificationArea}
	peer, ok := conv.currentPeer()
	if !ok {
		// ???
	}
	v.verificationWarning = v.buildStartVerificationNotification(peer, conv.transientParent, conv.currentResource())
	v.notifier.notify(v.verificationWarning)
	return v
}

func (v *verifier) buildStartVerificationNotification(peer *rosters.Peer, parent gtki.Window, resource string) gtki.InfoBar {
	builder := newBuilder("StartVerificationNotification")
	infoBar := builder.getObj("infobar").(gtki.InfoBar)
	message := builder.getObj("message").(gtki.Label)
	message.SetText(i18n.Local("Make sure no one else is reading your messages."))
	button := builder.getObj("button_verify").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			d := v.showNewPinDialog(peer, parent, resource, infoBar)
			d.Run()
			d.Destroy()
		})
	})
	infoBar.ShowAll()
	return infoBar
}

func (v *verifier) showNewPinDialog(peer *rosters.Peer, parent gtki.Window, resource string, infoBar gtki.InfoBar) gtki.Dialog {
	pinBuilder := newBuilder("GeneratePIN")
	sharePINDialog := pinBuilder.getObj("dialog").(gtki.Dialog)
	msg := pinBuilder.getObj("SharePinLabel").(gtki.Label)
	msg.SetText(fmt.Sprintf(i18n.Local("Share the one-time PIN below with %s"), peer.NameForPresentation()))
	var pinLabel gtki.Label
	pinBuilder.getItems(
		"PinLabel", &pinLabel,
	)
	pin, err := createPIN()
	if err != nil {
		if v.verificationWarning != nil {
			v.verificationWarning.Hide()
		}
		v.showNotificationWhenWeCannotGeneratePINForSMP(err, sharePINDialog)
		return sharePINDialog
	}
	pinBuilder.ConnectSignals(map[string]interface{}{
		"on_gen_pin": func() {
			pin, err = createPIN()
			if err != nil {
				if v.verificationWarning != nil {
					v.verificationWarning.Hide()
				}
				v.showNotificationWhenWeCannotGeneratePINForSMP(err, sharePINDialog)
				return
			}
			pinLabel.SetText(pin)
		},
		"close_share_pin": func() {
			if v.peerRequestsSMP != nil {
				showSMPHasAlreadyStarted(peer.NameForPresentation(), sharePINDialog)
				return
			}
			infoBar.Hide()
			v.showWaitingForPeerToCompleteSMPDialog(peer.NameForPresentation(), sharePINDialog)
			v.session.StartSMP(peer.Jid, resource, i18n.Local("Please enter the PIN that your contact shared with you."), pin)
		},
	})
	pinLabel.SetText(pin)
	sharePINDialog.SetTransientFor(parent)
	sharePINDialog.ShowAll()
	return sharePINDialog
}

func (v *verifier) showWaitingForPeerToCompleteSMPDialog(peer string, sharePINDialog gtki.Dialog) {
	builderWaitingSMP := newBuilder("WaitingSMPComplete")
	waitingInfoBar := builderWaitingSMP.getObj("smp_waiting_infobar").(gtki.InfoBar)
	waitingSMPMessage := builderWaitingSMP.getObj("message").(gtki.Label)
	waitingSMPMessage.SetText(fmt.Sprintf(i18n.Local("Waiting for %s to finish securing the channel..."), peer))
	waitingInfoBar.ShowAll()
	v.waitingForSMP = waitingInfoBar
	v.notifier.notify(waitingInfoBar)
	sharePINDialog.Destroy()
}

func (v *verifier) showNotificationWhenWeCannotGeneratePINForSMP(err error, pinDialog gtki.Dialog) {
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
	v.notifier.notify(errInfoBar)
}

func (v *verifier) displayRequestForSecret(peer *rosters.Peer, parent gtki.Window, resource string) {
	if v.verificationWarning != nil {
		v.verificationWarning.Hide()
	}
	b := newBuilder("PeerRequestsSMP")
	infobar := b.getObj("peer_requests_smp").(gtki.InfoBar)
	infobarMsg := b.getObj("message").(gtki.Label)
	verificationButton := b.getObj("verification_button").(gtki.Button)
	verificationButton.Connect("clicked", func() {
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
				v.session.FinishSMP(peer.Jid, resource, pin)
				d.Destroy()
			},
		})
		d.SetTransientFor(parent)
		d.ShowAll()
		d.Run()
		d.Destroy()
	})
	message := fmt.Sprintf("%s is waiting for you to finish verifying the security of this channel...", peer.NameForPresentation())
	infobarMsg.SetText(i18n.Local(message))
	infobar.ShowAll()
	v.peerRequestsSMP = infobar
	v.notifier.notify(infobar)
}

func (v *verifier) displayVerificationSuccess(peer *rosters.Peer, parent gtki.Window) {
	builder := newBuilder("VerificationSucceeded")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("Horray! No one is listening in on your conversations with %s", peer.NameForPresentation())))
	button := builder.getObj("button_ok").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			d.Destroy()
		})
	})
	d.SetTransientFor(parent)
	d.ShowAll()
	if v.waitingForSMP != nil {
		v.waitingForSMP.Destroy()
	}
	if v.verificationWarning != nil {
		v.verificationWarning.Destroy()
	}
	if v.peerRequestsSMP != nil {
		v.peerRequestsSMP.Destroy()
		v.peerRequestsSMP = nil
	}
}

func (v *verifier) displayVerificationFailure(peer *rosters.Peer, parent gtki.Window) {
	builder := newBuilder("VerificationFailed")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Local(fmt.Sprintf("We failed to verify this channel with %s\n\n Maybe:", peer.NameForPresentation())))
	tryLaterButton := builder.getObj("try_later").(gtki.Button)
	tryLaterButton.Connect("clicked", func() {
		doInUIThread(func() {
			// TODO: This is hacky and the checks will only apply to one of the peers at a time. We should do something better.
			if v.peerRequestsSMP != nil {
				v.peerRequestsSMP.Destroy()
				v.peerRequestsSMP = nil
			}
			if v.waitingForSMP != nil {
				v.waitingForSMP.Destroy()
			}
			v.verificationWarning.Show()
			d.Destroy()
		})
	})
	d.SetTransientFor(parent)
	d.ShowAll()

}

func (v *verifier) handle(ev events.SMP, peer *rosters.Peer, parent gtki.Window, resource string) {
	switch ev.Type {
	case events.SecretNeeded:
		v.displayRequestForSecret(peer, parent, resource)
	case events.Success:
		v.displayVerificationSuccess(peer, parent)
	case events.Failure:
		v.displayVerificationFailure(peer, parent)
	}
}
