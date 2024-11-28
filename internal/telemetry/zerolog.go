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

func InitLogger(ctx context.Context) (context.Context, *zerolog.Logger) {
	return GetLogger(ctx, os.Stdout)
}

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

func SetGlobalLogLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

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
