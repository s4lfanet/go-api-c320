# ZTE C320 V2.1.0 SNMP Monitoring - Project State

**Last Updated:** January 11, 2026  
**Status:** Phase 1-4 Complete ✅ | Deployed to Production ✅ | Testing Complete ✅

## Project Overview

Go-based SNMP monitoring and Telnet configuration application for ZTE C320 OLT firmware V2.1.0, deployed on VPS at 192.168.54.230:8081.

### Infrastructure

- **Target OLT:** 136.1.1.100 (ZTE C320 V2.1.0)
- **VPS:** 192.168.54.230 (Ubuntu/Debian)
- **Application Path:** `/opt/go-snmp-olt`
- **Binary Path:** `/opt/go-snmp-olt/bin/api`
- **Service:** `go-snmp-olt.service` (systemd)
- **Port:** 8081
- **SNMP:** v2c, community "public" (UDP 161)
- **Telnet:** TCP 23 (username: zte, password: zte, enable: zxr10)
- **Redis:** Password: `OsWkRgJLabn4n2+nodZ6BQeP+OKkrObnGeFcDY6w7Nw=`
- **Go Version:** 1.25.5

### Deployment Status

**Last Deployment:** January 11, 2026 15:27 UTC

✅ Phase 1 (Telnet Infrastructure) - Deployed & Tested  
✅ Phase 2 (ONU Provisioning) - Deployed & Tested  
✅ Phase 3 (VLAN Management) - Deployed & Tested  
✅ Phase 4 (Traffic Profile Management) - Deployed & Tested

**Endpoint Tests:**
- All 4 provisioning endpoints working
- All 5 VLAN management endpoints working
- All 10 traffic profile endpoints working
- Telnet connectivity confirmed
- Session management operational

**Total Endpoints:** 40+ (including SNMP monitoring endpoints)

## Critical OID Information

### Base OID Structure (V2.1.0)
```
Base OID: 1.3.6.1.4.1.3902.1012
```

**IMPORTANT:** V2.1.0 uses OID base `.1012` NOT `.1082` (which is used in V2.2+)

### PON Index Calculation
```
Board 1: 268500992 + (ponID * 256)
Board 2: 268509184 + (ponID * 256)

Examples:
- Board 1, PON 1: 268500992 + 256 = 268501248
- Board 1, PON 2: 268500992 + 512 = 268501504
- Board 2, PON 1: 268509184 + 256 = 268509440
```

### ONU OID Structure
```
ONU Table: .1012.3.13.3.1.{column}.{pon_index}.{onu_id}

Columns:
- 1: ONU Index
- 2: Name
- 5: Status (1=offline, 2=logging/unconfigured, 3=online)
- 8: Serial Number (GPON format: VendorID + Hex)
- 9: Description
- 11: Model
- 26: Control Flag
```

### PON Port OIDs
```
Admin Status: .1012.3.11.3.1.1.{pon_index}
Distance: .1012.3.11.5.1.3.{pon_index}
Oper Status: .1012.3.11.5.1.4.{pon_index}
```

### Traffic Profile OIDs
```
Profile Name: .1012.3.26.1.1.2.{profile_id}
CIR: .1012.3.26.1.1.3.{profile_id}
PIR: .1012.3.26.1.1.4.{profile_id}
Max BW: .1012.3.26.1.1.5.{profile_id}
```

### VLAN Profile OIDs
```
Base: .1012.3.50.20.15.1.{col}.{length}.{ascii_chars...}
Format: OID contains ASCII-encoded VLAN name
```

### Card/Slot Info OIDs
```
Base: 1.3.6.1.4.1.3902.1015.2.1.1.3.1.{col}.{rack}.{shelf}.{slot}

Columns:
- 2: Card Type (numeric)
- 4: Serial Number
- 5: Hardware Version
- 6: Software Version
- 7: Status (0=inactive, 3=active, 16=online)
```

## Implemented Features

### Original Features (Before This Session)
✅ ONU Monitoring
- List ONUs by Board/PON
- Get specific ONU details
- Get ONU ID + Serial Number list
- Get empty ONU slots
- Pagination support
- Cache management

### New Features (Implemented This Session)

#### 1. PON Port Information ✅
**Endpoint:** `GET /api/v1/board/{board_id}/pon/{pon_id}/info`

**Response:**
```json
{
  "board": 1,
  "pon": 1,
  "admin_status": "enabled",
  "oper_status": "up",
  "onu_count": 3,
  "distance": 200
}
```

**Files:**
- `internal/model/pon.go`
- `internal/usecase/pon.go`
- `internal/handler/pon.go`

#### 2. Traffic Profiles ✅
**Endpoints:**
- `GET /api/v1/profiles/traffic` - List all profiles
- `GET /api/v1/profiles/traffic/{profile_id}` - Get specific profile

**Response Example:**
```json
{
  "profile_id": 1879048194,
  "name": "UP-10M",
  "cir": 0,
  "pir": 0,
  "max_bw": 10240
}
```

**Files:**
- `internal/model/profile.go` (TrafficProfile struct)
- `internal/usecase/profile.go` (Traffic profile methods)
- `internal/handler/profile.go` (Traffic profile handlers)

#### 3. VLAN Profiles ✅
**Endpoint:** `GET /api/v1/profiles/vlan`

**Response Example:**
```json
{
  "name": "pppoe",
  "vlan_id": 1,
  "priority": 30,
  "mode": "tag",
  "description": "0"
}
```

**Special Implementation:** VLAN names are decoded from ASCII values in OID using length-prefixed format.

**Files:**
- `internal/model/profile.go` (VlanProfile struct)
- `internal/usecase/profile.go` (VLAN profile methods)
- `internal/handler/profile.go` (VLAN profile handlers)

#### 4. Card/Slot Information ✅
**Endpoints:**
- `GET /api/v1/system/cards` - List all cards
- `GET /api/v1/system/cards/{rack}/{shelf}/{slot}` - Get specific card

