package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/twstrike/coyim/xmpp"
	"golang.org/x/net/proxy"
)

func init() {
	proxy.RegisterDialerType("socks5+unix", func(u *url.URL, d proxy.Dialer) (proxy.Dialer, error) {
		var auth *proxy.Auth
		if u.User != nil {
			auth = &proxy.Auth{
				User: u.User.Username(),
			}

			if p, ok := u.User.Password(); ok {
				auth.Password = p
			}
		}

		return proxy.SOCKS5("unix", u.Path, auth, d)
	})
}

func resolveXMPPServerOverTor(domain string) ([]string, error) {
	u, err := url.Parse(newTorProxy(detectTor()))
	if err != nil {
		return nil, errors.New("Failed to resolve XMPP server: " + err.Error())
	}

	dnsProxy, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		return nil, errors.New("Failed to resolve XMPP server: " + err.Error())
	}

	ret, err := xmpp.ResolveProxy(dnsProxy, domain)
	if err != nil {
		return nil, errors.New("Failed to resolve XMPP server: " + err.Error())
	}

	return ret, nil
}

func buildProxyChain(proxies []string) (dialer proxy.Dialer, err error) {
	for i := len(proxies) - 1; i >= 0; i-- {
		u, e := url.Parse(proxies[i])
		if e != nil {
			err = errors.New("Failed to parse " + proxies[i] + " as a URL: " + e.Error())
			return
		}

		if dialer == nil {
			dialer = proxy.Direct
		}

		if dialer, err = proxy.FromURL(u, dialer); err != nil {
			err = errors.New("Failed to parse " + proxies[i] + " as a proxy: " + err.Error())
			return
		}
	}

	return
}

// NewXMPPConn creates a new XMPP connection based on the given information
func NewXMPPConn(config *Config, password string, createCallback xmpp.FormCallback, logger io.Writer) (*xmpp.Conn, error) {
	parts := strings.SplitN(config.Account, "@", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid username (want user@domain): " + config.Account)
	}

	user := parts[0]
	domain := parts[1]

	var xmppAddrs []string
	addrTrusted := false

	if len(config.Server) > 0 && config.Port > 0 {
		xmppAddrs = []string{fmt.Sprintf("%s:%d", config.Server, config.Port)}
		addrTrusted = true
	} else {
		if len(config.Proxies) > 0 && len(detectTor()) == 0 {
			return nil, errors.New("Cannot connect via a proxy without Server and Port being set in the config file as an SRV lookup would leak information.")
		}

		var err error
		if xmppAddrs, err = resolveXMPPServerOverTor(domain); err != nil {
			return nil, err
		}

		// Fallback to using the domain at default port
		if len(xmppAddrs) == 0 {
			xmppAddrs = []string{domain + "5222"}
		}
	}

	dialer, err := buildProxyChain(config.Proxies)
	if err != nil {
		return nil, err
	}

	var certSHA256 []byte
	if len(config.ServerCertificateSHA256) > 0 {
		certSHA256, err = hex.DecodeString(config.ServerCertificateSHA256)
		if err != nil {
			return nil, errors.New("Failed to parse ServerCertificateSHA256 (should be hex string): " + err.Error())
		}

		if len(certSHA256) != 32 {
			return nil, errors.New("ServerCertificateSHA256 is not 32 bytes long")
		}
	}

	xmppConfig := &xmpp.Config{
		Log:                     logger,
		CreateCallback:          createCallback,
		TrustedAddress:          addrTrusted,
		Archive:                 false,
		ServerCertificateSHA256: certSHA256,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS10,
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
		},
	}

	if domain == "jabber.ccc.de" {
		// jabber.ccc.de uses CACert but distros are removing that root
		// certificate.
		roots := x509.NewCertPool()
		caCertRoot, err := x509.ParseCertificate(caCertRootDER)
		if err == nil {
			//TODO: UI should have a Alert() method
			//alert(term, "Temporarily trusting only CACert root for CCC Jabber server")
			roots.AddCert(caCertRoot)
			xmppConfig.TLSConfig.RootCAs = roots
		} else {
			//TODO
			//alert(term, "Tried to add CACert root for jabber.ccc.de but failed: "+err.Error())
		}
	}

	//TODO: It may be locking
	//Also, move this defered functions
	//if len(config.RawLogFile) > 0 {
	//	rawLog, err := os.OpenFile(config.RawLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	//	if err != nil {
	//		return nil, errors.New("Failed to open raw log file: " + err.Error())
	//	}

	//	lock := new(sync.Mutex)
	//	in := rawLogger{
	//		out:    rawLog,
	//		prefix: []byte("<- "),
	//		lock:   lock,
	//	}
	//	out := rawLogger{
	//		out:    rawLog,
	//		prefix: []byte("-> "),
	//		lock:   lock,
	//	}
	//	in.other, out.other = &out, &in

	//	xmppConfig.InLog = &in
	//	xmppConfig.OutLog = &out

	//	defer in.flush()
	//	defer out.flush()
	//}

	var addr string
	xmppConfig.Conn, addr, err = connectToFirstAvailable(xmppAddrs, dialer)
	if err != nil {
		return nil, err
	}

	return connect(addr, user, domain, password, xmppConfig)
}

func connectToFirstAvailable(xmppAddrs []string, dialer proxy.Dialer) (net.Conn, string, error) {
	if dialer == nil {
		dialer = proxy.Direct
	}

	for _, addr := range xmppAddrs {
		log.Println("Connecting to " + addr)

		conn, err := dialer.Dial("tcp", addr)
		if err == nil {
			return conn, addr, nil
		}

		log.Println("Failed to connect to", addr, "\n\t", err)
	}

	// should NOT attempt the fallback described in XMPP section 3.2.2
	return nil, "", errors.New("Failed to connect to XMPP server: exhausted list of XMPP SRV for server")
}

func connect(addr, user, domain, password string, xmppConfig *xmpp.Config) (*xmpp.Conn, error) {
	conn, err := xmpp.Dial(addr, user, domain, password, xmppConfig)
	if err != nil {
		return nil, errors.New("Failed to connect to XMPP server: " + err.Error())
	}

	return conn, nil
}