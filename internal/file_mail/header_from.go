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
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("HeaderFromTransformer Init")
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
	logger.Info().Msg("HeaderFromTransformer")
	var err error

	from := t.From
	// 1. check if the header is present
	if t.FromType == config.FromTypeHeaders {
		fromValue, ok := myMail.Headers[HeaderFromKey]
		if !ok {
			return myMail, nil
		}
		// 2. parse from value
		fromValueStr := string(fromValue)
		// 3. parse the header
		from, err = smtp.ParseAddress(fromValueStr)
		if err != nil {
			return nil, err
		}
	}
	myMail.From = from

	return myMail, nil
}
