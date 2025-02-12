package mail

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeBodyProcessor(t *testing.T) {
	tests := []struct {
		name     string
		inMail   *Mail
		wantMail *Mail
		wantErr  bool
	}{
		{
			name: "happy",
			inMail: &Mail{
				BodyHeaders: map[string][]byte{"From": []byte("test@example.com")},
				Body:        []byte("Hello, world!"),
			},
			wantMail: &Mail{
				BodyHeaders: map[string][]byte{"From": []byte("test@example.com")},
				Body:        []byte("From: test@example.com\r\n\r\nHello, world!"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &MergeBodyProcessor{}
			err := processor.Init(context.Background(), config.MailProcessorConfig{})
			require.NoError(t, err)
			gotMail, err := processor.Process(context.Background(), tt.inMail)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantMail, gotMail)
		})
	}
}
