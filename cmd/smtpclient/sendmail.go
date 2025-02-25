/*
Copyright © 2024 Lim Swee Tat <st_lim@stlim.net>
*/
package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
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
	*GenericSvc
}

func newSendMailSvc(
	cmd *cobra.Command,
	args []string,
) *SendMailSvc {
	result := &SendMailSvc{}
	result.GenericSvc = newGenericSvc(cmd, args)
	return result
}

func (s *SendMailSvc) Run(
	cmd *cobra.Command,
	_ []string,
) error {
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)
	var err error

	// refresh file list
	_, err = s.FileReader.RefreshList(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("FileMailService.RefreshList")
		return err
	}

	// read a file
	fileInfo, myMail, err := s.SendMailService.ReadNextMail(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("FileMailService.ReadNextMail")
		return err
	}
	logger.Info().
		Str("fileInfo", fileInfo.ID).
		Str("from", myMail.From.String()).
		Msg("ReadNextMail")

	return nil
}
