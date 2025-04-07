package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandInt(t *testing.T) {
	tests := []struct {
		name       string
		max        int64
		wantErr    bool
		checkRange bool
	}{
		{
			name:       "positive range",
			max:        10,
			wantErr:    false,
			checkRange: true,
		},
		{
			name:       "large range",
			max:        1000000,
			wantErr:    false,
			checkRange: true,
		},
		// {
		// 	name:       "zero range",
		// 	max:        0,
		// 	wantErr:    true,
		// 	checkRange: false,
		// },
		// {
		// 	name:       "negative range",
		// 	max:        -1,
		// 	wantErr:    true,
		// 	checkRange: false,
		// },
		{
			name:       "single value range",
			max:        1,
			wantErr:    false,
			checkRange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandInt(tt.max)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(-1), got)
				return
			}

			require.NoError(t, err)
			if tt.checkRange {
				assert.GreaterOrEqual(t, got, int64(0))
				assert.Less(t, got, tt.max)
			}
		})
	}
}

// TestRandIntDistribution tests the distribution of random numbers
// to ensure they are reasonably uniform across the range.
func TestRandIntDistribution(t *testing.T) {
	const (
		maxValue = 10
		samples  = 10000
	)

	// Create a histogram to track the distribution
	histogram := make([]int, maxValue)

	// Generate samples and update histogram
	for i := 0; i < samples; i++ {
		value, err := RandInt(maxValue)
		require.NoError(t, err)
		histogram[value]++
	}

	// Calculate expected frequency (should be roughly equal)
	expected := float64(samples) / float64(maxValue)
	tolerance := 0.2 // 20% tolerance

	// Check if each bin is within tolerance
	for i, count := range histogram {
		actual := float64(count)
		deviation := (actual - expected) / expected
		assert.True(t, deviation < tolerance,
			"Bin %d has deviation %.2f%% (count: %d, expected: %.0f)",
			i, deviation*100, count, expected)
	}
}
