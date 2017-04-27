package gui

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

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
