# 📚 ZTE C320 OLT API - Documentation

Comprehensive documentation for ZTE C320 OLT Management API with SNMP & Telnet integration.

## 🚀 Getting Started

Perfect for first-time users and quick deployment.

- **[Installation Guide](getting-started/INSTALLATION.md)** - Complete installation walkthrough
  - Automated VPS installation (Ubuntu, Debian, CentOS, Rocky Linux)
  - Manual installation steps
  - Docker deployment
  
- **[Quick Start](getting-started/QUICK_START.md)** - Deploy in 5 minutes
  - One-command installation
  - Basic configuration
  - First API calls
  
- **[Proxmox Deployment](getting-started/PROXMOX_GUIDE.md)** - For Proxmox VE users
  - Container (CT) deployment
  - Virtual Machine (VM) deployment
  - Network configuration
  - Resource allocation

## 🏗️ Architecture

Understanding the system design and components.

- **[System Overview](architecture/OVERVIEW.md)** - High-level architecture
  - Component diagram
  - Technology stack
  - Data flow
  
- **[Workflow & Processes](architecture/WORKFLOW.md)** - How it works
  - SNMP data collection
  - Telnet command execution
  - Caching strategy
  - Error handling
  
- **[SNMP OID Mapping](architecture/OID_MAPPING.md)** - ZTE C320 V2.1.0 specifics
  - Available OIDs
  - Limitations (optical power)
  - Telnet fallback strategy

## ✨ Features

Explore available functionality and capabilities.

- **[Feature Phases & Roadmap](features/PHASES.md)** - Development timeline
  - Phase 1-7 (Completed)
  - Current capabilities
  - Future enhancements
  
- **[API Reference](features/API_REFERENCE.md)** - Complete endpoint documentation
  - 50+ REST endpoints
  - Request/response examples
  - Authentication & rate limiting
  
- **[Real-time Monitoring](features/MONITORING.md)** - Phase 7.1 & 7.2
  - ONU monitoring with optical power
  - PON port aggregation
  - OLT-wide statistics
  - Alert thresholds

## 🚢 Deployment

Production deployment guides for various environments.

- **[Docker Deployment](deployment/DOCKER.md)** - Containerized deployment
  - docker-compose setup
  - Environment variables
  - Multi-container architecture
  
- **[VPS Deployment](deployment/VPS.md)** - Native Linux deployment
  - Public VPS (DigitalOcean, AWS, etc.)
  - Local VPS (Proxmox, VMware)
  - Systemd service configuration
  - Nginx reverse proxy
  
- **[Troubleshooting](deployment/TROUBLESHOOTING.md)** - Common issues & solutions
  - Connection problems
  - Performance tuning
  - Debugging guide

## 👨‍💻 Development

For contributors and advanced users.

- **[Contributing Guide](development/CONTRIBUTING.md)** - How to contribute
  - Code standards
  - Pull request process
  - Testing requirements
  
- **[Command Reference](development/COMMAND_REFERENCE.md)** - Telnet commands
  - ZTE C320 CLI commands
  - Configuration syntax
  - Best practices
  
- **[Testing Guide](development/TESTING.md)** - Quality assurance
  - Unit tests
  - Integration tests
  - Load testing with k6

## 📖 Additional Resources

- **[Main README](../README.md)** - Project overview
- **[Changelog](../CHANGELOG.md)** - Version history
- **[Project State](../PROJECT_STATE.md)** - Current deployment status
- **[Security Policy](../SECURITY.md)** - Security guidelines
- **[License](../LICENSE)** - MIT License

## 🆘 Need Help?

- **Issues**: [GitHub Issues](https://github.com/s4lfanet/go-api-c320/issues)
- **Discussions**: [GitHub Discussions](https://github.com/s4lfanet/go-api-c320/discussions)
- **Email**: wardian370@gmail.com

---

**Version**: 1.7.2  
**Last Updated**: January 12, 2026  
**Status**: Production Ready ✅
