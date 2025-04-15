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
		errMsg     string
		expectNew  bool
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
			name:       "create new file",
			createFile: true,
			path:       filepath.Join(tmpDir, "newfile.txt"),
			fileNotDir: true,
			wantErr:    true,
			errMsg:     "ToIgnore: create new file",
			expectNew:  true,
		},
		{
			name:       "non-existent file without create",
			createFile: false,
			path:       filepath.Join(tmpDir, "nonexistent.txt"),
			fileNotDir: true,
			wantErr:    true,
			errMsg:     "failed to get file info",
		},
		{
			name:       "empty path",
			createFile: false,
			path:       "",
			fileNotDir: true,
			wantErr:    true,
			errMsg:     "path is required",
		},
		{
			name:       "file when directory expected",
			createFile: true,
			path:       filepath.Join(tmpDir, "test.txt"),
			fileNotDir: false,
			wantErr:    true,
			errMsg:     "path is not a directory",
		},
		{
			name:       "directory when file expected",
			createFile: false,
			path:       tmpDir,
			fileNotDir: true,
			wantErr:    true,
			errMsg:     "path is not a file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			var err error
			path := tt.path
			if tt.createFile && !tt.expectNew {
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
				// Create parent directory if it doesn't exist
				dir := filepath.Dir(path)
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					require.NoError(t, os.MkdirAll(dir, 0755))
				}
				_, err = os.Create(path)
				require.NoError(t, err)
			}
			gotPath, err := ValidateIO(ctx, tt.path, tt.fileNotDir, tt.createFile)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				if tt.expectNew {
					// Verify the file was created
					_, err := os.Stat(gotPath)
					assert.NoError(t, err)
					// Clean up the created file
					err = os.Remove(gotPath)
					require.NoError(t, err)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, path, gotPath)
			if tt.createFile && !tt.expectNew {
				err = os.Remove(path)
				require.NoError(t, err)
			}
		})
	}
}

// func TestGetWorkingDirRelativeToSourceRoot(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		setup       func(t *testing.T) (string, func())
// 		wantErr     bool
// 		errMsg      string
// 		expectBazel bool
// 	}{
// 		{
// 			name: "normal source root",
// 			setup: func(t *testing.T) (string, func()) {
// 				tmpDir := t.TempDir()
// 				require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".git"), 0755))
// 				oldWd, err := os.Getwd()
// 				require.NoError(t, err)
// 				require.NoError(t, os.Chdir(tmpDir))
// 				return tmpDir, func() {
// 					require.NoError(t, os.Chdir(oldWd))
// 				}
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "bazel environment",
// 			setup: func(t *testing.T) (string, func()) {
// 				tmpDir := t.TempDir()
// 				bazelDir := filepath.Join(tmpDir, "bazel-out")
// 				require.NoError(t, os.MkdirAll(bazelDir, 0755))
// 				oldWd, err := os.Getwd()
// 				require.NoError(t, err)
// 				require.NoError(t, os.Chdir(bazelDir))
// 				return bazelDir, func() {
// 					require.NoError(t, os.Chdir(oldWd))
// 				}
// 			},
// 			wantErr:     false,
// 			expectBazel: true,
// 		},
// 		{
// 			name: "no source root found",
// 			setup: func(t *testing.T) (string, func()) {
// 				tmpDir := t.TempDir()
// 				oldWd, err := os.Getwd()
// 				require.NoError(t, err)
// 				require.NoError(t, os.Chdir(tmpDir))
// 				return tmpDir, func() {
// 					require.NoError(t, os.Chdir(oldWd))
// 				}
// 			},
// 			wantErr: true,
// 			errMsg:  "source root not found",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, _ := telemetry.InitLogger(context.Background())
// 			_, cleanup := tt.setup(t)
// 			defer cleanup()

// 			got, err := getWorkingDirRelativeToSourceRoot(ctx)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				if tt.errMsg != "" {
// 					assert.Contains(t, err.Error(), tt.errMsg)
// 				}
// 				return
// 			}
// 			require.NoError(t, err)
// 			if tt.expectBazel {
// 				assert.Contains(t, got, "bazel-out")
// 			} else {
// 				assert.DirExists(t, filepath.Join(got, ".git"))
// 			}
// 		})
// 	}
// }
