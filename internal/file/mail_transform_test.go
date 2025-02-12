package file

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMailTransformer(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.ReadFileConfig
		fileInfo *FileInfo
		wantMail *mail.Mail
		wantErr  bool
	}{
		{
			name: "happy header from",
			cfg: config.ReadFileConfig{
				FromType:    config.FromTypeHeaders,
				DefaultFrom: "defaultFrom@example.com",
			},
			fileInfo: &FileInfo{
				DfReader: bytes.NewReader([]byte(
					`Test Body`,
				)),
				ID: "123",
				QfReader: bytes.NewReader([]byte(
					`From: sender@example.com
To: recipient@example.com
Subject: test
`,
				)),
				Status: input.FILE_STATUS_INIT,
			},
			wantMail: &mail.Mail{
				From: smtp.Address{
					Localpart: "sender",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				To: smtp.Address{
					Localpart: "recipient",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				Body: []byte("Test Body"),
			},
			wantErr: false,
		}, {
			name: "happy default from",
			cfg: config.ReadFileConfig{
				FromType:    config.FromTypeDefault,
				DefaultFrom: "defaultFrom@example.com",
			},
			fileInfo: &FileInfo{
				DfReader: bytes.NewReader([]byte(
					`Test Body`,
				)),
				ID: "123",
				QfReader: bytes.NewReader([]byte(
					`From: sender@example.com
To: recipient@example.com
Subject: test
					`,
				)),
				Status: input.FILE_STATUS_INIT,
			},
			wantMail: &mail.Mail{
				From: smtp.Address{
					Localpart: "defaultFrom",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				To: smtp.Address{
					Localpart: "recipient",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				Body: []byte("Test Body"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)
			mailTransformer := NewMailTransformer(ctx, tt.cfg)
			got, err := mailTransformer.Transform(ctx, tt.fileInfo)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantMail.From, got.From)
				assert.Equal(t, tt.wantMail.To, got.To)
				// assert.Equal(t, tt.wantMail.Body, got.Body)
			}
		})
	}
}

func TestMailTransformer_ReadBody(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.ReadFileConfig
		fileInfo *FileInfo
		want     []byte
		wantErr  bool
	}{
		{
			name: "happy",
			cfg: config.ReadFileConfig{
				FromType:    config.FromTypeHeaders,
				DefaultFrom: "defaultFrom@example.com",
			},
			fileInfo: &FileInfo{
				DfReader: bytes.NewReader([]byte("Test Body")),
			},
			want:    []byte("Test Body"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)
			mailTransformer := NewMailTransformer(ctx, tt.cfg)
			got, err := mailTransformer.ReadBody(ctx, tt.fileInfo)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMailTransformer_ReadHeaders(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.ReadFileConfig
		fileInfo *FileInfo
		want     map[string][]byte
		wantErr  bool
	}{
		{
			name: "happy - from headers",
			cfg: config.ReadFileConfig{
				FromType:    config.FromTypeHeaders,
				DefaultFrom: "defaultFrom@example.com",
			},
			fileInfo: &FileInfo{
				QfReader: bytes.NewReader([]byte(
					`From: sender@example.com
To: recipient@example.com
Subject: test
`,
				)),
			},
			want: map[string][]byte{
				"From":    []byte("sender@example.com"),
				"To":      []byte("recipient@example.com"),
				"Subject": []byte("test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)

			mailTransformer := NewMailTransformer(ctx, tt.cfg)
			got, err := mailTransformer.ReadHeaders(ctx, tt.fileInfo)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
