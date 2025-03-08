package dnsprober

import (
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/dns"
)

type Dnsprober struct {
	dnsclient *dns.Dnsclient
}

type Options struct {
	Resolvers  []string
	Timeout    int
	Maxretries int
}

type WildcardResponse struct {
	dns.DNSResponse
	Wildcard bool `json:"wildcard"`
}

var KnownResolvers = []string{"1.1.1.1", "8.8.8.8", "8.8.4.4", "1.0.0.1", "9.9.9.9", "149.112.112.112", "208.67.222.222", "208.67.220.220", "185.228.168.9", "185.228.169.9"}

func New(options Options) *Dnsprober {
	newclient := dns.New(options.Resolvers, options.Maxretries, options.Timeout)
	return &Dnsprober{dnsclient: newclient}
}

func Default() *Dnsprober {
	newclient := dns.New(KnownResolvers, 2, 20)
	return &Dnsprober{dnsclient: newclient}
}

func (dnsprober *Dnsprober) ReverseDns(ip string) (*dns.DNSResponse, error) {
	return dnsprober.dnsclient.Query(ip, "ptr")
}

func (dnsprober *Dnsprober) Query(hostname string, query string) (*dns.DNSResponse, error) {
	return dnsprober.dnsclient.Query(hostname, query)
}

func (dnsprober *Dnsprober) QueryWithResolver(hostname string, query string, resolver string) (*dns.DNSResponse, error) {
	return dnsprober.dnsclient.Query(hostname, query, resolver)
}

func(dnsprober *Dnsprober) AxfrQuery(hostname string) (*dns.AXFRResponse, error){
	return dnsprober.dnsclient.AXFR(hostname)
}
