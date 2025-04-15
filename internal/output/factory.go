// Package output provides functionality for writing mail processing results to various output destinations.
// It supports multiple output types including file-based output and file tracker output.
package output

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

// OutputFactory is responsible for creating and managing output instances.
// It implements the factory pattern to create different types of outputs based on configuration.
//
// Fields:
//   - Cfgs: List of output configurations that define how each output should be created
//   - Outputs: List of initialized output instances
//   - FileTracker: Interface for tracking file read operations
type OutputFactory struct {
	Cfgs        []config.OutputConfig
	Outputs     []IOutput
	FileTracker file.IFileReadTracker
}

// NewOutputFactory creates a new instance of OutputFactory.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - fileTracker: Interface for tracking file read operations
//
// Returns:
//   - *OutputFactory: A new instance of OutputFactory
func NewOutputFactory(
	_ context.Context,
	fileTracker file.IFileReadTracker,
) *OutputFactory {
	result := &OutputFactory{
		FileTracker: fileTracker,
	}
	return result
}

// NewOutputs initializes multiple output instances based on the provided configurations.
// It handles the creation of different output types and manages their lifecycle.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - cfgs: List of output configurations
//
// Returns:
//   - []IOutput: List of initialized output instances
//   - error: Non-nil if initialization fails
//
// Error conditions:
//   - No configurations provided
//   - Invalid output type
//   - Failed to create output instance
func (f *OutputFactory) NewOutputs(
	ctx context.Context,
	cfgs []config.OutputConfig,
) ([]IOutput, error) {
	logger := zerolog.Ctx(ctx)
	if len(cfgs) == 0 {
		logger.Error().
			Msg("No output configurations provided")
		return nil, fmt.Errorf("no output configurations provided")
	}
	f.Cfgs = cfgs
	f.Outputs = make([]IOutput, 0)
	for _, cfg := range cfgs {
		output, err := f.NewOutput(ctx, cfg)
		if err != nil {
			return nil, err
		}
		if output == nil {
			logger.Error().
				Interface("cfg", cfg).
				Msg("Failed to create output")
			continue
		}
		f.Outputs = append(f.Outputs, output)
	}
	return f.Outputs, nil
}

// NewOutput creates a single output instance based on the provided configuration.
// It supports different output types and handles their specific initialization requirements.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - cfg: Output configuration
//
// Returns:
//   - IOutput: Initialized output instance
//   - error: Non-nil if initialization fails
//
// Supported output types:
//   - file: Writes output to files
//   - file_tracker: Tracks file read operations
func (f *OutputFactory) NewOutput(
	ctx context.Context,
	cfg config.OutputConfig,
) (IOutput, error) {
	logger := zerolog.Ctx(ctx)
	var result IOutput
	var err error
	switch cfg.Type {
	case config.ConfigOutputTypeFile:
		logger.Debug().
			Interface("cfg", cfg).
			Msg("Creating file output")
		result, err = NewFileOutput(ctx, cfg)
		if err != nil {
			return nil, err
		}
	case config.ConfigOutputTypeFileTracker:
		logger.Debug().
			Interface("cfg", cfg).
			Msg("Creating file tracker output")
		result, err = NewFileTrackerOutput(ctx, cfg, f.FileTracker)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown output type: %s", cfg.Type)
	}
	return result, nil
}

// Write sends data to all configured output instances.
// It handles the distribution of data to multiple outputs and manages any errors that occur.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - fileInfo: Information about the file being processed
//   - myMail: The mail content to be written
//   - responses: Map of SMTP responses for different recipients
//
// Returns:
//   - error: Non-nil if any output fails to write the data
func (f *OutputFactory) Write(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *pmail.Mail,
	responses map[string][]pmail.Response,
) error {
	for _, output := range f.Outputs {
		err := output.Write(ctx, fileInfo, myMail, responses)
		if err != nil {
			return err
		}
	}
	return nil
}
