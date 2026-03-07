package resolver

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

type Resolver struct {
	targetIP net.IP
	zone     string
}

func New(zone string, targetIP string) *Resolver {
	if !strings.HasSuffix(zone, ".") {
		zone += "."
	}
	return &Resolver{
		targetIP: net.ParseIP(targetIP),
		zone:     zone,
	}
}

func (r *Resolver) Resolve(q dns.Question) []dns.RR {
	if q.Qtype != dns.TypeA && q.Qtype != dns.TypeAAAA {
		return nil
	}

	var rr dns.RR
	if q.Qtype == dns.TypeA && r.targetIP.To4() != nil {
		rr = &dns.A{
			Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   r.targetIP.To4(),
		}
	} else if q.Qtype == dns.TypeAAAA && r.targetIP.To4() == nil {
		rr = &dns.AAAA{
			Hdr:  dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
			AAAA: r.targetIP,
		}
	}

	if rr != nil {
		return []dns.RR{rr}
	}
	return nil
}

func (r *Resolver) Handles(name string) bool {
	return strings.HasSuffix(name, r.zone)
}
