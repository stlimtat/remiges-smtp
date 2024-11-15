/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type serverCmd struct {
	cmd *cobra.Command
	// cfg serverConfig
}

func newServerCmd(ctx context.Context) (*serverCmd, *cobra.Command) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Testing")

	serverCmd := &serverCmd{}

	// serverCmd represents the server command
	serverCmd.cmd = &cobra.Command{
		Use:   "server",
		Short: "Run the smtpclient",
		Long:  `Runs the smtp client which performs several tasks`,
		// Run:   newServer,
	}
	return serverCmd, serverCmd.cmd
}
