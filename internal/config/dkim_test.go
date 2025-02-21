package config

import (
	"context"
	"testing"
	"time"

	moxConfig "github.com/mjl-/mox/config"
	"github.com/mjl-/mox/dns"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDKIMConfig(t *testing.T) {
	ctx, _ := telemetry.InitLogger(context.Background())
	tests := []struct {
		name     string
		dkim     *DKIMConfig
		wantDKIM moxConfig.DKIM
		wantErr  bool
	}{
		{
			name: "default",
			dkim: DefaultDKIMConfig(ctx),
			wantDKIM: moxConfig.DKIM{
				Selectors: map[string]moxConfig.Selector{
					"google": {
						Algorithm: "rsa-sha256",
						Canonicalization: moxConfig.Canonicalization{
							HeaderRelaxed: true,
							BodyRelaxed:   true,
						},
						Domain:            dns.Domain{ASCII: "stlim.net"},
						DontSealHeaders:   true,
						Expiration:        "24h",
						ExpirationSeconds: int((time.Hour * 24).Seconds()),
						Hash:              "sha256",
						Headers:           nil,
					},
				},
				Sign: []string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dkim.Transform(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantDKIM, tt.dkim.DKIM)
		})
	}
}
