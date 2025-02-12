package config

import (
	"context"
	"strings"

	"github.com/mjl-/mox/smtp"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type SendMailConfig struct {
	Debug          bool                  `mapstructure:"debug"`
	From           string                `mapstructure:"from"`
	FromAddr       smtp.Address          `mapstructure:",omitempty"`
	To             string                `mapstructure:"to"`
	ToAddr         smtp.Address          `mapstructure:",omitempty"`
	Msg            string                `mapstructure:"msg"`
	MsgBytes       []byte                `mapstructure:",omitempty"`
	ReadFileConfig ReadFileConfig        `mapstructure:"read_file"`
	MailProcessors []MailProcessorConfig `mapstructure:"mail_processors"`
}

func NewSendMailConfig(ctx context.Context) SendMailConfig {
	logger := zerolog.Ctx(ctx)
	var err error

	viper.SetDefault("msg", "Test message\r\n")

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

	if !strings.HasSuffix(result.Msg, "\r\n") {
		result.Msg += "\r\n"
	}

	result.MsgBytes = []byte(result.Msg)

	return result
}
