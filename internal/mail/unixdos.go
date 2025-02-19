package mail

import (
	"context"
	"regexp"

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
	logger.Debug().Msg("UnixDosProcessor")

	re := regexp.MustCompile(`\r?\n`)
	inMail.Body = re.ReplaceAll(inMail.Body, []byte("\r\n"))

	return inMail, nil
}
