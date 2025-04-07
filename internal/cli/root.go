/*
Copyright © 2024 Swee Tat Lim <st_lim@stlim.net>
*/

// Package cli provides the command-line interface for the remiges-smtp application.
// It implements various subcommands for SMTP operations using the Cobra CLI framework.
// The package handles command initialization, configuration, and execution flow.
package cli

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

// rootCmd represents the root command structure for the CLI application.
// It encapsulates the Cobra command and provides methods for command execution.
type rootCmd struct {
	cmd *cobra.Command
}

// NewRootCmd initializes and returns a new root command for the CLI application.
// It sets up the command structure, flags, and subcommands.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//
// Returns:
//   - *rootCmd: The initialized root command structure
func NewRootCmd(ctx context.Context) *rootCmd {
	logger := zerolog.Ctx(ctx)
	var err error

	result := &rootCmd{}
	// rootCmd represents the base command when called without any subcommands
	result.cmd = &cobra.Command{
		Use:   "smtpclient",
		Short: "An smtp client for remigres",
		Long:  `An smtp client for remigres`,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			cmdCtx, _ := telemetry.GetLogger(cmd.Context(), cmd.OutOrStdout())
			cmd.SetContext(cmdCtx)
		},
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(config.RootConfigInit)

	result.cmd.PersistentFlags().BoolP(
		"debug", "d",
		false,
		"Run the application in debug mode",
	)
	err = viper.BindPFlag("debug", result.cmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag.debug")
	}
	_, genDKIMCmd := newGenDKIMCmd(ctx)
	_, lookupMXCmd := newLookupMXCmd(ctx)
	_, readFileCmd := newReadFileCmd(ctx)
	_, sendMailCmd := newSendMailCmd(ctx)
	_, serverCmd := newServerCmd(ctx)

	result.cmd.AddCommand(
		genDKIMCmd,
		lookupMXCmd,
		readFileCmd,
		sendMailCmd,
		serverCmd,
	)

	return result
}

// ExecuteContext executes the root command with the provided context.
// It handles the execution of all subcommands and propagates the context.
//
// Parameters:
//   - ctx: Context for command execution
//
// Returns:
//   - error: Non-nil if command execution fails
func (r *rootCmd) ExecuteContext(ctx context.Context) error {
	return r.cmd.ExecuteContext(ctx)
}
