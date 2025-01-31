package config

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Concurrency    int            `mapstructure:"concurrency"`
	Debug          bool           `mapstructure:"debug"`
	InPath         string         `mapstructure:"in_path"`
	PollInterval   time.Duration  `mapstructure:"poll_interval"`
	ReadFileConfig ReadFileConfig `mapstructure:"read_file_config"`
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
