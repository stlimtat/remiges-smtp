package intmail

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
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
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Interface("from", inMail.From).
		Bytes("msgid", inMail.MsgID).
		Bytes("subject", inMail.Subject).
		Interface("to", inMail.To).
		Msg("BodyHeadersProcessor")

	if inMail.BodyHeaders == nil {
		inMail.BodyHeaders = make(map[string][]byte)
	}
	inMail.BodyHeaders[input.HeaderContentTypeKey] = inMail.ContentType
	now := time.Now().Format(time.RFC1123Z)
	inMail.BodyHeaders[input.HeaderDateKey] = []byte(now)
	inMail.BodyHeaders[input.HeaderFromKey] = []byte(inMail.From.String())
	inMail.BodyHeaders[input.HeaderMsgIDKey] = inMail.MsgID
	inMail.BodyHeaders[input.HeaderSubjectKey] = inMail.Subject

	toBytes := []byte{}
	for _, to := range inMail.To {
		toBytes = append(toBytes, to.String()...)
		toBytes = append(toBytes, ',')
	}
	inMail.BodyHeaders[input.HeaderToKey] = toBytes[:len(toBytes)-1]

	return inMail, nil
}
