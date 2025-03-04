package utils

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateIO(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		createFile bool
		path       string
		fileNotDir bool
		wantErr    bool
	}{
		{
			name:       "valid file",
			createFile: true,
			path:       filepath.Join(tmpDir, "test.txt"),
			fileNotDir: true,
			wantErr:    false,
		},
		{
			name:       "valid directory",
			createFile: false,
			path:       tmpDir,
			fileNotDir: false,
			wantErr:    false,
		},
		{
			name:       "relative path with ./",
			createFile: true,
			path:       "./test.txt",
			fileNotDir: true,
			wantErr:    false,
		},
		{
			name:       "relative path with ~/",
			createFile: true,
			path:       "~/test.txt",
			fileNotDir: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			var err error
			path := tt.path
			if tt.createFile {
				if strings.HasPrefix(path, "~/") {
					var home string
					home, err = os.UserHomeDir()
					require.NoError(t, err)
					path = filepath.Join(home, path[2:])
				} else if strings.HasPrefix(path, "./") {
					var wd string
					wd, err = getWorkingDirRelativeToSourceRoot(ctx)
					require.NoError(t, err)
					path = filepath.Join(wd, path[2:])
				}
				_, err = os.Create(path)
				require.NoError(t, err)
			}
			err = ValidateIO(ctx, tt.path, tt.fileNotDir)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.createFile {
				err = os.Remove(path)
				require.NoError(t, err)
			}
		})
	}
}
