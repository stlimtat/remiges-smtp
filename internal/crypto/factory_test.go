package crypto

import (
	"context"
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCryptoFactory_GenerateKey(t *testing.T) {
	tests := []struct {
		name             string
		keyType          string
		bitSize          int
		id               string
		wantGeneratorErr bool
		wantGenKeyErr    bool
		wantWriteKeyErr  bool
		wantLoadKeyErr   bool
	}{
		{
			name:             "happy - rsa",
			keyType:          KeyTypeRSA,
			bitSize:          2048,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
			wantWriteKeyErr:  false,
			wantLoadKeyErr:   false,
		},
		{
			name:             "happy - ed25519",
			keyType:          KeyTypeEd25519,
			bitSize:          256,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
			wantWriteKeyErr:  false,
			wantLoadKeyErr:   false,
		},
		{
			name:             "default - will default to rsa",
			keyType:          "invalid",
			bitSize:          2048,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
			wantWriteKeyErr:  false,
			wantLoadKeyErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tmpDir := t.TempDir()
			keyWriter := NewKeyWriter(ctx, tmpDir)
			factory := &CryptoFactory{}
			generators, err := factory.Init(ctx, keyWriter)
			if tt.wantGeneratorErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, generators)
			assert.Len(t, generators, 2)
			assert.Contains(t, generators, KeyTypeRSA)
			assert.Contains(t, generators, KeyTypeEd25519)
			assert.IsType(t, &RsaKeyGenerator{}, generators[KeyTypeRSA])
			assert.IsType(t, &Ed25519KeyGenerator{}, generators[KeyTypeEd25519])

			publicKeyPEM, privateKeyPEM, err := factory.GenerateKey(ctx, tt.bitSize, tt.id, tt.keyType)
			if tt.wantGenKeyErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, publicKeyPEM)
			require.NotNil(t, privateKeyPEM)

			publicKeyPath, privateKeyPath, err := factory.WriteKey(ctx, tt.id, publicKeyPEM, privateKeyPEM)
			if tt.wantWriteKeyErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("%s/%s.pub", tmpDir, tt.id), publicKeyPath)
			require.Equal(t, fmt.Sprintf("%s/%s.pem", tmpDir, tt.id), privateKeyPath)
			assert.FileExists(t, publicKeyPath)
			assert.FileExists(t, privateKeyPath)

			gotPrivateKey, err := factory.LoadPrivateKey(ctx, tt.keyType, privateKeyPath)
			if tt.wantLoadKeyErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotPrivateKey)

			switch tt.keyType {
			case KeyTypeEd25519:
				assert.IsType(t, ed25519.PrivateKey{}, gotPrivateKey)
			default:
				assert.IsType(t, &rsa.PrivateKey{}, gotPrivateKey)
			}
		})
	}
}
