package ds

import (
	"sync"

	"github.com/miekg/dns"
)

func initiateLocalQuery(m *dns.Msg, cache *CacheManager) *dns.Msg {
	wg := &sync.WaitGroup{}
	rch := make(chan *queryResult, 1)
	defer close(rch)
	wg.Add(1)
	go ResolveWithLocalNS(wg, rch, m)
	wg.Wait()

	localNSResult := <-rch
	cache.Set(localNSResult.Answer)
	return localNSResult.Answer
}

func initiateDohQuery(m *dns.Msg, cache *CacheManager) *dns.Msg {
	wg := &sync.WaitGroup{}
	rch := make(chan *queryResult, 1)
	defer close(rch)
	wg.Add(1)
	go ResolveWithDoh(wg, rch, m)
	wg.Wait()

	dohResult := <-rch
	cache.Set(dohResult.Answer)
	return dohResult.Answer
}

func initiateQuery(m *dns.Msg, cache *CacheManager, whitelist *CacheManager, blacklist *CacheManager) *dns.Msg {
	wg := &sync.WaitGroup{}
	rch := make(chan *queryResult, 6)
	defer close(rch)
	wg.Add(6)
	go ResolveWithHoney(wg, rch, m)
	go ResolveWithHoney(wg, rch, m)
	go ResolveWithHoney(wg, rch, m)
	go ResolveWithHoney(wg, rch, m)
	go ResolveWithLocalNS(wg, rch, m)
	go ResolveWithDoh(wg, rch, m)
	wg.Wait()

	var r *queryResult
	var dohResult *queryResult
	var localNSResult *queryResult

	score := 0

	for i := 0; i < 6; i++ {
		r = <-rch
		switch r.QueryType {
		case QueryTypeHoneyPot:
			if r.Answer != nil {
				score -= 1
			} else {
				score += 1
			}
		case QueryTypeDoH:
			dohResult = r
		case QueryTypeLocalNS:
			localNSResult = r
		}
	}

	if score < 0 {
		lg.Debugf("to blacklist: %s", m.Question[0].Name)
		cache.Set(dohResult.Answer)
		blacklist.SetWithTTL(m, 11800)
		whitelist.Delete(m)
		return dohResult.Answer
	} else {
		lg.Debugf("to whitelist: %s", m.Question[0].Name)
		cache.Set(localNSResult.Answer)
		whitelist.SetWithTTL(m, 11800)
		blacklist.Delete(m)
		return localNSResult.Answer
	}
}
