package intmail

import (
	"bytes"
	"context"
	"log/slog"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	moxDkim "github.com/mjl-/mox/dkim"
	"github.com/mjl-/mox/mox-"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/crypto"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	DKIMProcessorType = "dkim"
)

type DKIMProcessor struct {
	Cfg       config.MailProcessorConfig
	DomainCfg *config.DomainConfig
	SLogger   *slog.Logger
}

func (p *DKIMProcessor) Init(
	ctx context.Context,
	cfg config.MailProcessorConfig,
) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("DKIMProcessor Init")
	p.Cfg = cfg

	p.DomainCfg = &config.DomainConfig{}
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Metadata:   nil,
			DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
			Result:     &p.DomainCfg,
		},
	)
	if err != nil {
		logger.Error().Err(err).Msg("DKIMProcessor: NewDecoder")
		return err
	}
	err = decoder.Decode(p.Cfg.Args)
	if err != nil {
		logger.Error().Err(err).Msg("DKIMProcessor: decode")
		return err
	}

	err = p.DomainCfg.Transform(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("DKIMProcessor: transform")
		return err
	}

	return nil
}

func (p *DKIMProcessor) Index() int {
	return p.Cfg.Index
}

func (p *DKIMProcessor) InitDKIMCrypto(
	ctx context.Context,
	loader crypto.IKeyLoader,
) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("InitDKIMCrypto")

	moxSelectors := p.DomainCfg.DKIM.MoxSelectors

	for selectorName, moxSelector := range moxSelectors {
		privateKeyPath, err := utils.ValidateIO(
			ctx,
			filepath.Clean(moxSelector.PrivateKeyFile),
			true,
		)
		if err != nil {
			logger.Error().Err(err).Msg("InitDKIMCrypto: ValidateIO")
			return err
		}
		signer, err := loader.LoadPrivateKey(
			ctx,
			moxSelector.Algorithm,
			privateKeyPath,
		)
		if err != nil {
			logger.Error().Err(err).Msg("InitDKIMCrypto: LoadPrivateKey")
			return err
		}
		moxDkimSelector := p.DomainCfg.DKIM.Selectors[selectorName]
		moxDkimSelector.PrivateKey = signer
		p.DomainCfg.DKIM.Selectors[selectorName] = moxDkimSelector
	}

	return nil
}

func (p *DKIMProcessor) Process(
	ctx context.Context,
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Interface("from", inMail.From).
		Msg("DKIMProcessor")

	canonical := mox.CanonicalLocalpart(
		inMail.From.Localpart,
		p.DomainCfg.MoxDomain,
	)

	mailMsg := inMail.Headers
	mailMsg = append(mailMsg, inMail.Body...)
	mailMsg = append(mailMsg, []byte("\r\n\r\n")...)
	selectors := slices.Collect(
		maps.Values(p.DomainCfg.DKIM.Selectors),
	)

	dkimHeaders, err := moxDkim.Sign(
		ctx,
		p.SLogger,
		canonical,
		inMail.From.Domain,
		selectors,
		true,
		bytes.NewReader(mailMsg),
	)
	if err != nil {
		logger.Error().Err(err).Msg("DKIMProcessor: sign")
		return inMail, err
	}
	// add dkim headers to the mail
	dkimHeaderParts := strings.SplitN(dkimHeaders, ":", 2)
	dkimHeaderKey := dkimHeaderParts[0]
	dkimHeaderValue := strings.Join(dkimHeaderParts[1:], "")
	dkimHeaderValue = strings.TrimSpace(dkimHeaderValue)
	if inMail.HeadersMap == nil {
		inMail.HeadersMap = make(map[string][]byte)
	}
	inMail.HeadersMap[dkimHeaderKey] = []byte(dkimHeaderValue)

	logger.Info().
		Str("dkim_header_key", dkimHeaderKey).
		Str("dkim_header_value", dkimHeaderValue).
		Msg("DKIMProcessor: Process.Done")

	return inMail, nil
}
