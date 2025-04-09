// Package output provides functionality for writing mail processing results to various output destinations.
// It includes implementations for different output types (e.g., file, HTTP, etc.) and a factory
// for creating output instances based on configuration.
package output

import (
	"context"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

// IOutput defines the interface for writing mail processing results to an output destination.
// Implementations of this interface handle the actual writing of mail data and processing responses
// to their respective destinations (e.g., files, HTTP endpoints, etc.).
type IOutput interface {
	// Write processes and writes mail data and responses to the output destination.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - fileInfo: Information about the source file being processed
	//   - myMail: The mail content being processed
	//   - responses: The processing results and responses
	//
	// Returns:
	//   - error: Non-nil if the write operation fails
	Write(ctx context.Context, fileInfo *file.FileInfo, myMail *pmail.Mail, responses []pmail.Response) error
}

// IOutputFactory defines the interface for creating output instances based on configuration.
// Implementations of this interface handle the creation and initialization of output
// destinations according to the provided configuration.
type IOutputFactory interface {
	// NewOutputs creates and initializes output instances based on the provided configurations.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - cfgs: List of output configurations specifying the type and settings of each output
	//
	// Returns:
	//   - []IOutput: List of initialized output instances
	//   - error: Non-nil if any output fails to initialize
	NewOutputs(ctx context.Context, cfgs []config.OutputConfig) ([]IOutput, error)
}

//go:generate mockgen -destination=mock.go -package=output github.com/stlimtat/remiges-smtp/internal/output IOutput,IOutputFactory
