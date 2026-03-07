package resolver

import (
	"testing"

	"github.com/miekg/dns"
)

func TestResolver_Handles(t *testing.T) {
	r := New("my.os.", "192.168.1.100")

	if !r.Handles("app.my.os.") {
		t.Errorf("Expected Handles to return true for 'app.my.os.'")
	}
	if r.Handles("google.com.") {
		t.Errorf("Expected Handles to return false for 'google.com.'")
	}
}

func TestResolver_Resolve_A_Record(t *testing.T) {
	r := New("my.os.", "192.168.1.100")

	q := dns.Question{Name: "app.my.os.", Qtype: dns.TypeA, Qclass: dns.ClassINET}
	answers := r.Resolve(q)

	if len(answers) != 1 {
		t.Fatalf("Expected 1 answer, got %d", len(answers))
	}

	aRecord, ok := answers[0].(*dns.A)
	if !ok {
		t.Fatalf("Expected dns.A record, got %T", answers[0])
	}

	if aRecord.A.String() != "192.168.1.100" {
		t.Errorf("Expected IP 192.168.1.100, got %s", aRecord.A.String())
	}
}

func TestResolver_Resolve_AAAA_Record(t *testing.T) {
    // Resolver should return nil since targetIP is IPv4
	r := New("my.os.", "192.168.1.100")

	q := dns.Question{Name: "app.my.os.", Qtype: dns.TypeAAAA, Qclass: dns.ClassINET}
	answers := r.Resolve(q)

	if len(answers) != 0 {
		t.Fatalf("Expected 0 answers for IPv4 target when asking for AAAA, got %d", len(answers))
	}
}

func TestResolver_Resolve_UnsupportedType(t *testing.T) {
	r := New("my.os.", "192.168.1.100")

	q := dns.Question{Name: "app.my.os.", Qtype: dns.TypeTXT, Qclass: dns.ClassINET}
	answers := r.Resolve(q)

	if answers != nil {
		t.Errorf("Expected nil answer for unsupported type (TXT), got %v", answers)
	}
}
