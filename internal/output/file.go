package output

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	DEFAULT_FILE_NAME string = "output-%s.csv"
)

type FileOutput struct {
	Cfg          config.OutputConfig
	FileNameType string
	Path         string
}

func NewFileOutput(
	ctx context.Context,
	cfg config.OutputConfig,
) (*FileOutput, error) {
	logger := zerolog.Ctx(ctx).With().Interface("cfg", cfg).Logger()
	var err error

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

	result.Path, err = utils.ValidateIO(ctx, result.Path, false)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to validate path")
		return nil, err
	}

	fileNameType, ok := cfg.Args[config.ConfigArgFileNameType]
	if !ok {
		logger.Error().
			Msg("FileNameType not found in config")
		fileNameType = config.ConfigArgFileNameTypeDate
	}
	result.FileNameType = fileNameType.(string)

	return result, nil
}

func (f *FileOutput) GetFileName(
	_ context.Context,
	myMail *pmail.Mail,
) string {
	var fileName string
	now := time.Now()
	switch f.FileNameType {
	case config.ConfigArgFileNameTypeMailID:
		fileName = fmt.Sprintf(DEFAULT_FILE_NAME, myMail.MsgID)
	case config.ConfigArgFileNameTypeHour:
		hour := now.Format("2006-01-02-15")
		fileName = fmt.Sprintf(DEFAULT_FILE_NAME, hour)
	case config.ConfigArgFileNameTypeQuarterHour:
		hour := now.Format("2006-01-02-15")
		minute := now.Minute()
		quarter := minute / 15
		hour = fmt.Sprintf("%s-%d", hour, quarter)
		fileName = fmt.Sprintf(DEFAULT_FILE_NAME, hour)
	default:
		date := time.Now().Format("2006-01-02")
		fileName = fmt.Sprintf(DEFAULT_FILE_NAME, date)
	}
	return filepath.Join(f.Path, fileName)
}

func (f *FileOutput) Write(
	ctx context.Context,
	myMail *pmail.Mail,
	resp []pmail.Response,
) error {
	fileName := f.GetFileName(ctx, myMail)
	logger := zerolog.Ctx(ctx).
		With().
		Str("file", fileName).
		Bytes("mail", myMail.MsgID).
		Logger()
	logger.Debug().Msg("FileOutput: Write")

	file, err := os.Create(fileName)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to create file")
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to close file")
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

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

	err = writer.Error()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to flush writer")
		return err
	}
	logger.Info().Msg("FileOutput: Write success")
	return nil
}
