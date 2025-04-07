package cli

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	moxDns "github.com/mjl-/mox/dns"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/dns"
	"github.com/stlimtat/remiges-smtp/internal/sendmail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

// lookupMXCmd represents the command for looking up MX (Mail Exchange) DNS records
// for a given domain. It provides functionality to query and display the mail server
// configuration for email domains.
type lookupMXCmd struct {
	cmd *cobra.Command
}

// newLookupMXCmd creates and initializes a new MX lookup command.
// It sets up command flags, validation, and execution logic.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//
// Returns:
//   - *lookupMXCmd: The initialized command structure
//   - *cobra.Command: The Cobra command for CLI integration
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
		RunE: func(cmd *cobra.Command, args []string) error {
			result := newLookupMXSvc(cmd, args)
			err = result.Run(cmd, args)
			if err != nil {
				logger.Error().Err(err).Msg("lookupmx.Run")
				return err
			}
			return nil
		},
	}

	result.cmd.Flags().StringP("lookup-domain", "l", "", "Domain to lookup mx entries")
	err = viper.BindPFlag("domain", result.cmd.Flags().Lookup("lookup-domain"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag")
	}
	return result, result.cmd
}

// LookupMXSvc handles the service layer for MX record lookups.
// It manages DNS resolution and provides functionality to query MX records
// for email domains.
type LookupMXSvc struct {
	Cfg           config.LookupMXConfig
	DialerFactory sendmail.INetDialerFactory
	MailSender    sendmail.IMailSender
	MoxResolver   moxDns.Resolver
	MyResolver    dns.IResolver
	Slogger       *slog.Logger
}

// newLookupMXSvc creates a new MX lookup service instance.
// It initializes the service with the provided configuration and sets up
// DNS resolution components.
//
// Parameters:
//   - cmd: The Cobra command instance
//   - args: Command arguments
//
// Returns:
//   - *LookupMXSvc: The initialized service instance
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
	return result
}

// Run executes the MX record lookup process.
// It queries the DNS system for MX records of the specified domain
// and logs the results.
//
// Parameters:
//   - cmd: The Cobra command instance
//   - args: Command arguments
//
// Returns:
//   - error: Non-nil if the lookup process fails
func (l *LookupMXSvc) Run(
	cmd *cobra.Command,
	_ []string,
) error {
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)

	if l.Cfg.Domain == "" {
		logger.Warn().Msg("domain is empty")
		return fmt.Errorf("domain is empty")
	}
	sublogger := logger.With().Str("domain", l.Cfg.Domain).Logger()

	domain, err := moxDns.ParseDomain(l.Cfg.Domain)
	if err != nil {
		sublogger.Error().Err(err).Msg("moxDns.ParseDomain")
		return err
	}
	if !strings.HasSuffix(domain.ASCII, ".") {
		domain.ASCII += "."
	}

	result, err := l.MyResolver.LookupMX(ctx, domain)
	if err != nil {
		sublogger.Error().Err(err).Msg("mailSender.LookupMX")
		return err
	}

	sublogger.Info().Interface("result", result).Msg("LookupMX")
	return nil
}
