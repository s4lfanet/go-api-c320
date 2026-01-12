# API Reference

Complete REST API documentation for ZTE C320 OLT Management API.

**Base URL:** `http://your-server:8081/api/v1`

**Content-Type:** `application/json`

---

## Table of Contents

1. [ONU Monitoring (SNMP)](#onu-monitoring-snmp)
2. [Real-time Monitoring (SNMP + Telnet)](#real-time-monitoring)
3. [ONU Provisioning (Telnet)](#onu-provisioning)
4. [VLAN Management (Telnet)](#vlan-management)
5. [Traffic Profiles (Telnet)](#traffic-profiles)
6. [ONU Management (Telnet)](#onu-management)
7. [Batch Operations (Telnet)](#batch-operations)
8. [Configuration Backup/Restore](#configuration-backup-restore)
9. [System Information (SNMP)](#system-information)
10. [Error Responses](#error-responses)

---

## ONU Monitoring (SNMP)

### List All ONUs on PON Port

Get all ONUs connected to a specific PON port.

**Endpoint:** `GET /board/{board_id}/pon/{pon_id}/`

**Parameters:**
- `board_id` (path, integer, required) - Board/card ID (1-20)
- `pon_id` (path, integer, required) - PON port ID (1-16)

**Example Request:**
```bash
curl http://localhost:8081/api/v1/board/2/pon/7
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "board": 2,
      "pon": 7,
      "onu_id": 3,
      "name": "Customer-001",
      "onu_type": "F670LV7.1",
      "serial_number": "ZTEGC1234ABCD",
      "rx_power": "-22.22",
      "status": "Online"
    },
    {
      "board": 2,
      "pon": 7,
      "onu_id": 4,
      "name": "Customer-002",
      "onu_type": "F670LV7.1",
      "serial_number": "ZTEGC5678EFGH",
      "rx_power": "-21.08",
      "status": "Online"
    }
  ]
}
```

---

### Get Specific ONU Details

Get detailed information for a specific ONU.

**Endpoint:** `GET /board/{board_id}/pon/{pon_id}/onu/{onu_id}`

**Parameters:**
- `board_id` (path, integer, required) - Board/card ID (1-20)
- `pon_id` (path, integer, required) - PON port ID (1-16)
- `onu_id` (path, integer, required) - ONU ID (1-128)

**Example Request:**
```bash
curl http://localhost:8081/api/v1/board/2/pon/7/onu/4
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "board": 2,
    "pon": 7,
    "onu_id": 4,
    "name": "Customer-002",
    "description": "Location Description",
    "onu_type": "F670LV7.1",
    "serial_number": "ZTEGC5678EFGH",
    "rx_power": "-20.71",
    "tx_power": "2.57",
    "status": "Online",
    "ip_address": "10.10.10.5",
    "last_online": "2024-08-11 10:09:37",
    "last_offline": "2024-08-11 10:08:35",
    "uptime": "5 days 13 hours 10 minutes 50 seconds",
    "last_down_time_duration": "0 days 0 hours 1 minutes 2 seconds",
    "offline_reason": "PowerOff",
    "gpon_optical_distance": "6701"
  }
}
```

---

### Get PON Port Information

Get PON port statistics and configuration.

**Endpoint:** `GET /board/{board_id}/pon/{pon_id}/info`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/board/2/pon/7/info
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "board": 2,
    "pon": 7,
    "admin_status": "up",
    "operational_status": "up",
    "total_onus": 8,
    "online_onus": 6,
    "offline_onus": 2,
    "rx_power": "-15.5",
    "tx_power": "3.2"
  }
}
```

---

### Get Available ONU IDs

Get list of available/unused ONU IDs on a PON port.

**Endpoint:** `GET /board/{board_id}/pon/{pon_id}/onu_id/empty`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/board/2/pon/5/onu_id/empty
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "board": 2,
      "pon": 5,
      "onu_id": 123
    },
    {
      "board": 2,
      "pon": 5,
      "onu_id": 124
    },
    {
      "board": 2,
      "pon": 5,
      "onu_id": 125
    }
  ]
}
```

---

### Get ONU IDs and Serial Numbers

Get mapping of all ONU IDs and their serial numbers.

**Endpoint:** `GET /board/{board_id}/pon/{pon_id}/onu_id_sn`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/board/2/pon/7/onu_id_sn
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "onu_id": 3,
      "serial_number": "ZTEGC1234ABCD"
    },
    {
      "onu_id": 4,
      "serial_number": "ZTEGC5678EFGH"
    }
  ]
}
```

---

### List ONUs with Pagination

Get paginated list of ONUs on a PON port.

**Endpoint:** `GET /paginate/board/{board_id}/pon/{pon_id}/`

**Query Parameters:**
- `page` (integer, optional, default: 1) - Page number
- `limit` (integer, optional, default: 10) - Items per page

**Example Request:**
```bash
curl "http://localhost:8081/api/v1/paginate/board/2/pon/7/?page=1&limit=10"
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "onus": [
      {
        "board": 2,
        "pon": 7,
        "onu_id": 3,
        "name": "Customer-001",
        "serial_number": "ZTEGC1234ABCD",
        "status": "Online"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 8,
      "total_pages": 1
    }
  }
}
```

---

## Real-time Monitoring

### Get Real-time ONU Monitoring

Get real-time monitoring data including optical power (Phase 7.2).

**Endpoint:** `GET /monitoring/onu/{pon}/{onuId}`

**Parameters:**
- `pon` (path, string, required) - PON port (format: "1/1/1")
- `onuId` (path, integer, required) - ONU ID

**Example Request:**
```bash
curl http://localhost:8081/api/v1/monitoring/onu/1/5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "pon_port": "1",
    "onu_id": 5,
    "serial_number": "ZTEG1234ABCD",
    "model": "ZTE-F660",
    "firmware_version": "V8.0.10P3",
    "online_status": 1,
    "statistics": {
      "rx_packets": 123456789,
      "rx_bytes": 98765432100,
      "rx_rate": "250.5 Mbps"
    },
    "optical": {
      "rx_power": -18.45,
      "tx_power": 2.35,
      "olt_rx_power": -18.45,
      "temperature": 42.5,
      "voltage": 3.28,
      "bias_current": 15.2,
      "rx_power_status": "normal",
      "tx_power_status": "normal",
      "temperature_status": "normal"
    },
    "last_update": "2026-01-12T03:30:00Z"
  }
}
```

**Optical Power Status Values:**
- `normal` - Within acceptable range
- `low` - Below threshold
- `high` - Above threshold

---

### Get PON Port Monitoring

Get aggregated monitoring for all ONUs on a PON port.

**Endpoint:** `GET /monitoring/pon/{pon}`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/monitoring/pon/1
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "pon_port": "1",
    "total_onus": 8,
    "online_onus": 6,
    "offline_onus": 2,
    "onus": [
      {
        "onu_id": 5,
        "serial_number": "ZTEG1234ABCD",
        "online_status": 1,
        "rx_power": -18.45,
        "tx_power": 2.35
      }
    ],
    "last_update": "2026-01-12T03:30:00Z"
  }
}
```

---

### Get OLT-wide Monitoring

Get OLT summary with all PON ports.

**Endpoint:** `GET /monitoring/olt`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/monitoring/olt
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "total_pon_ports": 16,
    "total_onus": 128,
    "online_onus": 115,
    "offline_onus": 13,
    "pon_ports": [
      {
        "pon_port": "1",
        "total_onus": 8,
        "online_onus": 6,
        "offline_onus": 2
      }
    ],
    "last_update": "2026-01-12T03:30:00Z"
  }
}
```

---

## ONU Provisioning

### List Unconfigured ONUs

Get all unconfigured ONUs detected on the OLT.

**Endpoint:** `GET /onu/unconfigured`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/onu/unconfigured
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "pon_port": "1/1/1",
      "serial_number": "ZTEGC9999XXXX",
      "loid": "",
      "detection_time": "2026-01-12 10:30:00"
    },
    {
      "pon_port": "1/1/2",
      "serial_number": "ZTEGC8888YYYY",
      "loid": "",
      "detection_time": "2026-01-12 10:35:00"
    }
  ]
}
```

