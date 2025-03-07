package output

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

type OutputFactory struct {
	Cfgs        []config.OutputConfig
	Outputs     []IOutput
	FileTracker file.IFileReadTracker
}

func NewOutputFactory(
	_ context.Context,
	fileTracker file.IFileReadTracker,
) *OutputFactory {
	result := &OutputFactory{
		FileTracker: fileTracker,
	}
	return result
}

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

func (f *OutputFactory) Write(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *pmail.Mail,
	resp []pmail.Response,
) error {
	for _, output := range f.Outputs {
		err := output.Write(ctx, fileInfo, myMail, resp)
		if err != nil {
			return err
		}
	}
	return nil
}
