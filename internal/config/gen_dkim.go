package config

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type GenDKIMConfig struct {
	BitSize  int    `mapstructure:"bit-size,omitempty"`
	Domain   string `mapstructure:"dkim-domain"`
	Hash     string `mapstructure:"hash,omitempty"`
	KeyType  string `mapstructure:"key-type,omitempty"`
	OutPath  string `mapstructure:"out-path"`
	Selector string `mapstructure:"selector"`
}

func NewGenDKIMConfig(ctx context.Context) GenDKIMConfig {
	logger := zerolog.Ctx(ctx)
	var err error
	var result GenDKIMConfig
	viper.SetDefault("bit-size", 2048)
	viper.SetDefault("hash", "sha256")
	viper.SetDefault("key-type", "rsa")
	viper.SetDefault("out-path", "./config")
	err = viper.Unmarshal(&result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unmarshal")
	}
	allSettings := viper.AllSettings()

	logger.Info().
		Interface("allSettings", allSettings).
		Interface("result", result).
		Msg("GenDKIMConfig init")

	return result
}
