package crypto

import (
	"context"
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
	}{
		{
			name:             "happy - rsa",
			keyType:          KeyTypeRSA,
			bitSize:          2048,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
			wantWriteKeyErr:  false,
		},
		{
			name:             "happy - ed25519",
			keyType:          KeyTypeEd25519,
			bitSize:          256,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
			wantWriteKeyErr:  false,
		},
		{
			name:             "default - will default to rsa",
			keyType:          "invalid",
			bitSize:          2048,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
			wantWriteKeyErr:  false,
		},
		{
			name:             "error - key writer",
			keyType:          KeyTypeRSA,
			bitSize:          2048,
			id:               "test",
			wantGeneratorErr: false,
			wantGenKeyErr:    false,
			wantWriteKeyErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			keyWriter := NewMockIKeyWriter(ctrl)
			factory := &CryptoFactory{}
			generator, err := factory.Init(ctx, tt.keyType, keyWriter)
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

			// replace the generator with a mock generator
			mockGenerator := NewMockIKeyGenerator(ctrl)
			mockGenerator.EXPECT().
				GenerateKey(ctx, tt.bitSize, tt.id).
				DoAndReturn(func(_ context.Context, _ int, _ string) ([]byte, []byte, error) {
					if tt.wantGenKeyErr {
						return nil, nil, fmt.Errorf("error generating key")
					}
					return []byte("public-key"), []byte("private-key"), nil
				})
			factory.keyGenerator = mockGenerator

			publicKeyPEM, privateKeyPEM, err := factory.GenerateKey(ctx, tt.bitSize, tt.id)
			if tt.wantGenKeyErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, publicKeyPEM)
			require.NotNil(t, privateKeyPEM)
			assert.Contains(t, string(publicKeyPEM), "public-key")
			assert.Contains(t, string(privateKeyPEM), "private-key")

			keyWriter.EXPECT().
				WriteKey(ctx, tt.id, publicKeyPEM, privateKeyPEM).
				DoAndReturn(func(_ context.Context, _ string, _, _ []byte) error {
					if tt.wantWriteKeyErr {
						return fmt.Errorf("error writing key")
					}
					return nil
				})
			err = factory.WriteKey(ctx, tt.id, publicKeyPEM, privateKeyPEM)
			if tt.wantWriteKeyErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
