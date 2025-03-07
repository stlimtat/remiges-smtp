package dns

import (
	"context"
	"log/slog"

	"github.com/mjl-/mox/dns"
	"github.com/mjl-/mox/smtpclient"
	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/pkg/dn"
)

type Resolver struct {
	dns.Resolver
	Slogger *slog.Logger
}

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

func (r *Resolver) LookupMX(
	ctx context.Context,
	domain dns.Domain,
) ([]string, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("domain", domain.ASCII).
		Logger()
	var result dn.MXRecord

	// 2. resolve the mx record for the domain
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
	// 3. convert from dns.IPDomain to string
	hostStrSlice := []string{}
	for _, host := range hosts {
		hostStrSlice = append(hostStrSlice, host.String())
	}
	if expandedNextHop.ASCII != domain.ASCII {
		result = dn.MXRecord{
			Domain:  expandedNextHop.ASCII,
			Entries: hosts,
			Hosts:   hostStrSlice,
		}
	}
	result = dn.MXRecord{
		Domain:  domain.ASCII,
		Entries: hosts,
		Hosts:   hostStrSlice,
	}
	logger.Info().
		Interface("result", result).
		Strs("hosts", result.Hosts).
		Msg("lookupMX")
	return result.Hosts, nil
}
