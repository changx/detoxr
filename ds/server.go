package ds

import (
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	"github.com/miekg/dns"
)

var lg logger.FieldLogger
var whitelist *CacheManager
var blacklist *CacheManager

func Serve() {
	if serverCfg == nil {
		InitServerConfig()
	}

	// init memory cache
	cache := NewCacheManager()
	go cache.CacheSuperVisor()

	whitelist = NewCacheManager()
	go whitelist.CacheSuperVisor()

	blacklist = NewCacheManager()
	go blacklist.CacheSuperVisor()

	var ENV = envy.Get("GO_ENV", "development")

	lg = JSONLogger(logger.DebugLevel)
	if ENV == "production" {
		lg = JSONLogger(logger.InfoLevel)
	}

	createHoneyPotClient()
	createLocalNSClient()
	createDoHClient()

	persistentBlackList := GetPersistentBlackList()
	for i := 0; i < len(persistentBlackList); i++ {
		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(persistentBlackList[i]), dns.TypeA)
		m.Opcode = dns.OpcodeQuery
		m.RecursionDesired = true
		blacklist.SetCacheWithKey(m, -1)
	}

	server := dns.Server{Addr: serverCfg.NSAddr, Net: "udp", ReusePort: true}

	server.NotifyStartedFunc = func() {
		lg.Info("ds started")
	}
	dns.HandleFunc(".", func(w dns.ResponseWriter, m *dns.Msg) {
		var answer *dns.Msg

		if len(m.Question) > 1 {
			answer = new(dns.Msg)
			answer.SetReply(m)
			answer.SetRcode(m, dns.RcodeServerFailure)
			w.WriteMsg(answer)
			return
		}

		answer = cache.Get(m)
		if answer != nil {
			answer.Id = m.Id
			w.WriteMsg(answer)
			return
		}

		// in whitelist
		answer = whitelist.Get(m)
		if answer != nil {
			lg.Debugf("whitelist hit: %s", m.Question[0].Name)
			if answer = initiateLocalQuery(m, cache); answer == nil {
				answer = new(dns.Msg)
				answer.SetReply(m)
				answer.SetRcode(m, dns.RcodeServerFailure)
				w.WriteMsg(answer)
				return
			} else {
				w.WriteMsg(answer)
				return
			}
		}

		// in blacklist
		answer = blacklist.Get(m)
		if answer != nil {
			lg.Debugf("blacklist hit: %s", m.Question[0].Name)
			if answer = initiateDohQuery(m, cache); answer == nil {
				answer = new(dns.Msg)
				answer.SetReply(m)
				answer.SetRcode(m, dns.RcodeServerFailure)
				w.WriteMsg(answer)
				return
			} else {
				w.WriteMsg(answer)
				return
			}
		}

		lg.Debugf("detecting %s", m.Question[0].Name)
		if answer = initiateQuery(m, cache, whitelist, blacklist); answer == nil {
			answer = new(dns.Msg)
			answer.SetReply(m)
			answer.SetRcode(m, dns.RcodeServerFailure)
			w.WriteMsg(answer)
			return
		} else {
			w.WriteMsg(answer)
			return
		}
	})

	server.ListenAndServe()
}
