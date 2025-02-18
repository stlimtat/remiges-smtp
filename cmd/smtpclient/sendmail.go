/*
Copyright © 2024 Lim Swee Tat <st_lim@stlim.net>
*/
package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mjl-/mox/dns"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/file_mail"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/sendmail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

type sendMailCmd struct {
	cmd *cobra.Command
}

func newSendMailCmd(ctx context.Context) (*sendMailCmd, *cobra.Command) {
	logger := zerolog.Ctx(ctx)
	var err error

	result := &sendMailCmd{}

	// sendMailCmd represents the server command
	result.cmd = &cobra.Command{
		Use:   "sendmail",
		Short: "Send a mail from a sender email, to a destination email, with a test message",
		Long:  `Runs the smtp client which will run sendMail`,
		Args: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			cmdLogger := zerolog.Ctx(ctx)
			cfg := config.NewSendMailConfig(ctx)
			if len(cfg.From) < 1 || len(cfg.To) < 1 {
				cmdLogger.Fatal().
					Err(fmt.Errorf("missing fields")).
					Interface("cfg", cfg).
					Msg("Missing fields")
			}
			ctx = config.SetContextConfig(ctx, cfg)
			cmd.SetContext(ctx)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sendMailSvc := newSendMailSvc(cmd, args)
			err = sendMailSvc.Run(cmd, args)
			return err
		},
	}

	result.cmd.Flags().StringP(
		"from", "f",
		"", "Send the test message from - sender email address",
	)
	result.cmd.Flags().StringP(
		"path", "p",
		"", "Path to the directory containing the df and qf files",
	)
	result.cmd.Flags().StringP(
		"to", "t",
		"", "Send the test message to - destination email address",
	)
	result.cmd.Flags().StringP(
		"msg", "m",
		"", "Test message",
	)
	result.cmd.Flags().StringP(
		"redis_addr", "r",
		"", "Redis address",
	)
	err = viper.BindPFlag("from", result.cmd.Flags().Lookup("from"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag.from")
	}
	err = viper.BindPFlag("in_path", result.cmd.Flags().Lookup("path"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag.inpath")
	}
	err = viper.BindPFlag("to", result.cmd.Flags().Lookup("to"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag.to")
	}
	err = viper.BindPFlag("msg", result.cmd.Flags().Lookup("msg"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag.msg")
	}

	return result, result.cmd
}

type SendMailSvc struct {
	Cfg                    config.SendMailConfig
	DialerFactory          sendmail.INetDialerFactory
	FileReader             file.IFileReader
	FileReadTracker        file.IFileReadTracker
	FileMailService        *file_mail.FileMailService
	MailProcessorFactory   *mail.DefaultMailProcessorFactory
	MailSender             sendmail.IMailSender
	MailTransformerFactory *file_mail.MailTransformerFactory
	RedisClient            *redis.Client
	Resolver               dns.Resolver
	Slogger                *slog.Logger
}

func newSendMailSvc(
	cmd *cobra.Command,
	_ []string,
) *SendMailSvc {
	var err error
	result := &SendMailSvc{}
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

	result.FileMailService = file_mail.NewFileMailService(
		ctx,
		result.Cfg.ReadFileConfig.Concurrency,
		result.FileReader,
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

func (s *SendMailSvc) Run(
	cmd *cobra.Command,
	_ []string,
) error {
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)
	var err error

	hosts, err := s.MailSender.LookupMX(ctx, s.Cfg.ToAddr.Domain)
	if err != nil {
		logger.Error().Err(err).Msg("MailSender.LookupMX")
		return err
	}

	// refresh file list
	_, err = s.FileReader.RefreshList(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("FileMailService.RefreshList")
		return err
	}

	// read a file
	fileInfo, myMail, err := s.FileMailService.ReadNextMail(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("FileMailService.ReadNextMail")
		return err
	}
	logger.Info().
		Str("fileInfo", fileInfo.ID).
		Str("from", myMail.From.String()).
		Msg("ReadNextMail")

	// process the mail
	myMail, err = s.MailProcessorFactory.Process(ctx, myMail)
	if err != nil {
		logger.Error().Err(err).Msg("MailProcessorFactory.Process")
		return err
	}

	conn, err := s.MailSender.NewConn(ctx, hosts)
	if err != nil {
		logger.Error().Err(err).Msg("MailSender.NewConn")
		return err
	}

	err = s.MailSender.SendMail(ctx, conn, myMail)
	if err != nil {
		logger.Error().Err(err).Msg("MailSender.SendMail")
		return err
	}

	return nil
}
