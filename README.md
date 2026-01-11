# ZTE C320 OLT Management API
[![Go Report Card](https://goreportcard.com/badge/github.com/s4lfanet/go-api-c320)](https://goreportcard.com/report/github.com/s4lfanet/go-api-c320)

Production-ready REST API for ZTE C320 OLT monitoring and configuration with SNMP & Telnet integration.

## 🚀 Features

### Phase 1-5 Complete ✅
- **ONU Monitoring** (SNMP) - Real-time status, signal levels, models
- **ONU Provisioning** (Telnet) - Auto-registration, configuration
- **VLAN Management** (Telnet) - Service-port creation, VLAN assignment
- **Traffic Profiles** (Telnet) - DBA profiles, T-CONT, GEM ports
- **ONU Management** (Telnet) - Reboot, block/unblock, description, delete

### Total: 45+ REST API Endpoints

## 📋 Technology Stack
* [Go](https://go.dev/) - Programming language
* [Chi](https://github.com/go-chi/chi/) - HTTP Server
* [GoSNMP](https://github.com/gosnmp/gosnmp) - SNMP library for Go
* [Redis](https://github.com/redis/go-redis/v9) - Redis client for Go
* [Zerolog](https://github.com/rs/zerolog) - Logger
* [Viper](https://github.com/spf13/viper) - Configuration management
* [Docker](https://www.docker.com/) - Containerization
* [Task](https://github.com/go-task/task) - Task runner
* [Air](https://github.com/cosmtrek/air) - Live reload for Go apps


## 📋 Technology Stack
* [Go 1.24+](https://go.dev/) - Programming language
* [Chi Router](https://github.com/go-chi/chi/) - HTTP server & routing
* [GoSNMP](https://github.com/gosnmp/gosnmp) - SNMP library
* [Telnet](https://github.com/ziutek/telnet) - Telnet client for OLT configuration
* [Redis](https://github.com/redis/go-redis/v9) - Caching layer
* [Zerolog](https://github.com/rs/zerolog) - Structured logging
* [Docker](https://www.docker.com/) - Containerization

## 🔧 Prerequisites

- Go 1.24 or higher
- Redis 7.2+
- ZTE C320 OLT with:
  - SNMP v2c enabled (UDP 161)
  - Telnet enabled (TCP 23)
  - Default credentials: `zte/zte` (enable: `zxr10`)

## 📦 Installation

### Quick Start (Production)

1. Clone repository:
```bash
git clone https://github.com/s4lfanet/go-api-c320.git
cd go-api-c320
```

2. Set environment variables:
```bash
export OLT_IP_ADDRESS=192.168.1.1
export OLT_SNMP_PORT=161
export OLT_SNMP_COMMUNITY=public
export OLT_TELNET_HOST=192.168.1.1
export OLT_TELNET_PORT=23
export OLT_TELNET_USERNAME=zte
export OLT_TELNET_PASSWORD=zte
export OLT_TELNET_ENABLE_PASSWORD=zxr10
export REDIS_HOST=localhost
export REDIS_PORT=6379
```

3. Build and run:
```bash
go build -o api cmd/api/main.go
./api
```

### Docker Deployment

```bash
docker network create olt-network
docker run -d --name redis --network olt-network redis:7.2
docker build -t go-api-c320 .
docker run -d -p 8081:8081 --name olt-api \
  --network olt-network \
  -e REDIS_HOST=redis \
  -e OLT_IP_ADDRESS=192.168.1.1 \
  go-api-c320
```

## 🌐 API Endpoints

Base URL: `http://localhost:8081/api/v1`

### ONU Monitoring (SNMP)
- `GET /board/{board_id}/pon/{pon_id}/` - List all ONUs on PON port
- `GET /board/{board_id}/pon/{pon_id}/onu/{onu_id}` - Get specific ONU
- `GET /board/{board_id}/pon/{pon_id}/info` - Get PON port info
- `GET /board/{board_id}/pon/{pon_id}/onu_id/empty` - Get available ONU IDs

### ONU Provisioning (Telnet)
- `GET /onu/unconfigured` - List unconfigured ONUs
- `GET /onu/unconfigured/{pon}` - List unconfigured ONUs by PON
- `POST /onu/register` - Register new ONU
- `DELETE /onu/{pon}/{onu_id}` - Delete ONU (legacy)

### VLAN Management (Telnet)
- `GET /vlan/onu/{pon}/{onu_id}` - Get ONU VLAN config
- `GET /vlan/service-ports` - List all service-ports
- `POST /vlan/onu` - Configure ONU VLAN
- `PUT /vlan/onu` - Modify ONU VLAN
- `DELETE /vlan/onu/{pon}/{onu_id}` - Delete VLAN config

### Traffic Profiles (Telnet)
- `GET /traffic/dba-profiles` - List DBA profiles
- `GET /traffic/dba-profile/{name}` - Get DBA profile
- `POST /traffic/dba-profile` - Create DBA profile
- `PUT /traffic/dba-profile` - Modify DBA profile
- `DELETE /traffic/dba-profile/{name}` - Delete DBA profile
- `GET /traffic/tcont/{pon}/{onu_id}/{tcont_id}` - Get T-CONT
- `POST /traffic/tcont` - Configure T-CONT
- `DELETE /traffic/tcont/{pon}/{onu_id}/{tcont_id}` - Delete T-CONT
- `POST /traffic/gemport` - Configure GEM port
- `DELETE /traffic/gemport/{pon}/{onu_id}/{gemport_id}` - Delete GEM port

### ONU Management (Telnet)
- `POST /onu-management/reboot` - Reboot ONU
- `POST /onu-management/block` - Block/disable ONU
- `POST /onu-management/unblock` - Unblock/enable ONU
- `PUT /onu-management/description` - Update ONU description
- `DELETE /onu-management/{pon}/{onu_id}` - Delete ONU configuration

### System Info (SNMP)
- `GET /system/cards` - List all cards/slots
- `GET /system/cards/{rack}/{shelf}/{slot}` - Get card info
- `GET /profiles/traffic` - List traffic profiles
- `GET /profiles/traffic/{profile_id}` - Get traffic profile
- `GET /profiles/vlan` - List VLAN profiles

## 📝 Example Usage

### Register New ONU
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

### Configure VLAN
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

### Create DBA Profile
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

### Reboot ONU
```bash
curl -X POST http://localhost:8081/api/v1/onu-management/reboot \
  -H "Content-Type: application/json" \
  -d '{
    "pon_port": "1/1/1",
    "onu_id": 5
  }'
```

## ⚙️ Configuration

All configuration via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `OLT_IP_ADDRESS` | - | OLT IP address |
| `OLT_SNMP_PORT` | 161 | SNMP port |
| `OLT_SNMP_COMMUNITY` | public | SNMP community string |
| `OLT_TELNET_HOST` | - | Telnet host (usually same as OLT IP) |
| `OLT_TELNET_PORT` | 23 | Telnet port |
| `OLT_TELNET_USERNAME` | zte | Telnet username |
| `OLT_TELNET_PASSWORD` | zte | Telnet password |
| `OLT_TELNET_ENABLE_PASSWORD` | zxr10 | Telnet enable password |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `REDIS_PASSWORD` | - | Redis password (optional) |
| `REDIS_DB` | 0 | Redis database number |
| `APP_PORT` | 8081 | API server port |
| `LOG_LEVEL` | info | Log level (debug/info/warn/error) |

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP API (Chi Router)                 │
│                       Port 8081                          │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Handler Layer                         │
│  ONU │ PON │ Profile │ Provision │ VLAN │ Traffic │ Mgmt│
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Usecase Layer                         │
│         Business Logic │ Validation │ Orchestration     │
└─────────────────────────────────────────────────────────┘
                          │
            ┌─────────────┴─────────────┐
            ▼                           ▼
┌────────────────────┐      ┌───────────────────────────┐
│  SNMP Repository   │      │   Telnet Repository       │
│  (Read-Only)       │      │   (Read/Write Config)     │
└────────────────────┘      └───────────────────────────┘
            │                           │
            ▼                           ▼
┌─────────────────────────────────────────────────────────┐
│                    ZTE C320 OLT                          │
│              SNMP (UDP 161) │ Telnet (TCP 23)           │
└─────────────────────────────────────────────────────────┘
```

## 🧪 Testing

Run tests:
```bash
go test ./... -v
```

Run with coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📚 Documentation

See `docs/` directory for:
- [COMMAND_REFERENCE.md](docs/COMMAND_REFERENCE.md) - ZTE C320 Telnet commands
- [TELNET_CONFIG_ROADMAP.md](docs/TELNET_CONFIG_ROADMAP.md) - Implementation roadmap
- [PROJECT_STATE.md](docs/PROJECT_STATE.md) - Current project state & progress

## 🔒 Security

- **Never expose Telnet credentials** in code or logs
- Use environment variables for sensitive data
- Run API behind reverse proxy (nginx/traefik) in production
- Enable HTTPS/TLS termination at proxy level
- Use firewall rules to restrict OLT access
- Implement rate limiting for API endpoints (built-in)

## 🐛 Known Limitations

- Single Telnet session (sequential command execution)
- No authentication/authorization on API (add via reverse proxy)
- SNMP read-only (by design)
- Telnet timeout: 30 seconds per command
- No support for SNMP v3 (currently v2c only)

## 🛠️ Development

### Local Development
```bash
# Install dependencies
go mod download

# Run with hot reload
go install github.com/cosmtrek/air@latest
air

# Or manual run
go run cmd/api/main.go
```

### Build for Production
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o api cmd/api/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o api.exe cmd/api/main.go
```

## 📄 License

This project is licensed under the MIT License.

## 👥 Contributors

- **s4lfanet** - Initial work & development

## 🙏 Acknowledgments

- ZTE for C320 OLT platform
- Go community for excellent libraries
- Redis for caching infrastructure

## 📞 Support

For issues and questions:
- Create issue on [GitHub](https://github.com/s4lfanet/go-api-c320/issues)
- Email: wardian370@gmail.com

---

**Production Status:** ✅ Ready for deployment (Phase 1-5 complete)

**Last Updated:** January 11, 2026
-e REDIS_POOL_TIMEOUT=240 -e SNMP_HOST=x.x.x.x \
-e SNMP_PORT=161 -e SNMP_COMMUNITY=xxxx \
cepatkilatteknologi/snmp-olt-zte-c320:latest
```

### Production usage without external redis:
```shell
docker run -d -p 8081:8081 --name go-snmp-olt-zte-c320 \
-e REDIS_HOST=redis_host \
-e REDIS_PORT=redis_port \
-e REDIS_DB=redis_db \
-e REDIS_MIN_IDLE_CONNECTIONS=redis_min_idle_connection \
-e REDIS_POOL_SIZE=redis_pool_size \
-e REDIS_POOL_TIMEOUT=redis_pool_timeout \
-e SNMP_HOST=snmp_host \
-e SNMP_PORT=snmp_port \
-e SNMP_COMMUNITY=snmp_community \
cepatkilatteknologi/snmp-olt-zte-c320:latest
```


### Available Tasks for this Project:

Run `task --list` or `task help` to see all available tasks.

#### Development Tasks

| Task               | Description                                                     |
|--------------------|-----------------------------------------------------------------|
| `task init`        | Initialize the development environment                          |
| `task dev`         | Start local development (Redis in Docker + App with hot reload)|
| `task dev-docker`  | Start full development environment in Docker (with hot reload) |
| `task dev-down`    | Stop local development environment                              |
| `task dl-deps`     | Install tools required to run/build this app                    |

#### Testing Tasks

| Task                  | Description                                                  |
|-----------------------|--------------------------------------------------------------|
| `task test`           | Run all unit tests                                           |
| `task test-verbose`   | Run all unit tests with verbose output                       |
| `task test-coverage`  | Run tests and generate coverage report (text)                |
| `task test-html`      | Generate HTML coverage report and open in browser            |
| `task test-race`      | Run tests with race detection                                |
| `task test-short`     | Run short tests (excluding integration tests)                |
| `task load-test`      | Run k6 load testing                                          |
| `task benchmark`      | Run benchmarks                                               |

#### Build Tasks

| Task                    | Description                                                |
|-------------------------|------------------------------------------------------------|
| `task app-build`        | Build the app binary                                       |
| `task build-image`      | Build the docker image (local)                             |
| `task build-image-prod` | Build production docker image (multi-arch)                 |
| `task push-image`       | Build and push docker image with multi-arch to Docker Hub  |
| `task pull-image`       | Pull docker image from Docker Hub                          |

#### Docker Development Tasks

| Task             | Description                                                     |
|------------------|-----------------------------------------------------------------|
| `task up`        | Start development Docker environment                            |
| `task down`      | Stop development Docker environment                             |
| `task restart`   | Restart development Docker environment                          |
| `task logs`      | View development container logs                                 |
| `task logs-redis`| View development Redis logs                                     |
| `task ps`        | Show development containers status                              |

#### Production Deployment Tasks

| Task                 | Description                                                  |
|----------------------|--------------------------------------------------------------|
| `task prod-up`       | Start production containers (requires .env.prod)             |
| `task prod-down`     | Stop production containers                                   |
| `task prod-restart`  | Restart production containers                                |
| `task prod-rebuild`  | Rebuild and restart production containers                    |
| `task prod-logs`     | View production container logs                               |
| `task prod-logs-redis` | View production Redis logs                                 |
| `task prod-ps`       | Show production containers status                            |

#### Cleanup Tasks

| Task                | Description                                                   |
|---------------------|---------------------------------------------------------------|
| `task tidy`         | Clean up Go dependencies                                      |
| `task clean`        | Clean up all containers, volumes, and build artifacts         |
| `task clean-cache`  | Clean Go module cache and build cache                         |

### Test with curl GET method Board 2 Pon 7
``` shell
curl -sS localhost:8081/api/v1/board/2/pon/7 | jq
```
### Result
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
      "serial_number": "ZTEGC*******",
      "rx_power": "-22.22",
      "status": "Online"
    },
    {
      "board": 2,
      "pon": 7,
      "onu_id": 4,
      "name": "Customer-002",
      "onu_type": "F670LV7.1",
      "serial_number": "ZTEGC*******",
      "rx_power": "-21.08",
      "status": "Online"
    },
    {
      "board": 2,
      "pon": 7,
      "onu_id": 5,
      "name": "Customer-003",
      "onu_type": "F670LV7.1",
      "serial_number": "ZTEGC*******",
      "rx_power": "-19.956",
      "status": "Online"
    }
  ]
}
```

### Test with curl GET method Board 2 Pon 7 Onu 4
```shell
 curl -sS localhost:8081/api/v1/board/2/pon/7/onu/4 | jq
```

### Result
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
    "serial_number": "ZTEGC*******",
    "rx_power": "-20.71",
    "tx_power": "2.57",
    "status": "Online",
    "ip_address": "10.x.x.x",
    "last_online": "2024-08-11 10:09:37",
    "last_offline": "2024-08-11 10:08:35",
    "uptime": "5 days 13 hours 10 minutes 50 seconds",
    "last_down_time_duration": "0 days 0 hours 1 minutes 2 seconds",
    "offline_reason": "PowerOff",
    "gpon_optical_distance": "6701"
  }
}
```

### Test with curl GET method Get Empty ONU_ID in Board 2 Pon 5
```shell
curl -sS localhost:8081/api/v1/board/2/pon/5/onu_id/empty | jq
```

### Result
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
    },
    {
      "board": 2,
      "pon": 5,
      "onu_id": 126
    }
  ]
}
```

### Test with curl GET method Get Empty ONU_ID After Add ONU in Board 2 Pon 5
```shell
curl -sS localhost:8081/api/v1/board/2/pon/5/onu_id/update | jq
```

```json
{
  "code": 200,
  "status": "OK",
  "data": "Success Update Empty ONU_ID"
}
```

### Test with curl GET method Get Onu Information in Board 2 Pon 8 with paginate
```shell
curl -sS 'http://localhost:8081/api/v1/paginate/board/2/pon/8?limit=3&page=2' | jq
```
### Result
```json
{
  "code": 200,
  "status": "OK",
  "page": 2,
  "limit": 3,
  "page_count": 23,
  "total_rows": 69,
  "data": [
    {
      "board": 2,
      "pon": 8,
      "onu_id": 4,
      "name": "Customer-004",
      "onu_type": "F670LV7.1",
      "serial_number": "ZTEGC*******",
      "rx_power": "-19.17",
      "status": "Online"
    },
    {
      "board": 2,
      "pon": 8,
      "onu_id": 5,
      "name": "Customer-005",
      "onu_type": "F660V6.0",
      "serial_number": "ZTEGD*******",
      "rx_power": "-19.54",
      "status": "Online"
    },
    {
      "board": 2,
      "pon": 8,
      "onu_id": 6,
      "name": "Customer-006",
      "onu_type": "F670LV7.1",
      "serial_number": "ZTEGC*******",
      "rx_power": "-21.81",
      "status": "Online"
    }
  ]
}
```

### Description of Paginate
| Syntax             | Description                                                     |
|--------------------|-----------------------------------------------------------------|
| page               | Page number                                                     |
| limit              | Limit data per page                                             |
| page_count         | Total page                                                      |
| total_rows         | Total rows                                                      |
| data               | Data of onu                                                     |

#### Default paginate
``` go
var (
	DefaultPageSize = 10 // default page size
	MaxPageSize     = 100 // max page size
	PageVar         = "page"
	PageSizeVar     = "limit"
)
```


### LICENSE
[MIT License](https://github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/blob/main/LICENSE)
