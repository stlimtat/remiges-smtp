package mail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnixDosProcessor(t *testing.T) {
	tests := []struct {
		name     string
		mail     *Mail
		wantMail *Mail
		wantErr  bool
	}{
		{
			name:     "happy1",
			mail:     &Mail{Body: []byte("Hello\rWorld")},
			wantMail: &Mail{Body: []byte("Hello\r\nWorld")},
			wantErr:  false,
		},
		{
			name:     "happy2",
			mail:     &Mail{Body: []byte("Hello\r\nWorld")},
			wantMail: &Mail{Body: []byte("Hello\r\nWorld")},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := UnixDosProcessor{}
			err := processor.Init(context.Background(), config.MailProcessorConfig{})
			require.NoError(t, err)
			got, err := processor.Process(context.Background(), tt.mail)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantMail, got)
		})
	}
}