---

### List Unconfigured ONUs by PON

Get unconfigured ONUs on a specific PON port.

**Endpoint:** `GET /onu/unconfigured/{pon}`

**Parameters:**
- `pon` (path, string, required) - PON port (format: "1/1/1")

**Example Request:**
```bash
curl http://localhost:8081/api/v1/onu/unconfigured/1/1/1
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "pon_port": "1/1/1",
      "serial_number": "ZTEGC9999XXXX",
      "loid": "",
      "detection_time": "2026-01-12 10:30:00"
    }
  ]
}
```

---

### Register New ONU

Register and configure a new ONU.

**Endpoint:** `POST /onu/register`

**Request Body:**
```json
{
  "pon_port": "1/1/1",
  "onu_id": 5,
  "serial_number": "ZTEG1234ABCD",
  "onu_type": "ZTE-F660",
  "name": "Customer_001",
  "description": "Customer at Location X"
}
```

**Parameters:**
- `pon_port` (string, required) - PON port (format: "1/1/1")
- `onu_id` (integer, required) - ONU ID (1-128)
- `serial_number` (string, required) - ONU serial number (12 characters)
- `onu_type` (string, required) - ONU model (e.g., "ZTE-F660", "F670LV7.1")
- `name` (string, optional) - ONU name/identifier
- `description` (string, optional) - ONU description

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/onu/register \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5,
    "serial_number": "ZTEG1234ABCD",
    "onu_type": "ZTE-F660",
    "name": "Customer_001"
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "ONU registered successfully",
  "data": {
    "pon_port": "1/1/1",
    "onu_id": 5,
    "serial_number": "ZTEG1234ABCD",
    "name": "Customer_001"
  }
}
```

---

### Delete ONU (Legacy)

Delete ONU configuration (legacy endpoint).

**Endpoint:** `DELETE /onu/{pon}/{onu_id}`

**Parameters:**
- `pon` (path, string, required) - PON port (format: "1/1/1")
- `onu_id` (path, integer, required) - ONU ID

**Example Request:**
```bash
curl -X DELETE http://localhost:8081/api/v1/onu/1/1/1/5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "ONU deleted successfully"
}
```

---

## VLAN Management

### Get ONU VLAN Configuration

Get VLAN configuration for a specific ONU.

**Endpoint:** `GET /vlan/onu/{pon}/{onu_id}`

**Parameters:**
- `pon` (path, string, required) - PON port (format: "1/1/1")
- `onu_id` (path, integer, required) - ONU ID

**Example Request:**
```bash
curl http://localhost:8081/api/v1/vlan/onu/1/1/1/5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "pon_port": "1/1/1",
    "onu_id": 5,
    "service_ports": [
      {
        "index": 1,
        "vport": 1,
        "user_vlan": 100,
        "vlan_mode": "tag",
        "svlan": 100,
        "cvlan": 200,
        "priority": 0
      }
    ]
  }
}
```

---

### Get All Service Ports

Get all service-port configurations on the OLT.

**Endpoint:** `GET /vlan/service-ports`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/vlan/service-ports
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "index": 1,
      "pon_port": "1/1/1",
      "onu_id": 5,
      "vport": 1,
      "user_vlan": 100,
      "vlan_mode": "tag",
      "svlan": 100,
      "cvlan": 200
    }
  ]
}
```

