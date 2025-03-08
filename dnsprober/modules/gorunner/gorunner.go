package gorunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/cli"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/dns"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/dnsprober"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/logger"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/progressbar"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/reader"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/utils"
	"github.com/projectdiscovery/hmap/store/hybrid"
	"github.com/projectdiscovery/ratelimit"
	sliceutil "github.com/projectdiscovery/utils/slice"
)

type Gorunner struct {
	Dnsclient        *dnsprober.Dnsprober
	Args             *cli.Argsparser
	Outputchan       chan string
	SubdomainChan    chan string
	DnsJobChan       chan string
	ResultsChan      chan string
	Domainchan       chan string
	Resolverwg       *sync.WaitGroup
	Readerwg         *sync.WaitGroup
	Resultswg        *sync.WaitGroup
	Subdomainswg     *sync.WaitGroup
	Ratelimiter      *ratelimit.Limiter
	Hmap             *hybrid.HybridMap
	TotalJobs        int64
	Dnsrequest       string
	Tmpfile          string
	Dnscodes         []string
	Progress         *progressbar.ProgressBar
	WildcardResponse *dns.DNSResponse
	Logger           *logger.Logger
}

func New(args *cli.Argsparser) (*Gorunner, error) {
	var (
		resolvers   []string
		recordtype  string
		ratelimiter *ratelimit.Limiter
		dnscodes    []string
	)

	if args.Resolvers != "" {
		if utils.FileExists(args.Resolvers) {
			res, err := reader.Reader(args.Resolvers)
			if err != nil {
				return nil, err
			}
			resolvers = append(resolvers, res...)
		} else {
			res := utils.SplitStrings(args.Resolvers)
			resolvers = append(resolvers, res...)
		}
	} else {
		resolvers = append(resolvers, dnsprober.KnownResolvers...)
	}

	if args.Dnscode != "" {
		dnscodes = append(dnscodes, strings.Split(args.Dnscode, ",")...)
	}

	dnsuts := utils.DnsUtils{Options: args}
	recordtype = dnsuts.GetQuery()
	inputsources := dnsuts.GetInputSources()
	if len(inputsources) == 0 {
		return nil, fmt.Errorf("no input provided for dnsprober")
	}
	ratelimiter = ratelimit.NewUnlimited(context.Background())
	if args.Ratelimit != 0 {
		ratelimiter = ratelimit.New(context.Background(), uint(args.Ratelimit), time.Second)
	}
	hmap, err := hybrid.New(hybrid.DefaultDiskOptions)
	if err != nil {
		return nil, err
	}
	dnsproberclient := dnsprober.New(dnsprober.Options{
		Resolvers:  resolvers,
		Timeout:    args.Timeout,
		Maxretries: args.Retries,
	})

	logobj := logger.New(!args.NoColor)
	newrunner := &Gorunner{
		Args:          args,
		Hmap:          hmap,
		Ratelimiter:   ratelimiter,
		Dnsrequest:    recordtype,
		Dnsclient:     dnsproberclient,
		Resolverwg:    &sync.WaitGroup{},
		Resultswg:     &sync.WaitGroup{},
		Readerwg:      &sync.WaitGroup{},
		Subdomainswg:  &sync.WaitGroup{},
		DnsJobChan:    make(chan string),
		ResultsChan:   make(chan string),
		SubdomainChan: make(chan string),
		Outputchan:    make(chan string),
		TotalJobs:     0,
		Dnscodes:      dnscodes,
		Logger:        logobj,
	}
	return newrunner, nil
}

func (gorunner *Gorunner) ListIntoHmap(hosts []string) {
	for _, host := range hosts {
		if _, ok := gorunner.Hmap.Get(host); ok {
			continue
		}
		gorunner.Hmap.Set(host, nil)
		atomic.AddInt64(&gorunner.TotalJobs, 1)
	}
}

func (gorunner *Gorunner) HmapScanner() {
	gorunner.Hmap.Scan(func(jobs, _ []byte) error {
		job := string(jobs)
		gorunner.DnsJobChan <- job
		return nil
	})
	close(gorunner.DnsJobChan)
}

