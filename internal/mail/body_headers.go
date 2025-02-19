package mail

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
)

const (
	BodyHeadersProcessorType = "bodyHeaders"
)

type BodyHeadersProcessor struct {
	cfg config.MailProcessorConfig
}

func (p *BodyHeadersProcessor) Init(
	ctx context.Context,
	cfg config.MailProcessorConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", BodyHeadersProcessorType).
		Int("index", p.cfg.Index).
		Interface("args", p.cfg.Args).
		Logger()
	logger.Debug().Msg("BodyHeadersProcessor Init")
	p.cfg = cfg
	return nil
}

func (p *BodyHeadersProcessor) Index() int {
	return p.cfg.Index
}

func (_ *BodyHeadersProcessor) Process(
	ctx context.Context,
	inMail *Mail,
) (*Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Interface("from", inMail.From).
		Interface("to", inMail.To).
		Bytes("subject", inMail.Subject).
		Msg("BodyHeadersProcessor")

	if inMail.BodyHeaders == nil {
		inMail.BodyHeaders = make(map[string][]byte)
	}
	inMail.BodyHeaders["Content-Type"] = inMail.ContentType
	inMail.BodyHeaders["From"] = []byte(inMail.From.String())
	inMail.BodyHeaders["Subject"] = inMail.Subject

	toBytes := []byte{}
	for _, to := range inMail.To {
		toBytes = append(toBytes, to.String()...)
		toBytes = append(toBytes, ',')
	}
	inMail.BodyHeaders["To"] = toBytes[:len(toBytes)-1]

	return inMail, nil
}
