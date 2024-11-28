package config

import (
	"context"
	"os"

	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

const CTX_KEY_CONFIG = "config"

func RootConfigInit() {
	ctx := context.Background()
	_, logger := telemetry.GetLogger(ctx, os.Stdout)

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

func SetContextConfig(ctx context.Context, cfg any) context.Context {
	return context.WithValue(ctx, CTX_KEY_CONFIG, cfg)
}
func GetContextConfig(ctx context.Context) any {
	return ctx.Value(CTX_KEY_CONFIG)
}
