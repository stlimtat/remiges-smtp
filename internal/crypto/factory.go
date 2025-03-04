package crypto

import (
	"context"
	"crypto"

	"github.com/rs/zerolog"
)

type CryptoFactory struct {
	Generator IKeyGenerator
	Writer    IKeyWriter
}

func (c *CryptoFactory) Init(
	ctx context.Context,
	keyType string,
	writer IKeyWriter,
) (IKeyGenerator, error) {
	logger := zerolog.Ctx(ctx).With().Str("key_type", keyType).Logger()

	c.Writer = writer

	switch keyType {
	case KeyTypeEd25519:
		c.Generator = &Ed25519KeyGenerator{}
	default:
		c.Generator = &RsaKeyGenerator{}
	}

	logger.Info().Msg("new key generator created")

	return c.Generator, nil
}

func (c *CryptoFactory) GenerateKey(
	ctx context.Context,
	bitSize int,
	id string,
) (publicKeyPEM, privateKeyPEM []byte, err error) {
	publicKeyPEM, privateKeyPEM, err = c.Generator.GenerateKey(ctx, bitSize, id)
	if err != nil {
		return nil, nil, err
	}

	return publicKeyPEM, privateKeyPEM, nil
}

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

func (c *CryptoFactory) LoadPrivateKey(
	ctx context.Context,
	privateKeyPath string,
) (privateKey crypto.Signer, err error) {
	loader := c.Generator.(IKeyLoader)
	privateKey, err = loader.LoadPrivateKey(ctx, privateKeyPath)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
