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
	HeaderFromTransformerType  = "header_from"
	HeaderFromKey              = "From"
	HeaderFromConfigArgType    = "type"
	HeaderFromConfigArgDefault = "default"
)

type HeaderFromTransformer struct {
	Cfg      config.FileMailConfig
	From     smtp.Address
	FromStr  string
	FromType config.FromType
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
	fromTypeStr, ok := t.Cfg.Args[HeaderFromConfigArgType]
	if !ok {
		fromTypeStr = config.FromTypeHeadersStr
	}
	switch fromTypeStr {
	case config.FromTypeDefaultStr:
		t.FromType = config.FromTypeDefault
	case config.FromTypeHeadersStr:
		t.FromType = config.FromTypeHeaders
	default:
		t.FromType = config.FromTypeHeaders
	}

	fromStr, ok := t.Cfg.Args[HeaderFromConfigArgDefault]
	if t.FromType == config.FromTypeDefault && ok {
		t.FromStr = fromStr
		t.From, err = smtp.ParseAddress(fromStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *HeaderFromTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderFromTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderFromTransformer")
	var err error

	from := t.From
	// 1. check if the header is present
	if t.FromType == config.FromTypeHeaders {
		fromValue, ok := myMail.Metadata[HeaderFromKey]
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
		from, err = smtp.ParseAddress(fromValueStr)
		if err != nil {
			return nil, err
		}
	}
	myMail.From = from
	logger.Debug().
		Interface("from", myMail.From).
		Msg("HeaderFromTransformer")
	return myMail, nil
}
