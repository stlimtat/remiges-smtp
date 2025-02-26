package dn

import (
	"github.com/mjl-/adns"
	"github.com/mjl-/mox/dns"
)

type MXRecord struct {
	ADNSResult adns.Result
	Domain     string
	Entries    []dns.IPDomain
	Hosts      []string
}
