package sendmail

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/file_mail"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/pkg/input"
)

type SendMailService struct {
	Concurrency     int
	FileReader      file.IFileReader
	MailProcessor   mail.IMailProcessor
	MailSender      IMailSender
	MailTransformer file_mail.IMailTransformer
	PollInterval    time.Duration
	ticker          *time.Ticker
}

func NewSendMailService(
	_ context.Context,
	concurrency int,
	fileReader file.IFileReader,
	mailProcessor mail.IMailProcessor,
	mailSender IMailSender,
	mailTransformer file_mail.IMailTransformer,
	pollInterval time.Duration,
) *SendMailService {
	result := &SendMailService{
		Concurrency:     concurrency,
		FileReader:      fileReader,
		MailProcessor:   mailProcessor,
		MailSender:      mailSender,
		MailTransformer: mailTransformer,
		PollInterval:    pollInterval,
		ticker:          time.NewTicker(pollInterval),
	}
	return result
}

func (s *SendMailService) Run(
	ctx context.Context,
) error {
	defer s.ticker.Stop()
	logger := zerolog.Ctx(ctx)

	var wg sync.WaitGroup
	for range s.Concurrency {
		wg.Add(1)
		go s.ProcessFileLoop(ctx, &wg)
	}

	// https://blog.devtrovert.com/p/select-and-for-range-channel-i-bet
outerloop:
	for {
		select {
		case t := <-s.ticker.C:
			// check file exists
			logger.Info().Time("t", t).Msg("Run.ticker.C")
			_, err := s.FileReader.RefreshList(ctx)
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

func (s *SendMailService) ProcessFileLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	logger := zerolog.Ctx(ctx)
	defer wg.Done()

outerloop:
	for {
		select {
		case t := <-s.ticker.C:
			logger.Info().Time("t", t).Msg("ProcessFileLoop.ticker.C")
			fileInfo, _, err := s.ReadNextMail(ctx)
			if err != nil {
				continue
			}
			if fileInfo == nil {
				logger.Debug().Msg("no fileInfo found")
				continue
			}
			fileInfo.Status = input.FILE_STATUS_DONE
		case <-ctx.Done():
			logger.Debug().Msg("ctx.Done")
			break outerloop
		}
	}
}

func (s *SendMailService) ReadNextMail(
	ctx context.Context,
) (*file.FileInfo, *mail.Mail, error) {
	logger := zerolog.Ctx(ctx)

	fileInfo, err := s.FileReader.ReadNextFile(ctx)
	if err != nil {
		return nil, nil, err
	}
	if fileInfo == nil {
		return nil, nil, nil
	}
	logger.Debug().
		Str("fileInfo", fileInfo.ID).
		Msg("ReadNextFile")
	fileInfo.Status = input.FILE_STATUS_PROCESSING
	myMail, err := s.MailTransformer.Transform(
		ctx, fileInfo, &mail.Mail{},
	)
	if err != nil {
		return nil, nil, err
	}
	fileInfo.Status = input.FILE_STATUS_BODY_READ

	// process the mail
	myMail, err = s.MailProcessor.Process(ctx, myMail)
	if err != nil {
		return nil, nil, err
	}
	fileInfo.Status = input.FILE_STATUS_MAIL_PROCESS
	responses, errs := s.MailSender.SendMail(ctx, myMail)
	if errs != nil {
		return nil, nil, err
	}
	for to, response := range responses {
		if errs[to] != nil {
			logger.Error().
				AnErr("err", errs[to]).
				Interface("response", response).
				Str("to", to).
				Msg("Delivery failed")
			continue
		}
		logger.Info().
			Interface("response", response).
			Str("to", to).
			Msg("Delivery done")
	}
	fileInfo.Status = input.FILE_STATUS_DELIVERED

	return fileInfo, myMail, nil
}
