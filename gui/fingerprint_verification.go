package gui

import (
	"fmt"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3/gtk"
)

func verifyFingerprintDialog(account *Account, uid string) {
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle(i18n.Local("Fingerprint verification"))
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox, _ := dialog.GetContentArea()

	//TODO: errors
	conversation := account.GetConversationWith(uid)
	fpr := conversation.GetTheirKey().DefaultFingerprint()

	// message copied from libpurple
	message := fmt.Sprintf(i18n.Local(`
	Fingerprint for you (%[1]s): %[2]x

	Purported fingerprint for %[3]s: %[4]x

	Is this the verifiably correct fingerprint for %[3]s?
	`), account.Account, account.Session.PrivateKey.DefaultFingerprint(), uid, fpr)
	l, _ := gtk.LabelNew(message)
	vbox.Add(l)

	button, _ := gtk.ButtonNewWithLabel(i18n.Local("Verify"))
	vbox.Add(button)

	button.Connect("clicked", func() {
		defer dialog.Destroy()

		err := account.AuthorizeFingerprint(uid, fpr)
		if err != nil {
			//TODO: Error
			return
		}

		//TODO: error
		account.configManager.Save()
	})

	dialog.ShowAll()
}