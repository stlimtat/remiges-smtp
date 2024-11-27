package input

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type FileStatus uint8

const (
	FILE_STATUS_DONE       FileStatus = 0
	FILE_STATUS_PROCESSING FileStatus = 1
)

type FileInfo struct {
	Name   string
	FileID string
	Status FileStatus
}

type FileReader struct {
	InPath       string
	PollInterval time.Duration
	ReadFiles    map[string]FileInfo
	ticker       *time.Ticker
}

func NewFileReader(
	_ context.Context,
	inPath string,
	pollInterval time.Duration,
) *FileReader {
	result := &FileReader{
		InPath:       inPath,
		PollInterval: pollInterval,
		ReadFiles:    make(map[string]FileInfo, 0),
		ticker:       time.NewTicker(pollInterval),
	}

	return result
}

func (fr *FileReader) Run(
	ctx context.Context,
) error {
	defer fr.ticker.Stop()
	logger := zerolog.Ctx(ctx)

	// https://blog.devtrovert.com/p/select-and-for-range-channel-i-bet
outerloop:
	for {
		select {
		case t := <-fr.ticker.C:
			// check file exists
			logger.Info().Time("t", t).Msg("ticker.C")
		case <-ctx.Done():
			logger.Debug().Msg("ctx.Done")
			break outerloop
		}
	}
	return nil
}

func (fr *FileReader) Process(
	ctx context.Context,
) error {
	logger := zerolog.Ctx(ctx)
	// 1. read list of files in directory
	entries, err := os.ReadDir(fr.InPath)
	if err != nil {
		logger.Error().Err(err).Msg("os.ReadDir")
	}
	// 2. check with are new files
	for _, e := range entries {
		// 3. read newest file - message, rcpt
		_, ok := fr.ReadFiles[e.Name()]
		if !ok {
			logger.Error().Err(fmt.Errorf("fileInfo error")).Msg("fr.ReadFiles")
			continue
		}
		if strings.HasPrefix(e.Name(), "df") {
			fileID := e.Name()[2:]
			fr.ReadFiles[e.Name()] = FileInfo{
				FileID: fileID,
				Name:   e.Name(),
				Status: FILE_STATUS_PROCESSING,
			}
			// 4. check if qf file exists
			qfFileName := strings.Replace(e.Name(), "df", "qf", 1)
			_, err := os.Stat(qfFileName)
			if err != nil {
				logger.Error().Err(err).Msg("qfExists")
				continue
			}
			fileBytes, err := os.ReadFile(e.Name())
			if err != nil {
				logger.Error().Err(fmt.Errorf("readfile error")).Msg("os.ReadFile")
				continue
			}
			logger.Debug().Bytes("fileBytes", fileBytes)
			// 5. send message to rcpt via smtpclient
		}
	}
	return nil
}
