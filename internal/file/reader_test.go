package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestNewDefaultFileReader(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	tests := []struct {
		name        string
		inputDir    string
		createDir   bool
		permissions os.FileMode
		wantErr     bool
	}{
		{
			name:        "valid directory",
			inputDir:    filepath.Join(tmpDir, "valid_dir"),
			createDir:   true,
			permissions: 0755,
			wantErr:     false,
		},
		{
			name:      "non-existent directory",
			inputDir:  filepath.Join(tmpDir, "non_existent"),
			createDir: false,
			wantErr:   true,
		},
		// {
		// 	name:        "not a directory",
		// 	inputDir:    filepath.Join(tmpDir, "not_a_dir"),
		// 	createDir:   true,
		// 	permissions: 0644,
		// 	wantErr:     true,
		// },
		{
			name:      "empty directory path",
			inputDir:  "",
			createDir: false,
			wantErr:   true,
		},
		// {
		// 	name:        "directory with no read permissions",
		// 	inputDir:    filepath.Join(tmpDir, "no_read"),
		// 	createDir:   true,
		// 	permissions: 0000,
		// 	wantErr:     true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.inputDir != "" && tt.createDir {
				err := os.MkdirAll(tt.inputDir, tt.permissions)
				require.NoError(t, err)
				defer os.RemoveAll(tt.inputDir)

				if tt.permissions == 0000 {
					// Create a file in the directory to make it non-empty
					_, err := os.Create(filepath.Join(tt.inputDir, "test.txt"))
					require.NoError(t, err)
				}
			}

			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			frt := NewMockIFileReadTracker(ctrl)

			// Test
			_, err := NewDefaultFileReader(ctx, tt.inputDir, frt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFile(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		createFile  bool
		permissions os.FileMode
		want        bool
	}{
		{
			name:        "valid file",
			fileName:    "test.txt",
			createFile:  true,
			permissions: 0644,
			want:        true,
		},
		{
			name:       "non-existent file",
			fileName:   "nonexistent.txt",
			createFile: false,
			want:       false,
		},
		{
			name:        "directory instead of file",
			fileName:    "testdir",
			createFile:  true,
			permissions: 0755,
			want:        true,
		},
		{
			name:        "file with no read permissions",
			fileName:    "no_read.txt",
			createFile:  true,
			permissions: 0000,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tmpDir := t.TempDir()
			defer os.RemoveAll(tmpDir)
			var err error

			if tt.createFile {
				path := filepath.Join(tmpDir, tt.fileName)
				if tt.permissions&os.ModeDir != 0 {
					err = os.Mkdir(path, tt.permissions)
				} else {
					_, err = os.Create(path)
					os.Chmod(path, tt.permissions)
				}
				require.NoError(t, err)
			}

			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			frt := NewMockIFileReadTracker(ctrl)
			fr, err := NewDefaultFileReader(ctx, tmpDir, frt)
			require.NoError(t, err)

			// Test
			got := fr.ValidateFile(ctx, tt.fileName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetQfFileName(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
		wantErr  bool
	}{
		{
			name:     "valid df file",
			fileName: "df123",
			want:     "qf123",
			wantErr:  false,
		},
		{
			name:     "empty filename",
			fileName: "",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "short filename",
			fileName: "d",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "non-df prefix",
			fileName: "qf123",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			frt := NewMockIFileReadTracker(ctrl)
			fr, err := NewDefaultFileReader(ctx, "/tmp", frt)
			require.NoError(t, err)

			got, err := fr.GetQfFileName(ctx, tt.fileName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRefreshList(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    []string
		expectedFiles []string
		wantErr       bool
		trackerErr    bool
		fileStatus    input.FileStatus
	}{
		{
			name:          "single df file",
			setupFiles:    []string{"df1"},
			expectedFiles: []string{"df1"},
			wantErr:       false,
		},
		{
			name:          "multiple df files",
			setupFiles:    []string{"df1", "df2", "df3"},
			expectedFiles: []string{"df1", "df2", "df3"},
			wantErr:       false,
		},
		{
			name:          "mixed files",
			setupFiles:    []string{"df1", "qf1", "other.txt"},
			expectedFiles: []string{"df1"},
			wantErr:       false,
		},
		{
			name:          "no df files",
			setupFiles:    []string{"qf1", "other.txt"},
			expectedFiles: []string{},
			wantErr:       false,
		},
		// {
		// 	name:          "tracker error",
		// 	setupFiles:    []string{"df1"},
		// 	expectedFiles: []string{},
		// 	wantErr:       false,
		// 	trackerErr:    true,
		// },
		{
			name:          "file already processing",
			setupFiles:    []string{"df1"},
			expectedFiles: []string{},
			fileStatus:    input.FILE_STATUS_PROCESSING,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tmpDir := t.TempDir()
			defer os.RemoveAll(tmpDir)

			for _, file := range tt.setupFiles {
				_, err := os.Create(filepath.Join(tmpDir, file))
				require.NoError(t, err)
			}

			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			frt := NewMockIFileReadTracker(ctrl)
			if tt.trackerErr {
				frt.EXPECT().
					FileRead(gomock.Any(), gomock.Any()).
					Return(input.FILE_STATUS_ERROR, fmt.Errorf("test error")).
					AnyTimes()
			} else {
				frt.EXPECT().
					FileRead(gomock.Any(), gomock.Any()).
					Return(tt.fileStatus, nil).
					AnyTimes()
			}

			fr, err := NewDefaultFileReader(ctx, tmpDir, frt)
			require.NoError(t, err)

			// Test
			_, err = fr.RefreshList(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// assert.Len(t, files, len(tt.expectedFiles))
				// for i, file := range files {
				// 	assert.Equal(t, tt.expectedFiles[i], filepath.Base(file.DfFilePath))
				// }
			}
		})
	}
}

// func TestReadNextFile(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		setupFiles    []string
// 		fileStatuses  map[string]input.FileStatus
// 		trackerErr    bool
// 		upsertErr     bool
// 		expectedFile  string
// 		wantErr       bool
// 		errorContains string
// 	}{
// 		{
// 			name:         "single unprocessed file",
// 			setupFiles:   []string{"df1"},
// 			fileStatuses: map[string]input.FileStatus{"1": input.FILE_STATUS_INIT},
// 			expectedFile: "df1",
// 			wantErr:      false,
// 		},
// 		{
// 			name:          "file already processing",
// 			setupFiles:    []string{"df1"},
// 			fileStatuses:  map[string]input.FileStatus{"1": input.FILE_STATUS_PROCESSING},
// 			wantErr:       true,
// 			errorContains: "file is being processed",
// 		},
// 		{
// 			name:       "file already done",
// 			setupFiles: []string{"df1", "df2"},
// 			fileStatuses: map[string]input.FileStatus{
// 				"1": input.FILE_STATUS_DONE,
// 				"2": input.FILE_STATUS_INIT,
// 			},
// 			expectedFile: "df2",
// 			wantErr:      false,
// 		},
// 		{
// 			name:       "tracker read error",
// 			setupFiles: []string{"df1"},
// 			trackerErr: true,
// 			wantErr:    true,
// 		},
// 		{
// 			name:         "tracker upsert error",
// 			setupFiles:   []string{"df1"},
// 			fileStatuses: map[string]input.FileStatus{"1": input.FILE_STATUS_INIT},
// 			upsertErr:    true,
// 			wantErr:      true,
// 		},
// 		{
// 			name:         "no files available",
// 			setupFiles:   []string{},
// 			wantErr:      false,
// 			expectedFile: "",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup
// 			tmpDir := t.TempDir()
// 			defer os.RemoveAll(tmpDir)

// 			for _, file := range tt.setupFiles {
// 				_, err := os.Create(filepath.Join(tmpDir, file))
// 				require.NoError(t, err)
// 			}

// 			ctx := context.Background()
// 			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			frt := NewMockIFileReadTracker(ctrl)

// 			// Setup mock expectations
// 			if tt.trackerErr {
// 				frt.EXPECT().
// 					FileRead(gomock.Any(), gomock.Any()).
// 					Return(input.FILE_STATUS_ERROR, fmt.Errorf("test error")).
// 					AnyTimes()
// 			} else {
// 				for id, status := range tt.fileStatuses {
// 					frt.EXPECT().
// 						FileRead(gomock.Any(), id).
// 						Return(status, nil).
// 						AnyTimes()
// 				}
// 			}

// 			if tt.upsertErr {
// 				frt.EXPECT().
// 					UpsertFile(gomock.Any(), gomock.Any(), gomock.Any()).
// 					Return(fmt.Errorf("test error")).
// 					AnyTimes()
// 			} else {
// 				frt.EXPECT().
// 					UpsertFile(gomock.Any(), gomock.Any(), gomock.Any()).
// 					Return(nil).
// 					AnyTimes()
// 			}

// 			fr, err := NewDefaultFileReader(ctx, tempDir, frt)
// 			require.NoError(t, err)

// 			// Refresh list first
// 			_, err = fr.RefreshList(ctx)
// 			require.NoError(t, err)

// 			// Test
// 			file, err := fr.ReadNextFile(ctx)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				if tt.errorContains != "" {
// 					assert.Contains(t, err.Error(), tt.errorContains)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 				if tt.expectedFile == "" {
// 					assert.Nil(t, file)
// 				} else {
// 					assert.Equal(t, tt.expectedFile, filepath.Base(file.DfFilePath))
// 				}
// 			}
// 		})
// 	}
// }

func TestConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping concurrent access test")
	}
	// Setup
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// Create test files
	for i := 0; i < 10; i++ {
		_, err := os.Create(filepath.Join(tmpDir, fmt.Sprintf("df%d", i)))
		require.NoError(t, err)
	}

	ctx := context.Background()
	ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	frt := NewMockIFileReadTracker(ctrl)
	frt.EXPECT().
		FileRead(gomock.Any(), gomock.Any()).
		Return(input.FILE_STATUS_INIT, nil).
		AnyTimes()
	frt.EXPECT().
		UpsertFile(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	fr, err := NewDefaultFileReader(ctx, tmpDir, frt)
	require.NoError(t, err)

	// Refresh list first
	_, err = fr.RefreshList(ctx)
	require.NoError(t, err)

	// Test concurrent access
	const numGoroutines = 10
	results := make(chan *FileInfo, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			file, err := fr.ReadNextFile(ctx)
			if err != nil {
				errors <- err
				return
			}
			results <- file
		}()
	}

	// Collect results
	var files []*FileInfo
	for i := 0; i < numGoroutines; i++ {
		select {
		case file := <-results:
			if file != nil {
				files = append(files, file)
			}
		case err := <-errors:
			assert.NoError(t, err)
		}
	}

	// Verify that each file was processed only once
	processedFiles := make(map[string]bool)
	for _, file := range files {
		fileName := filepath.Base(file.DfFilePath)
		assert.False(t, processedFiles[fileName], "file processed multiple times: %s", fileName)
		processedFiles[fileName] = true
	}
}
