package crypto

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"

	"github.com/rs/zerolog"
)

type Ed25519KeyGenerator struct{}

func (_ *Ed25519KeyGenerator) GenerateKey(
	ctx context.Context,
	bitSize int,
	id string,
) (publicKeyPEM, privateKeyPEM []byte, err error) {
	logger := zerolog.Ctx(ctx).With().Int("bit_size", bitSize).Str("id", id).Logger()

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
