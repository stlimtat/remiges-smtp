package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

func ValidateIO(
	ctx context.Context,
	path string,
	fileNotDir bool,
) error {
	logger := zerolog.Ctx(ctx).
		With().
		Str("path", path).
		Bool("fileNotDir", fileNotDir).
		Logger()

	if path == "" {
		logger.Error().Msg("path is required")
		return fmt.Errorf("path is required")
	}

	if strings.HasPrefix(path, "./") {
		wd, err := getWorkingDirRelativeToSourceRoot(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get working directory")
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		path = filepath.Join(wd, path[2:])
	} else if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get user home directory")
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = filepath.Join(home, path[2:])
	}

	filePath := filepath.Clean(path)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get file info")
		return fmt.Errorf("failed to get file info: %w", err)
	}
	if fileNotDir && fileInfo.IsDir() {
		logger.Error().Msg("path is not a file")
		return fmt.Errorf("path is not a file")
	} else if !fileNotDir && !fileInfo.IsDir() {
		logger.Error().Msg("path is not a directory")
		return fmt.Errorf("path is not a directory")
	}
	return nil
}

func getWorkingDirRelativeToSourceRoot(
	_ context.Context,
) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	wd = filepath.Clean(wd)
	for strings.Contains(wd, "cmd") || strings.Contains(wd, "internal") || strings.Contains(wd, "pkg") {
		wd = filepath.Dir(wd)
	}

	// We know this is the source root
	// only if .git exists
	gitDir := filepath.Join(wd, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return "", fmt.Errorf("source root not found")
	}

	return wd, nil
}
