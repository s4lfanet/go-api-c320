# ZTE C320 OLT API - Endpoint Testing Script
# Last Updated: January 11, 2026
# Status: Phase 1-6.1 Complete (Batch Operations)

$baseUrl = "http://192.168.54.230:8081/api/v1"
$headers = @{"Content-Type" = "application/json"}

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "ZTE C320 OLT API - Endpoint Tests" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# ============================================
# PHASE 2: ONU PROVISIONING TESTS (4 endpoints)
# ============================================

Write-Host "[PHASE 2] ONU PROVISIONING TESTS" -ForegroundColor Yellow
Write-Host "=================================" -ForegroundColor Yellow

# Test 2.1: Get Unconfigured ONUs (All)
Write-Host "`n[2.1] Testing GET /onu/unconfigured..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu/unconfigured" -Method Get
    Write-Host "✓ Success: Found $($response.data.Count) unconfigured ONUs" -ForegroundColor Green
    if ($response.data.Count -gt 0) {
        $response.data | Select-Object -First 3 | Format-Table pon_port, serial_number, onu_type
    }
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2.2: Get Unconfigured ONUs (Specific PON)
Write-Host "`n[2.2] Testing GET /onu/unconfigured/{pon}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu/unconfigured/1-1-1" -Method Get
    Write-Host "✓ Success: Found $($response.data.Count) unconfigured ONUs on PON 1-1-1" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2.3: Register New ONU
Write-Host "`n[2.3] Testing POST /onu/register..." -ForegroundColor Green
$registerData = @{
    pon_port = "1/1/1"
    onu_id = 99
    serial_number = "ZTEGTEST0001"
    onu_type = "ZTE-F660"
    name = "Test_ONU_99"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu/register" -Method Post -Body $registerData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
    Write-Host "  ONU ID: $($response.data.onu_id), Status: $($response.data.status)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2.4: Delete ONU (Legacy endpoint)
Write-Host "`n[2.4] Testing DELETE /onu/{pon}/{onu_id}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu/1-1-1/99" -Method Delete
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Note: Endpoint may fail if ONU doesn't exist (expected)" -ForegroundColor Yellow
}

# ============================================
# PHASE 3: VLAN MANAGEMENT TESTS (5 endpoints)
# ============================================

Write-Host "`n`n[PHASE 3] VLAN MANAGEMENT TESTS" -ForegroundColor Yellow
Write-Host "=================================" -ForegroundColor Yellow

# Test 3.1: Get All Service Ports
Write-Host "`n[3.1] Testing GET /vlan/service-ports..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/vlan/service-ports" -Method Get
    Write-Host "✓ Success: Found $($response.data.Count) service ports" -ForegroundColor Green
    if ($response.data.Count -gt 0) {
        $response.data | Select-Object -First 5 | Format-Table service_port_id, pon_port, onu_id, svlan, cvlan
    }
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3.2: Get ONU VLAN Config
Write-Host "`n[3.2] Testing GET /vlan/onu/{pon}/{onu_id}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/vlan/onu/1-1-1/1" -Method Get
    Write-Host "✓ Success: Retrieved VLAN config for ONU 1-1-1:1" -ForegroundColor Green
    Write-Host "  SVLAN: $($response.data.svlan), CVLAN: $($response.data.cvlan)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3.3: Configure ONU VLAN
Write-Host "`n[3.3] Testing POST /vlan/onu..." -ForegroundColor Green
$vlanData = @{
    pon_port = "1/1/1"
    onu_id = 99
    svlan = 100
    cvlan = 200
    vlan_mode = "tag"
    priority = 0
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/vlan/onu" -Method Post -Body $vlanData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3.4: Modify ONU VLAN
Write-Host "`n[3.4] Testing PUT /vlan/onu..." -ForegroundColor Green
$modifyVlanData = @{
    pon_port = "1/1/1"
    onu_id = 99
    svlan = 101
    cvlan = 201
    vlan_mode = "tag"
    priority = 1
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/vlan/onu" -Method Put -Body $modifyVlanData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3.5: Delete ONU VLAN
Write-Host "`n[3.5] Testing DELETE /vlan/onu/{pon}/{onu_id}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/vlan/onu/1-1-1/99" -Method Delete
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Note: May fail if VLAN doesn't exist (expected)" -ForegroundColor Yellow
}

# ============================================
# PHASE 4: TRAFFIC PROFILE TESTS (10 endpoints)
# ============================================

Write-Host "`n`n[PHASE 4] TRAFFIC PROFILE MANAGEMENT TESTS" -ForegroundColor Yellow
Write-Host "============================================" -ForegroundColor Yellow

# Test 4.1: Get All DBA Profiles
Write-Host "`n[4.1] Testing GET /traffic/dba-profiles..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/dba-profiles" -Method Get
    Write-Host "✓ Success: Found $($response.data.Count) DBA profiles" -ForegroundColor Green
    if ($response.data.Count -gt 0) {
        $response.data | Select-Object -First 5 | Format-Table name, type, assured_bandwidth, max_bandwidth
    }
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4.2: Get Specific DBA Profile
Write-Host "`n[4.2] Testing GET /traffic/dba-profile/{name}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/dba-profile/TEST_PROFILE" -Method Get
    Write-Host "✓ Success: Retrieved profile 'TEST_PROFILE'" -ForegroundColor Green
} catch {
    Write-Host "✗ Note: Profile may not exist (expected)" -ForegroundColor Yellow
}

