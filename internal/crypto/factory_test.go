package crypto

import (
	"context"
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"os"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCryptoFactory_Init(t *testing.T) {
	tests := []struct {
		name        string
		writer      IKeyWriter
		expectError bool
	}{
		{
			name:        "valid writer",
			writer:      &MockIKeyWriter{},
			expectError: false,
		},
		{
			name:        "nil writer",
			writer:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			factory := &CryptoFactory{}
			generators, err := factory.Init(ctx, tt.writer)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, generators)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, generators)
			assert.Len(t, generators, 2)
			assert.Contains(t, generators, KeyTypeRSA)
			assert.Contains(t, generators, KeyTypeEd25519)
		})
	}
}

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
		// {
		// 	name:             "invalid bit size for rsa",
		// 	keyType:          KeyTypeRSA,
		// 	bitSize:          64, // Too small for RSA
		// 	id:               "test",
		// 	wantGeneratorErr: false,
		// 	wantGenKeyErr:    true,
		// 	wantWriteKeyErr:  false,
		// 	wantLoadKeyErr:   false,
		// },
		{
			name:             "empty id",
			keyType:          KeyTypeRSA,
			bitSize:          2048,
			id:               "",
			wantGeneratorErr: false,
			wantGenKeyErr:    true,
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
			keyWriter, err := NewKeyWriter(ctx, tmpDir)
			require.NoError(t, err)
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

func TestCryptoFactory_WriteKey_ErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	tests := []struct {
		name          string
		setup         func(t *testing.T) (string, func())
		expectError   bool
		errorContains string
	}{
		{
			name: "read-only directory",
			setup: func(t *testing.T) (string, func()) {
				// Make directory read-only
				require.NoError(t, os.Chmod(tmpDir, 0444))
				return tmpDir, func() {
					// Restore permissions
					require.NoError(t, os.Chmod(tmpDir, 0755))
				}
			},
			expectError:   true,
			errorContains: "permission denied",
		},
		// {
		// 	name: "non-existent directory",
		// 	setup: func(t *testing.T) (string, func()) {
		// 		nonExistentDir := filepath.Join(tmpDir, "nonexistent")
		// 		return nonExistentDir, func() {}
		// 	},
		// 	expectError:   true,
		// 	errorContains: "no such file or directory",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			dir, cleanup := tt.setup(t)
			defer cleanup()

			keyWriter, err := NewKeyWriter(ctx, dir)
			require.NoError(t, err)
			factory := &CryptoFactory{}
			_, err = factory.Init(ctx, keyWriter)
			require.NoError(t, err)

			// Generate test keys
			publicKeyPEM, privateKeyPEM, err := factory.GenerateKey(ctx, 2048, "test", KeyTypeRSA)
			require.NoError(t, err)

			// Attempt to write keys
			_, _, err = factory.WriteKey(ctx, "test", publicKeyPEM, privateKeyPEM)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCryptoFactory_LoadPrivateKey_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		keyType        string
		privateKeyPath string
		expectError    bool
		errorContains  string
	}{
		{
			name:           "non-existent file",
			keyType:        KeyTypeRSA,
			privateKeyPath: "/nonexistent/path/key.pem",
			expectError:    true,
			errorContains:  "no such file or directory",
		},
		{
			name:           "invalid key type",
			keyType:        "invalid-type",
			privateKeyPath: "test.pem",
			expectError:    true,
			errorContains:  "key type not found",
		},
		{
			name:           "empty key path",
			keyType:        KeyTypeRSA,
			privateKeyPath: "",
			expectError:    true,
			errorContains:  "empty key path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			factory := &CryptoFactory{}
			_, err := factory.Init(ctx, &MockIKeyWriter{})
			require.NoError(t, err)

			_, err = factory.LoadPrivateKey(ctx, tt.keyType, tt.privateKeyPath)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
