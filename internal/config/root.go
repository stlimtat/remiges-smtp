package config

import (
	"context"
	"crypto/x509"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

const (
	CTX_KEY_CONFIG = "config"
	CONFIG_PATH    = "/app/config"
)

func RootConfigInit() {
	ctx := context.Background()
	_, logger := telemetry.GetLogger(ctx, os.Stdout)

	home, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal().Err(err).Msg("homedir.Dir")
	}
	viper.AddConfigPath(home)
	viper.AddConfigPath(CONFIG_PATH)
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

func GetCertPool(ctx context.Context) *x509.CertPool {
	logger := zerolog.Ctx(ctx)
	// basic cert pool
	result, err := x509.SystemCertPool()
	if err != nil {
		logger.Fatal().Err(err).Msg("x509.SystemCertPool")
	}
	return result
}
