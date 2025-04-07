package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRootCmd(t *testing.T) {
	tests := []struct {
		name            string
		ctx             context.Context
		expectedSubcmds []string
	}{
		{
			name: "normal context",
			ctx:  context.Background(),
			expectedSubcmds: []string{
				"gendkim",
				"lookupmx",
				"readfile",
				"sendmail",
				"server",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd(tt.ctx)
			require.NotNil(t, cmd)
			require.NotNil(t, cmd.cmd)

			// Verify subcommands
			subcmds := cmd.cmd.Commands()
			require.Len(t, subcmds, len(tt.expectedSubcmds))

			subcmdNames := make([]string, len(subcmds))
			for i, subcmd := range subcmds {
				subcmdNames[i] = subcmd.Name()
			}
			assert.ElementsMatch(t, tt.expectedSubcmds, subcmdNames)

			// Verify persistent flags
			debugFlag := cmd.cmd.PersistentFlags().Lookup("debug")
			require.NotNil(t, debugFlag)
			assert.Equal(t, "debug", debugFlag.Name)
			assert.Equal(t, "d", debugFlag.Shorthand)
			assert.Equal(t, "bool", debugFlag.Value.Type())
		})
	}
}

func TestExecuteContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		args        []string
		expectError bool
	}{
		{
			name:        "valid command",
			ctx:         context.Background(),
			args:        []string{"--help"},
			expectError: false,
		},
		{
			name:        "invalid command",
			ctx:         context.Background(),
			args:        []string{"nonexistent-command"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd(tt.ctx)
			cmd.cmd.SetArgs(tt.args)

			// Capture output
			var out bytes.Buffer
			cmd.cmd.SetOut(&out)
			cmd.cmd.SetErr(&out)

			err := cmd.ExecuteContext(tt.ctx)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDebugFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "debug flag set",
			args:     []string{"--debug"},
			expected: true,
		},
		{
			name:     "debug flag not set",
			args:     []string{},
			expected: false,
		},
		{
			name:     "debug short flag set",
			args:     []string{"-d"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd(context.Background())
			cmd.cmd.SetArgs(tt.args)

			// Execute the command to trigger flag parsing
			var out bytes.Buffer
			cmd.cmd.SetOut(&out)
			cmd.cmd.SetErr(&out)
			_ = cmd.cmd.Execute()

			// Verify the flag value in viper
			assert.Equal(t, tt.expected, viper.GetBool("debug"))
		})
	}
}
