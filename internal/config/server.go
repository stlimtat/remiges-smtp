package config

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	InPath          string `mapstructure:"in_path"`
	PollIntervalInt int    `mapstructure:"poll_interval"`
}

func NewServerConfig(ctx context.Context) ServerConfig {
	logger := zerolog.Ctx(ctx)

	var result ServerConfig
	err := viper.Unmarshal(&result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unmarshal")
	}

	logger.Debug().
		Interface("result", result).
		Msg("ServerConfig init")

	return result
}
