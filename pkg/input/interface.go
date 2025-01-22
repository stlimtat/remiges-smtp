package input

import (
	"context"
	"io"
)

type FileStatus uint8

const (
	FILE_STATUS_INIT          FileStatus = 0
	FILE_STATUS_PROCESSING    FileStatus = 1
	FILE_STATUS_BODY_READ     FileStatus = 2
	FILE_STATUS_HEADERS_READ  FileStatus = 3
	FILE_STATUS_HEADERS_PARSE FileStatus = 4
	FILE_STATUS_DONE          FileStatus = 99
)

type FileInfo struct {
	DfFilePath string
	DfReader   io.Reader
	ID         string
	QfFilePath string
	QfReader   io.Reader
	Status     FileStatus
}

//go:generate mockgen -destination=mock.go -package=input -source=interface.go
type IFileReader interface {
	RefreshList(ctx context.Context) ([]*FileInfo, error)
	ReadNextFile(ctx context.Context) (*FileInfo, error)
}
