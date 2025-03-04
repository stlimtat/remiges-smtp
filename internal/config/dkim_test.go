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
		selector string
		dkim     *DKIMConfig
		wantDKIM moxConfig.DKIM
		wantErr  bool
	}{
		{
			name:     "default",
			selector: "key001",
			dkim:     DefaultDKIMConfig(ctx),
			wantDKIM: moxConfig.DKIM{
				Selectors: map[string]moxConfig.Selector{
					"key001": {
						Algorithm: "rsa",
						Canonicalization: moxConfig.Canonicalization{
							HeaderRelaxed: true,
							BodyRelaxed:   true,
						},
						Domain:            dns.Domain{ASCII: "key001"},
						DontSealHeaders:   true,
						Expiration:        "24h",
						ExpirationSeconds: int((time.Hour * 24).Seconds()),
						Hash:              "sha256",
						Headers:           nil,
					},
				},
				Sign: []string{"key001"},
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
			require.Contains(t, tt.dkim.DKIM.Selectors, tt.selector)
			gotSelector := tt.dkim.DKIM.Selectors[tt.selector]
			require.Contains(t, tt.wantDKIM.Selectors, tt.selector)
			wantSelector := tt.wantDKIM.Selectors[tt.selector]
			assert.Equal(t, wantSelector.Algorithm, gotSelector.Algorithm)
			assert.Equal(t, wantSelector.Canonicalization, gotSelector.Canonicalization)
			assert.Equal(t, wantSelector.Domain, gotSelector.Domain)
			assert.Equal(t, wantSelector.DontSealHeaders, gotSelector.DontSealHeaders)
			assert.Equal(t, wantSelector.Expiration, gotSelector.Expiration)
			assert.Equal(t, wantSelector.ExpirationSeconds, gotSelector.ExpirationSeconds)
			assert.Equal(t, wantSelector.Hash, gotSelector.Hash)
			assert.Equal(t, wantSelector.Hash, gotSelector.HashEffective)
			assert.Equal(t, wantSelector.Headers, gotSelector.Headers)
		})
	}
}
