package input

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

type DefaultFileReader struct {
	FileIndex int
	Files     []FileInfo
	InPath    string
	mutex     sync.Mutex
}

func NewDefaultFileReader(
	ctx context.Context,
	inPath string,
) (*DefaultFileReader, error) {
	logger := zerolog.Ctx(ctx)
	var err error
	result := &DefaultFileReader{
		FileIndex: 0,
		Files:     make([]FileInfo, 0),
		InPath:    inPath,
		mutex:     sync.Mutex{},
	}

	// 1. check that directory exists
	err = result.ValidateInPath(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("ValidateInPath")
		return nil, err
	}

	return result, nil
}

func (r *DefaultFileReader) ValidateInPath(
	ctx context.Context,
) error {
	logger := zerolog.Ctx(ctx)
	// 1. check that directory exists
	fileInfo, err := os.Stat(r.InPath)
	if err != nil {
		logger.Error().Err(err).Msg("os.Stat")
		return err
	}
	// 2. check that it is a directory
	if !fileInfo.IsDir() {
		logger.Error().Msg("InPath is not a directory")
		return fmt.Errorf("InPath is not a directory")
	}
	return nil
}

func (r *DefaultFileReader) ValidateFile(
	ctx context.Context,
	fileName string,
) error {
	// 1. file path
	filePath := filepath.Join(r.InPath, fileName)
	err := r.ValidateFilePath(ctx, filePath)
	return err
}

func (_ *DefaultFileReader) ValidateFilePath(
	ctx context.Context,
	filePath string,
) error {
	logger := zerolog.Ctx(ctx)
	// 1. check that file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logger.Error().Err(err).Msg("os.Stat")
		return err
	}
	// 2. check that it is a file
	if fileInfo.IsDir() {
		logger.Error().Msg("fileName is not a file")
		return fmt.Errorf("fileName is not a file")
	}
	return nil
}

func (r *DefaultFileReader) GetQfFileName(
	ctx context.Context,
	dfFileName string,
) (string, error) {
	logger := zerolog.Ctx(ctx)
	result := ""
	if strings.HasPrefix(dfFileName, "df") {
		result = strings.Replace(dfFileName, "df", "qf", 1)
		err := r.ValidateFile(ctx, result)
		if err != nil {
			logger.Error().Err(err).Msg("ValidateFile")
			return result, err
		}
	}
	return result, nil
}

func (r *DefaultFileReader) RefreshList(
	ctx context.Context,
) ([]FileInfo, error) {
	logger := zerolog.Ctx(ctx)
	// 0. Reset the current list of files
	result := make([]FileInfo, 0)
	// 1. validate directory
	err := r.ValidateInPath(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("ValidateInPath")
		return nil, err
	}
	// 2. read list of files in directory
	entries, err := os.ReadDir(r.InPath)
	if err != nil {
		logger.Error().Err(err).Msg("os.ReadDir")
		return nil, err
	}
	// 3. no files found
	if len(entries) < 1 {
		logger.Error().Msg("no files found in directory")
		return nil, fmt.Errorf("no files found in directory")
	}
	// 4. check with new files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		dfFileName := entry.Name()
		if !strings.HasPrefix(dfFileName, "df") {
			continue
		}
		qfFileName, err := r.GetQfFileName(ctx, dfFileName)
		if err != nil {
			logger.Error().Err(err).Msg("GetQfFileName")
			continue
		}
		logger.Info().
			Str("dfFileName", dfFileName).
			Str("qfFileName", qfFileName).
			Msg("RefreshList")
		id := dfFileName[2:]
		fileInfo := FileInfo{
			DfFilePath: filepath.Join(r.InPath, dfFileName),
			ID:         id,
			QfFilePath: filepath.Join(r.InPath, qfFileName),
			Status:     FILE_STATUS_INIT,
		}
		result = append(result, fileInfo)
	}
	// 5. update the list of files - note this needs
	// to be thread safe
	r.mutex.Lock()
	r.Files = result
	r.FileIndex = 0
	defer r.mutex.Unlock()
	return r.Files, nil
}

func (r *DefaultFileReader) ReadNextFile(
	ctx context.Context,
) (dfReader, qfReader io.Reader, err error) {
	logger := zerolog.Ctx(ctx)
	// 1. check that there are files
	if len(r.Files) == 0 {
		logger.Error().Msg("no files found")
		return nil, nil, fmt.Errorf("no files found")
	}
	// 2. read the next file
	r.mutex.Lock()
	dfFilePath := r.Files[r.FileIndex].DfFilePath
	qfFilePath := r.Files[r.FileIndex].QfFilePath
	r.FileIndex++
	r.mutex.Unlock()
	// 3. validate the df file
	err = r.ValidateFilePath(ctx, dfFilePath)
	if err != nil {
		logger.Error().Err(err).
			Str("dfFilePath", dfFilePath).
			Msg("ValidateFile")
		return nil, nil, err
	}
	// 4. validate the qf file
	err = r.ValidateFilePath(ctx, qfFilePath)
	if err != nil {
		logger.Error().Err(err).
			Str("qfFilePath", qfFilePath).
			Msg("ValidateFile")
		return nil, nil, err
	}
	// 6. read the df file
	dfFileBytes, err := os.ReadFile(dfFilePath)
	if err != nil {
		logger.Error().Err(err).
			Str("dfFilePath", dfFilePath).
			Msg("os.ReadFile")
		return nil, nil, err
	}
	// 7. read the qf file
	qfFileBytes, err := os.ReadFile(qfFilePath)
	if err != nil {
		logger.Error().Err(err).
			Str("qfFilePath", qfFilePath).
			Msg("os.ReadFile")
		return nil, nil, err
	}

	return bytes.NewReader(dfFileBytes), bytes.NewReader(qfFileBytes), nil
}
