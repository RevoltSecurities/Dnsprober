package dns

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/utils"
	"github.com/miekg/dns"
)

type Dnsclient struct {
	Resolvers  []string
	Maxretries int
	Client     *dns.Client
	TCPClient  *dns.Client
}


func New(resolvers []string, maxretries int, timeout int) *Dnsclient {
	client := &dns.Client{
		Net:     "",
		Timeout: time.Duration(timeout) * time.Second,
	}

	tcpclient := &dns.Client{
		Net:     "tcp",
		Timeout: 3 * time.Second,
	}

	return &Dnsclient{
		Resolvers:  resolvers,
		Maxretries: maxretries,
		Client:     client,
		TCPClient:  tcpclient,
	}
}

func (d *Dnsclient) Query(hostname string, dnsrequest string, resolvers ...string) (*DNSResponse, error) {
	var (
		err         error
		dnsresponse DNSResponse
		dnsTypes    []uint16
	)

	dnsmsg := new(dns.Msg)
	dnsmsg.Id = dns.Id()
	dnsmsg.SetEdns0(4096, false)

	dnsTypes, err = d.QueryType(dnsrequest)
	if err != nil {
		return &dnsresponse, err
	}

	for _, dnsType := range dnsTypes {
		domain := dns.Fqdn(hostname)
		dnsmsg.Question = make([]dns.Question, 1)

		switch dnsType {
		case dns.TypeAXFR:
			dnsmsg.SetAxfr(domain)
		case dns.TypePTR:
			domain, err = d.ReverseAddress(DOTtrimers(domain))
			if err != nil {
				return nil, err
			}
			fallthrough
		default:
			dnsmsg.RecursionDesired = true
			dnsquestion := dns.Question{
				Name:   domain,
				Qtype:  dnsType,
				Qclass: dns.ClassINET,
			}
			dnsmsg.Question[0] = dnsquestion
		}

		var (
			i        int
			response *dns.Msg
		)

		for i = 0; i < d.Maxretries; i++ {
			var resolver string
			if len(resolvers) > 0 {
				resolver = resolvers[i%len(resolvers)] + ":53"
			} else {
				resolver = d.GetResolver()
			}

			if dnsType == dns.TypeAXFR {
				connection, err := d.TCPClient.Dial(resolver)
				if err != nil {
					continue
				}
				defer connection.Close()

				transfer := &dns.Transfer{Conn: connection}
				axfrResp, err := transfer.In(dnsmsg, resolver)
				if err != nil {
					continue
				}
				err = dnsresponse.EnvelopeParser(axfrResp)
				break

			} else {
				response, _, err = d.Client.Exchange(dnsmsg, resolver)
				if response != nil && response.Truncated {
					response, _, err = d.TCPClient.Exchange(dnsmsg, resolver)
				}

				if response == nil {
					continue
				}

				err = dnsresponse.MSGParser(response)

				dnsresponse.RawResponse = response
				dnsresponse.Host = hostname
				dnsresponse.StatusCode = dns.RcodeToString[response.Rcode]
				dnsresponse.RawStatusCode = response.Rcode
				dnsresponse.Raw += response.String()
				dnsresponse.Resolver = append(dnsresponse.Resolver, resolver)
				dnsresponse.SliceDeduplication()

				if response.Rcode == dns.RcodeSuccess {
					break
				}
			}
		}
	}
	return &dnsresponse, err
}

func (d *Dnsclient) Queries(hostnames []string, dnsrequest string, resolvers ...string) ([]*DNSResponse, error) { // implemented variadic parameters
	var wg sync.WaitGroup
	responses := make([]*DNSResponse, len(hostnames))
	errors := make([]error, len(hostnames))

	for i, hostname := range hostnames {
		wg.Add(1)
		var (
			resp *DNSResponse
			err  error
		)
		go func(i int, hostname string) {
			defer wg.Done()
			if len(resolvers) == 0 {
				resp, err = d.Query(hostname, dnsrequest)
			} else {
				resp, err = d.Query(hostname, dnsrequest, resolvers...)
			}
			responses[i] = resp
			errors[i] = err
		}(i, hostname)
	}
	wg.Wait()

	for _, err := range errors {
		if err != nil {
			return responses, err
		}
	}
	return responses, nil
}

