package output

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOutputs(t *testing.T) {
	tests := []struct {
		name    string
		cfgs    []config.OutputConfig
		wantLen int
		wantErr bool
	}{
		{
			name: "happy - file",
			cfgs: []config.OutputConfig{
				{
					Type: ConfigOutputTypeFile,
					Args: map[string]any{
						ConfigArgPath: "/tmp",
					},
				},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "alternate - file path does not exist",
			cfgs: []config.OutputConfig{
				{
					Type: ConfigOutputTypeFile,
					Args: map[string]any{
						ConfigArgPath: "/tmp/does-not-exist",
					},
				},
			},
			wantLen: 1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			factory := OutputFactory{}
			got, err := factory.NewOutputs(ctx, tt.cfgs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
		})
	}
}
