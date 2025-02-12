package config

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type FileMailConfig struct {
	Args  []string
	Index int
	Type  string
}

func NewFileMailConfig(ctx context.Context) FileMailConfig {
	logger := zerolog.Ctx(ctx)
	var err error

	var result FileMailConfig
	err = viper.Unmarshal(&result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unmarshal")
	}

	logger.Info().
		Interface("viper.AllSettings", viper.AllSettings()).
		Interface("result", result).
		Msg("FileMailConfig init")

	return result
}