---

### Configure ONU VLAN

Configure VLAN for an ONU.

**Endpoint:** `POST /vlan/onu`

**Request Body:**
```json
{
  "pon_port": "1/1/1",
  "onu_id": 5,
  "svlan": 100,
  "cvlan": 200,
  "vlan_mode": "tag",
  "priority": 0,
  "vport": 1
}
```

**Parameters:**
- `pon_port` (string, required) - PON port
- `onu_id` (integer, required) - ONU ID
- `svlan` (integer, required) - Service VLAN (1-4094)
- `cvlan` (integer, optional) - Customer VLAN (1-4094)
- `vlan_mode` (string, required) - VLAN mode: "tag", "untag", "translation"
- `priority` (integer, optional) - Priority (0-7), default: 0
- `vport` (integer, optional) - Virtual port, default: 1

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/vlan/onu \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5,
    "svlan": 100,
    "cvlan": 200,
    "vlan_mode": "tag",
    "priority": 0
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "VLAN configured successfully"
}
```

---

### Modify ONU VLAN

Modify existing VLAN configuration.

**Endpoint:** `PUT /vlan/onu`

**Request Body:** (same as Configure VLAN)

**Example Request:**
```bash
curl -X PUT http://localhost:8081/api/v1/vlan/onu \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5,
    "svlan": 200,
    "cvlan": 300,
    "vlan_mode": "tag"
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "VLAN modified successfully"
}
```

---

### Delete ONU VLAN

Delete VLAN configuration for an ONU.

**Endpoint:** `DELETE /vlan/onu/{pon}/{onu_id}`

**Example Request:**
```bash
curl -X DELETE http://localhost:8081/api/v1/vlan/onu/1/1/1/5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "VLAN deleted successfully"
}
```

---

## Traffic Profiles

### List All DBA Profiles

Get all Dynamic Bandwidth Allocation profiles.

**Endpoint:** `GET /traffic/dba-profiles`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/traffic/dba-profiles
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "name": "100M_Profile",
      "type": 3,
      "assured_bandwidth": 51200,
      "max_bandwidth": 102400,
      "reference_count": 5
    },
    {
      "name": "200M_Profile",
      "type": 3,
      "assured_bandwidth": 102400,
      "max_bandwidth": 204800,
      "reference_count": 2
    }
  ]
}
```

