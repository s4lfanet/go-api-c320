#!/bin/bash

###############################################################################
# Installer Lengkap untuk Go SNMP OLT ZTE C320
# Untuk VPS Linux (Tanpa Docker)
# Repository: https://github.com/s4lfanet/go-api-c320
###############################################################################

set -e  # Exit on error

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variabel konfigurasi
APP_NAME="go-snmp-olt-zte-c320"
APP_DIR="/opt/go-snmp-olt"
APP_USER="olt-service"
SERVICE_NAME="go-snmp-olt"
REPO_URL="https://github.com/s4lfanet/go-api-c320.git"
GO_VERSION="1.25.5"
REDIS_VERSION="7.2.6"

# Fungsi helper
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Fungsi untuk mengecek apakah script dijalankan sebagai root
check_root() {
    if [ "$EUID" -ne 0 ]; then 
        print_error "Script ini harus dijalankan sebagai root atau dengan sudo"
        exit 1
    fi
}

# Deteksi OS dan versi
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$NAME
        VER=$VERSION_ID
        print_info "Sistem operasi terdeteksi: $OS $VER"
    else
        print_error "Tidak dapat mendeteksi sistem operasi"
        exit 1
    fi
}

# Install dependencies berdasarkan OS
install_dependencies() {
    print_info "Menginstall dependencies..."
    
    if [[ "$OS" == *"Ubuntu"* ]] || [[ "$OS" == *"Debian"* ]]; then
        apt-get update
        apt-get install -y wget curl git build-essential
    elif [[ "$OS" == *"CentOS"* ]] || [[ "$OS" == *"Red Hat"* ]] || [[ "$OS" == *"Rocky"* ]]; then
        yum install -y wget curl git gcc make
    else
        print_error "OS tidak didukung: $OS"
        exit 1
    fi
    
    print_success "Dependencies berhasil diinstall"
}

# Install Go
install_go() {
    print_info "Mengecek instalasi Go..."
    
    if command -v go &> /dev/null; then
        CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_info "Go versi $CURRENT_GO_VERSION sudah terinstall"
        
        if [[ "$CURRENT_GO_VERSION" < "$GO_VERSION" ]]; then
            print_warning "Versi Go terlalu lama, mengupgrade ke $GO_VERSION..."
            remove_go
            install_go_binary
        fi
    else
        print_info "Go belum terinstall, menginstall Go $GO_VERSION..."
        install_go_binary
    fi
}

remove_go() {
    rm -rf /usr/local/go
}

install_go_binary() {
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)
            GO_ARCH="amd64"
            ;;
        aarch64)
            GO_ARCH="arm64"
            ;;
        armv7l)
            GO_ARCH="armv6l"
            ;;
        *)
            print_error "Arsitektur tidak didukung: $ARCH"
            exit 1
            ;;
    esac
    
    cd /tmp
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    
    if [ $? -ne 0 ]; then
        print_error "Gagal mendownload Go"
        exit 1
    fi
    
    tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    
    # Set Go path
    if ! grep -q "/usr/local/go/bin" /etc/profile; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
        echo 'export GOPATH=$HOME/go' >> /etc/profile
    fi
    
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    
    print_success "Go $GO_VERSION berhasil diinstall"
}

# Install Redis
install_redis() {
    print_info "Mengecek instalasi Redis..."
    
    if command -v redis-server &> /dev/null; then
        print_info "Redis sudah terinstall"
        return
    fi
    
    print_info "Menginstall Redis..."
    
    if [[ "$OS" == *"Ubuntu"* ]] || [[ "$OS" == *"Debian"* ]]; then
        apt-get install -y redis-server
    elif [[ "$OS" == *"CentOS"* ]] || [[ "$OS" == *"Red Hat"* ]] || [[ "$OS" == *"Rocky"* ]]; then
        yum install -y redis
    fi
    
    # Konfigurasi Redis
    systemctl enable redis-server 2>/dev/null || systemctl enable redis
    systemctl start redis-server 2>/dev/null || systemctl start redis
    
    print_success "Redis berhasil diinstall dan dijalankan"
}

# Konfigurasi Redis untuk production
configure_redis() {
    print_info "Mengkonfigurasi Redis untuk production..."
    
    REDIS_CONF="/etc/redis/redis.conf"
    if [ ! -f "$REDIS_CONF" ]; then
        REDIS_CONF="/etc/redis.conf"
    fi
    
    if [ -f "$REDIS_CONF" ]; then
        # Backup konfigurasi asli
        cp "$REDIS_CONF" "${REDIS_CONF}.backup"
        
        # Set password jika belum ada
        if ! grep -q "^requirepass" "$REDIS_CONF"; then
            REDIS_PASSWORD=$(openssl rand -base64 32)
            echo "requirepass $REDIS_PASSWORD" >> "$REDIS_CONF"
            print_success "Redis password telah diset"
            echo "$REDIS_PASSWORD" > /root/.redis_password
            chmod 600 /root/.redis_password
        fi
        
        # Restart Redis
        systemctl restart redis-server 2>/dev/null || systemctl restart redis
        print_success "Redis berhasil dikonfigurasi"
    else
        print_warning "File konfigurasi Redis tidak ditemukan"
    fi
}

