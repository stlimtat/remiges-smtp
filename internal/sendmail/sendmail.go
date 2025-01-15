package sendmail

import (
	"bytes"
	"context"
	"log/slog"
	"net"

	"github.com/mjl-/adns"
	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
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
	DialerFactory INetDialerFactory
	Resolver      dns.Resolver
	Slogger       *slog.Logger
	SmtpOpts      smtpclient.Opts
}

func NewMailSender(
	ctx context.Context,
	dialerFactory INetDialerFactory,
	resolver dns.Resolver,
	slogger *slog.Logger,
) *MailSender {
	result := &MailSender{
		CachedMX:      make(map[string]MXRecord, 0),
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
	mxRecord, ok := m.CachedMX[domain.ASCII]
	if !ok {
		// 2. resolve the mx record for the domain
		ipDomain := dns.IPDomain{
			Domain: domain,
		}

		_, _, _, expandedNextHop, hosts, _, err := smtpclient.GatherDestinations(ctx, m.Slogger, m.Resolver, ipDomain)
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
	}
	logger.Info().
		Strs("hosts", mxRecord.Hosts).
		Msg("lookupMX")
	return mxRecord.Hosts, nil
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
	result, err := dialer.Dial(TCP_NETWORK, addr)
	if err != nil {
		logger.Error().Err(err).Msg("d.Dial")
		return nil, err
	}
	return result, nil
}

func (m *MailSender) SendMail(
	ctx context.Context,
	conn net.Conn,
	from smtp.Address,
	to smtp.Address,
	msg []byte,
) error {
	logger := zerolog.Ctx(ctx).
		With().
		Str("from", from.String()).
		Str("to", to.String()).
		Bytes("msg", msg).
		Logger()

	client, err := smtpclient.New(
		ctx,
		m.Slogger,
		conn,
		smtpclient.TLSOpportunistic,
		false,
		to.Domain,
		to.Domain,
		m.SmtpOpts,
	)
	if err != nil {
		logger.Error().Err(err).Msg("smtpclient.New")
		return err
	}
	err = client.Deliver(
		ctx,
		from.String(),
		to.String(),
		int64(len(msg)),
		bytes.NewReader(msg),
		true, false, false,
	)
	if err != nil {
		if smtpclientErr, ok := err.(smtpclient.Error); ok {
			logger.Error().
				Err(err).
				Interface("smtpclient_err", smtpclientErr).
				Msg("smtpclient.Deliver")
		} else {
			logger.Error().Err(err).Msg("smtpclient.Deliver")
		}
	}
	return err
}
