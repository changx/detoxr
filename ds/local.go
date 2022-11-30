package ds

import (
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var localNSClient *dns.Client

func createLocalNSClient() {
	localNSClient = new(dns.Client)
	localNSClient.Dialer = &net.Dialer{
		Timeout: time.Second,
	}
}

func ResolveWithLocalNS(wg *sync.WaitGroup, rc chan *queryResult, m *dns.Msg) {
	defer wg.Done()

	in, _, err := localNSClient.Exchange(m, resolverConfig.localNS)

	if in != nil && err == nil {
		in.SetReply(m)
	} else {
		in = new(dns.Msg)
		in.SetRcode(m, dns.RcodeNameError)
	}

	qs := queryResult{
		QueryType: QueryTypeLocalNS,
		Answer:    in,
	}

	rc <- &qs
}
