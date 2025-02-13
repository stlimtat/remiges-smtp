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
		fileInfo *file.FileInfo
		wantMail *mail.Mail
		wantErr  bool
	}{
		{
			name: "happy",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			transformer := &HeadersTransformer{}
			err := transformer.Init(ctx, config.FileMailConfig{})
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
