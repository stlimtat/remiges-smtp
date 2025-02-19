package file_mail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMailTransformerFactory(t *testing.T) {
	tests := []struct {
		name           string
		cfgs           []config.FileMailConfig
		wantInitErr    bool
		wantProcessErr bool
	}{
		{
			name: "happy - single transformer",
			cfgs: []config.FileMailConfig{
				{
					Type: HeaderFromTransformerType,
					Args: map[string]string{
						HeaderConfigArgType:    config.ConfigTypeDefaultStr,
						HeaderConfigArgDefault: "test@example.com",
					},
					Index: 0,
				},
			},
			wantInitErr:    false,
			wantProcessErr: false,
		},
		{
			name: "happy - multiple transformers",
			cfgs: []config.FileMailConfig{
				{
					Type: HeaderFromTransformerType,
					Args: map[string]string{
						HeaderConfigArgType: config.ConfigTypeHeadersStr,
					},
					Index: 80,
				},
				{
					Type: HeaderToTransformerType,
					Args: map[string]string{
						HeaderToConfigArgType: "headers",
					},
					Index: 99,
				},
			},
			wantInitErr:    false,
			wantProcessErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			factory := NewMailTransformerFactory(ctx, tt.cfgs)

			err := factory.Init(ctx, config.FileMailConfig{})
			if tt.wantInitErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, len(tt.cfgs), len(factory.transformers))
		})
	}
}
