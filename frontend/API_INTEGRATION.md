# Frontend API Integration

Frontend Dashboard terintegrasi penuh dengan Backend API ZTE C320 OLT.

## ✅ Fitur yang Sudah Diimplementasikan

### 1. **Dashboard (Real-time Statistics)**

**Endpoint yang digunakan:**
- `GET /api/v1/board/{board}/pon/{pon}/` - Untuk setiap kombinasi Board 1-2 dan PON 1-16

**Data yang ditampilkan:**
- **Total ONUs**: Jumlah total ONU di semua board dan PON ports
- **Online ONUs**: Jumlah ONU yang online (status = 'online' atau '1')
- **Offline ONUs**: Jumlah ONU yang offline
- **Alerts**: ONUs dengan masalah (offline atau RX power < -27 dBm)
- **Uptime Percentage**: Persentase ONU yang online
- **Active PON Ports**: Tabel PON ports yang memiliki ONU

**Auto-refresh**: Setiap 30 detik

**Contoh response yang diproses:**
```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "board": 1,
      "pon": 8,
      "onu_id": 1,
      "name": "ONU-Customer-001",
      "serial_number": "ZTEG12345678",
      "onu_type": "F660",
      "status": "online",
      "rx_power": "-24.5",
      "gpon_optical_distance": "1250"
    }
  ]
}
```

### 2. **Monitoring Page (ONU List)**

**Endpoint yang digunakan:**
- `GET /api/v1/board/{board}/pon/{pon}/` - Berdasarkan board dan PON yang dipilih user

**Data yang ditampilkan:**
- ONU ID
- Name
- Serial Number
- ONU Type
- Status (Online/Offline dengan badge berwarna)
- RX Power (dengan color coding: merah < -27, kuning < -25, hijau >= -25)
- Optical Distance

**Fitur:**
- Selector untuk memilih Board (1-2) dan PON (1-16)
- Auto-refresh setiap 15 detik
- Manual refresh button
- Error handling untuk PON ports tanpa ONU

## 🔧 Konfigurasi API Client

**Base URL:** `http://192.168.54.230/api/v1`
- Frontend: Port 80 (Nginx)
- Backend API: Port 8081 (di-proxy oleh Nginx)

**Axios Configuration:**
```typescript
axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});
```

**Nginx Proxy Configuration:**
```nginx
location /api/ {
    proxy_pass http://localhost:8081/api/;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
}
```

## 📊 State Management

**React Query Configuration:**
- Stale Time: 5 menit
- Cache Time: 10 menit
- Retry: 3 kali
- Retry Delay: Exponential backoff

**Query Keys:**
- `dashboard-stats` - Dashboard statistics
- `['onu-list', board, pon]` - ONU list per board/PON

## ⚡ Performance

**Dashboard Load:**
- Queries: 32 API calls (2 boards × 16 PONs)
- Parallel execution dengan error handling
- Skips PON ports tanpa data (404)
- Total load time: ~2-5 detik (tergantung jumlah active PONs)

**Optimizations:**
- React Query caching
- Auto-refresh intervals
- Lazy loading untuk PON ports

## 🚀 Next Steps (Belum Diimplementasikan)

### Pages yang Perlu API Integration:

1. **Provisioning** - ONU Registration
   - `GET /api/v1/board/{board}/pon/{pon}/unconfigured`
   - `POST /api/v1/provision`

2. **VLAN Management**
   - `GET /api/v1/board/{board}/pon/{pon}/onu/{onu_id}/vlan`
   - `POST /api/v1/vlan`

3. **Traffic Control**
   - `GET /api/v1/board/{board}/pon/{pon}/onu/{onu_id}/traffic`
   - `POST /api/v1/traffic`

4. **Config Backup**
   - `GET /api/v1/config/backup`
   - `POST /api/v1/config/restore`

## 📝 Testing

**Test Endpoints:**
```bash
# Check ONU on Board 1 PON 8
curl http://192.168.54.230/api/v1/board/1/pon/8

# Check ONU on Board 2 PON 7
curl http://192.168.54.230/api/v1/board/2/pon/7

# Check specific ONU
curl http://192.168.54.230/api/v1/board/1/pon/8/onu/1
```

## 🔍 Debugging

**Browser Console:**
- React Query DevTools (development mode)
- Network tab untuk API calls
- Console logs untuk API errors

**Common Issues:**
1. **404 Not Found** - PON port tidak memiliki ONU (normal, di-skip)
2. **CORS errors** - Sudah di-handle oleh Nginx
3. **Timeout** - Backend lambat / tidak respond (10s timeout)

## 🎯 Status Codes

- **200 OK** - Data berhasil diambil
- **404 NOT_FOUND** - ONU tidak ditemukan (normal untuk empty PON)
- **500 INTERNAL_SERVER_ERROR** - Backend error
- **503 SERVICE_UNAVAILABLE** - Backend tidak running

---

**Last Updated:** 2026-01-12  
**Frontend Version:** 1.0.0  
**Backend API Version:** Compatible with ZTE C320 V2.1.0
