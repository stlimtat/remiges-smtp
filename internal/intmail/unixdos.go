package intmail

import (
	"context"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
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
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("UnixDosProcessor")
	if inMail == nil || inMail.Body == nil {
		logger.Error().Msg("UnixDosProcessor: inMail is nil or inMail.Body is nil")
		return inMail, nil
	}

	re := regexp.MustCompile(`\r?\n`)
	inMail.Body = re.ReplaceAll(inMail.Body, []byte("\r\n"))

	return inMail, nil
}
