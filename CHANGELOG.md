# Changelog

## [v2.0.1] - 2026-04-18

### Added
- Unit tests for scanner, detector, utils, config
- Structured logging with levels (DEBUG, INFO, WARN, ERROR)
- YAML config file support
- Recovery mechanism in main()
- Generate config flag (`-generate-config`)
- Documentation (QUICK_START.md, EXAMPLES.md, PROJECT_STRUCTURE.md, CONTRIBUTING.md)

### Fixed
- Service detection bug (SMTP vs FTP)
- Port validation

## [v2.0.0] - 2026-04-17

### Added
- Initial release
- TCP/UDP scanning
- Service detection (20+ services)
- Multiple output formats (JSON, CSV, TXT)
- TUI interface
- Worker pool with rate limiting