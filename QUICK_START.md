# Quick Start Guide

## Installation

### From Source
```bash
git clone https://github.com/yourusername/goscan.git
cd goscan
go build -o goscan .
```

### Using Install Script
```bash
./install.sh
```

## Basic Usage

### Scan Localhost
```bash
./goscan -host 127.0.0.1 -p 1-1024
```

### Scan Specific Ports
```bash
./goscan -host example.com -p 80,443,8080
```

### Fast Scan (More Workers)
```bash
./goscan -host example.com -p 1-65535 -w 100 -t 1s
```

## Configuration

### Generate Default Config
```bash
./goscan -generate-config
```

### Use Config File
```bash
./goscan -config goscan.yaml
```

## Output Formats

- JSON: `./goscan -o json -output results`
- CSV: `./goscan -o csv -output results`
- TXT: `./goscan -o txt -output results`

## Verbose Mode

```bash
./goscan -host example.com -p 1-100 -v
```