# Buat user untuk aplikasi
create_app_user() {
    print_info "Membuat user aplikasi..."
    
    if id "$APP_USER" &>/dev/null; then
        print_info "User $APP_USER sudah ada"
    else
        useradd -r -s /bin/false "$APP_USER"
        print_success "User $APP_USER berhasil dibuat"
    fi
}

# Clone dan build aplikasi
build_application() {
    print_info "Clone repository dan build aplikasi..."
    
    # Buat direktori aplikasi
    mkdir -p "$APP_DIR"
    cd "$APP_DIR"
    
    # Clone repository
    if [ -d "$APP_DIR/.git" ]; then
        print_info "Repository sudah ada, melakukan pull..."
        cd "$APP_DIR"
        git pull origin main
    else
        print_info "Clone repository..."
        git clone "$REPO_URL" "$APP_DIR"
    fi
    
    cd "$APP_DIR"
    
    # Download dependencies
    print_info "Download Go dependencies..."
    /usr/local/go/bin/go mod download
    
    # Build aplikasi
    print_info "Build aplikasi..."
    CGO_ENABLED=0 GOOS=linux /usr/local/go/bin/go build \
        -ldflags="-s -w -X main.Version=1.0.0" \
        -o "$APP_DIR/bin/api" \
        ./cmd/api
    
    if [ $? -eq 0 ]; then
        print_success "Aplikasi berhasil di-build"
    else
        print_error "Gagal build aplikasi"
        exit 1
    fi
    
    # Set permissions
    chmod +x "$APP_DIR/bin/api"
}

# Konfigurasi environment
configure_environment() {
    print_info "Konfigurasi environment variables..."
    
    # Buat file .env dari template
    if [ -f "$APP_DIR/.env.example" ]; then
        cp "$APP_DIR/.env.example" "$APP_DIR/.env"
    fi
    
    # Prompt untuk konfigurasi
    echo ""
    print_info "Silakan masukkan konfigurasi SNMP OLT:"
    
    read -p "SNMP Host (IP OLT): " SNMP_HOST
    read -p "SNMP Port [161]: " SNMP_PORT
    SNMP_PORT=${SNMP_PORT:-161}
    read -p "SNMP Community String [public]: " SNMP_COMMUNITY
    SNMP_COMMUNITY=${SNMP_COMMUNITY:-public}
    
    echo ""
    print_info "Pilih versi firmware ZTE C320:"
    echo "  1) V2.1.x (default - untuk firmware lama)"
    echo "  2) V2.2.x atau lebih baru"
    read -p "Pilihan [1]: " FW_CHOICE
    FW_CHOICE=${FW_CHOICE:-1}
    
    if [ "$FW_CHOICE" == "2" ]; then
        ZTE_FIRMWARE_VERSION="v2.2"
    else
        ZTE_FIRMWARE_VERSION="v2.1"
    fi
    
    read -p "Server Port [8081]: " SERVER_PORT
    SERVER_PORT=${SERVER_PORT:-8081}
    
    # Dapatkan Redis password jika ada
    if [ -f /root/.redis_password ]; then
        REDIS_PASSWORD=$(cat /root/.redis_password)
    else
        REDIS_PASSWORD=""
    fi
    
    # Buat file .env
    cat > "$APP_DIR/.env" <<EOF
# Application Environment
APP_ENV=production

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=$SERVER_PORT
SERVER_MODE=release

# ZTE C320 Firmware Version
# v2.1 = Firmware V2.1.x (menggunakan OID structure yang berbeda)
# v2.2 = Firmware V2.2.x dan lebih baru
ZTE_FIRMWARE_VERSION=$ZTE_FIRMWARE_VERSION

# SNMP Configuration
SNMP_HOST=$SNMP_HOST
SNMP_PORT=$SNMP_PORT
SNMP_COMMUNITY=$SNMP_COMMUNITY

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=$REDIS_PASSWORD
REDIS_DB=0
REDIS_MIN_IDLE_CONNECTIONS=200
REDIS_POOL_SIZE=12000
REDIS_POOL_TIMEOUT=240

# TLS/HTTPS Configuration (Optional)
USE_TLS=false

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
EOF
    
    chmod 600 "$APP_DIR/.env"
    print_success "Environment berhasil dikonfigurasi"
}

# Buat systemd service
create_systemd_service() {
    print_info "Membuat systemd service..."
    
    cat > "/etc/systemd/system/${SERVICE_NAME}.service" <<EOF
[Unit]
Description=Go SNMP OLT ZTE C320 Monitoring Service
After=network.target redis.service
Wants=redis.service

[Service]
Type=simple
User=$APP_USER
Group=$APP_USER
WorkingDirectory=$APP_DIR
EnvironmentFile=$APP_DIR/.env
ExecStart=$APP_DIR/bin/api
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$SERVICE_NAME

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$APP_DIR

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
EOF
    
    print_success "Systemd service berhasil dibuat"
}

