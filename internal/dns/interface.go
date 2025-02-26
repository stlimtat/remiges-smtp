package dns

import (
	"context"

	moxDns "github.com/mjl-/mox/dns"
)

//go:generate mockgen -destination=mox_mock.go -package=dns github.com/mjl-/mox/dns Resolver
//go:generate mockgen -destination=mock.go -package=dns . IResolver
type IResolver interface {
	LookupMX(ctx context.Context, domain moxDns.Domain) ([]string, error)
}
