package http

import (
	"context"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(
	_ context.Context,
	engine *gin.Engine,
) error {
	debugGroup := engine.Group("/debug", gin.BasicAuth(gin.Accounts{
		"foo": "bar",
	}))
	pprof.RouteRegister(debugGroup, "pprof")
	return nil
}
