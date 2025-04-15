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
	"github.com/stlimtat/remiges-smtp/internal/errors"
	"github.com/stlimtat/remiges-smtp/internal/utils"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

const (
	DKIMProcessorType = "dkim"
)

// DKIMProcessor handles DKIM signing of emails
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
			true, false,
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

// Process handles DKIM signing of an email
func (p *DKIMProcessor) Process(ctx context.Context, mail *pmail.Mail) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)

	// Validate input
	if mail == nil {
		return nil, errors.NewError(errors.ErrMailProcessing, "mail cannot be nil", nil)
	}

	if p.DomainCfg == nil || p.DomainCfg.DKIM == nil {
		return nil, errors.NewError(errors.ErrDKIMConfig, "DKIM configuration not initialized", nil)
	}

	// Validate From address
	if mail.From.String() == "" {
		return nil, errors.NewError(errors.ErrMailProcessing, "from address required for DKIM signing", nil)
	}

	logger.Debug().
		Str("from", mail.From.String()).
		Msg("Starting DKIM signing process")

	// Process DKIM signing
	signedMail, err := p.signMail(ctx, mail)
	if err != nil {
		return nil, errors.NewError(errors.ErrMailProcessing, "failed to sign mail", err)
	}

	return signedMail, nil
}

// signMail performs the actual DKIM signing
func (p *DKIMProcessor) signMail(ctx context.Context, mail *pmail.Mail) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().
		Interface("from", mail.From).
		Msg("DKIMProcessor")

	canonical := mox.CanonicalLocalpart(
		mail.From.Localpart,
		p.DomainCfg.MoxDomain,
	)

	mailMsg := mail.Headers
	mailMsg = append(mailMsg, mail.Body...)
	mailMsg = append(mailMsg, []byte("\r\n\r\n")...)
	selectors := slices.Collect(
		maps.Values(p.DomainCfg.DKIM.Selectors),
	)

	dkimHeaders, err := moxDkim.Sign(
		ctx,
		p.SLogger,
		canonical,
		mail.From.Domain,
		selectors,
		true,
		bytes.NewReader(mailMsg),
	)
	if err != nil {
		logger.Error().Err(err).Msg("DKIMProcessor: sign")
		return mail, err
	}
	// add dkim headers to the mail
	dkimHeaderParts := strings.SplitN(dkimHeaders, ":", 2)
	dkimHeaderKey := dkimHeaderParts[0]
	dkimHeaderValue := strings.Join(dkimHeaderParts[1:], "")
	dkimHeaderValue = strings.TrimSpace(dkimHeaderValue)
	if mail.HeadersMap == nil {
		mail.HeadersMap = make(map[string][]byte)
	}
	mail.HeadersMap[dkimHeaderKey] = []byte(dkimHeaderValue)

	logger.Info().
		Str("dkim_header_key", dkimHeaderKey).
		Str("dkim_header_value", dkimHeaderValue).
		Msg("DKIMProcessor: Process.Done")

	return mail, nil
}
