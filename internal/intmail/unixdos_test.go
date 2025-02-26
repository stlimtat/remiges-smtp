package intmail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnixDosProcessor(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		wantBody []byte
		wantErr  bool
	}{
		{
			name:     "happy1",
			body:     []byte("Hello\nWorld"),
			wantBody: []byte("Hello\r\nWorld"),
			wantErr:  false,
		},
		{
			name:     "happy2",
			body:     []byte("Hello\r\nWorld"),
			wantBody: []byte("Hello\r\nWorld"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			processor := UnixDosProcessor{}
			err := processor.Init(ctx, config.MailProcessorConfig{})
			require.NoError(t, err)
			got, err := processor.Process(ctx, &mail.Mail{Body: tt.body})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantBody, got.Body)
		})
	}
}
