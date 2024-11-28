/*
Copyright © 2024 Swee Tat Lim <st_lim@stlim.net>
*/
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ctx, logger := telemetry.InitLogger(ctx)
	rootCmd := newRootCmd(ctx)
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("ExecuteContext")
	}
}
