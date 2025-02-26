package intmail

import (
	"context"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
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
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("UnixDosProcessor")

	re := regexp.MustCompile(`\r?\n`)
	inMail.Body = re.ReplaceAll(inMail.Body, []byte("\r\n"))

	return inMail, nil
}
