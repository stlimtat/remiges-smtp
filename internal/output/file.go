package output

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/sendmail"
)

const (
	DEFAULT_FILE_NAME    string = "remiges-smtp-%s.csv"
	ConfigOutputTypeFile string = "file"
	ConfigArgPath        string = "path"
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

	path, ok := cfg.Args[ConfigArgPath]
	if !ok {
		logger.Error().
			Msg("Path not found in config")
		return nil, fmt.Errorf("path not found in config")
	}
	result.Path = path.(string)
	if len(result.Path) < 1 {
		logger.Error().
			Msg("Path is empty")
		return nil, fmt.Errorf("path is empty")
	}
	if strings.HasPrefix(result.Path, "./") {
		wd, err := os.Getwd()
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to get working directory")
			return nil, fmt.Errorf("failed to get working directory")
		}
		result.Path = filepath.Join(wd, result.Path[2:])
	}
	if strings.HasPrefix(result.Path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to get home directory")
			return nil, fmt.Errorf("failed to get home directory")
		}
		result.Path = filepath.Join(home, result.Path[2:])
	}
	result.Path = filepath.Clean(result.Path)
	filePathInfo, err := os.Stat(result.Path)
	if os.IsNotExist(err) {
		logger.Error().
			Err(err).
			Msg("Path does not exist")
		return nil, fmt.Errorf("path does not exist")
	}
	if !filePathInfo.IsDir() {
		logger.Error().
			Msg("Path is not a directory")
		return nil, fmt.Errorf("path is not a directory")
	}

	return result, nil
}

func (f *FileOutput) Write(
	ctx context.Context,
	myMail *mail.Mail,
	resp []sendmail.Response,
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
