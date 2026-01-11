#!/bin/bash

###############################################################################
# Quick Start Installer - Go SNMP OLT ZTE C320
# One-line installation script
###############################################################################

set -e

# Warna
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
cat << "EOF"
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║         Go SNMP OLT ZTE C320 - Quick Installer            ║
║                                                           ║
║         Repository: s4lfanet/go-api-c320                  ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Check root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root or with sudo"
    exit 1
fi

echo -e "${GREEN}Downloading full installer...${NC}"

# Download installer lengkap
curl -fsSL https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install.sh -o /tmp/go-snmp-installer.sh

# Jalankan installer
chmod +x /tmp/go-snmp-installer.sh
/tmp/go-snmp-installer.sh install

# Cleanup
rm -f /tmp/go-snmp-installer.sh

echo ""
echo -e "${GREEN}Installation completed!${NC}"
echo ""
echo "Quick commands:"
echo "  systemctl status go-snmp-olt    # Check status"
echo "  journalctl -u go-snmp-olt -f    # View logs"
echo "  curl http://localhost:8081/     # Test API"
echo ""
