package crypto

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/rs/zerolog"
)

type RsaKeyGenerator struct{}

func (_ *RsaKeyGenerator) GenerateKey(
	ctx context.Context,
	bitSize int,
	id string,
) (publicKeyPEM, privateKeyPEM []byte, err error) {
	logger := zerolog.Ctx(ctx).With().Int("bit_size", bitSize).Str("id", id).Logger()

	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate private key")
		return nil, nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		logger.Error().Err(err).Msg("failed to validate private key")
		return nil, nil, err
	}

	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyDER,
		},
	)
	publicKeyDER := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicKeyPEM = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyDER,
		},
	)
	logger.Info().
		Bytes("public_key", publicKeyPEM).
		Bytes("private_key", privateKeyPEM).
		Msg("generated key pair")
	return publicKeyPEM, privateKeyPEM, nil
}

func (_ *RsaKeyGenerator) LoadPrivateKey(
	ctx context.Context,
	privateKeyPath string,
) (crypto.Signer, error) {
	logger := zerolog.Ctx(ctx).With().Str("private_key_path", privateKeyPath).Logger()

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

	result, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse private key")
		return nil, err
	}

	return result, nil
}
