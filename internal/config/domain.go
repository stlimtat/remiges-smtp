package config

import (
	"context"

	moxConfig "github.com/mjl-/mox/config"
)

type DomainConfig struct {
	moxConfig.Domain           `mapstructure:",omitempty"`
	ClientSettingsDomain       string      `mapstructure:"client_settings_domain,omitempty"`
	Description                string      `mapstructure:"description,omitempty"`
	DomainStr                  string      `mapstructure:"domain,omitempty"`
	DKIM                       *DKIMConfig `mapstructure:"dkim,omitempty"`
	LocalpartCaseSensitive     bool        `mapstructure:"localpart_case_sensitive,omitempty"`
	LocalpartCatchallSeparator string      `mapstructure:"localpart_catchall_separator,omitempty"`
	ReportsOnly                bool        `mapstructure:"reports_only,omitempty"`
	// Aliases                    map[string]Alias
	// DMARC                      *DMARC
	// MTASTS                     *MTASTS
	// Routes                     []Route
	// TLSRPT                     *TLSRPT
}

func DefaultDomainConfig(
	ctx context.Context,
) *DomainConfig {
	result := &DomainConfig{
		ClientSettingsDomain:       "",
		DKIM:                       DefaultDKIMConfig(ctx),
		Domain:                     moxConfig.Domain{},
		DomainStr:                  "",
		LocalpartCaseSensitive:     false,
		LocalpartCatchallSeparator: "+",
		ReportsOnly:                false,
	}
	return result
}
