package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	slogzerolog "github.com/samber/slog-zerolog/v2"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

func InitLogger(ctx context.Context) (context.Context, *zerolog.Logger) {
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zerolog.FloatingPointPrecision = 2
	zerolog.ErrorFieldName = "e"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "t"

	wr := diode.NewWriter(
		os.Stderr,
		1000,
		10*time.Millisecond,
		func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	result := zerolog.New(wr).
		With().
		Timestamp().
		Caller().
		Logger()

	ctx = result.WithContext(ctx)

	_ = slog.New(
		slogzerolog.Option{
			Level:  slog.LevelInfo,
			Logger: &result,
		}.NewZerologHandler(),
	)
	slogzerolog.ErrorKeys = []string{"error", "err"}

	return ctx, &result
}
