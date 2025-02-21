package mail

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
)

const (
	DKIMProcessorType = "dkim"
)

type DKIMProcessor struct {
	cfg     config.MailProcessorConfig
	dkimCfg config.DKIMConfig
}

func (p *DKIMProcessor) Init(
	ctx context.Context,
	cfg config.MailProcessorConfig,
) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("DKIMProcessor Init")
	p.cfg = cfg
	dkimCfgAny, ok := cfg.Args[DKIMProcessorType]
	if !ok {
		logger.Fatal().Msg("DKIMProcessor: no config")
	}
	logger.Debug().Interface("dkimCfgAny", dkimCfgAny).Msg("DKIMProcessor: dkimCfgAny")
	p.dkimCfg = config.DKIMConfig{}
	err := mapstructure.Decode(dkimCfgAny, &p.dkimCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("DKIMProcessor: decode")
	}
	err = p.dkimCfg.Transform(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("DKIMProcessor: transform")
	}
	return nil
}

func (p *DKIMProcessor) Index() int {
	return p.cfg.Index
}

func (_ *DKIMProcessor) Process(
	ctx context.Context,
	inMail *Mail,
) (*Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Interface("from", inMail.From).
		Msg("DKIMProcessor")

	// selectors := mox.DKIMSelectors(confDom.DKIM)
	// if len(selectors) > 0 {
	// 	canonical := mox.CanonicalLocalpart(msgFrom.Localpart, confDom)
	// 	dkimHeaders, err := dkim.Sign(
	// 		ctx,
	// 		c.log.Logger,
	// 		canonical,
	// 		msgFrom.Domain,
	// 		selectors,
	// 		c.msgsmtputf8,
	// 		store.FileMsgReader(msgPrefix, dataFile))
	//  if err != nil {
	// 		c.log.Errorx("dkim sign for domain", err, slog.Any("domain", msgFrom.Domain))
	// 		metricServerErrors.WithLabelValues("dkimsign").Inc()
	// 	} else {
	// 		msgPrefix = append(msgPrefix, []byte(dkimHeaders)...)
	// 	}
	// }

	return inMail, nil
}
