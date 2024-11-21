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
	"github.com/stlimtat/remiges-smtp/internal/config"
	rhttp "github.com/stlimtat/remiges-smtp/internal/http"
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
	Cfg     config.ServerConfig
	Gin     *gin.Engine
	HTTPSvr *http.Server
	InPath  string
}

func newServer(
	cmd *cobra.Command,
	_ []string,
) *Server {
	ctx := cmd.Context()
	logger := zerolog.Ctx(ctx)
	var err error

	cfg := config.NewServerConfig(ctx)

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
	err = rhttp.RegisterAdminRoutes(ctx, result.Gin)
	if err != nil {
		logger.Fatal().Err(err).Msg("http.NewAdminRoutes")
	}

	result.HTTPSvr = &http.Server{
		Addr:              ":8000",
		Handler:           result.Gin,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return result
}

func (s *Server) Run(ctx context.Context) error {
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
	return err
}
