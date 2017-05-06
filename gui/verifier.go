package gui

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

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
	enterPinDialog      gtki.Dialog
	waitingForSMP       gtki.InfoBar
	peerRequestsSMP     gtki.InfoBar
	verificationWarning gtki.InfoBar
	displayController   *displaySettings
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
	finishedAnsweringPeer
	success
	failure
	smpErr
	cancelingSMP
)

func newVerifier(conv *conversationPane, ds *displaySettings) *verifier {
	v := &verifier{session: conv.account.session}
	v.notifier = &notifier{conv.notificationArea}
	peer, ok := conv.currentPeer()
	if !ok {
		// ???
	}
	v.peer = peer
	v.parent = conv.transientParent
	v.currentResource = conv.currentResource()
	v.displayController = ds
	return v
}

func (v *verifier) buildStartVerificationNotification() {
	builder := newBuilder("StartVerificationNotification")
	v.verificationWarning = builder.getObj("infobar").(gtki.InfoBar)
	message := builder.getObj("message").(gtki.Label)
	message.SetText(i18n.Local("Tenha certeza que ninguém está lendo suas mensagens."))
	sc, _ := message.GetStyleContext()
	sc.AddProvider(v.displayController.provider, 9999)
	sc.AddClass("start-verification-msg")
	v.displayController.provider.LoadFromData(`
	    .start-verification-msg {
	        font-size: 16px;
	    }`)
	builder.ConnectSignals(map[string]interface{}{
		"close_verification": func() {
			v.verificationWarning.Hide()
		},
	})
	button := builder.getObj("button_verify").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			if v.newPinDialog != nil {
				v.newPinDialog.Present()
			} else {
				v.showNewPinDialog()
			}
		})
	})
	v.verificationWarning.ShowAll()
	v.notifier.notify(v.verificationWarning)
}

func (v *verifier) smpError(err error) {
	v.state = smpErr
	v.disableNotifications()
	v.showNotificationWhenWeCannotGeneratePINForSMP(err)
}

func (v *verifier) showNewPinDialog() {
	pinBuilder := newBuilder("GeneratePIN")
	v.newPinDialog = pinBuilder.getObj("dialog").(gtki.Dialog)
	v.newPinDialog.HideOnDelete()
	msg := pinBuilder.getObj("SharePinLabel").(gtki.Label)
	msg.SetText(i18n.Localf("Compartilhe o código temporário que está abaixo com %s", v.peer.NameForPresentation()))
	var pinLabel gtki.Label
	pinBuilder.getItems(
		"PinLabel", &pinLabel,
	)
	sc, _ := pinLabel.GetStyleContext()
	sc.AddProvider(v.displayController.provider, 9999)
	sc.AddClass("new-pin")
	v.displayController.provider.LoadFromData(
		`.new-pin {
			font-size: 30px;
			font-weight: 800;
		}`)
	pin, err := createPIN()
	if err != nil {
		v.newPinDialog.Destroy()
		v.smpError(err)
	}
	pinBuilder.ConnectSignals(map[string]interface{}{
		"on_gen_pin": func() {
			pin, err = createPIN()
			if err != nil {
				v.newPinDialog.Destroy()
				v.smpError(err)
			}
			pinLabel.SetText(pin)
		},
		"close_share_pin": func() {
			if v.peerRequestsSMP != nil {
				v.showSMPHasAlreadyStarted(v.peer.NameForPresentation())
				return
			}
			v.showWaitingForPeerToCompleteSMPDialog()
			v.session.StartSMP(v.peer.Jid, v.currentResource, i18n.Local("Por favor digite o código que o teu contato te enviou"), pin)
			v.newPinDialog.Destroy()
			v.newPinDialog = nil
		},
	})
	pinLabel.SetText(pin)
	v.newPinDialog.SetTransientFor(v.parent)
	v.newPinDialog.ShowAll()
}

func createPIN() (string, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(1000000)))
	if err != nil {
		log.Printf("Error encountered when creating a new PIN: %v", err)
		return "", err
	}
	return fmt.Sprintf("%06d", val), nil
}

func (v *verifier) showSMPHasAlreadyStarted(peerName string) {
	b := newBuilder("SMPHasAlreadyStarted")
	d := b.getObj("dialog").(gtki.Dialog)
	msg := b.getObj("smp_has_already_started").(gtki.Label)
	msg.SetText(i18n.Localf("%s já começou a verificação e gerou um código.\nPor favor pergunte que código foi gerado.", peerName))
	button := b.getObj("button_ok").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(d.Destroy)
	})
	d.SetTransientFor(v.newPinDialog)
	d.ShowAll()
	v.newPinDialog.Destroy()
	v.newPinDialog = nil
}

