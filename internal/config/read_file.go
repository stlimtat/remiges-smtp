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
	Concurrency  int              `json:"concurrency" mapstructure:"concurrency"`
	DefaultFrom  string           `json:"from" mapstructure:"from"`
	FileMails    []FileMailConfig `json:"file_mails" mapstructure:"file_mails"`
	InPath       string           `json:"in_path" mapstructure:"in_path"`
	PollInterval time.Duration    `json:"poll_interval" mapstructure:"poll_interval"`
	RedisAddr    string           `json:"redis_addr" mapstructure:"redis_addr"`
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
