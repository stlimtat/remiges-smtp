// Package output provides functionality for writing mail processing results to various output destinations.
// It includes implementations for different output types (e.g., file, HTTP, etc.) and a factory
// for creating output instances based on configuration.
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
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

// DefaultFileName is the default format for output file names.
// The %s placeholder is replaced with a timestamp or mail ID based on configuration.
const (
	DefaultFileName string = "output-%s.csv"
)

// FileOutput implements the IOutput interface for writing mail processing results to CSV files.
// It supports different file naming strategies and writes mail processing results in a structured format.
type FileOutput struct {
	// Cfg contains the output configuration specifying the output type and settings
	Cfg config.OutputConfig

	// FileNameType determines how output files are named (e.g., by date, hour, mail ID)
	FileNameType string

	// Path is the directory where output files will be written
	Path string
}

// NewFileOutput creates a new FileOutput instance with the provided configuration.
// It validates the output path and sets up the file naming strategy.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - cfg: Output configuration containing path and file naming settings
//
// Returns:
//   - *FileOutput: A new FileOutput instance
//   - error: Non-nil if:
//   - The path is missing from the configuration
//   - The path is invalid or inaccessible
//   - The file naming type is invalid
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
		logger.Warn().
			Msg("Output file parameter FileNameType not found in config")
		fileNameType = config.ConfigArgFileNameTypeDate
	}
	result.FileNameType = fileNameType.(string)

	return result, nil
}

// GetFileName generates an output file name based on the configured naming strategy.
// The file name is constructed using the current time or mail ID, depending on the configuration.
//
// Parameters:
//   - ctx: Context for logging (currently unused)
//   - myMail: The mail being processed, used when naming by mail ID
//
// Returns:
//   - string: The full path to the output file
func (f *FileOutput) GetFileName(
	_ context.Context,
	myMail *pmail.Mail,
) string {
	var fileName string
	now := time.Now()
	switch f.FileNameType {
	case config.ConfigArgFileNameTypeMailID:
		fileName = fmt.Sprintf(DefaultFileName, myMail.MsgID)
	case config.ConfigArgFileNameTypeHour:
		hour := now.Format("2006-01-02-15")
		fileName = fmt.Sprintf(DefaultFileName, hour)
	case config.ConfigArgFileNameTypeQuarterHour:
		hour := now.Format("2006-01-02-15")
		minute := now.Minute()
		quarter := minute / 15
		hour = fmt.Sprintf("%s-%d", hour, quarter)
		fileName = fmt.Sprintf(DefaultFileName, hour)
	default:
		date := time.Now().Format("2006-01-02")
		fileName = fmt.Sprintf(DefaultFileName, date)
	}
	return filepath.Join(f.Path, fileName)
}

// Write implements the IOutput interface by writing mail processing results to a CSV file.
// It creates a new CSV file (or overwrites an existing one) and writes the mail ID,
// processing status, and any error messages.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - fileInfo: Information about the source file being processed
//   - myMail: The mail content being processed
//   - resp: The processing results and responses
//
// Returns:
//   - error: Non-nil if:
//   - The output file cannot be created
//   - Writing the CSV header fails
//   - Writing any response line fails
//   - Flushing the writer fails
//   - Closing the file fails
func (f *FileOutput) Write(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *pmail.Mail,
	responses map[string][]pmail.Response,
) error {
	fileName := f.GetFileName(ctx, myMail)
	logger := zerolog.Ctx(ctx).
		With().
		Str("output.file", fileName).
		Str("fileInfo.id", fileInfo.ID).
		Bytes("mail", myMail.MsgID).
		Logger()
	logger.Debug().Msg("FileOutput: Write")

	outputFile, err := os.Create(fileName)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to create file")
		return err
	}
	defer func() {
		err = outputFile.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to close file")
		}
	}()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	err = writer.Write([]string{"msg_id", "status", "error"})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write header")
		return err
	}
	for _, resp := range responses {
		for _, r := range resp {
			err = writer.Write([]string{
				string(myMail.MsgID),
				fmt.Sprintf("%d", r.Code),
				r.Line,
			})
			if err != nil {
				logger.Error().Err(err).Msg("Failed to write line")
				return err
			}
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
