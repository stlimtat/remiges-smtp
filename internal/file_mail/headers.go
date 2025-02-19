package file_mail

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"github.com/stlimtat/remiges-smtp/pkg/input"
)

const (
	HeadersTransformerType = "headers"
	HeadersConfigArgPrefix = "prefix"
)

type HeadersTransformer struct {
	Cfg         config.FileMailConfig
	PrefixStr   string
	PrefixBytes []byte
}

func (h *HeadersTransformer) Init(
	_ context.Context,
	cfg config.FileMailConfig,
) error {
	h.Cfg = cfg
	var ok bool
	h.PrefixStr, ok = cfg.Args[HeadersConfigArgPrefix]
	if !ok {
		h.PrefixStr = ""
	}
	h.PrefixBytes = []byte(h.PrefixStr)
	return nil
}

func (h *HeadersTransformer) Index() int {
	return h.Cfg.Index
}

func (h *HeadersTransformer) Transform(
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
	inMail.Metadata = make(map[string][]byte)

	// 2. read all the bytes from the qf file
	byteSlice, err := io.ReadAll(fileInfo.QfReader)
	if err != nil {
		return nil, err
	}
	fileInfo.Status = input.FILE_STATUS_HEADERS_READ
	// Replace all \n with \r\n
	re := regexp.MustCompile(`\r?\n`)
	byteSlice = re.ReplaceAll(byteSlice, []byte("\r\n"))

	// 3. split the bytes into lines
	lines := bytes.Split(byteSlice, []byte("\r\n"))

	// 4. iterate over the lines and add them to the result map
	var key []byte
	var keyStr string
	var keyPrefixStr string
	var value []byte
	for _, line := range lines {
		// 4. split the line into key and value with the colon as the delimiter
		if len(line) < 1 {
			continue
		}
		if !bytes.Contains(line, []byte(":")) {
			// 5. if the line does not have a colon, attach the next line to the value
			line = bytes.TrimSpace(line)
			value = append(value, line...)
			inMail.Metadata[keyStr] = value
			if bytes.HasPrefix(key, h.PrefixBytes) {
				inMail.Metadata[keyPrefixStr] = value
			}
			continue
		}
		kvPair := bytes.Split(line, []byte(":"))
		key = bytes.TrimSpace(kvPair[0])
		value = bytes.TrimSpace(kvPair[1])
		keyStr = string(key)

		inMail.Metadata[keyStr] = value
		// 6. if the key starts with the prefix, add it to the result map
		if bytes.HasPrefix(key, h.PrefixBytes) {
			keyPrefix := bytes.TrimPrefix(key, h.PrefixBytes)
			keyPrefixStr = string(keyPrefix)
			inMail.Metadata[keyPrefixStr] = value
		}
	}
	fileInfo.Status = input.FILE_STATUS_HEADERS_PARSE

	return inMail, nil
}
