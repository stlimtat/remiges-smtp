package config

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/go-viper/mapstructure/v2"
	moxDkim "github.com/mjl-/mox/dkim"
	"github.com/mjl-/mox/dns"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDKIMConfig(t *testing.T) {
	tests := []struct {
		name         string
		selectorName string
		moxSelector  MoxSelector
		wantSelector moxDkim.Selector
		wantErr      bool
	}{
		{
			name:         "default",
			selectorName: "key001",
			moxSelector: MoxSelector{
				Algorithm:      "rsa",
				BodyRelaxed:    true,
				Expiration:     time.Hour * 24,
				Hash:           "sha256",
				HeaderRelaxed:  true,
				Headers:        []string{"from", "to", "subject", "date", "message-id", "content-type"},
				SealHeaders:    true,
				SelectorDomain: "key001",
			},
			wantSelector: moxDkim.Selector{
				BodyRelaxed:   true,
				Domain:        dns.Domain{ASCII: "key001"},
				Expiration:    time.Hour * 24,
				Hash:          "sha256",
				HeaderRelaxed: true,
				Headers:       []string{"from", "to", "subject", "date", "message-id", "content-type"},
				SealHeaders:   true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			dkimCfg := DKIMConfig{
				MoxSelectors: map[string]MoxSelector{
					tt.selectorName: tt.moxSelector,
				},
			}
			err := dkimCfg.Transform(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Contains(t, dkimCfg.Selectors, tt.selectorName)
			gotSelector := dkimCfg.Selectors[tt.selectorName]
			assert.Equal(t, tt.wantSelector.BodyRelaxed, gotSelector.BodyRelaxed)
			assert.Equal(t, tt.wantSelector.Domain, gotSelector.Domain)
			assert.Equal(t, tt.wantSelector.Expiration, gotSelector.Expiration)
			assert.Equal(t, tt.wantSelector.Hash, gotSelector.Hash)
			assert.Equal(t, tt.wantSelector.HeaderRelaxed, gotSelector.HeaderRelaxed)
			assert.Equal(t, tt.wantSelector.Headers, gotSelector.Headers)
			assert.Equal(t, tt.wantSelector.SealHeaders, gotSelector.SealHeaders)
		})
	}
}

func TestDKIMConfigFromViper(t *testing.T) {
	tests := []struct {
		name                string
		cfg                 []byte
		wantMoxSelectors    map[string]MoxSelector
		wantMapstructureErr bool
		wantViperErr        bool
	}{
		{
			name: "happy",
			cfg: []byte(`
selectors:
  key001:
    algorithm: rsa
    body-relaxed: true
    expiration: 72h
    hash: sha256
    header-relaxed: true
    headers:
      - from
      - to
      - subject
      - date
      - message-id
      - content-type
    private-key-file: /tmp/file001.pem
    seal-headers: false
    selector-domain: key001
`),
			wantMoxSelectors: map[string]MoxSelector{
				"key001": {
					Algorithm:      "rsa",
					BodyRelaxed:    true,
					Expiration:     72 * time.Hour,
					Hash:           "sha256",
					HeaderRelaxed:  true,
					Headers:        []string{"from", "to", "subject", "date", "message-id", "content-type"},
					PrivateKeyFile: "/tmp/file001.pem",
					SealHeaders:    false,
					SelectorDomain: "key001",
				},
			},
			wantMapstructureErr: false,
			wantViperErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.SetConfigType("yaml")
			err := viper.ReadConfig(bytes.NewBuffer(tt.cfg))
			require.NoError(t, err)

			settings := viper.AllSettings()
			assert.Contains(t, settings, "selectors")

			mapGotDkimCfg := DKIMConfig{}
			decoder, err := mapstructure.NewDecoder(
				&mapstructure.DecoderConfig{
					Metadata:   nil,
					DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
					Result:     &mapGotDkimCfg,
				},
			)
			require.NoError(t, err)
			err = decoder.Decode(settings)
			if tt.wantMapstructureErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantMoxSelectors, mapGotDkimCfg.MoxSelectors)

			gotDkimCfg := DKIMConfig{}
			err = viper.Unmarshal(&gotDkimCfg)
			if tt.wantViperErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantMoxSelectors, gotDkimCfg.MoxSelectors)
		})
	}
}
