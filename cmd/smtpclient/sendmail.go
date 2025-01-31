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
	"github.com/stlimtat/remiges-smtp/internal/sendmail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
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
	Cfg             config.SendMailConfig
	DialerFactory   sendmail.INetDialerFactory
	FileReader      input.IFileReader
	FileReadTracker input.IFileReadTracker
	MailSender      sendmail.IMailSender
	RedisClient     *redis.Client
	Resolver        dns.Resolver
	Slogger         *slog.Logger
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
	result.FileReadTracker = input.NewFileReadTracker(ctx, result.RedisClient)
	result.FileReader, err = input.NewDefaultFileReader(ctx, result.Cfg.InPath, result.FileReadTracker)
	if err != nil {
		logger.Fatal().Err(err).Msg("newSendMailSvc.FileReader")
	}
	result.Resolver = dns.StrictResolver{
		Log: result.Slogger,
	}
	result.MailSender = sendmail.NewMailSender(
		ctx,
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

	_, err = s.MailSender.NewConn(ctx, hosts)
	if err != nil {
		logger.Error().Err(err).Msg("MailSender.NewConn")
		return err
	}

	// read a file
	// files, err := s.FileReader.Process(ctx)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("FileReader.Process")
	// 	return err
	// }
	// if len(files) < 1 {
	// 	logger.Fatal().Msg("No files found")
	// }
	// msgBytes := files[0].BodyBytes
	// if !strings.HasSuffix(string(msgBytes), "\r\n") {
	// 	msgBytes = append(msgBytes, []byte("\r\n")...)
	// }

	// err = s.MailSender.SendMail(ctx, conn, s.Cfg.FromAddr, s.Cfg.ToAddr, msgBytes)
	// // do nothing. underlying has handled the error
	return err
}
