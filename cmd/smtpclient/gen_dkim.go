package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/crypto"
	"github.com/stlimtat/remiges-smtp/internal/dkim"
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

	result.cmd.Flags().String("algorithm", "rsa", "Key type to generate DKIM keys, dns record and config")
	result.cmd.Flags().Int("bit-size", 2048, "Bit size of the DKIM keys")
	result.cmd.Flags().String("dkim-domain", "", "Domain to generate DKIM keys, dns record and config")
	result.cmd.Flags().String("hash", "sha256", "Hash algorithm to use for DKIM keys, dns record and config")
	result.cmd.Flags().String("out-path", "./config", "Path to write DKIM keys, dns record and config")
	result.cmd.Flags().String("selector", "key001", "Selector for DKIM keys")
	err = viper.BindPFlag("algorithm", result.cmd.Flags().Lookup("algorithm"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - algorithm")
	}
	err = viper.BindPFlag("bit-size", result.cmd.Flags().Lookup("bit-size"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - bit-size")
	}
	err = viper.BindPFlag("dkim-domain", result.cmd.Flags().Lookup("dkim-domain"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - dkim-domain")
	}
	err = viper.BindPFlag("hash", result.cmd.Flags().Lookup("hash"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - hash")
	}
	err = viper.BindPFlag("algorithm", result.cmd.Flags().Lookup("algorithm"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - algorithm")
	}
	err = viper.BindPFlag("out-path", result.cmd.Flags().Lookup("out-path"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - out-path")
	}
	err = viper.BindPFlag("selector", result.cmd.Flags().Lookup("selector"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag - selector")
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

	// Initialize the system
	cfg := config.GetContextConfig(ctx).(config.GenDKIMConfig)

	factory := &crypto.CryptoFactory{}
	keyWriter, err := crypto.NewKeyWriter(ctx, cfg.OutPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("crypto.NewKeyWriter")
	}
	txtGen := &dkim.TxtGen{}

	// Perform the running
	_, err = factory.Init(ctx, keyWriter)
	if err != nil {
		logger.Fatal().Err(err).Msg("crypto.CryptoFactory.Init")
	}

	publicKeyPEM, privateKeyPEM, err := factory.GenerateKey(ctx, cfg.BitSize, cfg.Domain, cfg.Algorithm)
	if err != nil {
		logger.Fatal().Err(err).Msg("crypto.CryptoFactory.GenerateKey")
	}

	_, privateKeyPath, err := factory.WriteKey(ctx, cfg.Domain, publicKeyPEM, privateKeyPEM)
	if err != nil {
		logger.Fatal().Err(err).Msg("crypto.CryptoFactory.WriteKey")
	}

	txtEntry, err := txtGen.Generate(ctx, cfg.Domain, cfg.Algorithm, cfg.Selector, publicKeyPEM)
	if err != nil {
		logger.Fatal().Err(err).Msg("dkim.TxtGen.Generate")
	}

	fmt.Printf(
		GenDKIMResult,
		cfg.Domain,
		txtEntry,
		cfg.Domain,
		cfg.Selector,
		cfg.Algorithm,
		cfg.Hash,
		privateKeyPath,
		cfg.Selector,
	)

	return nil
}

const GenDKIMResult = `To enable DKIM for %s, add the following TXT record to your DNS:

%s

To ensure that DKIM is working for the smtpclient, you need to add the following to
the smtpclient config, and this should be after merge_body:

` + "```" + `yaml
mail_processors:
  - type: dkim
    index: 100
    args:
      domain-str: %s
      dkim:
        selectors:
          %s:
            algorithm: %s
            body-relaxed: true
            expiration: 72h
            hash: %s
            header-relaxed: true
            headers:
              - from
              - to
              - subject
              - date
              - message-id
              - content-type
            private-key-file: %s
            seal-headers: false
            selector-domain: %s
` + "```" + `
Then restart the smtpclient.
`
