package sendmail

import (
	"context"
	"net"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/internal/mail"
)

//go:generate mockgen -destination=smtpclient_mock.go -package=sendmail github.com/mjl-/mox/smtpclient Dialer
//go:generate mockgen -destination=dns_mock.go -package=sendmail github.com/mjl-/mox/dns Resolver
//go:generate mockgen -destination=mock.go -package=sendmail -source=interface.go

type INetDialerFactory interface {
	NewDialer() smtpclient.Dialer
}

type IMailSender interface {
	Deliver(ctx context.Context, conn net.Conn, mail *mail.Mail) error
	LookupMX(ctx context.Context, domain dns.Domain) ([]string, error)
	NewConn(ctx context.Context, hosts []string) (net.Conn, error)
	SendMail(ctx context.Context, mail *mail.Mail) error
}
