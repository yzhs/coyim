package gui

import (
	"github.com/twstrike/coyim/i18n"
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

func buildVerifyIdentityNotification(acc *account, peer, resource string, win gtki.Window) gtki.InfoBar {
	builder := newBuilder("VerifyIdentityNotification")

	obj := builder.getObj("infobar")
	infoBar := obj.(gtki.InfoBar)

	obj = builder.getObj("message")
	message := obj.(gtki.Label)
	message.SetSelectable(true)

	text := i18n.Local("This conversation may not be secure.")
	message.SetText(text)

	obj = builder.getObj("button_verify")
	button := obj.(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			secureChannel(acc, peer, resource, win)
			// Handle hiding or keeping the infobar here
			//ok := verifyFingerprintDialog(acc, peer, resource, win)
			//if ok != gtki.RESPONSE_YES {
			//	infoBar.Hide()
			//	infoBar.Destroy()
			//}
		})
	})

	infoBar.ShowAll()

	return infoBar
}

func secureChannel(account *account, uid, resource string, parent gtki.Window) {
	builder := newBuilder("ChooseVerificationType")
	d := builder.getObj("dialog").(gtki.Dialog)
	cancelButton := builder.getObj("cancel_button").(gtki.Button)
	cancelButton.Connect("clicked", func() {
		d.Destroy()
	})
	validateButton := builder.getObj("validate_button").(gtki.Button)
	validateButton.Connect("clicked", func() {
		doInUIThread(func() {
			smpValidationDialog(parent)
			d.Destroy()
		})
	})
	d.SetTransientFor(parent)
	d.ShowAll()
}

func smpValidationDialog(parent gtki.Window) {
	builder := newBuilder("ValidateSecureChannel")
	d := builder.getObj("dialog").(gtki.Dialog)
	submit := builder.getObj("button_submit").(gtki.Button)
	submit.Connect("clicked", func() {
		doInUIThread(func() {
			e := builder.getObj("pin").(gtki.Entry)
			e.GetText()
			// SUBMIT PIN TO SMP BACKEND HERE
			// check if success or failure
			// if success
			//     showSecureChannelCreated(d)
			// else
			showPINWasIncorrect(parent)
			d.Destroy()
		})
	})
	d.SetTransientFor(parent)
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
			parent.Destroy()
		})
	})
	d.SetTransientFor(parent)
	d.ShowAll()
}

func showSecureChannelCreated(parent gtki.Window) {
	builder := newBuilder("SecureChannelCreated")
	d := builder.getObj("dialog").(gtki.Dialog)
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
