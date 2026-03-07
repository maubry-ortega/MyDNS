package server

import (
	"testing"
	"time"

	"MyDNS/pkg/cache"
	"MyDNS/pkg/forwarder"
	"MyDNS/pkg/resolver"

	"github.com/miekg/dns"
)

func TestServer_Integration(t *testing.T) {
	// Setup components
	c := cache.New()
	f := forwarder.New([]string{"8.8.8.8:53"})
	r := resolver.New("my.os.", "192.168.1.100")
	
	// Start server on a high random port
	addr := "127.0.0.1:53530"
	s := New(addr, r, f, c)

	go func() {
		_ = s.ListenAndServe()
	}()

	// Give server a moment to start
	time.Sleep(200 * time.Millisecond)

	// Test Internal Resolution
	client := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion("test.my.os.", dns.TypeA)

	resp, _, err := client.Exchange(msg, addr)
	if err != nil {
		t.Fatalf("Failed to query internal domain: %v", err)
	}

	if len(resp.Answer) != 1 {
		t.Fatalf("Expected 1 internal answer, got %d", len(resp.Answer))
	}
	
	aRecord, ok := resp.Answer[0].(*dns.A)
	if !ok || aRecord.A.String() != "192.168.1.100" {
		t.Errorf("Expected 192.168.1.100, got %v", resp.Answer[0])
	}

	// Test External Resolution (Forwarding)
	msgExt := new(dns.Msg)
	msgExt.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)
	msgExt.RecursionDesired = true

	respExt, _, err := client.Exchange(msgExt, addr)
	if err != nil {
		t.Fatalf("Failed to query external domain: %v", err)
	}

	if len(respExt.Answer) == 0 {
		t.Fatalf("Expected external answers, got 0")
	}
}
