package file_mail

import (
	"context"

	"github.com/mjl-/mox/smtp"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
)

const (
	HeaderToTransformerType  = "header_to"
	HeaderToKey              = "To"
	HeaderToConfigArgType    = "type"
	HeaderToConfigArgDefault = "default"
)

type HeaderToTransformer struct {
	Cfg    config.FileMailConfig
	To     smtp.Address
	ToStr  string
	ToType config.FromType
}

func (t *HeaderToTransformer) Init(
	ctx context.Context,
	cfg config.FileMailConfig,
) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HeaderToTransformer Init")
	var err error

	t.Cfg = cfg
	toType := t.Cfg.Args[HeaderToConfigArgType]
	if toType == "" {
		toType = config.FromTypeHeadersStr
	}
	switch toType {
	case config.FromTypeDefaultStr:
		t.ToType = config.FromTypeDefault
	case config.FromTypeHeadersStr:
		t.ToType = config.FromTypeHeaders
	default:
		t.ToType = config.FromTypeHeaders
	}
	toStr, ok := t.Cfg.Args[HeaderToConfigArgDefault]
	if !ok {
		toStr = ""
	}
	t.ToStr = toStr
	if t.ToStr != "" {
		t.To, err = smtp.ParseAddress(t.ToStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *HeaderToTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderToTransformer) Transform(
	ctx context.Context,
	_ *file.FileInfo,
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HeaderToTransformer Transform")
	var err error

	if t.ToType == config.FromTypeDefault {
		inMail.To = t.To
		return inMail, nil
	}
	// Handling if the totype is headers
	headerTo, ok := inMail.Headers[HeaderToKey]
	if ok {
		headerToStr := string(headerTo)
		inMail.To, err = smtp.ParseAddress(headerToStr)
		if err != nil {
			return nil, err
		}
	}

	return inMail, nil
}
