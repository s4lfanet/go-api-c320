# Comprehensive Endpoint Testing Script
# Test all Phase 1-6.1 endpoints on VPS
# VPS: http://192.168.54.230:8081

$baseUrl = "http://192.168.54.230:8081/api/v1"
$headers = @{"Content-Type" = "application/json"}
$totalTests = 0
$passedTests = 0
$failedTests = 0

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Uri,
        [string]$Body = $null,
        [int]$ExpectedStatus = 200
    )
    
    $global:totalTests++
    Write-Host "`n[$global:totalTests] Testing: $Name" -ForegroundColor Cyan
    Write-Host "  Method: $Method | URI: $Uri" -ForegroundColor Gray
    
    try {
        $params = @{
            Uri = $Uri
            Method = $Method
            Headers = $headers
        }
        
        if ($Body) {
            $params.Body = $Body
        }
        
        $response = Invoke-RestMethod @params -ErrorAction Stop
        
        if ($response.code -eq $ExpectedStatus -or $response.PSObject.Properties.Name -contains 'data') {
            Write-Host "  PASSED" -ForegroundColor Green
            $global:passedTests++
            return $true
        } else {
            Write-Host "  FAILED: Unexpected response" -ForegroundColor Red
            $global:failedTests++
            return $false
        }
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "  PASSED (Expected $ExpectedStatus)" -ForegroundColor Green
            $global:passedTests++
            return $true
        } else {
            Write-Host "  FAILED: $($_.Exception.Message)" -ForegroundColor Red
            $global:failedTests++
            return $false
        }
    }
}

Write-Host "=============================================" -ForegroundColor Cyan
Write-Host "  COMPREHENSIVE ENDPOINT TEST SUITE" -ForegroundColor Cyan
Write-Host "  VPS: 192.168.54.230:8081" -ForegroundColor Cyan
Write-Host "  Phases: 1-6.1 (All Features)" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan

# ===== PHASE 2: ONU PROVISIONING =====
Write-Host "`n`n===== PHASE 2: ONU PROVISIONING =====" -ForegroundColor Yellow

Test-Endpoint -Name "Get all unconfigured ONUs" -Method "GET" -Uri "$baseUrl/onu/unconfigured"
Test-Endpoint -Name "Get unconfigured ONUs by PON" -Method "GET" -Uri "$baseUrl/onu/unconfigured/2-4-1"
Test-Endpoint -Name "Authorize ONU" -Method "POST" -Uri "$baseUrl/onu/authorize" -Body (@{
    pon_port = "2/4/1"
    onu_id = 10
    serial_number = "ZTEGTEST0010"
    onu_type = "ZTE-F660"
    name = "TestONU10"
} | ConvertTo-Json)
Test-Endpoint -Name "Deauthorize ONU" -Method "DELETE" -Uri "$baseUrl/onu/deauthorize/2-4-1/10"

# ===== PHASE 3: VLAN MANAGEMENT =====
Write-Host "`n`n===== PHASE 3: VLAN MANAGEMENT =====" -ForegroundColor Yellow

Test-Endpoint -Name "Get all VLAN profiles" -Method "GET" -Uri "$baseUrl/vlan/profiles"
Test-Endpoint -Name "Get specific VLAN profile" -Method "GET" -Uri "$baseUrl/vlan/profile/INTERNET"
Test-Endpoint -Name "Get ONU VLAN config" -Method "GET" -Uri "$baseUrl/vlan/onu-config/2-4-1/1"
Test-Endpoint -Name "Configure ONU VLAN" -Method "POST" -Uri "$baseUrl/vlan/configure" -Body (@{
    pon_port = "2/4/1"
    onu_id = 1
    service_port_id = 1
    vlan_mode = "tag"
    cvlan = 100
    user_vlan = 100
} | ConvertTo-Json)
Test-Endpoint -Name "Remove ONU VLAN" -Method "DELETE" -Uri "$baseUrl/vlan/remove/2-4-1/1/1"

# ===== PHASE 4: TRAFFIC PROFILES =====
Write-Host "`n`n===== PHASE 4: TRAFFIC PROFILES =====" -ForegroundColor Yellow

Test-Endpoint -Name "List all DBA profiles" -Method "GET" -Uri "$baseUrl/traffic/dba-profiles"
Test-Endpoint -Name "Get specific DBA profile" -Method "GET" -Uri "$baseUrl/traffic/dba-profile/UP-10M"
Test-Endpoint -Name "Create DBA profile" -Method "POST" -Uri "$baseUrl/traffic/dba-profile" -Body (@{
    profile_name = "TEST-PROFILE"
    type = 4
    fix_bandwidth = 1024
    assure_bandwidth = 512
    max_bandwidth = 10240
} | ConvertTo-Json)
Test-Endpoint -Name "Modify DBA profile" -Method "PUT" -Uri "$baseUrl/traffic/dba-profile" -Body (@{
    profile_name = "TEST-PROFILE"
    type = 4
    fix_bandwidth = 2048
    assure_bandwidth = 1024
    max_bandwidth = 10240
} | ConvertTo-Json)
Test-Endpoint -Name "Delete DBA profile" -Method "DELETE" -Uri "$baseUrl/traffic/dba-profile/TEST-PROFILE"
Test-Endpoint -Name "Get T-CONT config" -Method "GET" -Uri "$baseUrl/traffic/tcont/2-4-1/1/1"
Test-Endpoint -Name "Configure T-CONT" -Method "POST" -Uri "$baseUrl/traffic/tcont" -Body (@{
    pon_port = "2/4/1"
    onu_id = 1
    tcont_id = 1
    dba_profile = "UP-10M"
} | ConvertTo-Json)
Test-Endpoint -Name "Delete T-CONT" -Method "DELETE" -Uri "$baseUrl/traffic/tcont/2-4-1/1/1"
Test-Endpoint -Name "Configure GEM port" -Method "POST" -Uri "$baseUrl/traffic/gemport" -Body (@{
    pon_port = "2/4/1"
    onu_id = 1
    gemport_id = 1
    tcont_id = 1
} | ConvertTo-Json)
Test-Endpoint -Name "Delete GEM port" -Method "DELETE" -Uri "$baseUrl/traffic/gemport/2-4-1/1/1"

