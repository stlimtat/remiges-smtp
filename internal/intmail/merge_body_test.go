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

func TestMergeBodyProcessor(t *testing.T) {
	tests := []struct {
		name     string
		inMail   *pmail.Mail
		wantBody []byte
		wantErr  bool
	}{
		{
			name: "happy",
			inMail: &pmail.Mail{
				Headers: []byte("From: test@example.com\r\n\r\n"),
				Body:    []byte("Hello, world!"),
			},
			wantBody: []byte("From: test@example.com\r\n\r\nHello, world!\r\n\r\n"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			processor := &MergeBodyProcessor{}
			err := processor.Init(ctx, config.MailProcessorConfig{})
			require.NoError(t, err)
			gotMail, err := processor.Process(ctx, tt.inMail)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantBody, gotMail.FinalBody)
		})
	}
}
