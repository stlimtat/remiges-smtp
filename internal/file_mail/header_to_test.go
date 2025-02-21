package file_mail

import (
	"context"
	"testing"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderToTransformer(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.FileMailConfig
		header  map[string][]byte
		wantTo  []smtp.Address
		wantErr bool
	}{
		{
			name: "happy - default",
			cfg: config.FileMailConfig{
				Type: HeaderToTransformerType,
				Args: map[string]any{
					HeaderConfigArgType:    config.ConfigTypeDefaultStr,
					HeaderConfigArgDefault: "default@example.com",
				},
			},
			header: map[string][]byte{
				input.HeaderToKey: []byte("test@example.com"),
			},
			wantTo: []smtp.Address{
				{Localpart: "default", Domain: dns.Domain{ASCII: "example.com"}},
			},
			wantErr: false,
		},
		{
			name: "happy - headers",
			cfg: config.FileMailConfig{
				Type: HeaderToTransformerType,
				Args: map[string]any{
					HeaderConfigArgType: config.ConfigTypeHeadersStr,
				},
			},
			header: map[string][]byte{
				input.HeaderToKey: []byte("test@example.com"),
			},
			wantTo: []smtp.Address{
				{Localpart: "test", Domain: dns.Domain{ASCII: "example.com"}},
			},
			wantErr: false,
		},
		{
			name: "happy - multiple headers",
			cfg: config.FileMailConfig{
				Type: HeaderToTransformerType,
				Args: map[string]any{
					HeaderConfigArgType: config.ConfigTypeHeadersStr,
				},
			},
			header: map[string][]byte{
				input.HeaderToKey: []byte("Example User <test1@example.com>, Example User 2 <test2@example.com>"),
			},
			wantTo: []smtp.Address{
				{Localpart: "test1", Domain: dns.Domain{ASCII: "example.com"}},
				{Localpart: "test2", Domain: dns.Domain{ASCII: "example.com"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			transformer := &HeaderToTransformer{}
			err := transformer.Init(ctx, tt.cfg)
			require.NoError(t, err)
			got, err := transformer.Transform(
				ctx,
				&file.FileInfo{},
				&mail.Mail{
					Metadata: tt.header,
				},
			)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantTo, got.To)
		})
	}
}
