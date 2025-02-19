package file_mail

import (
	"context"
	"fmt"

	"github.com/mcnijman/go-emailaddress"
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
	To     []smtp.Address
	ToStr  string
	ToType config.FromType
}

func (t *HeaderToTransformer) Init(
	ctx context.Context,
	cfg config.FileMailConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", HeaderToTransformerType).
		Int("index", t.Cfg.Index).
		Interface("args", t.Cfg.Args).
		Logger()
	logger.Debug().Msg("HeaderToTransformer Init")

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
		toAddress, err := smtp.ParseAddress(t.ToStr)
		if err != nil {
			return err
		}
		t.To = []smtp.Address{toAddress}
	}

	return nil
}

func (t *HeaderToTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderToTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderToTransformer")

	if t.ToType == config.FromTypeDefault {
		inMail.To = t.To
		return inMail, nil
	}
	// Handling if the totype is headers
	headerTo, ok := inMail.Metadata[HeaderToKey]
	if !ok {
		return nil, fmt.Errorf("header %s not found", HeaderToKey)
	}
	emails := emailaddress.FindWithIcannSuffix(headerTo, false)

	inMail.To = make([]smtp.Address, 0)
	for _, email := range emails {
		emailStr := email.String()
		to, err := smtp.ParseAddress(emailStr)
		if err != nil {
			return nil, err
		}
		inMail.To = append(inMail.To, to)
	}
	logger.Debug().
		Interface("to", inMail.To).
		Msg("HeaderToTransformer")

	return inMail, nil
}
