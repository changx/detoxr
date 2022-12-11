package ds

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

/*
{
  "Status": 0,  // NOERROR - Standard DNS response code (32 bit integer).
  "TC": false,  // Whether the response is truncated
  "RD": true,   // Always true for Google Public DNS
  "RA": true,   // Always true for Google Public DNS
  "AD": false,  // Whether all response data was validated with DNSSEC
  "CD": false,  // Whether the client asked to disable DNSSEC
  "Question":
  [
    {
      "name": "apple.com.",  // FQDN with trailing dot
      "type": 1              // A - Standard DNS RR type
    }
  ],
  "Answer":
  [
    {
      "name": "apple.com.",   // Always matches name in the Question section
      "type": 1,              // A - Standard DNS RR type
      "TTL": 3599,            // Record's time-to-live in seconds
      "data": "17.178.96.59"  // Data for A - IP address as text
    },
    {
      "name": "apple.com.",
      "type": 1,
      "TTL": 3599,
      "data": "17.172.224.47"
    },
    {
      "name": "apple.com.",
      "type": 1,
      "TTL": 3599,
      "data": "17.142.160.59"
    }
  ],
  "edns_client_subnet": "12.34.56.78/0"  // IP address / scope prefix-length
}
*/

const DohResponseNoError = 0
const DohResponseServFail = 2

type DohQueryQuestion struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type DohQueryAnswer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

type DoHMessage struct {
	Status           int                `json:"Status"`
	TC               bool               `json:"TC"`
	RD               bool               `json:"RD"`
	RA               bool               `json:"RA"`
	AD               bool               `json:"AD"`
	CD               bool               `json:"CD"`
	Questions        []DohQueryQuestion `json:"Question"`
	Answers          []DohQueryAnswer   `json:"Answer"`
	Comment          string             `json:"Comment"`
	EDnsClientSubnet string             `json:"edns_client_subnet"`
}

var dohClient *http.Client

func createDoHClient() {
	transport := &http.Transport{
		Proxy: nil,
		DialContext: func() func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   time.Second,
				KeepAlive: 600 * time.Second,
			}
			return dialer.DialContext
		}(),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   100,
		DisableKeepAlives:     false,
		DisableCompression:    false,
		ResponseHeaderTimeout: 2 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	dohClient = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

func toUint16(s string) uint16 {
	x, _ := strconv.ParseUint(s, 10, 0)
	return uint16(x)
}

func fillRRHeader(r *DohQueryAnswer, hdr *dns.RR_Header) {
	switch r.Type {
	case RR_A:
		hdr.Rrtype = dns.TypeA
	case RR_CNAME:
		hdr.Rrtype = dns.TypeCNAME
	case RR_SOA:
		hdr.Rrtype = dns.TypeSOA
	case RR_NS:
		hdr.Rrtype = dns.TypeNS
	case RR_MX:
		hdr.Rrtype = dns.TypeMX
	}
	hdr.Rdlength = uint16(len(r.Data))
	hdr.Class = dns.ClassINET
	hdr.Ttl = uint32(r.TTL)
	hdr.Name = r.Name
}

func ResolveWithDoh(wg *sync.WaitGroup, rc chan *queryResult, m *dns.Msg) *dns.Msg {
	if wg != nil {
		defer wg.Done()
	}

	params := url.Values{}
	params.Add("name", m.Question[0].Name)
	params.Add("type", fmt.Sprintf("%d", m.Question[0].Qtype))

	req, _ := http.NewRequest(http.MethodGet, GetDohService(), nil)
	req.URL.RawQuery = params.Encode()

	answer := new(dns.Msg)

	resp, err := dohClient.Do(req)
	if err != nil {
		answer.SetRcode(m, dns.RcodeServerFailure)
	} else {
		defer resp.Body.Close()
		rbody, _ := io.ReadAll(resp.Body)
		dohMsg := DoHMessage{
			Questions: make([]DohQueryQuestion, 0),
			Answers:   make([]DohQueryAnswer, 0),
		}

		json.Unmarshal(rbody, &dohMsg)

		answer.SetReply(m)

		for _, r := range dohMsg.Answers {
			switch r.Type {
			case RR_A:
				a := new(dns.A)
				fillRRHeader(&r, &a.Hdr)
				a.A = net.ParseIP(r.Data)
				answer.Answer = append(answer.Answer, a)
			case RR_CNAME:
				a := new(dns.CNAME)
				fillRRHeader(&r, &a.Hdr)
				a.Target = r.Data
				answer.Answer = append(answer.Answer, a)
			case RR_MX:
				a := new(dns.MX)
				fillRRHeader(&r, &a.Hdr)
				mx := strings.Split(r.Data, " ")
				pref, _ := strconv.Atoi(mx[0])
				a.Preference = uint16(pref)
				a.Mx = mx[1]
				answer.Answer = append(answer.Answer, a)
			case RR_NS:
				a := new(dns.NS)
				fillRRHeader(&r, &a.Hdr)
				a.Ns = r.Data
				answer.Answer = append(answer.Answer, a)
			case RR_AAAA:
				a := new(dns.AAAA)
				fillRRHeader(&r, &a.Hdr)
				a.AAAA = net.ParseIP(r.Data)
				answer.Answer = append(answer.Answer, a)
			case RR_SOA:
				a := new(dns.SOA)
				fillRRHeader(&r, &a.Hdr)
				answer.Answer = append(answer.Answer, a)
			case RR_SRV:
				a := new(dns.SRV)
				fillRRHeader(&r, &a.Hdr)

				srvParts := strings.Split(r.Data, " ")
				a.Priority = toUint16(srvParts[0])
				a.Weight = toUint16(srvParts[1])
				a.Port = toUint16(srvParts[2])
				a.Target = srvParts[3]

				answer.Answer = append(answer.Answer, a)
			}
		}

		lg.Debugf("doh answer: %s", answer.String())
	}

	if rc != nil {
		rc <- &queryResult{
			QueryType: QueryTypeDoH,
			Answer:    answer,
		}
	}
	return answer
}
