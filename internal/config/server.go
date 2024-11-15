package config

import (
	"context"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	InPath string `mapstructure:"in_path"`
}

func NewServerConfig(ctx context.Context) ServerConfig {
	logger := zerolog.Ctx(ctx)

	home, err := homedir.Dir()

	viper.AddConfigPath(home)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("REM")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		logger.Panic().Err(err).Msg("ReadInConfig")
	}
	var result ServerConfig
	err = viper.Unmarshal(&result)
	if err != nil {
		logger.Panic().Err(err).Msg("Unmarshal")
	}

	return result
}
