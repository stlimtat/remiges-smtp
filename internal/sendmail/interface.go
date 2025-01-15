package sendmail

import (
	"context"
	"net"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
)

//go:generate mockgen -destination=sendmail_mock.go -package=sendmail . IResolver,IMailSender
type IMailSender interface {
	LookupMX(ctx context.Context, domain dns.Domain) ([]string, error)
	NewConn(ctx context.Context, hosts []string) (net.Conn, error)
	SendMail(ctx context.Context, conn net.Conn, from smtp.Address, to smtp.Address, msg []byte) error
}
