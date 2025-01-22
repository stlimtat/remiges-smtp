package config

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type FromType uint8

const (
	FromTypeHeaders FromType = 0
	FromTypeDefault FromType = 1
)

type ReadFileConfig struct {
	InPath      string   `json:"in_path" mapstructure:"in_path"`
	FromType    FromType `json:"from_type" mapstructure:"from_type"`
	DefaultFrom string   `json:"from" mapstructure:"from"`
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
