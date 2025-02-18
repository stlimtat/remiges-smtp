package file_mail

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
)

const (
	BodyTransformerType = "body"
)

type BodyTransformer struct {
	Cfg config.FileMailConfig
}

func (t *BodyTransformer) Init(
	_ context.Context,
	cfg config.FileMailConfig,
) error {
	t.Cfg = cfg
	return nil
}

func (t *BodyTransformer) Index() int {
	return t.Cfg.Index
}

func (_ *BodyTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	inMail *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("BodyTransformer")
	var err error

	if fileInfo.DfReader == nil {
		logger.Error().Msg("DfReader is nil")
		return nil, fmt.Errorf("DfReader is nil")
	}

	inMail.Body, err = io.ReadAll(fileInfo.DfReader)
	if err != nil {
		return nil, err
	}
	// Handling of unix new line to dos new line is done in mail Processor

	return inMail, nil
}
