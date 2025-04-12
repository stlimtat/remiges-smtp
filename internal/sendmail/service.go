// Package sendmail provides functionality for sending emails via SMTP servers.
// It includes interfaces and implementations for dialing SMTP connections,
// handling email delivery, and managing the entire sending process.
package sendmail

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/file_mail"
	"github.com/stlimtat/remiges-smtp/internal/intmail"
	"github.com/stlimtat/remiges-smtp/internal/output"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

// SendMailService orchestrates the process of reading mail files,
// transforming them into mail objects, processing them, and sending
// them via SMTP. It manages concurrent processing of multiple files
// and handles the complete lifecycle of email delivery.
type SendMailService struct {
	// Concurrency specifies the number of concurrent mail processing goroutines
	Concurrency int

	// FileReader reads mail files from the filesystem
	FileReader file.IFileReader

	// MailProcessor handles mail processing tasks (e.g., DKIM signing)
	MailProcessor intmail.IMailProcessor

	// MailSender handles the actual SMTP delivery of emails
	MailSender IMailSender

	// MailTransformer converts file content into mail objects
	MailTransformer file_mail.IMailTransformer

	// MyOutput handles writing delivery results
	MyOutput output.IOutput

	// PollInterval specifies how often to check for new mail files
	PollInterval time.Duration

	// ticker is used for periodic file checking
	ticker *time.Ticker
}

// NewSendMailService creates a new SendMailService with the specified configuration.
// It initializes all necessary components for mail processing and delivery.
//
// Parameters:
//   - ctx: Context for service creation
//   - concurrency: Number of concurrent processing goroutines
//   - fileReader: Component for reading mail files
//   - mailProcessor: Component for processing mail (e.g., DKIM signing)
//   - mailSender: Component for SMTP delivery
//   - mailTransformer: Component for converting files to mail objects
//   - myOutput: Component for writing delivery results
//   - pollInterval: Interval between file checks
//
// Returns:
//   - *SendMailService: A new mail service instance
func NewSendMailService(
	_ context.Context,
	concurrency int,
	fileReader file.IFileReader,
	mailProcessor intmail.IMailProcessor,
	mailSender IMailSender,
	mailTransformer file_mail.IMailTransformer,
	myOutput output.IOutput,
	pollInterval time.Duration,
) *SendMailService {
	result := &SendMailService{
		Concurrency:     concurrency,
		FileReader:      fileReader,
		MailProcessor:   mailProcessor,
		MailSender:      mailSender,
		MailTransformer: mailTransformer,
		MyOutput:        myOutput,
		PollInterval:    pollInterval,
		ticker:          time.NewTicker(pollInterval),
	}
	return result
}

// Run starts the mail processing service. It launches multiple goroutines
// to handle concurrent processing of mail files and continuously checks
// for new files to process.
//
// Parameters:
//   - ctx: Context for controlling service lifecycle
//
// Returns:
//   - error: Any error that caused the service to stop
func (s *SendMailService) Run(
	ctx context.Context,
) error {
	defer s.ticker.Stop()
	logger := zerolog.Ctx(ctx)

	// Launch worker goroutines for concurrent processing
	var wg sync.WaitGroup
	for range s.Concurrency {
		wg.Add(1)
		go s.ProcessFileLoop(ctx, &wg)
	}

	// Main loop for checking new files
	// https://blog.devtrovert.com/p/select-and-for-range-channel-i-bet
outerloop:
	for {
		select {
		case t := <-s.ticker.C:
			// Check for new files at each tick
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

// ProcessFileLoop runs in a goroutine and continuously processes
// mail files as they become available. It handles the complete
// lifecycle of each mail file from reading to delivery.
//
// Parameters:
//   - ctx: Context for controlling the processing loop
//   - wg: WaitGroup for coordinating shutdown
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

// ReadNextMail processes a single mail file through the complete pipeline:
// reading the file, transforming it to a mail object, processing it,
// and sending it via SMTP.
//
// Parameters:
//   - ctx: Context for the processing operation
//
// Returns:
//   - *file.FileInfo: Information about the processed file
//   - *pmail.Mail: The processed mail object
//   - error: Any error encountered during processing
func (s *SendMailService) ReadNextMail(
	ctx context.Context,
) (*file.FileInfo, *pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)

	var fileInfo *file.FileInfo
	var myMail *pmail.Mail
	var err error

	// 1. Read the next available mail file
	// There is a mutex on the file reader to ensure that only one file is read at a time
	found := false
	// 2. Loop until a file is successfully processed
	for !found {
		fileInfo, err = s.FileReader.ReadNextFile(ctx)
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

		// Transform file content into a mail object
		myMail, err = s.MailTransformer.Transform(
			ctx, fileInfo, &pmail.Mail{},
		)
		if err != nil {
			if !strings.Contains(err.Error(), "ToIgnore") {
				return nil, nil, err
			}
		}
		fileInfo.Status = input.FILE_STATUS_BODY_READ

		// Process the mail (e.g., DKIM signing)
		myMail, err = s.MailProcessor.Process(ctx, myMail)
		if err != nil {
			if !strings.Contains(err.Error(), "ToIgnore") {
				return nil, nil, err
			}
		}
		fileInfo.Status = input.FILE_STATUS_MAIL_PROCESS

		// Send the mail via SMTP
		responses, errs := s.MailSender.SendMail(ctx, myMail)
		if errs != nil {
			return nil, nil, err
		}
		found = true

		// Log delivery results
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
	}

	return fileInfo, myMail, nil
}
