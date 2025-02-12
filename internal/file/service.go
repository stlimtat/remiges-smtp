package file

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/pkg/input"
)

type FileService struct {
	Concurrency     int
	FileReader      IFileReader
	MailTransformer IMailTransformer
	PollInterval    time.Duration
	ticker          *time.Ticker
}

func NewFileService(
	_ context.Context,
	concurrency int,
	fileReader IFileReader,
	mailTransformer IMailTransformer,
	pollInterval time.Duration,
) *FileService {
	result := &FileService{
		Concurrency:     concurrency,
		FileReader:      fileReader,
		MailTransformer: mailTransformer,
		PollInterval:    pollInterval,
		ticker:          time.NewTicker(pollInterval),
	}

	return result
}

func (fs *FileService) Run(
	ctx context.Context,
) error {
	defer fs.ticker.Stop()
	logger := zerolog.Ctx(ctx)

	var wg sync.WaitGroup
	for range fs.Concurrency {
		wg.Add(1)
		go fs.ProcessFileLoop(ctx, &wg)
	}

	// https://blog.devtrovert.com/p/select-and-for-range-channel-i-bet
outerloop:
	for {
		select {
		case t := <-fs.ticker.C:
			// check file exists
			logger.Info().Time("t", t).Msg("Run.ticker.C")
			_, err := fs.FileReader.RefreshList(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("RefreshList")
				continue
			}
		case <-ctx.Done():
			logger.Debug().Msg("ctx.Done")
			break outerloop
		}
	}
	wg.Wait()
	return nil
}

func (fs *FileService) ProcessFileLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	logger := zerolog.Ctx(ctx)
	defer wg.Done()

outerloop:
	for {
		select {
		case t := <-fs.ticker.C:
			logger.Info().Time("t", t).Msg("ProcessFileLoop.ticker.C")
			fileInfo, _, err := fs.ReadNextMail(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("ProcessFile")
				continue
			}
			if fileInfo == nil {
				logger.Debug().Msg("no fileInfo found")
				continue
			}
			fileInfo.Status = input.FILE_STATUS_BODY_READ
		case <-ctx.Done():
			logger.Debug().Msg("ctx.Done")
			break outerloop
		}
	}
}

func (fs *FileService) ReadNextMail(
	ctx context.Context,
) (*FileInfo, *mail.Mail, error) {
	logger := zerolog.Ctx(ctx)

	fileInfo, err := fs.FileReader.ReadNextFile(ctx)
	if err != nil {
		return nil, nil, err
	}
	if fileInfo == nil {
		logger.Info().Msg("no fileInfo found")
		return nil, nil, nil
	}
	logger.Info().
		Str("fileInfo", fileInfo.ID).
		Msg("ReadNextFile")
	fileInfo.Status = input.FILE_STATUS_PROCESSING
	myMail, err := fs.MailTransformer.Transform(
		ctx, fileInfo,
	)
	if err != nil {
		return nil, nil, err
	}
	fileInfo.Status = input.FILE_STATUS_BODY_READ

	return fileInfo, myMail, nil
}
