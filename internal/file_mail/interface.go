package file_mail

import (
	"context"

	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	HeaderConfigArgType    = "type"
	HeaderConfigArgDefault = "default"
)

//go:generate mockgen -destination=mock.go -package=file_mail . IMailTransformer
type IMailTransformer interface {
	Init(ctx context.Context, cfg config.FileMailConfig) error
	Index() int
	Transform(ctx context.Context, fileInfo *file.FileInfo, myMail *pmail.Mail) (*pmail.Mail, error)
}

type IMailTransformerFactory interface {
	NewMailTransformer(ctx context.Context, cfg config.FileMailConfig) IMailTransformer
}
