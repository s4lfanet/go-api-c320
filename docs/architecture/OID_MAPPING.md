# ZTE C320 V2.1.0 OID Mapping

Dokumen ini berisi hasil SNMP walk dan pemetaan OID untuk ZTE C320 firmware V2.1.0.

## Struktur OID Utama

ZTE C320 V2.1.0 menggunakan base OID yang berbeda dari V2.2:
- **Base OID**: `1.3.6.1.4.1.3902.1012` (bukan 1082 untuk data ONU)
- **GPON Base**: `1.3.6.1.4.1.3902.1012.3`

## PON Port Indexing

PON Index menggunakan format: `268501248 + (PON-1) * 256`
- PON 1 Board 1: 268501248
- PON 2 Board 1: 268501504 
- PON 3 Board 1: 268501760
- ... dst
- PON 16 Board 1: 268505088

## OID Tabel yang Ditemukan

### 1. PON Port Info - `1.3.6.1.4.1.3902.1012.3.13.1.1`
| OID Suffix | Tipe | Deskripsi |
|------------|------|-----------|
| .1.{pon_id} | STRING | PON Name (misal: "OLT-1") |
| .2.{pon_id} | STRING | PON Description |
| .3.{pon_id} | INTEGER | PON Status |
| .12.{pon_id} | INTEGER | Max ONU per PON (128) |

### 2. ONU Table - `1.3.6.1.4.1.3902.1012.3.13.3.1`
Format: `.X.{pon_id}.{onu_id}`

| OID Suffix X | Tipe | Deskripsi |
|--------------|------|-----------|
| .2 | Hex-STRING | ONU Serial Number (raw bytes) |
| .3 | STRING | ONU Password |
| .5 | STRING | ONU Device SN (misal: "GD824CDF3") |
| .10 | STRING | ONU Type/Model (misal: "F672YV9.1", "F670LV9.0") |
| .11 | STRING | ONU Firmware Version |

### 3. ONU Statistics - `1.3.6.1.4.1.3902.1012.3.31.4.1`
Format: `.X.{pon_id}.{onu_id}`

| OID Suffix X | Tipe | Deskripsi |
|--------------|------|-----------|
| .3 | Counter64 | Packets received |
| .6 | Counter64 | Bytes received |
| .100 | INTEGER | ONU Online Status (1=online) |

### 4. PON Port Statistics - `1.3.6.1.4.1.3902.1012.3.31.5.1`
Format: `.X.{pon_id}.1`

| OID Suffix X | Tipe | Deskripsi |
|--------------|------|-----------|
| .3 | Counter64 | PON port received packets |
| .6 | Counter64 | PON port received bytes |
| .100 | INTEGER | Status |

### 5. Traffic Profiles - `1.3.6.1.4.1.3902.1012.3.26`
| OID | Deskripsi |
|-----|-----------|
| .1.1.2.{id} | Upstream profile name |
| .2.1.2.{id} | Downstream profile name |

## ONU yang Terdeteksi

Dari SNMP walk, ditemukan 3 ONU di PON 1 (268501248):

| ONU ID | SN (Hex) | Device SN | Model | Firmware |
|--------|----------|-----------|-------|----------|
| 1 | 48 57 54 43 1F 14 CA AD | - | EG8041V5 | V5R021C00S192 |
| 2 | 5A 54 45 47 D8 24 CD F3 | GD824CDF3 | F672YV9.1 | V9.1.10P4N2 |
| 3 | 5A 54 45 47 DA 59 18 AC | GDA5918AC | F670LV9.0 | V9.0.11P2N38D |

## Perbedaan dengan V2.2

| Fitur | V2.1 OID Base | V2.2 OID Base |
|-------|---------------|---------------|
| ONU Table | 1.3.6.1.4.1.3902.1012.3.13.3.1 | 1.3.6.1.4.1.3902.1082.500.10.2.3.3.1 |
| ONU Name | Tidak tersedia (gunakan SN) | .500.10.2.3.3.1.2 |
| ONU Status | 1.3.6.1.4.1.3902.1012.3.31.4.1.100 | .500.10.2.3.3.1.4 |
| PON Info | 1.3.6.1.4.1.3902.1012.3.13.1.1 | .500.10.2.3.1.1 |

## OID Index Calculation

Untuk V2.1.0:
```
PON Index = 268500992 + (board * 8192) + (pon * 256)

Contoh Board 1:
- PON 1: 268500992 + (1 * 8192) + (1 * 256) = 268501248
- PON 2: 268500992 + (1 * 8192) + (2 * 256) = 268501504
```

## File SNMP Walk Tersimpan

Hasil SNMP walk tersimpan di VPS `/opt/go-snmp-olt/snmp-walks/`:
- `oid_1012.txt` - 3596 baris (data ONU utama)
- `oid_1015.txt` - 188 baris (system info)
- `oid_1082.txt` - 2496 baris (alarm dan config)

## Catatan Penting

1. **ONU Name tidak tersedia** - Pada V2.1.0, ONU Name tidak ada OID tersendiri. Gunakan `ONU Device SN` atau `ONU Model` sebagai identifier.

2. **RX/TX Power** - Belum ditemukan OID untuk optical power per ONU. Mungkin perlu dicari di OID lain atau tidak tersedia via SNMP di V2.1.0.

3. **ONU Run Status** - Gunakan `1.3.6.1.4.1.3902.1012.3.31.4.1.100.{pon_id}.{onu_id}` untuk cek status online (1=online).
