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
			name: "happy1",
			mail: &Mail{
				From: smtp.Address{
					Localpart: "john",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				To: smtp.Address{
					Localpart: "jane",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				Body: []byte("Hello\r\nWorld"),
			},
			wantMail: &Mail{
				From: smtp.Address{
					Localpart: "john",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				To: smtp.Address{
					Localpart: "jane",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				BodyHeaders: map[string][]byte{
					"From": []byte("john@example.com"),
					"To":   []byte("jane@example.com"),
				},
				Body: []byte("Hello\r\nWorld"),
			},
			wantErr: false,
		},
		{
			name: "happy2",
			mail: &Mail{
				From: smtp.Address{
					Localpart: "john",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				To: smtp.Address{
					Localpart: "jane",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				Body: []byte("From: a@example.com\r\nTo: b@example.com\r\n\r\nHello\r\nWorld"),
			},
			wantMail: &Mail{
				From: smtp.Address{
					Localpart: "john",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				To: smtp.Address{
					Localpart: "jane",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				BodyHeaders: map[string][]byte{
					"From": []byte("john@example.com"),
					"To":   []byte("jane@example.com"),
				},
				Body: []byte("Hello\r\nWorld"),
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
			require.Equal(t, tt.wantMail, got)
		})
	}
}