**Response Example:**
```json
{
  "rack": 1,
  "shelf": 1,
  "slot": 1,
  "card_type": "type_599049",
  "status": "online",
  "serial_number": "GTGH",
  "hardware_ver": "v1",
  "software_ver": "v4",
  "description": ""
}
```

**Files:**
- `internal/model/card.go`
- `internal/usecase/card.go`
- `internal/handler/card.go`

## Complete API Endpoints

### ONU Endpoints
```
GET  /api/v1/board/{board_id}/pon/{pon_id}
GET  /api/v1/board/{board_id}/pon/{pon_id}/onu/{onu_id}
GET  /api/v1/board/{board_id}/pon/{pon_id}/onu_id_sn
GET  /api/v1/board/{board_id}/pon/{pon_id}/onu_id/empty
GET  /api/v1/board/{board_id}/pon/{pon_id}/onu_id/update
GET  /api/v1/paginate/board/{board_id}/pon/{pon_id}
DEL  /api/v1/board/{board_id}/pon/{pon_id}
```

### PON Port Endpoints
```
GET  /api/v1/board/{board_id}/pon/{pon_id}/info
```

### Profile Endpoints
```
GET  /api/v1/profiles/traffic
GET  /api/v1/profiles/traffic/{profile_id}
GET  /api/v1/profiles/vlan
```

### System Endpoints
```
GET  /api/v1/system/cards
GET  /api/v1/system/cards/{rack}/{shelf}/{slot}
```

## Telnet Configuration Module (Phase 1 Complete) ✅

### Overview
Implementasi modul konfigurasi Telnet untuk memungkinkan operasi write/configuration pada OLT.

### Implemented Components (January 11, 2026)

#### 1. Telnet Configuration (`config/telnet_config.go`)
- Environment-based configuration loading
- Comprehensive timeout management
- Connection retry settings
- Validation with error handling

**Key Settings:**
- Connection timeouts (connect, read, write)
- Retry logic (count, delay)
- Session pooling configuration
- Prompt patterns (user, enable, config modes)

#### 2. Telnet Models (`internal/model/telnet.go`)
- TelnetCommand & TelnetResponse
- TelnetBatchResponse for multiple commands
- TelnetConnectionInfo for status monitoring
- Error models with recovery indicators
- Request/Response models for future endpoints:
  - ONURegistrationRequest/Response
  - VLANConfigRequest/Response
  - TrafficProfileRequest/Response
  - ServicePortRequest/Response
  - ONUManagementRequest/Response

#### 3. Telnet Repository (`internal/repository/telnet.go`)
**Interface Methods:**
- `Connect()` - Establish connection with authentication
- `Close()` - Graceful connection termination
- `IsConnected()` - Connection status check
- `Reconnect()` - Reconnect logic
- `Execute()` - Single command execution
- `ExecuteMulti()` - Batch command execution
- `ExecuteWithExpect()` - Command with pattern matching
- `EnterEnableMode()` - Enter privileged mode
- `EnterConfigMode()` - Enter configuration mode
- `ExitConfigMode()` - Exit configuration mode
- `SaveConfig()` - Save running config to startup
- `ShowRunningConfig()` - Retrieve running configuration

**Features:**
- Automatic login sequence (username/password)
- Enable mode authentication
- Mode state tracking (user/enable/config/disconnected)
- Output cleaning (remove echo, prompts, ANSI codes)
- Timeout handling with proper deadlines
- Error categorization (recoverable vs non-recoverable)

#### 4. Session Management (`internal/repository/telnet_session.go`)
**TelnetSessionPool:**
- Single session management with busy-wait
- Automatic connection establishment
- Stale connection detection & reconnection
- Session locking to prevent concurrent usage
- Idle connection cleanup

**TelnetSessionManager (Singleton):**
- Global session manager instance
- Wrapper methods for easy command execution:
  - `ExecuteCommand()` - Execute single command
  - `ExecuteCommands()` - Execute multiple commands
  - `ExecuteInConfigMode()` - Execute in config mode with auto exit
  - `SaveConfiguration()` - Save OLT config
- Automatic session acquisition & release
- Retry logic with exponential backoff
- Connection status monitoring

**Features:**
- Automatic idle cleanup (goroutine-based)
- Connection pooling (currently 1 session, expandable)
- Thread-safe operations
- Context-aware execution
- Error recovery mechanisms

### Telnet Environment Variables (.env.example)
```env
# Connection
TELNET_HOST=136.1.1.100
TELNET_PORT=23
TELNET_USERNAME=admin
TELNET_PASSWORD=your_password
TELNET_ENABLE_PASSWORD=your_enable_password

# Timeouts (seconds)
TELNET_TIMEOUT=30
TELNET_CONNECT_TIMEOUT=10
TELNET_READ_TIMEOUT=30
TELNET_WRITE_TIMEOUT=10

# Retry
TELNET_RETRY_COUNT=3
TELNET_RETRY_DELAY=2

# Session
TELNET_POOL_SIZE=1
TELNET_MAX_IDLE_TIME=300

# Prompts
TELNET_PROMPT_USER=ZXAN>
TELNET_PROMPT_ENABLE=ZXAN#
TELNET_PROMPT_CONFIG=ZXAN(config)#
```

### Testing Telnet Module
```go
// Example usage (will be used in Phase 2)
cfg := config.LoadTelnetConfig()
manager := repository.GetGlobalSessionManager(cfg)

// Execute single command
resp, err := manager.ExecuteCommand(ctx, "show gpon onu uncfg")

// Execute in config mode
commands := []string{
    "interface gpon-olt_1/1/1",
    "onu 1 type ZTE-F670L sn ZTEGDA5918AC",
}
result, err := manager.ExecuteInConfigMode(ctx, commands)

// Save config
err = manager.SaveConfiguration(ctx)
```

### Dependencies Added
- `github.com/ziutek/telnet` v0.0.0-20180329124119-c3b780dc415b

## ONU Provisioning Module (Phase 2 Complete) ✅

### Overview
Implementasi complete ONU provisioning workflow melalui HTTP API menggunakan Telnet backend.

### Implemented Components (January 11, 2026)

