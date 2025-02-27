package intmail

import (
	"context"
	"testing"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyHeadersProcessor(t *testing.T) {
	tests := []struct {
		name            string
		inMail          *pmail.Mail
		wantBodyHeaders map[string][]byte
		wantErr         bool
	}{
		{
			name: "happy - default",
			inMail: &pmail.Mail{
				ContentType: []byte("text/plain"),
				From:        smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				MsgID:       []byte("1234567890"),
				Subject:     []byte("test"),
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
			},
			wantBodyHeaders: map[string][]byte{
				input.HeaderContentTypeKey: []byte("text/plain"),
				input.HeaderFromKey:        []byte("sender@example.com"),
				input.HeaderMsgIDKey:       []byte("1234567890"),
				input.HeaderSubjectKey:     []byte("test"),
				input.HeaderToKey:          []byte("john@example.com,jane@example.com"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			processor := BodyHeadersProcessor{}
			err := processor.Init(ctx, config.MailProcessorConfig{})
			require.NoError(t, err)
			got, err := processor.Process(ctx, tt.inMail)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Subset(t, got.BodyHeaders, tt.wantBodyHeaders)
		})
	}
}
