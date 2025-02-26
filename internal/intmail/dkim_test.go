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
  dkim:
    selectors:
      google:
        domain: "stlim.net"
      blah:
        domain: "blah.com"
`),
			wantDKIMConfig: config.DKIMConfig{
				DKIM: moxConfig.DKIM{
					Selectors: map[string]moxConfig.Selector{
						"google": {
							Domain: dns.Domain{ASCII: "stlim.net"},
						},
					},
					Sign: []string{},
				},
				MoxSelectors: map[string]config.MoxSelector{
					"google": {
						Domain: "stlim.net",
					},
				},
				MoxSign: []string{},
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
			assert.Subset(t, processor.DkimCfg.Selectors, tt.wantDKIMConfig.Selectors)
			assert.Subset(t, processor.DkimCfg.Sign, tt.wantDKIMConfig.Sign)
			assert.Subset(t, processor.DkimCfg.MoxSelectors, tt.wantDKIMConfig.MoxSelectors)
			assert.Subset(t, processor.DkimCfg.MoxSign, tt.wantDKIMConfig.MoxSign)
		})
	}
}
