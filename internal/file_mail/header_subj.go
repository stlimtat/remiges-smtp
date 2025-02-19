package file_mail

import (
	"context"

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
	Cfg          config.FileMailConfig
	SubjectBytes []byte
	SubjectStr   string
}

func (t *HeaderSubjectTransformer) Init(
	_ context.Context,
	cfg config.FileMailConfig,
) error {
	t.Cfg = cfg
	var ok bool
	t.SubjectStr, ok = cfg.Args["default"]
	if !ok {
		t.SubjectStr = "no subject"
	}
	t.SubjectBytes = []byte(t.SubjectStr)
	return nil
}

func (t *HeaderSubjectTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderSubjectTransformer) Transform(
	ctx context.Context,
	_ *file.FileInfo,
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HeaderSubjectTransformer")

	subjectBytes, ok := inMail.Headers[HeaderSubjectKey]
	if !ok {
		subjectBytes = t.SubjectBytes
	}

	inMail.Subject = subjectBytes

	return inMail, nil
}