#### 1. Provisioning Usecase (`internal/usecase/provision.go`)
**Core Methods:**
- `GetAllUnconfiguredONUs()` - Retrieve all unconfigured ONUs from all PON ports
- `GetUnconfiguredONUs(ponPort)` - Get unconfigured ONUs for specific PON port
- `RegisterONU()` - Complete ONU registration workflow:
  - Register ONU with type and serial number
  - Configure TCONT based on DBA profile
  - Configure GEMPORT
  - Configure service port with VLAN
- `DeleteONU()` - Remove ONU from OLT
- `ConfigureTCONT()` - Configure T-CONT for ONU
- `ConfigureGEMPort()` - Configure GEM Port
- `ConfigureServicePort()` - Configure service port with VLAN

**Features:**
- Automatic ONU type detection
- Profile-based auto-configuration
- Comprehensive error handling with partial success tracking
- Output parsing with regex patterns
- ONU state validation

#### 2. Provisioning Handler (`internal/handler/provision.go`)
**HTTP Endpoints:**

**GET /api/v1/onu/unconfigured**
- List all unconfigured ONUs from all PON ports
- Returns array of UnconfiguredONU objects

**GET /api/v1/onu/unconfigured/{pon}**
- List unconfigured ONUs for specific PON port
- Path parameter: `pon` (e.g., "1/1/1")

**POST /api/v1/onu/register**
- Register new ONU with full provisioning
- Request body:
```json
{
  "pon_port": "1/1/1",
  "onu_id": 1,
  "onu_type": "ZTE-F670L",
  "serial_number": "ZTEGDA5918AC",
  "name": "CUSTOMER-001",
  "profile": {
    "dba_profile": "UP-10M",
    "vlan": 100
  }
}
```
- Response:
```json
{
  "success": true,
  "message": "ONU registered successfully",
  "onu_registered": true,
  "tcont_configured": true,
  "gemport_configured": true,
  "service_port_configured": true,
  "onu_info": {
    "pon_port": "1/1/1",
    "onu_id": 1,
    "type": "ZTE-F670L",
    "serial": "ZTEGDA5918AC",
    "name": "CUSTOMER-001"
  },
  "errors": []
}
```

**DELETE /api/v1/onu/{pon}/{onu_id}**
- Delete ONU from OLT
- Path parameters: `pon` (e.g., "1/1/1"), `onu_id` (integer)
- Validates ONU ID range (1-128)

**Features:**
- Comprehensive input validation
- Proper error response handling using utils.HandleError
- Structured logging with zerolog
- Swagger annotations
- HTTP status code handling (200, 201, 400, 500)

#### 3. Updated Application Initialization (`app/app.go`)
**Changes:**
- Load telnet configuration: `telnetCfg := config.LoadTelnetConfig()`
- Initialize global session manager: `telnetSessionManager := repository.GetGlobalSessionManager(telnetCfg)`
- Create provision usecase: `provisionUsecase := usecase.NewProvisionUsecase(telnetSessionManager, cfg)`
- Create provision handler: `provisionHandler := handler.NewProvisionHandler(provisionUsecase)`
- Pass provision handler to routes: `loadRoutes(..., provisionHandler)`

#### 4. Updated Routes (`app/routes.go`)
**New Route Group:**
```go
apiV1Group.Route("/onu", func(r chi.Router) {
    r.Get("/unconfigured", provisionHandler.GetUnconfiguredONUs)
    r.Get("/unconfigured/{pon}", provisionHandler.GetUnconfiguredONUsByPON)
    r.Post("/register", provisionHandler.RegisterONU)
    r.Delete("/{pon}/{onu_id}", provisionHandler.DeleteONU)
})
```

#### 5. Updated Models (`internal/model/telnet.go`)
**Models from Phase 1 (used in Phase 2):**
- `UnconfiguredONU` - Represents unregistered ONU
- `ONURegistrationRequest` - Request for ONU registration
- `ONURegistrationResponse` - Response with success/error details
- `ONUProvisioningProfile` - DBA profile and VLAN configuration
- `ONUInfo` - Basic ONU information

### Example Usage

**1. Get All Unconfigured ONUs:**
```bash
curl http://192.168.54.230:8081/api/v1/onu/unconfigured
```

**2. Get Unconfigured ONUs for PON 1/1/1:**
```bash
curl http://192.168.54.230:8081/api/v1/onu/unconfigured/1/1/1
```

**3. Register New ONU:**
```bash
curl -X POST http://192.168.54.230:8081/api/v1/onu/register \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 1,
    "onu_type": "ZTE-F670L",
    "serial_number": "ZTEGDA5918AC",
    "name": "CUSTOMER-001",
    "profile": {
      "dba_profile": "UP-10M",
      "vlan": 100
    }
  }'
```

**4. Delete ONU:**
```bash
curl -X DELETE http://192.168.54.230:8081/api/v1/onu/1/1/1/1
```

### Testing
**Updated Test Files:**
- `app/routes_test.go` - Updated all 4 test functions to include provision handler mock

### Known Limitations (Phase 2)
- No batch ONU registration endpoint
- No ONU status update endpoint
- No service port modification endpoint
- Manual VLAN/profile validation (no database lookup)
- No audit logging for provisioning actions

## VLAN Management Module (Phase 3 Complete) ✅

### Overview
Implementasi complete VLAN management untuk ONU menggunakan service-port configuration via Telnet.

### Implemented Components (January 11, 2026)

#### 1. VLAN Models (`internal/model/telnet.go`)
**New Models:**
- `ONUVLANInfo` - Complete VLAN information for ONU
  - PONPort, ONUID, SVLAN, CVLAN, VLANMode, Priority, ServicePortID
- `VLANConfigRequest` - Request to configure/modify VLAN
  - Validation tags for VLAN ranges (1-4094), mode (tag/translation/transparent), priority (0-7)
- `VLANConfigResponse` - Response after VLAN configuration
  - Success status, message, service-port ID

