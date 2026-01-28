# GoScan - Professional Port Scanner

🎯 **Complete, production-ready port scanner written in Go**

## 📦 What's Included

```
goscan/
├── 📂 Source Code (Go)
│   ├── cmd/goscan/        - Main application
│   ├── pkg/scanner/       - TCP/UDP scanning
│   ├── pkg/detector/      - Service detection
│   ├── pkg/output/        - JSON/CSV/TXT export
│   └── pkg/utils/         - Utilities
│
├── 📚 Documentation
│   ├── README.md          - Full documentation
│   ├── QUICK_START.md     - 5-minute setup guide
│   ├── EXAMPLES.md        - Usage examples
│   ├── PROJECT_STRUCTURE.md - Architecture
│   ├── CONTRIBUTING.md    - Developer guide
│   └── CHANGELOG.md       - Version history
│
├── 🛠️ Build Tools
│   ├── install.sh         - Auto installer
│   └── go.mod             - Dependencies
│
└── ⚙️ Configuration
    └── configs/config.example.yaml
```

## ✨ Features

- ✅ **TCP & UDP Scanning** - Full protocol support
- 🎯 **Flexible Port Ranges** - Single, multiple, ranges
- 🔍 **Service Detection** - Auto-identify 20+ services
- 📊 **Multiple Outputs** - JSON, CSV, TXT formats
- ⚡ **High Performance** - Configurable worker pools
- 🛡️ **Rate Limiting** - Prevent target overload
- 📈 **Statistics** - Detailed scan metrics
- 🎨 **Color Output** - Beautiful console display

## 🚀 Quick Start (30 seconds)

```bash
# Extract archive
tar -xzf goscan.tar.gz
cd goscan

# Install (automatic)
./install.sh

# First scan
./bin/goscan -host 127.0.0.1 -p 1-1024
```

## 💡 Usage Examples

### Basic Scan
```bash
./bin/goscan -host example.com -p 1-1024
```

### Web Server Scan
```bash
./bin/goscan -host example.com -p 80,443,8080,8443
```

### Fast Full Scan
```bash
./bin/goscan -host 192.168.1.1 -p 1-65535 -w 200 -t 2s
```

### UDP Services
```bash
./bin/goscan -host example.com -proto udp -p 53,123,161
```

### CSV Export
```bash
./bin/goscan -host example.com -p 1-1024 -o csv
```

## 📊 Sample Output

```
[*] Starting tcp scan on example.com
[*] Port range: 1-1024 (1024 ports)
[*] Workers: 25 | Timeout: 5s
============================================================
[+] Port    22 | SSH          | SSH-2.0-OpenSSH_8.2p1
[+] Port    80 | HTTP         | nginx/1.18.0
[+] Port   443 | HTTPS        | No banner
============================================================
[+] Found 3 open port(s)
[+] Results saved to: scan_results.json
============================================================
Scan Statistics:
  Total Ports:    1024
  Open Ports:     3 (0.29%)
  Closed Ports:   1021
  Duration:       15.2s
  Ports/second:   67.37
============================================================
```

## 🏗️ Architecture Highlights

### Concurrency Model
- Worker pool pattern for controlled concurrency
- Channel-based communication (no shared memory)
- Rate limiting to prevent network overload
- WaitGroups for synchronization

### Design Patterns
- Interface pattern for extensibility
- Factory pattern for clean initialization
- Channel-based communication (Go idiomatic)

## 🛠️ Building & Development

### Quick Build
```bash
make build          # Build binary
make run            # Build and run
make test           # Run tests
```

### Multi-Platform Build
```bash
make build-all      # Linux, macOS, Windows
```

### Development
```bash
make fmt            # Format code
make lint           # Run linter
make clean          # Clean artifacts
```

## 📋 Requirements

- **Go 1.24+** (for building)
- **Linux/macOS/Windows** (cross-platform)
- **Network access** (for scanning)

## 🔒 Security & Legal

⚠️ **IMPORTANT**: Only scan systems you own or have explicit permission to test.

Unauthorized port scanning may be:
- Illegal in your jurisdiction
- Violation of terms of service
- Triggering of security alerts
- Cause for legal action

**Use responsibly!**

## 📝 License

MIT License - See [LICENSE](goscan/LICENSE) file

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](goscan/CONTRIBUTING.md)

## 📞 Support

- **Issues**: Report bugs on GitHub
- **Questions**: Check documentation first
- **Features**: Submit feature requests

## 🌟 What Makes This Special?

1. **Production-Ready Code**
   - Clean architecture
   - Comprehensive error handling
   - Extensive testing
   - Professional documentation

2. **Educational Value**
   - Well-commented code
   - Clear design patterns
   - Architecture documentation
   - Learning examples

3. **Extensible Design**
   - Easy to add new scanners
   - Pluggable service detection
   - Multiple output formats
   - Configuration support

4. **Performance Optimized**
   - Worker pool pattern
   - Rate limiting
   - Buffered channels
   - Efficient resource usage

## 📈 Version History

**v2.0.0** (Current)
- Initial release
- TCP/UDP scanning
- Service detection
- Multiple output formats
- Complete documentation

See [CHANGELOG.md](goscan/CHANGELOG.md) for details

## 🚀 Next Steps

1. **Extract**: `tar -xzf goscan.tar.gz`
2. **Read**: Open `goscan/QUICK_START.md`
3. **Build**: Run `./install.sh`
4. **Scan**: Try your first scan!
5. **Learn**: Explore documentation
6. **Extend**: Add your own features

---

**Happy Scanning! 🎯**

*Built with ❤️ using Go*