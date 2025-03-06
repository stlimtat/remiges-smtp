package config

import (
	"context"

	moxConfig "github.com/mjl-/mox/config"
	moxDns "github.com/mjl-/mox/dns"
	"github.com/rs/zerolog"
)

const (
	DomainStlimNet = "stlim.net"
)

type DomainConfig struct {
	DKIM                       *DKIMConfig      `mapstructure:"dkim,omitempty"`
	Domain                     moxDns.Domain    `mapstructure:",omitempty"`
	DomainStr                  string           `mapstructure:"domain-str,omitempty"`
	LocalpartCaseSensitive     bool             `mapstructure:"localpart-case-sensitive,omitempty"`
	LocalpartCatchallSeparator string           `mapstructure:"localpart-catchall-separator,omitempty"`
	MoxDomain                  moxConfig.Domain `mapstructure:",omitempty"`
	ReportsOnly                bool             `mapstructure:"reports-only,omitempty"`
}

func DefaultDomainConfig(
	ctx context.Context,
) map[string]*DomainConfig {
	// yaml cannot have a map key with a dot, so we use a string key
	result := map[string]*DomainConfig{
		"stlimnet": {
			DKIM:                       DefaultDKIMConfig(ctx),
			Domain:                     moxDns.Domain{ASCII: DomainStlimNet},
			DomainStr:                  DomainStlimNet,
			LocalpartCaseSensitive:     false,
			LocalpartCatchallSeparator: "+",
			ReportsOnly:                false,
		},
	}
	return result
}

func (c *DomainConfig) Transform(
	ctx context.Context,
) error {
	logger := zerolog.Ctx(ctx)
	if c.DKIM != nil {
		err := c.DKIM.Transform(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("DomainConfig.Transform.DKIM")
			return err
		}
	}
	if c.DomainStr != "" {
		c.Domain = moxDns.Domain{ASCII: c.DomainStr}
	}
	c.MoxDomain = moxConfig.Domain{
		Domain:                     c.Domain,
		LocalpartCaseSensitive:     c.LocalpartCaseSensitive,
		LocalpartCatchallSeparator: c.LocalpartCatchallSeparator,
		ReportsOnly:                c.ReportsOnly,
	}
	return nil
}
