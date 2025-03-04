package crypto

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRsaKeyGenerator_GenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		bitSize int
		id      string
		wantErr bool
	}{
		{
			name:    "happy",
			bitSize: 2048,
			id:      "stlim.net",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			r := &RsaKeyGenerator{}
			gotPublic, gotPrivate, err := r.GenerateKey(ctx, tt.bitSize, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, gotPublic)
			assert.NotNil(t, gotPrivate)
			assert.Contains(t, string(gotPublic), "BEGIN RSA PUBLIC KEY")
			assert.Contains(t, string(gotPrivate), "BEGIN RSA PRIVATE KEY")
		})
	}
}

func TestRsaKeyGenerator_LoadPrivateKey(t *testing.T) {
	tests := []struct {
		name          string
		bitSize       int
		id            string
		wantGenKeyErr bool
		wantLoadErr   bool
		wantWriteErr  bool
	}{
		{
			name:          "happy",
			bitSize:       2048,
			id:            "stlim.net",
			wantGenKeyErr: false,
			wantLoadErr:   false,
			wantWriteErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			r := &RsaKeyGenerator{}
			gotGeneratedPublicKeyPEM, gotGeneratedPrivateKeyPEM, err := r.GenerateKey(ctx, tt.bitSize, tt.id)
			if tt.wantGenKeyErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, gotGeneratedPublicKeyPEM)
			assert.NotNil(t, gotGeneratedPrivateKeyPEM)
			assert.Contains(t, string(gotGeneratedPublicKeyPEM), "BEGIN RSA PUBLIC KEY")
			assert.Contains(t, string(gotGeneratedPrivateKeyPEM), "BEGIN RSA PRIVATE KEY")

			gotBlock, _ := pem.Decode(gotGeneratedPrivateKeyPEM)
			assert.Equal(t, "RSA PRIVATE KEY", gotBlock.Type)
			gotGeneratedPrivateKey, err := x509.ParsePKCS1PrivateKey(gotBlock.Bytes)
			require.NoError(t, err)

			tmpDir := t.TempDir()
			keyWriter := NewKeyWriter(ctx, tmpDir)
			gotPublicKeyPath, gotPrivateKeyPath, err := keyWriter.WriteKey(ctx, tt.id, gotGeneratedPublicKeyPEM, gotGeneratedPrivateKeyPEM)
			if tt.wantWriteErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, gotPublicKeyPath)
			assert.NotEmpty(t, gotPrivateKeyPath)
			assert.FileExists(t, gotPublicKeyPath)
			assert.FileExists(t, gotPrivateKeyPath)

			gotLoadedPrivateKey, err := r.LoadPrivateKey(ctx, gotPrivateKeyPath)
			if tt.wantLoadErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, gotLoadedPrivateKey)
			assert.IsType(t, &rsa.PrivateKey{}, gotLoadedPrivateKey)
			gotLoadedPrivateKeyBytes := gotLoadedPrivateKey.(*rsa.PrivateKey)
			assert.Equal(t, gotGeneratedPrivateKey.D, gotLoadedPrivateKeyBytes.D)
		})
	}
}
