/*
Copyright © 2024 Swee Tat Lim <st_lim@stlim.net>
*/
package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
)

type rootCmd struct {
	cmd     *cobra.Command
	cfgFile string
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func newRootCmd(ctx context.Context) *rootCmd {
	logger := zerolog.Ctx(ctx)
	var err error

	result := &rootCmd{}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(config.RootConfigInit)

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
	err = viper.BindPFlag("config", result.cmd.PersistentFlags().Lookup("config"))
	if err != nil {
		logger.Fatal().Err(err).Msg("viper.BindPFlag")
	}
	err = result.cmd.ExecuteContext(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("rootCmd.Execute")
	}
	_, serverCmd := newServerCmd(ctx)

	result.cmd.AddCommand(
		serverCmd,
	)
	return result
}
