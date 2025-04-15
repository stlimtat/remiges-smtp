package file_mail

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/input"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
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
	prefixAny, ok := cfg.Args[HeadersConfigArgPrefix]
	if !ok {
		prefixAny = ""
	}
	h.PrefixStr = prefixAny.(string)
	h.PrefixBytes = []byte(h.PrefixStr)
	return nil
}

func (h *HeadersTransformer) Index() int {
	return h.Cfg.Index
}

func (h *HeadersTransformer) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	inMail *pmail.Mail,
) (outMail *pmail.Mail, err error) {
	logger := zerolog.Ctx(ctx).With().Str("qf_file_path", fileInfo.QfFilePath).Logger()
	logger.Debug().Msg("HeadersTransformer")

	// 1. check fileInfo.QfReader is not nil
	if fileInfo.QfFilePath == "" {
		logger.Error().Msg("ToSkip: fileInfo.QfFilePath is empty")
		return nil, fmt.Errorf("ToSkip: fileInfo.QfFilePath is empty")
	}

	// 2. validate the qf file exists and is readable
	_, err = utils.ValidateIO(ctx, fileInfo.QfFilePath, true, false)
	if err != nil {
		logger.Error().Err(err).Msg("utils.ValidateIO")
		return nil, err
	}

	// 3. read all the bytes from the qf file
	byteSlice, err := os.ReadFile(fileInfo.QfFilePath)
	if err != nil {
		logger.Error().Err(err).Msg("os.ReadFile")
		return nil, err
	}
	fileInfo.Status = input.FILE_STATUS_HEADERS_READ

	// 4. initialize the headers map in the mail
	inMail.Metadata = make(map[string][]byte)
	// 5. replace all \n with \r\n
	re := regexp.MustCompile(`\r?\n`)
	byteSlice = re.ReplaceAll(byteSlice, []byte("\r\n"))

	// 6. split the bytes into lines
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
