package config

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type LookupMXConfig struct {
	Domain string
}

func NewLookupMXConfig(ctx context.Context) LookupMXConfig {
	logger := zerolog.Ctx(ctx)
	var err error

	var result LookupMXConfig
	err = viper.Unmarshal(&result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unmarshal")
	}

	logger.Info().
		Interface("viper.AllSettings", viper.AllSettings()).
		Interface("result", result).
		Msg("LookupMXConfig init")

	return result
}
