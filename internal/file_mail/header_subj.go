package file_mail

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
)

const (
	HeaderSubjectTransformerType = "header_subject"
	HeaderSubjectKey             = "Subject"
)

type HeaderSubjectTransformer struct {
	Cfg        config.FileMailConfig
	SubjectStr string
}

func (t *HeaderSubjectTransformer) Init(
	_ context.Context,
	cfg config.FileMailConfig,
) error {
	t.Cfg = cfg
	return nil
}

func (t *HeaderSubjectTransformer) Index() int {
	return t.Cfg.Index
}

func (_ *HeaderSubjectTransformer) Transform(
	ctx context.Context,
	_ *file.FileInfo,
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HeaderSubjectTransformer")

	subjectBytes, ok := inMail.Headers[HeaderSubjectKey]
	if !ok {
		logger.Error().Msg("subject header not found")
		return nil, fmt.Errorf("subject header not found")
	}

	subjectStr := string(subjectBytes)
	inMail.Subject = subjectStr

	return inMail, nil
}
