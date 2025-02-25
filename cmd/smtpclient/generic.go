package main

import (
	"log/slog"

	"github.com/mjl-/mox/dns"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/file_mail"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/sendmail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

type GenericSvc struct {
	Cfg                    config.SendMailConfig
	DialerFactory          sendmail.INetDialerFactory
	FileReader             file.IFileReader
	FileReadTracker        file.IFileReadTracker
	SendMailService        *sendmail.SendMailService
	MailProcessorFactory   *mail.DefaultMailProcessorFactory
	MailSender             sendmail.IMailSender
	MailTransformerFactory *file_mail.MailTransformerFactory
	RedisClient            *redis.Client
	Resolver               dns.Resolver
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

	result.SendMailService = sendmail.NewSendMailService(
		ctx,
		result.Cfg.ReadFileConfig.Concurrency,
		result.FileReader,
		result.MailProcessorFactory,
		result.MailSender,
		result.MailTransformerFactory,
		result.Cfg.ReadFileConfig.PollInterval,
	)
	result.MailProcessorFactory, err = mail.NewDefaultMailProcessorFactory(ctx, result.Cfg.MailProcessors)
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.MailProcessorFactory")
	}
	err = result.MailProcessorFactory.Init(ctx, config.MailProcessorConfig{})
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.MailProcessorFactory.Init")
	}
	result.Resolver = dns.StrictResolver{
		Log: result.Slogger,
	}
	result.MailSender = sendmail.NewMailSender(
		ctx,
		result.Cfg.Debug,
		result.DialerFactory,
		result.Resolver,
		result.Slogger,
	)
	return result
}
