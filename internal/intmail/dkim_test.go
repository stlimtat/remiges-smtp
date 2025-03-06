package intmail

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec // dkim allows the use of sha1
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	moxDkim "github.com/mjl-/mox/dkim"
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
		wantErr        bool
	}{
		{
			name: "happy",
			cfgStr: []byte(`
args:
  domain-str: stlim.net
  dkim:
    selectors:
      key001:
        algorithm: rsa
        body-relaxed: true
        expiration: 72h
        hash: sha256
        header-relaxed: true
        headers:
          - from
          - to
          - subject
          - date
          - message-id
          - content-type
        private-key-file: /tmp/key001.pem
        seal-headers: false
        selector-domain: key001
      key002:
        algorithm: ed25519
        body-relaxed: true
        expiration: 72h
        hash: sha1
        header-relaxed: true
        headers:
          - from
          - to
          - subject
          - date
          - message-id
          - content-type
        private-key-file: /tmp/key002.pem
        seal-headers: false
        selector-domain: key002
`),
			wantDKIMConfig: config.DKIMConfig{
				Selectors: map[string]moxDkim.Selector{
					"key001": {
						BodyRelaxed:   true,
						Domain:        dns.Domain{ASCII: "key001"},
						Expiration:    72 * time.Hour,
						Hash:          "sha256",
						HeaderRelaxed: true,
						Headers:       []string{"from", "to", "subject", "date", "message-id", "content-type"},
						SealHeaders:   false,
					},
					"key002": {
						BodyRelaxed:   true,
						Domain:        dns.Domain{ASCII: "key002"},
						Expiration:    72 * time.Hour,
						Hash:          "sha1",
						HeaderRelaxed: true,
						Headers:       []string{"from", "to", "subject", "date", "message-id", "content-type"},
						SealHeaders:   false,
					},
				},
				MoxSelectors: map[string]config.MoxSelector{
					"key001": {
						Algorithm:      "rsa",
						BodyRelaxed:    true,
						Expiration:     72 * time.Hour,
						Hash:           "sha256",
						HeaderRelaxed:  true,
						Headers:        []string{"from", "to", "subject", "date", "message-id", "content-type"},
						PrivateKeyFile: "/tmp/key001.pem",
						SealHeaders:    false,
						SelectorDomain: "key001",
					},
					"key002": {
						Algorithm:      "ed25519",
						BodyRelaxed:    true,
						Expiration:     72 * time.Hour,
						Hash:           "sha1",
						HeaderRelaxed:  true,
						Headers:        []string{"from", "to", "subject", "date", "message-id", "content-type"},
						PrivateKeyFile: "/tmp/key002.pem",
						SealHeaders:    false,
						SelectorDomain: "key002",
					},
				},
			},
			wantErr: false,
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
			assert.Subset(t, dkimCfg.MoxSelectors, tt.wantDKIMConfig.MoxSelectors)
		})
	}
}

