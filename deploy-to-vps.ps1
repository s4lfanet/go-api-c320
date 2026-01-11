# Deploy Phase 7.1 to VPS
# Usage: .\deploy-to-vps.ps1

$VPS_HOST = "192.168.54.230"
$VPS_USER = "root"
$APP_PATH = "/opt/go-snmp-olt"

Write-Host "=== Deploying Phase 7.1 to VPS ===" -ForegroundColor Cyan
Write-Host ""

# Step 1: Build for Linux
Write-Host "[1/5] Building for Linux AMD64..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o api cmd/api/main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "Build successful!" -ForegroundColor Green

# Step 2: Upload binary
Write-Host ""
Write-Host "[2/5] Uploading binary to VPS..." -ForegroundColor Yellow
scp api ${VPS_USER}@${VPS_HOST}:/tmp/api-phase71

if ($LASTEXITCODE -ne 0) {
    Write-Host "Upload failed!" -ForegroundColor Red
    exit 1
}
Write-Host "Upload successful!" -ForegroundColor Green

# Step 3: Deploy on VPS
Write-Host ""
Write-Host "[3/5] Deploying on VPS..." -ForegroundColor Yellow
ssh ${VPS_USER}@${VPS_HOST} @"
systemctl stop go-snmp-olt && \
mv /tmp/api-phase71 ${APP_PATH}/bin/api && \
chmod +x ${APP_PATH}/bin/api && \
systemctl start go-snmp-olt && \
sleep 2 && \
systemctl status go-snmp-olt --no-pager
"@

if ($LASTEXITCODE -ne 0) {
    Write-Host "Deployment failed!" -ForegroundColor Red
    exit 1
}
Write-Host "Deployment successful!" -ForegroundColor Green

# Step 4: Test endpoints
Write-Host ""
Write-Host "[4/5] Testing monitoring endpoints..." -ForegroundColor Yellow

Write-Host "Testing ONU monitoring..." -ForegroundColor Gray
$response = Invoke-RestMethod -Uri "http://${VPS_HOST}:8081/api/v1/monitoring/onu/1/1" -ErrorAction SilentlyContinue
if ($response) {
    Write-Host "  [OK] ONU endpoint working" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] ONU endpoint failed" -ForegroundColor Red
}

Write-Host "Testing PON monitoring..." -ForegroundColor Gray
$response = Invoke-RestMethod -Uri "http://${VPS_HOST}:8081/api/v1/monitoring/pon/1" -ErrorAction SilentlyContinue
if ($response) {
    Write-Host "  [OK] PON endpoint working" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] PON endpoint failed" -ForegroundColor Red
}

Write-Host "Testing OLT monitoring..." -ForegroundColor Gray
$response = Invoke-RestMethod -Uri "http://${VPS_HOST}:8081/api/v1/monitoring/olt" -ErrorAction SilentlyContinue
if ($response) {
    Write-Host "  [OK] OLT endpoint working" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] OLT endpoint failed" -ForegroundColor Red
}

# Step 5: Summary
Write-Host ""
Write-Host "[5/5] Deployment Summary" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host "Phase 7.1 - Real-time ONU Monitoring" -ForegroundColor White
Write-Host ""
Write-Host "Endpoints:" -ForegroundColor Yellow
Write-Host "  GET http://${VPS_HOST}:8081/api/v1/monitoring/onu/{pon}/{onuId}" -ForegroundColor Gray
Write-Host "  GET http://${VPS_HOST}:8081/api/v1/monitoring/pon/{pon}" -ForegroundColor Gray
Write-Host "  GET http://${VPS_HOST}:8081/api/v1/monitoring/olt" -ForegroundColor Gray
Write-Host ""
Write-Host "Test commands:" -ForegroundColor Yellow
Write-Host '  curl http://192.168.54.230:8081/api/v1/monitoring/onu/1/1 | jq' -ForegroundColor Gray
Write-Host '  curl http://192.168.54.230:8081/api/v1/monitoring/pon/1 | jq' -ForegroundColor Gray
Write-Host '  curl http://192.168.54.230:8081/api/v1/monitoring/olt | jq' -ForegroundColor Gray
Write-Host ""
Write-Host "Note: V2.1.0 has NO optical power monitoring" -ForegroundColor Magenta
Write-Host "      Available: Status, Traffic, Device Info" -ForegroundColor Magenta
Write-Host "=============================================" -ForegroundColor Cyan
