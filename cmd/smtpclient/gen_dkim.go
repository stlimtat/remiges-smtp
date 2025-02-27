package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/crypto"
)

type genDKIMCmd struct {
	cmd *cobra.Command
}

func newGenDKIMCmd(
	ctx context.Context,
) (*genDKIMCmd, *cobra.Command) {
	logger := zerolog.Ctx(ctx)
	var err error

	result := &genDKIMCmd{}
	result.cmd = &cobra.Command{
		Use:   "gendkim",
		Short: "Generate DKIM keys, dns record and config for provided domain",
		Long:  `Generate DKIM keys, dns record and config for provided domain`,
		Args: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			cmdLogger := zerolog.Ctx(ctx)
			cfg := config.NewGenDKIMConfig(ctx)
			if len(cfg.Domain) < 1 {
				cmdLogger.Fatal().
					Err(fmt.Errorf("domain fail")).
					Interface("cfg", cfg).
					Msg("Missing fields")
			}
			if len(cfg.OutPath) < 1 {
				cmdLogger.Fatal().
					Err(fmt.Errorf("out-path fail")).
					Interface("cfg", cfg).
					Msg("Missing fields")
			}
			ctx = config.SetContextConfig(ctx, cfg)
			cmd.SetContext(ctx)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			result := newGenDKIMSvc(cmd, args)
			err = result.Run(cmd, args)
			if err != nil {
				logger.Fatal().Err(err).Msg("genDKIM.Run")
			}
		},
	}

	result.cmd.Flags().Int(
		"bit-size",
		2048, "Bit size of the DKIM keys",
	)
	err = viper.BindPFlag("bit-size", result.cmd.Flags().Lookup("bit-size"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - bit-size")
	}

	result.cmd.Flags().String(
		"dkim-domain",
		"", "Domain to generate DKIM keys, dns record and config",
	)
	err = viper.BindPFlag("dkim-domain", result.cmd.Flags().Lookup("dkim-domain"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - dkim-domain")
	}

	result.cmd.Flags().String(
		"out-path",
		"~/config", "Path to write DKIM keys, dns record and config",
	)
	err = viper.BindPFlag("out-path", result.cmd.Flags().Lookup("out-path"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag")
	}
	return result, result.cmd
}

type GenDKIMSvc struct {
	Cfg config.GenDKIMConfig
}

func newGenDKIMSvc(
	cmd *cobra.Command,
	_ []string,
) *GenDKIMSvc {
	result := &GenDKIMSvc{}
	ctx := cmd.Context()
	result.Cfg = config.GetContextConfig(ctx).(config.GenDKIMConfig)
	return result
}

func (_ *GenDKIMSvc) Run(
	cmd *cobra.Command,
	_ []string,
) error { //nolint:unparam // result 0 is always nil in this
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)

	cfg := config.GetContextConfig(ctx).(config.GenDKIMConfig)

	factory := &crypto.CryptoFactory{}
	keyWriter := crypto.NewKeyWriter(ctx, cfg.OutPath)
	err := keyWriter.Validate(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("crypto.KeyWriter.Validate")
	}
	_, err = factory.Init(ctx, cfg.KeyType, keyWriter)
	if err != nil {
		logger.Fatal().Err(err).Msg("crypto.CryptoFactory.Init")
	}

	publicKeyPEM, privateKeyPEM, err := factory.GenerateKey(ctx, cfg.BitSize, cfg.Domain)
	if err != nil {
		logger.Fatal().Err(err).Msg("crypto.CryptoFactory.GenerateKey")
	}

	err = factory.WriteKey(ctx, cfg.Domain, publicKeyPEM, privateKeyPEM)
	if err != nil {
		logger.Fatal().Err(err).Msg("crypto.CryptoFactory.WriteKey")
	}

	return nil
}
