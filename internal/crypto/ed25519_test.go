package crypto

import (
	"context"
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
			publicKeyPEM, privateKeyPEM, err := generator.GenerateKey(ctx, tt.bitSize, tt.id)
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
