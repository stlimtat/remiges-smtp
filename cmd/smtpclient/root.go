/*
Copyright © 2024 Swee Tat Lim <st_lim@stlim.net>
*/
package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type rootCmd struct {
	cmd     *cobra.Command
	cfgFile string
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func newRootCmd(ctx context.Context) *rootCmd {
	logger := zerolog.Ctx(ctx)

	result := &rootCmd{}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// cobra.OnInitialize(NewConfig)

	// rootCmd represents the base command when called without any subcommands
	result.cmd = &cobra.Command{
		Use:   "smtpclient",
		Short: "An smtp client for remigres",
		Long:  `An smtp client for remigres`,
	}
	result.cmd.PersistentFlags().StringVar(
		&result.cfgFile,
		"config",
		"",
		"config file (default is $HOME/.smtpclient.yaml)",
	)
	err := result.cmd.ExecuteContext(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("rootCmd.Execute")
	}
	_, serverCobraCmd := newServerCmd(ctx)

	result.cmd.AddCommand(
		serverCobraCmd,
	)
	return result
}