func (v *verifier) showWaitingForPeerToCompleteSMPDialog() {
	v.state = waitingForAnswerFromPeer
	v.disableNotifications()
	builderWaitingSMP := newBuilder("WaitingSMPComplete")
	waitingInfoBar := builderWaitingSMP.getObj("smp_waiting_infobar").(gtki.InfoBar)
	waitingSMPMessage := builderWaitingSMP.getObj("message").(gtki.Label)
	waitingSMPMessage.SetText(i18n.Localf("Esperando %s terminar a verificação do canal...", v.peer.NameForPresentation()))
	b := builderWaitingSMP.getObj("cancel_button").(gtki.Button)
	b.Connect("clicked", func() {
		//v.session.AbortSMP() TODO: allow this
		v.state = cancelingSMP
		v.disableNotifications()
		v.verificationWarning.Show()
	})
	waitingInfoBar.ShowAll()
	v.waitingForSMP = waitingInfoBar
	v.notifier.notify(waitingInfoBar)
}

func (v *verifier) showNotificationWhenWeCannotGeneratePINForSMP(err error) {
	log.Printf("Cannot recover from error: %v. Quitting verification using SMP.", err)
	errBuilder := newBuilder("CannotVerifyWithSMP")
	errInfoBar := errBuilder.getObj("error_verifying_smp").(gtki.InfoBar)
	message := errBuilder.getObj("message").(gtki.Label)
	message.SetText(i18n.Local("Não é possível verificar o canal agora."))
	button := errBuilder.getObj("try_later_button").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(errInfoBar.Destroy)
	})
	errInfoBar.ShowAll()
	v.notifier.notify(errInfoBar)
}

func (v *verifier) buildEnterPinDialog() {
	builder := newBuilder("EnterPIN")
	v.enterPinDialog = builder.getObj("dialog").(gtki.Dialog)
	v.enterPinDialog.SetTransientFor(v.parent)
	v.enterPinDialog.HideOnDelete()
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Localf("Digite o código que %s enviou para você", v.peer.NameForPresentation()))
	button := builder.getObj("button_submit").(gtki.Button)
	button.SetSensitive(false)
	builder.ConnectSignals(map[string]interface{}{
		"text_changing": func() {
			e := builder.getObj("pin").(gtki.Entry)
			pin, _ := e.GetText()
			button.SetSensitive(len(pin) > 0)
		},
		"close_share_pin": func() {
			e := builder.getObj("pin").(gtki.Entry)
			pin, _ := e.GetText()
			v.state = finishedAnsweringPeer
			v.disableNotifications()
			v.showWaitingForPeerToCompleteSMPDialog()
			v.session.FinishSMP(v.peer.Jid, v.currentResource, pin)
			v.enterPinDialog.Destroy()
			v.enterPinDialog = nil
		},
	})
	v.enterPinDialog.ShowAll()
}

func (v *verifier) displayRequestForSecret() {
	v.disableNotifications()
	b := newBuilder("PeerRequestsSMP")
	infobar := b.getObj("peer_requests_smp").(gtki.InfoBar)
	infobarMsg := b.getObj("message").(gtki.Label)
	verificationButton := b.getObj("verification_button").(gtki.Button)
	verificationButton.Connect("clicked", func() {
		if v.enterPinDialog != nil {
			v.enterPinDialog.Present()
		} else {
			v.buildEnterPinDialog()
		}
	})
	cancelButton := b.getObj("cancel_button").(gtki.Button)
	cancelButton.Connect("clicked", func() {
		v.peerRequestsSMP.Destroy()
		v.peerRequestsSMP = nil
		v.verificationWarning.Show()
	})
	infobarMsg.SetText(i18n.Localf("%s está esperando você terminar de verificar a segurança desse canal...", v.peer.NameForPresentation()))
	infobar.ShowAll()
	v.peerRequestsSMP = infobar
	v.notifier.notify(infobar)
}

func (v *verifier) displayVerificationSuccess() {
	builder := newBuilder("VerificationSucceeded")
	d := builder.getObj("dialog").(gtki.Dialog)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Localf("Yeeeee! Ninguém está espiando suas conversas com %s", v.peer.NameForPresentation()))
	button := builder.getObj("button_ok").(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(d.Destroy)
	})
	d.SetTransientFor(v.parent)
	d.ShowAll()
	v.disableNotifications()
}

func (v *verifier) displayVerificationFailure() {
	builder := newBuilder("VerificationFailed")
	d := builder.getObj("dialog").(gtki.Dialog)
	h := builder.getObj("header").(gtki.Label)
	sc, _ := h.GetStyleContext()
	sc.AddProvider(v.displayController.provider, 9999)
	sc.AddClass("verification-failed-header")
	v.displayController.provider.LoadFromData(
		`.verification-failed-header{
			font-weight: bold;
			font-size: 30px;
		}`)
	msg := builder.getObj("verification_message").(gtki.Label)
	msg.SetText(i18n.Localf("A verificação desse canal com %s falhou.", v.peer.NameForPresentation()))
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
	case finishedAnsweringPeer, cancelingSMP:
		v.removeInProgressNotifications()
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
