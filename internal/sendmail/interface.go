package sendmail

import (
	"context"
	"net"

	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
)

//go:generate mockgen -destination=mox_mock.go -package=sendmail github.com/mjl-/mox/smtpclient Dialer
//go:generate mockgen -destination=mock.go -package=sendmail . INetDialerFactory,IMailSender

type INetDialerFactory interface {
	NewDialer() smtpclient.Dialer
}

type IMailSender interface {
	Deliver(ctx context.Context, conn net.Conn, mail *mail.Mail, to smtp.Address) ([]mail.Response, error)
	NewConn(ctx context.Context, hosts []string) (net.Conn, error)
	SendMail(ctx context.Context, mail *mail.Mail) (map[string][]mail.Response, map[string]error)
}
