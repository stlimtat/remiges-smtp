package mail

import (
	"context"
	"testing"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyHeadersProcessor(t *testing.T) {
	tests := []struct {
		name     string
		mail     *Mail
		wantMail *Mail
		wantErr  bool
	}{
		{
			name: "happy - default",
			mail: &Mail{
				ContentType: []byte("text/plain"),
				From:        smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Subject:     []byte("test"),
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
			},
			wantMail: &Mail{
				ContentType: []byte("text/plain"),
				From:        smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Subject:     []byte("test"),
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				BodyHeaders: map[string][]byte{
					"Content-Type": []byte("text/plain"),
					"From":         []byte("sender@example.com"),
					"Subject":      []byte("test"),
					"To":           []byte("john@example.com,jane@example.com"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := BodyHeadersProcessor{}
			err := processor.Init(context.Background(), config.MailProcessorConfig{})
			require.NoError(t, err)
			got, err := processor.Process(context.Background(), tt.mail)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantMail.From, got.From)
			require.Equal(t, tt.wantMail.To, got.To)
			require.Equal(t, tt.wantMail.Subject, got.Subject)
			require.Equal(t, tt.wantMail.BodyHeaders, got.BodyHeaders)
		})
	}
}
