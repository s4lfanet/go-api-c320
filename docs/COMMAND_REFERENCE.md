# ZTE C320 V2.1 - Complete Command Reference

**Firmware:** V2.1.0  
**Last Updated:** January 11, 2026

---

## 📚 Table of Contents

1. [Authentication & Navigation](#authentication--navigation)
2. [ONU Registration Commands](#onu-registration-commands)
3. [ONU Configuration Commands](#onu-configuration-commands)
4. [VLAN Commands](#vlan-commands)
5. [Traffic Profile Commands](#traffic-profile-commands)
6. [Service Port Commands](#service-port-commands)
7. [ONU Management Commands](#onu-management-commands)
8. [Show Commands (Diagnostic)](#show-commands-diagnostic)
9. [System Commands](#system-commands)
10. [Troubleshooting Commands](#troubleshooting-commands)

---

## 🔐 Authentication & Navigation

### Login Sequence
```
telnet 136.1.1.100

Username: admin
Password: *****

ZXAN>
```

### Mode Navigation
```bash
# User Mode → Enable Mode
ZXAN> enable
Password: *****
ZXAN#

# Enable Mode → Configure Mode
ZXAN# configure terminal
ZXAN(config)#

# Configure Mode → Interface Mode
ZXAN(config)# interface gpon-olt_1/1/1
ZXAN(config-if)#

# Configure Mode → ONU Interface Mode
ZXAN(config)# interface gpon-onu_1/1/1:1
ZXAN(gpon-onu-mng)#

# Exit one level
exit

# Exit to Enable Mode
end
```

### Prompt Reference
| Prompt | Mode | Description |
|--------|------|-------------|
| `ZXAN>` | User Mode | Read-only, limited commands |
| `ZXAN#` | Enable Mode | Read-only, full show commands |
| `ZXAN(config)#` | Configure Mode | Global configuration |
| `ZXAN(config-if)#` | Interface Mode | Port-level configuration |
| `ZXAN(gpon-onu-mng)#` | ONU Management | ONU-specific configuration |
| `ZXAN(config-vlan)#` | VLAN Mode | VLAN configuration |

---

## 📝 ONU Registration Commands

### Show Unconfigured ONUs
```bash
# Global
ZXAN# show gpon onu uncfg

# Per PON port
ZXAN# show gpon onu uncfg gpon-olt_1/1/1

# Output example:
# OltId        OnuId    Serial-Number   Password   Loid
# gpon-olt_1/1/1   N/A    ZTEGDA5918AC   N/A        N/A
```

### Register ONU
```bash
# Enter interface
ZXAN(config)# interface gpon-olt_1/1/1

# Register ONU with type and serial number
ZXAN(config-if)# onu 1 type ZTE-F670L sn ZTEGDA5918AC

# Register with specific ONU ID
ZXAN(config-if)# onu 5 type ZTE-F660 sn ZTEGD824CDF3

# Set ONU name/description
ZXAN(config-if)# onu 1 name "Pelanggan_Rumah_001"
```

### ONU Type List (Common)
```
ZTE-F601       - 1 GE + WiFi
ZTE-F609       - 4 GE + WiFi
ZTE-F660       - 4 GE + 2 POTS + WiFi
ZTE-F670L      - 4 GE + 2 POTS + WiFi (Dual Band)
ZTE-F680       - 4 GE + 2 POTS + WiFi (Gigabit)
Huawei-EG8145V5 - 4 GE + 2 POTS + WiFi
Huawei-EG8041V5 - 1 GE
FiberHome-AN5506-04-F - 4 GE + 2 POTS + WiFi
```

### Delete ONU
```bash
ZXAN(config)# interface gpon-olt_1/1/1
ZXAN(config-if)# no onu 1
```

---

## ⚙️ ONU Configuration Commands

### Enter ONU Configuration
```bash
ZXAN(config)# interface gpon-onu_1/1/1:1
ZXAN(gpon-onu-mng)#
```

### TCONT Configuration
```bash
# Create TCONT with DBA profile
ZXAN(gpon-onu-mng)# tcont 1 name TCONT_DATA profile UP-10M

# Multiple TCONTs
ZXAN(gpon-onu-mng)# tcont 1 name TCONT_DATA profile UP-10M
ZXAN(gpon-onu-mng)# tcont 2 name TCONT_VOIP profile UP-VOIP
ZXAN(gpon-onu-mng)# tcont 3 name TCONT_MGMT profile UP-MNG
```

### GEMPORT Configuration
```bash
# Create GEMPORT linked to TCONT
ZXAN(gpon-onu-mng)# gemport 1 name GEM_DATA tcont 1

# GEMPORT with traffic management
ZXAN(gpon-onu-mng)# gemport 1 name GEM_DATA tcont 1 queue 1
ZXAN(gpon-onu-mng)# gemport 2 name GEM_VOIP tcont 2 queue 2
```

### ONU VLAN Configuration
```bash
# Service port on ONU (VLAN tagging)
ZXAN(gpon-onu-mng)# service-port 1 vport 1 user-vlan 100 vlan 100

# Dengan translation
ZXAN(gpon-onu-mng)# service-port 1 vport 1 user-vlan untagged vlan 100

# Multiple services
ZXAN(gpon-onu-mng)# service-port 1 vport 1 user-vlan 100 vlan 100  # Internet
ZXAN(gpon-onu-mng)# service-port 2 vport 2 user-vlan 200 vlan 200  # VOIP
ZXAN(gpon-onu-mng)# service-port 3 vport 3 user-vlan 300 vlan 300  # IPTV
```

### ONU Port Configuration
```bash
# Enable/disable ONU port
ZXAN(gpon-onu-mng)# onu port eth 1 state enable
ZXAN(gpon-onu-mng)# onu port eth 1 state disable

# Set port speed
ZXAN(gpon-onu-mng)# onu port eth 1 speed 1000 duplex full

# POTS configuration
ZXAN(gpon-onu-mng)# onu port pots 1 state enable
```

---

## 🏷️ VLAN Commands

### Create VLAN
```bash
ZXAN(config)# vlan 100
ZXAN(config-vlan)# name INTERNET
ZXAN(config-vlan)# exit

# Quick create
ZXAN(config)# vlan 100 name INTERNET
```

### Delete VLAN
```bash
ZXAN(config)# no vlan 100
```

### VLAN Description
```bash
ZXAN(config)# vlan 100
ZXAN(config-vlan)# description "VLAN untuk layanan Internet"
```

### Show VLAN
```bash
ZXAN# show vlan all
ZXAN# show vlan 100
```

### PON Port VLAN
```bash
ZXAN(config)# interface gpon-olt_1/1/1

# Add VLAN to port
ZXAN(config-if)# port vlan 100

# Remove VLAN from port
ZXAN(config-if)# no port vlan 100

# Set native VLAN
ZXAN(config-if)# port native-vlan 1

# VLAN mode
ZXAN(config-if)# port vlan-mode tag
ZXAN(config-if)# port vlan-mode untag
```

---

## 📊 Traffic Profile Commands

### Create DBA Profile
```bash
ZXAN(config)# gpon-onu-profile dba-profile UP-10M
ZXAN(config-dba-profile)# type 4 assure 10240 max 10240
ZXAN(config-dba-profile)# exit
```

### DBA Profile Types
```
Type 1: Fixed bandwidth
  - Guaranteed bandwidth, always reserved
  - type 1 fix <bandwidth>

Type 2: Assured bandwidth
  - Guaranteed minimum
  - type 2 assure <bandwidth>

Type 3: Assured + Maximum
  - Guaranteed min with burst to max
  - type 3 assure <min> max <max>

Type 4: Maximum bandwidth
  - Best effort up to max
  - type 4 max <bandwidth>

Type 5: Assured + Maximum (with priority)
  - type 5 assure <min> max <max>
```

### DBA Profile Examples
```bash
# 10M Symmetric
gpon-onu-profile dba-profile UP-10M
type 4 assure 10240 max 10240
exit

# 50M Burst
gpon-onu-profile dba-profile UP-50M
type 3 assure 10240 max 51200
exit

# 100M Best Effort
gpon-onu-profile dba-profile UP-100M
type 4 max 102400
exit

# 1G Unlimited
gpon-onu-profile dba-profile UP-1G
type 4 max 1048576
exit

# VOIP Priority
gpon-onu-profile dba-profile UP-VOIP
type 1 fix 1024
exit
```

### Modify DBA Profile
```bash
ZXAN(config)# gpon-onu-profile dba-profile UP-10M
ZXAN(config-dba-profile)# no type
ZXAN(config-dba-profile)# type 4 assure 20480 max 20480
ZXAN(config-dba-profile)# exit
```

### Delete DBA Profile
```bash
ZXAN(config)# no gpon-onu-profile dba-profile UP-10M
```

### Show DBA Profiles
```bash
ZXAN# show gpon-onu-profile dba-profile
ZXAN# show gpon-onu-profile dba-profile UP-10M
```

---

## 🔌 Service Port Commands

### Create Service Port (Global)
```bash
# Basic service port
ZXAN(config)# service-port 1 gpon 1/1/1 onu 1 gemport 1 cos 0 vlan 100

# With user-vlan translation
ZXAN(config)# service-port 1 gpon 1/1/1 onu 1 gemport 1 cos 0 user-vlan 100 vlan 100

# Untagged user traffic
ZXAN(config)# service-port 1 gpon 1/1/1 onu 1 gemport 1 cos 0 user-vlan untagged vlan 100
```

### Service Port with QoS
```bash
# With specific CoS
ZXAN(config)# service-port 1 gpon 1/1/1 onu 1 gemport 1 cos 5 vlan 100

# With rate limit
ZXAN(config)# service-port 1 gpon 1/1/1 onu 1 gemport 1 cos 0 vlan 100 rx-cttr 10 tx-cttr 10
```

### Delete Service Port
```bash
ZXAN(config)# no service-port 1
```

### Show Service Ports
```bash
ZXAN# show service-port all
ZXAN# show service-port gpon 1/1/1
ZXAN# show service-port 1
```

---

## 🔧 ONU Management Commands

### Reboot ONU
```bash
ZXAN(config)# interface gpon-olt_1/1/1
ZXAN(config-if)# onu reset 1
```

### Block/Unblock ONU
```bash
# Block ONU (disable)
ZXAN(config)# interface gpon-olt_1/1/1
ZXAN(config-if)# onu 1 state disable

# Unblock ONU (enable)
ZXAN(config-if)# onu 1 state enable
```

### Change ONU Type
```bash
ZXAN(config)# interface gpon-olt_1/1/1
ZXAN(config-if)# onu 1 type ZTE-F680
```

### ONU Distance Limit
```bash
ZXAN(config)# interface gpon-olt_1/1/1
ZXAN(config-if)# pon-distance-limit 40
```

---

## 📋 Show Commands (Diagnostic)

### ONU Status
```bash
# All ONUs on PON
ZXAN# show gpon onu state gpon-olt_1/1/1

# Specific ONU detail
ZXAN# show gpon onu detail-info gpon-olt_1/1/1 1

# ONU running config
ZXAN# show gpon onu running-config gpon-onu_1/1/1:1
```

### ONU Optical Power
```bash
ZXAN# show gpon onu optical-info gpon-olt_1/1/1 1
```

### ONU Statistics
```bash
ZXAN# show gpon onu traffic gpon-olt_1/1/1 1
```

### PON Port Status
```bash
ZXAN# show gpon onu-info gpon-olt_1/1/1
ZXAN# show interface gpon-olt_1/1/1
```

### Running Configuration
```bash
# Full running config
ZXAN# show running-config

# Interface specific
ZXAN# show running-config interface gpon-olt_1/1/1

# ONU specific
ZXAN# show running-config interface gpon-onu_1/1/1:1
```

### System Info
```bash
ZXAN# show system-info
ZXAN# show card all
ZXAN# show version
ZXAN# show cpu-usage
ZXAN# show memory-usage
```

---

## 💾 System Commands

### Save Configuration
```bash
# Save running to startup
ZXAN# write

# Alternative
ZXAN# copy running-config startup-config
```

### Backup Configuration
```bash
# Save to flash
ZXAN# copy running-config flash:backup_20260111.cfg

# Show saved configs
ZXAN# dir flash:
```

### Restore Configuration
```bash
ZXAN# copy flash:backup_20260111.cfg running-config
```

### System Reboot
```bash
ZXAN# reboot

# Reboot specific card
ZXAN# reboot slot 1
```

---

## 🔍 Troubleshooting Commands

### PON Link Issues
```bash
# Check PON transceiver
ZXAN# show interface gpon-olt_1/1/1 optical-info

# Check ONU registration
ZXAN# show gpon onu uncfg gpon-olt_1/1/1

# Check ONU optical
ZXAN# show gpon onu optical-info gpon-olt_1/1/1 1
```

### Traffic Issues
```bash
# Check service port stats
ZXAN# show service-port traffic 1

# Check GEMPORT stats
ZXAN# show gpon onu gemport-info gpon-olt_1/1/1 1

# Check TCONT stats
ZXAN# show gpon onu tcont-info gpon-olt_1/1/1 1
```

### VLAN Issues
```bash
# Check VLAN on port
ZXAN# show interface gpon-olt_1/1/1 vlan

# Check service VLAN
ZXAN# show service-port 1
```

### Debug Commands
```bash
# Enable debug
ZXAN# debug gpon onu

# Show debug log
ZXAN# show log

# Disable debug
ZXAN# no debug all
```

---

## 📐 Complete ONU Provisioning Example

### Scenario: Add new Internet customer
```bash
# 1. Login dan masuk configure mode
enable
configure terminal

# 2. Cek ONU yang belum terdaftar
show gpon onu uncfg gpon-olt_1/1/1

# 3. Daftarkan ONU
interface gpon-olt_1/1/1
onu 1 type ZTE-F670L sn ZTEGDA5918AC
onu 1 name "Customer_Internet_001"
exit

# 4. Konfigurasi ONU
interface gpon-onu_1/1/1:1

# 5. Setup TCONT dengan profile bandwidth
tcont 1 name TCONT_DATA profile UP-10M

# 6. Setup GEMPORT
gemport 1 name GEM_DATA tcont 1

# 7. Setup service VLAN pada ONU
service-port 1 vport 1 user-vlan untagged vlan 100
exit

# 8. Buat service port global
service-port 1 gpon 1/1/1 onu 1 gemport 1 cos 0 user-vlan untagged vlan 100

# 9. Simpan konfigurasi
end
write

# 10. Verifikasi
show gpon onu state gpon-olt_1/1/1
show service-port 1
```

### Scenario: Internet + VOIP Customer
```bash
# After ONU registration...
interface gpon-onu_1/1/1:1

# Internet TCONT & GEMPORT
tcont 1 name TCONT_DATA profile UP-20M
gemport 1 name GEM_DATA tcont 1
service-port 1 vport 1 user-vlan 100 vlan 100

# VOIP TCONT & GEMPORT
tcont 2 name TCONT_VOIP profile UP-VOIP
gemport 2 name GEM_VOIP tcont 2
service-port 2 vport 2 user-vlan 200 vlan 200
exit

# Global service ports
service-port 10 gpon 1/1/1 onu 1 gemport 1 cos 0 vlan 100
service-port 11 gpon 1/1/1 onu 1 gemport 2 cos 5 vlan 200

end
write
```

---

## ⚠️ Common Errors & Solutions

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `% Invalid input detected` | Typo atau syntax salah | Periksa syntax command |
| `% Incomplete command` | Parameter kurang | Tambah parameter yang diperlukan |
| `% ONU not exist` | ONU belum terdaftar | Register ONU terlebih dahulu |
| `% Profile not exist` | Profile belum dibuat | Buat profile terlebih dahulu |
| `% VLAN not exist` | VLAN belum dibuat | Buat VLAN terlebih dahulu |
| `% Service port already exist` | ID sudah dipakai | Gunakan ID lain |
| `% ONU is offline` | ONU tidak online | Cek fisik/optical |

---

## 🔗 Quick Reference Card

```
╔═══════════════════════════════════════════════════════════════════╗
║                    ZTE C320 QUICK REFERENCE                       ║
╠═══════════════════════════════════════════════════════════════════╣
║ REGISTER ONU:                                                     ║
║   interface gpon-olt_1/1/1                                        ║
║   onu <id> type <type> sn <serial>                               ║
║                                                                   ║
║ CONFIGURE ONU:                                                    ║
║   interface gpon-onu_1/1/1:<onu_id>                              ║
║   tcont 1 name TCONT profile <profile>                           ║
║   gemport 1 name GEM tcont 1                                     ║
║   service-port 1 vport 1 user-vlan <vlan> vlan <vlan>           ║
║                                                                   ║
║ SERVICE PORT:                                                     ║
║   service-port <id> gpon 1/1/1 onu <id> gemport 1 cos 0 vlan <v>║
║                                                                   ║
║ DELETE ONU:                                                       ║
║   interface gpon-olt_1/1/1                                        ║
║   no onu <id>                                                     ║
║                                                                   ║
║ SAVE CONFIG:                                                      ║
║   write                                                           ║
╚═══════════════════════════════════════════════════════════════════╝
```