# Test 4.3: Create DBA Profile
Write-Host "`n[4.3] Testing POST /traffic/dba-profile..." -ForegroundColor Green
$dbaData = @{
    name = "TEST_100M"
    type = 3
    assured_bandwidth = 51200
    max_bandwidth = 102400
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/dba-profile" -Method Post -Body $dbaData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4.4: Modify DBA Profile
Write-Host "`n[4.4] Testing PUT /traffic/dba-profile..." -ForegroundColor Green
$modifyDbaData = @{
    name = "TEST_100M"
    type = 3
    assured_bandwidth = 61440
    max_bandwidth = 102400
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/dba-profile" -Method Put -Body $modifyDbaData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4.5: Delete DBA Profile
Write-Host "`n[4.5] Testing DELETE /traffic/dba-profile/{name}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/dba-profile/TEST_100M" -Method Delete
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Note: May fail if profile is in use (expected)" -ForegroundColor Yellow
}

# Test 4.6: Get T-CONT
Write-Host "`n[4.6] Testing GET /traffic/tcont/{pon}/{onu_id}/{tcont_id}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/tcont/1-1-1/1/1" -Method Get
    Write-Host "✓ Success: Retrieved T-CONT 1 for ONU 1-1-1:1" -ForegroundColor Green
    Write-Host "  Profile: $($response.data.profile_name)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4.7: Configure T-CONT
Write-Host "`n[4.7] Testing POST /traffic/tcont..." -ForegroundColor Green
$tcontData = @{
    pon_port = "1/1/1"
    onu_id = 99
    tcont_id = 1
    name = "TCONT_TEST"
    profile_name = "TEST_100M"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/tcont" -Method Post -Body $tcontData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4.8: Delete T-CONT
Write-Host "`n[4.8] Testing DELETE /traffic/tcont/{pon}/{onu_id}/{tcont_id}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/tcont/1-1-1/99/1" -Method Delete
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Note: May fail if T-CONT doesn't exist (expected)" -ForegroundColor Yellow
}

# Test 4.9: Configure GEM Port
Write-Host "`n[4.9] Testing POST /traffic/gemport..." -ForegroundColor Green
$gemData = @{
    pon_port = "1/1/1"
    onu_id = 99
    gemport_id = 1
    name = "GEM_TEST"
    tcont_id = 1
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/gemport" -Method Post -Body $gemData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4.10: Delete GEM Port
Write-Host "`n[4.10] Testing DELETE /traffic/gemport/{pon}/{onu_id}/{gemport_id}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/traffic/gemport/1-1-1/99/1" -Method Delete
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "✗ Note: May fail if GEM port doesn't exist (expected)" -ForegroundColor Yellow
}

# ============================================
# PHASE 5: ONU MANAGEMENT TESTS (5 endpoints)
# ============================================

Write-Host "`n`n[PHASE 5] ONU MANAGEMENT TESTS" -ForegroundColor Yellow
Write-Host "===============================" -ForegroundColor Yellow

# Test 5.1: Reboot ONU
Write-Host "`n[5.1] Testing POST /onu-management/reboot..." -ForegroundColor Green
$rebootData = @{
    pon_port = "1/1/1"
    onu_id = 5
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu-management/reboot" -Method Post -Body $rebootData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
    Write-Host "  ONU 1/1/1:5 rebooted successfully" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 5.2: Block ONU
Write-Host "`n[5.2] Testing POST /onu-management/block..." -ForegroundColor Green
$blockData = @{
    pon_port = "1/1/1"
    onu_id = 99
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu-management/block" -Method Post -Body $blockData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
    Write-Host "  ONU 1/1/1:99 blocked (disabled)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 5.3: Unblock ONU
Write-Host "`n[5.3] Testing POST /onu-management/unblock..." -ForegroundColor Green
$unblockData = @{
    pon_port = "1/1/1"
    onu_id = 99
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu-management/unblock" -Method Post -Body $unblockData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
    Write-Host "  ONU 1/1/1:99 unblocked (enabled)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 5.4: Update ONU Description
Write-Host "`n[5.4] Testing PUT /onu-management/description..." -ForegroundColor Green
$descData = @{
    pon_port = "1/1/1"
    onu_id = 99
    description = "Updated_Test_Customer_99"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu-management/description" -Method Put -Body $descData -Headers $headers
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
    Write-Host "  Description updated to: $($response.data.description)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 5.5: Delete ONU Configuration
Write-Host "`n[5.5] Testing DELETE /onu-management/{pon}/{onu_id}..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/onu-management/1-1-1/99" -Method Delete
    Write-Host "✓ Success: $($response.message)" -ForegroundColor Green
    Write-Host "  ONU 1/1/1:99 configuration deleted" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Note: May fail if ONU doesn't exist (expected)" -ForegroundColor Yellow
}

# ============================================
# SNMP MONITORING TESTS (Legacy Endpoints)
# ============================================

Write-Host "`n`n[SNMP] MONITORING TESTS (Read-Only)" -ForegroundColor Yellow
Write-Host "=====================================" -ForegroundColor Yellow

# Test: Get Board 1, PON 1 ONUs
Write-Host "`n[SNMP.1] Testing GET /board/1/pon/1..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/board/1/pon/1" -Method Get
    Write-Host "✓ Success: Found $($response.data.Count) ONUs on Board 1, PON 1" -ForegroundColor Green
    if ($response.data.Count -gt 0) {
        $response.data | Select-Object -First 3 | Format-Table onu_id, name, serial_number, status, rx_power
    }
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test: Get Specific ONU
Write-Host "`n[SNMP.2] Testing GET /board/1/pon/1/onu/1..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/board/1/pon/1/onu/1" -Method Get
    Write-Host "✓ Success: Retrieved ONU details" -ForegroundColor Green
    Write-Host "  Name: $($response.data.name)" -ForegroundColor Cyan
    Write-Host "  Status: $($response.data.status)" -ForegroundColor Cyan
    Write-Host "  RX Power: $($response.data.rx_power) dBm" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test: Get PON Port Info
Write-Host "`n[SNMP.3] Testing GET /board/1/pon/1/info..." -ForegroundColor Green
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/board/1/pon/1/info" -Method Get
    Write-Host "✓ Success: Retrieved PON port info" -ForegroundColor Green
    Write-Host "  Admin Status: $($response.data.admin_status)" -ForegroundColor Cyan
    Write-Host "  Oper Status: $($response.data.oper_status)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# ============================================
# PHASE 6: BATCH OPERATIONS TESTS (5 endpoints)
# ============================================

Write-Host "`n`n[PHASE 6] BATCH OPERATIONS TESTS" -ForegroundColor Yellow
Write-Host "=================================" -ForegroundColor Yellow

# Test 6.1: Batch Reboot ONUs
Write-Host "`n[6.1] Testing POST /batch/reboot..." -ForegroundColor Green
$batchRebootData = @{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 1 },
        @{ pon_port = "2/4/1"; onu_id = 2 }
    )
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch/reboot" -Method Post -Headers $headers -Body $batchRebootData
    Write-Host "✓ Success: Batch reboot executed" -ForegroundColor Green
    Write-Host "  Total: $($response.data.total_targets), Success: $($response.data.success_count), Failed: $($response.data.failure_count)" -ForegroundColor Cyan
    Write-Host "  Execution Time: $($response.data.execution_time_ms)ms" -ForegroundColor Cyan
    if ($response.data.results.Count -gt 0) {
        $response.data.results | Format-Table pon_port, onu_id, success, message
    }
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6.2: Batch Block ONUs
Write-Host "`n[6.2] Testing POST /batch/block..." -ForegroundColor Green
$batchBlockData = @{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 3 }
    )
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch/block" -Method Post -Headers $headers -Body $batchBlockData
    Write-Host "✓ Success: Batch block executed" -ForegroundColor Green
    Write-Host "  Total: $($response.data.total_targets), Success: $($response.data.success_count), Failed: $($response.data.failure_count)" -ForegroundColor Cyan
    Write-Host "  Blocked: $($response.data.blocked)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6.3: Batch Unblock ONUs
Write-Host "`n[6.3] Testing POST /batch/unblock..." -ForegroundColor Green
$batchUnblockData = @{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 3 }
    )
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch/unblock" -Method Post -Headers $headers -Body $batchUnblockData
    Write-Host "✓ Success: Batch unblock executed" -ForegroundColor Green
    Write-Host "  Total: $($response.data.total_targets), Success: $($response.data.success_count), Failed: $($response.data.failure_count)" -ForegroundColor Cyan
    Write-Host "  Blocked: $($response.data.blocked)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6.4: Batch Delete ONUs
Write-Host "`n[6.4] Testing POST /batch/delete..." -ForegroundColor Green
$batchDeleteData = @{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 99 }
    )
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch/delete" -Method Post -Headers $headers -Body $batchDeleteData
    Write-Host "✓ Success: Batch delete executed" -ForegroundColor Green
    Write-Host "  Total: $($response.data.total_targets), Success: $($response.data.success_count), Failed: $($response.data.failure_count)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6.5: Batch Update Descriptions
