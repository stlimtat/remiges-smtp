package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandInt(t *testing.T) {
	var tests = []struct {
		name    string
		max     int64
		wantErr bool
	}{
		{"happy", 10, false},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandInt(tt.max)
			assert.LessOrEqual(t, got, tt.max)
			if tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}
