package crypto

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

type KeyWriter struct {
	OutPath string
}

func NewKeyWriter(
	_ context.Context,
	outPath string,
) *KeyWriter {
	result := &KeyWriter{
		OutPath: outPath,
	}

	return result
}

func (k *KeyWriter) Validate(
	ctx context.Context,
) error {
	logger := zerolog.Ctx(ctx).With().Str("out_path", k.OutPath).Logger()

	if k.OutPath == "" {
		logger.Error().Msg("out_path is required")
		return fmt.Errorf("out_path is required")
	}
	if strings.HasPrefix(k.OutPath, "./") {
		wd, err := os.Getwd()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get working directory")
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		k.OutPath = filepath.Join(wd, k.OutPath[2:])
	}
	if strings.HasPrefix(k.OutPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get user home directory")
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		k.OutPath = filepath.Join(home, k.OutPath[2:])
	}

	filePath := filepath.Clean(k.OutPath)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get file info")
		return fmt.Errorf("failed to get file info: %w", err)
	}
	if !fileInfo.IsDir() {
		logger.Error().Msg("out_path is not a directory")
		return fmt.Errorf("out_path is not a directory")
	}

	return nil
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