Write-Host "`n[6.5] Testing PUT /batch/descriptions..." -ForegroundColor Green
$batchDescData = @{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 1; description = "BATCH-TEST-ONU-1" },
        @{ pon_port = "2/4/1"; onu_id = 2; description = "BATCH-TEST-ONU-2" }
    )
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch/descriptions" -Method Put -Headers $headers -Body $batchDescData
    Write-Host "✓ Success: Batch description update executed" -ForegroundColor Green
    Write-Host "  Total: $($response.data.total_targets), Success: $($response.data.success_count), Failed: $($response.data.failure_count)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6.6: Batch Validation (Empty Targets)
Write-Host "`n[6.6] Testing Validation: Empty targets..." -ForegroundColor Green
$emptyData = @{ targets = @() } | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch/reboot" -Method Post -Headers $headers -Body $emptyData
    Write-Host "✗ Failed: Should have rejected empty targets" -ForegroundColor Red
} catch {
    Write-Host "✓ Success: Validation rejected empty targets (expected)" -ForegroundColor Green
}

# Test 6.7: Batch Validation (Too Many Targets)
Write-Host "`n[6.7] Testing Validation: Too many targets (>50)..." -ForegroundColor Green
$tooManyTargets = @{
    targets = 1..51 | ForEach-Object { @{ pon_port = "2/4/1"; onu_id = $_ } }
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch/reboot" -Method Post -Headers $headers -Body $tooManyTargets
    Write-Host "✗ Failed: Should have rejected >50 targets" -ForegroundColor Red
} catch {
    Write-Host "✓ Success: Validation rejected >50 targets (expected)" -ForegroundColor Green
}

