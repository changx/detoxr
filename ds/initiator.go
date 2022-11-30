package ds

import (
	"fmt"
	"sync"

	"github.com/miekg/dns"
)

const (
	QueryTypeHoneyPot = iota
	QueryTypeDoH
	QueryTypeLocalNS
)

type queryResult struct {
	QueryType int
	Answer    *dns.Msg
}

func initiateQuery(m *dns.Msg, cache *cacheManager) *dns.Msg {
	wg := &sync.WaitGroup{}
	rch := make(chan *queryResult, 3)
	defer close(rch)
	wg.Add(3)
	go ResolveWithHoney(wg, rch, m)
	go ResolveWithLocalNS(wg, rch, m)
	go ResolveWithDoh(wg, rch, m)
	wg.Wait()

	var r *queryResult
	var dohResult *queryResult
	var localResult *queryResult
	var honeyPotResult *queryResult

	for i := 0; i < 3; i++ {
		r = <-rch
		switch r.QueryType {
		case QueryTypeHoneyPot:
			honeyPotResult = r
		case QueryTypeDoH:
			dohResult = r
		case QueryTypeLocalNS:
			localResult = r
		}
	}

	if honeyPotResult.Answer != nil {
		fmt.Printf("honeyPot got respond: \n%+v\n", honeyPotResult.Answer)
		cache.Set(dohResult.Answer)
		return dohResult.Answer
	} else {
		fmt.Printf("local ns respond: \n%+v\n", localResult.Answer)
		cache.Set(localResult.Answer)
		return localResult.Answer
	}
}