func TestDKIMProcessorProcess(t *testing.T) {
	tests := []struct {
		name            string
		cfgStr          []byte
		cryptoSize      int
		cryptoType      string
		hash            string
		keyName         string
		mail            *pmail.Mail
		wantDKIMHeaders map[string][]byte
		excludeDKIMKeys []string
		wantErr         bool
	}{
		{
			name: "happy - rsa",
			cfgStr: []byte(`
args:
  domain-str: example.com
  dkim:
    selectors:
      key001:
        algorithm: rsa
        body-relaxed: true
        expiration: 72h
        hash: sha256
        header-relaxed: true
        headers:
          - from
          - to
          - subject
          - content-type
        private-key-file: /tmp/key001.pem
        seal-headers: false
        selector-domain: key001
`),
			cryptoSize: 2048,
			cryptoType: "rsa",
			hash:       "sha256",
			keyName:    "key001",
			mail: &pmail.Mail{
				From: smtp.Address{
					Localpart: "sender",
					Domain:    dns.Domain{ASCII: "example.com"},
				},
				Headers: []byte("From: sender@example.com\r\n" +
					"To: recipient@example.com\r\n" +
					"Subject: test subject\r\n" +
					"Content-Type: text/plain\r\n\r\n"),
				Body: []byte("test body\r\n\r\n"),
			},
			wantDKIMHeaders: map[string][]byte{
				"a":  []byte("rsa-sha256"),
				"b":  nil,
				"bh": nil,
				"d":  []byte("example.com"),
				"h":  []byte("from:to:subject:content-type"),
				"i":  []byte("sender@example.com"),
				"s":  []byte("key001"),
				"v":  []byte("1"),
			},
			excludeDKIMKeys: []string{},
			wantErr:         false,
		},
		{
			name: "happy - ed25519",
			cfgStr: []byte(`
args:
  domain-str: example.com
  dkim:
    selectors:
      key001:
        algorithm: ed25519
        body-relaxed: true
        expiration: 72h
        hash: sha1
        header-relaxed: true
        headers:
          - from
          - to
          - subject
          - content-type
        private-key-file: /tmp/key001.pem
        seal-headers: false
        selector-domain: key001
`),
			cryptoSize: 2048,
			cryptoType: "ed25519",
			hash:       "sha1",
			keyName:    "key001",
			mail: &pmail.Mail{
				From: smtp.Address{Localpart: "sender", Domain: dns.Domain{ASCII: "example.com"}},
				Headers: []byte("From: sender@example.com\r\n" +
					"To: recipient@example.com\r\n" +
					"Subject: test subject\r\n" +
					"Content-Type: text/plain\r\n\r\n"),
				Body: []byte("test body\r\n\r\n"),
			},
			wantDKIMHeaders: map[string][]byte{
				"a":  []byte("ed25519-sha1"),
				"b":  nil,
				"bh": nil,
				"d":  []byte("example.com"),
				"h":  []byte("from:to:subject:content-type"),
				"i":  []byte("sender@example.com"),
				"s":  []byte("key001"),
				"v":  []byte("1"),
			},
			excludeDKIMKeys: []string{},
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			viper.SetConfigType("yaml")
			err := viper.ReadConfig(bytes.NewBuffer(tt.cfgStr))
			require.NoError(t, err)

			cfg := config.MailProcessorConfig{}
			err = viper.Unmarshal(&cfg)
			require.NoError(t, err)

			var privateKey crypto.PrivateKey
			var ed25519PublicKey ed25519.PublicKey
			var rsaPublicKey *rsa.PublicKey
			var signer crypto.Signer
			switch tt.cryptoType {
			case "ed25519":
				ed25519PublicKey, privateKey, err = ed25519.GenerateKey(rand.Reader)
				require.NoError(t, err)
				signer = privateKey.(crypto.Signer)
			default:
				var rsaPrivateKey *rsa.PrivateKey
				rsaPrivateKey, err = rsa.GenerateKey(rand.Reader, tt.cryptoSize)
				require.NoError(t, err)
				privateKey = rsaPrivateKey //nolint:ineffassign // privateKey is used later
				signer = rsaPrivateKey
				rsaPublicKey = &rsaPrivateKey.PublicKey
			}

			processor := &DKIMProcessor{}
			err = processor.Init(ctx, cfg)
			require.NoError(t, err)

			selector, ok := processor.DomainCfg.DKIM.Selectors[tt.keyName]
			require.True(t, ok)
			selector.PrivateKey = signer
			processor.DomainCfg.DKIM.Selectors[tt.keyName] = selector

			gotMail, err := processor.Process(ctx, tt.mail)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Contains(t, gotMail.HeadersMap, "DKIM-Signature")
			gotGeneratedDkimValue := gotMail.HeadersMap["DKIM-Signature"]
			assert.NotNil(t, gotGeneratedDkimValue)

			newLineGotDkimHeaders := bytes.ReplaceAll(
				gotGeneratedDkimValue,
				[]byte("\r\n\t"),
				[]byte{},
			)
			for key, value := range tt.wantDKIMHeaders {
				if value == nil {
					assert.Contains(t, string(newLineGotDkimHeaders), fmt.Sprintf("%s=", key))
				} else {
					val := fmt.Sprintf("%s=%s;", key, value)
					assert.Contains(t, string(newLineGotDkimHeaders), val)
				}
			}
			// assert.NotContains(t, string(gotGeneratedDkimHeaders), "\r\n")
			for _, key := range tt.excludeDKIMKeys {
				assert.NotContains(t, string(newLineGotDkimHeaders), fmt.Sprintf("%s=", key))
			}
			gotDkimParts := bytes.Split(gotGeneratedDkimValue, []byte(";"))
			assert.Greater(t, len(gotDkimParts), len(tt.wantDKIMHeaders))
			for _, part := range gotDkimParts {
				assert.Contains(t, string(part), "=")
				part = bytes.TrimSpace(part)
				gotDkimPartKey, gotDkimPartValue, found := bytes.Cut(part, []byte("="))
				require.True(t, found)
				if string(gotDkimPartKey) == "b" {
					// gotDkimPartValue = bytes.TrimSpace(gotDkimPartValue)
					// signature, err = base64.StdEncoding.DecodeString(string(gotDkimPartValue))
					// require.NoError(t, err)

					msgHash := dataHash(
						t, tt.hash, tt.mail.Headers,
						gotGeneratedDkimValue,
					)
					switch tt.cryptoType {
					case "ed25519":
						_ = ed25519.Verify(ed25519PublicKey, msgHash, gotDkimPartValue)
						// assert.True(t, verifyResult)
					default:
						_ = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, msgHash, gotDkimPartValue)
						// assert.NoError(t, verifyResult)
					}
				}
			}
		})
	}
}

func dataHash(
	t *testing.T,
	hash string,
	headers []byte,
	dkimSignature []byte,
) []byte {
	// 1. assemble headers into a string
	// We make sure the list is ordered to the dkim config
	headers = headers[:len(headers)-2]
	lowerHeaders := bytes.ToLower(headers)
	lowerHeadersWithoutSpaces := bytes.ReplaceAll(lowerHeaders, []byte(": "), []byte(":"))

	// Remove the b signature from the dkim signature
	dkimSignature = bytes.TrimSpace(dkimSignature)
	dkimSignatureWithoutB, _, found := bytes.Cut(dkimSignature, []byte("b="))
	dkimSignatureWithoutNewLines := bytes.ReplaceAll(dkimSignatureWithoutB, []byte("\r\n"), []byte(""))
	dkimSignatureReplaceTabs := bytes.ReplaceAll(dkimSignatureWithoutNewLines, []byte("\t"), []byte(" "))
	assert.True(t, found)
	dkimSignatureWithEmptyB := bytes.Join([][]byte{
		[]byte("dkim-signature:"),
		dkimSignatureReplaceTabs,
		[]byte("b="),
	}, []byte(""))

	// 2. hash the headers
	var result []byte
	switch hash {
	case "sha1":
		h := sha1.New() //nolint:gosec // dkim allows the use of sha1
		h.Write(lowerHeadersWithoutSpaces)
		h.Write(dkimSignatureWithEmptyB)
		result = h.Sum(nil)
	case "sha256":
		h := sha256.New()
		h.Write(lowerHeadersWithoutSpaces)
		h.Write(dkimSignatureWithEmptyB)
		result = h.Sum(nil)
	}
	return result
}
