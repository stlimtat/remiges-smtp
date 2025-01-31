package input

import (
	"context"
	"io"

	"github.com/stlimtat/remiges-smtp/internal/mail"
)

type FileStatus uint8

const (
	FILE_STATUS_INIT          FileStatus = 1
	FILE_STATUS_PROCESSING    FileStatus = 2
	FILE_STATUS_BODY_READ     FileStatus = 3
	FILE_STATUS_HEADERS_READ  FileStatus = 4
	FILE_STATUS_HEADERS_PARSE FileStatus = 5
	FILE_STATUS_DONE          FileStatus = 99
	FILE_STATUS_ERROR         FileStatus = 0
)

type FileInfo struct {
	DfFilePath string
	DfReader   io.Reader
	ID         string
	QfFilePath string
	QfReader   io.Reader
	Status     FileStatus
}

//go:generate mockgen -destination=mock.go -package=input . IFileReader,IMailTransformer,IFileReadTracker
type IFileReader interface {
	RefreshList(ctx context.Context) ([]*FileInfo, error)
	ReadNextFile(ctx context.Context) (*FileInfo, error)
}

type IMailTransformer interface {
	Transform(ctx context.Context, fileInfo *FileInfo) (*mail.Mail, error)
}

type IFileReadTracker interface {
	FileRead(ctx context.Context, id string) (FileStatus, error)
	UpsertFile(ctx context.Context, id string, status FileStatus) error
}
