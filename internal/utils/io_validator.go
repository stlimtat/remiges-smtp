// Package utils provides utility functions for common operations across the application.
// This package includes functions for input/output validation, path resolution,
// and working directory management.
package utils

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/errors"
)

// ValidateIO validates and resolves a given path, ensuring it exists and matches the expected type (file or directory).
// It handles special path prefixes:
//   - "./" - resolves relative to the source root directory
//   - "~/" - resolves relative to the user's home directory
//
// The function returns specific error types for different failure scenarios:
//   - ErrPathRequired: When the input path is empty
//   - ErrWorkingDir: When unable to get the working directory
//   - ErrHomeDir: When unable to get the user's home directory
//   - ErrFileStatFailed: When unable to get file information
//   - ErrNotFile: When a file is expected but a directory is found
//   - ErrNotDir: When a directory is expected but a file is found
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
	createFile bool,
) (string, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("path", path).
		Bool("fileNotDir", fileNotDir).
		Logger()

	if path == "" {
		logger.Error().Msg("path is required")
		return "", errors.NewError(errors.ErrPathRequired, "path is required", nil)
	}

	if strings.HasPrefix(path, "./") {
		wd, err := getWorkingDirRelativeToSourceRoot(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get working directory")
			return "", errors.NewError(errors.ErrWorkingDir, "failed to get working directory", err)
		}
		path = filepath.Join(wd, path[2:])
	} else if strings.HasPrefix(path, "~/") {
		// TODO: when running in bazel, $HOME is not defined
		// Workaround: add `test --test_env=HOME=$HOME` to $HOME/.bazelrc
		// https://github.com/bazelbuild/rules_apple/issues/877
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get user home directory")
			return "", errors.NewError(errors.ErrHomeDir, "failed to get user home directory", err)
		}
		path = filepath.Join(home, path[2:])
	}

	result := filepath.Clean(path)
	fileInfo, err := os.Stat(result)
	if err != nil {
		if createFile {
			file, err := os.Create(result)
			if err != nil {
				logger.Error().Err(err).Msg("failed to create file")
				return "", errors.NewError(errors.ErrFileStatFailed, "failed to create file", err)
			}
			defer func() {
				err = file.Close()
				if err != nil {
					logger.Error().Err(err).Msg("failed to close file")
				}
			}()
			return result, errors.NewError(errors.ErrNewlyCreatedFile, "ToIgnore: create new file", nil)
		}
		logger.Error().Err(err).Msg("failed to get file info")
		return result, errors.NewError(errors.ErrFileStatFailed, "failed to get file info", err)
	}
	if fileNotDir && fileInfo.IsDir() {
		logger.Error().Msg("path is not a file")
		return "", errors.NewError(errors.ErrNotFile, "path is not a file", nil)
	} else if !fileNotDir && !fileInfo.IsDir() {
		logger.Error().Msg("path is not a directory")
		return "", errors.NewError(errors.ErrNotDir, "path is not a directory", nil)
	}
	return result, nil
}

// getWorkingDirRelativeToSourceRoot determines the working directory relative to the source root.
// It handles special cases for Bazel builds and traverses up the directory tree
// until it finds the source root (indicated by the presence of a .git directory).
//
// The function returns specific error types for different failure scenarios:
//   - ErrWorkingDir: When unable to get the working directory
//   - ErrConfig: When unable to find the source root directory
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
		return "", errors.NewError(errors.ErrWorkingDir, "failed to get working directory", err)
	}
	wd = filepath.Clean(wd)
	if strings.Contains(wd, "bazel-out") {
		// When running in bazel, we need to traverse up to find the source root
		// The source root is typically the directory containing the .git directory
		for {
			gitDir := filepath.Join(wd, ".git")
			if _, err := os.Stat(gitDir); err == nil {
				return wd, nil
			}
			parent := filepath.Dir(wd)
			if parent == wd {
				logger.Error().Msg("failed to find source root directory")
				return "", errors.NewError(errors.ErrConfig, "failed to find source root directory", nil)
			}
			wd = parent
		}
	}
	return wd, nil
}
