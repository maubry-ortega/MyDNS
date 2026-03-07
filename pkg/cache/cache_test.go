package cache

import (
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestCache_SetAndGet(t *testing.T) {
	c := New()

	msg := new(dns.Msg)
	msg.SetQuestion("example.com.", dns.TypeA)
	msg.Answer = append(msg.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
	})

	c.Set("example.com.:A", msg, 60)

	cachedMsg := c.Get("example.com.:A")
	if cachedMsg == nil {
		t.Fatal("Expected cached message, got nil")
	}

	if len(cachedMsg.Answer) != 1 {
		t.Fatalf("Expected 1 answer, got %d", len(cachedMsg.Answer))
	}
}

func TestCache_Expiration(t *testing.T) {
	c := New()

	msg := new(dns.Msg)
	msg.SetQuestion("example.com.", dns.TypeA)

	c.Set("example.com.:A", msg, 1) // 1 second TTL

	// Should be available immediately
	cachedMsg := c.Get("example.com.:A")
	if cachedMsg == nil {
		t.Fatal("Expected cached message, got nil")
	}

	// Wait for expiration
	time.Sleep(2 * time.Second)

	cachedMsg = c.Get("example.com.:A")
	if cachedMsg != nil {
		t.Fatal("Expected expired message to be nil, but got a result")
	}
}

func TestCache_Miss(t *testing.T) {
	c := New()

	cachedMsg := c.Get("nonexistent.com.:A")
	if cachedMsg != nil {
		t.Fatal("Expected nil for non-existent key")
	}
}
