<h1 align="center">
  <img src="static/dnsprober.png" alt="dnsx" height="500px" width="500px">
  <br>
</h1>

<h4 align="center">Dnsprober - A concurrent, lightweight, scalable, and efficient DNS reconnaissance too</h4>

<div align="center">
  
![GitHub last commit](https://img.shields.io/github/last-commit/RevoltSecurities/Dnsprober)  ![GitHub release (latest by date)](https://img.shields.io/github/v/release/RevoltSecurities/Dnsprober)  [![GitHub license](https://img.shields.io/github/license/RevoltSecurities/Dnsprober)](https://github.com/RevoltSecurities/Dnsprober/blob/main/LICENSE) 

</div>

<p align="center">
  <a href="https://github.com/RevoltSecurities/Dnsprober/edit/main/README.md#features">Features</a> |
  <a href="https://github.com/RevoltSecurities/Dnsprober/edit/main/README.md#installation">Installation</a> |
  <a href="#usage">Usage</a> |
  <a href="https://github.com/RevoltSecurities/Dnsprober/edit/main/README.md#using-dnsprober">Using Dnsprober</a> |
  <a href="https://github.com/RevoltSecurities/Dnsprober/edit/main/README.md#wildcard-filteration">Wildcard Filteration</a> |
  <a href="#reverse-dns-lookup">Reverse DNS</a> |
  <a href="#dns-probe">DNS Probe</a> |
  <a href="#info">Info</a>
</p>

---

**dnsprober** is a fast and multipurpose DNS reconnaissance tool designed for efficient DNS probing and enumeration. It supports multiple DNS queries, custom resolvers, wildcard filtering, and retryable DNS lookups for accurate DNS Recon and analysis. Built for scalability and speed, dnsprober ensures lightweight and high-performance reconnaissance with advanced filtering and optimization options.

---

## Features:

- **Fast & Scalable** – High-performance **concurrent** DNS probing and enumeration.
- **Multipurpose DNS Reconnaissance** – Supports multiple DNS record queries like `A`, `AAAA`, `CNAME`, `MX`, `NS`, `TXT`, `PTR`, `SOA`, `SRV`, `CAA`, `AXFR`, and more.
- **Retryable DNS Lookups** – Automatically retries failed requests to improve accuracy.
- **Wildcard Filtering** – Smart wildcard detection and filtering to reduce false positives.
- **Custom Resolvers Support** – Use your own DNS resolvers for better control and accuracy.
- **Optimized for Speed** – Supports adjustable concurrency, rate-limiting, and timeouts.
- **AXFR Zone Transfers** – Check for misconfigured DNS servers allowing zone transfers.
- **Colorized Output & Progress Bar** – Better visualization and tracking of progress.
- **DNS Response Filtering** – Filter results based on DNS response codes.


## Installation:

To install dnsprober, ensure you have Go **1.23.4** or later, then run the following command:

```bash
go install -v github.com/RevoltSecurities/Dnsprober/dnsprober@latest
```

or you can install and build the dnsprober binary using **git** and **go**

```bash
git clone https://github.com/RevoltSecurities/Dnsprober.git && cd Dnsprober/dnsprober
go build -o dnsprober
./dnsprober -h
```

## Usage:

```bash
dnsprober -h
```

```console
       __                                            __
  ____/ /   ____    _____    ____    _____  ____    / /_   ___    _____
 / __  /   / __ \  / ___/   / __ \  / ___/ / __ \  / __ \ / _ \  / ___/
/ /_/ /   / / / / (__  )   / /_/ / / /    / /_/ / / /_/ //  __/ / /
\__,_/   /_/ /_/ /____/   / .___/ /_/     \____/ /_.___/ \___/ /_/
                         /_/

                    - RevoltSecurities

[DESCRIPTION]:  Dnsprober - a concurrent, scalable and efficient DNS reconnaissance tool


[USAGE]:  

    dnsprober [flags]
        
[FLAGS]:  

    [INPUT]:  

        -d,  --domain                   :  Specify a single target domain for brute-forcing subdomains and also supports comma seperated values (ex: -d hackerone.com,bugcrowd.com)
        -l,  --list                     :  Provide a file containing a list of target domains, one per line.
        -w,  --wordlist                 :  Supply a wordlist file to brute-force subdomains (one word per line or comma-separated) (ex: -w word.txt or -w api,admin)
		
    [DNS QUERIES]:  

        --a                             :  Query A records (IPv4 addresses) for the target domain.
        --aaaa                          :  Query AAAA records (IPv6 addresses) for the target domain.
        --cname                         :  Query CNAME records to retrieve canonical names.
        --ns                            :  Query NS records to find authoritative name servers.
        --txt                           :  Query TXT records for text-based information.
        --srv                           :  Query SRV records for service information.
        --ptr                           :  Query PTR records for reverse DNS lookups.
        --mx                            :  Query MX records to discover mail exchange servers.
        --soa                           :  Query SOA records to retrieve Start of Authority information.
        --caa                           :  Query CAA records to see Certification Authority Authorization data.
        --any                           :  Query ANY records to attempt to retrieve all available DNS record types.
        --axfr                          :  Attempt a zone transfer (AXFR) from the target's authoritative DNS servers (if permitted).
        --all                           :  Query all supported record types for comprehensive enumeration.
        --resolve                       :  Simultaneously query both A and AAAA records to resolve the target host.
		
    [RATE LIMITS]:  

        -C,  --concurrency              :  Set the number of concurrent DNS queries (default: 10).
        -R,  --rate-limit               :  Limit the number of DNS queries per second (default: 0). Set to 0 for unlimited requests.
        
    [OUTPUT]:  

        -o,  --output                   :  Specify an output file to save the results.
        -j,  --json                     :  Output results in JSON format for easier parsing and integration with other tools.
        
    [CONFIGURATION]:  

        -r,  --resolvers                :  Provide a custom resolvers file (list of resolver IPs, comma-separated or one per line).
        -P,  --no-progress              :  Disable the progress bar display during execution.
        -D,  --wildcard-domain          :  Provide a wildcard-subdomain to filter duplicate dns subdomains
		
    [FILTERS]:  

        --response                      :  Display a summary of the DNS response along with the domain.
        --dns-response                  :  Display the DNS records data of the resolved domain.
        --dns-code                      :  Filter output by specific DNS response codes (e.g., noerror, refused).
        --raw-response                  :  Display the complete raw DNS response (full packet details).
		
    [OPTMIZATION]:  

        -t,  --timeout                  :  Set the DNS request timeout in seconds (default: 3).
        -E,  --retries                  :  Number of retry attempts for each request in case of failures (default: 3).
        --wildcard-threshold            :  Define the number of similar responses to consider a domain as having a wildcard (default: 5) (works only when enabled --wildcard-domain).
		
    [DEBUG]:  

        -v,  --verbose                  :  Enable verbose logging to show detailed debugging information.
        -s,  --silent                   :  Run in silent mode; suppress banner and version logging for cleaner output.
        --disable-update                :  Disable automatic update checks for the dnsprober.
        --no-color                      :  Disable colored output for run-time and outputs.
```

## Using Dnsprober:

### Resolving A Records:

```console
subdominator -d hackerone.com -s | dnsprober --response -s

a.ns.hackerone.com [A] [162.159.0.31]
api.hackerone.com [A] [104.18.36.214]
api.hackerone.com [A] [172.64.151.42]
b.ns.hackerone.com [A] [162.159.1.31]
docs.hackerone.com [A] [104.18.36.214]
docs.hackerone.com [A] [172.64.151.42]
gslink.hackerone.com [A] [3.165.75.21]
gslink.hackerone.com [A] [3.165.75.18]
gslink.hackerone.com [A] [3.165.75.26]
gslink.hackerone.com [A] [3.165.75.103]
mta-sts.managed.hackerone.com [A] [185.199.110.153]
mta-sts.managed.hackerone.com [A] [185.199.109.153]
mta-sts.managed.hackerone.com [A] [185.199.108.153]
mta-sts.managed.hackerone.com [A] [185.199.111.153]
mta-sts.forwarding.hackerone.com [A] [185.199.110.153]
mta-sts.forwarding.hackerone.com [A] [185.199.111.153]
mta-sts.forwarding.hackerone.com [A] [185.199.109.153]
mta-sts.forwarding.hackerone.com [A] [185.199.108.153]
mta-sts.hackerone.com [A] [185.199.109.153]
mta-sts.hackerone.com [A] [185.199.110.153]
mta-sts.hackerone.com [A] [185.199.108.153]
mta-sts.hackerone.com [A] [185.199.111.153]
zendesk1.hackerone.com [A] [216.198.54.2]
zendesk1.hackerone.com [A] [216.198.53.2]
www.hackerone.com [A] [104.18.36.214]
www.hackerone.com [A] [172.64.151.42]
zendesk3.hackerone.com [A] [216.198.53.2]
zendesk3.hackerone.com [A] [216.198.54.2]
resources.hackerone.com [A] [52.60.160.16]
resources.hackerone.com [A] [3.98.63.202]
zendesk2.hackerone.com [A] [216.198.54.2]
support.hackerone.com [A] [172.66.0.145]
resources.hackerone.com [A] [52.60.165.183]
zendesk2.hackerone.com [A] [216.198.53.2]
support.hackerone.com [A] [162.159.140.147]
zendesk4.hackerone.com [A] [216.198.53.2]
zendesk4.hackerone.com [A] [216.198.54.2]
```

### Extract the *A* Dns Records data:

```console
subdominator -d hackerone.com -s | dnsprober --dns-response -s

162.159.0.31
162.159.1.31
172.64.151.42
104.18.36.214
172.64.151.42
104.18.36.214
3.165.75.18
3.165.75.21
3.165.75.26
3.165.75.103
185.199.110.153
185.199.109.153
185.199.111.153
185.199.108.153
185.199.108.153
185.199.109.153
185.199.110.153
185.199.111.153
185.199.110.153
185.199.111.153
185.199.108.153
185.199.109.153
172.66.0.145
162.159.140.147
216.198.54.2
216.198.53.2
216.198.53.2
216.198.54.2
3.98.63.202
52.60.165.183
52.60.160.16
216.198.53.2
216.198.54.2
104.18.36.214
172.64.151.42
216.198.53.2
216.198.54.2
```

### Extract *CNAME* records for the given subdomains:

```console
subdominator -d hackerone.com -s | dnsprober --response -s --cname

fsdkim.hackerone.com [CNAME] [spfmx3.domainkey.freshemail.io]
fwdkim1.hackerone.com [CNAME] [spfmx1.domainkey.freshemail.io]
gslink.hackerone.com [CNAME] [d3rxkn2g2bbsjp.cloudfront.net]
mta-sts.forwarding.hackerone.com [CNAME] [hacker0x01.github.io]
mta-sts.managed.hackerone.com [CNAME] [hacker0x01.github.io]
mta-sts.hackerone.com [CNAME] [hacker0x01.github.io]
resources.hackerone.com [CNAME] [read.uberflip.com]
zendesk1.hackerone.com [CNAME] [mail1.zendesk.com]
zendesk2.hackerone.com [CNAME] [mail2.zendesk.com]
zendesk4.hackerone.com [CNAME] [mail4.zendesk.com]
zendesk3.hackerone.com [CNAME] [mail3.zendesk.com]
support.hackerone.com [CNAME] [2fe254e58a0ea8096400b2fda121ee35.freshdesk.com]
```

### Extract *ALL* DNS records for the given domain or subdomains:

```console
dnsprober -d x.com --all --response

x.com [A] [104.244.42.193]
x.com [A] [104.244.42.129]
x.com [NS] [a.r10.twtrdns.net]
x.com [NS] [a.u10.twtrdns.net]
x.com [NS] [b.r10.twtrdns.net]
x.com [NS] [b.u10.twtrdns.net]
x.com [NS] [c.r10.twtrdns.net]
x.com [NS] [c.u10.twtrdns.net]
x.com [NS] [d.r10.twtrdns.net]
x.com [NS] [d.u10.twtrdns.net]
x.com [TXT] [3089463]
x.com [TXT] [_w548xs1kfxtlqk3jyx19bzwk34c473i]
x.com [TXT] [kkdl3qb3tcrmdhfsm803p67r0my0svs8]
x.com [TXT] [apple-domain-verification=sEij6tJOW11fVNrG]
x.com [TXT] [adobe-sign-verification=c693a744ee2d282a36a43e6e724c5ea]
x.com [TXT] [shopify-verification-code=cUZazKrqCWgcshrcGvgfFR1lieuhRF]
x.com [TXT] [slack-domain-verification=Csk4bjCPFnJaDLLaKFUwCTFuUpCVvnYlAm2Tba0i]
x.com [TXT] [google-site-verification=8yQmoVhQedzlt36RPeQP41ytrEFk9aHEnde_xm0626g]
x.com [TXT] [google-site-verification=F6u9mGL--d2lbLljvH3b1UUgXtevQPdcamKr9c8914A]
x.com [TXT] [atlassian-sending-domain-verification=bd424180-8645-4de5-bd6a-285479c7577a]
x.com [TXT] [stripe-verification=46F7B88485621DC18923B43D12E90E6CDBCE232F2FEBCF084E6EFA91F6BA707D]
x.com [TXT] [adobe-idp-site-verification=ab4d9ce3473a73e81f46238da34ea4967fd5ac80e5c43fbfa8dff46d06a5321c]
x.com [TXT] [atlassian-domain-verification=j6u0o1PTkobCXC84uEF/sWpIPtaZURBVYqKzmTvT8wugLcHT1vvrzzA63iP1qSLN]
x.com [TXT] [figma-domain-verification=ee8420edd01965ba297f3438c907cfc6fbbaa1ee90a07b28f28bcfca8e6017bb-1729630998]
x.com [TXT] [v=spf1 ip4:199.16.156.0/22 ip4:199.59.148.0/22 include:_spf.google.com include:_spf.salesforce.com include:_oerp.x.com include:phx1.rp.oracleemaildelivery.com include:iad1.rp.oracleemaildelivery.com -all]
x.com [MX] [aspmx.l.google.com]
x.com [MX] [alt3.aspmx.l.google.com]
x.com [MX] [alt4.aspmx.l.google.com]
x.com [MX] [alt1.aspmx.l.google.com]
x.com [MX] [alt2.aspmx.l.google.com]
x.com [SOA] [a.u10.twtrdns.net]
x.com [SOA] [noc.twitter.com]
x.com [SOA] [a.u10.twtrdns.net]
x.com [SOA] [noc.twitter.com]
x.com [SOA] [a.u10.twtrdns.net]
x.com [SOA] [noc.twitter.com]
x.com [SOA] [a.u10.twtrdns.net]
x.com [SOA] [noc.twitter.com]
x.com [SOA] [a.u10.twtrdns.net]
x.com [SOA] [noc.twitter.com]
x.com [SOA] [a.u10.twtrdns.net]
x.com [SOA] [noc.twitter.com]
```

## Wildcard Filteration:

### Handling Wildcard DNS for Accurate Reconnaissance with dnsprober

During active DNS reconnaissance, a common challenge arises when non-existent subdomains resolve to the same IP address, causing excessive false positives and cluttered results. This can severely impact the accuracy of subdomain enumeration.

To tackle this, dnsprober includes a powerful wildcard filtering feature with the --wildcard-domain option. By specifying the wildcard domain name, dnsprober intelligently analyzes DNS responses, detects patterns, and filters out duplicate or misleading results. It uses a maximum threshold to determine if multiple subdomains resolve to the same IP address, confirming the presence of a wildcard domain.

With dnsprober’s wildcard detection, you can focus on valid subdomains without sifting through irrelevant or falsely resolved entries. This ensures precise DNS reconnaissance, eliminating noise and improving the quality of discovered subdomains efficiently. 🚀

```console
dnsprober -d freshworks.com -w wordlist.txt --wildcard-domain wildcard-domain.freshworks.com

feet.freshworks.com
espo.freshworks.com
gtest.freshworks.com
household.freshworks.com
larsen.freshworks.com
monk.freshworks.com
mvcc.freshworks.com
```

## Reverse DNS Lookup

dnsprober supports **Reverse DNS Lookup** using the **PTR** query, allowing you to retrieve domain names associated with IP addresses. This helps in identifying hostnames linked to a given IP, aiding in reconnaissance and infrastructure mapping.

```console
prips 157.240.19.0/24 | dnsprober --ptr --response

157.240.19.137 [PTR] [edgeray-avalanche-shv-01-dfw5.facebook.com]
157.240.19.10 [PTR] [edge-fblite-tcp-p1-shv-01-dfw5.facebook.com]
157.240.19.130 [PTR] [edgeray-biz-wa-shv-01-dfw5.facebook.com]
157.240.19.12 [PTR] [edge-mqtt-shv-01-dfw5.facebook.com]
157.240.19.140 [PTR] [edgeray-wa-genai-shv-01-dfw5.facebook.com]
157.240.19.14 [PTR] [edge-secure-shv-01-dfw5.facebook.com]
157.240.19.129 [PTR] [edgeray-cloudgaming-shv-01-dfw5.facebook.com]
157.240.19.141 [PTR] [edgeray-latest-star-shv-01-dfw5.facebook.com]
157.240.19.134 [PTR] [whatsapp-chatd-igd-edge-shv-01-dfw5.facebook.com]
157.240.19.138 [PTR] [edge-mws-l-shv-01-dfw5.facebook.com]
157.240.19.143 [PTR] [edge-ig-mqtt-shv-01-dfw5.facebook.com]
157.240.19.1 [PTR] [livestream-edgetee-upload-shv-01-dfw5.facebook.com]
157.240.19.135 [PTR] [edge-iglite-p3-shv-01-dfw5.facebook.com]
157.240.19.11 [PTR] [edge-msgr-latest-shv-01-dfw5.facebook.com]
157.240.19.139 [PTR] [edgeray-wa-group-shv-01-dfw5.facebook.com]
157.240.19.15 [PTR] [edge-dgw-shv-01-dfw5.facebook.com]
157.240.19.161 [PTR] [edge-ig-dgw-p42-shv-01-dfw5.facebook.com]
157.240.19.160 [PTR] [edge-iglite-p42-shv-01-dfw5.facebook.com]
157.240.19.16 [PTR] [edge-shortwave-shv-01-dfw5.facebook.com]
157.240.19.168 [PTR] [edge-dgw-p4-shv-01-dfw5.facebook.com]
157.240.19.142 [PTR] [edge-ig-dgw-shv-01-dfw5.facebook.com]
157.240.19.167 [PTR] [instagram-p4-new-shv-01-dfw5.fbcdn.net]
157.240.19.169 [PTR] [edge-z-p4-shv-01-dfw5.facebook.com]
157.240.19.17 [PTR] [anycast-control-ext-shv-01-dfw5.facebook.com]
157.240.19.170 [PTR] [edge-mqtt-p4-shv-01-dfw5.facebook.com]
157.240.19.174 [PTR] [instagram-p42-shv-01-dfw5.fbcdn.net]
157.240.19.193 [PTR] [metaenterprise-pp-shv-01-dfw5.facebook.com]
```
## DNS Probe:

dnsprober allows you to probe a list of (sub)domains based on their DNS response status codes. This feature helps in identifying live domains, misconfigured DNS records, and filtering out unwanted responses. By analyzing DNS status codes like NOERROR, NXDOMAIN, REFUSED, SERVFAIL, etc., you can refine your enumeration results, focusing only on meaningful responses.

This is particularly useful for detecting active subdomains, identifying potential takeovers, and filtering out wildcard responses that could clutter your reconnaissance data.

```console
subdominator -d facebook.com | dnsprober --dns-code noerror

edge-fwdproxy-4-bgp-01-dfw5.facebook.com [NOERROR]
edge-star-mini-shv-01-dfw5.facebook.com [NOERROR]
dohproxy-mobile-shv-01-dfw5.facebook.com [NOERROR]
anycast-control-mini-shv-01-dfw5.facebook.com [NOERROR]
edge-dgw-shv-01-dfw5.facebook.com [NOERROR]
edge-extern-shv-01-dfw5.facebook.com [NOERROR]
edge-shortwave-shv-01-dfw5.facebook.com [NOERROR]
edge-iglite-p3-shv-01-dfw5.facebook.com [NOERROR]
edge-secure-shv-01-dfw5.facebook.com [NOERROR]
edge-ig-mqtt-shv-01-dfw5.facebook.com [NOERROR]
edge-fwdproxy-1-bgp-01-dfw5.facebook.com [NOERROR]
edge-sonar-mini-shv-01-dfw5.fbcdn.net [NOERROR]
edge-stun-shv-01-dfw5.facebook.com [NOERROR]
edge-ig-dgw-p42-shv-01-dfw5.facebook.com [NOERROR]
edge-fwdproxy-5-bgp-01-dfw5.facebook.com [NOERROR]
edge-dgw-p4-shv-01-dfw5.facebook.com [NOERROR]
edge-fwdproxy-2-bgp-01-dfw5.facebook.com [NOERROR]
anycast-control-ext-shv-01-dfw5.facebook.com [NOERROR]
edge-star-shv-01-dfw5.facebook.com [NOERROR]
edge-recursor002-bgp-01-dfw5.facebook.com [NOERROR]
edge-onevc-sip-shv-01-dfw5.facebook.com [NOERROR]
edge-graph-cert-shv-01-dfw5.facebook.com [NOERROR]
edge-atlas-shv-01-dfw5.facebook.com [NOERROR]
edge-recursor001-bgp-01-dfw5.facebook.com [NOERROR]
dohproxy-shv-01-dfw5.facebook.com [NOERROR]
edge-fblite-tcp-mini-shv-01-dfw5.facebook.com [NOERROR]
edge-mqtt-p4-shv-01-dfw5.facebook.com [NOERROR]
edge-fblite-tcp-p1-shv-01-dfw5.facebook.com [NOERROR]
```
## INFO:

### **Built with ❤️ by RevoltSecurities**  
**dnsprober** is crafted with passion by **RevoltSecurities** for the open-source cybersecurity community. Designed for fast, scalable, and concurrent DNS reconnaissance, this tool empowers security researchers and penetration testers with accurate and efficient domain enumeration.  

Inspired by **[ProjectDiscovery's dnsx](https://github.com/projectdiscovery/dnsx)**, **dnsprober** serves as an **alternative with focused on scalability, wildcard filteration and concurrency**, ensuring high-performance DNS probing without compromising accuracy.  

We believe in the power of **open-source collaboration** and **continuous improvement**. Contributions, feedback, and innovative ideas are always welcome to help make **dnsprober** even better. Join us in shaping the future of cybersecurity tools! 🚀💙  
