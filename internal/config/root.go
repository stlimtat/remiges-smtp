package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func RootConfigInit() {
	logger := log.Logger

	home, err := homedir.Dir()
	if err != nil {
		logger.Fatal().Err(err).Msg("homedir.Dir")
	}
	viper.AddConfigPath(home)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("REM")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("ReadInConfig")
	}
	logger.Debug().
		Interface("viper_AllSettings", viper.AllSettings()).
		Msg("RootConfigInitialize...Done")
}
