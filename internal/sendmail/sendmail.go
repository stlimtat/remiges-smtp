package sendmail

import (
	"bytes"
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/dns"
	rerrors "github.com/stlimtat/remiges-smtp/internal/errors"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/dn"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	// DefaultSMTPPort is the standard SMTP port
	DefaultSMTPPort = "25"

	// TCPNetwork specifies the network type for SMTP connections
	TCPNetwork = "tcp"

	// Maximum number of delivery attempts
	maxDeliveryAttempts = 3

	// Initial retry delay
	baseRetryDelay = 5 * time.Second
)

// MailSender handles the delivery of emails to SMTP servers
type MailSender struct {
	CachedMX      map[string]dn.MXRecord
	Debug         bool
	DialerFactory INetDialerFactory
	Resolver      dns.IResolver
	Slogger       *slog.Logger
	SmtpOpts      smtpclient.Opts
	maxRetries    int
	retryDelay    time.Duration
	metrics       *Metrics
}

type smtpConnection struct {
	conn     net.Conn
	client   *smtpclient.Client
	lastUsed time.Time
}

// deliveryResult represents the outcome of a mail delivery attempt
type deliveryResult struct {
	responses []pmail.Response
	err       error
}

// Add metrics collection
type Metrics struct {
	deliveryAttempts  prometheus.Counter
	deliverySuccesses prometheus.Counter
	deliveryFailures  prometheus.Counter
	deliveryDuration  prometheus.Histogram
	connectionErrors  prometheus.Counter
}

func (m *MailSender) recordMetrics(start time.Time, err error) {
	m.metrics.deliveryAttempts.Inc()
	if err != nil {
		m.metrics.deliveryFailures.Inc()
	} else {
		m.metrics.deliverySuccesses.Inc()
	}
	m.metrics.deliveryDuration.Observe(time.Since(start).Seconds())
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
		maxRetries: 3,
		retryDelay: 5 * time.Second,
		metrics:    &Metrics{},
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

	dialer, err := m.DialerFactory.NewDialer(ctx)
	if err != nil {
		return nil, err
	}
	addr := net.JoinHostPort(host, DefaultSMTPPort)
	result, err := dialer.DialContext(ctx, TCPNetwork, addr)
	if err != nil {
		logger.Error().Err(err).Msg("d.Dial")
		return nil, err
	}
	return result, nil
}

// SendMail attempts to deliver an email to all recipients
func (m *MailSender) SendMail(
	ctx context.Context,
	mail *pmail.Mail,
) (map[string][]pmail.Response, map[string]error) {
	if err := mail.Validate(); err != nil {
		return nil, map[string]error{
			mail.To[0].String(): rerrors.NewError(rerrors.ErrMailValidation, "invalid mail", err),
		}
	}

	logger := zerolog.Ctx(ctx)
	results := make(map[string][]pmail.Response)
	errs := make(map[string]error)

	// Create channels for concurrent delivery
	resultChan := make(chan struct {
		addr   string
		result deliveryResult
	}, len(mail.To))

	// Process each recipient concurrently
	for _, to := range mail.To {
		go func(addr smtp.Address) {
			result := m.deliverToRecipient(ctx, mail, addr)
			resultChan <- struct {
				addr   string
				result deliveryResult
			}{addr.String(), result}
		}(to)
	}

	// Collect results
	for i := 0; i < len(mail.To); i++ {
		result := <-resultChan
		if result.result.err != nil {
			errs[result.addr] = result.result.err
			logger.Error().
				Err(result.result.err).
				Str("recipient", result.addr).
				Msg("Failed to deliver mail")
		} else {
			results[result.addr] = result.result.responses
			logger.Info().
				Str("recipient", result.addr).
				Msg("Successfully delivered mail")
		}
	}

	if len(errs) > 0 {
		return results, map[string]error{
			mail.To[0].String(): rerrors.NewError(
				rerrors.ErrMailDelivery,
				"delivery failed for some recipients",
				nil,
			).
				WithContext("errors", errs),
		}
	}

	return results, nil
}

// deliverToRecipient handles delivery to a single recipient with retries
func (m *MailSender) deliverToRecipient(
	ctx context.Context,
	mail *pmail.Mail,
	to smtp.Address,
) deliveryResult {
	var lastErr error

	for attempt := 0; attempt < m.maxRetries; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return deliveryResult{nil, ctx.Err()}
		default:
		}

		// Lookup MX records
		hosts, err := m.Resolver.LookupMX(ctx, to.Domain)
		if err != nil {
			lastErr = rerrors.NewError(rerrors.ErrDNSLookup, "failed to lookup MX records", err).
				WithContext("domain", to.Domain)
			continue
		}

		// Attempt delivery
		conn, err := m.NewConn(ctx, hosts)
		if err != nil {
			lastErr = rerrors.NewError(rerrors.ErrSMTPConnection, "failed to establish connection", err).
				WithContext("hosts", hosts)
			continue
		}

		responses, err := m.Deliver(ctx, conn, mail, to)
		if err != nil {
			lastErr = err
			continue
		}

		return deliveryResult{responses, nil}
	}

	return deliveryResult{nil, rerrors.NewError(rerrors.ErrMailDelivery, "max retries exceeded", lastErr)}
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
		return nil, rerrors.NewError(rerrors.ErrMailDelivery, "no responses", nil)
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

// Add structured logging
func (m *MailSender) logDeliveryAttempt(
	ctx context.Context,
	mail *pmail.Mail,
	to smtp.Address,
	err error,
) {
	logger := zerolog.Ctx(ctx)
	event := logger.Info()
	if err != nil {
		event = logger.Error().Err(err)
	}

	event.
		Str("message_id", string(mail.MsgID)).
		Str("from", mail.From.String()).
		Str("to", to.String()).
		Str("subject", string(mail.Subject)).
		Int64("size", int64(len(mail.FinalBody))).
		Timestamp().
		Msg("Mail delivery attempt")
}