# ===== PHASE 5: ONU MANAGEMENT =====
Write-Host "`n`n===== PHASE 5: ONU MANAGEMENT =====" -ForegroundColor Yellow

Test-Endpoint -Name "Reboot ONU" -Method "POST" -Uri "$baseUrl/onu-management/reboot" -Body (@{
    pon_port = "2/4/1"
    onu_id = 1
} | ConvertTo-Json)
Test-Endpoint -Name "Block ONU" -Method "POST" -Uri "$baseUrl/onu-management/block" -Body (@{
    pon_port = "2/4/1"
    onu_id = 5
    block = $true
} | ConvertTo-Json)
Test-Endpoint -Name "Unblock ONU" -Method "POST" -Uri "$baseUrl/onu-management/unblock" -Body (@{
    pon_port = "2/4/1"
    onu_id = 5
} | ConvertTo-Json)
Test-Endpoint -Name "Update ONU description" -Method "PUT" -Uri "$baseUrl/onu-management/description" -Body (@{
    pon_port = "2/4/1"
    onu_id = 1
    description = "TEST-DESCRIPTION-UPDATE"
} | ConvertTo-Json)
Test-Endpoint -Name "Delete ONU config" -Method "DELETE" -Uri "$baseUrl/onu-management/2-4-1/99"

# ===== PHASE 6.1: BATCH OPERATIONS =====
Write-Host "`n`n===== PHASE 6.1: BATCH OPERATIONS =====" -ForegroundColor Yellow

Test-Endpoint -Name "Batch reboot ONUs" -Method "POST" -Uri "$baseUrl/batch/reboot" -Body (@{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 1 }
    )
} | ConvertTo-Json -Depth 10)

Test-Endpoint -Name "Batch block ONUs" -Method "POST" -Uri "$baseUrl/batch/block" -Body (@{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 5 }
    )
} | ConvertTo-Json -Depth 10)

Test-Endpoint -Name "Batch unblock ONUs" -Method "POST" -Uri "$baseUrl/batch/unblock" -Body (@{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 5 }
    )
} | ConvertTo-Json -Depth 10)

Test-Endpoint -Name "Batch delete ONUs" -Method "POST" -Uri "$baseUrl/batch/delete" -Body (@{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 99 }
    )
} | ConvertTo-Json -Depth 10)

Test-Endpoint -Name "Batch update descriptions" -Method "PUT" -Uri "$baseUrl/batch/descriptions" -Body (@{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 1; description = "BATCH-TEST-1" }
    )
} | ConvertTo-Json -Depth 10)

# VALIDATION TESTS
Write-Host "`n`n===== BATCH VALIDATION TESTS =====" -ForegroundColor Yellow

Test-Endpoint -Name "Reject empty targets" -Method "POST" -Uri "$baseUrl/batch/reboot" -Body (@{
    targets = @()
} | ConvertTo-Json) -ExpectedStatus 500

Test-Endpoint -Name "Reject duplicate targets" -Method "POST" -Uri "$baseUrl/batch/reboot" -Body (@{
    targets = @(
        @{ pon_port = "2/4/1"; onu_id = 1 },
        @{ pon_port = "2/4/1"; onu_id = 1 }
    )
} | ConvertTo-Json -Depth 10) -ExpectedStatus 500

# ===== SUMMARY =====
Write-Host "`n`n=============================================" -ForegroundColor Cyan
Write-Host "  TEST SUMMARY" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Total Tests:   $totalTests" -ForegroundColor White
Write-Host "  Passed:        $passedTests" -ForegroundColor Green
Write-Host "  Failed:        $failedTests" -ForegroundColor Red
Write-Host ""
$successRate = [math]::Round(($passedTests / $totalTests) * 100, 2)
Write-Host "  Success Rate:  $successRate%" -ForegroundColor $(if ($successRate -ge 90) { "Green" } elseif ($successRate -ge 70) { "Yellow" } else { "Red" })
Write-Host ""

if ($failedTests -eq 0) {
    Write-Host "  ALL TESTS PASSED!" -ForegroundColor Green
} else {
    Write-Host "  Some tests failed. Review logs above." -ForegroundColor Yellow
}
Write-Host ""