**DBA Profile Types:**
- Type 1: Fixed bandwidth
- Type 2: Assured bandwidth
- Type 3: Assured + maximum bandwidth
- Type 4: Maximum bandwidth
- Type 5: Fixed + assured

---

### Get Specific DBA Profile

Get details of a specific DBA profile.

**Endpoint:** `GET /traffic/dba-profile/{name}`

**Parameters:**
- `name` (path, string, required) - Profile name

**Example Request:**
```bash
curl http://localhost:8081/api/v1/traffic/dba-profile/100M_Profile
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "name": "100M_Profile",
    "type": 3,
    "assured_bandwidth": 51200,
    "max_bandwidth": 102400,
    "reference_count": 5
  }
}
```

---

### Create DBA Profile

Create a new DBA profile.

**Endpoint:** `POST /traffic/dba-profile`

**Request Body:**
```json
{
  "name": "100M_Profile",
  "type": 3,
  "assured_bandwidth": 51200,
  "max_bandwidth": 102400
}
```

**Parameters:**
- `name` (string, required) - Profile name (max 32 chars)
- `type` (integer, required) - DBA type (1-5)
- `assured_bandwidth` (integer, optional) - Assured bandwidth in kbps
- `max_bandwidth` (integer, optional) - Maximum bandwidth in kbps

**Bandwidth Values:**
- 100 Mbps = 102400 kbps
- 50 Mbps = 51200 kbps
- 200 Mbps = 204800 kbps

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/traffic/dba-profile \
  -H "Content-Type: application/json" \
  -d '{
    "name": "100M_Profile",
    "type": 3,
    "assured_bandwidth": 51200,
    "max_bandwidth": 102400
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "DBA profile created successfully"
}
```

---

### Modify DBA Profile

Modify existing DBA profile.

**Endpoint:** `PUT /traffic/dba-profile`

**Request Body:** (same as Create DBA Profile)

**Example Request:**
```bash
curl -X PUT http://localhost:8081/api/v1/traffic/dba-profile \
  -H "Content-Type: application/json" \
  -d '{
    "name": "100M_Profile",
    "type": 3,
    "assured_bandwidth": 76800,
    "max_bandwidth": 153600
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "DBA profile modified successfully"
}
```

---

### Delete DBA Profile

Delete a DBA profile.

**Endpoint:** `DELETE /traffic/dba-profile/{name}`

**Example Request:**
```bash
curl -X DELETE http://localhost:8081/api/v1/traffic/dba-profile/100M_Profile
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "DBA profile deleted successfully"
}
```

---

### Get ONU T-CONT Configuration

Get T-CONT configuration for an ONU.

**Endpoint:** `GET /traffic/tcont/{pon}/{onu_id}/{tcont_id}`

**Parameters:**
- `pon` (path, string, required) - PON port
- `onu_id` (path, integer, required) - ONU ID
- `tcont_id` (path, integer, required) - T-CONT ID (1-8)

**Example Request:**
```bash
curl http://localhost:8081/api/v1/traffic/tcont/1/1/1/5/1
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "pon_port": "1/1/1",
    "onu_id": 5,
    "tcont_id": 1,
    "dba_profile": "100M_Profile"
  }
}
```

---

### Configure T-CONT

Configure T-CONT for an ONU.

**Endpoint:** `POST /traffic/tcont`

**Request Body:**
```json
{
  "pon_port": "1/1/1",
  "onu_id": 5,
  "tcont_id": 1,
  "dba_profile": "100M_Profile"
}
```

**Parameters:**
- `pon_port` (string, required) - PON port
- `onu_id` (integer, required) - ONU ID
- `tcont_id` (integer, required) - T-CONT ID (1-8)
- `dba_profile` (string, required) - DBA profile name

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/traffic/tcont \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5,
    "tcont_id": 1,
    "dba_profile": "100M_Profile"
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "T-CONT configured successfully"
}
```

