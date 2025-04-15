// Package output provides functionality for writing mail processing results to various output destinations.
// It includes implementations for different output types (e.g., file, HTTP, etc.) and a factory
// for creating output instances based on configuration.
package output

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

// FileTrackerOutput implements the IOutput interface for tracking file processing status.
// It updates the file tracker to mark files as processed, allowing the system to track
// which files have been successfully processed and prevent duplicate processing.
type FileTrackerOutput struct {
	// Cfg contains the output configuration specifying the output type and settings
	Cfg config.OutputConfig

	// FileTracker is the interface for tracking file processing states
	FileTracker file.IFileReadTracker
}

// NewFileTrackerOutput creates a new FileTrackerOutput instance with the provided configuration
// and file tracker. This output type is used to mark files as processed in the file tracking system.
//
// Parameters:
//   - ctx: Context for logging and cancellation (currently unused)
//   - cfg: Output configuration
//   - fileTracker: The file tracker instance to use for updating file status
//
// Returns:
//   - *FileTrackerOutput: A new FileTrackerOutput instance
//   - error: Always nil in the current implementation
func NewFileTrackerOutput(
	_ context.Context,
	cfg config.OutputConfig,
	fileTracker file.IFileReadTracker,
) (*FileTrackerOutput, error) {
	result := &FileTrackerOutput{
		Cfg:         cfg,
		FileTracker: fileTracker,
	}
	return result, nil
}

// Write implements the IOutput interface by updating the file tracker to mark a file as processed.
// It sets the file status to FILE_STATUS_DONE in the file tracker, indicating that the file
// has been successfully processed.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - fileInfo: Information about the source file being processed
//   - myMail: The mail content being processed (used for logging)
//   - _: The processing responses (unused in this implementation)
//
// Returns:
//   - error: Non-nil if updating the file tracker fails
func (f *FileTrackerOutput) Write(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *pmail.Mail,
	_ map[string][]pmail.Response,
) error {
	logger := zerolog.Ctx(ctx).
		With().
		Str("fileInfo.id", fileInfo.ID).
		Bytes("mail", myMail.MsgID).
		Logger()
	logger.Debug().Msg("FileTrackerOutput: Write")

	err := f.FileTracker.UpsertFile(ctx, fileInfo.ID, input.FILE_STATUS_DONE)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to upsert file")
		return err
	}

	logger.Info().Msg("FileOutput: Write success")
	return nil
}
