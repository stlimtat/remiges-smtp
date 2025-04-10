package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenDKIMCmd(t *testing.T) {
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
			ctx, _ := telemetry.InitLogger(context.Background())
			telemetry.SetGlobalLogLevel(zerolog.ErrorLevel)
			cmd, cobraCmd := newGenDKIMCmd(ctx)
			if tt.expectedErr {
				assert.Nil(t, cmd)
				assert.Nil(t, cobraCmd)
				return
			}

			require.NotNil(t, cmd)
			require.NotNil(t, cobraCmd)

			// Verify command flags
			flags := []struct {
				name      string
				shorthand string
				valueType string
			}{
				{"algorithm", "", "string"},
				{"bit-size", "", "int"},
				{"dkim-domain", "", "string"},
				{"hash", "", "string"},
				{"out-path", "", "string"},
				{"selector", "", "string"},
			}

			for _, flag := range flags {
				f := cobraCmd.Flags().Lookup(flag.name)
				require.NotNil(t, f, "flag %s not found", flag.name)
				assert.Equal(t, flag.valueType, f.Value.Type(), "flag %s has wrong type", flag.name)
			}
		})
	}
}

func TestGenDKIMSvc_Run(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	tests := []struct {
		name        string
		algorithm   string
		tempDir     string
		expectError bool
	}{
		{
			name:        "valid configuration",
			algorithm:   "rsa",
			tempDir:     tmpDir,
			expectError: false,
		},
		{
			name:        "invalid algorithm - defaults to rsa",
			algorithm:   "invalid",
			tempDir:     tmpDir,
			expectError: false,
		},
		// {
		// 	name:        "invalid output path",
		// 	algorithm:   "rsa",
		// 	tempDir:     "/nonexistent/path",
		// 	expectError: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.InitLogger(ctx)
			telemetry.SetGlobalLogLevel(zerolog.ErrorLevel)

			// Setup configuration
			cfg := config.GenDKIMConfig{
				Domain:    "example.com",
				OutPath:   tt.tempDir,
				Algorithm: tt.algorithm,
				BitSize:   2048,
				Hash:      "sha256",
				Selector:  "key001",
			}
			ctx = config.SetContextConfig(ctx, cfg)

			cmd := &cobra.Command{}
			cmd.SetContext(ctx)

			svc := newGenDKIMSvc(cmd, nil)

			err := svc.Run(cmd, nil)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify output files
				expectedFiles := []string{
					filepath.Join(svc.Cfg.OutPath, svc.Cfg.Domain+".pem"),
					filepath.Join(svc.Cfg.OutPath, svc.Cfg.Domain+".pub"),
				}

				for _, file := range expectedFiles {
					_, err := os.Stat(file)
					assert.NoError(t, err, "file %s should exist", file)
				}
				defer os.Remove(filepath.Join(svc.Cfg.OutPath, svc.Cfg.Domain+".pem"))
				defer os.Remove(filepath.Join(svc.Cfg.OutPath, svc.Cfg.Domain+".pub"))
			}
		})
	}
}

// func TestGenDKIMCmd_ArgsValidation(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		args        []string
// 		envVars     map[string]string
// 		expectError bool
// 	}{
// 		{
// 			name: "valid arguments",
// 			args: []string{"--dkim-domain=example.com", "--out-path=/tmp"},
// 			envVars: map[string]string{
// 				"DKIM_DOMAIN": "example.com",
// 				"OUT_PATH":    "/tmp",
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name:        "missing domain",
// 			args:        []string{"--out-path=/tmp"},
// 			expectError: true,
// 		},
// 		{
// 			name:        "missing output path",
// 			args:        []string{"--dkim-domain=example.com"},
// 			expectError: true,
// 		},
// 		{
// 			name: "invalid bit size",
// 			args: []string{
// 				"--dkim-domain=example.com",
// 				"--out-path=/tmp",
// 				"--bit-size=1024", // Too small for RSA
// 			},
// 			expectError: true,
// 		},
// 		{
// 			name: "invalid hash algorithm",
// 			args: []string{
// 				"--dkim-domain=example.com",
// 				"--out-path=/tmp",
// 				"--hash=md5", // MD5 is not supported
// 			},
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Set environment variables
// 			for k, v := range tt.envVars {
// 				os.Setenv(k, v)
// 				defer os.Unsetenv(k)
// 			}

// 			ctx := context.Background()
// 			ctx, _ = telemetry.InitLogger(ctx)

// 			cmd, cobraCmd := newGenDKIMCmd(ctx)
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
