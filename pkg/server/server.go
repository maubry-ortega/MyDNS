package server

import (
	"log"

	"MyDNS/pkg/cache"
	"MyDNS/pkg/forwarder"
	"MyDNS/pkg/resolver"

	"github.com/miekg/dns"
)

type Server struct {
	addr      string
	resolver  *resolver.Resolver
	forwarder *forwarder.Forwarder
	cache     *cache.Cache
}

func New(addr string, r *resolver.Resolver, f *forwarder.Forwarder, c *cache.Cache) *Server {
	return &Server{
		addr:      addr,
		resolver:  r,
		forwarder: f,
		cache:     c,
	}
}

func (s *Server) ListenAndServe() error {
	dns.HandleFunc(".", s.handleDNSRequest)

	udpServer := &dns.Server{Addr: s.addr, Net: "udp"}
	tcpServer := &dns.Server{Addr: s.addr, Net: "tcp"}

	errChan := make(chan error, 1)

	go func() {
		log.Printf("Listening on UDP %s", s.addr)
		errChan <- udpServer.ListenAndServe()
	}()

	go func() {
		log.Printf("Listening on TCP %s", s.addr)
		errChan <- tcpServer.ListenAndServe()
	}()

	return <-errChan
}

func (s *Server) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		dns.HandleFailed(w, r)
		return
	}

	q := r.Question[0]
	cacheKey := q.Name + ":" + dns.TypeToString[q.Qtype]

	// 1. Check Cache
	if msg := s.cache.Get(cacheKey); msg != nil {
		msg.Id = r.Id
		w.WriteMsg(msg)
		return
	}

	var resp *dns.Msg
	var err error

	// 2. Check Internal Resolver (Authoritative for my.os)
	if s.resolver.Handles(q.Name) {
		resp = new(dns.Msg)
		resp.SetReply(r)
		resp.Authoritative = true
		resp.Answer = s.resolver.Resolve(q)
	} else {
		// 3. Forward External Queries
		resp, err = s.forwarder.Forward(r)
		if err != nil {
			log.Printf("Forward error: %v", err)
			dns.HandleFailed(w, r)
			return
		}

		// 4. Cache External Response
		if len(resp.Answer) > 0 {
			// Use the min TTL from answers for caching
			minTTL := resp.Answer[0].Header().Ttl
			for _, rr := range resp.Answer {
				if rr.Header().Ttl < minTTL {
					minTTL = rr.Header().Ttl
				}
			}
			s.cache.Set(cacheKey, resp, minTTL)
		}
	}

	w.WriteMsg(resp)
}
