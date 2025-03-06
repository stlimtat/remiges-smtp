package config

import (
	"context"
	"time"

	moxDkim "github.com/mjl-/mox/dkim"
	moxDns "github.com/mjl-/mox/dns"
	"github.com/rs/zerolog"
)

type DKIMConfig struct {
	Selectors    map[string]moxDkim.Selector `mapstructure:",omitempty"`
	MoxSelectors map[string]MoxSelector      `mapstructure:"selectors,omitempty"`
}
type MoxSelector struct {
	Algorithm      string        `mapstructure:"algorithm"`
	BodyRelaxed    bool          `mapstructure:"body-relaxed"`
	Expiration     time.Duration `mapstructure:"expiration,omitempty"`
	Hash           string        `mapstructure:"hash"`
	HeaderRelaxed  bool          `mapstructure:"header-relaxed"`
	Headers        []string      `mapstructure:"headers,omitempty"`
	PrivateKeyFile string        `mapstructure:"private-key-file,omitempty"`
	SealHeaders    bool          `mapstructure:"seal-headers,omitempty"`
	SelectorDomain string        `mapstructure:"selector-domain"`
}

func DefaultDKIMConfig(
	ctx context.Context,
) *DKIMConfig {
	logger := zerolog.Ctx(ctx)
	result := &DKIMConfig{
		Selectors: make(map[string]moxDkim.Selector),
		MoxSelectors: map[string]MoxSelector{
			"key001": {
				Algorithm:      "rsa",
				BodyRelaxed:    true,
				Hash:           "sha256",
				HeaderRelaxed:  true,
				Headers:        []string{"from", "to", "subject", "date", "message-id", "content-type"},
				Expiration:     time.Hour * 24,
				PrivateKeyFile: "~/go/src/github.com/stlimtat/remiges-smtp/config/key001.key",
				SealHeaders:    true,
				SelectorDomain: "key001",
			},
		},
	}
	err := result.Transform(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("DefaultDKIMConfig.Transform")
	}
	return result
}

func (c *DKIMConfig) Transform(
	ctx context.Context,
) error {
	logger := zerolog.Ctx(ctx)
	for selectorName, moxSelector := range c.MoxSelectors {
		sublogger := logger.With().Str("selector", selectorName).Logger()
		err := c.TransformSelector(ctx, selectorName, &moxSelector)
		if err != nil {
			sublogger.Error().Err(err).Msg("DKIMConfig.TransformSelector")
			return err
		}
	}
	return nil
}

func (c *DKIMConfig) TransformSelector(
	ctx context.Context,
	selectorName string,
	moxSelector *MoxSelector,
) error {
	logger := zerolog.Ctx(ctx).
		With().
		Str("selector", selectorName).
		Logger()
	var err error
	if c.Selectors == nil {
		c.Selectors = make(map[string]moxDkim.Selector)
	}

	result, ok := c.Selectors[selectorName]
	if !ok {
		result = moxDkim.Selector{}
	}
	result.BodyRelaxed = moxSelector.BodyRelaxed
	result.Domain, err = moxDns.ParseDomain(moxSelector.SelectorDomain)
	if err != nil {
		logger.Error().Err(err).Msg("TransformSelector.ParseDomain")
		return err
	}
	result.Expiration = moxSelector.Expiration
	result.Hash = moxSelector.Hash
	result.HeaderRelaxed = moxSelector.HeaderRelaxed
	// This is used in mox/dkimsign.go:L23 - DKIMSelectors
	result.Headers = moxSelector.Headers
	result.SealHeaders = moxSelector.SealHeaders

	c.Selectors[selectorName] = result
	return nil
}
