package sendmail

import (
	"net"
	"time"

	"github.com/mjl-/mox/smtpclient"
)

type DefaultNetDialerFactory struct{}

func NewDefaultDialerFactory() *DefaultNetDialerFactory {
	result := &DefaultNetDialerFactory{}
	return result
}

func (_ *DefaultNetDialerFactory) NewDialer() smtpclient.Dialer {
	return &net.Dialer{Timeout: 50 * time.Second}
}
