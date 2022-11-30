package ds

import (
	"fmt"

	"github.com/miekg/dns"
)

type config struct {
	honeypotIP string
	localNS    string
	dohService string
}

var resolverConfig *config

func Serve() {
	// init memory cache
	cache := NewCacheManager()
	go cache.CacheSuperVisor()

	// init server configuration
	resolverConfig = &config{}

	// cfip, err := net.LookupIP("b.root-servers.net")
	// if err != nil {
	// 	return
	// }

	// SetHoneypotIP(cfip[0].String() + ":53")
	SetHoneypotIP("47.90.125.25:53")
	createHoneyPotClient()

	createLocalNSClient()
	SetLocalNS("119.29.29.29:53")

	createDoHClient()
	SetDohService("https://dns.google/resolve")

	server := dns.Server{Addr: ":1053", Net: "udp", ReusePort: true}

	server.NotifyStartedFunc = func() {
		fmt.Printf("ds started\n")
	}
	dns.HandleFunc(".", func(w dns.ResponseWriter, m *dns.Msg) {
		fmt.Printf("calling HandleFunc\n")
		var answer *dns.Msg

		if len(m.Question) > 1 {
			answer = new(dns.Msg)
			answer.SetReply(m)
			answer.SetRcode(m, dns.RcodeServerFailure)
			w.WriteMsg(answer)
			return
		}

		answer = cache.Get(m)
		if answer == nil {
			if answer = initiateQuery(m, cache); answer == nil {
				answer = new(dns.Msg)
				answer.SetReply(m)
				answer.SetRcode(m, dns.RcodeServerFailure)
			}
		} else {
			fmt.Printf("cache hit %+v\n", answer)
			answer.Id = m.Id
		}

		w.WriteMsg(answer)
	})

	server.ListenAndServe()
}
