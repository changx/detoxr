package ds

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var honeyPotClient *dns.Client

func createHoneyPotClient() {
	honeyPotClient = new(dns.Client)
	honeyPotClient.Dialer = &net.Dialer{
		Timeout: 100 * time.Millisecond,
	}
}

func ResolveWithHoney(wg *sync.WaitGroup, rc chan *queryResult, m *dns.Msg) {
	defer wg.Done()

	fmt.Printf("using honeyPot addr: %s\n", resolverConfig.honeypotIP)
	in, _, err := honeyPotClient.Exchange(m, resolverConfig.honeypotIP)

	answer := in

	fmt.Printf("honeyPot respond: \n%+v err = %+v\n", in, err)
	if err == nil && answer != nil {
		answer.SetReply(m)
	} else {
		answer = nil
	}

	qr := queryResult{
		QueryType: QueryTypeHoneyPot,
		Answer:    answer,
	}

	rc <- &qr
}
