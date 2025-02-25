package sendmail

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net"

	"github.com/mjl-/adns"
	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/utils"
)

const (
	DEFAULT_SMTP_PORT_STR string = "25"
	TCP_NETWORK           string = "tcp"
)

type MXRecord struct {
	ADNSResult adns.Result
	Domain     string
	Entries    []dns.IPDomain
	Hosts      []string
}

type MailSender struct {
	CachedMX      map[string]MXRecord
	Debug         bool
	DialerFactory INetDialerFactory
	Resolver      dns.Resolver
	Slogger       *slog.Logger
	SmtpOpts      smtpclient.Opts
}

func NewMailSender(
	ctx context.Context,
	debug bool,
	dialerFactory INetDialerFactory,
	resolver dns.Resolver,
	slogger *slog.Logger,
) *MailSender {
	result := &MailSender{
		CachedMX:      make(map[string]MXRecord, 0),
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

func (m *MailSender) LookupMX(
	ctx context.Context,
	domain dns.Domain,
) ([]string, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("domain", domain.ASCII).
		Logger()

	// 1. check if mx already exists in cache
	result, ok := m.CachedMX[domain.ASCII]
	if !ok {
		// 2. resolve the mx record for the domain
		ipDomain := dns.IPDomain{
			Domain: domain,
		}

		_, _, _, expandedNextHop, hosts, _, err := smtpclient.GatherDestinations(
			ctx, m.Slogger, m.Resolver, ipDomain,
		)
		if err != nil {
			logger.Error().Err(err).Msg("smtpclient.GatherDestinations")
			return nil, err
		}
		// 3. convert from dns.IPDomain to string
		hostStrSlice := []string{}
		for _, host := range hosts {
			hostStrSlice = append(hostStrSlice, host.String())
		}
		if expandedNextHop.ASCII != domain.ASCII {
			m.CachedMX[expandedNextHop.ASCII] = MXRecord{
				Domain:  expandedNextHop.ASCII,
				Entries: hosts,
				Hosts:   hostStrSlice,
			}
		}
		m.CachedMX[domain.ASCII] = MXRecord{
			Domain:  domain.ASCII,
			Entries: hosts,
			Hosts:   hostStrSlice,
		}
		result = m.CachedMX[domain.ASCII]
	}
	logger.Info().
		Strs("hosts", result.Hosts).
		Msg("lookupMX")
	return result.Hosts, nil
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
	myMail *mail.Mail,
) (results map[string][]Response, errs map[string]error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("from", myMail.From.String()).
		Bytes("msgid", myMail.MsgID).
		Str("subject", string(myMail.Subject)).
		Logger()
	results = make(map[string][]Response, 0)
	errs = make(map[string]error, 0)

	for _, to := range myMail.To {
		toStr := to.String()
		hosts, err := m.LookupMX(ctx, to.Domain)
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
	myMail *mail.Mail,
	to smtp.Address,
) ([]Response, error) {
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
		int64(len(myMail.Body)),
		bytes.NewReader(myMail.Body),
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
	results := make([]Response, 0)
	for _, resp := range resps {
		logger.Info().
			Interface("resp", resp).
			Msg("smtpclient.Deliver response")
		resp := Response{
			Response: resp,
		}
		results = append(results, resp)
	}
	return results, nil
}