---

### Delete T-CONT

Delete T-CONT configuration.

**Endpoint:** `DELETE /traffic/tcont/{pon}/{onu_id}/{tcont_id}`

**Example Request:**
```bash
curl -X DELETE http://localhost:8081/api/v1/traffic/tcont/1/1/1/5/1
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "T-CONT deleted successfully"
}
```

---

### Configure GEM Port

Configure GEM port for an ONU.

**Endpoint:** `POST /traffic/gemport`

**Request Body:**
```json
{
  "pon_port": "1/1/1",
  "onu_id": 5,
  "gemport_id": 1,
  "tcont_id": 1,
  "direction": "both"
}
```

**Parameters:**
- `pon_port` (string, required) - PON port
- `onu_id` (integer, required) - ONU ID
- `gemport_id` (integer, required) - GEM port ID (1-4096)
- `tcont_id` (integer, required) - T-CONT ID (1-8)
- `direction` (string, optional) - "upstream", "downstream", "both" (default: "both")

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/traffic/gemport \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5,
    "gemport_id": 1,
    "tcont_id": 1
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "GEM port configured successfully"
}
```

---

### Delete GEM Port

Delete GEM port configuration.

**Endpoint:** `DELETE /traffic/gemport/{pon}/{onu_id}/{gemport_id}`

**Example Request:**
```bash
curl -X DELETE http://localhost:8081/api/v1/traffic/gemport/1/1/1/5/1
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "GEM port deleted successfully"
}
```

---

## ONU Management

### Reboot ONU

Reboot a specific ONU.

**Endpoint:** `POST /onu-management/reboot`

**Request Body:**
```json
{
  "pon_port": "1/1/1",
  "onu_id": 5
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/onu-management/reboot \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "ONU reboot command sent successfully"
}
```

---

### Block ONU

Block/disable an ONU.

**Endpoint:** `POST /onu-management/block`

**Request Body:**
```json
{
  "pon_port": "1/1/1",
  "onu_id": 5
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/onu-management/block \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "ONU blocked successfully"
}
```

---

### Unblock ONU

Unblock/enable an ONU.

**Endpoint:** `POST /onu-management/unblock`

**Request Body:**
```json
{
  "pon_port": "1/1/1",
  "onu_id": 5
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/onu-management/unblock \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "ONU unblocked successfully"
}
```

---

### Update ONU Description

Update ONU description/name.

**Endpoint:** `PUT /onu-management/description`

**Request Body:**
```json
{
  "pon_port": "1/1/1",
  "onu_id": 5,
  "description": "New customer location description"
}
```

**Example Request:**
```bash
curl -X PUT http://localhost:8081/api/v1/onu-management/description \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5,
    "description": "Customer at Building A"
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "ONU description updated successfully"
}
```

---

### Delete ONU

Delete ONU configuration completely.

**Endpoint:** `DELETE /onu-management/{pon}/{onu_id}`

**Example Request:**
```bash
curl -X DELETE http://localhost:8081/api/v1/onu-management/1/1/1/5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "ONU deleted successfully"
}
```

---

## Batch Operations

### Batch Reboot ONUs

Reboot multiple ONUs at once.

**Endpoint:** `POST /batch/reboot`

**Request Body:**
```json
{
  "onus": [
    {
      "pon_port": "1/1/1",
      "onu_id": 5
    },
    {
      "pon_port": "1/1/1",
      "onu_id": 6
    },
    {
      "pon_port": "1/1/2",
      "onu_id": 3
    }
  ]
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/batch/reboot \
  -H "Content-Type: application/json" \
  -d '{
    "onus": [
      {"pon_port": "1/1/1", "onu_id": 5},
      {"pon_port": "1/1/1", "onu_id": 6}
    ]
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "Batch reboot completed",
  "data": {
    "total": 2,
    "successful": 2,
    "failed": 0,
    "results": [
      {
        "pon_port": "1/1/1",
        "onu_id": 5,
        "status": "success"
      },
      {
        "pon_port": "1/1/1",
        "onu_id": 6,
        "status": "success"
      }
    ]
  }
}
```

---

### Batch Block ONUs

Block multiple ONUs at once.

**Endpoint:** `POST /batch/block`

**Request Body:**
```json
{
  "onus": [
    {
      "pon_port": "1/1/1",
      "onu_id": 5
    },
    {
      "pon_port": "1/1/1",
      "onu_id": 6
    }
  ]
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/batch/block \
  -H "Content-Type: application/json" \
  -d '{
    "onus": [
      {"pon_port": "1/1/1", "onu_id": 5},
      {"pon_port": "1/1/1", "onu_id": 6}
    ]
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "Batch block completed",
  "data": {
    "total": 2,
    "successful": 2,
    "failed": 0
  }
}
```

---

### Batch Unblock ONUs

Unblock multiple ONUs at once.

**Endpoint:** `POST /batch/unblock`

**Request Body:** (same format as batch block)

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/batch/unblock \
  -H "Content-Type: application/json" \
  -d '{
    "onus": [
      {"pon_port": "1/1/1", "onu_id": 5},
      {"pon_port": "1/1/1", "onu_id": 6}
    ]
  }'
```

