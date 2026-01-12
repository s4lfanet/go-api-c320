# ========================================
# Frontend Deployment Script untuk VPS
# PowerShell Version
# ZTE C320 OLT Management Dashboard
# ========================================
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
#   .\deploy-frontend.ps1 [-Environment production|staging]
#
# Prerequisites:
#   - SSH access ke VPS (key-based authentication)
#   - Nginx sudah terinstall di VPS
#   - Node.js dan npm terinstall (untuk build)
#   - OpenSSH Client (built-in di Windows 10+)
#
# ========================================

param(
    [string]$Environment = "production"
)

$ErrorActionPreference = "Stop"

# ===========================
# Configuration
# ===========================

$VPS_HOST = if ($env:VPS_HOST) { $env:VPS_HOST } else { "192.168.54.230" }
$VPS_USER = if ($env:VPS_USER) { $env:VPS_USER } else { "root" }
$VPS_PORT = if ($env:VPS_PORT) { $env:VPS_PORT } else { "22" }

$LOCAL_BUILD_DIR = "dist"
$REMOTE_BASE_DIR = "/var/www/olt-dashboard"
$REMOTE_FRONTEND_DIR = "$REMOTE_BASE_DIR/frontend"
$BACKUP_DIR = "$REMOTE_BASE_DIR/backups"

$TIMESTAMP = Get-Date -Format "yyyyMMdd_HHmmss"
$ARCHIVE_NAME = "frontend-$Environment-$TIMESTAMP.tar.gz"

# ===========================
# Functions
# ===========================

function Write-Step {
    param([string]$Message)
    Write-Host "==> " -NoNewline -ForegroundColor Blue
    Write-Host $Message
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ " -NoNewline -ForegroundColor Green
    Write-Host $Message
}

function Write-Failure {
    param([string]$Message)
    Write-Host "✗ " -NoNewline -ForegroundColor Red
    Write-Host $Message
}

function Write-Warning {
    param([string]$Message)
    Write-Host "! " -NoNewline -ForegroundColor Yellow
    Write-Host $Message
}

function Test-Command {
    param([string]$Command)
    $null = Get-Command $Command -ErrorAction SilentlyContinue
    return $?
}

# ===========================
# Validations
# ===========================

Write-Step "Validating prerequisites..."

# Check Node.js
if (-not (Test-Command "node")) {
    Write-Failure "Node.js is not installed"
    exit 1
}

# Check npm
if (-not (Test-Command "npm")) {
    Write-Failure "npm is not installed"
    exit 1
}

# Check SSH
if (-not (Test-Command "ssh")) {
    Write-Failure "SSH is not installed or not in PATH"
    Write-Warning "Install OpenSSH Client via Windows Features"
    exit 1
}

# Test SSH connection
$sshTest = ssh -q -o BatchMode=yes -o ConnectTimeout=5 -p $VPS_PORT "$VPS_USER@$VPS_HOST" "exit" 2>&1
if ($LASTEXITCODE -ne 0) {
    $vpsInfo = "$VPS_USER@$VPS_HOST" + ":" + $VPS_PORT
    Write-Failure "Cannot connect to VPS ($vpsInfo)"
    Write-Warning "Make sure SSH key-based authentication is configured"
    exit 1
}

Write-Success "All prerequisites met"

# ===========================
# Build Frontend
# ===========================

Write-Step "Building frontend for $Environment..."

# Clean previous build
if (Test-Path $LOCAL_BUILD_DIR) {
    Remove-Item -Recurse -Force $LOCAL_BUILD_DIR
}

# Install dependencies if needed
if (-not (Test-Path "node_modules")) {
    Write-Step "Installing dependencies..."
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Failure "npm install failed"
        exit 1
    }
}

# Build
if ($Environment -eq "production") {
    npm run build
} else {
    npm run build -- --mode $Environment
}

if ($LASTEXITCODE -ne 0) {
    Write-Failure "Build failed"
    exit 1
}

if (-not (Test-Path $LOCAL_BUILD_DIR)) {
    Write-Failure "Build failed - $LOCAL_BUILD_DIR not found"
    exit 1
}

Write-Success "Build completed"

# ===========================
# Create Archive
# ===========================

