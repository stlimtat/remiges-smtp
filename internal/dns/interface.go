// Package dns provides DNS resolution functionality for SMTP operations.
// It includes interfaces and implementations for DNS record lookups,
// particularly focused on MX (Mail Exchange) record resolution.
package dns

import (
	"context"

	moxDns "github.com/mjl-/mox/dns"
)

//go:generate mockgen -destination=mox_mock.go -package=dns github.com/mjl-/mox/dns Resolver
//go:generate mockgen -destination=mock.go -package=dns . IResolver

// IResolver defines the interface for DNS resolution operations.
// Implementations of this interface provide methods to look up various DNS records.
type IResolver interface {
	// LookupMX performs a DNS lookup for MX (Mail Exchange) records for the given domain.
	// It returns a list of hostnames that are configured to receive email for the domain,
	// ordered by preference (lower numbers indicate higher preference).
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - domain: The domain to look up MX records for
	//
	// Returns:
	//   - []string: List of hostnames configured as mail servers for the domain
	//   - error: Non-nil if the lookup fails, with specific error messages for:
	//     - Network errors
	//     - DNS resolution failures
	//     - Invalid domain names
	//     - Timeout errors
	LookupMX(ctx context.Context, domain moxDns.Domain) ([]string, error)
}
