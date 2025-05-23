package file_mail

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	BodyTransformerType = "body"
)

type BodyTransformer struct {
	Cfg config.FileMailConfig
}

func (t *BodyTransformer) Init(
	ctx context.Context,
	cfg config.FileMailConfig,
) error {
	logger := zerolog.Ctx(ctx).With().
		Str("type", BodyTransformerType).
		Int("index", t.Cfg.Index).
		Interface("args", t.Cfg.Args).
		Logger()
	logger.Debug().Msg("BodyTransformer Init")
	t.Cfg = cfg
	return nil
}

func (t *BodyTransformer) Index() int {
	return t.Cfg.Index
}

func (_ *BodyTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Str("df_file_path", fileInfo.DfFilePath).
		Msg("BodyTransformer")
	var err error

	// 1. validate the df file exists and is readable
	if fileInfo.DfFilePath == "" {
		return nil, fmt.Errorf("ToSkip: DfFilePath is empty")
	}

	// 2. read all the bytes from the df file
	inMail.Body, err = os.ReadFile(fileInfo.DfFilePath)
	if err != nil {
		logger.Error().Err(err).Msg("os.ReadFile")
		return nil, err
	}

	// 3. Handling of unix new line to dos new line is done in mail Processor
	re := regexp.MustCompile(`\r?\n`)
	inMail.Body = re.ReplaceAll(inMail.Body, []byte("\r\n"))
	inMail.Body = bytes.TrimSpace(inMail.Body)

	return inMail, nil
}
