package dkim

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"slices"
	"strings"
	"testing"

	mcrypto "github.com/stlimtat/remiges-smtp/internal/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxtGen_Generate(t *testing.T) {
	ctx := context.Background()
	gen := &TxtGen{}

	// Generate test keys
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	rsaPubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&rsaKey.PublicKey),
	})

	ed25519PubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	ed25519PubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: ed25519PubKey,
	})

	tests := []struct {
		name      string
		domain    string
		keyType   string
		selector  string
		pubKeyPEM []byte
		wantErr   bool
	}{
		{
			name:      "valid RSA key",
			domain:    "example.com",
			keyType:   "rsa",
			selector:  "default",
			pubKeyPEM: rsaPubKeyPEM,
			wantErr:   false,
		},
		{
			name:      "valid Ed25519 key",
			domain:    "example.com",
			keyType:   "ed25519",
			selector:  "default",
			pubKeyPEM: ed25519PubKeyPEM,
			wantErr:   false,
		},
		{
			name:      "empty public key",
			domain:    "example.com",
			keyType:   "rsa",
			selector:  "default",
			pubKeyPEM: []byte{},
			wantErr:   true,
		},
		{
			name:      "invalid PEM format",
			domain:    "example.com",
			keyType:   "rsa",
			selector:  "default",
			pubKeyPEM: []byte("invalid pem"),
			wantErr:   true,
		},
		{
			name:      "unsupported key type - defaults to rsa",
			domain:    "example.com",
			keyType:   "dsa",
			selector:  "default",
			pubKeyPEM: rsaPubKeyPEM,
			wantErr:   false,
		},
		{
			name:      "empty domain",
			domain:    "",
			keyType:   "rsa",
			selector:  "default",
			pubKeyPEM: rsaPubKeyPEM,
			wantErr:   true,
		},
		{
			name:      "empty selector",
			domain:    "example.com",
			keyType:   "rsa",
			selector:  "",
			pubKeyPEM: rsaPubKeyPEM,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gen.Generate(ctx, tt.domain, tt.keyType, tt.selector, tt.pubKeyPEM)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)

			// Verify the output format
			record := string(got)
			assert.Contains(t, record, "v=DKIM1")
			keyType := tt.keyType
			if !slices.Contains(mcrypto.ValidKeyTypes, tt.keyType) {
				keyType = mcrypto.KeyTypeRSA
			}
			assert.Contains(t, record, fmt.Sprintf("k=%s", keyType))
			assert.Contains(t, record, "p=")
			assert.Contains(t, record, tt.selector+"._domainkey."+tt.domain)
		})
	}
}

// func TestTxtGen_Generate_InvalidKey(t *testing.T) {
// 	ctx := context.Background()
// 	gen := &TxtGen{}

// 	// Create an invalid RSA key (too small)
// 	invalidKey, err := rsa.GenerateKey(rand.Reader, 512)
// 	require.NoError(t, err)
// 	invalidPubKeyPEM := pem.EncodeToMemory(&pem.Block{
// 		Type:  "RSA PUBLIC KEY",
// 		Bytes: x509.MarshalPKCS1PublicKey(&invalidKey.PublicKey),
// 	})

// 	_, err = gen.Generate(ctx, "example.com", "rsa", "default", invalidPubKeyPEM)
// 	assert.NoError(t, err)
// }

func TestTxtGen_Generate_EdgeCases(t *testing.T) {
	ctx := context.Background()
	gen := &TxtGen{}

	// Generate a valid RSA key for testing
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	})

	tests := []struct {
		name      string
		domain    string
		keyType   string
		selector  string
		pubKeyPEM []byte
	}{
		{
			name:      "domain with special characters",
			domain:    "sub-domain.example.com",
			keyType:   "rsa",
			selector:  "default",
			pubKeyPEM: pubKeyPEM,
		},
		{
			name:      "selector with special characters",
			domain:    "example.com",
			keyType:   "rsa",
			selector:  "selector-2024",
			pubKeyPEM: pubKeyPEM,
		},
		{
			name:      "long domain",
			domain:    "a" + strings.Repeat(".sub", 10) + ".example.com",
			keyType:   "rsa",
			selector:  "default",
			pubKeyPEM: pubKeyPEM,
		},
		{
			name:      "long selector",
			domain:    "example.com",
			keyType:   "rsa",
			selector:  "selector" + strings.Repeat("x", 50),
			pubKeyPEM: pubKeyPEM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gen.Generate(ctx, tt.domain, tt.keyType, tt.selector, tt.pubKeyPEM)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)

			record := string(got)
			assert.Contains(t, record, tt.selector+"._domainkey."+tt.domain)
			assert.Contains(t, record, "v=DKIM1")
			assert.Contains(t, record, fmt.Sprintf("k=%s", tt.keyType))
			assert.Contains(t, record, "p=")
		})
	}
}