func (gorunner *Gorunner) SetupJobs() error {
	var err error
	if utils.IsStdin() {
		gorunner.Tmpfile, err = utils.CreateTmpFile()
		if err != nil {
			return err
		}
		file, err := os.Create(gorunner.Tmpfile)
		if err != nil {
			return err
		}
		if _, err := io.Copy(file, os.Stdin); err != nil {
			return err
		}
		defer file.Close()
		defer os.Remove(gorunner.Tmpfile)
	}

	if gorunner.Domainchan == nil {
		if utils.IsStdin() {
			gorunner.Domainchan, err = reader.Streamer(gorunner.Tmpfile)
			if err != nil {
				return err
			}
		} else if gorunner.Args.List != "" {
			gorunner.Domainchan, err = reader.Streamer(gorunner.Args.List)
			if err != nil {
				return err
			}
		} else if gorunner.Args.Domain != "" {
			domains := strings.ReplaceAll(gorunner.Args.Domain, ",", "\n")
			gorunner.Domainchan, err = reader.IOStreamer(strings.NewReader(strings.TrimSpace(domains) + "\n"))
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("no inputs provided for dnsprober")
		}
	}

	for i := 0; i < gorunner.Args.Concurrency; i++ {
		gorunner.Subdomainswg.Add(1)
		go func(workerID int) {
			defer gorunner.Subdomainswg.Done()
			for subdomain := range gorunner.SubdomainChan {
				gorunner.ListIntoHmap([]string{subdomain})
			}
		}(i)
	}

	for domain := range gorunner.Domainchan {
		if gorunner.Args.Wordlist != "" {
			if utils.FileExists(gorunner.Args.Wordlist) {
				wordCh, err := reader.Streamer(gorunner.Args.Wordlist)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error streaming wordlist: %v\n", err)
				}
				for sub := range wordCh {
					var subdomain string
					if strings.Contains(domain, "FUZZ") {
						subdomain = strings.ReplaceAll(domain, "FUZZ", sub)
					} else {
						subdomain = fmt.Sprintf("%s.%s", sub, domain)
					}
					gorunner.SubdomainChan <- subdomain
				}
			} else {
				words := strings.Split(gorunner.Args.Wordlist, ",")
				for _, sub := range words {
					sub = strings.TrimSpace(sub)
					if sub == "" {
						continue
					}
					var subdomain string
					if strings.Contains(domain, "FUZZ") {
						subdomain = strings.ReplaceAll(domain, "FUZZ", sub)
					} else {
						subdomain = fmt.Sprintf("%s.%s", sub, domain)
					}
					gorunner.SubdomainChan <- subdomain
				}
			}
		} else {
			gorunner.SubdomainChan <- domain
		}
	}
	close(gorunner.SubdomainChan)
	gorunner.Subdomainswg.Wait()
	return nil
}

func (gorunner *Gorunner) ResolverPool() {
	defer gorunner.Resolverwg.Done()
	for domain := range gorunner.DnsJobChan {
		gorunner.Ratelimiter.Take()
		resp, err := gorunner.Dnsclient.Query(domain, gorunner.Dnsrequest)
		if err != nil {
			if gorunner.Args.Verbose && !gorunner.Args.Silent {
				gorunner.Logger.Logger(fmt.Sprintf("unable to make DNS query for %s due to: %s", domain, err.Error()), "error")
			}
			gorunner.Progress.Increment(1, 1)
			continue
		} else {
			gorunner.Progress.Increment(1, 0)
		}

		if gorunner.Args.WildcardDomain != "" {
			if gorunner.Dnsclient.IsWildcard(resp, domain, gorunner.Args.WildcardThreshold) {
				continue
			}
		}

		if gorunner.Args.Axfr {
			axfresp, err := gorunner.Dnsclient.AxfrQuery(domain)

			if axfresp != nil && err == nil {
				resp.AXFR = axfresp
			}

		}

		if gorunner.Args.Json {
			jsonized, _ := resp.JSONIZE()
			gorunner.ResultsChan <- jsonized
			continue
		}

		if gorunner.Args.Dnscode != "" {
			gorunner.ResponseCodePaser(resp, domain)
			continue
		}

		if gorunner.Args.Rawresponse {
			gorunner.ResultsChan <- resp.Raw
			continue
		}

		if gorunner.Args.A {
			gorunner.ResponseParser(domain, resp.A, "A", resp)
		}

		if gorunner.Args.Aaaa {
			gorunner.ResponseParser(domain, resp.AAAA, "AAAA", resp)
		}

		if gorunner.Args.Caa {
			gorunner.ResponseParser(domain, resp.CAA, "CAA", resp)
		}

		if gorunner.Args.Cname {
			gorunner.ResponseParser(domain, resp.CNAME, "CNAME", resp)
		}

		if gorunner.Args.Ns {
			gorunner.ResponseParser(domain, resp.NS, "NS", resp)
		}

		if gorunner.Args.Txt {
			gorunner.ResponseParser(domain, resp.TXT, "TXT", resp)
		}

		if gorunner.Args.Srv {
			gorunner.ResponseParser(domain, resp.SRV, "SRV", resp)

		}

		if gorunner.Args.Ptr {
			gorunner.ResponseParser(domain, resp.PTR, "PTR", resp)
		}

		if gorunner.Args.Mx {
			gorunner.ResponseParser(domain, resp.MX, "MX", resp)
		}

		if gorunner.Args.Soa {
			gorunner.ResponseParser(domain, resp.SOA, "SOA", resp)
		}

		if gorunner.Args.Resolve {
			resolverec := sliceutil.Merge(resp.A, resp.AAAA)
			gorunner.ResponseParser(domain, resolverec, "RESOLVE", resp)
		}

		if gorunner.Args.Any {
			soarecords := sliceutil.Merge(resp.A, resp.AAAA, resp.CNAME, resp.NS, resp.TXT, resp.SRV, resp.PTR, resp.MX, resp.SoaParser(), resp.CAA)
			go gorunner.ResponseParser(domain, soarecords, "ANY", resp)
		}
	}
}

