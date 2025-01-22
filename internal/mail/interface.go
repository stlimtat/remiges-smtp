package mail

import (
	"context"

	"github.com/stlimtat/remiges-smtp/internal/config"
)

type Mail struct {
	From    string
	To      string
	Headers map[string]string
	Body    []byte
}

//go:generate mockgen -destination=mock.go -package=mail github.com/stlimtat/remiges-smtp/internal/mail IMailProcessor,IMailProcessorFactory
type IMailProcessor interface {
	Index() int
	Init(ctx context.Context, cfg config.MailProcessorConfig) error
	Process(ctx context.Context, inMail *Mail) (outMail *Mail, err error)
}

type IMailProcessorFactory interface {
	// Having the cfg here allows us to create different types of mail processors
	NewMailProcessors(ctx context.Context, cfgs []config.MailProcessorConfig) ([]IMailProcessor, error)
	// Factory will need to be also an IMailProcessor
	// Process is a builder function to build through the mail processors in
	// the order they were created
}