# Setup log directory
setup_logging() {
    print_info "Setup logging..."
    
    mkdir -p "$APP_DIR/logs"
    
    # Logrotate configuration
    cat > "/etc/logrotate.d/${SERVICE_NAME}" <<EOF
$APP_DIR/logs/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0640 $APP_USER $APP_USER
    sharedscripts
    postrotate
        systemctl reload ${SERVICE_NAME} > /dev/null 2>&1 || true
    endscript
}
EOF
    
    print_success "Logging berhasil dikonfigurasi"
}

# Set permissions
set_permissions() {
    print_info "Mengatur permissions..."
    
    chown -R "$APP_USER:$APP_USER" "$APP_DIR"
    chmod 755 "$APP_DIR/bin/api"
    chmod 600 "$APP_DIR/.env"
    
    print_success "Permissions berhasil diatur"
}

# Setup firewall
setup_firewall() {
    print_info "Konfigurasi firewall..."
    
    if command -v ufw &> /dev/null; then
        ufw allow "$SERVER_PORT/tcp" comment "Go SNMP OLT"
        print_success "UFW firewall rule ditambahkan untuk port $SERVER_PORT"
    elif command -v firewall-cmd &> /dev/null; then
        firewall-cmd --permanent --add-port="$SERVER_PORT/tcp"
        firewall-cmd --reload
        print_success "Firewalld rule ditambahkan untuk port $SERVER_PORT"
    else
        print_warning "Firewall tidak terdeteksi, silakan buka port $SERVER_PORT secara manual"
    fi
}

# Start service
start_service() {
    print_info "Memulai service..."
    
    systemctl daemon-reload
    systemctl enable "$SERVICE_NAME"
    systemctl start "$SERVICE_NAME"
    
    sleep 3
    
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        print_success "Service berhasil dijalankan"
        systemctl status "$SERVICE_NAME" --no-pager
    else
        print_error "Service gagal dijalankan"
        journalctl -u "$SERVICE_NAME" -n 50 --no-pager
        exit 1
    fi
}

# Test API
test_api() {
    print_info "Testing API endpoint..."
    
    sleep 2
    
    if command -v curl &> /dev/null; then
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:$SERVER_PORT/")
        
        if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 404 ]; then
            print_success "API berjalan dengan baik (HTTP $HTTP_CODE)"
        else
            print_warning "API mungkin belum siap (HTTP $HTTP_CODE)"
        fi
    fi
}

# Print summary
print_summary() {
    echo ""
    echo "=========================================="
    print_success "INSTALASI SELESAI!"
    echo "=========================================="
    echo ""
    echo "Informasi Service:"
    echo "  - Service Name: $SERVICE_NAME"
    echo "  - App Directory: $APP_DIR"
    echo "  - Binary: $APP_DIR/bin/api"
    echo "  - Config: $APP_DIR/.env"
    echo "  - Port: $SERVER_PORT"
    echo ""
    echo "Command berguna:"
    echo "  - Status service: systemctl status $SERVICE_NAME"
    echo "  - Stop service: systemctl stop $SERVICE_NAME"
    echo "  - Start service: systemctl start $SERVICE_NAME"
    echo "  - Restart service: systemctl restart $SERVICE_NAME"
    echo "  - View logs: journalctl -u $SERVICE_NAME -f"
    echo ""
    echo "Test API:"
    echo "  curl http://localhost:$SERVER_PORT/"
    echo ""
    if [ -f /root/.redis_password ]; then
        echo "Redis Password tersimpan di: /root/.redis_password"
        echo ""
    fi
    echo "=========================================="
}

# Fungsi uninstall
uninstall() {
    print_warning "Menghapus instalasi..."
    
    # Stop dan disable service
    systemctl stop "$SERVICE_NAME" 2>/dev/null || true
    systemctl disable "$SERVICE_NAME" 2>/dev/null || true
    
    # Hapus service file
    rm -f "/etc/systemd/system/${SERVICE_NAME}.service"
    systemctl daemon-reload
    
    # Hapus aplikasi
    rm -rf "$APP_DIR"
    
    # Hapus user
    userdel "$APP_USER" 2>/dev/null || true
    
    # Hapus logrotate config
    rm -f "/etc/logrotate.d/${SERVICE_NAME}"
    
    print_success "Uninstall selesai"
}

# Main installation flow
main() {
    echo ""
    echo "=========================================="
    echo "  Go SNMP OLT ZTE C320 - Installer VPS"
    echo "  Version: 1.0.0"
    echo "=========================================="
    echo ""
    
    # Parse arguments
    case "${1:-install}" in
        install)
            check_root
            detect_os
            install_dependencies
            install_go
            install_redis
            configure_redis
            create_app_user
            build_application
            configure_environment
            create_systemd_service
            setup_logging
            set_permissions
            setup_firewall
            start_service
            test_api
            print_summary
            ;;
        uninstall)
            check_root
            uninstall
            ;;
        update)
            check_root
            print_info "Updating aplikasi..."
            systemctl stop "$SERVICE_NAME"
            build_application
            set_permissions
            systemctl start "$SERVICE_NAME"
            print_success "Update selesai"
            ;;
        *)
            echo "Usage: $0 {install|uninstall|update}"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
