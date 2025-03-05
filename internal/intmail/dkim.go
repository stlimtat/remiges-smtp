package intmail

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/mitchellh/mapstructure"
	moxDkim "github.com/mjl-/mox/dkim"
	"github.com/mjl-/mox/mox-"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/crypto"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	DKIMProcessorType = "dkim"
)

type DKIMProcessor struct {
	Cfg       config.MailProcessorConfig
	DomainCfg *config.DomainConfig
	SLogger   *slog.Logger
}

func (p *DKIMProcessor) Init(
	ctx context.Context,
	cfg config.MailProcessorConfig,
) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("DKIMProcessor Init")
	p.Cfg = cfg

	p.DomainCfg = &config.DomainConfig{}
	err := mapstructure.Decode(p.Cfg.Args, p.DomainCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("DKIMProcessor: decode")
	}

	err = p.DomainCfg.Transform(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("DKIMProcessor: transform")
	}

	return nil
}

func (p *DKIMProcessor) Index() int {
	return p.Cfg.Index
}

func (p *DKIMProcessor) InitDKIMCrypto(
	ctx context.Context,
	loader crypto.IKeyLoader,
) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("DKIMProcessor: InitDKIMCrypto")

	selectors := p.DomainCfg.DKIM.DKIM.Selectors

	for selectorName, selector := range selectors { //nolint:gocritic // This was inherited from mox
		privateKeyPath, err := utils.ValidateIO(ctx, selector.PrivateKeyFile, true)
		if err != nil {
			logger.Error().Err(err).Msg("DKIMProcessor: InitDKIMCrypto: ValidateIO")
			return err
		}
		signer, err := loader.LoadPrivateKey(ctx, selector.Algorithm, privateKeyPath)
		if err != nil {
			logger.Error().Err(err).Msg("DKIMProcessor: InitDKIMCrypto: LoadPrivateKey")
			return err
		}
		selector.Key = signer
		selector.PrivateKeyFile = privateKeyPath
		p.DomainCfg.DKIM.DKIM.Selectors[selectorName] = selector
	}

	return nil
}

func (p *DKIMProcessor) Process(
	ctx context.Context,
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Interface("from", inMail.From).
		Msg("DKIMProcessor")

	selectors := mox.DKIMSelectors(p.DomainCfg.DKIM.DKIM)
	if len(selectors) == 0 {
		logger.Debug().Msg("DKIMProcessor: no selectors")
		return inMail, nil
	}

	canonical := mox.CanonicalLocalpart(inMail.From.Localpart, p.DomainCfg.Domain)

	dkimHeaders, err := moxDkim.Sign(
		ctx,
		p.SLogger,
		canonical,
		inMail.From.Domain,
		selectors,
		true,
		bytes.NewReader(inMail.Body),
	)
	if err != nil {
		logger.Error().Err(err).Msg("DKIMProcessor: sign")
		return inMail, err
	}
	// m1 := regexp.MustCompile(`i=([^;]+);`)
	// result := m1.ReplaceAllString(dkimHeaders, "")

	// add dkim headers to the mail
	inMail.DKIMHeaders = []byte(dkimHeaders)
	inMail.Body = append(inMail.DKIMHeaders, inMail.Body...)

	return inMail, nil
}
