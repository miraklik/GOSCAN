# Project Structure

```
goscan/
├── main.go                 # Entry point
├── scanner/
│   ├── scanner.go        # TCP/UDP scanner implementation
│   └── scanner_test.go
├── detector/
│   ├── detector.go    # Service detection by banner/port
│   └── detector_test.go
├── output/
│   └── output.go    # JSON/CSV/TXT export
├── utils/
│   ├── flags.go    # CLI flag handling
│   ├── worker.go  # Worker pool
│   ├── utils.go   # Helper functions
│   └── utils_test.go
├── tui/
│   └── tui.go     # Terminal UI
├── pkg/
│   ├── logger/
│   │   └── logger.go    # Structured logging
│   ├── config/
│   │   ├── config.go   # YAML config support
│   │   └── config_test.go
│   └── getsysinfo.go   # System info
├── go.mod           # Go module
├── go.sum           # Dependencies
└── README.md        # Documentation
```

## Components

### scanner/
TCP and UDP port scanning with configurable timeout.

### detector/
Service detection based on banner grabbing and known ports.

### output/
Export results to JSON, CSV, or TXT format.

### utils/
CLI flags, worker pool, and helper utilities.

### tui/
Interactive terminal UI using bubbletea.

### pkg/
Additional packages for logging, config, and system info.