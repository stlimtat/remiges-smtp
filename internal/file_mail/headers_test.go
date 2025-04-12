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

func TestHeadersTransformer(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.FileMailConfig
		headers  []byte
		wantMail *pmail.Mail
		wantErr  bool
	}{
		{
			name: "happy - default headers",
			cfg: config.FileMailConfig{
				Type: HeadersTransformerType,
				Args: map[string]any{},
			},
			headers: []byte("From: test@example.com\nTo: test@example.com\nSubject: test"),
			wantMail: &pmail.Mail{
				Metadata: map[string][]byte{
					"From":    []byte("test@example.com"),
					"To":      []byte("test@example.com"),
					"Subject": []byte("test"),
				},
			},
			wantErr: false,
		},
		{
			name: "happy - multiple line header",
			cfg: config.FileMailConfig{
				Type: HeadersTransformerType,
				Args: map[string]any{},
			},
			headers: []byte("Content-Type: text/plain;\n\tcharset=utf-8\n"),
			wantMail: &pmail.Mail{
				Metadata: map[string][]byte{
					"Content-Type": []byte("text/plain;charset=utf-8"),
				},
			},
			wantErr: false,
		},
		{
			name: "happy - with prefix",
			cfg: config.FileMailConfig{
				Type: HeadersTransformerType,
				Args: map[string]any{
					HeadersConfigArgPrefix: "H??",
				},
			},
			headers: []byte("H??From: test1@example.com\nH??To: test1@example.com\nH??Subject: test1\n"),
			wantMail: &pmail.Mail{
				Metadata: map[string][]byte{
					"From":       []byte("test1@example.com"),
					"H??From":    []byte("test1@example.com"),
					"H??To":      []byte("test1@example.com"),
					"H??Subject": []byte("test1"),
					"To":         []byte("test1@example.com"),
					"Subject":    []byte("test1"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			transformer := &HeadersTransformer{}
			err := transformer.Init(ctx, tt.cfg)
			require.NoError(t, err)

			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.qf")
			err = os.WriteFile(tmpFile, tt.headers, 0644)
			require.NoError(t, err)
			defer os.Remove(tmpFile)
			fileInfo := &file.FileInfo{
				QfFilePath: tmpFile,
			}

			got, err := transformer.Transform(ctx, fileInfo, &pmail.Mail{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantMail.Metadata, got.Metadata)
		})
	}
}
