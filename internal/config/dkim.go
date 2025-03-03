package config

import (
	"context"
	"time"

	moxConfig "github.com/mjl-/mox/config"
	"github.com/mjl-/mox/dns"
	"github.com/rs/zerolog"
)

type DKIMConfig struct {
	moxConfig.DKIM `mapstructure:",omitempty"`
	MoxSelectors   map[string]MoxSelector `mapstructure:"selectors,omitempty"`
	MoxSign        []string               `mapstructure:"sign,omitempty"`
}

type MoxSelector struct {
	Algorithm        string              `mapstructure:"algorithm"`
	Canonicalization MoxCanonicalization `mapstructure:"canonicalization,omitempty"`
	Domain           string              `mapstructure:"domain"`
	DontSealHeaders  bool                `mapstructure:"dont-seal-headers,omitempty"`
	Expiration       string              `mapstructure:"expiration,omitempty"`
	Hash             string              `mapstructure:"hash"`
	Headers          []string            `mapstructure:"headers,omitempty"`
	PrivateKeyFile   string              `mapstructure:"private-key-file,omitempty"`
}

type MoxCanonicalization struct {
	HeaderRelaxed bool `mapstructure:"header-relaxed"`
	BodyRelaxed   bool `mapstructure:"body-relaxed"`
}

func DefaultDKIMConfig(
	ctx context.Context,
) *DKIMConfig {
	logger := zerolog.Ctx(ctx)
	result := &DKIMConfig{
		DKIM: moxConfig.DKIM{
			Selectors: make(map[string]moxConfig.Selector),
		},
		MoxSelectors: map[string]MoxSelector{
			"key001": {
				Algorithm: "rsa-sha256",
				Canonicalization: MoxCanonicalization{
					HeaderRelaxed: true,
					BodyRelaxed:   true,
				},
				Domain:          "stlim.net",
				DontSealHeaders: true,
				Expiration:      "24h",
				Hash:            "sha256",
				Headers:         nil,
			},
		},
		MoxSign: make([]string, 0),
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
	err := c.TransformSign(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("DKIMConfig.TransformSign")
		return err
	}
	return nil
}

func (c *DKIMConfig) TransformSelector(
	ctx context.Context,
	selectorName string,
	moxSelector *MoxSelector,
) error {
	logger := zerolog.Ctx(ctx).With().Str("selector", selectorName).Logger()
	var err error
	if c.DKIM.Selectors == nil {
		c.DKIM.Selectors = make(map[string]moxConfig.Selector)
	}

	result, ok := c.DKIM.Selectors[selectorName]
	if !ok {
		result = moxConfig.Selector{}
	}
	result.Algorithm = moxSelector.Algorithm
	result.Canonicalization = moxConfig.Canonicalization{
		HeaderRelaxed: moxSelector.Canonicalization.HeaderRelaxed,
		BodyRelaxed:   moxSelector.Canonicalization.BodyRelaxed,
	}
	result.Domain, err = dns.ParseDomain(moxSelector.Domain)
	if err != nil {
		logger.Error().Err(err).Msg("TransformSelector.ParseDomain")
		return err
	}
	result.DontSealHeaders = moxSelector.DontSealHeaders
	result.Expiration = moxSelector.Expiration
	if moxSelector.Expiration != "" {
		expiration, err := time.ParseDuration(moxSelector.Expiration)
		if err != nil {
			logger.Error().Err(err).Msg("TransformSelector.ParseDuration")
			return err
		}
		result.ExpirationSeconds = int(expiration.Seconds())
	}
	result.Hash = moxSelector.Hash
	// This is used in mox/dkimsign.go:L23 - DKIMSelectors
	result.HashEffective = moxSelector.Hash
	result.Headers = moxSelector.Headers

	c.DKIM.Selectors[selectorName] = result
	return nil
}

func (c *DKIMConfig) TransformSign(
	_ context.Context,
) error {
	if c.DKIM.Sign == nil {
		c.DKIM.Sign = make([]string, 0)
	}
	c.DKIM.Sign = append(c.DKIM.Sign, c.MoxSign...)
	return nil
}
