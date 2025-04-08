// Package telemetry provides logging and observability functionality for the application.
// It uses zerolog as the primary logging library and provides integration with the
// standard library's slog package. The package includes functions for logger initialization,
// configuration, and retrieval.
package telemetry

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	slogzerolog "github.com/samber/slog-zerolog/v2"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes a new zerolog logger with default configuration and returns
// a context with the logger embedded. The logger is configured to write to os.Stdout
// with a diode writer to prevent blocking during high load.
//
// Parameters:
//   - ctx: The context to embed the logger in
//
// Returns:
//   - context.Context: A new context with the logger embedded
//   - *zerolog.Logger: The initialized logger instance
func InitLogger(ctx context.Context) (context.Context, *zerolog.Logger) {
	return GetLogger(ctx, os.Stdout)
}

// GetLogger creates and configures a new zerolog logger with the specified writer.
// The logger is configured with:
//   - Caller information (file and line number)
//   - Timestamp in Unix format
//   - Custom field names for better log parsing
//   - Diode writer for non-blocking writes
//   - Integration with slog for compatibility
//
// Parameters:
//   - ctx: The context to embed the logger in
//   - writer: The io.Writer to write logs to
//
// Returns:
//   - context.Context: A new context with the logger embedded
//   - *zerolog.Logger: The configured logger instance
func GetLogger(ctx context.Context, writer io.Writer) (context.Context, *zerolog.Logger) {
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zerolog.FloatingPointPrecision = 2
	zerolog.ErrorFieldName = "e"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "t"

	wr := diode.NewWriter(
		writer,
		1000,
		10*time.Millisecond,
		func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})

	result := zerolog.New(wr).
		With().
		Timestamp().
		Caller().
		Logger()

	ctx = result.WithContext(ctx)
	log.Logger = result

	_ = slog.New(
		slogzerolog.Option{
			Level:  slog.LevelInfo,
			Logger: &result,
		}.NewZerologHandler(),
	)
	slogzerolog.ErrorKeys = []string{"error", "err"}

	return ctx, &result
}

// SetGlobalLogLevel sets the global logging level for all zerolog loggers.
// This affects the minimum level of logs that will be output.
//
// Parameters:
//   - level: The zerolog.Level to set as the global logging level
func SetGlobalLogLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

// GetSLogger retrieves a slog.Logger instance from the context, configured to use
// the embedded zerolog logger. This provides compatibility with the standard
// library's logging interface.
//
// Parameters:
//   - ctx: The context containing the zerolog logger
//
// Returns:
//   - *slog.Logger: A slog logger instance that writes to the embedded zerolog logger
func GetSLogger(ctx context.Context) *slog.Logger {
	logger := zerolog.Ctx(ctx)
	result := slog.New(
		slogzerolog.Option{
			Level:  slog.LevelInfo,
			Logger: logger,
		}.NewZerologHandler(),
	)
	return result
}
