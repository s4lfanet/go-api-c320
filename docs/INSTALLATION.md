# Panduan Instalasi Go SNMP OLT ZTE C320 untuk VPS

Installer lengkap untuk deployment Go SNMP OLT ZTE C320 di VPS Linux tanpa menggunakan Docker.

## 📋 Prasyarat

### Sistem Operasi yang Didukung
- Ubuntu 18.04, 20.04, 22.04, 24.04
- Debian 10, 11, 12
- CentOS 7, 8
- Rocky Linux 8, 9
- Red Hat Enterprise Linux 8, 9

### Spesifikasi Minimum VPS
- **CPU**: 2 cores
- **RAM**: 512 MB (Recommended: 2 GB)
- **Storage**: 5 GB
- **Network**: Akses ke OLT ZTE C320 (Port SNMP 161)

### Akses yang Diperlukan
- Root access atau sudo privileges
- Port 8081 (atau custom port) terbuka
- Akses internet untuk download dependencies
- Akses SNMP ke OLT device (UDP port 161)

## 🚀 Cara Instalasi

### 1. Download Installer

```bash
# Download installer script
wget https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install.sh

# Atau gunakan curl
curl -O https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install.sh

# Beri permission execute
chmod +x install.sh
```

### 2. Jalankan Installer

```bash
# Install sebagai root atau dengan sudo
sudo ./install.sh install
```

### 3. Konfigurasi Saat Instalasi

Installer akan meminta informasi berikut:

```
SNMP Host (IP OLT): 192.168.1.1
SNMP Port [161]: 161
SNMP Community String [public]: your_community_string
Server Port [8081]: 8081
```

**Contoh konfigurasi:**
- **SNMP Host**: IP address dari OLT ZTE C320 Anda
- **SNMP Port**: Default 161 (kecuali Anda mengubahnya)
- **SNMP Community**: Community string SNMP dari OLT (ganti dari "public" untuk keamanan)
- **Server Port**: Port untuk API server (default 8081)

### 4. Verifikasi Instalasi

Setelah instalasi selesai, verifikasi bahwa service berjalan:

```bash
# Cek status service
systemctl status go-snmp-olt

# Test API endpoint
curl http://localhost:8081/

# View logs
journalctl -u go-snmp-olt -f
```

## 📦 Apa yang Diinstall?

Installer akan secara otomatis:

1. ✅ Mendeteksi dan install dependencies sistem
2. ✅ Install Go 1.25.5
3. ✅ Install dan konfigurasi Redis 7.2+
4. ✅ Clone repository dari GitHub
5. ✅ Build aplikasi dari source code
6. ✅ Buat user system untuk aplikasi (`olt-service`)
7. ✅ Setup systemd service
8. ✅ Konfigurasi environment variables
9. ✅ Setup logging dan log rotation
10. ✅ Konfigurasi firewall (UFW/firewalld)
11. ✅ Set proper file permissions
12. ✅ Start dan enable service

## 📍 Lokasi File dan Direktori

```
/opt/go-snmp-olt/              # Direktori aplikasi utama
├── bin/
│   └── api                     # Binary aplikasi
├── .env                        # Environment configuration
├── logs/                       # Application logs
├── cmd/                        # Source code
├── pkg/                        # Packages
└── ...

/etc/systemd/system/
└── go-snmp-olt.service         # Systemd service file

/etc/logrotate.d/
└── go-snmp-olt                 # Log rotation config

/root/.redis_password           # Redis password (jika diset)
```

## 🔧 Management Service

### Start/Stop/Restart Service

```bash
# Start service
sudo systemctl start go-snmp-olt

# Stop service
sudo systemctl stop go-snmp-olt

# Restart service
sudo systemctl restart go-snmp-olt

# Reload configuration
sudo systemctl reload go-snmp-olt

# Enable service (auto-start on boot)
sudo systemctl enable go-snmp-olt

# Disable auto-start
sudo systemctl disable go-snmp-olt
```

### Monitoring dan Logs

```bash
# Check service status
systemctl status go-snmp-olt

# View real-time logs
journalctl -u go-snmp-olt -f

# View last 100 lines
journalctl -u go-snmp-olt -n 100

# View logs since today
journalctl -u go-snmp-olt --since today

# View logs with specific date
journalctl -u go-snmp-olt --since "2025-01-01"
```

## ⚙️ Konfigurasi

### Edit Configuration

```bash
# Edit environment file
sudo nano /opt/go-snmp-olt/.env

# Setelah edit, restart service
sudo systemctl restart go-snmp-olt
```

### Environment Variables

File `.env` berisi konfigurasi berikut:

