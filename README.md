# GoScan - Professional Port Scanner

A production-ready TCP/UDP port scanner written in Go with service detection.

## Features

- TCP & UDP port scanning
- Service detection (20+ services)
- Multiple output formats (JSON, CSV, TXT)
- Configurable worker pools with rate limiting
- Structured logging (DEBUG, INFO, WARN, ERROR)
- YAML configuration support

## Installation

```bash
go build -o goscan .
```

## Usage

```bash
# Basic scan
./goscan -host 127.0.0.1 -p 1-1024

# Specific ports
./goscan -host example.com -p 80,443,8080

# Config file
./goscan -config config.yaml

# Generate config
./goscan -generate-config

# TUI mode
./goscan -tui
```

## Options

```
-host string       Target host to scan (default "127.0.0.1")
-w int           Number of workers (default 25)
-p string        Port range (default "1-1024")
-proto string    Protocol: tcp or udp (default "tcp")
-t duration      Timeout (default 5s)
-o string        Output format: json, csv, txt (default "json")
-output string   Output file name (default "scan_results")
-progress       Show progress (default true)
-v              Verbose mode
-tui             Launch TUI
-config string   Config file path
-generate-config  Generate sample config
```

## Configuration

Create `goscan.yaml`:

```yaml
host: 127.0.0.1
protocol: tcp
port_range: 1-1024
workers: 25
timeout: 5s
output_format: json
output_file: scan_results
verbose: false
```

## Testing

```bash
go test ./...
```

## Security

Only scan systems you own or have permission to test.

## License

MIT