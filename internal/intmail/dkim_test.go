package intmail

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"

	moxConfig "github.com/mjl-/mox/config"
	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtp"
	"github.com/spf13/viper"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDKIMProcessorInit(t *testing.T) {
	tests := []struct {
		name           string
		cfgStr         []byte
		wantDKIMConfig config.DKIMConfig
	}{
		{
			name: "happy",
			cfgStr: []byte(`
args:
  domain-str: stlim.net
  dkim:
    selectors:
      key001:
        domain: stlim.net
      key002:
        domain: blah.com
    sign:
      - key001
      - key002
`),
			wantDKIMConfig: config.DKIMConfig{
				DKIM: moxConfig.DKIM{
					Selectors: map[string]moxConfig.Selector{
						"key001": {
							Domain: dns.Domain{ASCII: "stlim.net"},
						},
						"key002": {
							Domain: dns.Domain{ASCII: "blah.com"},
						},
					},
					Sign: []string{"key001", "key002"},
				},
				MoxSelectors: map[string]config.MoxSelector{
					"key001": {
						Domain: "stlim.net",
					},
					"key002": {
						Domain: "blah.com",
					},
				},
				MoxSign: []string{"key001", "key002"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			cfg := config.MailProcessorConfig{}
			viper.SetConfigType("yaml")
			err := viper.ReadConfig(bytes.NewBuffer(tt.cfgStr))
			require.NoError(t, err)
			settings := viper.AllSettings()
			assert.Contains(t, settings, "args")
			err = viper.Unmarshal(&cfg)
			require.NoError(t, err)
			processor := &DKIMProcessor{}
			err = processor.Init(ctx, cfg)
			require.NoError(t, err)
			dkimCfg := processor.DomainCfg.DKIM
			assert.Subset(t, dkimCfg.Selectors, tt.wantDKIMConfig.Selectors)
			assert.Subset(t, dkimCfg.Sign, tt.wantDKIMConfig.Sign)
			assert.Subset(t, dkimCfg.MoxSelectors, tt.wantDKIMConfig.MoxSelectors)
			assert.Subset(t, dkimCfg.MoxSign, tt.wantDKIMConfig.MoxSign)
		})
	}
}

func TestDKIMProcessorProcess(t *testing.T) {
	tests := []struct {
		name            string
		cfgStr          []byte
		cryptoSize      int
		cryptoType      string
		keyName         string
		mail            *pmail.Mail
		wantDKIMHeaders map[string][]byte
		wantErr         bool
	}{
		{
			name:       "happy - rsa",
			cryptoSize: 2048,
			cryptoType: "rsa",
			keyName:    "key001",
			cfgStr: []byte(`
args:
  domain-str: example.com
  dkim:
    selectors:
      key001:
        domain: key001
        hash: sha256
        private-key-file: /tmp/key001.pem
    sign:
      - key001
`),
			mail: &pmail.Mail{
				From: smtp.Address{
					Localpart: "sender",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				Body: []byte("From: sender@example.com\r\n" +
					"To: recipient@example.com\r\n" +
					"Subject: test subject\r\n" +
					"Content-Type: text/plain\r\n\r\n" +
					"test body\r\n\r\n"),
			},
			wantDKIMHeaders: map[string][]byte{
				"d": []byte("example.com"),
				"s": []byte("key001"),
				"i": []byte("sender@example.com"),
				"a": []byte("rsa-sha256"),
				"v": []byte("1"),
			},
			wantErr: false,
		},
		{
			name:       "happy - ed25519",
			cryptoSize: 2048,
			cryptoType: "ed25519",
			keyName:    "key001",
			cfgStr: []byte(`
args:
  domain-str: example.com
  dkim:
    selectors:
      key001:
        domain: key001
        hash: sha1
        private-key-file: /tmp/key001.pem
    sign:
      - key001
`),
			mail: &pmail.Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Body: []byte("From: sender@example.com\r\n" +
					"To: recipient@example.com\r\n" +
					"Subject: test subject\r\n" +
					"Content-Type: text/plain\r\n\r\n" +
					"test body\r\n\r\n"),
			},
			wantDKIMHeaders: map[string][]byte{
				"d": []byte("example.com"),
				"s": []byte("key001"),
				"i": []byte("sender@example.com"),
				"a": []byte("ed25519-sha1"),
				"v": []byte("1"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, logger := telemetry.InitLogger(context.Background())

			viper.SetConfigType("yaml")
			err := viper.ReadConfig(bytes.NewBuffer(tt.cfgStr))
			require.NoError(t, err)

			cfg := config.MailProcessorConfig{}
			err = viper.Unmarshal(&cfg)
			require.NoError(t, err)

			var privateKey crypto.PrivateKey
			var signer crypto.Signer
			switch tt.cryptoType {
			case "ed25519":
				_, privateKey, err = ed25519.GenerateKey(rand.Reader)
				require.NoError(t, err)
				signer = privateKey.(crypto.Signer)
			default:
				privateKey, err = rsa.GenerateKey(rand.Reader, tt.cryptoSize)
				require.NoError(t, err)
				signer = privateKey.(crypto.Signer)
			}

			processor := &DKIMProcessor{}
			err = processor.Init(ctx, cfg)
			require.NoError(t, err)

			selector := processor.DomainCfg.DKIM.DKIM.Selectors[tt.keyName]
			selector.Key = signer
			processor.DomainCfg.DKIM.DKIM.Selectors[tt.keyName] = selector

			gotMail, err := processor.Process(ctx, tt.mail)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, gotMail.DKIMHeaders)
			logger.Info().Bytes("dkimHeaders", gotMail.DKIMHeaders).Msg("dkimHeaders")
			for key, value := range tt.wantDKIMHeaders {
				val := fmt.Sprintf("%s=%s;", key, value)
				assert.Contains(t, string(gotMail.DKIMHeaders), val)
			}
		})
	}
}
