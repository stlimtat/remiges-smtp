package crypto

import (
	"context"
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
