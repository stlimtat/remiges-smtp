// Package crypto provides cryptographic operations for key generation, loading, and writing.
// It supports multiple key types including RSA and Ed25519, and provides interfaces
// for key management operations.
package crypto

import (
	"context"
	"crypto"
	"fmt"

	"github.com/rs/zerolog"
)

// CryptoFactory is a factory that manages cryptographic operations.
// It provides a unified interface for key generation, writing, and loading,
// supporting multiple key types through a plugin architecture.
type CryptoFactory struct {
	// Generators maps key types to their respective key generators.
	// Each generator implements the IKeyGenerator interface and optionally
	// the IKeyLoader interface for key loading operations.
	Generators map[string]IKeyGenerator

	// Writer handles the storage of generated keys.
	// It implements the IKeyWriter interface for writing keys to persistent storage.
	Writer IKeyWriter
}

// Init initializes the CryptoFactory with the provided key writer and
// registers the available key generators. It sets up the factory for
// subsequent cryptographic operations.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - writer: The key writer implementation for storing generated keys
//
// Returns:
//   - map[string]IKeyGenerator: The map of registered key generators
//   - error: Non-nil if initialization fails
func (c *CryptoFactory) Init(
	ctx context.Context,
	writer IKeyWriter,
) (map[string]IKeyGenerator, error) {
	logger := zerolog.Ctx(ctx)

	if writer == nil {
		logger.Error().Msg("writer is nil, using default writer")
		return nil, fmt.Errorf("writer is nil")
	}

	c.Writer = writer

	c.Generators = map[string]IKeyGenerator{
		KeyTypeEd25519: &Ed25519KeyGenerator{},
		KeyTypeRSA:     &RsaKeyGenerator{},
	}

	logger.Debug().Msg("CryptoFactory initialized")

	return c.Generators, nil
}

// GenerateKey generates a new key pair of the specified type and size.
// If the requested key type is not supported, it defaults to RSA.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - bitSize: The size of the key in bits (relevant for RSA)
//   - id: Unique identifier for the key pair
//   - keyType: The type of key to generate (KeyTypeEd25519 or KeyTypeRSA)
//
// Returns:
//   - publicKeyPEM: The public key in PEM format
//   - privateKeyPEM: The private key in PEM format
//   - error: Non-nil if key generation fails
func (c *CryptoFactory) GenerateKey(
	ctx context.Context,
	bitSize int,
	id string,
	keyType string,
) (publicKeyPEM, privateKeyPEM []byte, err error) {
	logger := zerolog.Ctx(ctx).
		With().
		Int("bit_size", bitSize).
		Str("id", id).
		Str("key_type", keyType).
		Logger()
	generator, ok := c.Generators[keyType]
	if !ok {
		logger.Warn().Msg("key type not found, defaulting to rsa")
		keyType = KeyTypeRSA
		generator = c.Generators[keyType]
	}
	if id == "" {
		logger.Error().Msg("id is empty")
		return nil, nil, fmt.Errorf("id is empty")
	}
	publicKeyPEM, privateKeyPEM, err = generator.GenerateKey(ctx, bitSize, id, keyType)
	if err != nil {
		return nil, nil, err
	}

	return publicKeyPEM, privateKeyPEM, nil
}

// WriteKey writes the generated key pair to storage using the configured writer.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - id: Unique identifier for the key pair
//   - publicKeyPEM: The public key in PEM format
//   - privateKeyPEM: The private key in PEM format
//
// Returns:
//   - publicKeyPath: Path where the public key was stored
//   - privateKeyPath: Path where the private key was stored
//   - error: Non-nil if key writing fails
func (c *CryptoFactory) WriteKey(
	ctx context.Context,
	id string,
	publicKeyPEM, privateKeyPEM []byte,
) (publicKeyPath, privateKeyPath string, err error) {
	publicKeyPath, privateKeyPath, err = c.Writer.WriteKey(ctx, id, publicKeyPEM, privateKeyPEM)
	if err != nil {
		return "", "", err
	}

	return publicKeyPath, privateKeyPath, nil
}

// LoadPrivateKey loads a private key from storage and returns a signer interface.
// If the requested key type is not supported, it defaults to RSA.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - keyType: The type of key to load (KeyTypeEd25519 or KeyTypeRSA)
//   - privateKeyPath: Path to the private key file
//
// Returns:
//   - crypto.Signer: A signer interface for the loaded private key
//   - error: Non-nil if key loading fails
func (c *CryptoFactory) LoadPrivateKey(
	ctx context.Context,
	keyType string,
	privateKeyPath string,
) (privateKey crypto.Signer, err error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("key_type", keyType).
		Str("private_key_path", privateKeyPath).
		Logger()
	generator, ok := c.Generators[keyType]
	if !ok {
		logger.Warn().Msg("key type not found, defaulting to rsa")
		keyType = KeyTypeRSA
		generator = c.Generators[keyType]
	}
	loader := generator.(IKeyLoader)
	privateKey, err = loader.LoadPrivateKey(ctx, keyType, privateKeyPath)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