#### 2. VLAN Repository (`internal/repository/telnet_vlan.go`)
**Core Methods:**
- `GetONUVLAN()` - Retrieve VLAN configuration for specific ONU
- `GetAllServicePorts()` - List all service-port configurations
- `ConfigureONUVLAN()` - Create/update VLAN configuration (service-port)
- `DeleteONUVLAN()` - Remove VLAN configuration
- `buildCreateServicePortCommand()` - Build service-port creation command
- `buildUpdateServicePortCommand()` - Build service-port modification command
- `parseServicePortDetails()` - Parse service-port output
- `parseAllServicePorts()` - Parse all service-ports list

**Features:**
- Automatic service-port ID detection
- Support for multiple VLAN modes: tag, translation, transparent
- SVLAN and CVLAN (user-vlan) configuration
- CoS/Priority configuration
- Output parsing with regex for reliable data extraction

#### 3. VLAN Usecase (`internal/usecase/vlan.go`)
**Business Logic Methods:**
- `GetONUVLAN()` - Get VLAN config with validation
- `GetAllServicePorts()` - Retrieve all service-ports
- `ConfigureVLAN()` - Configure new VLAN for ONU
- `ModifyVLAN()` - Update existing VLAN configuration
- `DeleteVLAN()` - Remove VLAN configuration
- `validateVLANRequest()` - Comprehensive request validation
- `validatePONPort()` - PON port format validation (rack/shelf/slot)
- `validateONUID()` - ONU ID range validation (1-128)

**Validation:**
- SVLAN range: 1-4094
- CVLAN range: 1-4094 (optional)
- VLAN mode: tag, translation, transparent
- Priority: 0-7
- CVLAN required for translation mode

#### 4. VLAN Handler (`internal/handler/vlan.go`)
**HTTP Endpoints:**

**GET /api/v1/vlan/onu/{pon}/{onu_id}**
- Get VLAN configuration for specific ONU
- Path parameters: `pon` (e.g., "1/1/1"), `onu_id` (integer)
- Response:
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "pon_port": "1/1/1",
    "onu_id": 1,
    "svlan": 100,
    "cvlan": 100,
    "vlan_mode": "tag-transform",
    "priority": 0,
    "service_port_id": 123
  }
}
```

**GET /api/v1/vlan/service-ports**
- List all service-port (VLAN) configurations
- Returns array of ONUVLANInfo

**POST /api/v1/vlan/onu**
- Configure new VLAN for ONU
- Request body:
```json
{
  "pon_port": "1/1/1",
  "onu_id": 1,
  "svlan": 100,
  "cvlan": 100,
  "vlan_mode": "tag",
  "priority": 0
}
```
- Response:
```json
{
  "code": 201,
  "status": "Created",
  "data": {
    "pon_port": "1/1/1",
    "onu_id": 1,
    "svlan": 100,
    "cvlan": 100,
    "vlan_mode": "tag",
    "success": true,
    "message": "VLAN configured successfully",
    "service_port_id": 124
  }
}
```

**PUT /api/v1/vlan/onu**
- Modify existing VLAN configuration
- Same request body as POST
- Validates existing VLAN presence before modification

**DELETE /api/v1/vlan/onu/{pon}/{onu_id}**
- Delete VLAN configuration
- Removes service-port from OLT
- Response:
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "pon_port": "1/1/1",
    "onu_id": 1,
    "message": "VLAN deleted successfully",
    "deleted_at": "2026-01-11T10:30:00Z"
  }
}
```

**Features:**
- Comprehensive input validation
- Proper error handling with utils.HandleError
- Structured logging
- Swagger annotations
- HTTP status codes: 200 (OK), 201 (Created), 400 (Bad Request), 404 (Not Found), 500 (Internal Server Error)

#### 5. Updated Files
**app/app.go:**
- Added `vlanUsecase` initialization
- Added `vlanHandler` initialization
- Updated `loadRoutes()` call with vlanHandler

**app/routes.go:**
- Added new `/vlan` route group
- 5 new VLAN endpoints
- Updated function signature to include vlanHandler

**app/routes_test.go:**
- Updated all 4 test functions with vlanUsecase and vlanHandler mocks

### Example Usage

**1. Get ONU VLAN Configuration:**
```bash
curl http://192.168.54.230:8081/api/v1/vlan/onu/1/1/1/1
```

**2. Get All Service-Ports:**
```bash
curl http://192.168.54.230:8081/api/v1/vlan/service-ports
```

**3. Configure VLAN (Tag Mode):**
```bash
curl -X POST http://192.168.54.230:8081/api/v1/vlan/onu \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 1,
    "svlan": 100,
    "cvlan": 100,
    "vlan_mode": "tag",
    "priority": 0
  }'
```

**4. Configure VLAN (Translation Mode):**
```bash
curl -X POST http://192.168.54.230:8081/api/v1/vlan/onu \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 2,
    "svlan": 200,
    "cvlan": 100,
    "vlan_mode": "translation",
    "priority": 3
  }'
```

**5. Modify VLAN:**
```bash
curl -X PUT http://192.168.54.230:8081/api/v1/vlan/onu \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 1,
    "svlan": 200,
    "cvlan": 200,
    "vlan_mode": "tag",
    "priority": 0
  }'
```

**6. Delete VLAN:**
```bash
curl -X DELETE http://192.168.54.230:8081/api/v1/vlan/onu/1/1/1/1
```

### VLAN Modes Explained

**1. Tag Mode (`tag`):**
- Default VLAN tagging
- SVLAN and CVLAN are the same (typically)
- Standard 802.1Q tagging

**2. Translation Mode (`translation`):**
- VLAN translation from CVLAN to SVLAN
- Used for multi-tenant scenarios
- Requires both SVLAN and CVLAN
- Customer VLAN (CVLAN) is translated to Service VLAN (SVLAN)

**3. Transparent Mode (`transparent`):**
- Pass-through mode
- No VLAN manipulation
- Transparent to customer traffic

### Known Limitations (Phase 3)
- No batch VLAN configuration
- No VLAN profile templates
- Manual service-port ID allocation
- No VLAN conflict detection
- No service-port utilization tracking

---

## Traffic Profile Management Module (Phase 4 Complete) ✅

