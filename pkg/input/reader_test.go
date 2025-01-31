package input

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestRefreshList(t *testing.T) {
	var tests = []struct {
		name               string
		inPath             string
		createDirInPath    bool
		createDfFileInPath bool
		createQfFileInPath bool
		wantInitErr        bool
		wantTrackerErr     bool
		wantRefreshErr     bool
	}{
		{"happy", "/tmp", true, true, true, false, false, false},
		{"no-file-exists", "/tmp/no-file-exists", true, false, false, false, true, true},
		{"not-a-directory", "/tmp/not-a-directory", false, false, false, true, false, false},
		{"no-df-file", "/tmp/no-df-file", true, false, true, false, false, false},
		{"no-qf-file", "/tmp/no-qf-file", true, true, false, false, false, false},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			// create the directory and file if requested
			if tt.createDirInPath {
				inPathInfo, err := os.Stat(tt.inPath)
				if err != nil || inPathInfo == nil {
					_ = os.MkdirAll(tt.inPath, 0755)
				}
			}
			if tt.createDfFileInPath {
				dfFilePath := filepath.Join(tt.inPath, "df123")
				dfFileInfo, err := os.Stat(dfFilePath)
				if err != nil || dfFileInfo == nil {
					_, _ = os.Create(dfFilePath)
				}
			}
			if tt.createQfFileInPath {
				qfFilePath := filepath.Join(tt.inPath, "qf123")
				qfFileInfo, err := os.Stat(qfFilePath)
				if err != nil || qfFileInfo == nil {
					_, _ = os.Create(qfFilePath)
				}
			}

			frt := NewMockIFileReadTracker(ctrl)
			frt.EXPECT().
				UpsertFile(gomock.Any(), "123", FILE_STATUS_INIT).
				Return(nil).
				AnyTimes()
			if tt.wantTrackerErr {
				frt.EXPECT().
					FileRead(gomock.Any(), "123").
					Return(FILE_STATUS_ERROR, fmt.Errorf("test error")).
					AnyTimes()
			} else {
				frt.EXPECT().
					FileRead(gomock.Any(), "123").
					Return(FILE_STATUS_INIT, nil).
					AnyTimes()
			}

			fr, err := NewDefaultFileReader(ctx, tt.inPath, frt)
			if tt.wantInitErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				_, err := fr.RefreshList(ctx)
				assert.Equal(t, tt.wantRefreshErr, err != nil)
			}
		})
	}
}

func TestReadNextFile(t *testing.T) {
	var tests = []struct {
		name               string
		inPath             string
		createDirInPath    bool
		createDfFileInPath bool
		createQfFileInPath bool
		wantInitErr        bool
		wantErr            bool
	}{
		{"happy", "/tmp", true, true, true, false, false},
		{"no-file-exists", "/tmp/no-file-exists", true, false, false, false, true},
		{"not-a-directory", "/tmp/not-a-directory", false, false, false, true, false},
		{"no-df-file", "/tmp/no-df-file", true, false, true, false, false},
		{"no-qf-file", "/tmp/no-qf-file", true, true, false, false, false},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			// create the directory and file if requested
			if tt.createDirInPath {
				inPathInfo, err := os.Stat(tt.inPath)
				if err != nil || inPathInfo == nil {
					_ = os.MkdirAll(tt.inPath, 0755)
				}
			}
			if tt.createDfFileInPath {
				dfFilePath := filepath.Join(tt.inPath, "df123")
				dfFileInfo, err := os.Stat(dfFilePath)
				if err != nil || dfFileInfo == nil {
					_, _ = os.Create(dfFilePath)
				}
			}
			if tt.createQfFileInPath {
				qfFilePath := filepath.Join(tt.inPath, "qf123")
				qfFileInfo, err := os.Stat(qfFilePath)
				if err != nil || qfFileInfo == nil {
					_, _ = os.Create(qfFilePath)
				}
			}

			frt := NewMockIFileReadTracker(ctrl)
			frt.EXPECT().
				UpsertFile(gomock.Any(), "123", FILE_STATUS_INIT).
				Return(nil).
				AnyTimes()

			fr, err := NewDefaultFileReader(ctx, tt.inPath, frt)
			if tt.wantInitErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				_, err := fr.RefreshList(ctx)
				assert.Equal(t, tt.wantErr, err != nil)
			}
		})
	}
}