---

### Batch Delete ONUs

Delete multiple ONUs at once.

**Endpoint:** `POST /batch/delete`

**Request Body:**
```json
{
  "onus": [
    {
      "pon_port": "1/1/1",
      "onu_id": 5
    },
    {
      "pon_port": "1/1/1",
      "onu_id": 6
    }
  ]
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/batch/delete \
  -H "Content-Type: application/json" \
  -d '{
    "onus": [
      {"pon_port": "1/1/1", "onu_id": 5},
      {"pon_port": "1/1/1", "onu_id": 6}
    ]
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "Batch delete completed",
  "data": {
    "total": 2,
    "successful": 2,
    "failed": 0
  }
}
```

---

### Batch Update Descriptions

Update descriptions for multiple ONUs at once.

**Endpoint:** `PUT /batch/descriptions`

**Request Body:**
```json
{
  "onus": [
    {
      "pon_port": "1/1/1",
      "onu_id": 5,
      "description": "Customer A - Building 1"
    },
    {
      "pon_port": "1/1/1",
      "onu_id": 6,
      "description": "Customer B - Building 2"
    }
  ]
}
```

**Example Request:**
```bash
curl -X PUT http://localhost:8081/api/v1/batch/descriptions \
  -H "Content-Type: application/json" \
  -d '{
    "onus": [
      {
        "pon_port": "1/1/1",
        "onu_id": 5,
        "description": "Customer A - Building 1"
      }
    ]
  }'
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "Batch update completed",
  "data": {
    "total": 1,
    "successful": 1,
    "failed": 0
  }
}
```

---

## Configuration Backup/Restore

### Backup Single ONU

Create backup of single ONU configuration.

**Endpoint:** `POST /config/backup/onu/{pon}/{onuId}`

**Parameters:**
- `pon` (path, string, required) - PON port
- `onuId` (path, integer, required) - ONU ID

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/config/backup/onu/1/1/1/5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "ONU backup created successfully",
  "data": {
    "backup_id": "backup_20260112_103045_onu_1_1_1_5",
    "pon_port": "1/1/1",
    "onu_id": 5,
    "timestamp": "2026-01-12T10:30:45Z"
  }
}
```

---

### Backup Entire OLT

Create backup of entire OLT configuration.

**Endpoint:** `POST /config/backup/olt`

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/config/backup/olt
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "OLT backup created successfully",
  "data": {
    "backup_id": "backup_20260112_103045_olt_full",
    "total_onus": 128,
    "timestamp": "2026-01-12T10:30:45Z"
  }
}
```

---

### List All Backups

Get list of all configuration backups.

