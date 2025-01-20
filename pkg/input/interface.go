package input

import (
	"context"
	"io"
)

type FileInfo struct {
	DfFilePath string
	ID         string
	QfFilePath string
	Status     FileStatus
}

//go:generate mockgen -destination=mock.go -package=input -source=interface.go
type IFileReader interface {
	RefreshList(ctx context.Context) ([]FileInfo, error)
	ReadNextFile(ctx context.Context) (io.Reader, io.Reader, error)
}
