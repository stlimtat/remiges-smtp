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
	HeaderFromTransformerType = "header_from"
)

type HeaderFromTransformer struct {
	Cfg      config.FileMailConfig
	From     smtp.Address
	FromStr  string
	FromType config.ConfigType
}

func (t *HeaderFromTransformer) Init(
	ctx context.Context,
	cfg config.FileMailConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", HeaderFromTransformerType).
		Int("index", t.Cfg.Index).
		Interface("args", t.Cfg.Args).
		Logger()
	logger.Debug().Msg("HeaderFromTransformer Init")
	var err error

	t.Cfg = cfg
	fromTypeAny, ok := t.Cfg.Args[HeaderConfigArgType]
	if !ok {
		fromTypeAny = config.ConfigTypeHeadersStr
	}
	fromTypeStr := fromTypeAny.(string)
	switch fromTypeStr {
	case config.ConfigTypeDefaultStr:
		t.FromType = config.ConfigTypeDefault
		fromAny, ok := t.Cfg.Args[HeaderConfigArgDefault]
		if ok {
			t.FromStr = fromAny.(string)
			t.From, err = smtp.ParseAddress(t.FromStr)
			if err != nil {
				return err
			}
		}
	default:
		t.FromType = config.ConfigTypeHeaders
	}

	return nil
}

func (t *HeaderFromTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderFromTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderFromTransformer")
	var err error
	var result smtp.Address
	switch t.FromType {
	case config.ConfigTypeDefault:
		result = t.From
	default:
		// 1. check if the header is present
		fromValue, ok := myMail.Metadata[input.HeaderFromKey]
		if !ok {
			return myMail, nil
		}
		logger.Debug().Bytes("fromValue", fromValue).Msg("HeaderFromTransformer")
		// 2. parse from value
		// We have "From: Name of user <from@example.com>"
		emails := emailaddress.FindWithIcannSuffix(fromValue, false)
		if len(emails) == 0 {
			return nil, fmt.Errorf("no email address found in from header")
		}
		var fromValueStr string
		for _, email := range emails {
			fromValueStr = email.String()
		}
		// 3. parse the header
		result, err = smtp.ParseAddress(fromValueStr)
		if err != nil {
			return nil, err
		}
	}
	myMail.From = result
	logger.Debug().
		Interface(input.HeaderFromKey, myMail.From).
		Msg("HeaderFromTransformer")
	return myMail, nil
}
