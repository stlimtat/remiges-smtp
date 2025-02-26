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
	CachedMX map[string]dn.MXRecord
	Slogger  *slog.Logger
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
	result.CachedMX = make(map[string]dn.MXRecord)
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

	// 1. check if mx already exists in cache
	result, ok := r.CachedMX[domain.ASCII]
	if !ok {
		// 2. resolve the mx record for the domain
		ipDomain := dns.IPDomain{
			Domain: domain,
		}

		_, _, _, expandedNextHop, hosts, _, err := smtpclient.GatherDestinations(
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
			r.CachedMX[expandedNextHop.ASCII] = dn.MXRecord{
				Domain:  expandedNextHop.ASCII,
				Entries: hosts,
				Hosts:   hostStrSlice,
			}
		}
		r.CachedMX[domain.ASCII] = dn.MXRecord{
			Domain:  domain.ASCII,
			Entries: hosts,
			Hosts:   hostStrSlice,
		}
		result = r.CachedMX[domain.ASCII]
	}
	logger.Info().
		Strs("hosts", result.Hosts).
		Msg("lookupMX")
	return result.Hosts, nil
}
