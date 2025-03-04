package crypto

import (
	"context"
	"crypto"

	"github.com/rs/zerolog"
)

type CryptoFactory struct {
	Generators map[string]IKeyGenerator
	Writer     IKeyWriter
}

func (c *CryptoFactory) Init(
	ctx context.Context,
	writer IKeyWriter,
) (map[string]IKeyGenerator, error) {
	logger := zerolog.Ctx(ctx)

	c.Writer = writer

	c.Generators = map[string]IKeyGenerator{
		KeyTypeEd25519: &Ed25519KeyGenerator{},
		KeyTypeRSA:     &RsaKeyGenerator{},
	}

	logger.Info().Msg("new key generators created")

	return c.Generators, nil
}

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
		logger.Error().Msg("key type not found")
		keyType = KeyTypeRSA
		generator = c.Generators[keyType]
	}
	publicKeyPEM, privateKeyPEM, err = generator.GenerateKey(ctx, bitSize, id, keyType)
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
		logger.Error().Msg("key type not found")
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
