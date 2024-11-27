package config

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func RootConfigInit() {
	logger := log.Logger

	home, err := os.UserHomeDir()
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
	logger.Info().
		Interface("viper_AllSettings", viper.AllSettings()).
		Msg("RootConfigInitialize...Done")
}
