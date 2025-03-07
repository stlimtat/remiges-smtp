package output

import (
	"context"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

//go:generate mockgen -destination=mock.go -package=output github.com/stlimtat/remiges-smtp/internal/output IOutput,IOutputFactory
type IOutput interface {
	Write(ctx context.Context, fileInfo *file.FileInfo, myMail *pmail.Mail, responses []pmail.Response) error
}

type IOutputFactory interface {
	NewOutputs(ctx context.Context, cfgs []config.OutputConfig) ([]IOutput, error)
}
