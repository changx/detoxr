package actions

import (
	"net"
	"time"

	"github.com/changx/detoxr/ds"
	"github.com/gobuffalo/buffalo"
	"github.com/miekg/dns"
)

type SettingsQuery struct {
	Server   string   `json:"server"`
	Servers  []string `json:"servers"`
	Host     string   `json:"host,omitempty"`
	Response []string `json:"answer,omitempty"`
	Error    string   `json:"error,omitempty"`
}

func GetHoneypot(c buffalo.Context) error {
	r := SettingsQuery{
		Server: ds.GetHoneypotIP(),
	}

	return successResponse(c, r)
}

func TryHoneypot(c buffalo.Context) error {
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

func SaveHoneypot(c buffalo.Context) error {
	r := SettingsQuery{}

	if err := c.Bind(&r); err != nil {
		return serverErrorWithShortMessage(c, "illegal request", err.Error())
	}

	if r.Server != "" {
		ds.SetHoneypotIP(r.Server)
	}

	return successResponse(c, r)
}
