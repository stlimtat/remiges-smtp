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
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

// DefaultFileName is the default format for output file names.
// The %s placeholder is replaced with a timestamp or mail ID based on configuration.
// Example: "output-2024-01-01.csv" for date-based naming
const (
	DefaultFileName string = "output-%s.csv"
)

// FileOutput implements the IOutput interface for writing mail processing results to CSV files.
// It supports different file naming strategies and writes mail processing results in a structured format.
//
// The implementation handles:
// - File creation and management
// - CSV formatting of mail processing results
// - Different file naming strategies (by date, hour, quarter-hour, or mail ID)
// - Concurrent access to output files
//
// Example usage:
//
//	output, err := NewFileOutput(ctx, config.OutputConfig{
//	    Type: config.ConfigOutputTypeFile,
//	    Args: map[string]any{
//	        config.ConfigArgPath: "/path/to/output",
//	        config.ConfigArgFileNameType: config.ConfigArgFileNameTypeDate,
//	    },
//	})
type FileOutput struct {
	// Cfg contains the output configuration specifying the output type and settings
	Cfg config.OutputConfig

	// FileNameType determines how output files are named:
	// - config.ConfigArgFileNameTypeDate: Files named by date (e.g., "output-2024-01-01.csv")
	// - config.ConfigArgFileNameTypeHour: Files named by hour (e.g., "output-2024-01-01-15.csv")
	// - config.ConfigArgFileNameTypeQuarterHour: Files named by quarter-hour (e.g., "output-2024-01-01-15-0.csv")
	// - config.ConfigArgFileNameTypeMailID: Files named by mail ID (e.g., "output-<mail-id>.csv")
	FileNameType string

	// Path is the directory where output files will be written.
	// The directory must exist and be writable.
	Path string
}

// NewFileOutput creates a new FileOutput instance with the provided configuration.
// It validates the output path and sets up the file naming strategy.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - cfg: Output configuration containing:
//   - Path: Directory for output files
//   - FileNameType: Strategy for naming output files
//
// Returns:
//   - *FileOutput: A new FileOutput instance
//   - error: Non-nil if:
//   - The path is missing from the configuration
//   - The path is invalid or inaccessible
//   - The file naming type is invalid
//
// Example:
//
//	cfg := config.OutputConfig{
//	    Type: config.ConfigOutputTypeFile,
//	    Args: map[string]any{
//	        config.ConfigArgPath: "/path/to/output",
//	        config.ConfigArgFileNameType: config.ConfigArgFileNameTypeDate,
//	    },
//	}
//	output, err := NewFileOutput(ctx, cfg)
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

	result.Path, err = utils.ValidateIO(ctx, result.Path, false, false)
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
//
// Examples:
//   - Date-based: "/path/to/output/output-2024-01-01.csv"
//   - Hour-based: "/path/to/output/output-2024-01-01-15.csv"
//   - Quarter-hour: "/path/to/output/output-2024-01-01-15-0.csv"
//   - Mail ID: "/path/to/output/output-<mail-id>.csv"
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

// InitAndWriteHeader initializes the output file and writes the CSV header if the file is new.
// It handles file creation and validation, ensuring the file is ready for writing.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - myMail: The mail being processed, used for file naming
//
// Returns:
//   - string: The full path to the output file
//   - error: Non-nil if:
//   - File validation fails
//   - File creation fails
//   - Header writing fails
func (f *FileOutput) InitAndWriteHeader(
	ctx context.Context,
	myMail *pmail.Mail,
) (string, error) {
	fileName := f.GetFileName(ctx, myMail)
	logger := zerolog.Ctx(ctx).
		With().
		Str("output.file", fileName).
		Logger()
	logger.Debug().Msg("FileOutput: InitAndWriteHeader")
	var err error

	fileName, err = utils.ValidateIO(ctx, fileName, true, true)
	if err == nil {
		return fileName, nil
	}
	if !strings.Contains(err.Error(), "ToIgnore") {
		logger.Error().Err(err).Msg("Failed to validate file name")
		return fileName, err
	}
	outputFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to create file")
		return fileName, err
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
		return fileName, err
	}

	return fileName, nil
}

// Write implements the IOutput interface by writing mail processing results to a CSV file.
// It creates a new CSV file (or appends to an existing one) and writes the mail ID,
// processing status, and any error messages.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - fileInfo: Information about the source file being processed
//   - myMail: The mail content being processed
//   - responses: Map of SMTP responses for different recipients
//
// Returns:
//   - error: Non-nil if:
//   - File initialization fails
//   - File opening fails
//   - Writing any response line fails
//   - Flushing the writer fails
//   - Closing the file fails
//
// The CSV output format is:
//
//	msg_id,status,error
//	<mail-id>,<status-code>,<response-line>
//
// Example output:
//
//	msg_id,status,error
//	abc123,250,250 2.0.0 OK
//	def456,550,550 5.1.1 User unknown
func (f *FileOutput) Write(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *pmail.Mail,
	responses map[string][]pmail.Response,
) error {
	fileName, err := f.InitAndWriteHeader(ctx, myMail)
	if err != nil {
		return err
	}
	logger := zerolog.Ctx(ctx).
		With().
		Str("output.file", fileName).
		Str("fileInfo.id", fileInfo.ID).
		Bytes("mail", myMail.MsgID).
		Logger()
	logger.Debug().Msg("FileOutput: Write")

	outputFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to open file")
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
