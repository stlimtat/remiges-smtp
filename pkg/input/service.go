package input

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type FileService struct {
	FileReader   *FileReader
	PollInterval time.Duration
	ticker       *time.Ticker
}

func NewFileService(
	_ context.Context,
	fileReader *FileReader,
	pollInterval time.Duration,
) *FileService {
	result := &FileService{
		FileReader:   fileReader,
		PollInterval: pollInterval,
		ticker:       time.NewTicker(pollInterval),
	}

	return result
}

func (fs *FileService) Run(
	ctx context.Context,
) error {
	defer fs.ticker.Stop()
	logger := zerolog.Ctx(ctx)

	// https://blog.devtrovert.com/p/select-and-for-range-channel-i-bet
outerloop:
	for {
		select {
		case t := <-fs.ticker.C:
			// check file exists
			logger.Info().Time("t", t).Msg("ticker.C")
		case <-ctx.Done():
			logger.Debug().Msg("ctx.Done")
			break outerloop
		}
	}
	return nil
}
