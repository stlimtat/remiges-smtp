package input

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestRefreshList(t *testing.T) {
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
			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			// create the directory and file if requested
			if tt.createDirInPath {
				inPathInfo, err := os.Stat(tt.inPath)
				if err != nil || inPathInfo == nil {
					os.MkdirAll(tt.inPath, 0755)
				}
			}
			if tt.createDfFileInPath {
				dfFilePath := filepath.Join(tt.inPath, "df123")
				dfFileInfo, err := os.Stat(dfFilePath)
				if err != nil || dfFileInfo == nil {
					os.Create(dfFilePath)
				}
			}
			if tt.createQfFileInPath {
				qfFilePath := filepath.Join(tt.inPath, "qf123")
				qfFileInfo, err := os.Stat(qfFilePath)
				if err != nil || qfFileInfo == nil {
					os.Create(qfFilePath)
				}
			}

			fr, err := NewDefaultFileReader(ctx, tt.inPath)
			if tt.wantInitErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := fr.RefreshList(ctx)
				assert.Equal(t, tt.wantErr, err != nil)
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
			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			// create the directory and file if requested
			if tt.createDirInPath {
				inPathInfo, err := os.Stat(tt.inPath)
				if err != nil || inPathInfo == nil {
					os.MkdirAll(tt.inPath, 0755)
				}
			}
			if tt.createDfFileInPath {
				dfFilePath := filepath.Join(tt.inPath, "df123")
				dfFileInfo, err := os.Stat(dfFilePath)
				if err != nil || dfFileInfo == nil {
					os.Create(dfFilePath)
				}
			}
			if tt.createQfFileInPath {
				qfFilePath := filepath.Join(tt.inPath, "qf123")
				qfFileInfo, err := os.Stat(qfFilePath)
				if err != nil || qfFileInfo == nil {
					os.Create(qfFilePath)
				}
			}

			fr, err := NewDefaultFileReader(ctx, tt.inPath)
			if tt.wantInitErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := fr.RefreshList(ctx)
				assert.Equal(t, tt.wantErr, err != nil)
			}
		})
	}
}
