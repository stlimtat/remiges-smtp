package crypto

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Ed25519KeyGenerator struct{}

func (_ *Ed25519KeyGenerator) GenerateKey(
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

	if keyType != KeyTypeEd25519 {
		logger.Error().Msg("key type not supported")
		return nil, nil, fmt.Errorf("key type %s not supported", keyType)
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate ed25519 key pair")
		return nil, nil, err
	}

	publicKeyPEM = pem.EncodeToMemory(
		&pem.Block{
			Type:  "ED25519 PUBLIC KEY",
			Bytes: publicKey,
		},
	)
	privateKeyPEM = pem.EncodeToMemory(
		&pem.Block{
			Type:  "ED25519 PRIVATE KEY",
			Bytes: privateKey,
		},
	)

	logger.Info().
		Bytes("public_key", publicKeyPEM).
		Bytes("private_key", privateKeyPEM).
		Msg("generated ed25519 key pair")

	return publicKeyPEM, privateKeyPEM, nil
}

func (_ *Ed25519KeyGenerator) LoadPrivateKey(
	ctx context.Context,
	keyType string,
	privateKeyPath string,
) (crypto.Signer, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("key_type", keyType).
		Str("private_key_path", privateKeyPath).
		Logger()

	if keyType != KeyTypeEd25519 {
		logger.Error().Msg("key type not supported")
		return nil, fmt.Errorf("key type %s not supported", keyType)
	}

	rawFileData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read private key")
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(rawFileData)
	if privateKeyBlock == nil {
		logger.Error().Msg("failed to decode private key")
		return nil, err
	}

	result := ed25519.PrivateKey(privateKeyBlock.Bytes)
	return result, nil
}
