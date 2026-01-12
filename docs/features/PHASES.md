# ZTE C320 Telnet Configuration Module - Roadmap & Documentation

**Version:** 1.0.0  
**Firmware Target:** ZTE C320 V2.1.0  
**Last Updated:** January 11, 2026  

---

## 📋 Table of Contents

1. [Executive Summary](#executive-summary)
2. [Current State Analysis](#current-state-analysis)
3. [Proposed Architecture](#proposed-architecture)
4. [Roadmap & Phases](#roadmap--phases)
5. [Telnet Command Reference](#telnet-command-reference)
6. [API Endpoints Design](#api-endpoints-design)
7. [Workflow Diagrams](#workflow-diagrams)
8. [Security Considerations](#security-considerations)
9. [Error Handling Strategy](#error-handling-strategy)
10. [Testing Strategy](#testing-strategy)

---

## 📊 Executive Summary

### Objective
Menambahkan modul konfigurasi OLT ZTE C320 via Telnet untuk melengkapi fitur monitoring SNMP yang sudah ada.

### Current Capabilities (SNMP - Read Only)
- ✅ ONU Monitoring (status, serial number, model)
- ✅ PON Port Information
- ✅ Traffic Profiles (view)
- ✅ VLAN Profiles (view)
- ✅ Card/Slot Information

### Proposed Capabilities (Telnet - Read/Write)
- 🔄 ONU Registration & Provisioning
- 🔄 ONU VLAN Configuration
- 🔄 Traffic Profile Assignment
- 🔄 Service Port Configuration
- 🔄 ONU Management (reboot, reset, delete)
- 🔄 System Configuration

---

## 🔍 Current State Analysis

### Existing Architecture
```
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer (HTTP)                          │
│                          Port 8081                                │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Handler Layer                              │
│   OnuHandler │ PonHandler │ ProfileHandler │ CardHandler        │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Usecase Layer                              │
│   OnuUsecase │ PonUsecase │ ProfileUsecase │ CardUsecase        │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Repository Layer                            │
│              SnmpRepository │ RedisRepository                    │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                         OLT C320                                 │
│                      SNMP (UDP 161)                              │
└─────────────────────────────────────────────────────────────────┘
```

### Proposed Architecture (with Telnet)
```
┌─────────────────────────────────────────────────────────────────┐
│                        API Layer (HTTP)                          │
│                          Port 8081                                │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Handler Layer                              │
│   OnuHandler │ PonHandler │ ProfileHandler │ CardHandler        │
│   ConfigHandler │ ProvisionHandler │ ServiceHandler             │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Usecase Layer                              │
│   OnuUsecase │ PonUsecase │ ProfileUsecase │ CardUsecase        │
│   ConfigUsecase │ ProvisionUsecase │ ServiceUsecase             │
└─────────────────────────────────────────────────────────────────┘
                                │
                ┌───────────────┴───────────────┐
                ▼                               ▼
┌───────────────────────────┐   ┌───────────────────────────────┐
│     SNMP Repository       │   │      Telnet Repository        │
│    (Read Operations)      │   │   (Write/Config Operations)   │
└───────────────────────────┘   └───────────────────────────────┘
                │                               │
                ▼                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                         OLT C320                                 │
│              SNMP (UDP 161) │ Telnet (TCP 23)                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🗺️ Roadmap & Phases

### Phase 1: Foundation (Week 1-2)
**Goal:** Setup Telnet infrastructure dan basic connection

| Task | Priority | Status | Description |
|------|----------|--------|-------------|
| 1.1 | HIGH | ⬜ | Create Telnet Repository dengan connection pooling |
| 1.2 | HIGH | ⬜ | Implement authentication handler (login/enable) |
| 1.3 | HIGH | ⬜ | Create command executor dengan timeout handling |
| 1.4 | MEDIUM | ⬜ | Setup error detection dari command output |
| 1.5 | MEDIUM | ⬜ | Create session manager untuk concurrent access |
| 1.6 | LOW | ⬜ | Unit tests untuk Telnet repository |

**Deliverables:**
- `internal/repository/telnet.go`
- `internal/repository/telnet_session.go`
- `config/telnet_config.go`

### Phase 2: ONU Provisioning (Week 3-4)
**Goal:** Implement ONU auto-provisioning workflow

| Task | Priority | Status | Description |
|------|----------|--------|-------------|
| 2.1 | HIGH | ⬜ | ONU Registration (onu-type, auth) |
| 2.2 | HIGH | ⬜ | ONU Name Configuration |
| 2.3 | HIGH | ⬜ | Default VLAN Assignment |
| 2.4 | MEDIUM | ⬜ | Traffic Profile Assignment |
| 2.5 | MEDIUM | ⬜ | Service Port Creation |
| 2.6 | LOW | ⬜ | Batch provisioning support |

**Deliverables:**
- `internal/usecase/provision.go`
- `internal/handler/provision.go`
- `internal/model/provision.go`

### Phase 3: VLAN Management (Week 5-6)
**Goal:** Complete VLAN configuration capabilities

| Task | Priority | Status | Description |
|------|----------|--------|-------------|
| 3.1 | HIGH | ⬜ | Create VLAN Profile |
| 3.2 | HIGH | ⬜ | Modify VLAN Profile |
| 3.3 | HIGH | ⬜ | Delete VLAN Profile |
| 3.4 | MEDIUM | ⬜ | ONU VLAN Port Configuration |
| 3.5 | MEDIUM | ⬜ | Service VLAN Assignment |
| 3.6 | LOW | ⬜ | VLAN Statistics |

**Deliverables:**
- `internal/usecase/vlan_config.go`
- `internal/handler/vlan_config.go`

### Phase 4: Traffic Management (Week 7-8)
**Goal:** DBA profile dan bandwidth management

| Task | Priority | Status | Description |
|------|----------|--------|-------------|
| 4.1 | HIGH | ⬜ | Create Traffic Profile |
| 4.2 | HIGH | ⬜ | Modify Traffic Profile |
| 4.3 | HIGH | ⬜ | Delete Traffic Profile |
| 4.4 | MEDIUM | ⬜ | Assign Profile to ONU |
| 4.5 | MEDIUM | ⬜ | DBA Configuration |
| 4.6 | LOW | ⬜ | QoS Settings |

**Deliverables:**
- `internal/usecase/traffic_config.go`
- `internal/handler/traffic_config.go`

### Phase 5: ONU Management (Week 9-10)
**Goal:** ONU lifecycle management

| Task | Priority | Status | Description |
|------|----------|--------|-------------|
| 5.1 | HIGH | ⬜ | Reboot ONU |
| 5.2 | HIGH | ⬜ | Delete ONU |
| 5.3 | MEDIUM | ⬜ | Reset ONU to Factory |
| 5.4 | MEDIUM | ⬜ | Update ONU Description |
| 5.5 | MEDIUM | ⬜ | Block/Unblock ONU |
| 5.6 | LOW | ⬜ | ONU Firmware Upgrade |

**Deliverables:**
- `internal/usecase/onu_management.go`
- `internal/handler/onu_management.go`

### Phase 6: Advanced Features (Week 11-12)
**Goal:** Fitur advanced dan optimasi

| Task | Priority | Status | Description |
|------|----------|--------|-------------|
| 6.1 | MEDIUM | ⬜ | Configuration Backup/Restore |
| 6.2 | MEDIUM | ⬜ | Batch Configuration |
| 6.3 | MEDIUM | ⬜ | Configuration Templates |
| 6.4 | LOW | ⬜ | Audit Logging |
| 6.5 | LOW | ⬜ | Rollback Support |
| 6.6 | LOW | ⬜ | Webhook Notifications |

**Deliverables:**
- `internal/usecase/config_management.go`
- `internal/service/template.go`

---

## 📖 Telnet Command Reference (ZTE C320 V2.1)

### Authentication & Mode Commands
```bash
# Login
Username: <username>
Password: <password>

# Enter Enable Mode
enable
Password: <enable_password>

# Enter Configure Mode
configure terminal

# Exit Configure Mode
exit
end
```

### ONU Registration & Provisioning
```bash
# Masuk ke interface GPON
interface gpon-olt_1/1/1

# Lihat ONU yang belum terdaftar
show gpon onu uncfg

# Registrasi ONU dengan type dan auth
onu <onu_id> type <onu_type> sn <serial_number>

# Contoh registrasi ONU
onu 1 type ZTE-F660 sn ZTEGD824CDF3

# Set ONU name/description
onu <onu_id> name <name>

# Contoh
onu 1 name Pelanggan_001
```

### ONU Type Configuration (per ONU)
```bash
# Masuk ke konfigurasi ONU
interface gpon-onu_1/1/1:1

# Set TCONT
tcont 1 name TCONT_1 profile UP-10M

# Set GEMPORT
gemport 1 name GEM_1 tcont 1

# Set service port (VLAN tagging)
service-port 1 vport 1 user-vlan 100 vlan 100

# Keluar dari interface ONU
exit
```

### VLAN Configuration
```bash
# Masuk configure terminal
configure terminal

# Buat VLAN baru
vlan <vlan_id>

# Set VLAN name
name <vlan_name>

# Contoh lengkap
vlan 100
name INTERNET
exit
```

### PON Port VLAN
```bash
# Masuk ke interface GPON
interface gpon-olt_1/1/1

# Set port VLAN mode
port vlan-mode tag

# Add VLAN to port
port vlan 100

# Native VLAN
port native-vlan 1
```

### Service Port Configuration
```bash
# Global service port configuration
service-port <port_id> gpon 1/1/1 onu <onu_id> gemport <gem_id> cos <cos> vlan <vlan_id>

# Contoh
service-port 1 gpon 1/1/1 onu 1 gemport 1 cos 0 vlan 100

# Dengan user-vlan translation
service-port 1 gpon 1/1/1 onu 1 gemport 1 cos 0 user-vlan 100 vlan 100
```

### Traffic Profile (DBA Profile)
```bash
# Masuk configure terminal
configure terminal

# Buat DBA profile
gpon-onu-profile dba-profile <profile_name>
type 4 assure <cir> max <pir>
exit

# Contoh
gpon-onu-profile dba-profile UP-10M
type 4 assure 10240 max 10240
exit
```

### ONU Management Commands
```bash
# Reboot ONU
interface gpon-olt_1/1/1
onu reset <onu_id>

# Delete ONU
no onu <onu_id>

# Block ONU
onu <onu_id> state disable

# Unblock ONU
onu <onu_id> state enable
```

### Show Commands (Read Only via Telnet)
```bash
# Show registered ONUs
show gpon onu state gpon-olt_1/1/1

# Show unconfigured ONUs
show gpon onu uncfg gpon-olt_1/1/1

# Show ONU detail
show gpon onu detail-info gpon-olt_1/1/1 <onu_id>

# Show service ports
show service-port all

# Show running config
show running-config

# Show traffic profiles
show gpon-onu-profile dba-profile

# Show VLAN config
show vlan all
```

### Configuration Save
```bash
# Save running config to startup
write
# atau
copy running-config startup-config

# Save to file
copy running-config flash:backup.cfg
```

---

## 🔌 API Endpoints Design

### Provisioning Endpoints
```
POST   /api/v1/provision/onu
       Request: { board_id, pon_id, onu_id, onu_type, serial_number, name }
       Response: { success, message, onu_config }

DELETE /api/v1/provision/onu/{board_id}/{pon_id}/{onu_id}
       Response: { success, message }

PUT    /api/v1/provision/onu/{board_id}/{pon_id}/{onu_id}/name
       Request: { name }
       Response: { success, message }
```

### VLAN Configuration Endpoints
```
POST   /api/v1/config/vlan
       Request: { vlan_id, name, description }
       Response: { success, message, vlan }

PUT    /api/v1/config/vlan/{vlan_id}
       Request: { name, description }
       Response: { success, message }

DELETE /api/v1/config/vlan/{vlan_id}
       Response: { success, message }

POST   /api/v1/config/onu/{board_id}/{pon_id}/{onu_id}/vlan
       Request: { vlan_id, mode, gem_port }
       Response: { success, message }
```

### Traffic Profile Endpoints
```
POST   /api/v1/config/traffic-profile
       Request: { name, type, cir, pir, max_bw }
       Response: { success, message, profile }

PUT    /api/v1/config/traffic-profile/{name}
       Request: { cir, pir, max_bw }
       Response: { success, message }

DELETE /api/v1/config/traffic-profile/{name}
       Response: { success, message }

POST   /api/v1/config/onu/{board_id}/{pon_id}/{onu_id}/traffic-profile
       Request: { profile_name, tcont_id }
       Response: { success, message }
```

### Service Port Endpoints
```
POST   /api/v1/config/service-port
       Request: { 
         board_id, pon_id, onu_id, 
         gemport, vlan_id, user_vlan, cos 
       }
       Response: { success, message, service_port_id }

DELETE /api/v1/config/service-port/{service_port_id}
       Response: { success, message }

GET    /api/v1/config/service-port
       Query: ?board_id=1&pon_id=1&onu_id=1
       Response: { service_ports: [...] }
```

### ONU Management Endpoints
```
POST   /api/v1/manage/onu/{board_id}/{pon_id}/{onu_id}/reboot
       Response: { success, message }

POST   /api/v1/manage/onu/{board_id}/{pon_id}/{onu_id}/reset
       Response: { success, message }

PUT    /api/v1/manage/onu/{board_id}/{pon_id}/{onu_id}/state
       Request: { state: "enable" | "disable" }
       Response: { success, message }
```

### System Endpoints
```
POST   /api/v1/system/save-config
       Response: { success, message }

POST   /api/v1/system/backup
       Response: { success, backup_file, timestamp }

POST   /api/v1/system/restore
       Request: { backup_file }
       Response: { success, message }
```

---

## 🔄 Workflow Diagrams

### 1. ONU Auto-Provisioning Workflow
```
┌─────────────────────────────────────────────────────────────────────┐
│                    ONU AUTO-PROVISIONING WORKFLOW                    │
└─────────────────────────────────────────────────────────────────────┘

    ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
    │   API    │────▶│  Usecase │────▶│  Telnet  │────▶│   OLT    │
    │ Request  │     │  Layer   │     │   Repo   │     │  C320    │
    └──────────┘     └──────────┘     └──────────┘     └──────────┘
         │                │                │                │
         │   1. POST      │                │                │
         │   /provision   │                │                │
         │───────────────▶│                │                │
         │                │  2. Validate   │                │
         │                │     Input      │                │
         │                │───────────────▶│                │
         │                │                │  3. Connect    │
         │                │                │───────────────▶│
         │                │                │  4. Login      │
         │                │                │───────────────▶│
         │                │                │  5. Enable     │
         │                │                │───────────────▶│
         │                │                │  6. Configure  │
         │                │                │───────────────▶│
         │                │                │     Terminal   │
         │                │                │───────────────▶│
         │                │                │  7. Interface  │
         │                │                │     gpon-olt   │
         │                │                │───────────────▶│
         │                │                │  8. Register   │
         │                │                │     ONU        │
         │                │                │───────────────▶│
         │                │                │  9. Set Name   │
         │                │                │───────────────▶│
         │                │                │  10. Exit &    │
         │                │                │      Save      │
         │                │                │◀───────────────│
         │                │  11. Parse     │                │
         │                │      Response  │                │
         │                │◀───────────────│                │
         │  12. Return    │                │                │
         │      Result    │                │                │
         │◀───────────────│                │                │
         │                │                │                │
         ▼                ▼                ▼                ▼
```

### 2. Complete ONU Setup Workflow
```
┌─────────────────────────────────────────────────────────────────────┐
│                    COMPLETE ONU SETUP SEQUENCE                       │
└─────────────────────────────────────────────────────────────────────┘

Step 1: ONU Detection (SNMP)
┌─────────────────────────────────────────────────────────┐
│  GET /api/v1/board/1/pon/1/onu_id/empty                │
│  Response: { empty_slots: [1, 2, 3...] }               │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
Step 2: Check Unconfigured ONU (SNMP/Telnet)
┌─────────────────────────────────────────────────────────┐
│  GET /api/v1/board/1/pon/1/unconfigured                │
│  Response: { uncfg_onus: [{sn, type}...] }             │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
Step 3: Register ONU (Telnet)
┌─────────────────────────────────────────────────────────┐
│  POST /api/v1/provision/onu                            │
│  Body: {                                               │
│    board_id: 1,                                        │
│    pon_id: 1,                                          │
│    onu_id: 1,                                          │
│    onu_type: "ZTE-F670L",                              │
│    serial_number: "ZTEGDA5918AC",                      │
│    name: "Customer_001"                                │
│  }                                                     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
Step 4: Configure TCONT & GEMPORT (Telnet)
┌─────────────────────────────────────────────────────────┐
│  POST /api/v1/config/onu/1/1/1/tcont                   │
│  Body: { tcont_id: 1, profile_name: "UP-10M" }         │
│                                                        │
│  POST /api/v1/config/onu/1/1/1/gemport                 │
│  Body: { gemport_id: 1, tcont_id: 1 }                  │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
Step 5: Configure Service Port (Telnet)
┌─────────────────────────────────────────────────────────┐
│  POST /api/v1/config/service-port                      │
│  Body: {                                               │
│    board_id: 1, pon_id: 1, onu_id: 1,                  │
│    gemport: 1, vlan_id: 100, user_vlan: 100, cos: 0    │
│  }                                                     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
Step 6: Save Configuration (Telnet)
┌─────────────────────────────────────────────────────────┐
│  POST /api/v1/system/save-config                       │
│  Response: { success: true }                           │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
Step 7: Verify ONU Status (SNMP)
┌─────────────────────────────────────────────────────────┐
│  GET /api/v1/board/1/pon/1/onu/1                       │
│  Response: { status: "online", ... }                   │
└─────────────────────────────────────────────────────────┘
```

### 3. Error Handling Workflow
```
┌─────────────────────────────────────────────────────────────────────┐
│                      ERROR HANDLING WORKFLOW                         │
└─────────────────────────────────────────────────────────────────────┘

                    ┌────────────────┐
                    │  API Request   │
                    └───────┬────────┘
                            │
                            ▼
                    ┌────────────────┐
                    │   Validation   │
                    └───────┬────────┘
                            │
              ┌─────────────┼─────────────┐
              │             │             │
         ❌ Invalid    ✅ Valid      ⚠️ Warning
              │             │             │
              ▼             ▼             ▼
       ┌──────────┐  ┌──────────┐  ┌──────────┐
       │  Return  │  │  Execute │  │  Log &   │
       │  400     │  │  Command │  │ Continue │
       └──────────┘  └────┬─────┘  └──────────┘
                          │
              ┌───────────┼───────────┐
              │           │           │
         ❌ Timeout   ✅ Success  ❌ Error
              │           │           │
              ▼           ▼           ▼
       ┌──────────┐ ┌──────────┐ ┌──────────┐
       │  Retry   │ │  Parse   │ │  Parse   │
       │  (3x)    │ │  Output  │ │  Error   │
       └────┬─────┘ └────┬─────┘ └────┬─────┘
            │            │            │
            ▼            ▼            ▼
       ┌──────────┐ ┌──────────┐ ┌──────────┐
       │  Return  │ │  Return  │ │  Return  │
       │  504     │ │  200     │ │  500     │
       └──────────┘ └──────────┘ └──────────┘
```

---

## 🔐 Security Considerations

### Authentication
```yaml
Telnet Credentials:
  - Store encrypted in environment variables
  - Support multiple credential levels (admin, operator, viewer)
  - Session timeout after inactivity

API Security:
  - JWT/API Key authentication
  - Role-based access control (RBAC)
  - Rate limiting per endpoint
  - IP whitelist for config endpoints
```

### Configuration
```go
// config/telnet_config.go
type TelnetConfig struct {
    Host            string        `env:"TELNET_HOST"`
    Port            int           `env:"TELNET_PORT" default:"23"`
    Username        string        `env:"TELNET_USER"`
    Password        string        `env:"TELNET_PASS"`      // Encrypted
    EnablePassword  string        `env:"TELNET_ENABLE"`    // Encrypted
    Timeout         time.Duration `env:"TELNET_TIMEOUT" default:"30s"`
    MaxRetries      int           `env:"TELNET_RETRIES" default:"3"`
    SessionPoolSize int           `env:"TELNET_POOL_SIZE" default:"5"`
}
```

### Audit Logging
```yaml
Log Events:
  - All configuration changes
  - User identification
  - Timestamp
  - Before/After values
  - Command executed
  - Result status
```

---

## ⚠️ Error Handling Strategy

### Telnet Error Types
```go
type TelnetError struct {
    Code    string // ERR_CONNECTION, ERR_AUTH, ERR_COMMAND, ERR_TIMEOUT
    Message string
    Command string
    Output  string
}
```

### Error Patterns Detection
```go
var errorPatterns = map[string]string{
    "% Invalid input":           "ERR_INVALID_COMMAND",
    "% Incomplete command":      "ERR_INCOMPLETE_COMMAND",
    "% Ambiguous command":       "ERR_AMBIGUOUS_COMMAND",
    "% Unknown command":         "ERR_UNKNOWN_COMMAND",
    "Error: ONU not exist":      "ERR_ONU_NOT_FOUND",
    "Error: VLAN not exist":     "ERR_VLAN_NOT_FOUND",
    "Error: Profile not exist":  "ERR_PROFILE_NOT_FOUND",
    "Connection refused":        "ERR_CONNECTION_REFUSED",
    "Authentication failed":     "ERR_AUTH_FAILED",
}
```

### Retry Strategy
```go
type RetryConfig struct {
    MaxRetries  int           // Default: 3
    Delay       time.Duration // Default: 1s
    MaxDelay    time.Duration // Default: 10s
    Multiplier  float64       // Default: 2.0 (exponential backoff)
}
```

---

## 🧪 Testing Strategy

### Unit Tests
```
internal/repository/telnet_test.go
- TestConnect
- TestAuthenticate
- TestExecuteCommand
- TestParseOutput
- TestErrorDetection
```

### Integration Tests
```
internal/usecase/provision_test.go
- TestRegisterONU_Success
- TestRegisterONU_DuplicateSN
- TestRegisterONU_InvalidType
- TestDeleteONU_Success
- TestDeleteONU_NotFound
```

### Mock Server
```go
// Create mock Telnet server for testing
type MockTelnetServer struct {
    responses map[string]string
}

func (m *MockTelnetServer) HandleCommand(cmd string) string {
    if resp, ok := m.responses[cmd]; ok {
        return resp
    }
    return "% Unknown command"
}
```

---

## 📁 File Structure (Proposed)

```
go-snmp-olt-zte-c320/
├── cmd/api/
│   └── main.go
├── app/
│   ├── app.go
│   └── routes.go              # Updated with new routes
├── config/
│   ├── config.go
│   ├── oid_generator.go
│   └── telnet_config.go       # NEW
├── internal/
│   ├── model/
│   │   ├── onu.go
│   │   ├── pon.go
│   │   ├── profile.go
│   │   ├── card.go
│   │   ├── provision.go       # NEW: Provisioning models
│   │   ├── service_port.go    # NEW: Service port models
│   │   └── config_request.go  # NEW: Config request models
│   ├── usecase/
│   │   ├── onu.go
│   │   ├── pon.go
│   │   ├── profile.go
│   │   ├── card.go
│   │   ├── provision.go       # NEW: ONU provisioning logic
│   │   ├── vlan_config.go     # NEW: VLAN config logic
│   │   ├── traffic_config.go  # NEW: Traffic profile logic
│   │   └── onu_management.go  # NEW: ONU management logic
│   ├── handler/
│   │   ├── onu.go
│   │   ├── pon.go
│   │   ├── profile.go
│   │   ├── card.go
│   │   ├── provision.go       # NEW
│   │   ├── vlan_config.go     # NEW
│   │   ├── traffic_config.go  # NEW
│   │   └── onu_management.go  # NEW
│   ├── repository/
│   │   ├── snmp.go
│   │   ├── redis.go
│   │   ├── telnet.go          # NEW: Telnet repository
│   │   └── telnet_session.go  # NEW: Session management
│   ├── middleware/
│   ├── utils/
│   │   └── telnet_parser.go   # NEW: Telnet output parser
│   └── errors/
│       └── telnet_errors.go   # NEW: Telnet-specific errors
├── pkg/
│   └── telnet/                # NEW: Telnet client library
│       ├── client.go
│       ├── pool.go
│       └── parser.go
├── docs/
│   ├── TELNET_CONFIG_ROADMAP.md   # This file
│   ├── API_REFERENCE.md
│   └── COMMAND_REFERENCE.md
└── templates/                 # NEW: Config templates
    ├── onu_basic.tmpl
    ├── onu_internet.tmpl
    ├── onu_voip.tmpl
    └── onu_iptv.tmpl
```

---

## 🎯 Priority Matrix

| Feature | Business Value | Technical Complexity | Priority |
|---------|---------------|---------------------|----------|
| ONU Registration | HIGH | LOW | P1 |
| Service Port Config | HIGH | MEDIUM | P1 |
| VLAN Assignment | HIGH | LOW | P1 |
| Traffic Profile Assign | MEDIUM | LOW | P2 |
| ONU Reboot | MEDIUM | LOW | P2 |
| ONU Delete | MEDIUM | LOW | P2 |
| Create Traffic Profile | LOW | MEDIUM | P3 |
| Create VLAN | LOW | MEDIUM | P3 |
| Config Backup | LOW | HIGH | P4 |
| Batch Operations | LOW | HIGH | P4 |

---

## 📊 Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Provisioning Time | < 10s | API response time |
| Command Success Rate | > 99% | Success/Total commands |
| Connection Uptime | > 99.5% | Session availability |
| Error Recovery | < 3 retries | Average retries per error |
| Config Save Time | < 5s | Write command response |

---

## 🔗 Dependencies

### Go Libraries (to add)
```go
// go.mod additions
require (
    github.com/ziutek/telnet v0.0.0-20180329124119-c3b780dc415b
    // or
    github.com/reiver/go-telnet v0.0.0-20180421082511-9ff0b2ab096e
)
```

### External Requirements
- OLT accessible via Telnet (port 23)
- Valid admin credentials
- Network connectivity from VPS to OLT

---

## 📝 Notes & Considerations

1. **Concurrent Access**: Telnet connections are stateful. Need connection pooling to handle multiple simultaneous requests.

2. **Command Ordering**: Some commands depend on others (e.g., must create TCONT before GEMPORT). Usecase layer must enforce ordering.

3. **Transaction Support**: OLT doesn't support transactions. Implement compensation logic for rollback.

4. **Rate Limiting**: Don't overload OLT with too many commands. Implement command queue if needed.

5. **Firmware Compatibility**: Command syntax may vary between firmware versions. Abstract command generation.

---

**Next Steps:**
1. Review and approve roadmap
2. Start Phase 1: Telnet Repository implementation
3. Setup development environment with test OLT
4. Create API contract documentation

**Questions/Clarifications Needed:**
- [ ] Telnet credentials for OLT
- [ ] Specific ONU types to support
- [ ] Default VLAN configuration requirements
- [ ] Traffic profile templates
- [ ] Priority for batch operations
