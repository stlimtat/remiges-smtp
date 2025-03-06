package dkim

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

type TxtGen struct{}

func (_ *TxtGen) Generate(
	ctx context.Context,
	domain, keyType, selector string,
	pubKeyPEM []byte,
) ([]byte, error) {
	selector = strings.TrimSpace(selector)
	domain = strings.TrimSpace(domain)
	keyType = strings.TrimSpace(keyType)

	logger := zerolog.Ctx(ctx).
		With().
		Str("domain", domain).
		Str("key_type", keyType).
		Str("selector", selector).
		Logger()

	if len(pubKeyPEM) < 1 {
		logger.Error().Msg("pubKeyPEM is empty")
		return nil, fmt.Errorf("pubKeyPEM is empty")
	}

	block, _ := pem.Decode(pubKeyPEM)
	if block == nil {
		logger.Error().Msg("pubKeyPEM is not PEM formatted")
		return nil, fmt.Errorf("pubKeyPEM is not PEM formatted")
	}

	switch keyType {
	case "ed25519":
		break
	default:
		_, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			logger.Error().Err(err).Msg("x509.ParsePKCS1PublicKey")
			return nil, fmt.Errorf("x509.ParsePKCS1PublicKey: %w", err)
		}
	}

	resultKey := base64.StdEncoding.EncodeToString(block.Bytes)

	result := fmt.Sprintf("\"v=DKIM1; k=%s; p=%s\"", keyType, resultKey)
	result = fmt.Sprintf("%s._domainkey.%s IN TXT %s", selector, domain, result)

	return []byte(result), nil
}
