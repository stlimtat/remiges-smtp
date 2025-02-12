package input

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/mjl-/mox/smtp"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/mail"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MailTransformer struct {
	Casers cases.Caser
	Cfg    config.ReadFileConfig
	ToAddr string
}

func NewMailTransformer(
	_ context.Context,
	cfg config.ReadFileConfig,
) *MailTransformer {
	result := &MailTransformer{
		Casers: cases.Title(language.English),
		Cfg:    cfg,
	}

	return result
}

func (t *MailTransformer) WithToAddr(
	toAddr string,
) *MailTransformer {
	t.ToAddr = toAddr
	return t
}

func (t *MailTransformer) Transform(
	ctx context.Context,
	fileInfo *FileInfo,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().
		Str("id", fileInfo.ID).
		Msg("Transforming mail")

	result := &mail.Mail{}
	headers, err := t.ReadHeaders(ctx, fileInfo)
	if err != nil {
		return nil, err
	}

	result.From, err = t.ReadFrom(ctx, headers)
	if err != nil {
		return nil, err
	}

	result.To, err = t.ReadTo(ctx, headers)
	if err != nil {
		return nil, err
	}

	result.Body, err = t.ReadBody(ctx, fileInfo)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *MailTransformer) ReadFrom(
	_ context.Context,
	headers map[string][]byte,
) (smtp.Address, error) {
	var from string
	switch t.Cfg.FromType {
	case config.FromTypeHeaders:
		// combine the df and qf files
		from = string(headers["From"])
		from = strings.TrimSpace(from)
	case config.FromTypeDefault:
		// use the default from
		from = t.Cfg.DefaultFrom
	}
	result, err := smtp.ParseAddress(from)
	return result, err
}

func (t *MailTransformer) ReadTo(
	_ context.Context,
	headers map[string][]byte,
) (smtp.Address, error) {
	to := string(headers["To"])
	if to == "" {
		to = t.ToAddr
	}
	to = strings.TrimSpace(to)
	result, err := smtp.ParseAddress(to)
	return result, err
}

func (_ *MailTransformer) ReadBody(
	ctx context.Context,
	fileInfo *FileInfo,
) ([]byte, error) {
	logger := zerolog.Ctx(ctx).With().Str("df_file_path", fileInfo.DfFilePath).Logger()
	logger.Info().Msg("Reading body")

	result, err := io.ReadAll(fileInfo.DfReader)
	if err != nil {
		return nil, err
	}
	fileInfo.Status = FILE_STATUS_BODY_READ

	return result, nil
}

func (t *MailTransformer) ReadHeaders(
	ctx context.Context,
	fileInfo *FileInfo,
) (map[string][]byte, error) {
	logger := zerolog.Ctx(ctx).With().Str("qf_file_path", fileInfo.QfFilePath).Logger()
	logger.Info().Msg("Reading headers")

	// 1. read all the bytes from the qf file
	result := make(map[string][]byte)
	byteSlice, err := io.ReadAll(fileInfo.QfReader)
	if err != nil {
		return nil, err
	}
	fileInfo.Status = FILE_STATUS_HEADERS_READ

	// 2. split the bytes into lines
	var lines [][]byte
	if bytes.Contains(byteSlice, []byte("\r\n")) {
		lines = bytes.Split(byteSlice, []byte("\r\n"))
	} else {
		lines = bytes.Split(byteSlice, []byte("\n"))
	}

	// 3. iterate over the lines and add them to the result map
	for _, line := range lines {
		// 4. split the line into key and value with the colon as the delimiter
		if len(line) < 1 || !bytes.Contains(line, []byte(":")) {
			continue
		}
		kvPair := bytes.Split(line, []byte(":"))
		key := bytes.TrimSpace(kvPair[0])
		keyStr := t.Casers.String(string(key))

		value := bytes.TrimSpace(kvPair[1])
		result[keyStr] = value
	}
	fileInfo.Status = FILE_STATUS_HEADERS_PARSE

	return result, nil
}
