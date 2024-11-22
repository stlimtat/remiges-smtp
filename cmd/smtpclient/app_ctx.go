package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	zerologgin "github.com/go-mods/zerolog-gin"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	rhttp "github.com/stlimtat/remiges-smtp/internal/http"
	"github.com/stlimtat/remiges-smtp/pkg/input"
)

type AppCtx struct {
	AdminSvr   *http.Server
	Cfg        config.ServerConfig
	FileReader *input.FileReader
	Gin        *gin.Engine
}

func NewAppCtx(
	ctx context.Context,
) *AppCtx {
	logger := zerolog.Ctx(ctx)
	var err error
	result := &AppCtx{}

	result.Cfg = config.NewServerConfig(ctx)
	result.FileReader = input.NewFileReader(
		ctx,
		result.Cfg.InPath,
		result.Cfg.PollInterval,
	)

	result.Gin = gin.New()
	result.Gin.Use(gin.Recovery())
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

	result.AdminSvr = &http.Server{
		Addr:              ":8000",
		Handler:           result.Gin,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return result
}
