// Package file provides file reading and tracking functionality for mail processing.
// It includes implementations for tracking file processing states, reading file contents,
// and managing file processing workflows.
package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/pkg/input"
)

// DefaultFileReader implements the IFileReader interface for managing a list of files
// to be processed. It tracks which files have been read and provides thread-safe
// access to file contents.
//
// The reader maintains an internal list of files and their processing status,
// ensuring that each file is processed only once and in the correct order.
type DefaultFileReader struct {
	// inputDir is the directory containing files to be processed
	inputDir string

	// files is the list of files to be processed
	files []*FileInfo

	// fileIndex is the current position in the files list
	fileIndex int

	// mu protects concurrent access to files and fileIndex
	mu sync.Mutex

	// fileReadTracker tracks which files have been read
	fileReadTracker IFileReadTracker
}

// NewDefaultFileReader creates a new instance of DefaultFileReader.
// It validates the input directory and initializes the file tracking system.
//
// Parameters:
//   - ctx: Context for initialization and logging
//   - inputDir: The directory containing files to be processed
//   - fileReadTracker: The tracker for file processing states
//
// Returns:
//   - *DefaultFileReader: A new reader instance
//   - error: Non-nil if the input directory is invalid or inaccessible
func NewDefaultFileReader(
	ctx context.Context,
	inputDir string,
	fileReadTracker IFileReadTracker,
) (*DefaultFileReader, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Str("inputDir", inputDir).
		Msg("NewDefaultFileReader")

	// Validate input directory
	info, err := os.Stat(inputDir)
	if err != nil {
		logger.Error().Err(err).Msg("NewDefaultFileReader: os.Stat")
		return nil, err
	}
	if !info.IsDir() {
		logger.Error().Msg("NewDefaultFileReader: not a directory")
		return nil, errors.New("not a directory")
	}

	return &DefaultFileReader{
		inputDir:        inputDir,
		files:           make([]*FileInfo, 0),
		fileIndex:       0,
		fileReadTracker: fileReadTracker,
	}, nil
}

// ValidateFile checks if a file exists and is accessible in the input directory.
// It verifies that the file:
// 1. Exists in the input directory
// 2. Is a regular file (not a directory)
// 3. Has read permissions
//
// Parameters:
//   - ctx: Context for logging
//   - fileName: The name of the file to validate
//
// Returns:
//   - bool: true if the file is valid and accessible, false otherwise
func (f *DefaultFileReader) ValidateFile(
	ctx context.Context,
	fileName string,
) bool {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Str("fileName", fileName).
		Msg("ValidateFile")

	filePath := filepath.Join(f.inputDir, fileName)
	info, err := os.Stat(filePath)
	if err != nil {
		logger.Error().Err(err).Msg("ValidateFile: os.Stat")
		return false
	}
	if info.IsDir() {
		logger.Error().Msg("ValidateFile: is a directory")
		return false
	}
	return true
}

// GetQfFileName derives the queue file name from a data file name.
// It replaces the "df" prefix with "qf" in the file name.
//
// Parameters:
//   - ctx: Context for logging
//   - fileName: The data file name to convert
//
// Returns:
//   - string: The corresponding queue file name
//   - error: Non-nil if the file name doesn't start with "df"
func (f *DefaultFileReader) GetQfFileName(
	ctx context.Context,
	fileName string,
) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Str("fileName", fileName).
		Msg("GetQfFileName")

	if len(fileName) < 2 {
		logger.Error().Msg("GetQfFileName: fileName too short")
		return "", errors.New("fileName too short")
	}
	if fileName[:2] != "df" {
		logger.Error().Msg("GetQfFileName: fileName does not start with df")
		return "", errors.New("fileName does not start with df")
	}
	return "qf" + fileName[2:], nil
}

// RefreshList scans the input directory for new files to process.
// It updates the internal file list and resets the file index.
// Only files starting with "df" are considered for processing.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//
// Returns:
//   - []*FileInfo: The list of files found in the input directory
//   - error: Non-nil if directory scanning fails
func (f *DefaultFileReader) RefreshList(
	ctx context.Context,
) ([]*FileInfo, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("RefreshList")

	f.mu.Lock()
	defer f.mu.Unlock()

	entries, err := os.ReadDir(f.inputDir)
	if err != nil {
		logger.Error().Err(err).Msg("RefreshList: os.ReadDir")
		return nil, err
	}

	f.files = make([]*FileInfo, 0)
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) >= 2 && entry.Name()[:2] == "df" {
			dfFilePath := filepath.Join(f.inputDir, entry.Name())
			qfFileName := "qf" + entry.Name()[2:]
			qfFilePath := filepath.Join(f.inputDir, qfFileName)
			f.files = append(f.files, &FileInfo{
				DfFilePath: dfFilePath,
				ID:         entry.Name()[2:],
				QfFilePath: qfFilePath,
				Status:     input.FILE_STATUS_INIT,
			})
		}
	}

	f.fileIndex = 0
	return f.files, nil
}

// ReadNextFile retrieves the next unprocessed file from the list.
// It skips files that have already been processed or are currently being processed.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//
// Returns:
//   - *FileInfo: Information about the next file to process
//   - error: Non-nil if:
//   - No more files are available
//   - File tracking operations fail
//   - The file is already being processed
func (f *DefaultFileReader) ReadNextFile(
	ctx context.Context,
) (*FileInfo, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("ReadNextFile")

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.fileIndex >= len(f.files) {
		logger.Debug().Msg("ReadNextFile: no more files")
		return nil, nil
	}

	file := f.files[f.fileIndex]
	status, err := f.fileReadTracker.FileRead(ctx, file.ID)
	if err != nil {
		logger.Error().Err(err).Msg("ReadNextFile: FileRead")
		return nil, err
	}

	if status == input.FILE_STATUS_PROCESSING {
		logger.Debug().
			Str("fileName", file.DfFilePath).
			Msg("ReadNextFile: file is being processed")
		return nil, fmt.Errorf("file is being processed: %s", file.DfFilePath)
	}

	if status == input.FILE_STATUS_DONE {
		logger.Debug().
			Str("fileName", file.DfFilePath).
			Msg("ReadNextFile: file is already done")
		f.fileIndex++
		return f.ReadNextFile(ctx)
	}

	err = f.fileReadTracker.UpsertFile(ctx, file.ID, input.FILE_STATUS_PROCESSING)
	if err != nil {
		logger.Error().Err(err).Msg("ReadNextFile: UpsertFile")
		return nil, err
	}

	f.fileIndex++
	return file, nil
}