Write-Step "Creating archive..."

# Use tar if available (Windows 10 1803+)
if (Test-Command "tar") {
    tar -czf $ARCHIVE_NAME -C $LOCAL_BUILD_DIR .
    if ($LASTEXITCODE -ne 0) {
        Write-Failure "Failed to create archive"
        exit 1
    }
} else {
    # Fallback to PowerShell compression
    Compress-Archive -Path "$LOCAL_BUILD_DIR\*" -DestinationPath $ARCHIVE_NAME -Force
}

$archiveSize = (Get-Item $ARCHIVE_NAME).Length / 1MB
$archiveSizeStr = $archiveSize.ToString('0.00')
$archiveInfo = "Archive created: " + $ARCHIVE_NAME + " (" + $archiveSizeStr + " MB)"
Write-Success $archiveInfo

# ===========================
# Backup Current Version
# ===========================

Write-Step "Backing up current version on VPS..."

$backupScript = @'
mkdir -p {0}
if [ -d "{1}" ]; then
    BACKUP_NAME="frontend-backup-{2}.tar.gz"
    tar -czf {0}/$BACKUP_NAME -C {1} . 2>/dev/null || true
    cd {0}
    ls -t frontend-backup-*.tar.gz | tail -n +6 | xargs -r rm
fi
'@ -f $BACKUP_DIR, $REMOTE_FRONTEND_DIR, $TIMESTAMP

ssh -p $VPS_PORT "$VPS_USER@$VPS_HOST" $backupScript

Write-Success "Backup completed"

# ===========================
# Transfer to VPS
# ===========================

Write-Step "Transferring files to VPS..."

# Create temp directory
ssh -p $VPS_PORT "$VPS_USER@$VPS_HOST" "mkdir -p /tmp/frontend-deploy"

# Transfer archive
scp -P $VPS_PORT $ARCHIVE_NAME "${VPS_USER}@${VPS_HOST}:/tmp/frontend-deploy/"
if ($LASTEXITCODE -ne 0) {
    Write-Failure "Transfer failed"
    exit 1
}

Write-Success "Transfer completed"

# ===========================
# Deploy on VPS
# ===========================

Write-Step "Deploying on VPS..."

$deployScript = @'
set -e
mkdir -p {0}
rm -rf {0}/*
tar -xzf /tmp/frontend-deploy/{1} -C {0}
chown -R www-data:www-data {0}
chmod -R 755 {0}
rm -rf /tmp/frontend-deploy
nginx -t
systemctl reload nginx || service nginx reload
'@ -f $REMOTE_FRONTEND_DIR, $ARCHIVE_NAME

ssh -p $VPS_PORT "$VPS_USER@$VPS_HOST" $deployScript
if ($LASTEXITCODE -ne 0) {
    Write-Failure "Deployment failed"
    exit 1
}

Write-Success "Deployment completed"

# ===========================
# Cleanup Local Files
# ===========================

Write-Step "Cleaning up local files..."
Remove-Item -Force $ARCHIVE_NAME
Write-Success "Cleanup completed"

# ===========================
# Health Check
# ===========================

Write-Step "Running health check..."
Start-Sleep -Seconds 2

try {
    $response = Invoke-WebRequest -Uri "http://$VPS_HOST/" -UseBasicParsing -TimeoutSec 10
    $statusCode = $response.StatusCode
    
    if ($statusCode -eq 200) {
        Write-Success "Health check passed (HTTP $statusCode)"
    } else {
        Write-Warning "Health check returned HTTP $statusCode"
    }
} catch {
    Write-Warning "Health check failed: $($_.Exception.Message)"
    Write-Warning "Please verify manually: http://$VPS_HOST/"
}

# ===========================
# Completion
# ===========================

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "Deployment Successful!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Environment: $Environment"
Write-Host "VPS Host: $VPS_HOST"
Write-Host "Deployed to: $REMOTE_FRONTEND_DIR"
Write-Host "Backup location: $BACKUP_DIR"
Write-Host ""
Write-Host "Frontend URL: http://$VPS_HOST/"
Write-Host ""
Write-Host "Note: Frontend files are NOT pushed to GitHub" -ForegroundColor Blue
Write-Host ""
