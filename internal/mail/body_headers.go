package mail

import (
	"bytes"
	"context"
	"fmt"

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
	_ context.Context,
	cfg config.MailProcessorConfig,
) error {
	p.cfg = cfg
	return nil
}

func (p *BodyHeadersProcessor) Index() int {
	return p.cfg.Index
}

func (_ *BodyHeadersProcessor) Process(
	ctx context.Context,
	inMail *Mail,
) (outMail *Mail, err error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("BodyHeadersProcessor")

	inMail.BodyHeaders = make(map[string][]byte)
	// separate the mail header from the body
	mailSections := bytes.Split(inMail.Body, []byte("\r\n\r\n"))
	if len(mailSections) > 2 {
		return nil, fmt.Errorf("invalid mail body")
	}
	if len(mailSections) == 2 {
		bodyHeadersBytes := mailSections[0]
		for _, header := range bytes.Split(bodyHeadersBytes, []byte("\r\n")) {
			headerParts := bytes.Split(header, []byte(": "))
			if len(headerParts) != 2 {
				return nil, fmt.Errorf("invalid header: %s", header)
			}
			inMail.BodyHeaders[string(headerParts[0])] = headerParts[1]
		}
		// Add the headers to the mail
		inMail.Body = mailSections[1]
	}
	inMail.BodyHeaders["From"] = []byte(inMail.From.String())
	inMail.BodyHeaders["To"] = []byte(inMail.To.String())

	return inMail, nil
}
