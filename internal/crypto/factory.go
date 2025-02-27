package crypto

import (
	"context"

	"github.com/rs/zerolog"
)

type CryptoFactory struct {
	keyGenerator IKeyGenerator
}

func (c *CryptoFactory) NewGenerator(
	ctx context.Context,
	keyType string,
) (IKeyGenerator, error) {
	logger := zerolog.Ctx(ctx).With().Str("key_type", keyType).Logger()

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
	return c.keyGenerator.GenerateKey(ctx, bitSize, id)
}
