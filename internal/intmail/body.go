package intmail

import (
	"bytes"
	"context"
	"fmt"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
)

const (
	BodyProcessorType = "body"
)

type BodyProcessor struct {
	cfg config.MailProcessorConfig
}

func (p *BodyProcessor) Init(
	ctx context.Context,
	cfg config.MailProcessorConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", BodyProcessorType).
		Int("index", p.cfg.Index).
		Interface("args", p.cfg.Args).
		Logger()
	logger.Debug().Msg("BodyProcessor Init")
	p.cfg = cfg
	return nil
}

func (p *BodyProcessor) Index() int {
	return p.cfg.Index
}

func (_ *BodyProcessor) Process(
	ctx context.Context,
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("BodyProcessor")

	if inMail.BodyHeaders == nil {
		inMail.BodyHeaders = make(map[string][]byte)
	}

	// Forced operation to replace all \n with \r\n
	re := regexp.MustCompile(`\r?\n`)
	result := re.ReplaceAll(inMail.Body, []byte("\r\n"))
	// Remove leading and trailing whitespace
	result = bytes.TrimSpace(result)
	inMail.Body = result

	// If the body starts with --, it is a multipart message
	if bytes.HasPrefix(result, []byte("--")) {
		logger.Debug().Bytes("result", result).Msg("BodyHeadersProcessor - multipart message")
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
	inMail.Body = bytes.TrimSpace(result)

	return inMail, nil
}
