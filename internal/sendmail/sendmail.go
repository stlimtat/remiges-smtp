package sendmail

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net"

	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/dns"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/dn"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	DEFAULT_SMTP_PORT_STR string = "25"
	TCP_NETWORK           string = "tcp"
)

type MailSender struct {
	CachedMX      map[string]dn.MXRecord
	Debug         bool
	DialerFactory INetDialerFactory
	Resolver      dns.IResolver
	Slogger       *slog.Logger
	SmtpOpts      smtpclient.Opts
}

func NewMailSender(
	ctx context.Context,
	debug bool,
	dialerFactory INetDialerFactory,
	resolver dns.IResolver,
	slogger *slog.Logger,
) *MailSender {
	result := &MailSender{
		CachedMX:      make(map[string]dn.MXRecord, 0),
		Debug:         debug,
		DialerFactory: dialerFactory,
		Resolver:      resolver,
		Slogger:       slogger,
		SmtpOpts: smtpclient.Opts{
			// Auth is nil, because we don't need authentication for the smtp server/relay
			Auth:    nil,
			RootCAs: config.GetCertPool(ctx),
		},
	}
	return result
}

func (m *MailSender) NewConn(
	ctx context.Context,
	hosts []string,
) (net.Conn, error) {
	logger := zerolog.Ctx(ctx).With().Strs("hosts", hosts).Logger()
	var err error
	randomInt, err := utils.RandInt(int64(len(hosts)))
	if err != nil {
		logger.Error().Err(err).Msg("utils.RandInt")
		return nil, err
	}
	host := hosts[randomInt]

	dialer := m.DialerFactory.NewDialer()
	addr := net.JoinHostPort(host, DEFAULT_SMTP_PORT_STR)
	result, err := dialer.DialContext(ctx, TCP_NETWORK, addr)
	if err != nil {
		logger.Error().Err(err).Msg("d.Dial")
		return nil, err
	}
	return result, nil
}

func (m *MailSender) SendMail(
	ctx context.Context,
	myMail *pmail.Mail,
) (results map[string][]pmail.Response, errs map[string]error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("from", myMail.From.String()).
		Bytes("msgid", myMail.MsgID).
		Str("subject", string(myMail.Subject)).
		Logger()
	results = make(map[string][]pmail.Response, 0)
	errs = make(map[string]error, 0)

	for _, to := range myMail.To {
		toStr := to.String()
		hosts, err := m.Resolver.LookupMX(ctx, to.Domain)
		if err != nil {
			errs[toStr] = err
			continue
		}
		conn, err := m.NewConn(ctx, hosts)
		if err != nil {
			errs[toStr] = err
			continue
		}
		result, err := m.Deliver(ctx, conn, myMail, to)
		if err != nil {
			errs[toStr] = err
			continue
		}
		results[toStr] = result
		logger.Info().
			AnErr("err", errs[toStr]).
			Interface("result", result).
			Str("to", toStr).
			Msg("Delivery attempted")
	}

	return results, nil
}

func (m *MailSender) Deliver(
	ctx context.Context,
	conn net.Conn,
	myMail *pmail.Mail,
	to smtp.Address,
) ([]pmail.Response, error) {
	toStr := to.String()
	logger := zerolog.Ctx(ctx).
		With().
		Str("from", myMail.From.String()).
		Str("to", toStr).
		Bytes("msgid", myMail.MsgID).
		Bytes("subject", myMail.Subject).
		Bytes("content_type", myMail.ContentType).
		Logger()
	if m.Debug {
		logger.Info().Msg("debug mode, not sending mail")
		return nil, nil
	}

	client, err := smtpclient.New(
		ctx,
		m.Slogger,
		conn,
		smtpclient.TLSOpportunistic,
		false,
		myMail.From.Domain,
		to.Domain,
		m.SmtpOpts,
	)
	if err != nil {
		logger.Error().Err(err).
			Msg("smtpclient.New")
		return nil, err
	}
	resps, err := client.DeliverMultiple(
		ctx,
		myMail.From.String(),
		[]string{toStr},
		int64(len(myMail.FinalBody)),
		bytes.NewReader(myMail.FinalBody),
		true, false, false,
	)
	if err != nil {
		if smtpclientErr, ok := err.(smtpclient.Error); ok {
			logger.Error().
				Err(err).
				Interface("smtpclient_err", smtpclientErr).
				Msg("smtpclient.Deliver")
		} else {
			logger.Error().Err(err).
				Msg("smtpclient.Deliver")
		}
		return nil, err
	}

	if len(resps) == 0 {
		return nil, errors.New("no responses")
	}
	results := make([]pmail.Response, 0)
	for _, resp := range resps {
		logger.Info().
			Interface("resp", resp).
			Msg("smtpclient.Deliver response")
		resp := pmail.Response{
			Response: resp,
		}
		results = append(results, resp)
	}
	return results, nil
}
