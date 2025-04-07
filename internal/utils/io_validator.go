// Package utils provides utility functions for common operations across the application.
// This package includes functions for input/output validation, path resolution,
// and working directory management.
package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

// ValidateIO validates and resolves a given path, ensuring it exists and matches the expected type (file or directory).
// It handles special path prefixes:
//   - "./" - resolves relative to the source root directory
//   - "~/" - resolves relative to the user's home directory
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - path: The path to validate and resolve
//   - fileNotDir: If true, expects a file; if false, expects a directory
//
// Returns:
//   - string: The resolved and cleaned absolute path
//   - error: Non-nil if validation fails or path resolution encounters an error
func ValidateIO(
	ctx context.Context,
	path string,
	fileNotDir bool,
) (string, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("path", path).
		Bool("fileNotDir", fileNotDir).
		Logger()

	if path == "" {
		logger.Error().Msg("path is required")
		return "", fmt.Errorf("path is required")
	}

	if strings.HasPrefix(path, "./") {
		wd, err := getWorkingDirRelativeToSourceRoot(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get working directory")
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
		path = filepath.Join(wd, path[2:])
	} else if strings.HasPrefix(path, "~/") {
		// TODO: when running in bazel, $HOME is not defined
		// Workaround: add `test --test_env=HOME=$HOME` to $HOME/.bazelrc
		// https://github.com/bazelbuild/rules_apple/issues/877
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get user home directory")
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = filepath.Join(home, path[2:])
	}

	result := filepath.Clean(path)
	fileInfo, err := os.Stat(result)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get file info")
		return "", fmt.Errorf("failed to get file info: %w", err)
	}
	if fileNotDir && fileInfo.IsDir() {
		logger.Error().Msg("path is not a file")
		return "", fmt.Errorf("path is not a file")
	} else if !fileNotDir && !fileInfo.IsDir() {
		logger.Error().Msg("path is not a directory")
		return "", fmt.Errorf("path is not a directory")
	}
	return result, nil
}

// getWorkingDirRelativeToSourceRoot determines the working directory relative to the source root.
// It handles special cases for Bazel builds and traverses up the directory tree
// until it finds the source root (indicated by the presence of a .git directory).
//
// Parameters:
//   - ctx: Context for logging and cancellation
//
// Returns:
//   - string: The absolute path to the source root directory
//   - error: Non-nil if the source root cannot be determined
func getWorkingDirRelativeToSourceRoot(
	ctx context.Context,
) (string, error) {
	logger := zerolog.Ctx(ctx)
	wd, err := os.Getwd()
	if err != nil {
		logger.Error().Err(err).Msg("failed to get working directory")
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	wd = filepath.Clean(wd)
	if strings.Contains(wd, "bazel-out") {
		logger.Warn().Str("wd", wd).Msg("working directory is inside bazel-out")
		return wd, nil
	}
	for strings.Contains(wd, "cmd") || strings.Contains(wd, "internal") || strings.Contains(wd, "pkg") {
		wd = filepath.Dir(wd)
	}
	logger.Warn().Str("wd", wd).Msg("working directory")
	fmt.Println("working directory", wd)

	// We know this is the source root
	// only if .git exists
	gitDir := filepath.Join(wd, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		logger.Error().Err(err).Msg("source root not found")
		return "", fmt.Errorf("source root not found")
	}

	return wd, nil
}
