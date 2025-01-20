/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/input"
)

type readFileCmd struct {
	cmd *cobra.Command
}

func newReadFileCmd(ctx context.Context) (*readFileCmd, *cobra.Command) {
	logger := zerolog.Ctx(ctx)
	var err error

	result := &readFileCmd{}

	// sendMailCmd represents the server command
	result.cmd = &cobra.Command{
		Use:   "readfile",
		Short: "Reads a df file from the testdata directory",
		Long:  `Reads a df file from the testdata directory, and also reads the corresponding qf file`,
		Args: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			_ = zerolog.Ctx(ctx)
			cfg := config.NewReadFileConfig(ctx)
			ctx = config.SetContextConfig(ctx, cfg)
			cmd.SetContext(ctx)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := newReadFileSvc(cmd, args)
			err = svc.Run(cmd, args)
			return err
		},
	}

	result.cmd.Flags().StringP(
		"path", "p",
		"", "Path to the directory containing the df and qf files",
	)
	err = viper.BindPFlag("in_path", result.cmd.Flags().Lookup("path"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag.inpath")
	}

	return result, result.cmd
}

type ReadFileSvc struct {
	Cfg        config.ReadFileConfig
	FileReader input.IFileReader
}

func newReadFileSvc(
	cmd *cobra.Command,
	_ []string,
) *ReadFileSvc {
	var err error
	result := &ReadFileSvc{}
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)
	result.Cfg = config.GetContextConfig(ctx).(config.ReadFileConfig)
	result.FileReader, err = input.NewDefaultFileReader(ctx, result.Cfg.InPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("newReadFileSvc.FileReader")
	}
	return result
}

func (s *ReadFileSvc) Run(
	cmd *cobra.Command,
	_ []string,
) error {
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)
	var err error

	_, err = s.FileReader.RefreshList(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("s.FileReader.RefreshList")
		return err
	}

	dfReader, qfReader, err := s.FileReader.ReadNextFile(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("s.FileReader.ReadNextFile")
		return err
	}

	logger.Info().Interface("dfReader", dfReader).Interface("qfReader", qfReader).Msg("ReadNextFile")

	return err
}
