package config

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ConfigType int

const (
	ConfigTypeHeaders ConfigType = 0
	ConfigTypeDefault ConfigType = 1

	ConfigTypeHeadersStr = "headers"
	ConfigTypeDefaultStr = "default"
)

type ReadFileConfig struct {
	Concurrency  int              `mapstructure:"concurrency"`
	DefaultFrom  string           `mapstructure:"from"`
	FileMails    []FileMailConfig `mapstructure:"file-mails"`
	InPath       string           `mapstructure:"in-path"`
	PollInterval time.Duration    `mapstructure:"poll-interval"`
	RedisAddr    string           `mapstructure:"redis-addr"`
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
