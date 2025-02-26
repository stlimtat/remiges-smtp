package output

import (
	"context"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/mail"
)

//go:generate mockgen -destination=mock.go -package=output github.com/stlimtat/remiges-smtp/internal/output IOutput,IOutputFactory
type IOutput interface {
	Write(ctx context.Context, myMail *mail.Mail, responses []mail.Response) error
}

type IOutputFactory interface {
	NewOutputs(ctx context.Context, cfgs []config.OutputConfig) ([]IOutput, error)
}
