/*
Copyright © 2024 Swee Tat Lim <st_lim@stlim.net>
*/
package main

import (
	"context"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
)

func main() {
	ctx := context.Background()
	ctx, logger := telemetry.InitLogger(ctx)
	rootCmd := newRootCmd(ctx)
	err := rootCmd.cmd.ExecuteContext(ctx)
	if err != nil {
		logger.Panic().Err(err).Msg("ExecuteContext")
	}
}
