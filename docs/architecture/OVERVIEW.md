# 🏗️ System Architecture Overview

Comprehensive overview of ZTE C320 OLT API architecture, components, and design decisions.

## 📐 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Client Applications                      │
│            (Web UI, Mobile App, Integration Services)        │
└───────────────────────────┬───────────────────────────────────┘
                            │ HTTP/REST
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway Layer                        │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ Chi Router (Go)                                       │   │
│  │ - Request Routing                                     │   │
│  │ - Middleware (CORS, Auth, Rate Limit, Logging)       │   │
│  └──────────────────────────────────────────────────────┘   │
└───────────────────────────┬───────────────────────────────────┘
                            │
                ┌───────────┴───────────┐
                │                       │
                ▼                       ▼
┌──────────────────────┐     ┌──────────────────────┐
│   Handler Layer      │     │   Middleware Layer   │
│  - OnuHandler        │     │  - RequestID         │
│  - PonHandler        │     │  - Logger            │
│  - VlanHandler       │     │  - RateLimiter       │
│  - TrafficHandler    │     │  - SecurityHeaders   │
│  - MonitoringHandler │     │  - RequestTimeout    │
└──────────┬───────────┘     └──────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Business Logic Layer                     │
│  (Usecase Layer)                                             │
│  ┌────────────────┐  ┌────────────────┐  ┌───────────────┐  │
│  │ ONU Usecase    │  │ VLAN Usecase   │  │ Monitoring    │  │
│  ├────────────────┤  ├────────────────┤  ├───────────────┤  │
│  │ Provision UC   │  │ Traffic UC     │  │ Batch UC      │  │
│  ├────────────────┤  ├────────────────┤  ├───────────────┤  │
│  │ ONUMgmt UC     │  │ ConfigBackup   │  │ Profile UC    │  │
│  └────────────────┘  └────────────────┘  └───────────────┘  │
└───────────┬──────────────────────┬───────────────────────────┘
            │                      │
    ┌───────┴────────┐    ┌───────┴────────┐
    │                │    │                │
    ▼                ▼    ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│   SNMP Repo  │  │ Telnet Repo  │  │  Redis Cache │
│              │  │              │  │              │
│ - GoSNMP     │  │ - Session    │  │ - go-redis   │
│ - ONU Data   │  │   Manager    │  │ - TTL Cache  │
│ - PON Info   │  │ - Commands   │  │ - Distributed│
│ - Card Info  │  │ - Config     │  │              │
└──────┬───────┘  └──────┬───────┘  └──────────────┘
       │                 │
       │                 │
       ▼                 ▼
┌─────────────────────────────────────────┐
│         ZTE C320 OLT Device             │
│  ┌──────────────┐  ┌──────────────┐    │
│  │ SNMP Agent   │  │ Telnet CLI   │    │
│  │ Port: 161    │  │ Port: 23     │    │
│  │ v2c Public   │  │ zte/zte      │    │
│  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────┘
```

## 🔧 Core Components

### 1. API Gateway Layer

**Technology**: Chi Router (Go)

**Responsibilities**:
- HTTP request routing
- Middleware orchestration
- Request/response formatting
- Error handling

**Key Features**:
- RESTful routing
- Versioning support (`/api/v1`)
- Middleware chaining
- Context propagation

### 2. Handler Layer

**Pattern**: HTTP Handlers (Controllers)

**Components**:
- `OnuHandler` - ONU monitoring endpoints
- `PonHandler` - PON port information
- `ProfileHandler` - Profile management
- `VlanHandler` - VLAN configuration
- `TrafficHandler` - Traffic profiles
- `ProvisionHandler` - ONU provisioning
- `ONUManagementHandler` - ONU lifecycle
- `BatchOperationsHandler` - Bulk operations
- `ConfigBackupHandler` - Backup/restore
- `MonitoringHandler` - Real-time monitoring

**Responsibilities**:
- Request validation
- Response formatting
- Error handling
- HTTP status codes

### 3. Business Logic Layer (Usecase)

**Pattern**: Clean Architecture Use Cases

**Components**:
- `OnuUsecase` - ONU data processing
- `PonUsecase` - PON port logic
- `VLANUsecase` - VLAN management
- `TrafficUsecase` - Traffic profile logic
- `ProvisionUsecase` - Provisioning workflow
- `ONUManagementUsecase` - ONU operations
- `BatchOperationsUsecase` - Batch processing
- `ConfigBackupUsecase` - Configuration management
- `MonitoringUsecase` - Real-time data aggregation

**Responsibilities**:
- Business rules enforcement
- Data transformation
- Cross-cutting concerns
- Transaction coordination

### 4. Repository Layer

**Pattern**: Repository Pattern

#### SNMP Repository
- `SNMPRepository` - SNMP data access
- `OnuRepository` - ONU SNMP queries
- `PonRepository` - PON SNMP queries
- `CardRepository` - Card/system info

**Technology**: GoSNMP
**Operations**:
- SNMP Walk
- SNMP Get
- OID parsing
- Data mapping

#### Telnet Repository
- `TelnetSessionManager` - Session pooling
- `TelnetCommands` - Command execution
- `TelnetOptical` - Optical power queries
- `TelnetProvision` - Configuration commands

**Technology**: github.com/ziutek/telnet
**Operations**:
- Session management
- Command execution
- Response parsing
- Error recovery

### 5. Caching Layer

**Technology**: Redis (go-redis/v9)

**Strategy**:
- Read-through caching
- TTL-based expiration
- Cache invalidation
- Distributed caching

**Cache Keys**:
```
onu:board:{board}:pon:{pon}           TTL: 5 minutes
onu:board:{board}:pon:{pon}:onu:{id}  TTL: 5 minutes
pon:info:{board}:{pon}                TTL: 10 minutes
profile:traffic:{id}                  TTL: 60 minutes
profile:vlan:{id}                     TTL: 60 minutes
```

---

## 📊 Data Flow Diagrams

### ONU Monitoring Flow (SNMP)

```
Client Request
    │
    ▼
