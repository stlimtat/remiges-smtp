package dkim

import (
	"context"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxtGen(t *testing.T) {
	tests := []struct {
		name       string
		keyType    string
		pubKeyPEM  []byte
		wantResult []byte
		wantErr    bool
	}{
		{
			name:    "happy - rsa",
			keyType: "rsa",
			pubKeyPEM: []byte(`-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAtI17ucNsuiV4hUxDeLjIutR6hhR5RDcm7xeeBWTJbiNYzbr6Bt0O
f4AXl0VJn1dj/lFXDKqq82ytV6aY8a0gQQkObvbh7uDo2+/fEo6e/7LLXL1dSq7N
GttKrEGQjNxQQTSemeqptYY+t9MapywaU8PzSve+urgTJeuEsGUUZQ0SV4vEXuV0
qAghANUoxBKk0nOUUknT1MzgVAl2bCrStQG7k5eg55okDfL5LW6zboM5oXB+/cWf
UsaPSYFgsaUBaMluFjj9SG/fuj4+a7KQ1uR72x+Di9fXKO32PTqhCmSc2Xy8Lznw
VAArBbos2eD1kbmJIOxYoKiZZBrQnrEMhQIDAQAB
-----END RSA PUBLIC KEY-----`),
			wantResult: []byte("selector._domainkey.example.com IN TXT \"v=DKIM1; k=rsa; p=MIIBCgKCAQEAtI17ucNsuiV4hUxDeLjIutR6hhR5RDcm7xeeBWTJbiNYzbr6Bt0Of4AXl0VJn1dj/lFXDKqq82ytV6aY8a0gQQkObvbh7uDo2+/fEo6e/7LLXL1dSq7NGttKrEGQjNxQQTSemeqptYY+t9MapywaU8PzSve+urgTJeuEsGUUZQ0SV4vEXuV0qAghANUoxBKk0nOUUknT1MzgVAl2bCrStQG7k5eg55okD\" \"fL5LW6zboM5oXB+/cWfUsaPSYFgsaUBaMluFjj9SG/fuj4+a7KQ1uR72x+Di9fXKO32PTqhCmSc2Xy8LznwVAArBbos2eD1kbmJIOxYoKiZZBrQnrEMhQIDAQAB\""),
			wantErr:    false,
		},
		{
			name:    "happy - ed25519",
			keyType: "ed25519",
			pubKeyPEM: []byte(`-----BEGIN ED25519 PUBLIC KEY-----
xgIPwhTr75HLiXOz2EcEokZpbE/wEOT1TCtM4yexrUM=
-----END ED25519 PUBLIC KEY-----`),
			wantResult: []byte("selector._domainkey.example.com IN TXT \"v=DKIM1; k=ed25519; p=xgIPwhTr75HLiXOz2EcEokZpbE/wEOT1TCtM4yexrUM=\""),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := telemetry.InitLogger(context.Background())

			txtGen := &TxtGen{}
			got, err := txtGen.Generate(ctx, "example.com", tt.keyType, "selector", tt.pubKeyPEM)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}
