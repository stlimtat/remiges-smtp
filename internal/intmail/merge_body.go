package intmail

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
)

const (
	MergeBodyProcessorType = "mergeBody"
)

type MergeBodyProcessor struct {
	cfg config.MailProcessorConfig
}

func (p *MergeBodyProcessor) Init(
	ctx context.Context,
	cfg config.MailProcessorConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", MergeBodyProcessorType).
		Int("index", p.cfg.Index).
		Interface("args", p.cfg.Args).
		Logger()
	logger.Debug().Msg("MergeBodyProcessor Init")
	p.cfg = cfg
	return nil
}

func (p *MergeBodyProcessor) Index() int {
	return p.cfg.Index
}

func (_ *MergeBodyProcessor) Process(
	ctx context.Context,
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("MergeBodyProcessor")

	mailBodyHeaders := make([]byte, 0)
	for key, value := range inMail.BodyHeaders {
		mailBodyHeaders = append(
			mailBodyHeaders,
			[]byte(key+": "+string(value)+"\r\n")...,
		)
	}

	mailBodyHeaders = append(mailBodyHeaders, []byte("\r\n")...)
	inMail.Body = append(mailBodyHeaders, inMail.Body...)
	inMail.Body = append(inMail.Body, []byte("\r\n\r\n")...)

	logger.Info().Bytes("body", inMail.Body).Msg("MergeBodyProcessor")
	return inMail, nil
}
