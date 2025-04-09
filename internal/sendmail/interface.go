// Package sendmail provides functionality for sending emails via SMTP servers.
// It includes interfaces and implementations for dialing SMTP connections,
// handling email delivery, and managing the entire sending process.
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

// INetDialerFactory defines an interface for creating network dialers.
// It abstracts the creation of SMTP connection dialers, allowing for different
// dialing strategies (e.g., direct TCP, SOCKS5 proxy).
type INetDialerFactory interface {
	// NewDialer creates and returns a new SMTP client dialer.
	// The dialer is configured based on the factory's settings and can handle
	// both direct TCP connections and proxy connections (if configured).
	//
	// Parameters:
	//   - ctx: Context for the dialer creation operation
	//
	// Returns:
	//   - smtpclient.Dialer: A configured dialer for SMTP connections
	//   - error: Any error encountered during dialer creation
	NewDialer(ctx context.Context) (smtpclient.Dialer, error)
}

// IMailSender defines the interface for sending emails via SMTP.
// It provides methods for establishing connections to SMTP servers,
// delivering emails, and managing the entire sending process.
type IMailSender interface {
	// Deliver sends an email to a specific recipient through an established connection.
	// It handles the SMTP protocol interaction for a single recipient.
	//
	// Parameters:
	//   - ctx: Context for the delivery operation
	//   - conn: Established network connection to the SMTP server
	//   - myMail: Email to be delivered
	//   - to: Recipient's SMTP address
	//
	// Returns:
	//   - []pmail.Response: Array of SMTP server responses
	//   - error: Any error encountered during delivery
	Deliver(ctx context.Context, conn net.Conn, myMail *pmail.Mail, to smtp.Address) ([]pmail.Response, error)

	// NewConn establishes a new connection to one of the provided SMTP hosts.
	// It randomly selects a host from the provided list and attempts to connect.
	//
	// Parameters:
	//   - ctx: Context for the connection operation
	//   - hosts: List of SMTP server hostnames to try connecting to
	//
	// Returns:
	//   - net.Conn: Established network connection
	//   - error: Any error encountered during connection
	NewConn(ctx context.Context, hosts []string) (net.Conn, error)

	// SendMail handles the complete process of sending an email to all recipients.
	// It manages concurrent delivery attempts and collects results.
	//
	// Parameters:
	//   - ctx: Context for the sending operation
	//   - mail: Email to be sent
	//
	// Returns:
	//   - map[string][]pmail.Response: Map of recipient addresses to their SMTP responses
	//   - map[string]error: Map of recipient addresses to any errors encountered
	SendMail(ctx context.Context, mail *pmail.Mail) (map[string][]pmail.Response, map[string]error)
}
