package file_mail

import (
	"context"
	"testing"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderFromTransformer(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.FileMailConfig
		headers  map[string][]byte
		wantFrom smtp.Address
		wantErr  bool
	}{
		{
			name: "happy - header",
			cfg: config.FileMailConfig{
				Type: HeaderFromTransformerType,
				Args: map[string]any{
					HeaderConfigArgType: config.ConfigTypeHeadersStr,
				},
			},
			headers: map[string][]byte{
				input.HeaderFromKey: []byte("test@example.com"),
			},
			wantFrom: smtp.Address{Localpart: "test", Domain: dns.Domain{ASCII: "example.com"}},
			wantErr:  false,
		},
		{
			name: "happy - default",
			cfg: config.FileMailConfig{
				Type: HeaderFromTransformerType,
				Args: map[string]any{
					HeaderConfigArgType:    config.ConfigTypeDefaultStr,
					HeaderConfigArgDefault: "default@example.com",
				},
			},
			headers: map[string][]byte{
				input.HeaderFromKey: []byte("test@example.com"),
			},
			wantFrom: smtp.Address{Localpart: "default", Domain: dns.Domain{ASCII: "example.com"}},
			wantErr:  false,
		},
		{
			name: "happy - long from header",
			cfg: config.FileMailConfig{
				Type: HeaderFromTransformerType,
				Args: map[string]any{
					HeaderConfigArgType: config.ConfigTypeHeadersStr,
				},
			},
			headers: map[string][]byte{
				input.HeaderFromKey: []byte("Name of user <test@example.com>"),
			},
			wantFrom: smtp.Address{Localpart: "test", Domain: dns.Domain{ASCII: "example.com"}},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			transformer := &HeaderFromTransformer{}
			err := transformer.Init(ctx, tt.cfg)
			require.NoError(t, err)

			got, err := transformer.Transform(
				ctx,
				&file.FileInfo{},
				&pmail.Mail{
					Metadata: tt.headers,
				},
			)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantFrom, got.From)
		})
	}
}
