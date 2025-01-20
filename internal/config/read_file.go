package config

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ReadFileConfig struct {
	InPath string `json:"in_path" mapstructure:"in_path"`
}

func NewReadFileConfig(ctx context.Context) ReadFileConfig {
	logger := zerolog.Ctx(ctx)
	var err error

	var result ReadFileConfig
	err = viper.Unmarshal(&result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unmarshal")
	}

	logger.Info().
		Interface("viper.AllSettings", viper.AllSettings()).
		Interface("result", result).
		Msg("ReadFileConfig init")

	return result
}
