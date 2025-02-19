package file_mail

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
)

const (
	HeaderContentTypeTransformerType  = "header_contenttype"
	HeaderContentTypeKey              = "Content-Type"
	HeaderContentTypeConfigArgType    = "type"
	HeaderContentTypeConfigArgDefault = "default"
)

type HeaderContentTypeTransformer struct {
	Cfg             config.FileMailConfig
	ContentType     string
	ContentTypeStr  string
	ContentTypeType config.FromType
}

func (t *HeaderContentTypeTransformer) Init(
	ctx context.Context,
	cfg config.FileMailConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", HeaderContentTypeTransformerType).
		Int("index", t.Cfg.Index).
		Interface("args", t.Cfg.Args).
		Logger()
	logger.Debug().Msg("HeaderContentTypeTransformer Init")

	t.Cfg = cfg
	fromTypeStr, ok := t.Cfg.Args[HeaderFromConfigArgType]
	if !ok {
		fromTypeStr = config.FromTypeHeadersStr
	}
	switch fromTypeStr {
	case config.FromTypeDefaultStr:
		t.ContentTypeType = config.FromTypeDefault
	case config.FromTypeHeadersStr:
		t.ContentTypeType = config.FromTypeHeaders
	default:
		t.ContentTypeType = config.FromTypeHeaders
	}

	contentTypeStr, ok := t.Cfg.Args[HeaderContentTypeConfigArgDefault]
	if t.ContentTypeType == config.FromTypeDefault && ok {
		t.ContentTypeStr = contentTypeStr
	}

	return nil
}

func (t *HeaderContentTypeTransformer) Index() int {
	return t.Cfg.Index
}

func (_ *HeaderContentTypeTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderContentTypeTransformer")

	contentType, ok := myMail.Metadata[HeaderContentTypeKey]
	if !ok {
		return myMail, nil
	}
	myMail.ContentType = contentType
	logger.Debug().
		Interface("contentType", myMail.ContentType).
		Msg("HeaderContentTypeTransformer")
	return myMail, nil
}
