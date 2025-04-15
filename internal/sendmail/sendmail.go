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

// MailSender handles the delivery of emails to SMTP servers.
// It manages connections, retries, and concurrent delivery to multiple recipients.
type MailSender struct {
	// CachedMX stores MX records for domains to reduce DNS lookups
	CachedMX map[string]dn.MXRecord

	// Debug enables debug mode which prevents actual mail sending
	Debug bool

	// DialerFactory creates network dialers for SMTP connections
	DialerFactory INetDialerFactory

	// Resolver handles DNS lookups for MX records
	Resolver dns.IResolver

	// Slogger is used for structured logging
	Slogger *slog.Logger

	// SmtpOpts contains SMTP client configuration options
	SmtpOpts smtpclient.Opts

	// maxRetries is the maximum number of delivery attempts per recipient
	maxRetries int

	// retryDelay is the base delay between retry attempts
	retryDelay time.Duration

	// metrics tracks various delivery statistics
	metrics *Metrics
}

// deliveryResult represents the outcome of a mail delivery attempt
type deliveryResult struct {
	responses []pmail.Response // SMTP server responses
	err       error            // Any error that occurred during delivery
}

// Metrics collects various statistics about mail delivery
type Metrics struct {
	deliveryAttempts  prometheus.Counter   // Total number of delivery attempts
	deliverySuccesses prometheus.Counter   // Number of successful deliveries
	deliveryFailures  prometheus.Counter   // Number of failed deliveries
	deliveryDuration  prometheus.Histogram // Distribution of delivery times
	connectionErrors  prometheus.Counter   // Number of connection failures
}

// recordMetrics updates delivery metrics after each attempt
func (m *MailSender) recordMetrics(start time.Time, err error) {
	m.metrics.deliveryAttempts.Inc()
	if err != nil {
		m.metrics.deliveryFailures.Inc()
	} else {
		m.metrics.deliverySuccesses.Inc()
	}
	m.metrics.deliveryDuration.Observe(time.Since(start).Seconds())
}

// NewMailSender creates a new MailSender with the specified configuration.
//
// Parameters:
//   - ctx: Context for the sender creation
//   - debug: Enable debug mode (prevents actual mail sending)
//   - dialerFactory: Factory for creating network dialers
//   - resolver: DNS resolver for MX lookups
//   - slogger: Structured logger
//
// Returns:
//   - *MailSender: A new mail sender instance
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

// NewConn establishes a new connection to one of the provided SMTP hosts.
// It randomly selects a host from the list to enable load balancing.
//
// Parameters:
//   - ctx: Context for the connection operation
//   - hosts: List of SMTP server hostnames
//
// Returns:
//   - net.Conn: Established network connection
//   - error: Any error encountered during connection
func (m *MailSender) NewConn(
	ctx context.Context,
	hosts []string,
) (net.Conn, error) {
	logger := zerolog.Ctx(ctx).With().Strs("hosts", hosts).Logger()
	var err error

	if m.Debug {
		logger.Debug().Msg("debug mode, not connecting to SMTP server")
		return nil, rerrors.NewError(rerrors.ErrSMTPConnection, "debug mode, not connecting to SMTP server", nil)
	}

	// Randomly select a host for load balancing
	randomInt, err := utils.RandInt(int64(len(hosts)))
	if err != nil {
		logger.Error().Err(err).Msg("utils.RandInt")
		return nil, err
	}
	host := hosts[randomInt]

	// Create a new dialer and establish connection
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

// SendMail attempts to deliver an email to all recipients concurrently.
// It manages the delivery process, including retries and error handling.
//
// Parameters:
//   - ctx: Context for the sending operation
//   - mail: Email to be sent
//
// Returns:
//   - map[string][]pmail.Response: Map of recipient addresses to their SMTP responses
//   - map[string]error: Map of recipient addresses to any errors encountered
func (m *MailSender) SendMail(
	ctx context.Context,
	mail *pmail.Mail,
) (map[string][]pmail.Response, map[string]error) {
	logger := zerolog.Ctx(ctx)

	// Validate the email before attempting delivery
	if err := mail.Validate(); err != nil {
		return nil, map[string]error{
			mail.To[0].String(): rerrors.NewError(rerrors.ErrMailValidation, "invalid mail", err),
		}
	}

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

	// Collect results from all delivery attempts
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

	// Return combined error if any deliveries failed
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

// deliverToRecipient handles delivery to a single recipient with retries.
// It implements exponential backoff for retries and handles various error conditions.
//
// Parameters:
//   - ctx: Context for the delivery operation
//   - mail: Email to be delivered
//   - to: Recipient's SMTP address
//
// Returns:
//   - deliveryResult: The result of the delivery attempt including responses and errors
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

		// Lookup MX records for the recipient's domain
		hosts, err := m.Resolver.LookupMX(ctx, to.Domain)
		if err != nil {
			lastErr = rerrors.NewError(rerrors.ErrDNSLookup, "failed to lookup MX records", err).
				WithContext("domain", to.Domain)
			continue
		}

		// Attempt to establish connection and deliver
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

// Deliver sends an email to a specific recipient through an established connection.
// It handles the SMTP protocol interaction and collects server responses.
//
// Parameters:
//   - ctx: Context for the delivery operation
//   - conn: Established network connection
//   - myMail: Email to be delivered
//   - to: Recipient's SMTP address
//
// Returns:
//   - []pmail.Response: Array of SMTP server responses
//   - error: Any error encountered during delivery
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

	// Skip actual delivery in debug mode
	if m.Debug {
		logger.Info().Msg("debug mode, not sending mail")
		return nil, nil
	}

	// Create SMTP client with TLS support
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

	// Deliver the email and collect responses
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

	// Ensure we received responses
	if len(resps) == 0 {
		return nil, rerrors.NewError(rerrors.ErrMailDelivery, "no responses", nil)
	}

	// Convert and collect responses
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