```bash
# Application Environment
APP_ENV=production

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8081
SERVER_MODE=release

# SNMP Configuration
SNMP_HOST=192.168.1.1
SNMP_PORT=161
SNMP_COMMUNITY=public

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0
REDIS_MIN_IDLE_CONNECTIONS=200
REDIS_POOL_SIZE=12000
REDIS_POOL_TIMEOUT=240

# TLS/HTTPS (Optional)
USE_TLS=false
TLS_CERT_FILE=/path/to/cert.pem
TLS_KEY_FILE=/path/to/key.pem

# CORS
CORS_ALLOWED_ORIGINS=*
```

## 🔒 Keamanan

### Rekomendasi Keamanan Production

1. **Ganti SNMP Community String**
   ```bash
   # Jangan gunakan "public" di production
   SNMP_COMMUNITY=your_secure_community_string
   ```

2. **Set Redis Password**
   ```bash
   # Password otomatis di-generate saat instalasi
   # Tersimpan di /root/.redis_password
   cat /root/.redis_password
   ```

3. **Enable HTTPS/TLS**
   ```bash
   # Generate SSL certificate (contoh dengan Let's Encrypt)
   sudo apt install certbot
   sudo certbot certonly --standalone -d your-domain.com
   
   # Edit .env
   USE_TLS=true
   TLS_CERT_FILE=/etc/letsencrypt/live/your-domain.com/fullchain.pem
   TLS_KEY_FILE=/etc/letsencrypt/live/your-domain.com/privkey.pem
   ```

4. **Konfigurasi Firewall**
   ```bash
   # UFW (Ubuntu/Debian)
   sudo ufw allow 8081/tcp
   sudo ufw enable
   
   # Firewalld (CentOS/Rocky)
   sudo firewall-cmd --permanent --add-port=8081/tcp
   sudo firewall-cmd --reload
   ```

5. **Restrict CORS**
   ```bash
   # Batasi access dari domain tertentu
   CORS_ALLOWED_ORIGINS=https://your-domain.com,https://app.your-domain.com
   ```

## 🔄 Update Aplikasi

### Update ke Versi Terbaru

```bash
# Gunakan installer dengan parameter update
sudo ./install.sh update
```

Atau manual:

```bash
# Stop service
sudo systemctl stop go-snmp-olt

# Pull latest code
cd /opt/go-snmp-olt
sudo git pull origin main

# Rebuild
sudo /usr/local/go/bin/go build -ldflags="-s -w" -o bin/api ./cmd/api

# Restart service
sudo systemctl start go-snmp-olt
```

## 🗑️ Uninstall

### Hapus Aplikasi Sepenuhnya

```bash
# Jalankan uninstaller
sudo ./install.sh uninstall
```

Atau manual:

```bash
# Stop dan disable service
sudo systemctl stop go-snmp-olt
sudo systemctl disable go-snmp-olt

# Hapus service file
sudo rm /etc/systemd/system/go-snmp-olt.service
sudo systemctl daemon-reload

# Hapus aplikasi
sudo rm -rf /opt/go-snmp-olt

# Hapus user
sudo userdel olt-service

# Hapus logrotate config
sudo rm /etc/logrotate.d/go-snmp-olt
```

**Note:** Uninstall TIDAK akan menghapus Redis atau Go. Jika ingin hapus:

```bash
# Hapus Redis
sudo apt remove redis-server  # Ubuntu/Debian
sudo yum remove redis          # CentOS/Rocky

# Hapus Go
sudo rm -rf /usr/local/go
```

## 🧪 Testing API

### Test Endpoint Dasar

```bash
# Root endpoint
curl http://localhost:8081/

# Get ONU information (Board 2, PON 7)
curl http://localhost:8081/api/v1/board/2/pon/7

# Get specific ONU (Board 2, PON 7, ONU 4)
curl http://localhost:8081/api/v1/board/2/pon/7/onu/4

# Get empty ONU IDs (Board 2, PON 5)
curl http://localhost:8081/api/v1/board/2/pon/5/onu_id/empty

# With pagination
curl 'http://localhost:8081/api/v1/paginate/board/2/pon/8?limit=10&page=1'
```

### Test dengan jq untuk format JSON

```bash
# Install jq
sudo apt install jq  # Ubuntu/Debian
sudo yum install jq  # CentOS/Rocky

# Query dengan format yang rapi
curl -s http://localhost:8081/api/v1/board/2/pon/7 | jq
```

## 📊 Monitoring dan Performance

### Resource Monitoring

```bash
# Monitor CPU dan Memory usage
top -p $(pgrep -f go-snmp-olt)

# Atau gunakan htop
htop -p $(pgrep -f go-snmp-olt)

# Check port listening
sudo netstat -tulpn | grep 8081
# Atau
sudo ss -tulpn | grep 8081
```

### Redis Monitoring

