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
			name: "happy - default no body headers",
			mail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Body: []byte("Hello\r\nWorld"),
			},
			wantMail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				BodyHeaders: map[string][]byte{
					"From": []byte("sender@example.com"),
					"To":   []byte("john@example.com,jane@example.com"),
				},
				Body: []byte("Hello\r\nWorld"),
			},
			wantErr: false,
		},
		{
			name: "happy - body already has from/to",
			mail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Body: []byte("From: a@example.com\r\nTo: b@example.com\r\n\r\nHello\r\nWorld"),
			},
			wantMail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				BodyHeaders: map[string][]byte{
					"From": []byte("sender@example.com"),
					"To":   []byte("john@example.com,jane@example.com"),
				},
				Body: []byte("Hello\r\nWorld"),
			},
			wantErr: false,
		},
		{
			name: "happy - body already has from/to - new lines",
			mail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Body: []byte("From: a@example.com\nTo: b@example.com\n\nHello\r\nWorld"),
			},
			wantMail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				BodyHeaders: map[string][]byte{
					"From": []byte("sender@example.com"),
					"To":   []byte("john@example.com,jane@example.com"),
				},
				Body: []byte("Hello\r\nWorld"),
			},
			wantErr: false,
		},
		{
			name: "happy - multipart message - with new line first line",
			mail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Body: []byte(`
------=_Part_123
Content-Type: text/plain
Content-Transfer-Encoding: 7bit

Hello
World
------=_Part_123--
				`),
			},
			wantMail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				BodyHeaders: map[string][]byte{
					"From": []byte("sender@example.com"),
					"To":   []byte("john@example.com,jane@example.com"),
				},
				Body: []byte("------=_Part_123\r\nContent-Type: text/plain\r\nContent-Transfer-Encoding: 7bit\r\n\r\nHello\r\nWorld\r\n------=_Part_123--"),
			},
			wantErr: false,
		},
		{
			name: "happy - multipart message",
			mail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Body: []byte(`------=_Part_123
Content-Type: text/plain
Content-Transfer-Encoding: 7bit

Hello
World
------=_Part_123--
				`),
			},
			wantMail: &Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				BodyHeaders: map[string][]byte{
					"From": []byte("sender@example.com"),
					"To":   []byte("john@example.com,jane@example.com"),
				},
				Body: []byte("------=_Part_123\r\nContent-Type: text/plain\r\nContent-Transfer-Encoding: 7bit\r\n\r\nHello\r\nWorld\r\n------=_Part_123--"),
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
			require.Equal(t, tt.wantMail.Body, got.Body)
			require.Equal(t, tt.wantMail.BodyHeaders, got.BodyHeaders)
		})
	}
}
