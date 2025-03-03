package intmail

import (
	"bytes"
	"context"
	"testing"

	moxConfig "github.com/mjl-/mox/config"
	"github.com/mjl-/mox/dns"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDKIMProcessorInit(t *testing.T) {
	tests := []struct {
		name           string
		cfgStr         []byte
		wantDKIMConfig config.DKIMConfig
	}{
		{
			name: "happy",
			cfgStr: []byte(`
args:
  domain-str: stlim.net
  dkim:
    selectors:
      key001:
        domain: stlim.net
      key002:
        domain: blah.com
    sign:
      - key001
      - key002
`),
			wantDKIMConfig: config.DKIMConfig{
				DKIM: moxConfig.DKIM{
					Selectors: map[string]moxConfig.Selector{
						"key001": {
							Domain: dns.Domain{ASCII: "stlim.net"},
						},
						"key002": {
							Domain: dns.Domain{ASCII: "blah.com"},
						},
					},
					Sign: []string{"key001", "key002"},
				},
				MoxSelectors: map[string]config.MoxSelector{
					"key001": {
						Domain: "stlim.net",
					},
					"key002": {
						Domain: "blah.com",
					},
				},
				MoxSign: []string{"key001", "key002"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			cfg := config.MailProcessorConfig{}
			viper.SetConfigType("yaml")
			err := viper.ReadConfig(bytes.NewBuffer(tt.cfgStr))
			require.NoError(t, err)
			settings := viper.AllSettings()
			assert.Contains(t, settings, "args")
			err = viper.Unmarshal(&cfg)
			require.NoError(t, err)
			processor := &DKIMProcessor{}
			err = processor.Init(ctx, cfg)
			require.NoError(t, err)
			dkimCfg := processor.DomainCfg.DKIM
			assert.Subset(t, dkimCfg.Selectors, tt.wantDKIMConfig.Selectors)
			assert.Subset(t, dkimCfg.Sign, tt.wantDKIMConfig.Sign)
			assert.Subset(t, dkimCfg.MoxSelectors, tt.wantDKIMConfig.MoxSelectors)
			assert.Subset(t, dkimCfg.MoxSign, tt.wantDKIMConfig.MoxSign)
		})
	}
}
