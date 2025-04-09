package sendmail

import (
	"context"
	"net"

	"github.com/mjl-/mox/smtpclient"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"golang.org/x/net/proxy"
)

// DefaultNetDialerFactory implements INetDialerFactory to create network dialers
// for SMTP connections. It supports both direct TCP connections and SOCKS5 proxy
// connections based on the provided configuration.
type DefaultNetDialerFactory struct {
	// cfg holds the dialer configuration including timeout settings and proxy details
	cfg config.DialerConfig
}

// NewDefaultDialerFactory creates a new DefaultNetDialerFactory with the specified
// configuration. It initializes a factory that can create dialers for SMTP connections.
//
// Parameters:
//   - ctx: Context for the factory creation
//   - cfg: Configuration for the dialer including timeout and proxy settings
//
// Returns:
//   - *DefaultNetDialerFactory: A new dialer factory instance
func NewDefaultDialerFactory(
	ctx context.Context,
	cfg config.DialerConfig,
) *DefaultNetDialerFactory {
	result := &DefaultNetDialerFactory{
		cfg: cfg,
	}
	return result
}

// NewDialer creates a new SMTP client dialer based on the factory's configuration.
// If SOCKS5 proxy settings are provided, it creates a proxy-enabled dialer;
// otherwise, it returns a standard TCP dialer.
//
// Parameters:
//   - ctx: Context for the dialer creation
//
// Returns:
//   - smtpclient.Dialer: A configured dialer (either direct TCP or SOCKS5 proxy)
//   - error: Any error encountered during dialer creation
func (n *DefaultNetDialerFactory) NewDialer(
	ctx context.Context,
) (smtpclient.Dialer, error) {
	logger := zerolog.Ctx(ctx)
	var result smtpclient.Dialer

	// Create base TCP dialer with configured timeout
	baseDialer := &net.Dialer{
		Timeout: n.cfg.Timeout,
	}

	// If SOCKS5 proxy is configured, create a proxy-enabled dialer
	if n.cfg.Socks5 != "" {
		socks5Dialer, err := proxy.SOCKS5(
			TCPNetwork,
			n.cfg.Socks5,
			&n.cfg.Auth,
			baseDialer,
		)
		if err != nil {
			logger.Error().Err(err).Msg("NewDialer.proxy.SOCKS5")
			return nil, err
		}
		result = socks5Dialer.(proxy.ContextDialer)
	} else {
		// Use standard TCP dialer if no proxy is configured
		result = baseDialer
	}
	return result, nil
}