func (gorunner *Gorunner) ResponseParser(domain string, dnsresponsed interface{}, dnsrequest string, resp *dns.DNSResponse) {
	var dnsrecords []string
	switch dnsresponsed := dnsresponsed.(type) {
	case []string:
		if len(dnsresponsed) == 0 {
			return
		}
		dnsrecords = dnsresponsed
	case []dns.SOARecords:
		for _, dnsreponse := range dnsresponsed {
			dnsrecords = append(dnsrecords, dnsreponse.NS, dnsreponse.Mbox)
		}
	}

	for _, record := range dnsrecords {
		if gorunner.Args.Response {
			gorunner.ResultsChan <- fmt.Sprintf("%s [%s] [%s]", gorunner.Logger.Bolder(domain), gorunner.Logger.Colorizer(strings.ToUpper(dnsrequest), "blue"), gorunner.Logger.Colorizer(record, "yellow"))
		} else if gorunner.Args.DnsResponse {
			gorunner.ResultsChan <- gorunner.Logger.Bolder(record)
		} else {
			gorunner.ResultsChan <- fmt.Sprint(gorunner.Logger.Bolder(domain))
			break
		}
	}
}

func (gorunner *Gorunner) ResponseCodePaser(response *dns.DNSResponse, domain string) {
	if gorunner.IsStatusCode(response) {
		gorunner.ResultsChan <- fmt.Sprintf("%s [%s]", gorunner.Logger.Bolder(domain), gorunner.Logger.Colorizer(response.StatusCode, "green"))
	}
}

func (gorunner *Gorunner) IsStatusCode(resp *dns.DNSResponse) bool {
	for _, dnscode := range gorunner.Dnscodes {
		if strings.ToUpper(dnscode) == resp.StatusCode {
			return true
		}
	}
	return false
}

func (gorunner *Gorunner) InitateResolvers() {
	for i := 1; i < gorunner.Args.Concurrency; i++ {
		gorunner.Resolverwg.Add(1)
		go gorunner.ResolverPool()
	}
}

func (gorunner *Gorunner) SetOptions() {
	gorunner.Args.A = true
	gorunner.Args.Aaaa = true
	gorunner.Args.Cname = true
	gorunner.Args.Ns = true
	gorunner.Args.Txt = true
	gorunner.Args.Srv = true
	gorunner.Args.Ptr = true
	gorunner.Args.Mx = true
	gorunner.Args.Soa = true
	gorunner.Args.Caa = true
	gorunner.Args.Axfr = true
	gorunner.Args.Any = false
	gorunner.Args.DnsResponse = false
	gorunner.Args.Response = true
}

func (gorunner *Gorunner) Save(content string) {
	if gorunner.Args.Output == "" {
		return
	}
	file, err := os.OpenFile(gorunner.Args.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if gorunner.Args.Verbose && !gorunner.Args.Silent {
			gorunner.Logger.StdLogger(fmt.Sprintf("Exception occured in the save module method due to: %s", err.Error()), "warn")
		}
		return
	}
	defer file.Close()
	if _, err = file.WriteString(content + "\n"); err != nil {
		if gorunner.Args.Verbose && !gorunner.Args.Silent {
			gorunner.Logger.StdLogger(fmt.Sprintf("Exception occured in the save method when writing output due to: %s", err.Error()), "warn")
		}
		return
	}
}

func (gorunner *Gorunner) Sprint() error {
	if gorunner.Args.All {
		gorunner.SetOptions()
	}
	if err := gorunner.SetupJobs(); err != nil {
		return err
	}
	go gorunner.HmapScanner()

	if gorunner.Args.WildcardDomain != "" {
		wildcardResp, err := gorunner.Dnsclient.Query(gorunner.Args.WildcardDomain, gorunner.Dnsrequest)
		if err != nil {
			gorunner.Logger.Logger(fmt.Sprintf("Wildcard query failed for %s: %s", gorunner.Args.WildcardDomain, err.Error()), "error")
		} else {
			gorunner.WildcardResponse = wildcardResp
		}
	}

	gorunner.Progress = progressbar.New(atomic.LoadInt64(&gorunner.TotalJobs))

	if !gorunner.Args.NoProgress {
		go func() {
			for {
				gorunner.Progress.Render()
			}
		}()
	}
	go func() {
		for res := range gorunner.ResultsChan {
			if !gorunner.Args.NoProgress {
				gorunner.Logger.StdinLogger(res)
				gorunner.Outputchan <- res
			} else {
				fmt.Println(res)
			}
		}
	}()

	go func() {
		for result := range gorunner.Outputchan {
			if gorunner.Args.Output != "" {
				gorunner.Save(result)
			}
		}
	}()
	gorunner.InitateResolvers()
	gorunner.Resolverwg.Wait()
	if !gorunner.Args.NoProgress {
		gorunner.Progress.Render()
	}
	return nil
}

func (gorunner *Gorunner) Down() {
	gorunner.Hmap.Close()
}
