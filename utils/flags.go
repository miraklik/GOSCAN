package utils

import (
	"flag"
	"fmt"
	"os"
	"time"

	"goscan/pkg/config"
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
	TUI          bool
	ConfigFile   string
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

	var configFile string
	var genConfig bool
	flag.StringVar(&configFile, "config", "", "Config file path (YAML)")
	flag.BoolVar(&genConfig, "generate-config", false, "Generate sample config file")
	flag.StringVar(&f.Host, "host", "", "Target host to scan")
	flag.IntVar(&f.NumWorkers, "w", 0, "Number of concurrent workers")
	flag.StringVar(&f.PortRange, "p", "", "Port range (e.g., 80,443 or 1-1024)")
	flag.StringVar(&f.Protocol, "proto", "", "Protocol: tcp or udp")
	flag.DurationVar(&f.Timeout, "t", 0, "Connection timeout")
	flag.StringVar(&f.OutputFormat, "o", "", "Output format: json, csv, txt")
	flag.StringVar(&f.OutputFile, "output", "", "Output file name")
	flag.BoolVar(&f.ShowProgress, "progress", true, "Show progress indicator")
	flag.BoolVar(&f.Verbose, "v", false, "Verbose mode (show closed ports)")
	flag.BoolVar(&f.TUI, "tui", false, "Launch TUI interface")

	flag.Parse()

	if genConfig {
		cfg := config.Default()
		if err := config.Save("goscan.yaml", cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Config file created: goscan.yaml")
		os.Exit(0)
	}

	if configFile != "" {
		cfg, err := config.Load(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		f.applyConfig(cfg)
	}
}

func (f *Flags) applyConfig(cfg *config.Config) {
	if f.Host == "" {
		f.Host = cfg.Host
	}
	if f.NumWorkers == 0 {
		f.NumWorkers = cfg.Workers
	}
	if f.PortRange == "" {
		f.PortRange = cfg.PortRange
	}
	if f.Protocol == "" {
		f.Protocol = cfg.Protocol
	}
	if f.Timeout == 0 {
		f.Timeout = cfg.Timeout
	}
	if f.OutputFormat == "" {
		f.OutputFormat = cfg.OutputFormat
	}
	if f.OutputFile == "" {
		f.OutputFile = cfg.OutputFile
	}
	f.ShowProgress = cfg.ShowProgress
	f.Verbose = cfg.Verbose
	f.TUI = cfg.TUI
}
