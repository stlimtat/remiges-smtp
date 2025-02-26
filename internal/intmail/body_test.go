package intmail

import (
	"context"
	"testing"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyProcessor(t *testing.T) {
	tests := []struct {
		name         string
		inMail       *mail.Mail
		wantBody     []byte
		wantMetadata map[string][]byte
		wantErr      bool
	}{
		{
			name: "happy - default no body headers",
			inMail: &mail.Mail{
				From:    smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Subject: []byte("test"),
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Body: []byte("Hello\r\nWorld"),
			},
			wantBody:     []byte("Hello\r\nWorld"),
			wantMetadata: map[string][]byte{},
			wantErr:      false,
		},
		{
			name: "happy - body has from/to",
			inMail: &mail.Mail{
				From:    smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Subject: []byte("test"),
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Body: []byte("From: a@example.com\r\nTo: B User <b@example.com>, C User <c@example.com>\r\n\r\nHello\r\nWorld"),
			},
			wantBody: []byte("Hello\r\nWorld"),
			wantMetadata: map[string][]byte{
				input.HeaderFromKey: []byte("a@example.com"),
				input.HeaderToKey:   []byte("B User <b@example.com>, C User <c@example.com>"),
			},
			wantErr: false,
		},
		{
			name: "happy - body already has from/to - new lines",
			inMail: &mail.Mail{
				From:    smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Subject: []byte("test"),
				To: []smtp.Address{
					{Localpart: "john", Domain: dns.Domain{ASCII: "example.com"}},
					{Localpart: "jane", Domain: dns.Domain{ASCII: "example.com"}},
				},
				Body: []byte("From: a@example.com\r\nTo: B User <b@example.com>, C User <c@example.com>\r\nSubject: alt subject\r\n\r\nHello\r\nWorld"),
			},
			wantBody: []byte("Hello\r\nWorld"),
			wantMetadata: map[string][]byte{
				input.HeaderFromKey:    []byte("a@example.com"),
				input.HeaderToKey:      []byte("B User <b@example.com>, C User <c@example.com>"),
				input.HeaderSubjectKey: []byte("alt subject"),
			},
			wantErr: false,
		},
		{
			name: "happy - multipart message - with new line first line",
			inMail: &mail.Mail{
				From:    smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Subject: []byte("test"),
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
			wantBody:     []byte("------=_Part_123\r\nContent-Type: text/plain\r\nContent-Transfer-Encoding: 7bit\r\n\r\nHello\r\nWorld\r\n------=_Part_123--"),
			wantMetadata: map[string][]byte{},
			wantErr:      false,
		},
		{
			name: "happy - multipart message - with boundary first line",
			inMail: &mail.Mail{
				From:    smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Subject: []byte("test"),
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
			wantBody:     []byte("------=_Part_123\r\nContent-Type: text/plain\r\nContent-Transfer-Encoding: 7bit\r\n\r\nHello\r\nWorld\r\n------=_Part_123--"),
			wantMetadata: map[string][]byte{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			processor := BodyProcessor{}
			err := processor.Init(ctx, config.MailProcessorConfig{})
			require.NoError(t, err)
			got, err := processor.Process(ctx, tt.inMail)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantBody, got.Body)
			require.Equal(t, tt.wantMetadata, got.BodyHeaders)
		})
	}
}
