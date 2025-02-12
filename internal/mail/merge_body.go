package mail

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
)

const (
	MergeBodyProcessorType = "mergeBody"
)

type MergeBodyProcessor struct {
	cfg config.MailProcessorConfig
}

func (p *MergeBodyProcessor) Init(
	_ context.Context,
	cfg config.MailProcessorConfig,
) error {
	p.cfg = cfg
	return nil
}

func (p *MergeBodyProcessor) Index() int {
	return p.cfg.Index
}

func (_ *MergeBodyProcessor) Process(
	ctx context.Context,
	inMail *Mail,
) (outMail *Mail, err error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("MergeBodyProcessor")

	mailBodyHeaders := make([]byte, 0)
	for key, value := range inMail.BodyHeaders {
		mailBodyHeaders = append(
			mailBodyHeaders,
			[]byte(key+": "+string(value)+"\r\n")...,
		)
	}

	mailBodyHeaders = append(mailBodyHeaders, []byte("\r\n")...)
	inMail.Body = append(mailBodyHeaders, inMail.Body...)

	return inMail, nil
}
