#!/bin/bash

###############################################################################
# Quick Deployment Script untuk Go SNMP OLT ZTE C320
# Dengan Support Firmware V2.1.0
###############################################################################

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

APP_DIR="/opt/go-snmp-olt"
SERVICE_NAME="go-snmp-olt"

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check root
if [ "$EUID" -ne 0 ]; then 
    print_error "Script ini harus dijalankan sebagai root"
    exit 1
fi

echo "=============================================="
echo "  Go SNMP OLT ZTE C320 - V2.1.0 Support"
echo "=============================================="
echo ""

# Menu
echo "Pilih opsi:"
echo "1) Fresh Install (dari source)"
echo "2) Update ke versi V2.1.0 support"
echo "3) Ubah konfigurasi firmware version"
echo "4) Test SNMP connectivity"
echo "5) View logs"
echo "6) Restart service"
echo ""
read -p "Pilihan [1-6]: " CHOICE

case $CHOICE in
    1)
        print_info "Menjalankan fresh install..."
        # Download dan jalankan installer utama
        if [ -f "./install.sh" ]; then
            bash ./install.sh
        else
            curl -fsSL https://raw.githubusercontent.com/your-repo/install.sh | bash
        fi
        ;;
        
    2)
        print_info "Update aplikasi dengan V2.1.0 support..."
        
        # Backup
        if [ -d "$APP_DIR" ]; then
            cp "$APP_DIR/.env" "$APP_DIR/.env.backup" 2>/dev/null || true
        fi
        
        # Stop service
        systemctl stop $SERVICE_NAME 2>/dev/null || true
        
        # Re-clone dan build
        cd /tmp
        rm -rf go-snmp-olt-zte-c320-new
        
        print_info "Cloning repository dengan V2.1.0 support..."
        git clone https://github.com/s4lfanet/go-api-c320.git go-snmp-olt-zte-c320-new
        
        # Copy modified files (jika ada)
        if [ -f "./oid_generator.go" ]; then
            cp ./oid_generator.go go-snmp-olt-zte-c320-new/config/
        fi
        
        cd go-snmp-olt-zte-c320-new
        
        # Build
        print_info "Building aplikasi..."
        /usr/local/go/bin/go mod download
        CGO_ENABLED=0 GOOS=linux /usr/local/go/bin/go build -ldflags="-s -w" -o bin/api ./cmd/api
        
        # Copy binary
        mkdir -p "$APP_DIR/bin"
        cp bin/api "$APP_DIR/bin/api"
        chmod +x "$APP_DIR/bin/api"
        
        # Copy new files
        cp .env.example "$APP_DIR/.env.example.new"
        cp FIRMWARE_V21_SUPPORT.md "$APP_DIR/" 2>/dev/null || true
        
        # Restore config
        if [ -f "$APP_DIR/.env.backup" ]; then
            cp "$APP_DIR/.env.backup" "$APP_DIR/.env"
            # Add firmware version if not exists
            if ! grep -q "ZTE_FIRMWARE_VERSION" "$APP_DIR/.env"; then
                echo "" >> "$APP_DIR/.env"
                echo "# ZTE C320 Firmware Version (v2.1 or v2.2)" >> "$APP_DIR/.env"
                echo "ZTE_FIRMWARE_VERSION=v2.1" >> "$APP_DIR/.env"
            fi
        fi
        
        # Restart
        systemctl start $SERVICE_NAME
        
        # Cleanup
        rm -rf /tmp/go-snmp-olt-zte-c320-new
        
        print_success "Update selesai! Service sudah di-restart."
        ;;
        
    3)
        print_info "Konfigurasi firmware version..."
        
        echo ""
        echo "Versi firmware saat ini:"
        grep "ZTE_FIRMWARE_VERSION" "$APP_DIR/.env" 2>/dev/null || echo "Belum dikonfigurasi"
        echo ""
        
        echo "Pilih versi firmware:"
        echo "  1) V2.1.x (untuk firmware lama)"
        echo "  2) V2.2.x atau lebih baru"
        read -p "Pilihan [1]: " FW_CHOICE
        
        if [ "$FW_CHOICE" == "2" ]; then
            NEW_VERSION="v2.2"
        else
            NEW_VERSION="v2.1"
        fi
        
        # Update .env
        if grep -q "ZTE_FIRMWARE_VERSION" "$APP_DIR/.env"; then
            sed -i "s/ZTE_FIRMWARE_VERSION=.*/ZTE_FIRMWARE_VERSION=$NEW_VERSION/" "$APP_DIR/.env"
        else
            echo "" >> "$APP_DIR/.env"
            echo "# ZTE C320 Firmware Version" >> "$APP_DIR/.env"
            echo "ZTE_FIRMWARE_VERSION=$NEW_VERSION" >> "$APP_DIR/.env"
        fi
        
        print_success "Firmware version diubah ke $NEW_VERSION"
        
        # Restart service
        systemctl restart $SERVICE_NAME
        
        # Clear cache
        if [ -f /root/.redis_password ]; then
            REDIS_PASS=$(cat /root/.redis_password)
            redis-cli -a "$REDIS_PASS" FLUSHALL 2>/dev/null || true
        else
            redis-cli FLUSHALL 2>/dev/null || true
        fi
        
        print_success "Service di-restart dan cache di-clear"
        ;;
        
    4)
        print_info "Test SNMP connectivity..."
        
        SNMP_HOST=$(grep "SNMP_HOST" "$APP_DIR/.env" | cut -d'=' -f2)
        SNMP_COMMUNITY=$(grep "SNMP_COMMUNITY" "$APP_DIR/.env" | cut -d'=' -f2)
        
        echo ""
        echo "Target: $SNMP_HOST"
        echo "Community: $SNMP_COMMUNITY"
        echo ""
        
        echo "=== Device Info ==="
        snmpwalk -v2c -c "$SNMP_COMMUNITY" "$SNMP_HOST" 1.3.6.1.2.1.1.1.0 2>&1 || echo "Failed"
        
        echo ""
        echo "=== GPON ONU List (sample) ==="
        snmpwalk -v2c -c "$SNMP_COMMUNITY" "$SNMP_HOST" 1.3.6.1.4.1.3902.1082.500.10.2.3.3.1.2 2>&1 | head -20 || echo "No data or failed"
        ;;
        
    5)
        print_info "Menampilkan logs..."
        journalctl -u $SERVICE_NAME -f
        ;;
        
    6)
        print_info "Restart service..."
        systemctl restart $SERVICE_NAME
        
        # Clear cache
        if [ -f /root/.redis_password ]; then
            REDIS_PASS=$(cat /root/.redis_password)
            redis-cli -a "$REDIS_PASS" FLUSHALL 2>/dev/null || true
        else
            redis-cli FLUSHALL 2>/dev/null || true
        fi
        
        print_success "Service di-restart dan cache di-clear"
        systemctl status $SERVICE_NAME
        ;;
        
    *)
        print_error "Pilihan tidak valid"
        exit 1
        ;;
esac
