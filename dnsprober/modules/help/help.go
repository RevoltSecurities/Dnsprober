package help

import (
	"fmt"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/banner"
	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/logger"
)

var loggers = logger.New(true)

func Helper() string {
	msg := fmt.Sprintf(`
%s
                    %s

%s

%s
    %s
%s
    %s
        %s
    %s
        %s
    %s
        %s
    %s
        %s
    %s
        %s
    %s
        %s
    %s
        %s
    %s
        %s
`,
		banner.BannerGenerator("dnsprober"),
		loggers.Bolder("- RevoltSecurities"),
		loggers.Loader("Dnsprober - a concurrent, scalable and efficient DNS reconnaissance tool", "DESCRIPTION"),
		loggers.Loader("", "USAGE"),
		loggers.Bolder(`dnsprober [flags]
        `),
		loggers.Loader("", "FLAGS"),
		loggers.Loader("", "INPUT"),
		loggers.Bolder(`-d,  --domain                   :  Specify a single target domain for brute-forcing subdomains and also supports comma seperated values (ex: -d hackerone.com,bugcrowd.com)
        -l,  --list                     :  Provide a file containing a list of target domains, one per line.
        -w,  --wordlist                 :  Supply a wordlist file to brute-force subdomains (one word per line or comma-separated) (ex: -w word.txt or -w api,admin)
		`),
		loggers.Loader("", "DNS QUERIES"),
		loggers.Bolder(`--a                             :  Query A records (IPv4 addresses) for the target domain.
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
		`),
		loggers.Loader("", "RATE LIMITS"),
		loggers.Bolder(`-C,  --concurrency              :  Set the number of concurrent DNS queries (default: 10).
        -R,  --rate-limit               :  Limit the number of DNS queries per second (default: 0). Set to 0 for unlimited requests.
        `),
		loggers.Loader("", "OUTPUT"),
		loggers.Bolder(`-o,  --output                   :  Specify an output file to save the results.
        -j,  --json                     :  Output results in JSON format for easier parsing and integration with other tools.
        `),
		loggers.Loader("", "CONFIGURATION"),
		loggers.Bolder(`-r,  --resolvers                :  Provide a custom resolvers file (list of resolver IPs, comma-separated or one per line).
        -P,  --no-progress              :  Disable the progress bar display during execution.
        -D,  --wildcard-domain          :  Provide a wildcard-subdomain to filter duplicate dns subdomains
		`),
		loggers.Loader("", "FILTERS"),
		loggers.Bolder(`--response                      :  Display a summary of the DNS response along with the domain.
        --dns-response                  :  Display the DNS records data of the resolved domain.
        --dns-code                      :  Filter output by specific DNS response codes (e.g., noerror, refused).
        --raw-response                  :  Display the complete raw DNS response (full packet details).
		`),
		loggers.Loader("", "OPTMIZATION"),
		loggers.Bolder(`-t,  --timeout                  :  Set the DNS request timeout in seconds (default: 3).
        -E,  --retries                  :  Number of retry attempts for each request in case of failures (default: 3).
        --wildcard-threshold            :  Define the number of similar responses to consider a domain as having a wildcard (default: 5) (works only when enabled --wildcard-domain).
		`),
		loggers.Loader("", "DEBUG"),
		loggers.Bolder(`-v,  --verbose                  :  Enable verbose logging to show detailed debugging information.
        -s,  --silent                   :  Run in silent mode; suppress banner and version logging for cleaner output.
        --disable-update                :  Disable automatic update checks for the dnsprober.
        --no-color                      :  Disable colored output for run-time and outputs.
        `),
	)
	return msg
}
