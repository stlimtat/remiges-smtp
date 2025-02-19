package mail

import (
	"bytes"
	"context"
	"fmt"
	"regexp"

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

func (p *BodyHeadersProcessor) Process(
	ctx context.Context,
	inMail *Mail,
) (*Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Bytes("body", inMail.Body).Msg("BodyHeadersProcessor")

	inMail.BodyHeaders = make(map[string][]byte)

	// Forced operation to replace all \n with \r\n
	re := regexp.MustCompile(`\r?\n`)
	result := re.ReplaceAll(inMail.Body, []byte("\r\n"))
	// Remove leading and trailing whitespace
	result = bytes.TrimSpace(result)
	inMail.Body = result

	// If the body starts with --, it is a multipart message
	if bytes.HasPrefix(result, []byte("--")) {
		logger.Debug().Bytes("result", result).Msg("BodyHeadersProcessor - multipart message")
		inMail = p.PopulateFromTo(ctx, inMail)
		return inMail, nil
	}
	// separate the mail header from the body
	mailSections := bytes.Split(result, []byte("\r\n\r\n"))
	if len(mailSections) > 2 {
		logger.Error().Int("mailSections", len(mailSections)).Msg("invalid mail body")
		return nil, fmt.Errorf("invalid mail body")
	}
	if len(mailSections) == 2 {
		bodyHeadersBytes := mailSections[0]
		// This overrides the from and to headers
		for _, header := range bytes.Split(bodyHeadersBytes, []byte("\r\n")) {
			headerParts := bytes.Split(header, []byte(": "))
			if len(headerParts) != 2 {
				return nil, fmt.Errorf("invalid header: %s", header)
			}
			inMail.BodyHeaders[string(headerParts[0])] = headerParts[1]
		}
		// Add the headers to the mail
		result = mailSections[1]
	}
	// Re-replacing the from and to body headers
	inMail = p.PopulateFromTo(ctx, inMail)
	inMail.Body = bytes.TrimSpace(result)

	return inMail, nil
}

func (_ *BodyHeadersProcessor) PopulateFromTo(
	ctx context.Context,
	inMail *Mail,
) *Mail {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Interface("from", inMail.From).Interface("to", inMail.To).Msg("PopulateFromTo")

	inMail.BodyHeaders["From"] = []byte(inMail.From.String())
	toBytes := []byte{}
	for _, to := range inMail.To {
		toBytes = append(toBytes, to.String()...)
		toBytes = append(toBytes, ',')
	}
	inMail.BodyHeaders["To"] = toBytes[:len(toBytes)-1]

	return inMail
}
