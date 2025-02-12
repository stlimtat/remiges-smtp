package config

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type FromType uint8

const (
	FromTypeHeaders FromType = 0
	FromTypeDefault FromType = 1
)

type ReadFileConfig struct {
	Concurrency  int           `json:"concurrency" mapstructure:"concurrency"`
	DefaultFrom  string        `json:"from" mapstructure:"from"`
	FromType     FromType      `json:"from_type" mapstructure:"from_type"`
	InPath       string        `json:"in_path" mapstructure:"in_path"`
	PollInterval time.Duration `json:"poll_interval" mapstructure:"poll_interval"`
	RedisAddr    string        `json:"redis_addr" mapstructure:"redis_addr"`
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