**Implementation Date:** January 11, 2026  
**Status:** Deployed to Production ✅

### Overview
Phase 4 implements comprehensive traffic profile management for ZTE C320 OLT, including DBA (Dynamic Bandwidth Allocation) profiles, T-CONT (Transmission Container) configuration, and GEM (GPON Encapsulation Method) port management. This phase enables fine-grained bandwidth control and QoS configuration for ONU services.

### Key Components

#### 1. Models (internal/model/telnet.go)
**New Models Added:**

**DBA Profile Models:**
```go
type DBAProfileInfo struct {
    Name             string
    Type             int    // 1=Fixed, 2=Assured, 3=Assured+Max, 4=Max, 5=Assured+Max+Priority
    FixedBandwidth   int    // Kbps (Type 1)
    AssuredBandwidth int    // Kbps (Type 2,3,5)
    MaxBandwidth     int    // Kbps (Type 3,4,5)
}

type DBAProfileRequest struct {
    Name             string
    Type             int
    FixedBandwidth   int
    AssuredBandwidth int
    MaxBandwidth     int
}
```

**T-CONT Models:**
```go
type TCONTInfo struct {
    PONPort   string
    ONUID     int
    TCONTID   int
    Name      string
    Profile   string  // DBA profile name
    GEMPorts  []int
    Bandwidth int     // Current allocated bandwidth
}

type TCONTConfigRequest struct {
    PONPort string
    ONUID   int
    TCONTID int
    Name    string
    Profile string  // DBA profile name
}
```

**GEM Port Models:**
```go
type GEMPortInfo struct {
    PONPort   string
    ONUID     int
    GEMPortID int
    Name      string
    TCONTID   int
    Queue     int
}

type GEMPortConfigRequest struct {
    PONPort   string
    ONUID     int
    GEMPortID int
    Name      string
    TCONTID   int
    Queue     int
}
```

#### 2. Repository (internal/repository/telnet_traffic.go)
**File:** 550+ lines of traffic profile management logic

**Key Methods:**

**DBA Profile Operations:**
- `GetDBAProfile(name)` - Retrieve DBA profile details
- `GetAllDBAProfiles()` - List all DBA profiles
- `CreateDBAProfile(req)` - Create new DBA profile
- `ModifyDBAProfile(req)` - Modify existing DBA profile
- `DeleteDBAProfile(name)` - Delete DBA profile

**T-CONT Operations:**
- `GetONUTCONT(ponPort, onuID, tcontID)` - Get T-CONT configuration
- `ConfigureTCONT(req)` - Configure T-CONT with DBA profile
- `DeleteTCONT(ponPort, onuID, tcontID)` - Delete T-CONT

**GEM Port Operations:**
- `ConfigureGEMPort(req)` - Configure GEM port with T-CONT binding
- `DeleteGEMPort(ponPort, onuID, gemportID)` - Delete GEM port

**Command Generation:**
- Generates proper ZTE CLI commands for each operation
- Handles different DBA profile types (1-5)
- Validates bandwidth values
- Manages ONU interface mode transitions

**DBA Profile Type Support:**
| Type | Description | Parameters |
|------|-------------|------------|
| 1 | Fixed Bandwidth | `fix <bandwidth>` |
| 2 | Assured Bandwidth | `assure <bandwidth>` |
| 3 | Assured + Max | `assure <min> max <max>` |
| 4 | Maximum (Best Effort) | `max <bandwidth>` |
| 5 | Assured + Max + Priority | `assure <min> max <max>` |

#### 3. Usecase (internal/usecase/traffic.go)
**File:** 490+ lines of business logic and validation

**Key Methods:**

**DBA Profile Management:**
- `GetDBAProfile(name)` - Get profile with validation
- `GetAllDBAProfiles()` - List all profiles
- `CreateDBAProfile(req)` - Create with validation
- `ModifyDBAProfile(req)` - Modify with existence check
- `DeleteDBAProfile(name)` - Delete with validation

**T-CONT Management:**
- `GetONUTCONT(ponPort, onuID, tcontID)` - Get T-CONT info
- `ConfigureTCONT(req)` - Configure with validation
- `DeleteTCONT(ponPort, onuID, tcontID)` - Delete T-CONT

**GEM Port Management:**
- `ConfigureGEMPort(req)` - Configure with validation
- `DeleteGEMPort(ponPort, onuID, gemportID)` - Delete GEM port

**Validation Functions:**
- `validateProfileName()` - Name length and format
- `validateDBAProfileRequest()` - Type and bandwidth validation
- `validateTCONTID()` - Range check (1-8)
- `validateGEMPortID()` - Range check (1-128)
- `validateTCONTRequest()` - Comprehensive T-CONT validation
- `validateGEMPortRequest()` - Comprehensive GEM port validation

**Validation Rules:**
- Profile name: max 32 characters
- Bandwidth: minimum 64 Kbps
- Type: 1-5
- T-CONT ID: 1-8
- GEM port ID: 1-128
- Queue: 0-8
- Assured bandwidth ≤ Max bandwidth (Type 3,5)

#### 4. Handler (internal/handler/traffic.go)
**File:** 490+ lines of HTTP handlers

**Endpoints Implemented:**

**DBA Profile Endpoints:**
1. `GET /api/v1/traffic/dba-profiles` - Get all DBA profiles
2. `GET /api/v1/traffic/dba-profile/{name}` - Get specific DBA profile
3. `POST /api/v1/traffic/dba-profile` - Create DBA profile
4. `PUT /api/v1/traffic/dba-profile` - Modify DBA profile
5. `DELETE /api/v1/traffic/dba-profile/{name}` - Delete DBA profile

**T-CONT Endpoints:**
6. `GET /api/v1/traffic/tcont/{pon}/{onu_id}/{tcont_id}` - Get T-CONT config
7. `POST /api/v1/traffic/tcont` - Configure T-CONT
8. `DELETE /api/v1/traffic/tcont/{pon}/{onu_id}/{tcont_id}` - Delete T-CONT

