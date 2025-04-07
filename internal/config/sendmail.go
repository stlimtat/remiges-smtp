package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mjl-/mox/smtp"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/proxy"
)

type SendMailConfig struct {
	Debug          bool                  `mapstructure:"debug"`
	Dialer         DialerConfig          `mapstructure:"dialer"`
	From           string                `mapstructure:"from"`
	FromAddr       smtp.Address          `mapstructure:",omitempty"`
	To             string                `mapstructure:"to"`
	ToAddr         smtp.Address          `mapstructure:",omitempty"`
	Msg            string                `mapstructure:"msg"`
	MsgBytes       []byte                `mapstructure:",omitempty"`
	MailProcessors []MailProcessorConfig `mapstructure:"mail-processors"`
	Outputs        []OutputConfig        `mapstructure:"outputs"`
	PollInterval   time.Duration         `mapstructure:"poll-interval"`
	ReadFileConfig ReadFileConfig        `mapstructure:"read-file"`
}

type DialerConfig struct {
	Socks5  string        `mapstructure:"socks5"`
	Auth    proxy.Auth    `mapstructure:"auth"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func CobraSendMailArgsFunc(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	cmdLogger := zerolog.Ctx(ctx)
	cfg := NewSendMailConfig(ctx)
	if len(cfg.From) < 1 || len(cfg.To) < 1 {
		cmdLogger.Fatal().
			Err(fmt.Errorf("missing fields")).
			Interface("cfg", cfg).
			Msg("Missing fields")
	}
	ctx = SetContextConfig(ctx, cfg)
	cmd.SetContext(ctx)
	return nil
}

func NewSendMailConfig(ctx context.Context) SendMailConfig {
	logger := zerolog.Ctx(ctx)
	var err error

	viper.SetDefault("msg", "Test message\r\n")

	// setting up default values
	result := SendMailConfig{
		MailProcessors: DefaultMailProcessorConfigs(),
		Outputs:        DefaultOutputConfig(ctx),
		ReadFileConfig: ReadFileConfig{
			FileMails: DefaultFileMailConfigs(),
			InPath:    "inbox",
		},
	}

	err = viper.Unmarshal(&result)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unmarshal")
	}

	allSettings := viper.AllSettings()

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

	logger.Info().
		Interface("allSettings", allSettings).
		Interface("result", result).
		Msg("SendMailConfig init")

	return result
}

// Add configuration validation
func ValidateConfig(cfg *SendMailConfig) error {
	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Add custom validation logic
	// if err := validateMailProcessors(cfg.MailProcessors); err != nil {
	// 	return fmt.Errorf("invalid mail processors: %w", err)
	// }

	return nil
}
