package sendmail

import (
	"context"
	"net"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
)

//go:generate mockgen -destination=smtpclient_mock.go -package=sendmail github.com/mjl-/mox/smtpclient Dialer
//go:generate mockgen -destination=dns_mock.go -package=sendmail github.com/mjl-/mox/dns Resolver
//go:generate mockgen -destination=mock.go -package=sendmail -source=interface.go

type INetDialerFactory interface {
	NewDialer() smtpclient.Dialer
}

type IMailSender interface {
	LookupMX(ctx context.Context, domain dns.Domain) ([]string, error)
	NewConn(ctx context.Context, hosts []string) (net.Conn, error)
	SendMail(ctx context.Context, conn net.Conn, from smtp.Address, to smtp.Address, msg []byte) error
}
