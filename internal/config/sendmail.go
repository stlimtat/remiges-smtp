package config

import (
	"context"

	"github.com/mjl-/mox/smtp"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type SendMailConfig struct {
	Debug    bool         `mapstructure:"debug"`
	From     string       `mapstructure:"from"`
	FromAddr smtp.Address `mapstructure:",omitempty"`
	To       string       `mapstructure:"to"`
	ToAddr   smtp.Address `mapstructure:",omitempty"`
	Msg      string       `mapstructure:"msg"`
}

func NewSendMailConfig(ctx context.Context) SendMailConfig {
	logger := zerolog.Ctx(ctx)
	var err error

	var result SendMailConfig
	err = viper.Unmarshal(&result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unmarshal")
	}

	logger.Info().
		Interface("viper.AllSettings", viper.AllSettings()).
		Interface("result", result).
		Msg("SendMailConfig init")

	// converting from and to email address
	result.FromAddr, err = smtp.ParseAddress(result.From)
	if err != nil {
		logger.Fatal().Err(err).Msg("smtp.ParseAddress.From")
	}

	result.ToAddr, err = smtp.ParseAddress(result.To)
	if err != nil {
		logger.Fatal().Err(err).Msg("smtp.ParseAddress.To")
	}

	return result
}
