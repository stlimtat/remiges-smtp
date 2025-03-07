package output

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

type FileTrackerOutput struct {
	Cfg         config.OutputConfig
	FileTracker file.IFileReadTracker
}

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

func (f *FileTrackerOutput) Write(
	ctx context.Context,
	fileInfo *file.FileInfo,
	myMail *pmail.Mail,
	_ []pmail.Response,
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
