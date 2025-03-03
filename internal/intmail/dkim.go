package intmail

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/mitchellh/mapstructure"
	"github.com/mjl-/mox/dkim"
	"github.com/mjl-/mox/mox-"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
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

	dkimHeaders, err := dkim.Sign(
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
	// add dkim headers to the mail
	inMail.DKIMHeaders = []byte(dkimHeaders)
	inMail.Body = append(inMail.DKIMHeaders, inMail.Body...)

	return inMail, nil
}
