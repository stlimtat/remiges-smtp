package file_mail

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyTransformer(t *testing.T) {
	tests := []struct {
		name             string
		cfg              config.FileMailConfig
		body             []byte
		wantBody         []byte
		wantInitErr      bool
		wantTransformErr bool
	}{
		{
			name:             "happy - simple body",
			cfg:              config.FileMailConfig{},
			body:             []byte("test body"),
			wantBody:         []byte("test body"),
			wantInitErr:      false,
			wantTransformErr: false,
		},
		{
			name: "happy - mime",
			cfg:  config.FileMailConfig{},
			body: []byte(`------=_Part_123
Content-Type: text/plain
Content-Transfer-Encoding: 7bit

Hello
World
------=_Part_123--
			`),
			wantBody:         []byte("------=_Part_123\r\nContent-Type: text/plain\r\nContent-Transfer-Encoding: 7bit\r\n\r\nHello\r\nWorld\r\n------=_Part_123--"),
			wantInitErr:      false,
			wantTransformErr: false,
		},
		{
			name: "happy - mime with first line as new line",
			cfg:  config.FileMailConfig{},
			body: []byte(`
------=_Part_123
Content-Type: text/plain
Content-Transfer-Encoding: 7bit

Hello
World
------=_Part_123--
			`),
			wantBody:         []byte("------=_Part_123\r\nContent-Type: text/plain\r\nContent-Transfer-Encoding: 7bit\r\n\r\nHello\r\nWorld\r\n------=_Part_123--"),
			wantInitErr:      false,
			wantTransformErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			transformer := &BodyTransformer{}
			err := transformer.Init(ctx, tt.cfg)
			if tt.wantInitErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.txt")
			err = os.WriteFile(tmpFile, tt.body, 0644)
			require.NoError(t, err)
			defer os.Remove(tmpFile)
			fileInfo := &file.FileInfo{
				DfFilePath: tmpFile,
			}
			gotMail, err := transformer.Transform(ctx, fileInfo, &pmail.Mail{})
			if tt.wantTransformErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, gotMail.Body)
		})
	}
}
