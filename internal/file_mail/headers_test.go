package file_mail

import (
	"bytes"
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersTransformer(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.FileMailConfig
		fileInfo *file.FileInfo
		wantMail *mail.Mail
		wantErr  bool
	}{
		{
			name: "happy - default headers",
			cfg: config.FileMailConfig{
				Type: HeadersTransformerType,
				Args: map[string]string{},
			},
			fileInfo: &file.FileInfo{
				ID:         "1",
				QfFilePath: "testdata/test.qf",
				QfReader:   bytes.NewReader([]byte("From: test@example.com\nTo: test@example.com\nSubject: test")),
			},
			wantMail: &mail.Mail{
				Headers: map[string][]byte{
					"From":    []byte("test@example.com"),
					"To":      []byte("test@example.com"),
					"Subject": []byte("test"),
				},
			},
			wantErr: false,
		},
		{
			name: "happy - with prefix",
			cfg: config.FileMailConfig{
				Type: HeadersTransformerType,
				Args: map[string]string{
					HeadersConfigArgPrefix: "H??",
				},
			},
			fileInfo: &file.FileInfo{
				ID:         "1",
				QfFilePath: "testdata/test.qf",
				QfReader:   bytes.NewReader([]byte("H??From: test1@example.com\nH??To: test1@example.com\nH??Subject: test1\n")),
			},
			wantMail: &mail.Mail{
				Headers: map[string][]byte{
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

			got, err := transformer.Transform(ctx, tt.fileInfo, &mail.Mail{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantMail, got)
		})
	}
}
