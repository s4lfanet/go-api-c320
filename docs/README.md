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
  
- **[Firmware V2.1.0 Support](features/FIRMWARE_V21_SUPPORT.md)** - Firmware-specific implementation
  - SNMP limitations
  - Telnet workarounds
  - Optical power monitoring

## 📡 API Reference

Complete REST API documentation for frontend integration.

- **[API Reference](api/API_REFERENCE.md)** - Full endpoint documentation
  - 50+ REST endpoints with examples
  - Request/response formats
  - Error handling & status codes
  - Rate limiting & CORS
  - Authentication strategies
  - Code examples (cURL, JavaScript, Python)
  - Batch operations
  - Real-time monitoring endpoints

## 🚢 Deployment

Production deployment guides for various environments.

- **[Deployment Summary](deployment/DEPLOYMENT_SUMMARY.md)** - Overview of deployment options
- **[Phase 7.2 Deployment](deployment/PHASE_7.2_DEPLOYMENT.md)** - Latest deployment (Optical Power Monitoring)
  - Production VPS deployment
  - Configuration details
  - Testing results

## 👨‍💻 Development

For contributors and advanced users.

- **[Contributing Guide](development/CONTRIBUTING.md)** - How to contribute
  - Code standards
  - Pull request process
  - Testing requirements
  
- **[Command Reference](development/COMMAND_REFERENCE.md)** - ZTE C320 Telnet commands
  - CLI command syntax
  - Configuration examples
  - Best practices & tips

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
