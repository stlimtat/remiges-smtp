package input

import (
	"context"
	"os"
	"testing"

	"github.com/stlimtat/remiges-smtp/internal/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestParseHeaders(t *testing.T) {
	var tests = []struct {
		name        string
		headerBytes []byte
		want        map[string]string
		wantErr     bool
	}{
		{
			"happy",
			[]byte(`H??Header1: Value1
H??Header2: Value2`),
			map[string]string{
				"Header1": "Value1",
				"Header2": "Value2",
			},
			false,
		},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = telemetry.GetLogger(ctx, os.Stdout)
			fr := NewFileReader(ctx, "inPath")
			got, err := fr.ParseHeaders(ctx, tt.headerBytes)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}
