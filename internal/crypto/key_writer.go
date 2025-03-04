package crypto

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/utils"
)

type KeyWriter struct {
	OutPath string
}

func NewKeyWriter(
	ctx context.Context,
	outPath string,
) (*KeyWriter, error) {
	err := utils.ValidateIO(ctx, outPath, false)
	if err != nil {
		return nil, err
	}
	result := &KeyWriter{
		OutPath: outPath,
	}

	return result, nil
}

func (k *KeyWriter) WriteKey(
	ctx context.Context,
	id string,
	publicKeyPEM, privateKeyPEM []byte,
) (publicKeyPath, privateKeyPath string, err error) {
	logger := zerolog.Ctx(ctx).With().Str("out_path", k.OutPath).Str("id", id).Logger()

	publicKeyPath = filepath.Join(k.OutPath, fmt.Sprintf("%s.pub", id))
	privateKeyPath = filepath.Join(k.OutPath, fmt.Sprintf("%s.pem", id))

	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create public key file")
		return "", "", fmt.Errorf("failed to create public key file: %w", err)
	}
	defer func() {
		if err = publicKeyFile.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close public key file")
		}
	}()

	_, err = publicKeyFile.Write(publicKeyPEM)
	if err != nil {
		logger.Error().Err(err).Msg("failed to write public key")
		return "", "", fmt.Errorf("failed to write public key: %w", err)
	}

	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create private key file")
		return "", "", fmt.Errorf("failed to create private key file: %w", err)
	}
	defer func() {
		if err = privateKeyFile.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close private key file")
		}
	}()

	_, err = privateKeyFile.Write(privateKeyPEM)
	if err != nil {
		logger.Error().Err(err).Msg("failed to write private key")
		return "", "", fmt.Errorf("failed to write private key: %w", err)
	}

	logger.Info().Msg("keys written to file")

	return publicKeyPath, privateKeyPath, nil
}
