package utils

import (
	"io"
	"os"
	"strings"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/cli"
)

type DnsUtils struct {
	Options *cli.Argsparser
}

func (d *DnsUtils) GetQuery() string {
	opts := d.Options
	queries := []struct {
		flag  bool
		query string
	}{
		{opts.A, "a"},
		{opts.Aaaa, "aaaa"},
		{opts.All, "all"},
		{opts.Any, "any"},
		{opts.Caa, "caa"},
		{opts.Cname, "cname"},
		{opts.Mx, "mx"},
		{opts.Ns, "ns"},
		{opts.Ptr, "ptr"},
		{opts.Soa, "soa"},
		{opts.Srv, "srv"},
		{opts.Txt, "txt"},
		{opts.Resolve, "resolve"},
	}

	for _, item := range queries {
		if item.flag {
			return item.query
		}
	}
	opts.A = true
	return "a"
}

func (d *DnsUtils) GetInputSources() []io.Reader {
	var inputs []io.Reader

	if IsStdin() {
		inputs = append(inputs, os.Stdin)
	}

	if d.Options.List != "" {
		if file, err := os.Open(d.Options.List); err == nil {
			inputs = append(inputs, file)
		}
	}

	if d.Options.Domain != "" {
		domains := strings.Split(d.Options.Domain, ",")
		for _, domain := range domains {
			inputs = append(inputs, strings.NewReader(strings.TrimSpace(domain)+"\n"))
		}
	}
	return inputs
}