[Chi Router] → Validate params
    │
    ▼
[OnuHandler] → Parse request
    │
    ▼
[OnuUsecase] → Check cache
    │
    ├─ Cache HIT  → Return cached data
    │
    └─ Cache MISS → Query OLT
                    │
                    ▼
              [SNMP Repository]
                    │
                    ▼
              [ZTE C320 OLT]
                    │
                    ▼
              Parse OID response
                    │
                    ▼
              Store in cache
                    │
                    ▼
              Return to client
```

### ONU Provisioning Flow (Telnet)

```
Client Request (POST /onu/register)
    │
    ▼
[Chi Router] → Validate JSON body
    │
    ▼
[ProvisionHandler] → Parse request
    │
    ▼
[ProvisionUsecase] → Build command sequence
    │
    ▼
[Telnet Session Manager]
    │
    ├─ Get session from pool
    │
    ├─ Execute commands:
    │  1. configure terminal
    │  2. interface gpon-olt_1/{board}/{pon}
    │  3. onu {onu_id} type {model} sn {serial}
    │  4. exit
    │  5. show gpon onu state
    │
    ├─ Verify success
    │
    └─ Return session to pool
            │
            ▼
    [Invalidate cache]
            │
            ▼
    Return success response
```

### Real-time Monitoring with Optical Power

```
Client Request (/monitoring/onu/{pon}/{onu_id})
    │
    ▼
[MonitoringHandler]
    │
    ▼
[MonitoringUsecase]
    │
    ├─ Parallel Execution ─┐
    │                       │
    ▼                       ▼
[SNMP Repository]    [Telnet Repository]
    │                       │
    ▼                       ▼
Get ONU data         Get optical power
(serial, model,      (RX/TX power,
 status, stats)       temperature, voltage)
    │                       │
    └───────┬───────────────┘
            │
            ▼
    Merge data structures
            │
            ▼
    Apply status classification
    (normal/low/high)
            │
            ▼
    Return comprehensive response
```

---

## 🛡️ Security Architecture

### Authentication & Authorization

**Current**: Open API (no auth)
**Future**: JWT-based authentication

### Rate Limiting

```go
middleware.RateLimiter(100, 200)
// 100 requests/second
// Burst: 200 requests
```

### Security Headers

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
```

### Input Validation

- Parameter validation (board, pon, onu IDs)
- JSON schema validation
- SQL injection prevention (no SQL used)
- Command injection prevention (parameterized commands)

### Request Timeout

```
Default: 90 seconds
Allows cold-cache SNMP queries (up to 60s)
```

---

## 📈 Scalability Design

### Horizontal Scaling

**Load Balancer** → Multiple API instances → Shared Redis

```
        ┌─────────────┐
        │   Nginx LB  │
        └──────┬──────┘
               │
        ┌──────┴──────┐
        │             │
        ▼             ▼
   ┌─────────┐   ┌─────────┐
   │ API #1  │   │ API #2  │
   └────┬────┘   └────┬────┘
        │             │
        └──────┬──────┘
               │
               ▼
        ┌──────────────┐
        │ Redis Cluster│
        └──────────────┘
               │
               ▼
        ┌──────────────┐
        │  ZTE C320    │
        └──────────────┘
```

