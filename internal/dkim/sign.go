package dkim

import (
	"context"

	"github.com/mjl-/mox/dns"
	mox "github.com/mjl-/mox/mox-"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/mail"
)

type DKIMSigner struct {
	Cfg       config.DKIMConfig
	Domain    dns.Domain
	DomainCfg config.DomainConfig
}

func NewDKIMSigner(
	_ context.Context,
	cfg config.DKIMConfig,
	domain dns.Domain,
) (*DKIMSigner, error) {
	result := &DKIMSigner{
		Cfg:    cfg,
		Domain: domain,
	}
	return result, nil
}

func (s *DKIMSigner) Sign(
	ctx context.Context,
	msg *mail.Mail,
) (*mail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("DKIMSigner: Sign")

	_ = mox.CanonicalLocalpart(msg.From.Localpart, s.DomainCfg.Domain)

	return msg, nil
}
