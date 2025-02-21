package mail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultMailProcessorFactory(t *testing.T) {
	tests := []struct {
		name           string
		cfgs           []config.MailProcessorConfig
		wantInitErr    bool
		wantProcessErr bool
	}{
		{
			name: "happy - single processor",
			cfgs: []config.MailProcessorConfig{
				{
					Type:  UnixDosProcessorType,
					Args:  map[string]any{},
					Index: 0,
				},
			},
			wantInitErr:    false,
			wantProcessErr: false,
		},
		{
			name: "happy - index not in sequence",
			cfgs: []config.MailProcessorConfig{
				{
					Type:  UnixDosProcessorType,
					Args:  map[string]any{},
					Index: 50,
				},
				{
					Type:  BodyHeadersProcessorType,
					Args:  map[string]any{},
					Index: 99,
				},
			},
			wantInitErr:    false,
			wantProcessErr: false,
		},
		{
			name:           "no processors",
			cfgs:           []config.MailProcessorConfig{},
			wantInitErr:    true,
			wantProcessErr: false,
		},
		{
			name: "processor does not exist",
			cfgs: []config.MailProcessorConfig{
				{
					Type:  "processordoesnotexist",
					Args:  map[string]any{},
					Index: 0,
				},
			},
			wantInitErr:    true,
			wantProcessErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			factory, err := NewDefaultMailProcessorFactory(ctx, tt.cfgs)
			require.NoError(t, err)

			err = factory.Init(ctx, config.MailProcessorConfig{})
			if tt.wantInitErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, len(tt.cfgs), len(factory.processors))
		})
	}
}