```bash
# Connect ke Redis CLI
redis-cli

# Auth jika menggunakan password
AUTH your_redis_password

# Monitor Redis
MONITOR

# Info Redis
INFO

# Check keys
KEYS *

# Exit
exit
```

## 🐛 Troubleshooting

### Service Tidak Start

```bash
# Check logs detail
journalctl -u go-snmp-olt -xe

# Check konfigurasi
cat /opt/go-snmp-olt/.env

# Verify binary
ls -lah /opt/go-snmp-olt/bin/api

# Test manual
cd /opt/go-snmp-olt
sudo -u olt-service ./bin/api
```

### Connection Error ke OLT

```bash
# Test SNMP connectivity
sudo apt install snmp  # Ubuntu/Debian
sudo yum install net-snmp-utils  # CentOS/Rocky

# Test SNMP walk
snmpwalk -v2c -c your_community 192.168.1.1 system
```

### Redis Connection Error

```bash
# Check Redis status
systemctl status redis-server  # Ubuntu/Debian
systemctl status redis          # CentOS/Rocky

# Test Redis connection
redis-cli ping

# Check Redis logs
journalctl -u redis-server -n 50
```

### Port Already in Use

```bash
# Check apa yang menggunakan port 8081
sudo lsof -i :8081

# Kill process jika diperlukan
sudo kill -9 <PID>

# Atau ubah port di .env
sudo nano /opt/go-snmp-olt/.env
# Ubah SERVER_PORT=8082
```

### Permission Errors

```bash
# Fix ownership
sudo chown -R olt-service:olt-service /opt/go-snmp-olt

# Fix binary permission
sudo chmod +x /opt/go-snmp-olt/bin/api

# Fix env permission
sudo chmod 600 /opt/go-snmp-olt/.env
```

## 📖 API Documentation

### Endpoints

#### 1. Health Check
```bash
GET /
```

#### 2. Get All ONUs in PON
```bash
GET /api/v1/board/{board}/pon/{pon}

# Example:
curl http://localhost:8081/api/v1/board/2/pon/7
```

#### 3. Get Specific ONU
```bash
GET /api/v1/board/{board}/pon/{pon}/onu/{onu_id}

# Example:
curl http://localhost:8081/api/v1/board/2/pon/7/onu/4
```

#### 4. Get Empty ONU IDs
```bash
GET /api/v1/board/{board}/pon/{pon}/onu_id/empty

# Example:
curl http://localhost:8081/api/v1/board/2/pon/5/onu_id/empty
```

#### 5. Update Empty ONU IDs
```bash
GET /api/v1/board/{board}/pon/{pon}/onu_id/update

# Example:
curl http://localhost:8081/api/v1/board/2/pon/5/onu_id/update
```

#### 6. Get ONUs with Pagination
```bash
GET /api/v1/paginate/board/{board}/pon/{pon}?limit={limit}&page={page}

# Example:
curl 'http://localhost:8081/api/v1/paginate/board/2/pon/8?limit=10&page=1'
```

### Response Format

Semua response dalam format JSON:

```json
{
  "code": 200,
  "status": "OK",
  "data": [...]
}
```

## 💡 Tips dan Best Practices

1. **Backup Configuration**
   ```bash
   # Backup .env sebelum update
   sudo cp /opt/go-snmp-olt/.env /opt/go-snmp-olt/.env.backup
   ```

2. **Monitor Logs Regularly**
   ```bash
   # Setup cron untuk monitoring
   # Kirim alert jika ada error
   ```

3. **Auto-restart on Failure**
   ```bash
   # Sudah dikonfigurasi di systemd service
   Restart=always
   RestartSec=10
   ```

4. **Resource Limits**
   ```bash
   # Edit systemd service jika perlu adjust
   sudo nano /etc/systemd/system/go-snmp-olt.service
   
   # Tambahkan di [Service]
   LimitNOFILE=65536
   MemoryLimit=2G
   CPUQuota=200%
   ```

5. **Reverse Proxy dengan Nginx**
   ```bash
   # Install nginx
   sudo apt install nginx
   
   # Konfigurasi
   sudo nano /etc/nginx/sites-available/go-snmp-olt
   ```
   
   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       
       location / {
           proxy_pass http://localhost:8081;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```
   
   ```bash
   sudo ln -s /etc/nginx/sites-available/go-snmp-olt /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl restart nginx
   ```

## 📞 Support dan Dokumentasi

- **Repository**: https://github.com/s4lfanet/go-api-c320
- **Issues**: https://github.com/s4lfanet/go-api-c320/issues
- **Documentation**: https://github.com/s4lfanet/go-api-c320/blob/main/GUIDES.md

## 📝 License

MIT License - See LICENSE file in repository

## ✨ Credits

Developed by [Cepat Kilat Teknologi](https://ckt.co.id/)

---

**Happy Monitoring! 🚀**
