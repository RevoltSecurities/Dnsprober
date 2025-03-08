package dnsprober

import (
	"encoding/json"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/dns"
)

func (dnsprober *Dnsprober) WildcardDetect(hostname string, threshold int) (*WildcardResponse, error) {
	baseresponse, err := dnsprober.dnsclient.Query(hostname, "a")
	if err != nil {
		return nil, err
	}
	if len(baseresponse.A) == 0 {
		return &WildcardResponse{DNSResponse: *baseresponse, Wildcard: false}, nil
	}
	baseIPs := make(map[string]struct{})
	for _, ip := range baseresponse.A {
		baseIPs[ip] = struct{}{}
	}
	matchCount := 0
	for i := 0; i < threshold; i++ {
		response, err := dnsprober.dnsclient.Wildcard(hostname)
		if err != nil || len(response.A) == 0 {
			continue
		}
		matched := true
		for _, ip := range response.A {
			if _, ok := baseIPs[ip]; !ok {
				matched = false
				break
			}
		}

		if matched {
			matchCount++
		}
	}

	isWildcard := matchCount >= threshold

	return &WildcardResponse{
		DNSResponse: *baseresponse,
		Wildcard:    isWildcard,
	}, nil
}
func (dnsprober *Dnsprober) IsWildcard(baseresponse *dns.DNSResponse, hostname string, maxthreshold int) bool {

	baseIps := make(map[string]struct{})
	for _, ip := range baseresponse.A {
		baseIps[ip] = struct{}{}
	}

	matchCount := 0
	for i := 0; i < maxthreshold; i++ {
		response, err := dnsprober.dnsclient.Wildcard(hostname)
		if err != nil || len(response.A) == 0 {
			continue
		}

		matched := true
		for _, ip := range response.A {
			if _, ok := baseIps[ip]; !ok {
				matched = false
				break
			}
		}

		if matched {
			matchCount++
		}
	}

	// If a significant number of queries return the same base IPs, it's a wildcard domain
	return matchCount >= maxthreshold/2
}

func (response *WildcardResponse) JSONIZE() (string, error) {
	jsonized, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(jsonized), nil
}
