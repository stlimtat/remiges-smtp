package file_mail

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
)

const (
	HeaderContentTypeTransformerType  = "header_contenttype"
	HeaderContentTypeKey              = "Content-Type"
	HeaderContentTypeConfigArgType    = "type"
	HeaderContentTypeConfigArgDefault = "default"
)

type HeaderContentTypeTransformer struct {
	Cfg             config.FileMailConfig
	ContentType     []byte
	ContentTypeStr  string
	ContentTypeType config.ConfigType
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
	contentTypeTypeAny, ok := t.Cfg.Args[HeaderConfigArgType]
	if !ok {
		contentTypeTypeAny = config.ConfigTypeHeadersStr
	}
	contentTypeTypeStr := contentTypeTypeAny.(string)
	switch contentTypeTypeStr {
	case config.ConfigTypeDefaultStr:
		t.ContentTypeType = config.ConfigTypeDefault
		contentTypeAny, ok := t.Cfg.Args[HeaderConfigArgDefault]
		if ok {
			t.ContentTypeStr = contentTypeAny.(string)
			t.ContentType = []byte(t.ContentTypeStr)
		}
	default:
		t.ContentTypeType = config.ConfigTypeHeaders
	}

	return nil
}

func (t *HeaderContentTypeTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderContentTypeTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderContentTypeTransformer")
	var contentType []byte

	switch t.ContentTypeType {
	case config.ConfigTypeDefault:
		contentType = t.ContentType
	default:
		var ok bool
		contentType, ok = myMail.Metadata[input.HeaderContentTypeKey]
		if !ok {
			contentType = make([]byte, 0)
		}
	}
	myMail.ContentType = contentType
	logger.Debug().
		Bytes(HeaderContentTypeKey, myMail.ContentType).
		Msg("HeaderContentTypeTransformer")
	return myMail, nil
}
