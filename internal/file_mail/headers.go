package file_mail

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/pkg/input"
)

const (
	HeadersTransformerType = "headers"
)

type HeadersTransformer struct {
	Cfg config.FileMailConfig
}

func (h *HeadersTransformer) Init(
	_ context.Context,
	cfg config.FileMailConfig,
) error {
	h.Cfg = cfg
	return nil
}

func (h *HeadersTransformer) Index() int {
	return h.Cfg.Index
}

func (_ *HeadersTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	inMail *mail.Mail,
) (outMail *mail.Mail, err error) {
	logger := zerolog.Ctx(ctx).With().Str("qf_file_path", fileInfo.QfFilePath).Logger()
	logger.Info().Msg("HeadersTransformer")

	// 1. check fileInfo.QfReader is not nil
	if fileInfo.QfReader == nil {
		return nil, fmt.Errorf("fileInfo.QfReader is nil")
	}

	// 1. initialize the headers map in the mail
	inMail.Headers = make(map[string][]byte)

	// 2. read all the bytes from the qf file
	byteSlice, err := io.ReadAll(fileInfo.QfReader)
	if err != nil {
		return nil, err
	}
	fileInfo.Status = input.FILE_STATUS_HEADERS_READ

	// 3. split the bytes into lines
	var lines [][]byte
	if bytes.Contains(byteSlice, []byte("\r\n")) {
		lines = bytes.Split(byteSlice, []byte("\r\n"))
	} else {
		lines = bytes.Split(byteSlice, []byte("\n"))
	}

	// 4. iterate over the lines and add them to the result map
	for _, line := range lines {
		// 4. split the line into key and value with the colon as the delimiter
		if len(line) < 1 || !bytes.Contains(line, []byte(":")) {
			continue
		}
		kvPair := bytes.Split(line, []byte(":"))
		key := bytes.TrimSpace(kvPair[0])
		keyStr := string(key)

		value := bytes.TrimSpace(kvPair[1])
		inMail.Headers[keyStr] = value
	}
	fileInfo.Status = input.FILE_STATUS_HEADERS_PARSE

	return inMail, nil
}
