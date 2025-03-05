package intmail

import (
	"bytes"
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
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
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("MergeBodyProcessor")

	body := bytes.TrimSpace(inMail.Body)
	finalBody := inMail.Headers
	finalBody = append(finalBody, body...)
	finalBody = append(finalBody, []byte("\r\n\r\n")...)
	inMail.FinalBody = finalBody

	logger.Info().Bytes("final_body", inMail.FinalBody).Msg("MergeBodyProcessor")
	return inMail, nil
}
