package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/xmpp"
)

func (u *gtkUI) connectAccount(account *account) {
	switch p := account.session.CurrentAccount.Password; p {
	case "":
		u.askForPasswordAndConnect(account)
	default:
		go u.connectWithPassword(account, p)
	}
}

func (u *gtkUI) connectWithPassword(account *account, password string) error {
	u.showConnectAccountNotification(account)
	defer u.removeConnectAccountNotification(account)

	err := account.session.Connect(password)
	switch err {
	case config.ErrTorNotRunning:
		glib.IdleAdd(u.alertTorIsNotRunning)
	case xmpp.ErrTCPBindingFailed:
		u.askForServerDetailsAndConnect(account, password)
	case xmpp.ErrAuthenticationFailed:
		//TODO: notify authentication failure?
		u.askForPasswordAndConnect(account)
	case xmpp.ErrConnectionFailed:
		//TODO: notify connection failure?
	}

	return err
}

func (u *gtkUI) askForPasswordAndConnect(account *account) {
	accountName := account.session.CurrentAccount.Account
	glib.IdleAdd(func() {
		u.askForPassword(accountName, func(password string) error {
			return u.connectWithPassword(account, password)
		})
	})
}

func (u *gtkUI) askForServerDetailsAndConnect(account *account, password string) {
	conf := account.session.CurrentAccount
	glib.IdleAdd(func() {
		u.askForServerDetails(conf, func() error {
			return u.connectWithPassword(account, password)
		})
	})
}
