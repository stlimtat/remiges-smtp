package file

import (
	"context"
	"io"

	"github.com/stlimtat/remiges-smtp/pkg/input"
)

type FileInfo struct {
	DfFilePath string
	DfReader   io.Reader
	ID         string
	QfFilePath string
	QfReader   io.Reader
	Status     input.FileStatus
}

//go:generate mockgen -destination=mock.go -package=file . IFileReader,IFileReadTracker
type IFileReader interface {
	RefreshList(ctx context.Context) ([]*FileInfo, error)
	ReadNextFile(ctx context.Context) (*FileInfo, error)
}

type IFileReadTracker interface {
	FileRead(ctx context.Context, id string) (input.FileStatus, error)
	UpsertFile(ctx context.Context, id string, status input.FileStatus) error
}
