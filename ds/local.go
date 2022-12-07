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

func ResolveWithLocalNS(wg *sync.WaitGroup, rc chan *queryResult, m *dns.Msg) *dns.Msg {
	if wg != nil {
		defer wg.Done()
	}

	in, _, err := localNSClient.Exchange(m, GetLocalNS())

	if in != nil && err == nil {
		in.SetReply(m)
	} else {
		lg.Errorf("query local ns err: %s %s", m.String(), err.Error())
		in = new(dns.Msg)
		in.SetRcode(m, dns.RcodeNameError)
	}

	if rc != nil {
		qs := queryResult{
			QueryType: QueryTypeLocalNS,
			Answer:    in,
		}

		rc <- &qs
	}

	return in
}
