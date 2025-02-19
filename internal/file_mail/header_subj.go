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
	ctx context.Context,
	cfg config.FileMailConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", HeaderSubjectTransformerType).
		Int("index", t.Cfg.Index).
		Interface("args", t.Cfg.Args).
		Logger()
	logger.Debug().Msg("HeaderSubjectTransformer Init")
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
	fileInfo *file.FileInfo,
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderSubjectTransformer")

	subjectBytes, ok := inMail.Metadata[HeaderSubjectKey]
	if !ok {
		subjectBytes = t.SubjectBytes
	}

	inMail.Subject = subjectBytes
	logger.Debug().Bytes("subject", inMail.Subject).Msg("HeaderSubjectTransformer")

	return inMail, nil
}
