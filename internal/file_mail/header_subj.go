package file_mail

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
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
	subjectAny, ok := cfg.Args[HeaderConfigArgDefault]
	if !ok {
		subjectAny = "no subject"
	}
	t.SubjectStr = subjectAny.(string)
	t.SubjectBytes = []byte(t.SubjectStr)
	return nil
}

func (t *HeaderSubjectTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderSubjectTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderSubjectTransformer")

	subjectBytes, ok := inMail.Metadata[input.HeaderSubjectKey]
	if !ok {
		subjectBytes = t.SubjectBytes
	}

	inMail.Subject = subjectBytes
	logger.Debug().
		Bytes(HeaderSubjectKey, inMail.Subject).
		Msg("HeaderSubjectTransformer")

	return inMail, nil
}
