# Changelog

All notable changes to the ZTE C320 OLT API project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Automated VPS Installation Scripts**
  - Added `scripts/install.sh` - Full automated installer for Linux VPS (Ubuntu, Debian, CentOS, Rocky)
  - Added `scripts/install-quickstart.sh` - One-line installation command
  - Added `scripts/deploy-v21.sh` - Interactive deployment management script
  - Added `docs/INSTALLATION.md` - Comprehensive 600+ line installation guide
  - Auto-installs Go 1.25.5, Redis 7.2, and systemd service
  - Environment variable configuration wizard
  - Support for firmware version selection (v2.1/v2.2)

### Changed
- Updated README.md with automated installation section
- Updated repository URLs from old organization to s4lfanet

### Fixed
- **CI/CD golangci-lint Compatibility**
  - Fixed `routes_test.go` type mismatch for `trafficHandler` parameter
  - Changed from `*handler.TrafficHandler` to `handler.TrafficHandlerInterface` in all 7 test functions
  - Ensures strict type checking passes in GitHub Actions CI/CD
- **CI/CD Trivy Security Scanner**
  - Fixed hardcoded `:develop` tag in Trivy scan step
  - Implemented dynamic tag detection based on branch/tag
  - Main branch → scan `:latest`, develop branch → scan `:develop`, version tags → scan actual version
  - Resolves `MANIFEST_UNKNOWN` error when scanning non-existent image tags

## [1.7.2] - 2026-01-12

### Added
- **Phase 7.2: Optical Power Monitoring via Telnet**
  - Added `internal/repository/telnet_optical.go` - Optical power data retrieval via Telnet
  - Added `GetONUOpticalInfo()` - Retrieve optical info for specific ONU
  - Added `GetPONOpticalInfo()` - Retrieve optical info for all ONUs on PON port
  - Added `OpticalInfo` model in `internal/model/monitoring.go`
  - Integrated optical data with monitoring endpoints (GET /api/v1/monitoring/onu/{pon}/{onuId})
  - Environment file loader using godotenv for .env support
  - Optical power status classification (normal/low/high)

### Changed
- Updated `MonitoringUsecase` to inject `TelnetSessionManager` for optical data
- Updated `ONUMonitoringInfo` model to include `Optical` field
- Updated `cmd/api/main.go` to load .env file on startup
- Updated `.gitignore` to exclude build artifacts and diagnostic tools

### Fixed
- Fixed `routes_test.go` missing `monitoringHandler` parameter (11th param)
- Fixed environment variable loading for Redis authentication
- Fixed VPS project folder consistency (consolidated to `/opt/go-snmp-olt/`)

### Technical Notes
- **V2.1.0 SNMP Limitation**: Optical power OIDs NOT available via SNMP in firmware V2.1.0
- **Telnet Fallback**: Uses Telnet command `show gpon onu optical-info gpon-olt_1/{board}/{pon} {onu_id}`
- **Parser**: Regex-based parser for optical info (RX/TX power, temperature, voltage, bias current)
- **Deployment**: VPS cleaned up, single source of truth at `/opt/go-snmp-olt/`

## [1.7.1] - 2026-01-11

### Added
- **Phase 7.1: Real-time ONU Monitoring**
  - GET /api/v1/monitoring/onu/{pon}/{onuId} - Single ONU monitoring
  - GET /api/v1/monitoring/pon/{pon} - PON port aggregated monitoring
  - GET /api/v1/monitoring/olt - OLT-wide monitoring summary
  - `internal/usecase/monitoring.go` - Monitoring business logic
  - `internal/handler/monitoring.go` - Monitoring HTTP handlers
  - `internal/model/monitoring.go` - Monitoring data models

### Changed
- Updated `app/routes.go` to include monitoring routes

## [1.6.2] - 2025-12-31

### Added
- **Phase 6.2: Configuration Backup & Restore**
  - POST /api/v1/config/backup - Create configuration backup
  - POST /api/v1/config/restore - Restore from backup file
  - GET /api/v1/config/backups - List available backups
  - `internal/usecase/config_backup.go` - Backup/restore logic
  - `internal/handler/config_backup.go` - Backup HTTP handlers

### Changed
- Updated backup storage to `/opt/go-snmp-olt/backups/` directory

## [1.6.1] - 2025-12-30

### Added
- **Phase 6.1: Batch Operations**
  - POST /api/v1/batch/provision - Bulk ONU provisioning
  - POST /api/v1/batch/delete - Bulk ONU deletion
  - POST /api/v1/batch/reboot - Bulk ONU reboot
  - `internal/usecase/batch.go` - Batch operations logic
  - `internal/handler/batch.go` - Batch HTTP handlers

## [1.5.0] - 2025-12-30

### Added
- **Phase 5: ONU Lifecycle Management**
  - DELETE /api/v1/onu/{pon}/{onuId} - Delete ONU from PON
  - POST /api/v1/onu/{pon}/{onuId}/reboot - Reboot ONU
  - POST /api/v1/onu/{pon}/{onuId}/authorize - Authorize pending ONU
  - POST /api/v1/onu/{pon}/{onuId}/disable - Disable ONU
  - POST /api/v1/onu/{pon}/{onuId}/enable - Enable ONU
  - `internal/repository/telnet_onu_mgmt.go` - ONU management operations

## [1.4.0] - 2025-12-29

### Added
- **Phase 4: Traffic Profile Management**
  - POST /api/v1/traffic/{onuId}/assign - Assign traffic profile to ONU
  - DELETE /api/v1/traffic/{onuId}/remove - Remove traffic profile from ONU
  - `internal/repository/telnet_traffic.go` - Traffic profile operations

## [1.3.0] - 2025-12-29

### Added
- **Phase 3: VLAN Management**
  - POST /api/v1/vlan/create - Create VLAN configuration
  - PUT /api/v1/vlan/{vlanId} - Update VLAN configuration
  - DELETE /api/v1/vlan/{vlanId} - Delete VLAN configuration
  - `internal/repository/telnet_vlan.go` - VLAN configuration via Telnet

## [1.2.0] - 2025-12-28

### Added
- **Phase 2: ONU Provisioning via Telnet**
  - POST /api/v1/provision - Provision new ONU
  - `internal/repository/telnet_session.go` - Telnet session management
  - `internal/repository/telnet_provision.go` - ONU provisioning operations
  - Connection pooling for Telnet sessions

## [1.1.0] - 2025-12-27

### Added
- **Phase 1: Infrastructure & Basic SNMP**
  - GET /api/v1/onu - List all ONUs
  - GET /api/v1/onu/{pon} - List ONUs on specific PON
  - GET /api/v1/pon - List all PON ports
  - GET /api/v1/profiles/traffic - List traffic profiles
  - GET /api/v1/profiles/vlan - List VLAN profiles
  - SNMP repository with connection management
  - Redis caching layer
  - Configuration management via environment variables

## [1.0.0] - 2025-12-26

### Added
- Initial project setup
- Chi router with middleware
- Zerolog logging
- Health check endpoint
- Basic error handling
- Docker support
- CI/CD pipeline with GitHub Actions

---

## Legend

- **Added**: New features
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Removed features
- **Fixed**: Bug fixes
- **Security**: Security improvements
