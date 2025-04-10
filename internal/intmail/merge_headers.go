package intmail

import (
	"context"
	"slices"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	MergeHeadersProcessorType = "mergeHeaders"
)

type MergeHeadersProcessor struct {
	Cfg config.MailProcessorConfig
}

func (p *MergeHeadersProcessor) Init(
	_ context.Context,
	cfg config.MailProcessorConfig,
) error {
	p.Cfg = cfg
	return nil
}

func (p *MergeHeadersProcessor) Index() int {
	return p.Cfg.Index
}

func (_ *MergeHeadersProcessor) Process(
	ctx context.Context,
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("MergeHeadersProcessor: Process")

	result := make([]byte, 0)
	keys := make([]string, 0)
	for key := range inMail.HeadersMap {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for _, key := range keys {
		result = append(
			result,
			[]byte(key+": "+string(inMail.HeadersMap[key])+"\r\n")...,
		)
	}

	result = append(result, []byte("\r\n")...)

	inMail.Headers = result

	return inMail, nil
}