func (d *Dnsclient) Wildcard(domain string) (*DNSResponse, error) {
	randowords := utils.GenerateRandomWords(10)
	fulldomain := fmt.Sprintf("%s.%s", randowords, domain)
	response, err := d.Query(fulldomain, "A")
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (d *Dnsclient) QueryType(requestType string) ([]uint16, error) {
	switch strings.ToUpper(requestType) {
	case "A":
		return []uint16{dns.TypeA}, nil
	case "AAAA":
		return []uint16{dns.TypeAAAA}, nil
	case "MX":
		return []uint16{dns.TypeMX}, nil
	case "CNAME":
		return []uint16{dns.TypeCNAME}, nil
	case "TXT":
		return []uint16{dns.TypeTXT}, nil
	case "NS":
		return []uint16{dns.TypeNS}, nil
	case "SOA":
		return []uint16{dns.TypeSOA}, nil
	case "ANY":
		return []uint16{dns.TypeANY}, nil
	case "PTR":
		return []uint16{dns.TypePTR}, nil
	case "SRV":
		return []uint16{dns.TypeSRV}, nil
	case "AXFR":
		return []uint16{dns.TypeAXFR}, nil
	case "CAA":
		return []uint16{dns.TypeCAA}, nil
	case "RESOLVE":
		return []uint16{dns.TypeA, dns.TypeAAAA}, nil
	case "ALL":
		return []uint16{dns.TypeA, dns.TypeAAAA, dns.TypeMX, dns.TypeCNAME, dns.TypeTXT, dns.TypeNS, dns.TypeSOA, dns.TypeANY, dns.TypePTR, dns.TypeSRV, dns.TypeCAA}, nil
	default:
		return nil, fmt.Errorf("%s record is undefined", requestType)
	}
}

func (d *Dnsclient) ReverseAddress(hostname string) (string, error) {
	if ip := net.ParseIP(hostname); ip != nil {
		reversedAddr, err := dns.ReverseAddr(hostname)
		if err != nil {
			return "", fmt.Errorf("failed to generate reverse address: %s", err.Error())
		}
		return reversedAddr, nil
	}
	return dns.Fqdn(hostname), nil
}

func (d *Dnsclient) Resolve(hostname string) ([]string, error) {
	var ips []string
	response, err := d.Query(hostname, "resolve")
	if err != nil {
		return nil, err
	}
	ips = append(ips, response.A...)
	ips = append(ips, response.AAAA...)
	return utils.Deduplslice(ips), nil
}

func (d *Dnsclient) A(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "a")
	if err != nil {
		return nil, err
	}
	return response.A, nil
}

func (d *Dnsclient) AAAA(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "aaaa")
	if err != nil {
		return nil, err
	}
	return response.AAAA, nil
}

func (d *Dnsclient) CNAME(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "cname")
	if err != nil {
		return nil, err
	}
	return response.CNAME, nil
}

func (d *Dnsclient) SOA(hostname string) ([]SOARecords, error) {
	response, err := d.Query(hostname, "soa")
	if err != nil {
		return nil, err
	}
	return response.SOA, nil
}

func (d *Dnsclient) MX(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "mx")
	if err != nil {
		return nil, err
	}
	return response.MX, nil
}

func (d *Dnsclient) NS(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "ns")
	if err != nil {
		return nil, err
	}
	return response.NS, nil
}

func (d *Dnsclient) TXT(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "txt")
	if err != nil {
		return nil, err
	}
	return response.TXT, nil
}

func (d *Dnsclient) SRV(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "srv")
	if err != nil {
		return nil, err
	}
	return response.SRV, nil
}

func (d *Dnsclient) PTR(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "ptr")
	if err != nil {
		return nil, err
	}
	return response.PTR, nil
}

func (d *Dnsclient) CAA(hostname string) ([]string, error) {
	response, err := d.Query(hostname, "caa")
	if err != nil {
		return nil, err
	}
	return response.CAA, nil
}