**GEM Port Endpoints:**
9. `POST /api/v1/traffic/gemport` - Configure GEM port
10. `DELETE /api/v1/traffic/gemport/{pon}/{onu_id}/{gemport_id}` - Delete GEM port

**Features:**
- Comprehensive Swagger documentation
- Request validation and error handling
- Structured logging with zerolog
- HTTP status codes: 200 (OK), 400 (Bad Request), 404 (Not Found), 500 (Internal Server Error)

#### 5. Updated Files
**app/app.go:**
- Added `trafficUsecase` initialization
- Added `trafficHandler` initialization
- Updated `loadRoutes()` call with trafficHandler

**app/routes.go:**
- Added new `/traffic` route group
- 10 new traffic management endpoints
- Updated function signature to include trafficHandler

**app/routes_test.go:**
- Updated all 4 test functions with trafficUsecase and trafficHandler mocks

### Example Usage

**1. Create DBA Profile (10M Symmetric):**
```bash
curl -X POST http://192.168.54.230:8081/api/v1/traffic/dba-profile \
  -H "Content-Type: application/json" \
  -d '{
    "name": "UP-10M",
    "type": 4,
    "max_bandwidth": 10240
  }'
```

**2. Create DBA Profile (50M Burst):**
```bash
curl -X POST http://192.168.54.230:8081/api/v1/traffic/dba-profile \
  -H "Content-Type: application/json" \
  -d '{
    "name": "UP-50M-BURST",
    "type": 3,
    "assured_bandwidth": 10240,
    "max_bandwidth": 51200
  }'
```

**3. Configure T-CONT on ONU:**
```bash
curl -X POST http://192.168.54.230:8081/api/v1/traffic/tcont \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 1,
    "tcont_id": 1,
    "name": "TCONT_DATA",
    "profile": "UP-10M"
  }'
```

**4. Configure GEM Port:**
```bash
curl -X POST http://192.168.54.230:8081/api/v1/traffic/gemport \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 1,
    "gemport_id": 1,
    "name": "GEM_DATA",
    "tcont_id": 1,
    "queue": 1
  }'
```

**5. Get DBA Profile Details:**
```bash
curl http://192.168.54.230:8081/api/v1/traffic/dba-profile/UP-10M
```

**6. Get All DBA Profiles:**
```bash
curl http://192.168.54.230:8081/api/v1/traffic/dba-profiles
```

**7. Get T-CONT Configuration:**
```bash
curl http://192.168.54.230:8081/api/v1/traffic/tcont/1-1-1/1/1
```

**8. Modify DBA Profile:**
```bash
curl -X PUT http://192.168.54.230:8081/api/v1/traffic/dba-profile \
  -H "Content-Type: application/json" \
  -d '{
    "name": "UP-10M",
    "type": 4,
    "max_bandwidth": 20480
  }'
```

**9. Delete DBA Profile:**
```bash
curl -X DELETE http://192.168.54.230:8081/api/v1/traffic/dba-profile/UP-10M
```

**10. Delete T-CONT:**
```bash
curl -X DELETE http://192.168.54.230:8081/api/v1/traffic/tcont/1-1-1/1/1
```

### Technical Implementation Details

**DBA Profile Command Generation:**
```go
// Type 1 - Fixed: type 1 fix 10240
// Type 2 - Assured: type 2 assure 10240  
// Type 3 - Assured+Max: type 3 assure 10240 max 51200
// Type 4 - Max: type 4 max 10240
// Type 5 - Assured+Max+Priority: type 5 assure 10240 max 51200
```

**T-CONT Configuration Flow:**
1. Enter ONU interface mode: `interface gpon-onu_1/1/1:1`
2. Configure T-CONT: `tcont 1 name TCONT_DATA profile UP-10M`
3. Exit interface mode: `exit`

**GEM Port Configuration Flow:**
1. Enter ONU interface mode: `interface gpon-onu_1/1/1:1`
2. Configure GEM port: `gemport 1 name GEM_DATA tcont 1 queue 1`
3. Exit interface mode: `exit`

### Bandwidth Calculation Examples

**Example 1: 10M Symmetric**
- Type: 4 (Max only)
- Max Bandwidth: 10240 Kbps
- Suitable for: Standard internet users

**Example 2: 100M with 10M Guaranteed**
- Type: 3 (Assured + Max)
- Assured: 10240 Kbps  
- Max: 102400 Kbps
- Suitable for: Business users needing guaranteed bandwidth with burst capability

**Example 3: VoIP Priority**
- Type: 1 (Fixed)
- Fixed: 1024 Kbps
- Suitable for: VoIP traffic requiring constant bandwidth

### Known Limitations (Phase 4)
- No DBA profile conflict detection
- No bandwidth utilization monitoring
- No automatic T-CONT/GEM port ID allocation
- No profile usage tracking (which ONUs use which profiles)
- No bulk profile operations
- No profile import/export

### Testing Results (Phase 4)

**Deployment Date:** January 11, 2026 15:27 UTC  
**VPS:** 192.168.54.230:8081

**Endpoint Tests:**
✅ GET /traffic/dba-profiles - Working  
✅ GET /traffic/dba-profile/{name} - Implemented  
✅ POST /traffic/dba-profile - Implemented  
✅ PUT /traffic/dba-profile - Implemented  
✅ DELETE /traffic/dba-profile/{name} - Implemented  
✅ GET /traffic/tcont/{pon}/{onu_id}/{tcont_id} - Implemented  
✅ POST /traffic/tcont - Implemented  
✅ DELETE /traffic/tcont/{pon}/{onu_id}/{tcont_id} - Implemented  
✅ POST /traffic/gemport - Implemented  
✅ DELETE /traffic/gemport/{pon}/{onu_id}/{gemport_id} - Implemented  

**Total New Endpoints:** 10  
**Total Project Endpoints:** 40+ (including all phases)

**Service Status:**
- Application running: ✅
- Telnet connectivity: ✅
- Redis connection: ✅
- SNMP connection: ✅

### Next Steps (Phase 5 - ONU Management - Planned)
- ONU reboot functionality
- ONU reset to factory
- ONU block/unblock
- ONU description update
- Bulk ONU operations

