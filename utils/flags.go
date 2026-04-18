package utils

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Flags struct {
	Host         string
	NumWorkers   int
	PortRange    string
	Protocol     string
	Timeout      time.Duration
	OutputFormat string
	OutputFile   string
	ShowProgress bool
	Verbose      bool
}

func (f *Flags) InitializeFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `goscan Version 1.0.0
Usage: goscan <targets> [options]

TARGET SPECIFICATION:
  Valid targets are hostnames, IP addresses, networks, etc.
  Ex: scanme.nmap.org, 192.168.0.1, 192.168.0.0/24

OPTIONS:
`)
		flag.PrintDefaults()
	}

	flag.StringVar(&f.Host, "host", "127.0.0.1", "Target host to scan")
	flag.IntVar(&f.NumWorkers, "w", 25, "Number of concurrent workers")
	flag.StringVar(&f.PortRange, "p", "1-1024", "Port range (e.g., 80,443 or 1-1024)")
	flag.StringVar(&f.Protocol, "proto", "tcp", "Protocol: tcp or udp")
	flag.DurationVar(&f.Timeout, "t", 5*time.Second, "Connection timeout")
	flag.StringVar(&f.OutputFormat, "o", "json", "Output format: json, csv, txt")
	flag.StringVar(&f.OutputFile, "output", "scan_results", "Output file name")
	flag.BoolVar(&f.ShowProgress, "progress", true, "Show progress indicator")
	flag.BoolVar(&f.Verbose, "v", false, "Verbose mode (show closed ports)")

	flag.Parse()
}
