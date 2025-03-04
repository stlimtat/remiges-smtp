package output

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	DEFAULT_FILE_NAME string = "remiges-smtp-%s.csv"
)

type FileOutput struct {
	Cfg  config.OutputConfig
	Path string
}

func NewFileOutput(
	ctx context.Context,
	cfg config.OutputConfig,
) (*FileOutput, error) {
	logger := zerolog.Ctx(ctx).With().Interface("cfg", cfg).Logger()

	result := &FileOutput{
		Cfg: cfg,
	}

	path, ok := cfg.Args[config.ConfigArgPath]
	if !ok {
		logger.Error().
			Msg("Path not found in config")
		return nil, fmt.Errorf("path not found in config")
	}
	result.Path = path.(string)

	err := utils.ValidateIO(ctx, result.Path, false)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to validate path")
		return nil, err
	}

	return result, nil
}

func (f *FileOutput) Write(
	ctx context.Context,
	myMail *pmail.Mail,
	resp []pmail.Response,
) error {
	filePath := filepath.Join(f.Path, fmt.Sprintf(DEFAULT_FILE_NAME, myMail.MsgID))
	logger := zerolog.Ctx(ctx).
		With().
		Str("file", filePath).
		Bytes("mail", myMail.MsgID).
		Logger()

	file, err := os.Create(filePath)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to create file")
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to close file")
		}
	}()

	writer := csv.NewWriter(file)
	err = writer.Write([]string{"msg_id", "status", "error"})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write header")
		return err
	}
	for _, resp := range resp {
		err = writer.Write([]string{
			string(myMail.MsgID),
			fmt.Sprintf("%d", resp.Code),
			resp.Line,
		})
		if err != nil {
			logger.Error().Err(err).Msg("Failed to write line")
			return err
		}
	}
	writer.Flush()
	err = writer.Error()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to flush writer")
		return err
	}
	return nil
}
