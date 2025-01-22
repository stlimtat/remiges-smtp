package mail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultMailProcessorFactory(t *testing.T) {
	tests := []struct {
		name    string
		cfgs    []config.MailProcessorConfig
		wantErr bool
	}{
		{
			name: "happy",
			cfgs: []config.MailProcessorConfig{
				{
					Type:  UnixDosProcessorType,
					Args:  map[string]string{},
					Index: 0,
				},
			},
			wantErr: false,
		},
		{
			name:    "no processors",
			cfgs:    []config.MailProcessorConfig{},
			wantErr: true,
		},
		{
			name: "processor does not exist",
			cfgs: []config.MailProcessorConfig{
				{
					Type:  "processordoesnotexist",
					Args:  map[string]string{},
					Index: 0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := NewDefaultMailProcessorFactory(context.Background())
			require.NoError(t, err)

			processors, err := factory.NewMailProcessors(
				context.Background(),
				tt.cfgs,
			)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tt.cfgs), len(processors))
			}
		})
	}
}
