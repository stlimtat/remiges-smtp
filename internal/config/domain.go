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
	moxConfig.Domain           `mapstructure:",omitempty"`
	ClientSettingsDomain       string      `mapstructure:"client-settings-domain,omitempty"`
	Description                string      `mapstructure:"description,omitempty"`
	DomainStr                  string      `mapstructure:"domain-str,omitempty"`
	DKIM                       *DKIMConfig `mapstructure:"dkim,omitempty"`
	LocalpartCaseSensitive     bool        `mapstructure:"localpart-case-sensitive,omitempty"`
	LocalpartCatchallSeparator string      `mapstructure:"localpart-catchall-separator,omitempty"`
	ReportsOnly                bool        `mapstructure:"reports-only,omitempty"`
	// Aliases                    map[string]Alias
	// DMARC                      *DMARC
	// MTASTS                     *MTASTS
	// Routes                     []Route
	// TLSRPT                     *TLSRPT
}

func DefaultDomainConfig(
	ctx context.Context,
) map[string]*DomainConfig {
	// yaml cannot have a map key with a dot, so we use a string key
	result := map[string]*DomainConfig{
		"stlimnet": {
			ClientSettingsDomain:       "",
			DKIM:                       DefaultDKIMConfig(ctx),
			Domain:                     moxConfig.Domain{},
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
		c.Domain.DKIM = c.DKIM.DKIM
	}
	if c.ClientSettingsDomain != "" {
		c.Domain.ClientSettingsDomain = c.ClientSettingsDomain
	}
	if c.Description != "" {
		c.Domain.Description = c.Description
	}
	if c.DomainStr != "" {
		c.Domain.Domain = moxDns.Domain{ASCII: c.DomainStr}
	}
	if c.LocalpartCaseSensitive {
		c.Domain.LocalpartCaseSensitive = c.LocalpartCaseSensitive
	}
	if c.LocalpartCatchallSeparator != "" {
		c.Domain.LocalpartCatchallSeparator = c.LocalpartCatchallSeparator
	}
	if c.ReportsOnly {
		c.Domain.ReportsOnly = c.ReportsOnly
	}
	return nil
}