### Connection Pooling

**SNMP**: Single connection, sequential queries
**Telnet**: Session pool (max 10 sessions)
**Redis**: Connection pool (default 10)

### Caching Strategy

**Levels**:
1. Application cache (in-memory)
2. Redis cache (distributed)
3. CDN (future for static responses)

---

## 🔌 Integration Points

### External Systems

1. **Monitoring Systems**
   - Prometheus metrics export
   - Grafana dashboards
   - Alert manager integration

2. **Ticketing Systems**
   - Webhook support
   - REST API integration

3. **Billing Systems**
   - Traffic data export
   - Usage statistics

4. **Inventory Management**
   - ONU registration sync
   - Serial number tracking

---

## 🧩 Design Patterns Used

### 1. Repository Pattern
```
Usecase → Repository Interface → Concrete Implementation
```
Benefits: Testability, flexibility, separation of concerns

### 2. Dependency Injection
```go
func NewOnuHandler(usecase OnuUsecaseInterface) *OnuHandler {
    return &OnuHandler{usecase: usecase}
}
```
Benefits: Loose coupling, testability

### 3. Factory Pattern
```go
func NewSNMPRepository(config *Config) *SNMPRepository
```
Benefits: Centralized creation, configuration management

### 4. Strategy Pattern
```
ONU provisioning strategies per model (F660, F609, etc.)
```

### 5. Singleton Pattern
```
Redis connection pool, Telnet session manager
```

---

## 📦 Technology Stack Decisions

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **Language** | Go 1.24+ | Performance, concurrency, single binary |
| **Router** | Chi | Lightweight, idiomatic, middleware support |
| **SNMP** | GoSNMP | Pure Go, well-maintained, feature-rich |
| **Telnet** | ziutek/telnet | Simple, reliable, adequate for CLI |
| **Cache** | Redis | Fast, distributed, TTL support |
| **Logger** | Zerolog | Zero-allocation, structured, fast |
| **Config** | Viper + godotenv | Flexible, .env support |
| **HTTP** | net/http | Standard library, proven, stable |

---

## 🎯 Performance Considerations

### Benchmarks

- **Cold cache**: 1-5 seconds (SNMP walk)
- **Warm cache**: 10-50ms (Redis)
- **Telnet command**: 500ms-2s
- **Concurrent requests**: 100 req/s sustainable

### Optimization Strategies

1. **Aggressive caching** (5-60 min TTL)
2. **Connection reuse** (pooling)
3. **Parallel queries** (goroutines)
4. **Request timeout** (90s max)
5. **Rate limiting** (prevent OLT overload)

---

## 🚀 Deployment Architecture

### Production Setup

```
Internet
    │
    ▼
┌─────────────────┐
│  Nginx (SSL)    │  Port 443
│  Reverse Proxy  │
└────────┬────────┘
         │ Port 8081
         ▼
┌─────────────────┐
│  API Service    │  systemd
│  (Go binary)    │  /opt/go-snmp-olt
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌────────┐ ┌────────┐
│ Redis  │ │OLT VPN │
└────────┘ └────────┘
```

---

## 📝 Folder Structure

```
go-api-c320/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── handler/              # HTTP handlers
│   ├── usecase/              # Business logic
│   ├── repository/           # Data access (SNMP, Telnet)
│   ├── model/                # Data structures
│   └── middleware/           # HTTP middleware
├── app/
│   ├── app.go                # Application setup
│   └── routes.go             # Route configuration
├── config/
│   └── config.go             # Configuration management
├── pkg/
│   └── utils/                # Shared utilities
├── docs/                     # Documentation
└── scripts/                  # Installation scripts
```

---

## 🔄 Future Enhancements

1. **Authentication & Authorization**
   - JWT tokens
   - Role-based access control (RBAC)
   - API key management

2. **WebSocket Support**
   - Real-time ONU status updates
   - Live monitoring dashboards

3. **Multi-OLT Support**
   - Multiple OLT management
   - Aggregated monitoring
   - Cross-OLT operations

4. **Advanced Analytics**
   - Traffic pattern analysis
   - Predictive maintenance
   - Anomaly detection

5. **Microservices Architecture**
   - Separate services per domain
   - Event-driven communication
   - Service mesh (Istio)

---

**Last Updated**: January 12, 2026  
**Version**: 1.7.2  
**Architecture**: Monolithic → Microservices-ready
