package ds

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var localNSClient *dns.Client
var errNsTimeout = errors.New("local ns timeout")

func createLocalNSClient() {
	localNSClient = new(dns.Client)
	localNSClient.Dialer = &net.Dialer{
		Timeout: 500 * time.Millisecond,
	}
}

func ResolveWithLocalNS(wg *sync.WaitGroup, rc chan *queryResult, m *dns.Msg) {
	if wg != nil {
		defer wg.Done()
	}

	nsList := GetLocalNS()

	var err error
	var result *queryResult

	for i := 0; i < len(nsList); i++ {
		result, err = lookupWithServer(m, nsList[i])
		if err == nil {
			break
		}
	}

	if result == nil {
		answer := new(dns.Msg)
		answer.SetRcode(m, dns.RcodeServerFailure)
		result = &queryResult{
			QueryType: QueryTypeLocalNS,
			Answer:    answer,
		}
	}

	if rc != nil {
		rc <- result
	}
}

func lookupWithServer(m *dns.Msg, server string) (*queryResult, error) {
	in, _, err := localNSClient.Exchange(m, server)

	if in != nil && err == nil {
		in.SetReply(m)
	} else {
		lg.Errorf("query local ns err: %s %s", m.String(), err.Error())
		return nil, errNsTimeout
	}

	return &queryResult{
		QueryType: QueryTypeLocalNS,
		Answer:    in,
	}, nil
}
