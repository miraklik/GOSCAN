package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	content := `host: 192.168.1.1
protocol: tcp
port_range: 80,443
workers: 50
timeout: 3s
output_format: csv
output_file: results
show_progress: false
verbose: true
tui: false
`
	tmpFile, err := os.CreateTemp("", "goscan-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Host != "192.168.1.1" {
		t.Errorf("host = %q, want %q", cfg.Host, "192.168.1.1")
	}
	if cfg.Protocol != "tcp" {
		t.Errorf("protocol = %q, want %q", cfg.Protocol, "tcp")
	}
	if cfg.Workers != 50 {
		t.Errorf("workers = %d, want %d", cfg.Workers, 50)
	}
	if cfg.Verbose != true {
		t.Errorf("verbose = %v, want %v", cfg.Verbose, true)
	}
}

func TestLoad_Defaults(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "goscan-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Protocol != "tcp" {
		t.Errorf("protocol = %q, want %q (default)", cfg.Protocol, "tcp")
	}
	if cfg.Workers != 25 {
		t.Errorf("workers = %d, want %d (default)", cfg.Workers, 25)
	}
}

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Host != "127.0.0.1" {
		t.Errorf("host = %q, want %q", cfg.Host, "127.0.0.1")
	}
	if cfg.Protocol != "tcp" {
		t.Errorf("protocol = %q, want %q", cfg.Protocol, "tcp")
	}
	if cfg.Workers != 25 {
		t.Errorf("workers = %d, want %d", cfg.Workers, 25)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("timeout = %v, want %v", cfg.Timeout, 5*time.Second)
	}
}

func TestSave(t *testing.T) {
	cfg := Default()
	tmpFile, err := os.CreateTemp("", "goscan-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	if err := Save(tmpFile.Name(), cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	reloaded, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if reloaded.Host != cfg.Host {
		t.Errorf("reloaded host = %q, want %q", reloaded.Host, cfg.Host)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}