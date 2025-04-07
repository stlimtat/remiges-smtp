package cli

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestWithServerConfig(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.ServerConfig
		cmdArgs     []string
		expectError bool
	}{
		{
			name: "valid server config",
			cfg: &config.ServerConfig{
				Debug:        true,
				PollInterval: 5 * time.Second,
			},
			cmdArgs:     []string{},
			expectError: false,
		},
		{
			name:        "nil server config",
			cfg:         nil,
			cmdArgs:     []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test command
			cmd := &cobra.Command{
				Use: "test",
				Run: func(cmd *cobra.Command, args []string) {
					// This will be replaced by the wrapper
				},
			}

			// Track if the wrapped function was called
			called := false
			var receivedCfg *config.ServerConfig

			// Create the wrapped command function
			wrappedCmd := func(cmd *cobra.Command, args []string, cfg *config.ServerConfig) {
				called = true
				receivedCfg = cfg
			}

			// Create the wrapper
			wrapper := WithServerConfig(tt.cfg, wrappedCmd)

			// Set the wrapper as the command's Run function
			cmd.Run = wrapper

			// Execute the command
			cmd.SetArgs(tt.cmdArgs)
			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, called, "wrapped command was not called")
				assert.Equal(t, tt.cfg, receivedCfg, "received incorrect config")
			}
		})
	}
}

func TestWithServerConfig_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.ServerConfig
		wrappedCmd  func(*cobra.Command, []string, *config.ServerConfig)
		expectError bool
	}{
		{
			name: "wrapped command panics",
			cfg: &config.ServerConfig{
				Debug:        true,
				PollInterval: 5 * time.Second,
			},
			wrappedCmd: func(cmd *cobra.Command, args []string, cfg *config.ServerConfig) {
				panic("test panic")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use: "test",
			}

			wrapper := WithServerConfig(tt.cfg, tt.wrappedCmd)
			cmd.Run = wrapper

			// Execute the command and capture any panic
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			err := cmd.Execute()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
