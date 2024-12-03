package input

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type FileStatus uint8

const (
	FILE_STATUS_PROCESSING    FileStatus = 0
	FILE_STATUS_BODY_READ     FileStatus = 1
	FILE_STATUS_HEADERS_READ  FileStatus = 2
	FILE_STATUS_HEADERS_PARSE FileStatus = 3
	FILE_STATUS_DONE          FileStatus = 9
)

type FileInfo struct {
	BodyFileName   string
	BodyBytes      []byte
	FileID         string
	HeaderBytes    []byte
	HeaderFileName string
	Headers        map[string]string
	Status         FileStatus
}

type FileReader struct {
	InPath    string
	ReadFiles map[string]FileInfo
}

func NewFileReader(
	_ context.Context,
	inPath string,
) *FileReader {
	result := &FileReader{
		InPath:    inPath,
		ReadFiles: make(map[string]FileInfo, 0),
	}

	return result
}

func (fr *FileReader) Process(
	ctx context.Context,
) ([]FileInfo, error) {
	logger := zerolog.Ctx(ctx)
	// 1. read list of files in directory
	entries, err := os.ReadDir(fr.InPath)
	if err != nil {
		logger.Error().Err(err).Msg("os.ReadDir")
		return nil, err
	}
	result := make([]FileInfo, 0)
	// 2. check with are new files
	for _, e := range entries {
		fileName := e.Name()
		// 3. read newest file - message, rcpt
		_, ok := fr.ReadFiles[fileName]
		if ok {
			logger.Error().Err(fmt.Errorf("fileInfo error")).Msg("fr.ReadFiles")
			continue
		}
		if strings.HasPrefix(fileName, "df") {
			fileID := fileName[2:]
			currFileInfo := FileInfo{
				BodyFileName: fileName,
				FileID:       fileID,
				Status:       FILE_STATUS_PROCESSING,
			}
			// 4. check if qf file exists
			qfFileName := strings.Replace(fileName, "df", "qf", 1)
			qfFullFileName := fr.InPath + "/" + qfFileName
			_, err := os.Stat(qfFullFileName)
			if err != nil {
				logger.Error().Err(err).Msg("qfExists")
				continue
			}
			currFileInfo.HeaderFileName = qfFileName
			// 5. read mail body
			currFileInfo.BodyBytes, err = fr.ReadFile(ctx, fr.InPath+"/"+fileName)
			if err != nil {
				continue
			}
			currFileInfo.Status = FILE_STATUS_BODY_READ
			// 6. read mail headers
			currFileInfo.HeaderBytes, err = fr.ReadFile(ctx, qfFullFileName)
			if err != nil {
				continue
			}
			currFileInfo.Status = FILE_STATUS_HEADERS_READ
			// 7. parse headers
			currFileInfo.Headers, err = fr.ParseHeaders(ctx, currFileInfo.HeaderBytes)
			if err != nil {
				continue
			}
			currFileInfo.Status = FILE_STATUS_HEADERS_PARSE
			// 8. add to result
			fr.ReadFiles[fileName] = currFileInfo
			result = append(result, currFileInfo)
		}
	}
	return result, nil
}

func (_ *FileReader) ReadFile(ctx context.Context, fileName string) ([]byte, error) {
	logger := zerolog.Ctx(ctx).With().Str("fileName", fileName).Logger()
	result, err := os.ReadFile(fileName)
	if err != nil {
		logger.Error().Err(err).Msg("os.ReadFile")
		return nil, err
	}
	logger.Debug().Bytes("fileBytes", result).Msg("ReadFile")
	return result, nil
}

func (_ *FileReader) ParseHeaders(ctx context.Context, headerBytes []byte) (map[string]string, error) {
	logger := zerolog.Ctx(ctx)
	result := make(map[string]string, 0)
	lines := strings.Split(string(headerBytes), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "H??") {
			line = strings.TrimPrefix(line, "H??")
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headerKey := strings.TrimSpace(parts[0])
				headerValue := strings.TrimSpace(parts[1])
				result[headerKey] = headerValue
			}
		}
	}
	logger.Debug().Interface("headers", result).Msg("ParseHeaders")
	return result, nil
}