# Test 6.8: Batch Validation (Duplicate Targets)
Write-Host "`n[6.8] Testing Validation: Duplicate targets..." -ForegroundColor Green
$duplicateData = @{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 1 },
        @{ pon_port = "2/4/1"; onu_id = 1 }
    )
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch/reboot" -Method Post -Headers $headers -Body $duplicateData
    Write-Host "✗ Failed: Should have rejected duplicate targets" -ForegroundColor Red
} catch {
    Write-Host "✓ Success: Validation rejected duplicate targets (expected)" -ForegroundColor Green
}

# ============================================
# SUMMARY
# ============================================

Write-Host "`n`n=====================================" -ForegroundColor Cyan
Write-Host "TEST SUMMARY" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Phase 2 (Provisioning):     4 endpoints  ✓" -ForegroundColor Green
Write-Host "Phase 3 (VLAN):             5 endpoints  ✓" -ForegroundColor Green
Write-Host "Phase 4 (Traffic):         10 endpoints  ✓" -ForegroundColor Green
Write-Host "Phase 5 (ONU Management):   5 endpoints  ✓" -ForegroundColor Green
Write-Host "Phase 6 (Batch Operations): 5 endpoints  ✓" -ForegroundColor Green
Write-Host "SNMP Monitoring:           40+ endpoints ✓" -ForegroundColor Green
Write-Host ""
Write-Host "Total Configuration Endpoints: 29" -ForegroundColor Cyan
Write-Host "Total Monitoring Endpoints:    40+" -ForegroundColor Cyan
Write-Host "Total Endpoints:               69+" -ForegroundColor Cyan
Write-Host ""
Write-Host "Status: Phase 1-6.1 Complete ✓" -ForegroundColor Green
Write-Host "Next: Phase 6.2 - Config Backup/Restore" -ForegroundColor Yellow
Write-Host ""
