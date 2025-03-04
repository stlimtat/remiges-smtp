package crypto

import (
	"context"
	"crypto/ed25519"
	"encoding/pem"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEd25519KeyGenerator_GenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		bitSize int
		id      string
		wantErr bool
	}{
		{
			name:    "happy",
			bitSize: 256,
			id:      "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			generator := &Ed25519KeyGenerator{}
			publicKeyPEM, privateKeyPEM, err := generator.GenerateKey(ctx, tt.bitSize, tt.id, KeyTypeEd25519)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, publicKeyPEM)
			assert.NotNil(t, privateKeyPEM)
			assert.Contains(t, string(publicKeyPEM), "ED25519 PUBLIC KEY")
			assert.Contains(t, string(privateKeyPEM), "ED25519 PRIVATE KEY")
		})
	}
}

func TestEd25519KeyGenerator_WriteThenLoad(t *testing.T) {
	tests := []struct {
		name         string
		bitSize      int
		wantWriteErr bool
		wantLoadErr  bool
	}{
		{
			name:         "happy",
			bitSize:      256,
			wantWriteErr: false,
			wantLoadErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			generator := &Ed25519KeyGenerator{}
			gotPublicKeyPEM, gotPrivateKeyPEM, err := generator.GenerateKey(ctx, tt.bitSize, "test", KeyTypeEd25519)
			require.NoError(t, err)
			gotBlock, _ := pem.Decode(gotPrivateKeyPEM)
			assert.Equal(t, "ED25519 PRIVATE KEY", gotBlock.Type)
			gotGeneratedPrivateKey := ed25519.PrivateKey(gotBlock.Bytes)

			tmpDir := t.TempDir()
			keyWriter := NewKeyWriter(ctx, tmpDir)

			publicKeyPath, privateKeyPath, err := keyWriter.WriteKey(ctx, "test", gotPublicKeyPEM, gotPrivateKeyPEM)
			if tt.wantWriteErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, publicKeyPath)
			assert.NotEmpty(t, privateKeyPath)

			gotLoadedPrivateKey, err := generator.LoadPrivateKey(ctx, KeyTypeEd25519, privateKeyPath)
			if tt.wantLoadErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, gotLoadedPrivateKey)
			assert.IsType(t, ed25519.PrivateKey{}, gotLoadedPrivateKey)
			gotLoadedPrivateKeyBytes := gotLoadedPrivateKey.(ed25519.PrivateKey)
			assert.Equal(t, gotGeneratedPrivateKey, gotLoadedPrivateKeyBytes)
		})
	}
}
