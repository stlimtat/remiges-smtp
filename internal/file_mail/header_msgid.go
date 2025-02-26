package file_mail

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
)

type MsgIDType int

const (
	HeaderMsgIDTransformerType = "header_msgid"
	HeaderMsgIDConfigArgUuid   = "uuid"

	MsgIDTypeHeaders MsgIDType = 0
	MsgIDTypeDefault MsgIDType = 1
	MsgIDTypeUuid    MsgIDType = 2
)

type HeaderMsgIDTransformer struct {
	Cfg       config.FileMailConfig
	MsgID     string
	MsgIDStr  string
	MsgIDType MsgIDType
}

func (t *HeaderMsgIDTransformer) Init(
	ctx context.Context,
	cfg config.FileMailConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", HeaderMsgIDTransformerType).
		Int("index", t.Cfg.Index).
		Interface("args", t.Cfg.Args).
		Logger()
	logger.Debug().Msg("HeaderMsgIDTransformer Init")

	t.Cfg = cfg
	msgIDTypeAny, ok := t.Cfg.Args[HeaderConfigArgType]
	if !ok {
		msgIDTypeAny = config.ConfigTypeHeadersStr
	}
	msgIDTypeStr := msgIDTypeAny.(string)
	switch msgIDTypeStr {
	case HeaderConfigArgDefault:
		t.MsgIDType = MsgIDTypeDefault
		msgIDAny, ok := t.Cfg.Args[HeaderConfigArgDefault]
		if ok {
			t.MsgIDStr = msgIDAny.(string)
		}
	case HeaderMsgIDConfigArgUuid:
		t.MsgIDType = MsgIDTypeUuid
	default:
		t.MsgIDType = MsgIDTypeHeaders
	}

	return nil
}

func (t *HeaderMsgIDTransformer) Index() int {
	return t.Cfg.Index
}

func (t *HeaderMsgIDTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("id", fileInfo.ID).
		Logger()
	logger.Debug().Msg("HeaderMsgIDTransformer")

	var msgID []byte
	switch t.MsgIDType {
	case MsgIDTypeDefault:
		msgID = []byte(t.MsgIDStr)
	case MsgIDTypeUuid:
		msgID = t.GetMsgID(ctx)
	default:
		var ok bool
		msgID, ok = myMail.Metadata[input.HeaderMsgIDKey]
		if !ok {
			msgID = t.GetMsgID(ctx)
		}
	}

	myMail.MsgID = msgID
	logger.Debug().
		Interface(input.HeaderMsgIDKey, myMail.MsgID).
		Msg("HeaderMsgIDTransformer")
	return myMail, nil
}

func (_ *HeaderMsgIDTransformer) GetMsgID(
	ctx context.Context,
) []byte {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("GetMsgID")
	var result []byte
	for {
		rawUuid, err := uuid.NewV7()
		if err != nil {
			logger.Error().Err(err).Msg("uuid.NewV7")
			continue
		}
		result = rawUuid[:]
		break
	}
	return result
}
