package ds

import (
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var honeyPotClient *dns.Client

func createHoneyPotClient() {
	honeyPotClient = new(dns.Client)
	honeyPotClient.Dialer = &net.Dialer{
		Timeout: 220 * time.Millisecond,
	}
}

func ResolveWithHoney(wg *sync.WaitGroup, rc chan *queryResult, m *dns.Msg) *dns.Msg {
	if wg != nil {
		defer wg.Done()
	}

	in, _, err := honeyPotClient.Exchange(m, GetHoneypotIP())

	if err == nil && in != nil {
		hit := false
		for _, answer := range in.Answer {
			lg.Debugf("honeypot ans: %+v", answer)
			if answer.Header().Rrtype == dns.TypeA {
				hit = true
				break
			}
		}

		if hit {
			in.SetReply(m)
		} else {
			in = nil
		}
	} else {
		in = nil
	}

	if rc != nil {
		qr := queryResult{
			QueryType: QueryTypeHoneyPot,
			Answer:    in,
		}

		rc <- &qr
	}

	return in
}
