package mail

import (
	"bytes"
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
)

const (
	UnixDosProcessorType = "unixdos"
)

type UnixDosProcessor struct {
	cfg config.MailProcessorConfig
}

func (p *UnixDosProcessor) Init(
	_ context.Context,
	cfg config.MailProcessorConfig,
) error {
	p.cfg = cfg
	return nil
}

func (p *UnixDosProcessor) Index() int {
	return p.cfg.Index
}

func (_ *UnixDosProcessor) Process(
	ctx context.Context,
	inMail *Mail,
) (outMail *Mail, err error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("UnixDosProcessor")

	inMail.Body = bytes.ReplaceAll(inMail.Body, []byte("\n"), []byte("\r\n"))

	return inMail, nil
}
