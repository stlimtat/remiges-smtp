package file_mail

import (
	"context"
	"fmt"

	"github.com/mcnijman/go-emailaddress"
	"github.com/mjl-/mox/smtp"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	HeaderToTransformerType  = "header_to"
	HeaderToConfigArgType    = "type"
	HeaderToConfigArgDefault = "default"
)

type HeaderToTransformer struct {
	Cfg    config.FileMailConfig
	To     []smtp.Address
	ToStr  string
	ToType config.ConfigType
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
	toTypeAny, ok := t.Cfg.Args[HeaderToConfigArgType]
	if !ok {
		toTypeAny = config.ConfigTypeHeadersStr
	}
	toTypeStr := toTypeAny.(string)
	switch toTypeStr {
	case config.ConfigTypeDefaultStr:
		t.ToType = config.ConfigTypeDefault
		toAny, ok := t.Cfg.Args[HeaderToConfigArgDefault]
		if !ok {
			toAny = ""
		}
		t.ToStr = toAny.(string)
		if t.ToStr != "" {
			toAddress, err := smtp.ParseAddress(t.ToStr)
			if err != nil {
				return err
			}
			t.To = []smtp.Address{toAddress}
		}
	default:
		t.ToType = config.ConfigTypeHeaders
	}

	return nil
}

func (t *HeaderToTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderToTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderToTransformer")
	if inMail == nil {
		inMail = &pmail.Mail{}
	}

	var result []smtp.Address

	switch t.ToType {
	case config.ConfigTypeDefault:
		result = t.To
	default:
		// Handling if the totype is headers
		headerTo, ok := inMail.Metadata[input.HeaderToKey]
		if !ok {
			return nil, fmt.Errorf("header %s not found", input.HeaderToKey)
		}
		emails := emailaddress.FindWithIcannSuffix(headerTo, false)

		result = make([]smtp.Address, 0)
		for _, email := range emails {
			emailStr := email.String()
			to, err := smtp.ParseAddress(emailStr)
			if err != nil {
				return nil, err
			}
			result = append(result, to)
		}
	}
	inMail.To = result
	logger.Debug().
		Interface(input.HeaderToKey, inMail.To).
		Msg("HeaderToTransformer")

	return inMail, nil
}
