package intmail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeHeaders(t *testing.T) {
	tests := []struct {
		name        string
		headersMap  map[string][]byte
		wantHeaders []byte
		wantErr     bool
	}{
		{
			name: "test1",
			headersMap: map[string][]byte{
				"From": []byte("sender@example.com"),
				"To":   []byte("recipient@example.com"),
			},
			wantHeaders: []byte("From: sender@example.com\r\nTo: recipient@example.com\r\n\r\n"),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			processor := &MergeHeadersProcessor{}
			err := processor.Init(ctx, config.MailProcessorConfig{})
			require.NoError(t, err)

			got, err := processor.Process(ctx, &pmail.Mail{
				HeadersMap: tt.headersMap,
			})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantHeaders, got.Headers)
		})
	}
}