**Endpoint:** `GET /config/backups`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/config/backups
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "backup_id": "backup_20260112_103045_olt_full",
      "type": "olt",
      "timestamp": "2026-01-12T10:30:45Z",
      "size": "2.5 MB",
      "onu_count": 128
    },
    {
      "backup_id": "backup_20260112_100000_onu_1_1_1_5",
      "type": "onu",
      "pon_port": "1/1/1",
      "onu_id": 5,
      "timestamp": "2026-01-12T10:00:00Z",
      "size": "15 KB"
    }
  ]
}
```

---

### Get Specific Backup

Get details of a specific backup.

**Endpoint:** `GET /config/backup/{backupId}`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/config/backup/backup_20260112_103045_onu_1_1_1_5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "backup_id": "backup_20260112_103045_onu_1_1_1_5",
    "type": "onu",
    "pon_port": "1/1/1",
    "onu_id": 5,
    "timestamp": "2026-01-12T10:30:45Z",
    "configuration": {
      "serial_number": "ZTEG1234ABCD",
      "onu_type": "ZTE-F660",
      "name": "Customer_001",
      "vlan": {
        "svlan": 100,
        "cvlan": 200
      },
      "traffic": {
        "dba_profile": "100M_Profile",
        "tcont_id": 1
      }
    }
  }
}
```

---

### Delete Backup

Delete a backup.

**Endpoint:** `DELETE /config/backup/{backupId}`

**Example Request:**
```bash
curl -X DELETE http://localhost:8081/api/v1/config/backup/backup_20260112_103045_onu_1_1_1_5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "Backup deleted successfully"
}
```

---

### Export Backup

Export backup as downloadable file.

**Endpoint:** `GET /config/backup/{backupId}/export`

**Example Request:**
```bash
curl -O -J http://localhost:8081/api/v1/config/backup/backup_20260112_103045_onu_1_1_1_5/export
```

**Response:** Binary file download (JSON format)

---

### Import Backup

Import backup from file.

**Endpoint:** `POST /config/backup/import`

**Request:** Multipart form data with file upload

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/config/backup/import \
  -F "file=@backup_20260112_103045_onu_1_1_1_5.json"
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "Backup imported successfully",
  "data": {
    "backup_id": "backup_20260112_103045_onu_1_1_1_5"
  }
}
```

---

### Restore from Backup

Restore configuration from a backup.

**Endpoint:** `POST /config/restore/{backupId}`

**Parameters:**
- `backupId` (path, string, required) - Backup ID

**Example Request:**
```bash
curl -X POST http://localhost:8081/api/v1/config/restore/backup_20260112_103045_onu_1_1_1_5
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "message": "Configuration restored successfully",
  "data": {
    "backup_id": "backup_20260112_103045_onu_1_1_1_5",
    "restored_onus": 1
  }
}
```

---

## System Information

### Get All Cards/Slots

Get information about all cards/slots in the OLT.

**Endpoint:** `GET /system/cards`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/system/cards
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "rack": 1,
      "shelf": 1,
      "slot": 1,
      "card_type": "GTGH",
      "status": "online",
      "description": "PON Card 16-port"
    },
    {
      "rack": 1,
      "shelf": 1,
      "slot": 2,
      "card_type": "GTGH",
      "status": "online",
      "description": "PON Card 16-port"
    }
  ]
}
```

---

### Get Specific Card

Get information about a specific card/slot.

**Endpoint:** `GET /system/cards/{rack}/{shelf}/{slot}`

**Parameters:**
- `rack` (path, integer, required) - Rack number
- `shelf` (path, integer, required) - Shelf number
- `slot` (path, integer, required) - Slot number

**Example Request:**
```bash
curl http://localhost:8081/api/v1/system/cards/1/1/1
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "rack": 1,
    "shelf": 1,
    "slot": 1,
    "card_type": "GTGH",
    "status": "online",
    "description": "PON Card 16-port",
    "serial_number": "ZTEC12345678",
    "firmware_version": "V2.1.0"
  }
}
```

---

### Get All Traffic Profiles (SNMP)

Get all traffic profiles from SNMP.

**Endpoint:** `GET /profiles/traffic`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/profiles/traffic
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "profile_id": 1,
      "name": "Default_Profile",
      "upstream_rate": 1024000,
      "downstream_rate": 1024000
    }
  ]
}
```

---

### Get Specific Traffic Profile (SNMP)

Get specific traffic profile details.

**Endpoint:** `GET /profiles/traffic/{profile_id}`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/profiles/traffic/1
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "profile_id": 1,
    "name": "Default_Profile",
    "upstream_rate": 1024000,
    "downstream_rate": 1024000
  }
}
```

---

### Get All VLAN Profiles (SNMP)

Get all VLAN profiles from SNMP.

**Endpoint:** `GET /profiles/vlan`

**Example Request:**
```bash
curl http://localhost:8081/api/v1/profiles/vlan
```

