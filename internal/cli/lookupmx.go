package cli

import (
	"context"
	"fmt"
	"log/slog"

	moxDns "github.com/mjl-/mox/dns"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/dns"
	"github.com/stlimtat/remiges-smtp/internal/sendmail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

type lookupMXCmd struct {
	cmd *cobra.Command
}

func newLookupMXCmd(
	ctx context.Context,
) (*lookupMXCmd, *cobra.Command) {
	logger := zerolog.Ctx(ctx)
	var err error

	result := &lookupMXCmd{}
	result.cmd = &cobra.Command{
		Use:   "lookupmx",
		Short: "Lookup MX DNS records for provided domain",
		Long:  `Lookup MX DNS records for provided domain`,
		Args: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			cmdLogger := zerolog.Ctx(ctx)
			cfg := config.NewLookupMXConfig(ctx)
			if len(cfg.Domain) < 1 {
				cmdLogger.Fatal().
					Err(fmt.Errorf("domain fail")).
					Interface("cfg", cfg).
					Msg("Missing fields")
			}
			ctx = config.SetContextConfig(ctx, cfg)
			cmd.SetContext(ctx)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			result := newLookupMXSvc(cmd, args)
			err = result.Run(cmd, args)
			if err != nil {
				logger.Fatal().Err(err).Msg("lookupmx.Run")
			}
		},
	}

	result.cmd.Flags().StringP("lookup-domain", "l", "", "Domain to lookup mx entries")
	err = viper.BindPFlag("domain", result.cmd.Flags().Lookup("lookup-domain"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag")
	}
	return result, result.cmd
}

type LookupMXSvc struct {
	Cfg           config.LookupMXConfig
	DialerFactory sendmail.INetDialerFactory
	MailSender    sendmail.IMailSender
	MoxResolver   moxDns.Resolver
	MyResolver    dns.IResolver
	Slogger       *slog.Logger
}

func newLookupMXSvc(
	cmd *cobra.Command,
	_ []string,
) *LookupMXSvc {
	result := &LookupMXSvc{}
	ctx := cmd.Context()
	result.Cfg = config.GetContextConfig(ctx).(config.LookupMXConfig)
	result.Slogger = telemetry.GetSLogger(ctx)
	result.MoxResolver = moxDns.StrictResolver{
		Log: result.Slogger,
	}
	result.MyResolver = dns.NewResolver(
		ctx,
		result.MoxResolver,
		result.Slogger,
	)
	result.DialerFactory = sendmail.NewDefaultDialerFactory()
	return result
}

func (l *LookupMXSvc) Run(
	cmd *cobra.Command,
	_ []string,
) error {
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)

	result, err := l.MyResolver.LookupMX(ctx, moxDns.Domain{ASCII: l.Cfg.Domain})
	if err != nil {
		logger.Fatal().Err(err).Msg("mailSender.LookupMX")
		return err
	}

	logger.Info().Interface("result", result).Msg("LookupMX")
	return nil
}