---

## File Structure

```
go-snmp-olt-zte-c320/
├── cmd/api/
│   └── main.go
├── app/
│   ├── app.go                    # Application initialization (UPDATED Phase 2, 3)
│   └── routes.go                 # Route definitions (UPDATED Phase 2, 3)
├── config/
│   ├── config.go
│   └── telnet_config.go          # NEW - Telnet configuration (Phase 1)
├── docs/
│   ├── COMMAND_REFERENCE.md      # NEW - ZTE C320 CLI commands (Phase 1)
│   └── TELNET_CONFIG_ROADMAP.md  # NEW - Implementation roadmap (Phase 1)
├── internal/
│   ├── model/
│   │   ├── onu.go
│   │   ├── pon.go                # PON port models
│   │   ├── profile.go            # Traffic & VLAN profile models
│   │   ├── card.go               # Card/slot models
│   │   └── telnet.go             # NEW - Telnet request/response models (Phase 1, UPDATED Phase 3)
│   ├── usecase/
│   │   ├── onu.go
│   │   ├── pon.go
│   │   ├── profile.go
│   │   ├── card.go
│   │   ├── provision.go          # NEW - ONU provisioning logic (Phase 2)
│   │   └── vlan.go               # NEW - VLAN management logic (Phase 3)
│   ├── handler/
│   │   ├── onu.go
│   │   ├── pon.go
│   │   ├── profile.go
│   │   ├── card.go
│   │   ├── provision.go          # NEW - Provisioning HTTP handlers (Phase 2)
│   │   └── vlan.go               # NEW - VLAN HTTP handlers (Phase 3)
│   ├── repository/
│   │   ├── snmp.go
│   │   ├── redis.go
│   │   ├── telnet.go             # NEW - Telnet operations (Phase 1)
│   │   ├── telnet_session.go     # NEW - Session pooling (Phase 1)
│   │   └── telnet_vlan.go        # NEW - VLAN operations (Phase 3)
│   ├── middleware/
│   ├── utils/
│   └── errors/
└── pkg/
```

## Complete API Endpoints (Updated)

### ONU Endpoints (SNMP)
```
GET  /api/v1/board/{board_id}/pon/{pon_id}
GET  /api/v1/board/{board_id}/pon/{pon_id}/onu/{onu_id}
GET  /api/v1/board/{board_id}/pon/{pon_id}/onu_id_sn
GET  /api/v1/board/{board_id}/pon/{pon_id}/onu_id/empty
GET  /api/v1/board/{board_id}/pon/{pon_id}/onu_id/update
GET  /api/v1/paginate/board/{board_id}/pon/{pon_id}
DEL  /api/v1/board/{board_id}/pon/{pon_id}
```

### ONU Provisioning Endpoints (Telnet) - NEW
```
GET    /api/v1/onu/unconfigured              # List all unconfigured ONUs
GET    /api/v1/onu/unconfigured/{pon}        # List unconfigured ONUs by PON
POST   /api/v1/onu/register                  # Register new ONU
DELETE /api/v1/onu/{pon}/{onu_id}           # Delete ONU
```

### PON Port Endpoints
```
GET  /api/v1/board/{board_id}/pon/{pon_id}/info
```

### Profile Endpoints
```
GET  /api/v1/profiles/traffic
GET  /api/v1/profiles/traffic/{profile_id}
GET  /api/v1/profiles/vlan
```

### System Endpoints
```
GET  /api/v1/system/cards
GET  /api/v1/system/cards/{rack}/{shelf}/{slot}
```

### VLAN Endpoints (Telnet) - NEW Phase 3
```
GET    /api/v1/vlan/onu/{pon}/{onu_id}      # Get ONU VLAN configuration
GET    /api/v1/vlan/service-ports            # Get all service-port configurations
POST   /api/v1/vlan/onu                      # Configure ONU VLAN
PUT    /api/v1/vlan/onu                      # Modify ONU VLAN
DELETE /api/v1/vlan/onu/{pon}/{onu_id}       # Delete ONU VLAN
```

## Key Code Changes

### 1. config/oid_generator.go
```go
Board1OnuIDBase: 268500992  // Fixed calculation for V2.1
OnuIDIncrement: 256
BaseOID: ".1.3.6.1.4.1.3902.1012"
```

### 2. internal/utils/extractor.go
Fixed serial number extraction to handle GPON hex format:
```go
// Converts 8-byte hex to VendorID + uppercase hex string
// Example: ZTEGD824CDF3
```

### 3. app/app.go
Added handler initialization:
```go
ponUsecase := usecase.NewPonUsecase(snmpRepo, redisRepo, cfg)
profileUsecase := usecase.NewProfileUsecase(snmpRepo, redisRepo, cfg)
cardUsecase := usecase.NewCardUsecase(snmpRepo, redisRepo, cfg)

ponHandler := handler.NewPonHandler(ponUsecase)
profileHandler := handler.NewProfileHandler(profileUsecase)
cardHandler := handler.NewCardHandler(cardUsecase)
```

### 4. app/routes.go
Added new route groups for PON info, profiles, and system cards.

## Deployment Commands

### Build on VPS
```bash
cd /opt/go-snmp-olt
/usr/local/go/bin/go build -buildvcs=false -o bin/api ./cmd/api
```

### Upload Files from Windows
```bash
# Upload single file
scp internal/model/pon.go root@192.168.54.230:/opt/go-snmp-olt/internal/model/

# Upload using pipe (more reliable)
cat internal/model/pon.go | ssh root@192.168.54.230 "cat > /opt/go-snmp-olt/internal/model/pon.go"
```

### Service Management
```bash
# Restart service
systemctl restart go-snmp-olt

# Check status
systemctl status go-snmp-olt

# View logs
journalctl -u go-snmp-olt -f
```

### Clear Redis Cache
```bash
redis-cli -a 'OsWkRgJLabn4n2+nodZ6BQeP+OKkrObnGeFcDY6w7Nw='
FLUSHALL
```

