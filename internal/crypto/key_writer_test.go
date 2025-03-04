package crypto

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyWriter_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	tests := []struct {
		name            string
		outPath         string
		wantValidateErr bool
		wantWriteErr    bool
	}{
		{
			name:            "happy",
			outPath:         tmpDir,
			wantValidateErr: false,
			wantWriteErr:    false,
		},
		{
			name:            "invalid out path",
			outPath:         filepath.Join(tmpDir, "does-not-exist"),
			wantValidateErr: true,
			wantWriteErr:    false,
		},
		{
			name:            "out path is a file",
			outPath:         filepath.Join(tmpDir, "key_writer_test.go"),
			wantValidateErr: true,
			wantWriteErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())
			k, err := NewKeyWriter(ctx, tt.outPath)
			if tt.wantValidateErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			randName := uuid.New().String()
			gotPublicKeyPath, gotPrivateKeyPath, err := k.WriteKey(ctx, randName, []byte("test-public-key"), []byte("test-private-key"))
			if tt.wantWriteErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.FileExists(t, gotPublicKeyPath)
			assert.FileExists(t, gotPrivateKeyPath)
			files, err := os.ReadDir(tt.outPath)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(files), 2)
			names := make([]string, 0)
			for _, file := range files {
				names = append(names, file.Name())
			}
			assert.Contains(t, names, fmt.Sprintf("%s.pub", randName))
			assert.Contains(t, names, fmt.Sprintf("%s.pem", randName))
			pubKey, err := os.ReadFile(filepath.Join(tt.outPath, fmt.Sprintf("%s.pub", randName)))
			require.NoError(t, err)
			assert.Equal(t, "test-public-key", string(pubKey))
			privKey, err := os.ReadFile(filepath.Join(tt.outPath, fmt.Sprintf("%s.pem", randName)))
			require.NoError(t, err)
			assert.Equal(t, "test-private-key", string(privKey))
		})
	}
}
