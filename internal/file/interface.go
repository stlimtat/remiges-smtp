// Package file provides file reading and tracking functionality for mail processing.
// It includes implementations for tracking file processing states, reading file contents,
// and managing file processing workflows.
//
// The package is designed to handle mail processing files in a directory, where:
// - Files starting with "df" are data files containing mail content
// - Files starting with "qf" are queue files tracking processing status
// - The file tracker maintains the state of file processing to prevent duplicate processing
package file

import (
	"context"
	"io"

	"github.com/stlimtat/remiges-smtp/pkg/input"
)

// FileInfo represents a mail processing file and its associated metadata.
// It contains information about both the data file (df) and queue file (qf),
// along with their respective readers and processing status.
type FileInfo struct {
	// DfFilePath is the absolute path to the data file containing mail content
	DfFilePath string

	// DfReader is the reader for the data file content
	DfReader io.Reader

	// ID is the unique identifier for the file, derived from the filename
	// (e.g., "123" for "df123" and "qf123")
	ID string

	// QfFilePath is the absolute path to the queue file tracking processing status
	QfFilePath string

	// QfReader is the reader for the queue file content
	QfReader io.Reader

	// Status represents the current processing state of the file
	Status input.FileStatus
}

// IFileReader defines the interface for file reading operations in the mail processing system.
// Implementations of this interface provide functionality to:
// - Scan and list available files for processing
// - Read the next unprocessed file
// - Validate file existence and accessibility
// - Derive queue file names from data file names
type IFileReader interface {
	// RefreshList scans the input directory for new files to process.
	// It updates the internal list of files and resets the file index.
	// Only files starting with "df" are considered for processing.
	//
	// Returns:
	//   - []*FileInfo: The list of files found in the input directory
	//   - error: Non-nil if directory scanning fails
	RefreshList(ctx context.Context) ([]*FileInfo, error)

	// ReadNextFile retrieves the next unprocessed file from the list.
	// It skips files that have already been processed or are currently being processed.
	//
	// Returns:
	//   - *FileInfo: Information about the next file to process
	//   - error: Non-nil if:
	//     - No more files are available
	//     - File tracking operations fail
	//     - The file is already being processed
	ReadNextFile(ctx context.Context) (*FileInfo, error)
}

// IFileReadTracker defines the interface for tracking file processing states.
// Implementations of this interface provide functionality to:
// - Track which files have been read
// - Update file processing status
// - Prevent duplicate processing of files
type IFileReadTracker interface {
	// FileRead retrieves the current processing status of a file by its ID.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - id: The unique identifier of the file to check
	//
	// Returns:
	//   - input.FileStatus: The current processing status of the file
	//   - error: Non-nil if the operation fails
	FileRead(ctx context.Context, id string) (input.FileStatus, error)

	// UpsertFile updates or inserts a file's processing status.
	// It first checks if the status has changed before performing the update
	// to prevent unnecessary operations.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - id: The unique identifier of the file
	//   - status: The new processing status to set
	//
	// Returns:
	//   - error: Non-nil if:
	//     - The status hasn't changed
	//     - The operation fails
	//     - The provided status is invalid
	UpsertFile(ctx context.Context, id string, status input.FileStatus) error
}

//go:generate mockgen -destination=mock.go -package=file . IFileReader,IFileReadTracker
