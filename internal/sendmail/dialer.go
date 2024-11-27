package sendmail

import (
	"net"
	"time"
)

//go:generate mockgen -destination=dialer_mock.go -package=sendmail . INetDialer,INetDialerFactory
type INetDialer interface {
	Dial(network, address string) (net.Conn, error)
}

type INetDialerFactory interface {
	NewDialer() INetDialer
}

type DefaultNetDialerFactory struct{}

func NewDefaultDialerFactory() INetDialerFactory {
	result := &DefaultNetDialerFactory{}
	return result
}

func (_ *DefaultNetDialerFactory) NewDialer() INetDialer {
	return &net.Dialer{Timeout: 50 * time.Second}
}
