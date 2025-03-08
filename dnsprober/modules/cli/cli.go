package cli

import (
	"fmt"
	"os"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/help"
	"github.com/spf13/cobra"
)

type Argsparser struct {
	Domain            string
	List              string
	WildcardDomain    string
	Wordlist          string
	Resolvers         string
	Output            string
	Dnscode           string
	A                 bool
	Aaaa              bool
	Cname             bool
	Ns                bool
	Txt               bool
	Srv               bool
	Ptr               bool
	Mx                bool
	Soa               bool
	Caa               bool
	Any               bool
	Axfr              bool
	Resolve           bool
	All               bool
	Json              bool
	Verbose           bool
	Rawresponse       bool
	Response          bool
	DnsResponse       bool
	DisableUp         bool
	Silent            bool
	NoColor           bool
	NoProgress		  bool
	Concurrency       int
	Timeout           int
	Ratelimit         int
	Retries           int
	WildcardThreads   int
	Exceptions        error
	WildcardThreshold int
}

var Opts = &Argsparser{}
var rootCmd = &cobra.Command{
	Use: "dnsprober",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func customHelp(cmd *cobra.Command, args []string) {
	helpMessage := help.Helper()
	fmt.Println(helpMessage)
	os.Exit(0)
}

func CLI() Argsparser {
	err := rootCmd.Execute()
	if err != nil {
		return Argsparser{Exceptions: err}
	}
	return *Opts
}

func init() {
	rootCmd.SetHelpFunc(customHelp)

	rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		customHelp(cmd, nil)
		return nil
	})

	rootCmd.Flags().StringVarP(&Opts.Domain, "domain", "d", "", "Single domain to bruteforce (e.g., example.com)")
	rootCmd.Flags().StringVarP(&Opts.List, "list", "l", "", "File containing a list of domains")
	rootCmd.Flags().StringVarP(&Opts.WildcardDomain, "wildcard-domain", "D", "", "Comma-separated list of domains")
	rootCmd.Flags().StringVarP(&Opts.Wordlist, "wordlist", "w", "", "Wordlist file for brute-force")
	rootCmd.Flags().StringVarP(&Opts.Resolvers, "resolvers", "r", "", "Custom resolvers file")
	rootCmd.Flags().StringVarP(&Opts.Output, "output", "o", "", "Output file for results")
	rootCmd.Flags().StringVarP(&Opts.Dnscode, "dns-code", "", "", "Filter output by dns response code")
	rootCmd.Flags().BoolVarP(&Opts.NoProgress, "no-progress", "P", false, "disable the progress bar of dnsprober")
	rootCmd.Flags().BoolVarP(&Opts.A, "a", "", false, "Query A records")
	rootCmd.Flags().BoolVarP(&Opts.Aaaa, "aaaa", "", false, "Query AAAA records")
	rootCmd.Flags().BoolVarP(&Opts.Cname, "cname", "", false, "Query CNAME records")
	rootCmd.Flags().BoolVarP(&Opts.Ns, "ns", "", false, "Query NS records")
	rootCmd.Flags().BoolVarP(&Opts.Txt, "txt", "", false, "Query TXT records")
	rootCmd.Flags().BoolVarP(&Opts.Srv, "srv", "", false, "Query SRV records")
	rootCmd.Flags().BoolVarP(&Opts.Ptr, "ptr", "", false, "Query PTR records")
	rootCmd.Flags().BoolVarP(&Opts.Mx, "mx", "", false, "Query MX records")
	rootCmd.Flags().BoolVarP(&Opts.Soa, "soa", "", false, "Query SOA records")
	rootCmd.Flags().BoolVarP(&Opts.Caa, "caa", "", false, "Query CAA records")
	rootCmd.Flags().BoolVarP(&Opts.Any, "any", "", false, "Query ANY records")
	rootCmd.Flags().BoolVarP(&Opts.Axfr, "axfr", "", false, "Query Axfr records")
	rootCmd.Flags().BoolVarP(&Opts.Resolve, "resolve", "", false, "Query for A,AAAA records")
	rootCmd.Flags().BoolVarP(&Opts.All, "all", "", false, "Query all record types")
	rootCmd.Flags().BoolVarP(&Opts.Json, "json", "j", false, "Output results in JSON format")
	rootCmd.Flags().BoolVarP(&Opts.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVarP(&Opts.NoColor, "no-color", "n", false, "Enable to decolorize the output")
	rootCmd.Flags().BoolVarP(&Opts.Rawresponse, "raw-response", "", false, "show raw dns query response")
	rootCmd.Flags().BoolVarP(&Opts.Response, "response", "", false, "show domain with dns response")
	rootCmd.Flags().BoolVarP(&Opts.DnsResponse, "dns-response", "", false, "show dns response record details")
	rootCmd.Flags().BoolVarP(&Opts.DisableUp, "disable-update", "", false, "Disable update check")
	rootCmd.Flags().BoolVarP(&Opts.Silent, "silent", "s", false, "Enable silent mode (no banner/version logging)")
	rootCmd.Flags().IntVarP(&Opts.Concurrency, "concurrency", "C", 50, "Set concurrency level")
	rootCmd.Flags().IntVarP(&Opts.Timeout, "timeout", "t", 3, "Request timeout (in seconds)")
	rootCmd.Flags().IntVarP(&Opts.Ratelimit, "rate-limit", "R", 0, "Rate limit for requests")
	rootCmd.Flags().IntVarP(&Opts.Retries, "retries", "E", 2, "Number of retries per request")
	rootCmd.Flags().IntVarP(&Opts.WildcardThreshold, "wildcard-threshold", "", 5, "Number of threshold for wildcard detection")
}
