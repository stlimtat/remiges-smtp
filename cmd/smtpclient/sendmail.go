/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mjl-/mox/dns"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/sendmail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"golang.org/x/sync/errgroup"
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
		"to", "t",
		"", "Send the test message to - destination email address",
	)
	result.cmd.Flags().StringP(
		"msg", "m",
		"Testing", "Test message",
	)
	err = viper.BindPFlag("from", result.cmd.Flags().Lookup("from"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag.from")
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
	Cfg           config.SendMailConfig
	DialerFactory sendmail.INetDialerFactory
	MailSender    sendmail.IMailSender
	Resolver      sendmail.IResolver
	Slogger       *slog.Logger
}

func newSendMailSvc(
	cmd *cobra.Command,
	_ []string,
) *SendMailSvc {
	result := &SendMailSvc{}
	ctx := cmd.Context()
	cfg := config.GetContextConfig(ctx).(config.SendMailConfig)
	result.Cfg = cfg
	result.Slogger = telemetry.GetSLogger(ctx)
	result.DialerFactory = sendmail.NewDefaultDialerFactory()
	result.Resolver = dns.StrictResolver{
		Log: result.Slogger,
	}
	result.MailSender = sendmail.NewMailSender(ctx, result.DialerFactory, result.Resolver)
	return result
}

func (s *SendMailSvc) Run(
	cmd *cobra.Command,
	_ []string,
) error {
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)

	eg, _ := errgroup.WithContext(ctx)
	var err error

	eg.Go(func() error {
		logger.Info().Interface("config", s.Cfg).Msg("Running in errgroup")
		// opts := smtpclient.Opts{}
		// client, err := smtpclient.New(
		// 	ctx,
		// 	sLogger,
		// 	conn,
		// 	smtpclient.TLSOpportunistic,
		// 	false,
		// 	s.Cfg.FromAddr,
		// 	s.Cfg.ToAddr,
		// 	opts,
		// )
		return nil
	})

	err = eg.Wait()
	if err != nil {
		logger.Error().Err(err).Msg("errgroup Wait")
	}
	return err
}
