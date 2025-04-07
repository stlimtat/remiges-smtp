package sendmail

import (
	"context"
	"net"

	"github.com/mjl-/mox/smtpclient"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"golang.org/x/net/proxy"
)

type DefaultNetDialerFactory struct {
	cfg config.DialerConfig
}

func NewDefaultDialerFactory(
	ctx context.Context,
	cfg config.DialerConfig,
) *DefaultNetDialerFactory {
	result := &DefaultNetDialerFactory{
		cfg: cfg,
	}
	return result
}

func (n *DefaultNetDialerFactory) NewDialer(
	ctx context.Context,
) (smtpclient.Dialer, error) {
	logger := zerolog.Ctx(ctx)
	var result smtpclient.Dialer
	baseDialer := &net.Dialer{
		Timeout: n.cfg.Timeout,
	}
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
		result = baseDialer
	}
	return result, nil
}
