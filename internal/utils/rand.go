// Package utils provides utility functions for common operations across the application.
// This package includes functions for random number generation, input/output validation,
// and working directory management.
package utils

import (
	"crypto/rand"
	"math/big"
)

// RandInt generates a cryptographically secure random integer in the range [0, maxValue).
// It uses crypto/rand for secure random number generation, making it suitable for
// cryptographic operations and security-sensitive applications.
//
// Parameters:
//   - maxValue: The upper bound (exclusive) for the random number generation
//
// Returns:
//   - int64: A random integer in the range [0, maxValue)
//   - error: Non-nil if random number generation fails
//
// Example:
//
//	value, err := RandInt(100) // Generates a random number between 0 and 99
func RandInt(maxValue int64) (int64, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(maxValue))
	if err != nil {
		return -1, err
	}
	return value.Int64(), nil
}
