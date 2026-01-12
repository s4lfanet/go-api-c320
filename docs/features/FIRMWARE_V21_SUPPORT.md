# ZTE C320 Firmware V2.1.0 Support

## Perubahan yang dilakukan

Aplikasi ini telah dimodifikasi untuk mendukung **ZTE C320 dengan firmware V2.1.0** yang menggunakan struktur OID berbeda dari firmware V2.2.x ke atas.

## Perbedaan OID Structure

### Base OID
| Firmware | Base OID |
|----------|----------|
| V2.1.x   | `.1.3.6.1.4.1.3902.1082` |
| V2.2.x+  | `.1.3.6.1.4.1.3902.1082` |

Base OID sama, tetapi sub-OID untuk beberapa parameter berbeda.

### OID Prefix Differences

| Parameter | V2.1.x | V2.2.x+ |
|-----------|--------|---------|
| ONU Type | `.500.10.2.3.3.1.5` | `.3.50.11.2.1.17` |
| TX Power | `.500.20.2.2.2.1.11` | `.3.50.12.1.1.14` |
| IP Address | `.500.10.2.3.3.1.16` | `.3.50.16.1.1.10` |

## Cara Menggunakan

### 1. Set Firmware Version via Environment Variable

```bash
# Untuk firmware V2.1.x (default)
export ZTE_FIRMWARE_VERSION=v2.1

# Untuk firmware V2.2.x dan lebih baru
export ZTE_FIRMWARE_VERSION=v2.2
```

### 2. Konfigurasi di file .env

```env
# ZTE C320 Firmware Version
# v2.1 = Firmware V2.1.x (struktur OID berbeda)
# v2.2 = Firmware V2.2.x dan lebih baru
ZTE_FIRMWARE_VERSION=v2.1

# SNMP Configuration
SNMP_HOST=192.168.1.1
SNMP_PORT=161
SNMP_COMMUNITY=public
```

### 3. Custom OID Override (Advanced)

Jika OLT Anda memiliki struktur OID yang unik, Anda bisa override individual OID:

```env
# Override specific OIDs
OLT_BASE_OID=.1.3.6.1.4.1.3902.1082
ONU_ID_NAME_PREFIX=.500.10.2.3.3.1.2
ONU_TYPE_PREFIX=.500.10.2.3.3.1.5
ONU_SERIAL_NUMBER_PREFIX=.500.10.2.3.3.1.18
ONU_RX_POWER_PREFIX=.500.20.2.2.2.1.10
ONU_TX_POWER_PREFIX=.500.20.2.2.2.1.11
ONU_STATUS_ID_PREFIX=.500.10.2.3.8.1.4
```

## Verifikasi OID di OLT Anda

Untuk memastikan OID yang benar untuk OLT Anda, jalankan:

```bash
# Check device info
snmpwalk -v2c -c public <IP_OLT> 1.3.6.1.2.1.1.1.0

# Check GPON ONU entries
snmpwalk -v2c -c public <IP_OLT> 1.3.6.1.4.1.3902.1082.500.10.2.3.3.1.2

# Check available GPON MIBs
snmpwalk -v2c -c public <IP_OLT> 1.3.6.1.4.1.3902.1082.500
```

## MIB Groups yang Didukung V2.1.0

Berdasarkan SNMP walk, V2.1.0 mendukung:

- `zxAnGponProfileMgmtGroup` - Profile management
- `zxAnGponOltMgmtGroup` - OLT management
- `zxAnGponOltPerfMgmtGroup` - OLT performance
- `zxAnGponOntMgmtGroup` - ONT management
- `zxAnGponOntPerfMgmtGroup` - ONT performance
- `zxAnGponTrapGroup` - GPON traps
- `zxAnGponRmOntEquipGroup` - Remote ONT equipment
- `zxAnGponRmOntAniGroup` - Remote ONT ANI
- `zxAnGponRmOntEthernetGroup` - Remote ONT Ethernet

## Troubleshooting

### API returns 404 "ONU info not found"

1. **Verifikasi firmware version sudah benar:**
   ```bash
   grep ZTE_FIRMWARE_VERSION /opt/go-snmp-olt/.env
   ```

2. **Restart service setelah perubahan:**
   ```bash
   sudo systemctl restart go-snmp-olt
   ```

3. **Clear Redis cache:**
   ```bash
   redis-cli -a $(cat /root/.redis_password) FLUSHALL
   ```

4. **Test SNMP connectivity:**
   ```bash
   snmpwalk -v2c -c public <IP_OLT> 1.3.6.1.4.1.3902.1082.500.10.2.3.3.1.2
   ```

### Tidak ada data ONU

Pastikan:
1. Board ID benar (1 atau 2)
2. PON ID benar (1-16)
3. Ada ONU yang terdaftar di PON tersebut

### Check logs

```bash
journalctl -u go-snmp-olt -f
```

## Build dari Source

```bash
# Clone repository
git clone https://github.com/s4lfanet/go-api-c320.git
cd go-snmp-olt-zte-c320

# Build
go mod download
CGO_ENABLED=0 go build -o bin/api ./cmd/api

# Run dengan V2.1 support
ZTE_FIRMWARE_VERSION=v2.1 ./bin/api
```

## Catatan Penting

1. **Firmware V2.1.0 sudah cukup lama** - Pertimbangkan untuk upgrade ke firmware yang lebih baru jika memungkinkan.

2. **OID structure bisa berbeda** per-device - Jika masih tidak work, gunakan snmpwalk untuk menemukan OID yang benar.

3. **Board/PON base values** mungkin berbeda di instalasi Anda - Gunakan environment variable untuk override jika perlu.
