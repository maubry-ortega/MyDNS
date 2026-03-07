package forwarder

import (
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestForwarder_Forward(t *testing.T) {
	// Use Cloudflare and Google DNS for testing forwarding
	f := New([]string{"1.1.1.1:53", "8.8.8.8:53"})

	// Reduce timeout for the test to avoid long waits if network is down
	f.client.Timeout = 2 * time.Second

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)
	msg.RecursionDesired = true

	resp, err := f.Forward(msg)
	
	if err != nil {
		t.Fatalf("Failed to forward request: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if len(resp.Answer) == 0 {
		t.Fatal("Expected answers for example.com, got 0")
	}

	// Verify we got an A record
	hasARecord := false
	for _, ans := range resp.Answer {
		if _, ok := ans.(*dns.A); ok {
			hasARecord = true
			break
		}
	}

	if !hasARecord {
		t.Fatal("Expected at least one A record in response")
	}
}

func TestForwarder_Forward_FailsWithBadServer(t *testing.T) {
	// Use an unroutable IP
	f := New([]string{"192.0.2.1:53"})
	f.client.Timeout = 100 * time.Millisecond // fast fail

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)

	_, err := f.Forward(msg)
	if err == nil {
		t.Fatal("Expected forwarding to fail with bad server, but it succeeded")
	}
}
