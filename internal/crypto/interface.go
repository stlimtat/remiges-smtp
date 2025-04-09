// Package crypto provides cryptographic operations for key generation, loading, and writing.
// It supports multiple key types including RSA and Ed25519, and provides interfaces
// for key management operations.
package crypto

import (
	"context"
	"crypto"
)

const (
	// KeyTypeEd25519 represents the Ed25519 digital signature algorithm.
	// Ed25519 is a modern, high-performance signature scheme that provides
	// strong security guarantees and fast verification.
	KeyTypeEd25519 string = "ed25519"

	// KeyTypeRSA represents the RSA public-key cryptosystem.
	// RSA is a widely used algorithm for secure data transmission and
	// digital signatures.
	KeyTypeRSA string = "rsa"
)

var (
	// ValidKeyTypes is a list of valid key types.
	ValidKeyTypes []string = []string{KeyTypeEd25519, KeyTypeRSA}
)

// IKeyGenerator defines the interface for cryptographic key generation.
// Implementations should support generating keys of different types and sizes.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - bitSize: The size of the key in bits (relevant for RSA)
//   - id: Unique identifier for the key pair
//   - keyType: The type of key to generate (KeyTypeEd25519 or KeyTypeRSA)
//
// Returns:
//   - publicKeyPEM: The public key in PEM format
//   - privateKeyPEM: The private key in PEM format
//   - error: Non-nil if key generation fails
type IKeyGenerator interface {
	GenerateKey(ctx context.Context, bitSize int, id string, keyType string) (publicKeyPEM, privateKeyPEM []byte, err error)
}

// IKeyLoader defines the interface for loading private keys from storage.
// Implementations should support loading keys of different types and formats.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - keyType: The type of key to load (KeyTypeEd25519 or KeyTypeRSA)
//   - privateKeyPath: Path to the private key file
//
// Returns:
//   - crypto.Signer: A signer interface for the loaded private key
//   - error: Non-nil if key loading fails
type IKeyLoader interface {
	LoadPrivateKey(ctx context.Context, keyType string, privateKeyPath string) (crypto.Signer, error)
}

// IKeyWriter defines the interface for writing cryptographic keys to storage.
// Implementations should support writing keys in PEM format and return the paths
// where the keys were stored.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - id: Unique identifier for the key pair
//   - publicKeyPEM: The public key in PEM format
//   - privateKeyPEM: The private key in PEM format
//
// Returns:
//   - publicKeyPath: Path where the public key was stored
//   - privateKeyPath: Path where the private key was stored
//   - error: Non-nil if key writing fails
type IKeyWriter interface {
	WriteKey(ctx context.Context, id string, publicKeyPEM, privateKeyPEM []byte) (publicKeyPath, privateKeyPath string, err error)
}
