package config

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ReadFileConfig struct {
	InPath string
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

	// validate the config
	err = result.Validate(ctx)
	if err != nil {
		return ReadFileConfig{}
	}

	return result
}

func (cfg *ReadFileConfig) Validate(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	// validate the config
	if len(cfg.InPath) < 1 {
		logger.Fatal().
			Err(fmt.Errorf("missing fields")).
			Interface("cfg", cfg).
			Msg("Missing fields")
	}
	fileInfo, err := os.Stat(cfg.InPath)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("InPath does not exist")
	}
	if !fileInfo.IsDir() {
		logger.Fatal().
			Err(fmt.Errorf("InPath is not a directory")).
			Msg("InPath is not a directory")
	}

	return nil
}