**Success Response (200 OK):**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "profile_id": 1,
      "name": "Default_VLAN",
      "vlan_mode": "tag",
      "svlan": 100
    }
  ]
}
```

---

## Error Responses

### Common Error Codes

All error responses follow this format:

```json
{
  "code": 400,
  "status": "Bad Request",
  "message": "Error description here",
  "error": "Detailed error message"
}
```

### HTTP Status Codes

- **200 OK** - Request successful
- **400 Bad Request** - Invalid request parameters
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - OLT connection failed

### Example Error Responses

**Invalid Parameters (400):**
```json
{
  "code": 400,
  "status": "Bad Request",
  "message": "Invalid PON port format",
  "error": "PON port must be in format '1/1/1'"
}
```

**ONU Not Found (404):**
```json
{
  "code": 404,
  "status": "Not Found",
  "message": "ONU not found",
  "error": "No ONU found with ID 5 on PON port 1/1/1"
}
```

**OLT Connection Error (503):**
```json
{
  "code": 503,
  "status": "Service Unavailable",
  "message": "Failed to connect to OLT",
  "error": "SNMP timeout: no response from 192.168.1.1:161"
}
```

**Telnet Command Error (500):**
```json
{
  "code": 500,
  "status": "Internal Server Error",
  "message": "Failed to execute Telnet command",
  "error": "Command execution timeout after 30 seconds"
}
```

---

## Rate Limiting

API has built-in rate limiting:
- **100 requests/second** per IP
- **Burst up to 200 requests**

**Rate Limit Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1673456789
```

**Rate Limit Exceeded Response (429):**
```json
{
  "code": 429,
  "status": "Too Many Requests",
  "message": "Rate limit exceeded",
  "error": "Maximum 100 requests per second allowed"
}
```

---

## Authentication

Currently, the API does **not** implement built-in authentication. For production use:

1. **Deploy behind reverse proxy** (nginx, Traefig)
2. **Enable HTTPS/TLS** at proxy level
3. **Implement authentication** at proxy (Basic Auth, JWT, OAuth)
4. **Use firewall rules** to restrict access

**Example nginx configuration:**
```nginx
location /api/ {
    auth_basic "OLT API";
    auth_basic_user_file /etc/nginx/.htpasswd;
    proxy_pass http://localhost:8081;
}
```

---

## CORS Configuration

CORS is configurable via environment variables:

```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://app.example.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE
CORS_ALLOWED_HEADERS=Content-Type,Authorization
CORS_ALLOW_CREDENTIALS=true
```

**Default:** All origins allowed (`*`)

---

## Pagination

Endpoints supporting pagination use query parameters:

**Parameters:**
- `page` (integer, optional, default: 1) - Page number
- `limit` (integer, optional, default: 10, max: 100) - Items per page

**Response Format:**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 128,
      "total_pages": 13
    }
  }
}
```

---

## Webhook Support

**Coming in Phase 8:** Webhook notifications for ONU state changes.

Planned events:
- `onu.online` - ONU came online
- `onu.offline` - ONU went offline
- `onu.registered` - New ONU registered
- `onu.deleted` - ONU deleted
- `optical.warning` - Optical power warning

---

## SDK & Client Libraries

**Coming Soon:**
- JavaScript/TypeScript SDK
- Python SDK
- Go SDK
- PHP SDK

---

## Testing

Test the API using:

**cURL:**
```bash
curl http://localhost:8081/api/v1/board/2/pon/7
```

**Postman:** Import OpenAPI spec (coming soon)

**HTTPie:**
```bash
http GET localhost:8081/api/v1/board/2/pon/7
```

**JavaScript fetch:**
```javascript
fetch('http://localhost:8081/api/v1/board/2/pon/7')
  .then(res => res.json())
  .then(data => console.log(data));
```

**Python requests:**
```python
import requests
response = requests.get('http://localhost:8081/api/v1/board/2/pon/7')
print(response.json())
```

---

## Support

- **GitHub Issues:** [github.com/s4lfanet/go-api-c320/issues](https://github.com/s4lfanet/go-api-c320/issues)
- **Email:** wardian370@gmail.com
- **Documentation:** [docs/README.md](../README.md)

---

**Last Updated:** January 12, 2026
**API Version:** 1.0
**Phase:** 7.2 (Optical Power Monitoring)
