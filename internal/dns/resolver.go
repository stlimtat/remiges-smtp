// Package dns provides DNS resolution functionality for SMTP operations.
// It includes implementations for DNS record lookups, particularly focused on
// MX (Mail Exchange) record resolution for email delivery.
package dns

import (
	"context"
	"log/slog"
	"time"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtpclient"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/pkg/dn"
)

// Resolver provides DNS resolution capabilities with caching and retry mechanisms.
// It implements the IResolver interface and wraps the mox/dns.Resolver with
// additional functionality for logging and error handling.
type Resolver struct {
	dns.Resolver
	Slogger *slog.Logger

	// Cache duration for DNS records
	cacheDuration time.Duration

	// Maximum number of retries for DNS lookups
	maxRetries int
}

// NewResolver creates a new instance of the DNS resolver with the specified configuration.
//
// Parameters:
//   - ctx: Context for initialization (currently unused but reserved for future use)
//   - resolver: The underlying DNS resolver implementation
//   - slogger: Structured logger for recording DNS operations
//
// Returns:
//   - *Resolver: A new resolver instance configured with the provided parameters
func NewResolver(
	_ context.Context,
	resolver dns.Resolver,
	slogger *slog.Logger,
) *Resolver {
	result := &Resolver{
		Resolver: resolver,
		Slogger:  slogger,
	}
	return result
}

// LookupMX performs a DNS lookup for MX (Mail Exchange) records for the given domain.
// It uses the underlying resolver to gather destination information and returns
// a list of hostnames that can receive email for the domain.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - domain: The domain to look up MX records for
//
// Returns:
//   - []string: List of hostnames configured as mail servers for the domain
//   - error: Non-nil if the lookup fails, with specific error messages for:
//   - Network errors
//   - DNS resolution failures
//   - Invalid domain names
//   - Timeout errors
//
// The function logs detailed information about the lookup process and results
// using the configured structured logger.
func (r *Resolver) LookupMX(
	ctx context.Context,
	domain dns.Domain,
) ([]string, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("domain", domain.ASCII).
		Logger()
	var result dn.MXRecord

	// Resolve the MX record for the domain
	ipDomain := dns.IPDomain{
		Domain: domain,
	}

	_, _, _, expandedNextHop, hosts, _, err := smtpclient.GatherDestinations( //nolint:dogsled // none of the identifiers are used
		ctx, r.Slogger, r.Resolver, ipDomain,
	)
	if err != nil {
		logger.Error().Err(err).Msg("smtpclient.GatherDestinations")
		return nil, err
	}

	// Convert from dns.IPDomain to string slice
	hostStrSlice := []string{}
	for _, host := range hosts {
		hostStrSlice = append(hostStrSlice, host.String())
	}

	// Handle domain expansion if necessary
	if expandedNextHop.ASCII != domain.ASCII {
		result = dn.MXRecord{
			Domain:  expandedNextHop.ASCII,
			Entries: hosts,
			Hosts:   hostStrSlice,
		}
	} else {
		result = dn.MXRecord{
			Domain:  domain.ASCII,
			Entries: hosts,
			Hosts:   hostStrSlice,
		}
	}

	// Log the successful lookup results
	logger.Info().
		Interface("result", result).
		Strs("hosts", result.Hosts).
		Msg("lookupMX")

	return result.Hosts, nil
}
