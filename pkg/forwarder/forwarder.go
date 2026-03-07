package forwarder

import (
	"fmt"
	"github.com/miekg/dns"
)

type Forwarder struct {
	nameservers []string
	client      *dns.Client
}

func New(nameservers []string) *Forwarder {
	return &Forwarder{
		nameservers: nameservers,
		client:      &dns.Client{},
	}
}

func (f *Forwarder) Forward(msg *dns.Msg) (*dns.Msg, error) {
	for _, ns := range f.nameservers {
		resp, _, err := f.client.Exchange(msg, ns)
		if err == nil {
			return resp, nil
		}
	}
	return nil, fmt.Errorf("failed to forward query to any nameservers")
}
