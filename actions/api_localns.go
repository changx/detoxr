package actions

import (
	"net"
	"time"

	"github.com/changx/detoxr/ds"
	"github.com/gobuffalo/buffalo"
	"github.com/miekg/dns"
)

func GetLocalNS(c buffalo.Context) error {
	q := SettingsQuery{
		Server: ds.GetLocalNS(),
	}
	return successResponse(c, q)
}

func TryLocalNS(c buffalo.Context) error {
	r := SettingsQuery{}

	if err := c.Bind(&r); err != nil {
		return serverErrorWithShortMessage(c, "illegal request", err.Error())
	}

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(r.Host), dns.TypeA)
	msg.Opcode = dns.OpcodeQuery
	msg.RecursionDesired = true

	c.Logger().Debugf("query msg: %+v", msg)
	client := new(dns.Client)
	client.Dialer = &net.Dialer{
		Timeout: 220 * time.Millisecond,
	}

	response, _, err := client.Exchange(msg, r.Server)

	if response != nil && response.Answer != nil {
		for _, answer := range response.Answer {
			if answer.Header().Rrtype == dns.TypeA {
				c.Logger().Debugf("answer: %+v", answer)
				if r.Response == nil {
					r.Response = make([]string, 0)
				}
				r.Response = append(r.Response, answer.String())
			}
		}
	}
	if err != nil {
		r.Error = err.Error()
	}

	return successResponse(c, r)
}

func SaveLocalNS(c buffalo.Context) error {
	r := SettingsQuery{}

	if err := c.Bind(&r); err != nil {
		return serverErrorWithShortMessage(c, "illegal request", err.Error())
	}

	if r.Server != "" {
		ds.SetLocalNS(r.Server)
	}

	return successResponse(c, r)
}
