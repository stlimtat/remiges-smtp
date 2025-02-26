package main

import (
	"log/slog"

	moxDns "github.com/mjl-/mox/dns"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
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
	DialerFactory          sendmail.INetDialerFactory
	FileReader             file.IFileReader
	FileReadTracker        file.IFileReadTracker
	SendMailService        *sendmail.SendMailService
	MailProcessor          intmail.IMailProcessor
	MailSender             sendmail.IMailSender
	MailTransformerFactory *file_mail.MailTransformerFactory
	MyOutput               output.IOutput
	RedisClient            *redis.Client
	MoxResolver            moxDns.Resolver
	MyResolver             dns.IResolver
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

	result.MailProcessor, err = intmail.NewDefaultMailProcessorFactory(ctx, result.Cfg.MailProcessors)
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.MailProcessorFactory")
	}
	err = result.MailProcessor.Init(ctx, config.MailProcessorConfig{})
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.MailProcessorFactory.Init")
	}

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
	return result
}
