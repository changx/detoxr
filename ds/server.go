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
	initServerConfig()

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

	honeyMsg := new(dns.Msg)
	honeyMsg.SetQuestion(dns.Fqdn("m.root-servers.net"), dns.TypeA)
	honeyMsg.RecursionDesired = true

	r, _ := dns.Exchange(honeyMsg, "114.114.114.114:53")

	if r != nil && r.Rcode == dns.RcodeSuccess {
		for _, r_ := range r.Answer {
			if a, ok := r_.(*dns.A); ok {
				hip := a.A.String()
				SetHoneypotIP(hip + ":53")
				break
			}
		}
	}

	if GetHoneypotIP() == "" {
		SetHoneypotIP(envy.Get("HONEYPOT_ADDR", "202.12.27.3:53"))
	}

	createHoneyPotClient()

	createLocalNSClient()
	SetLocalNS(envy.Get("FASTNS_ADDR", "114.114.114.114:53"))

	createDoHClient()
	SetDohService(envy.Get("DOH_SERVICE", "https://dns.google/resolve"))

	server := dns.Server{Addr: envy.Get("NS_ADDR", ":1053"), Net: "udp", ReusePort: true}

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
