package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host         string        `yaml:"host"`
	Protocol     string        `yaml:"protocol"`
	PortRange    string        `yaml:"port_range"`
	Workers      int           `yaml:"workers"`
	Timeout      time.Duration `yaml:"timeout"`
	OutputFormat string        `yaml:"output_format"`
	OutputFile   string        `yaml:"output_file"`
	ShowProgress bool          `yaml:"show_progress"`
	Verbose      bool          `yaml:"verbose"`
	TUI          bool          `yaml:"tui"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Protocol == "" {
		cfg.Protocol = "tcp"
	}
	if cfg.PortRange == "" {
		cfg.PortRange = "1-1024"
	}
	if cfg.Workers == 0 {
		cfg.Workers = 25
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "json"
	}
	if cfg.OutputFile == "" {
		cfg.OutputFile = "scan_results"
	}

	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func Default() *Config {
	return &Config{
		Host:         "127.0.0.1",
		Protocol:     "tcp",
		PortRange:    "1-1024",
		Workers:      25,
		Timeout:      5 * time.Second,
		OutputFormat: "json",
		OutputFile:   "scan_results",
		ShowProgress: true,
		Verbose:      false,
		TUI:          false,
	}
}