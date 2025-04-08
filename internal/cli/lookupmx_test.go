package cli

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLookupMXCmd(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr bool
	}{
		{
			name:        "valid context",
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.InitLogger(ctx)
			cmd, cobraCmd := newLookupMXCmd(ctx)
			if tt.expectedErr {
				assert.Nil(t, cmd)
				assert.Nil(t, cobraCmd)
				return
			}

			require.NotNil(t, cmd)
			require.NotNil(t, cobraCmd)

			// Verify command flags
			flag := cobraCmd.Flags().Lookup("lookup-domain")
			require.NotNil(t, flag, "lookup-domain flag not found")
			assert.Equal(t, "string", flag.Value.Type(), "lookup-domain flag has wrong type")
		})
	}
}

func TestNewLookupMXSvc(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "valid configuration",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.InitLogger(ctx)

			cfg := config.LookupMXConfig{
				Domain: "example.com",
			}
			ctx = config.SetContextConfig(ctx, cfg)

			cmd := &cobra.Command{}
			cmd.SetContext(ctx)

			svc := newLookupMXSvc(cmd, nil)

			if tt.expectError {
				// Verify that critical components are nil
				assert.Nil(t, svc.MyResolver)
			} else {
				// Verify that all components are initialized
				assert.NotNil(t, svc.MyResolver)
				assert.NotNil(t, svc.MoxResolver)
				assert.NotNil(t, svc.Slogger)
				assert.Equal(t, "example.com", svc.Cfg.Domain)
			}
		})
	}
}

func TestLookupMXSvc_Run(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		expectError bool
	}{
		{
			name:        "successful lookup",
			domain:      "example.com",
			expectError: false,
		},
		{
			name:        "invalid domain",
			domain:      "invalid domain",
			expectError: true,
		},
		{
			name:        "empty domain",
			domain:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.InitLogger(ctx)

			cfg := config.LookupMXConfig{
				Domain: tt.domain,
			}
			ctx = config.SetContextConfig(ctx, cfg)

			cmd := &cobra.Command{}
			cmd.SetContext(ctx)

			svc := newLookupMXSvc(cmd, nil)

			err := svc.Run(cmd, nil)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// func TestLookupMXCmd_ArgsValidation(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		args        []string
// 		envVars     map[string]string
// 		expectError bool
// 	}{
// 		{
// 			name: "valid arguments",
// 			args: []string{"--lookup-domain=example.com"},
// 			envVars: map[string]string{
// 				"DOMAIN": "example.com",
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name:        "missing domain",
// 			args:        []string{},
// 			expectError: true,
// 		},
// 		{
// 			name: "empty domain",
// 			args: []string{"--lookup-domain="},
// 			envVars: map[string]string{
// 				"DOMAIN": "",
// 			},
// 			expectError: true,
// 		},
// 		{
// 			name: "invalid domain format",
// 			args: []string{"--lookup-domain=invalid@domain"},
// 			envVars: map[string]string{
// 				"DOMAIN": "invalid@domain",
// 			},
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Set environment variables
// 			for k, v := range tt.envVars {
// 				t.Setenv(k, v)
// 			}

// 			ctx := context.Background()
// 			ctx, _ = telemetry.InitLogger(ctx)

// 			cmd, cobraCmd := newLookupMXCmd(ctx)
// 			require.NotNil(t, cmd)
// 			require.NotNil(t, cobraCmd)

// 			cobraCmd.SetArgs(tt.args)

// 			err := cobraCmd.Execute()
// 			if tt.expectError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }
