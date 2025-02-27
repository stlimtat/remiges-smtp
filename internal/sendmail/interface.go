package sendmail

import (
	"context"
	"net"

	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

//go:generate mockgen -destination=mox_mock.go -package=sendmail github.com/mjl-/mox/smtpclient Dialer
//go:generate mockgen -destination=mock.go -package=sendmail . INetDialerFactory,IMailSender

type INetDialerFactory interface {
	NewDialer() smtpclient.Dialer
}

type IMailSender interface {
	Deliver(ctx context.Context, conn net.Conn, myMail *pmail.Mail, to smtp.Address) ([]pmail.Response, error)
	NewConn(ctx context.Context, hosts []string) (net.Conn, error)
	SendMail(ctx context.Context, mail *pmail.Mail) (map[string][]pmail.Response, map[string]error)
}
