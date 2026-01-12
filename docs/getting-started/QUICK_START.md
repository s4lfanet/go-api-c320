# ⚡ Quick Start Guide

Deploy ZTE C320 OLT API dalam 5 menit! Panduan singkat untuk yang ingin langsung action.

## 🎯 Prerequisites Check (1 menit)

```bash
# ✅ Linux OS (Ubuntu/Debian/CentOS/Rocky)
cat /etc/os-release

# ✅ Root/sudo access
sudo whoami

# ✅ Internet connection
ping -c 3 google.com

# ✅ Can reach OLT device
ping -c 3 <OLT_IP_ADDRESS>
```

---

## 🚀 Installation (2 menit)

### Option 1: One-Command Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install-quickstart.sh | sudo bash
```

Installer akan:
- ✅ Auto-detect OS
- ✅ Install Go 1.25.5
- ✅ Install Redis 7.2
- ✅ Clone & build application
- ✅ Create systemd service
- ✅ Start service automatically

### Option 2: Manual Install

```bash
# Download installer
wget https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install.sh
chmod +x install.sh

# Run installer
sudo ./install.sh
```

---

## ⚙️ Configuration (1 menit)

Edit konfigurasi:

```bash
sudo nano /opt/go-snmp-olt/.env
```

**Minimal Configuration** (hanya ganti 3 baris ini):

```bash
OLT_IP_ADDRESS=136.1.1.100          # ← Ganti dengan IP OLT Anda
OLT_SNMP_COMMUNITY=public           # ← Ganti jika beda
OLT_TELNET_HOST=136.1.1.100         # ← Sama dengan OLT_IP_ADDRESS
```

**Full Configuration** (optional):

```bash
# OLT Configuration
OLT_IP_ADDRESS=136.1.1.100
OLT_SNMP_PORT=161
OLT_SNMP_COMMUNITY=public
OLT_SNMP_VERSION=2c
OLT_SNMP_TIMEOUT=10
OLT_SNMP_RETRIES=3

# Telnet Configuration
OLT_TELNET_HOST=136.1.1.100
OLT_TELNET_PORT=23
OLT_TELNET_USERNAME=zte
OLT_TELNET_PASSWORD=zte
OLT_TELNET_ENABLE_PASSWORD=zxr10
OLT_TELNET_TIMEOUT=30

# API Server
SERVER_PORT=8081
SERVER_HOST=0.0.0.0
LOG_LEVEL=info

# Redis Cache
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=OsWkRgJLabn4n2+nodZ6BQeP+OKkrObnGeFcDY6w7Nw=
REDIS_DB=0

# Firmware Version
ZTE_FIRMWARE_VERSION=v2.1
```

Save file (`Ctrl+X`, `Y`, `Enter`).

Restart service:

```bash
sudo systemctl restart go-snmp-olt
```

---

## ✅ Verification (1 menit)

### Check Service Status

```bash
sudo systemctl status go-snmp-olt
```

Harus muncul: `Active: active (running)` ✅

### Test API Endpoint

```bash
# Test root endpoint
curl http://localhost:8081/

# Expected output:
# Hello, this is the root endpoint!
```

### Test Real Endpoint

```bash
# Get ONUs on board 1, pon 1
curl http://localhost:8081/api/v1/board/1/pon/1/
```

Expected response:
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "pon_port": "1",
      "onu_id": 1,
      "serial_number": "ZTEG1234ABCD",
      "model": "ZTE-F660",
      "online_status": 1
    }
  ]
}
```

---

## 🎨 Common Commands

```bash
# Check status
sudo systemctl status go-snmp-olt

# Start service
sudo systemctl start go-snmp-olt

# Stop service
sudo systemctl stop go-snmp-olt

# Restart service
sudo systemctl restart go-snmp-olt

# View logs (live)
sudo journalctl -u go-snmp-olt -f

# View logs (last 100 lines)
sudo journalctl -u go-snmp-olt -n 100

# Check if port is listening
sudo netstat -tulpn | grep 8081

# Test Redis connection
redis-cli ping
```

---

## 📖 Next Steps

### Explore API Endpoints

**ONU Monitoring:**
```bash
# List all ONUs on PON port
curl http://localhost:8081/api/v1/board/1/pon/1/

# Get specific ONU details
curl http://localhost:8081/api/v1/board/1/pon/1/onu/5

# Real-time monitoring with optical power
curl http://localhost:8081/api/v1/monitoring/onu/1/5
```

**System Information:**
```bash
# Get all cards
curl http://localhost:8081/api/v1/system/cards/

# Get traffic profiles
curl http://localhost:8081/api/v1/profiles/traffic/
```

### Full API Documentation

See [API Reference](../features/API_REFERENCE.md) for all 50+ endpoints.

### Advanced Configuration

See [Installation Guide](INSTALLATION.md) for:
- Custom port configuration
- Nginx reverse proxy setup
- SSL/TLS configuration
- Multiple OLT support

---

## 🚨 Troubleshooting

### Service Won't Start

```bash
# Check logs
sudo journalctl -u go-snmp-olt -n 50

# Common issues:
# 1. Port already in use
sudo netstat -tulpn | grep 8081

# 2. Redis not running
sudo systemctl status redis
sudo systemctl start redis

# 3. Invalid .env configuration
sudo nano /opt/go-snmp-olt/.env
```

### Can't Reach OLT

```bash
# Test network connectivity
ping <OLT_IP>

# Test SNMP
snmpwalk -v2c -c public <OLT_IP> system

# Test Telnet
telnet <OLT_IP> 23
```

### API Returns Empty Data

```bash
# Check OLT configuration
# 1. SNMP enabled?
# 2. Community string correct?
# 3. Firewall blocking?

# Check API logs
sudo journalctl -u go-snmp-olt -f

# Test SNMP manually
snmpwalk -v2c -c public <OLT_IP> 1.3.6.1.4.1.3902.1082.500.10.2.3.3.1.2
```

---

## 🎯 Common Use Cases

### 1. Monitor All ONUs

```bash
# Get all boards and PONs
for board in {1..8}; do
  for pon in {1..16}; do
    echo "Board $board, PON $pon:"
    curl -s http://localhost:8081/api/v1/board/$board/pon/$pon/ | jq
  done
done
```

### 2. Find Specific Serial Number

```bash
# Search all PONs for serial number
curl http://localhost:8081/api/v1/board/1/pon/1/onu_id/serial | \
  jq '.data[] | select(.serial_number == "ZTEG1234ABCD")'
```

### 3. Monitor Optical Power

```bash
# Get optical power for specific ONU
curl http://localhost:8081/api/v1/monitoring/onu/1/5 | \
  jq '.data.optical'
```

---

## 📞 Get Help

- **Documentation**: [docs/README.md](../README.md)
- **Troubleshooting**: [deployment/TROUBLESHOOTING.md](../deployment/TROUBLESHOOTING.md)
- **Issues**: https://github.com/s4lfanet/go-api-c320/issues
- **Email**: wardian370@gmail.com

---

## 🎉 Success Checklist

- [ ] Service running (`systemctl status go-snmp-olt`)
- [ ] API responding (`curl http://localhost:8081/`)
- [ ] Can get ONU list (`curl .../board/1/pon/1/`)
- [ ] Redis working (`redis-cli ping`)
- [ ] Logs clean (no errors in `journalctl -u go-snmp-olt`)

**If all checked:** Congratulations! 🎊 You're ready to use the API!

---

**Installation Time**: ~5 minutes  
**Difficulty**: Easy ⭐  
**Last Updated**: January 12, 2026