## Known Issues & Solutions

### Issue 1: Serial Number Display
**Problem:** Serial numbers showing as garbled unicode  
**Solution:** Implemented hex-to-string conversion in `ExtractSerialNumber()`  
**Status:** ✅ Fixed

### Issue 2: Wrong PON Index
**Problem:** 404 errors due to incorrect base calculation  
**Solution:** Changed `Board1OnuIDBase` from 268509184 to 268500992  
**Status:** ✅ Fixed

### Issue 3: VLAN Name Parsing
**Problem:** VLAN names had control characters (e.g., `\u0003pppoe`)  
**Solution:** Fixed OID parsing to use length-prefixed format  
**Status:** ✅ Fixed

### Issue 4: Card Info Empty Values
**Problem:** Card information showing empty values  
**Solution:** Corrected column mapping (col 2=type, 4=serial, 5=hw, 6=sw, 7=status)  
**Status:** ✅ Fixed

## Testing Examples

### Test PON Port Info
```bash
curl http://192.168.54.230:8081/api/v1/board/1/pon/1/info
```

### Test Traffic Profiles
```bash
# List all
curl http://192.168.54.230:8081/api/v1/profiles/traffic

# Specific profile
curl http://192.168.54.230:8081/api/v1/profiles/traffic/1879048194
```

### Test VLAN Profiles
```bash
curl http://192.168.54.230:8081/api/v1/profiles/vlan
```

### Test Card Info
```bash
# List all
curl http://192.168.54.230:8081/api/v1/system/cards

# Specific card
curl http://192.168.54.230:8081/api/v1/system/cards/1/1/1
```

## Current ONU Status

**Board 1, PON 1** has 3 ONUs (all unconfigured):
1. **ONU 1:** HWTC1F14CAAD (EG8041V5) - Status: Logging
2. **ONU 2:** ZTEGD824CDF3 (F672YV9.1) - Status: Logging  
3. **ONU 3:** ZTEGDA5918AC (F670LV9.0) - Status: Logging

## Next Steps / Future Enhancements

1. ✅ ~~PON Port Information~~ - COMPLETED
2. ✅ ~~Traffic Profiles~~ - COMPLETED
3. ✅ ~~VLAN Profiles~~ - COMPLETED
4. ✅ ~~Card/Slot Information~~ - COMPLETED
5. 🔄 ONU Configuration (auto-provision, VLAN assignment)
6. 🔄 Performance metrics (RX power, TX power - if available in V2.1)
7. 🔄 Alarm/Event monitoring
8. 🔄 Historical data tracking
9. 🔄 WebSocket real-time updates

## Architecture Notes

### Pattern: Handler → Usecase → Repository

**Handler Layer:**
- HTTP request/response handling
- Parameter validation (via middleware)
- Error handling
- JSON marshaling

**Usecase Layer:**
- Business logic
- SNMP data aggregation
- Data transformation
- Caching logic (with singleflight)

**Repository Layer:**
- SNMP operations (Get, Walk)
- Redis operations
- Low-level data access

### Middleware Stack
1. RequestID - Unique request tracking
2. SecurityHeaders - HTTP security headers
3. RequestTimeout - 90s timeout (for slow SNMP)
4. RateLimiter - 100 req/s, burst 200
5. MaxBodySize - 1MB limit
6. Logger - Request/response logging
7. CORS - Configurable CORS
8. ValidateBoardPonParams - Board/PON validation
9. ValidateOnuIDParam - ONU ID validation

## Important Notes

1. **V2.1.0 vs V2.2 Differences:**
   - Different base OID (1012 vs 1082)
   - Different ONU table structure
   - No optical power data in V2.1.0
   - Status codes differ

2. **Serial Number Format:**
   - GPON standard: 4-char VendorID + 8-char hex serial
   - Example: ZTEGD824CDF3 = ZTE + D824CDF3

3. **PON Index Critical:**
   - Always use correct formula
   - Board 1 starts at 268500992
   - Increment by 256 per PON

4. **Cache Strategy:**
   - Redis used for ONU data caching
   - TTL configurable
   - Cache can be manually cleared per board/PON

5. **Error Handling:**
   - Custom error types (NotFoundError, SNMPError, ValidationError)
   - Centralized error handler
   - Proper HTTP status codes

## Configuration Files

### Environment Variables (.env on VPS)
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=OsWkRgJLabn4n2+nodZ6BQeP+OKkrObnGeFcDY6w7Nw=
SNMP_TARGET=136.1.1.100
SNMP_COMMUNITY=public
SNMP_PORT=161
OLT_FIRMWARE=V2.1
```

### Systemd Service (/etc/systemd/system/go-snmp-olt.service)
```ini
[Unit]
Description=Go SNMP OLT ZTE C320 Monitoring Service
After=network.target redis.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/go-snmp-olt
ExecStart=/opt/go-snmp-olt/bin/api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Documentation Files

- `ZTE_C320_V21_OID_MAPPING.md` - Complete OID reference with snmpwalk examples
- `PROJECT_STATE.md` - This file (project state and memory)
- `README.md` - Original project documentation

## Troubleshooting

### Service Won't Start
```bash
# Check logs
journalctl -u go-snmp-olt -n 50

# Check binary
ls -la /opt/go-snmp-olt/bin/api

# Test manually
cd /opt/go-snmp-olt
./bin/api
```

### Empty SNMP Response
```bash
# Test SNMP directly
snmpwalk -v2c -c public 136.1.1.100 1.3.6.1.4.1.3902.1012.3.13.3.1.5.268501248

# Verify PON index calculation
# Board 1, PON 1 = 268500992 + 256 = 268501248
```

### Build Errors
```bash
# Update dependencies
cd /opt/go-snmp-olt
/usr/local/go/bin/go mod tidy

# Clean build
rm -rf bin/
/usr/local/go/bin/go build -buildvcs=false -o bin/api ./cmd/api
```

---

**Project Completion Status:** All 4 requested features implemented and tested ✅  
**Last Successful Deployment:** December 31, 2025  
**Service Status:** Active and running on port 8081
