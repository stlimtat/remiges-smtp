package intmail

import (
	"context"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

//go:generate mockgen -destination=mock.go -package=intmail . IMailProcessor,IMailProcessorFactory
type IMailProcessor interface {
	Index() int
	Init(ctx context.Context, cfg config.MailProcessorConfig) error
	Process(ctx context.Context, myMail *pmail.Mail) (*pmail.Mail, error)
}

type IMailProcessorFactory interface {
	// Having the cfg here allows us to create different types of mail processors
	NewMailProcessors(ctx context.Context, cfgs []config.MailProcessorConfig) ([]IMailProcessor, error)
	// Factory will need to be also an IMailProcessor
	// Process is a builder function to build through the mail processors in
	// the order they were created
}
