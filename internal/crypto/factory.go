package crypto

import (
	"context"

	"github.com/rs/zerolog"
)

type CryptoFactory struct {
	keyGenerator IKeyGenerator
	keyWriter    IKeyWriter
}

func (c *CryptoFactory) Init(
	ctx context.Context,
	keyType string,
	keyWriter IKeyWriter,
) (IKeyGenerator, error) {
	logger := zerolog.Ctx(ctx).With().Str("key_type", keyType).Logger()

	c.keyWriter = keyWriter

	switch keyType {
	case KeyTypeEd25519:
		c.keyGenerator = &Ed25519KeyGenerator{}
	default:
		c.keyGenerator = &RsaKeyGenerator{}
	}

	logger.Info().Msg("new key generator created")

	return c.keyGenerator, nil
}

func (c *CryptoFactory) GenerateKey(
	ctx context.Context,
	bitSize int,
	id string,
) (publicKeyPEM, privateKeyPEM []byte, err error) {
	publicKeyPEM, privateKeyPEM, err = c.keyGenerator.GenerateKey(ctx, bitSize, id)
	if err != nil {
		return nil, nil, err
	}

	return publicKeyPEM, privateKeyPEM, nil
}

func (c *CryptoFactory) WriteKey(
	ctx context.Context,
	id string,
	publicKeyPEM, privateKeyPEM []byte,
) error {
	err := c.keyWriter.WriteKey(ctx, id, publicKeyPEM, privateKeyPEM)
	if err != nil {
		return err
	}

	return nil
}
