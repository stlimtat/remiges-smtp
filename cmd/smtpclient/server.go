/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"net/http"
	"time"

	zerologgin "github.com/go-mods/zerolog-gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/cli"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
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
		Run: cli.WithServerConfig(
			serverCmd.cfg,
			func(cmd *cobra.Command, args []string, cfg config.ServerConfig) {
				serverCmd.server = newServer(cmd, args, cfg)
				serverCmd.server.Run(ctx)
			}),
	}

	return serverCmd, serverCmd.cmd
}

type Server struct {
	Cfg     config.ServerConfig
	Gin     *gin.Engine
	HTTPSvr *http.Server
	InPath  string
}

func newServer(
	cmd *cobra.Command,
	_ []string,
	cfg config.ServerConfig,
) *Server {
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)
	result := &Server{
		Cfg:    cfg,
		InPath: cfg.InPath,
	}

	result.Gin = gin.Default()
	result.Gin.Use(
		zerologgin.LoggerWithOptions(
			&zerologgin.Options{
				Name:   "remiges-smtp",
				Logger: logger,
			},
		),
	)
	result.HTTPSvr = &http.Server{
		Addr:              ":8000",
		Handler:           result.Gin,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return result
}

func (s *Server) Run(ctx context.Context) {
	logger := zerolog.Ctx(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	var err error

	eg.Go(func() error {
		err = s.HTTPSvr.ListenAndServe()
		if err != nil {
			logger.Error().Err(err).Msg("HTTPSvr.ListenAndServe")
		}
		return err
	})

	eg.Go(func() error {
		<-ctx.Done()
		ctx1, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		logger.Warn().Msg("Shutting down")
		err = s.HTTPSvr.Shutdown(ctx1)
		if err != nil {
			logger.Error().Err(err).Msg("HTTPSvr.Shutdown")
		}
		return err
	})

	err = eg.Wait()
	if err != nil {
		logger.Error().Err(err).Msg("errgroup Wait")
	}
}
