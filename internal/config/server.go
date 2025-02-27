package config

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Debug          bool           `mapstructure:"debug"`
	PollInterval   time.Duration  `mapstructure:"poll-interval"`
	ReadFileConfig ReadFileConfig `mapstructure:"read-file"`
}

func NewServerConfig(ctx context.Context) ServerConfig {
	logger := zerolog.Ctx(ctx)

	var result ServerConfig
	err := viper.Unmarshal(&result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unmarshal")
	}

	logger.Info().
		Interface("viper.AllSettings", viper.AllSettings()).
		Interface("result", result).
		Msg("ServerConfig init")

	return result
}
