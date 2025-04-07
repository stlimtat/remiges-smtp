package cli

import (
	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
)

// WithServerConfig creates a command wrapper that injects server configuration
// into command execution. It allows commands to access server configuration
// without directly depending on the configuration package.
//
// Parameters:
//   - cfg: Server configuration to be injected
//   - wrappedCmd: Command function that will receive the configuration
//
// Returns:
//   - func(*cobra.Command, []string): A wrapped command function that includes server configuration
func WithServerConfig(
	cfg *config.ServerConfig,
	wrappedCmd func(*cobra.Command, []string, *config.ServerConfig),
) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		wrappedCmd(cmd, args, cfg)
	}
}
