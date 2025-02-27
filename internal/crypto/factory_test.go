package crypto

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCryptoFactory_GenerateKey(t *testing.T) {
	tests := []struct {
		name             string
		keyType          string
		bitSize          int
		id               string
		wantGeneratorErr bool
		wantGenKeyErr    bool
	}{
		{
			name:             "happy - rsa",
			keyType:          KeyTypeRSA,
			bitSize:          2048,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
		},
		{
			name:             "happy - ed25519",
			keyType:          KeyTypeEd25519,
			bitSize:          256,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
		},
		{
			name:             "default - will default to rsa",
			keyType:          "invalid",
			bitSize:          2048,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			factory := &CryptoFactory{}
			generator, err := factory.NewGenerator(ctx, tt.keyType)
			if tt.wantGeneratorErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, generator)
			switch tt.keyType {
			case KeyTypeEd25519:
				assert.IsType(t, &Ed25519KeyGenerator{}, generator)
			case KeyTypeRSA:
				assert.IsType(t, &RsaKeyGenerator{}, generator)
			default:
				assert.IsType(t, &RsaKeyGenerator{}, generator)
			}

			publicKeyPEM, privateKeyPEM, err := generator.GenerateKey(ctx, tt.bitSize, tt.id)
			if tt.wantGenKeyErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, publicKeyPEM)
			require.NotNil(t, privateKeyPEM)
			assert.Contains(t, string(publicKeyPEM), "PUBLIC KEY")
			assert.Contains(t, string(privateKeyPEM), "PRIVATE KEY")
		})
	}
}
