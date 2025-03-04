package main

import (
	"log/slog"
	"os"
	"reflect"

	moxDns "github.com/mjl-/mox/dns"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/crypto"
	"github.com/stlimtat/remiges-smtp/internal/dns"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/file_mail"
	"github.com/stlimtat/remiges-smtp/internal/intmail"
	"github.com/stlimtat/remiges-smtp/internal/output"
	"github.com/stlimtat/remiges-smtp/internal/sendmail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

type GenericSvc struct {
	Cfg                    config.SendMailConfig
	CryptoFactory          *crypto.CryptoFactory
	DialerFactory          sendmail.INetDialerFactory
	FileReader             file.IFileReader
	FileReadTracker        file.IFileReadTracker
	KeyWriter              crypto.IKeyWriter
	MailProcessor          intmail.IMailProcessor
	MailSender             sendmail.IMailSender
	MailTransformerFactory *file_mail.MailTransformerFactory
	MoxResolver            moxDns.Resolver
	MyOutput               output.IOutput
	MyResolver             dns.IResolver
	RedisClient            *redis.Client
	SendMailService        *sendmail.SendMailService
	Slogger                *slog.Logger
}

func newGenericSvc(
	cmd *cobra.Command,
	_ []string,
) *GenericSvc {
	var err error
	result := &GenericSvc{}
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)
	result.Cfg = config.GetContextConfig(ctx).(config.SendMailConfig)
	result.Slogger = telemetry.GetSLogger(ctx)
	result.DialerFactory = sendmail.NewDefaultDialerFactory()
	result.RedisClient = redis.NewClient(&redis.Options{
		Addr: result.Cfg.ReadFileConfig.RedisAddr,
	})
	// Run a ping on the redis client to check if it's working
	_, err = result.RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.RedisClient.Ping")
	}
	result.FileReadTracker = file.NewFileReadTracker(
		ctx,
		result.RedisClient,
	)
	result.FileReader, err = file.NewDefaultFileReader(
		ctx,
		result.Cfg.ReadFileConfig.InPath,
		result.FileReadTracker,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.FileReader")
	}
	result.MailTransformerFactory = file_mail.NewMailTransformerFactory(
		ctx,
		result.Cfg.ReadFileConfig.FileMails,
	)
	err = result.MailTransformerFactory.Init(ctx, config.FileMailConfig{})
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.MailTransformerFactory.Init")
	}
	outputFactory := &output.OutputFactory{}
	_, err = outputFactory.NewOutputs(ctx, result.Cfg.Outputs)
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.MyOutput")
	}
	result.MyOutput = outputFactory

	mailProcessorFactory, err := intmail.NewDefaultMailProcessorFactory(ctx, result.Cfg.MailProcessors)
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.MailProcessorFactory")
	}
	_, err = mailProcessorFactory.NewMailProcessors(ctx, result.Cfg.MailProcessors)
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.MailProcessorFactory.Init")
	}
	result.MailProcessor = mailProcessorFactory

	result.MoxResolver = moxDns.StrictResolver{
		Log: result.Slogger,
	}
	result.MyResolver = dns.NewResolver(
		ctx,
		result.MoxResolver,
		result.Slogger,
	)
	result.MailSender = sendmail.NewMailSender(
		ctx,
		result.Cfg.Debug,
		result.DialerFactory,
		result.MyResolver,
		result.Slogger,
	)

	result.SendMailService = sendmail.NewSendMailService(
		ctx,
		result.Cfg.ReadFileConfig.Concurrency,
		result.FileReader,
		result.MailProcessor,
		result.MailSender,
		result.MailTransformerFactory,
		result.MyOutput,
		result.Cfg.ReadFileConfig.PollInterval,
	)

	// We need to init the crypto factory
	// for generic, we are not needing to write to a file
	// so we can just use a temp directory
	tempDir, err := os.MkdirTemp("", "remiges-smtp")
	if err != nil {
		logger.Fatal().Err(err).Msg("newGenericSvc.MkdirTemp")
	}
	result.CryptoFactory = &crypto.CryptoFactory{}
	result.KeyWriter = crypto.NewKeyWriter(ctx, tempDir)
	_, err = result.CryptoFactory.Init(ctx, result.KeyWriter)
	if err != nil {
		logger.Fatal().Err(err).Msg("newGenericSvc.CryptoFactory.Init")
	}
	// This is a hack to inject the crypto factory into the dkim processor
	for _, mailProcessor := range mailProcessorFactory.Processors {
		if reflect.TypeOf(mailProcessor) == reflect.TypeOf(&intmail.DKIMProcessor{}) {
			dkimProcessor := mailProcessor.(*intmail.DKIMProcessor)
			err = dkimProcessor.InitDKIMCrypto(ctx, result.CryptoFactory)
			if err != nil {
				logger.Fatal().Err(err).Msg("newGenericSvc.DKIMProcessor.InitDKIMCrypto")
			}
		}
	}
	return result
}
