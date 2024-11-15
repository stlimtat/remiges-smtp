package cli

import (
	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
)

func WithServerConfig(
	cfg config.ServerConfig,
	wrappedCmd func(*cobra.Command, []string, config.ServerConfig),
) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		wrappedCmd(cmd, args, cfg)
	}
}
