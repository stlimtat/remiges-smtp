package file_mail

import (
	"context"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
)

//go:generate mockgen -destination=mock.go -package=file_mail . IMailTransformer
type IMailTransformer interface {
	Transform(ctx context.Context, fileInfo *file.FileInfo, myMail *mail.Mail) (*mail.Mail, error)
}

type IMailTransformerFactory interface {
	NewMailTransformer(ctx context.Context, cfg config.FileMailConfig) IMailTransformer
}
