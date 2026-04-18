# Usage Examples

## Basic Scans

### Single Host Scan
```bash
./goscan -host scanme.nmap.org -p 1-1024
```

### Multiple Ports
```bash
./goscan -host example.com -p 22,80,443,3306,5432
```

### Port Range
```bash
./goscan -host 192.168.1.1 -p 1-10000
```

## Advanced Options

### Fast Scan (High Workers)
```bash
./goscan -host example.com -p 1-65535 -w 200 -t 2s
```

### UDP Scan
```bash
./goscan -host example.com -proto udp -p 53,123,161
```

### Verbose Mode (Show Closed Ports)
```bash
./goscan -host example.com -p 1-100 -v
```

## Output Formats

### JSON Output
```bash
./goscan -host example.com -p 1-1024 -o json -output myscan
```

### CSV Output
```bash
./goscan -host example.com -p 1-1024 -o csv -output myscan
```

### TXT Output
```bash
./goscan -host example.com -p 1-1024 -o txt -output myscan
```

## Configuration File

### Using Config
```bash
./goscan -config myconfig.yaml
```

### Generate Config
```bash
./goscan -generate-config
```

## TUI Mode

```bash
./goscan -tui
```