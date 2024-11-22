/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"golang.org/x/sync/errgroup"
)

type serverCmd struct {
	cmd    *cobra.Command
	cfg    config.ServerConfig
	server *Server
}

func newServerCmd(ctx context.Context) (*serverCmd, *cobra.Command) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Testing")

	serverCmd := &serverCmd{}
	serverCmd.cfg = config.NewServerConfig(ctx)

	// serverCmd represents the server command
	serverCmd.cmd = &cobra.Command{
		Use:   "server",
		Short: "Run the smtpclient",
		Long:  `Runs the smtp client which performs several tasks`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			serverCmd.server = newServer(cmd, args)
			err = serverCmd.server.Run(ctx)
			return err
		},
	}

	return serverCmd, serverCmd.cmd
}

type Server struct {
	AdminSvr   *http.Server
	Cfg        config.ServerConfig
	FileReader *input.FileReader
}

func newServer(
	cmd *cobra.Command,
	_ []string,
) *Server {
	appCtx := NewAppCtx(cmd.Context())

	result := &Server{
		AdminSvr:   appCtx.AdminSvr,
		Cfg:        appCtx.Cfg,
		FileReader: appCtx.FileReader,
	}
	return result
}

func (s *Server) Run(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	var err error

	eg.Go(func() error {
		err = s.AdminSvr.ListenAndServe()
		if err != nil {
			logger.Error().Err(err).Msg("AdminSvr.ListenAndServe")
		}
		return err
	})

	eg.Go(func() error {
		<-ctx.Done()
		ctx1, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		logger.Warn().Msg("Shutting down")
		err = s.AdminSvr.Shutdown(ctx1)
		if err != nil {
			logger.Error().Err(err).Msg("AdminSvr.Shutdown")
		}
		return err
	})

	eg.Go(func() error {
		// fileReader is able to stop based on ctx.Done
		return s.FileReader.Run(ctx)
	})

	err = eg.Wait()
	if err != nil {
		logger.Error().Err(err).Msg("errgroup Wait")
	}
	return err
}
