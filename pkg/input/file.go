package input

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type FileReader struct {
	InPath       string
	PollInterval time.Duration
	ticker       *time.Ticker
}

func NewFileReader(
	_ context.Context,
	inPath string,
	pollInterval time.Duration,
) *FileReader {
	result := &FileReader{
		InPath:       inPath,
		PollInterval: pollInterval,
		ticker:       time.NewTicker(pollInterval),
	}

	return result
}

func (fr *FileReader) Run(
	ctx context.Context,
) error {
	defer fr.ticker.Stop()
	logger := zerolog.Ctx(ctx)

	// https://blog.devtrovert.com/p/select-and-for-range-channel-i-bet
outerloop:
	for {
		select {
		case t := <-fr.ticker.C:
			// check file exists
			logger.Info().Time("t", t).Msg("ticker.C")
		case <-ctx.Done():
			logger.Debug().Msg("ctx.Done")
			break outerloop
		}
	}
	return nil
}
