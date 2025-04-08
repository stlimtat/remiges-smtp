// Package dkim provides functionality for generating and managing DomainKeys Identified Mail (DKIM)
// signatures and DNS records. This package handles the creation of DKIM TXT records for DNS
// configuration, supporting multiple key types and formats.
package dkim

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"slices"
	"strings"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/crypto"
)

// TxtGen is a service for generating DKIM TXT records for DNS configuration.
// It handles the conversion of public keys into the appropriate DNS record format
// and ensures proper formatting according to DKIM specifications.
type TxtGen struct{}

// Generate creates a DKIM TXT record for DNS configuration using the provided parameters.
// The generated record follows the format specified in RFC 6376 and includes the public key
// in the appropriate format for the specified key type.
//
// Parameters:
//   - ctx: Context for logging and cancellation
//   - domain: The domain for which the DKIM record is being generated
//   - keyType: The type of key being used (e.g., "rsa", "ed25519")
//   - selector: The DKIM selector used to identify the key
//   - pubKeyPEM: The public key in PEM format
//
// Returns:
//   - []byte: The generated DKIM TXT record in DNS format
//   - error: Non-nil if generation fails, with specific error messages for:
//   - Empty public key
//   - Invalid PEM format
//   - Unsupported key type
//   - Key parsing errors
//
// Example output format:
//
//	selector._domainkey.example.com IN TXT "v=DKIM1; k=rsa; p=base64encodedkey"
func (_ *TxtGen) Generate(
	ctx context.Context,
	domain, keyType, selector string,
	pubKeyPEM []byte,
) ([]byte, error) {
	domain = strings.TrimSpace(domain)
	keyType = strings.TrimSpace(keyType)
	selector = strings.TrimSpace(selector)

	if keyType == "" || !slices.Contains(crypto.ValidKeyTypes, keyType) {
		keyType = crypto.KeyTypeRSA
	}
	if selector == "" {
		selector = "default"
	}

	logger := zerolog.Ctx(ctx).
		With().
		Str("domain", domain).
		Str("key_type", keyType).
		Str("selector", selector).
		Logger()

	if domain == "" {
		logger.Error().Msg("domain is empty")
		return nil, fmt.Errorf("domain is empty")
	}

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