func (d *Dnsclient) ANY(hostname string) (*DNSResponse, error) {
	response, err := d.Query(hostname, "any")
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (d *Dnsclient) AXFR(hostname string) (*AXFRResponse, error) {
	nsdomains, err := d.NS(hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NS records: %v", err)
	}

	var resolvers []string
	for _, nsdomain := range nsdomains {
		_, err := d.A(nsdomain)
		if err != nil {
			continue
		}
		resolvers = append(resolvers, nsdomain)
	}

	var dnsresponses []*DNSResponse
	var success bool
	for _, resolver := range resolvers {
		resolver = dns.Fqdn(resolver)
		axfrdata, err := d.Query(hostname, "axfr", resolver)
		if err != nil || axfrdata == nil {
			continue
		}

		validRecords := 0
		if len(axfrdata.NS) != 0 {
			validRecords++
		}

		if validRecords > 0 {
			dnsresponses = append(dnsresponses, axfrdata)
			success = true
		}
	}

	if !success {
		return nil, fmt.Errorf("zone transfer failed: no valid AXFR records found")
	}

	return &AXFRResponse{Host: hostname, DNSResponses: dnsresponses}, nil
}

func (response *DNSResponse) JSONIZE() (string, error) {
	jsoned, err := json.Marshal(response)
	return string(jsoned), err
}

func (d *Dnsclient) GetResolver() string {
	resolver := d.Resolvers[rand.Intn(len(d.Resolvers))]
	return resolver + ":53"
}

func DOTtrimers(s string) string {
	return strings.TrimRight(s, ".")
}

type DNSResponse struct {
	Timestamp     time.Time     `json:"timestamp,omitempty"`
	Host          string        `json:"host,omitempty"`
	TTL           uint32        `json:"ttl,omitempty"`
	Resolver      []string      `json:"resolver_records,omitempty"`
	A             []string      `json:"a_records,omitempty"`
	AAAA          []string      `json:"aaaa_records,omitempty"`
	CNAME         []string      `json:"cname_records,omitempty"`
	MX            []string      `json:"mx_records,omitempty"`
	PTR           []string      `json:"ptr_records,omitempty"`
	SOA           []SOARecords  `json:"soa_records,omitempty"`
	NS            []string      `json:"ns_records,omitempty"`
	TXT           []string      `json:"txt_records,omitempty"`
	SRV           []string      `json:"srv_records,omitempty"`
	CAA           []string      `json:"caa_records,omitempty"`
	AXFR          *AXFRResponse `json:"axfr_records,omitempty"`
	AllRecords    []string      `json:"all_records,omitempty"`
	Raw           string        `json:"raw,omitempty"`
	StatusCode    string        `json:"status_code,omitempty"`
	RawStatusCode int           `json:"raw_status_code,omitempty"`
	RawResponse   *dns.Msg      `json:"raw_response,omitempty"`
}

type SOARecords struct {
	Name    string `json:"name,omitempty"`
	NS      string `json:"ns,omitempty"`
	Mbox    string `json:"mailbox,omitempty"`
	Serial  uint32 `json:"serial,omitempty"`
	Refresh uint32 `json:"refresh,omitempty"`
	Retry   uint32 `json:"retry,omitempty"`
	Expire  uint32 `json:"expire,omitempty"`
	Minttl  uint32 `json:"minttl,omitempty"`
}

type AXFRResponse struct {
	Host         string         `json:"host,omitempty"`
	DNSResponses []*DNSResponse `json:"chain_response,omitempty"`
}

func (dnsr *DNSResponse) MSGParser(response *dns.Msg) error {
	allRecords := append(response.Answer, response.Extra...)
	allRecords = append(allRecords, response.Ns...)
	return dnsr.MSG2RRParser(allRecords)
}

func (dnsr *DNSResponse) MSG2RRParser(records []dns.RR) error {
	for _, record := range records {
		if dnsr.TTL == 0 && record.Header().Ttl > 0 {
			dnsr.TTL = record.Header().Ttl
		}

		switch dnstype := record.(type) {

		case *dns.A:
			dnsr.A = append(dnsr.A, DOTtrimers(dnstype.A.String()))

		case *dns.AAAA:
			dnsr.AAAA = append(dnsr.AAAA, DOTtrimers(dnstype.AAAA.String()))

		case *dns.CNAME:
			dnsr.CNAME = append(dnsr.CNAME, DOTtrimers(dnstype.Target))

		case *dns.SOA:
			dnsr.SOA = append(dnsr.SOA, SOARecords{Name: DOTtrimers(dnstype.Hdr.Name),
				NS:      DOTtrimers(dnstype.Ns),
				Mbox:    DOTtrimers(dnstype.Mbox),
				Serial:  dnstype.Serial,
				Refresh: dnstype.Refresh,
				Retry:   dnstype.Retry,
				Expire:  dnstype.Expire,
				Minttl:  dnstype.Minttl})
		case *dns.NS:
			dnsr.NS = append(dnsr.NS, DOTtrimers(dnstype.Ns))
		case *dns.PTR:
			dnsr.PTR = append(dnsr.PTR, DOTtrimers(dnstype.Ptr))
		case *dns.MX:
			dnsr.MX = append(dnsr.MX, DOTtrimers(dnstype.Mx))
		case *dns.CAA:
			dnsr.CAA = append(dnsr.CAA, DOTtrimers(dnstype.Value))

		case *dns.TXT:
			dnsr.TXT = append(dnsr.TXT, strings.Join(dnstype.Txt, ""))

		case *dns.SRV:
			dnsr.SRV = append(dnsr.SRV, dnstype.Target)

		}
		dnsr.AllRecords = append(dnsr.AllRecords, record.String())
	}
	return nil
}

func (dnsr *DNSResponse) EnvelopeParser(envChan chan *dns.Envelope) error {
	var allRecords []dns.RR
	for env := range envChan {
		if env.Error != nil {
			return env.Error
		}
		allRecords = append(allRecords, env.RR...)
	}
	return dnsr.MSG2RRParser(allRecords)
}

func (dnsr *DNSResponse) SliceDeduplication() {
	dnsr.A = utils.Deduplslice(dnsr.A)
	dnsr.AAAA = utils.Deduplslice(dnsr.AAAA)
	dnsr.CNAME = utils.Deduplslice(dnsr.CNAME)
	dnsr.MX = utils.Deduplslice(dnsr.MX)
	dnsr.PTR = utils.Deduplslice(dnsr.PTR)
	dnsr.NS = utils.Deduplslice(dnsr.NS)
	dnsr.TXT = utils.Deduplslice(dnsr.TXT)
	dnsr.SRV = utils.Deduplslice(dnsr.SRV)
	dnsr.CAA = utils.Deduplslice(dnsr.CAA)
	dnsr.AllRecords = utils.Deduplslice(dnsr.AllRecords)
}

func (dnsr *DNSResponse) SoaParser() []string {
	var soarecs []string
	for _, soa := range dnsr.SOA {
		soarecs = append(soarecs, soa.NS, soa.Mbox)
	}
	return soarecs
}
