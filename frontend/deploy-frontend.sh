#!/bin/bash

# =============================================================================
# Frontend Deployment Script untuk VPS
# ZTE C320 OLT Management Dashboard
# =============================================================================
# 
# Script ini akan:
# 1. Build production bundle
# 2. Compress files
# 3. Transfer ke VPS via SCP
# 4. Extract dan deploy di VPS
# 5. Reload Nginx
# 6. Health check
#
# Usage:
#   ./deploy-frontend.sh [environment]
#   environment: production | staging (default: production)
#
# Prerequisites:
#   - SSH access ke VPS (key-based authentication)
#   - Nginx sudah terinstall di VPS
#   - Node.js dan npm terinstall (untuk build)
#
# =============================================================================

set -e  # Exit on error

# ===========================
# Configuration
# ===========================

# Environment (default: production)
ENV=${1:-production}

# VPS Configuration
VPS_HOST="${VPS_HOST:-192.168.54.230}"
VPS_USER="${VPS_USER:-root}"
VPS_PORT="${VPS_PORT:-22}"

# Paths
LOCAL_BUILD_DIR="dist"
REMOTE_BASE_DIR="/var/www/olt-dashboard"
REMOTE_FRONTEND_DIR="${REMOTE_BASE_DIR}/frontend"
BACKUP_DIR="${REMOTE_BASE_DIR}/backups"

# Archive name with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
ARCHIVE_NAME="frontend-${ENV}-${TIMESTAMP}.tar.gz"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ===========================
# Functions
# ===========================

print_step() {
    echo -e "${BLUE}==>${NC} ${1}"
}

print_success() {
    echo -e "${GREEN}✓${NC} ${1}"
}

print_error() {
    echo -e "${RED}✗${NC} ${1}"
}

print_warning() {
    echo -e "${YELLOW}!${NC} ${1}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# ===========================
# Validations
# ===========================

print_step "Validating prerequisites..."

# Check Node.js
if ! command_exists node; then
    print_error "Node.js is not installed"
    exit 1
fi

# Check npm
if ! command_exists npm; then
    print_error "npm is not installed"
    exit 1
fi

# Check SSH connectivity
if ! ssh -q -o BatchMode=yes -o ConnectTimeout=5 -p ${VPS_PORT} ${VPS_USER}@${VPS_HOST} exit; then
    print_error "Cannot connect to VPS (${VPS_USER}@${VPS_HOST}:${VPS_PORT})"
    print_warning "Make sure SSH key-based authentication is configured"
    exit 1
fi

print_success "All prerequisites met"

# ===========================
# Build Frontend
# ===========================

print_step "Building frontend for ${ENV}..."

# Clean previous build
if [ -d "${LOCAL_BUILD_DIR}" ]; then
    rm -rf "${LOCAL_BUILD_DIR}"
fi

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    print_step "Installing dependencies..."
    npm install
fi

# Build
if [ "${ENV}" = "production" ]; then
    npm run build
else
    npm run build -- --mode ${ENV}
fi

if [ ! -d "${LOCAL_BUILD_DIR}" ]; then
    print_error "Build failed - ${LOCAL_BUILD_DIR} not found"
    exit 1
fi

print_success "Build completed"

# ===========================
# Create Archive
# ===========================

print_step "Creating archive..."

# Create compressed archive
tar -czf "${ARCHIVE_NAME}" -C "${LOCAL_BUILD_DIR}" .

if [ ! -f "${ARCHIVE_NAME}" ]; then
    print_error "Failed to create archive"
    exit 1
fi

ARCHIVE_SIZE=$(du -h "${ARCHIVE_NAME}" | cut -f1)
print_success "Archive created: ${ARCHIVE_NAME} (${ARCHIVE_SIZE})"

# ===========================
# Backup Current Version
# ===========================

print_step "Backing up current version on VPS..."

ssh -p ${VPS_PORT} ${VPS_USER}@${VPS_HOST} << EOF
    # Create backup directory if not exists
    mkdir -p ${BACKUP_DIR}
    
    # Backup current frontend if exists
    if [ -d "${REMOTE_FRONTEND_DIR}" ]; then
        BACKUP_NAME="frontend-backup-${TIMESTAMP}.tar.gz"
        tar -czf ${BACKUP_DIR}/\${BACKUP_NAME} -C ${REMOTE_FRONTEND_DIR} . 2>/dev/null || true
        
        # Keep only last 5 backups
        cd ${BACKUP_DIR}
        ls -t frontend-backup-*.tar.gz | tail -n +6 | xargs -r rm
    fi
EOF

print_success "Backup completed"

# ===========================
# Transfer to VPS
# ===========================

print_step "Transferring files to VPS..."

# Create temp directory on VPS
ssh -p ${VPS_PORT} ${VPS_USER}@${VPS_HOST} "mkdir -p /tmp/frontend-deploy"

# Transfer archive
scp -P ${VPS_PORT} "${ARCHIVE_NAME}" ${VPS_USER}@${VPS_HOST}:/tmp/frontend-deploy/

print_success "Transfer completed"

# ===========================
# Deploy on VPS
# ===========================

print_step "Deploying on VPS..."

ssh -p ${VPS_PORT} ${VPS_USER}@${VPS_HOST} << EOF
    set -e
    
    # Create frontend directory
    mkdir -p ${REMOTE_FRONTEND_DIR}
    
    # Remove old files
    rm -rf ${REMOTE_FRONTEND_DIR}/*
    
    # Extract new files
    tar -xzf /tmp/frontend-deploy/${ARCHIVE_NAME} -C ${REMOTE_FRONTEND_DIR}
    
    # Set permissions
    chown -R www-data:www-data ${REMOTE_FRONTEND_DIR}
    chmod -R 755 ${REMOTE_FRONTEND_DIR}
    
    # Cleanup temp files
    rm -rf /tmp/frontend-deploy
    
    # Test Nginx configuration
    nginx -t
    
    # Reload Nginx
    systemctl reload nginx || service nginx reload
EOF

print_success "Deployment completed"

# ===========================
# Cleanup Local Files
# ===========================

print_step "Cleaning up local files..."

# Remove archive
rm -f "${ARCHIVE_NAME}"

print_success "Cleanup completed"

# ===========================
# Health Check
# ===========================

print_step "Running health check..."

# Wait a moment for service to be ready
sleep 2

# Try to access the frontend
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://${VPS_HOST}/ || echo "000")

if [ "${HTTP_STATUS}" = "200" ]; then
    print_success "Health check passed (HTTP ${HTTP_STATUS})"
else
    print_warning "Health check returned HTTP ${HTTP_STATUS}"
    print_warning "Please verify manually: http://${VPS_HOST}/"
fi

# ===========================
# Completion
# ===========================

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Deployment Successful!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Environment: ${ENV}"
echo "VPS Host: ${VPS_HOST}"
echo "Deployed to: ${REMOTE_FRONTEND_DIR}"
echo "Backup location: ${BACKUP_DIR}"
echo ""
echo "Frontend URL: http://${VPS_HOST}/"
echo ""
echo -e "${BLUE}Note: Frontend files are NOT pushed to GitHub${NC}"
echo ""
