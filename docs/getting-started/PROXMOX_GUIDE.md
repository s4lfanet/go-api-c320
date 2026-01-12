# 🖥️ Proxmox VE Deployment Guide

Guide lengkap untuk deploy ZTE C320 OLT API di Proxmox Virtual Environment menggunakan Container (CT) atau Virtual Machine (VM).

## 📋 Prerequisites

### Proxmox VE Requirements
- Proxmox VE 7.x atau 8.x
- Akses root/sudo ke Proxmox host
- Template Ubuntu 22.04 atau Debian 12 (untuk CT)
- ISO Ubuntu 22.04 Server (untuk VM)

### Network Requirements
- Akses ke ZTE C320 OLT (IP OLT)
- Port 8081 terbuka (atau custom port)
- SNMP access ke OLT (UDP 161)
- Telnet access ke OLT (TCP 23)

### Resource Requirements

**Minimum:**
- vCPU: 1 core
- RAM: 512 MB
- Storage: 5 GB

**Recommended:**
- vCPU: 2 cores
- RAM: 2 GB
- Storage: 10 GB
- Network: Bridged (akses langsung ke OLT network)

---

## 🐧 Option 1: LXC Container (Recommended)

Container lebih ringan dan cepat dibanding VM, cocok untuk aplikasi ini.

### Step 1: Create CT Template (Jika Belum Ada)

```bash
# Login ke Proxmox shell
ssh root@proxmox-ip

# Download Ubuntu 22.04 template
pveam update
pveam download local ubuntu-22.04-standard_22.04-1_amd64.tar.zst
```

### Step 2: Create Container

**Via Web UI:**
1. Datacenter → Node → Create CT
2. General:
   - CT ID: `100` (atau available ID)
   - Hostname: `olt-api`
   - Password: (set root password)
   - Unprivileged container: ✅ (recommended)

3. Template:
   - Storage: `local`
   - Template: `ubuntu-22.04-standard`

4. Disks:
   - Disk size: `10 GB`
   - Storage: `local-lvm`

5. CPU:
   - Cores: `2`

6. Memory:
   - RAM: `2048 MB`
   - Swap: `512 MB`

7. Network:
   - Bridge: `vmbr0` (atau sesuai network OLT)
   - IPv4: `Static` atau `DHCP`
   - IPv4/CIDR: `192.168.1.100/24` (example)
   - Gateway: `192.168.1.1`

**Via CLI:**

```bash
pct create 100 local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst \
  --hostname olt-api \
  --password YourSecurePassword \
  --memory 2048 \
  --swap 512 \
  --cores 2 \
  --storage local-lvm \
  --rootfs local-lvm:10 \
  --net0 name=eth0,bridge=vmbr0,ip=192.168.1.100/24,gw=192.168.1.1 \
  --features nesting=1 \
  --unprivileged 1 \
  --onboot 1
```

### Step 3: Start & Access Container

```bash
# Start container
pct start 100

# Enter container
pct enter 100
```

### Step 4: Install Application

**One-command installation:**

```bash
# Update system
apt update && apt upgrade -y

# Download and run installer
curl -fsSL https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install-quickstart.sh | bash
```

**Or manual installation:**

```bash
# Download full installer
wget https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install.sh
chmod +x install.sh

# Run installer
./install.sh
```

### Step 5: Configure Application

Edit konfigurasi:

```bash
nano /opt/go-snmp-olt/.env
```

Atur sesuai environment Proxmox:

```bash
# OLT Configuration
OLT_IP_ADDRESS=192.168.1.10          # IP OLT (harus reachable dari CT)
OLT_SNMP_PORT=161
OLT_SNMP_COMMUNITY=public
OLT_TELNET_HOST=192.168.1.10
OLT_TELNET_PORT=23
OLT_TELNET_USERNAME=zte
OLT_TELNET_PASSWORD=zte
OLT_TELNET_ENABLE_PASSWORD=zxr10

# API Server
SERVER_PORT=8081

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=OsWkRgJLabn4n2+nodZ6BQeP+OKkrObnGeFcDY6w7Nw=
```

Restart service:

```bash
systemctl restart go-snmp-olt
systemctl status go-snmp-olt
```

---

## 💻 Option 2: Virtual Machine (Full VM)

Untuk kasus yang memerlukan isolasi penuh atau custom kernel.

### Step 1: Create VM

**Via Web UI:**
1. Datacenter → Node → Create VM
2. General:
   - VM ID: `200`
   - Name: `olt-api`

3. OS:
   - ISO image: `ubuntu-22.04-live-server-amd64.iso`
   - Type: `Linux`
   - Version: `6.x - 2.6 Kernel`

4. System:
   - SCSI Controller: `VirtIO SCSI`
   - Qemu Agent: ✅

5. Disks:
   - Bus/Device: `VirtIO Block`
   - Storage: `local-lvm`
   - Disk size: `10 GB`
   - Cache: `Write back`
   - Discard: ✅

6. CPU:
   - Sockets: `1`
   - Cores: `2`
   - Type: `host`

7. Memory:
   - Memory: `2048 MB`
   - Minimum: `512 MB`
   - Ballooning Device: ✅

8. Network:
   - Bridge: `vmbr0`
   - Model: `VirtIO (paravirtualized)`

**Via CLI:**

```bash
qm create 200 \
  --name olt-api \
  --memory 2048 \
  --cores 2 \
  --net0 virtio,bridge=vmbr0 \
  --scsi0 local-lvm:10 \
  --ide2 local:iso/ubuntu-22.04-live-server-amd64.iso,media=cdrom \
  --boot order=scsi0;ide2 \
  --ostype l26 \
  --agent 1 \
  --onboot 1
```

### Step 2: Install Ubuntu Server

1. Start VM dan buka console
2. Install Ubuntu Server 22.04:
   - Language: English
   - Network: Configure static IP
   - Storage: Use entire disk
   - Profile: 
     - Name: `admin`
     - Server: `olt-api`
     - Username: `admin`
     - Password: (secure password)
   - SSH: Install OpenSSH server ✅
   - Featured snaps: None

3. Reboot setelah instalasi selesai

### Step 3: Post-Installation Setup

SSH ke VM:

```bash
ssh admin@192.168.1.101
```

Update system:

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install qemu-guest-agent -y
sudo systemctl enable qemu-guest-agent
sudo systemctl start qemu-guest-agent
```

### Step 4: Install Application

```bash
# One-command installation
curl -fsSL https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install-quickstart.sh | sudo bash
```

Configure seperti di Option 1 Step 5.

---

## 🌐 Network Configuration

### Bridged Network (Recommended)

VM/CT mendapat IP di network yang sama dengan OLT:

```
Proxmox Host: 192.168.1.1/24
OLT Device:   192.168.1.10
CT/VM:        192.168.1.100  ← Bisa akses langsung ke OLT
```

### NAT Network

Jika menggunakan NAT, perlu port forwarding dari Proxmox host:

```bash
# Di Proxmox host, forward port 8081
iptables -t nat -A PREROUTING -p tcp --dport 8081 -j DNAT --to-destination 192.168.100.100:8081
iptables -t nat -A POSTROUTING -j MASQUERADE
```

### Firewall Rules (UFW)

Di dalam CT/VM:

```bash
# Install UFW
sudo apt install ufw -y

# Allow SSH
sudo ufw allow 22/tcp

# Allow API port
sudo ufw allow 8081/tcp

# Enable firewall
sudo ufw enable
```

---

## 🔧 Proxmox-Specific Optimizations

### Container Nesting (Untuk Docker di CT)

Jika ingin run Docker di dalam CT:

```bash
# Edit CT config di Proxmox host
nano /etc/pve/lxc/100.conf

# Add these lines
lxc.apparmor.profile: unconfined
lxc.cap.drop:
features: nesting=1
```

### Resource Limits

Edit `/etc/pve/lxc/100.conf`:

```bash
# CPU limit (cores)
cores: 2

# Memory limit
memory: 2048
swap: 512

# I/O priority
onboot: 1
startup: order=2
```

### Auto-start Configuration

```bash
# Via CLI
pct set 100 --onboot 1 --startup order=2,up=30

# order=2: Start after network (order=1)
# up=30: Wait 30 seconds before starting next service
```

---

## 📊 Monitoring & Maintenance

### Check Container Status

```bash
# Via Proxmox host
pct status 100
pct list

# Resource usage
pct df 100
```

### Inside CT/VM Monitoring

```bash
# Application status
systemctl status go-snmp-olt

# Logs
journalctl -u go-snmp-olt -f

# Resource usage
htop
```

### Backup & Restore

**Backup CT/VM:**

```bash
# Via Proxmox host
vzdump 100 --storage local --mode snapshot --compress zstd

# Backup location: /var/lib/vz/dump/
```

**Restore:**

```bash
pct restore 100 /var/lib/vz/dump/vzdump-lxc-100-*.tar.zst
```

---

## 🚨 Troubleshooting

### CT Won't Start

```bash
# Check config
pct config 100

# Check logs
journalctl -xe | grep pve

# Force stop and start
pct stop 100
pct start 100
```

### Network Issues

```bash
# Inside CT/VM
ip addr show
ping 192.168.1.10  # Ping OLT

# Check routing
ip route

# Test SNMP
snmpwalk -v2c -c public 192.168.1.10 system

# Test Telnet
telnet 192.168.1.10 23
```

### Application Not Starting

```bash
# Check service
systemctl status go-snmp-olt

# Check ports
netstat -tulpn | grep 8081

# Test Redis
redis-cli ping

# Check logs
tail -f /opt/go-snmp-olt/logs/api.log
```

---

## 🎯 Best Practices

### Security

- ✅ Use unprivileged containers when possible
- ✅ Enable firewall (ufw)
- ✅ Change default passwords
- ✅ Regular backups (automated)
- ✅ Keep system updated

### Performance

- ✅ Use VirtIO drivers for VM
- ✅ Enable CPU host passthrough
- ✅ Use SSD storage if available
- ✅ Allocate enough memory (min 2GB)

### Maintenance

- ✅ Enable auto-start
- ✅ Set up monitoring (Prometheus/Grafana)
- ✅ Regular system updates
- ✅ Log rotation configured

---

## 📝 Example: Complete CT Deployment

```bash
# 1. Create CT
pct create 100 local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst \
  --hostname olt-api \
  --password SecurePass123 \
  --memory 2048 \
  --swap 512 \
  --cores 2 \
  --storage local-lvm \
  --rootfs local-lvm:10 \
  --net0 name=eth0,bridge=vmbr0,ip=192.168.1.100/24,gw=192.168.1.1 \
  --features nesting=1 \
  --unprivileged 1 \
  --onboot 1

# 2. Start CT
pct start 100

# 3. Enter CT
pct enter 100

# 4. Update system
apt update && apt upgrade -y

# 5. Install application
curl -fsSL https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install-quickstart.sh | bash

# 6. Configure (edit .env)
nano /opt/go-snmp-olt/.env

# 7. Restart service
systemctl restart go-snmp-olt

# 8. Test API
curl http://localhost:8081/
```

---

## ✅ Installer Status Confirmation

**Auto Installer:** ✅ **FULLY TESTED & WORKING**

Supported environments:
- ✅ **Ubuntu** 18.04, 20.04, 22.04, 24.04
- ✅ **Debian** 10, 11, 12
- ✅ **CentOS** 7, 8
- ✅ **Rocky Linux** 8, 9
- ✅ **Proxmox CT** (LXC containers)
- ✅ **Proxmox VM** (Virtual machines)
- ✅ **Public VPS** (DigitalOcean, AWS, Vultr, dll)
- ✅ **Local VPS** (VMware, VirtualBox, Hyper-V)

**Features:**
- Auto-detects OS and version
- Installs Go 1.25.5
- Installs Redis 7.2
- Creates systemd service
- Configuration wizard
- Zero manual configuration needed

**Installation method:**
```bash
# One-line installation (any supported OS)
curl -fsSL https://raw.githubusercontent.com/s4lfanet/go-api-c320/main/scripts/install-quickstart.sh | sudo bash
```

---

## 📞 Support

Jika ada masalah saat deployment di Proxmox:

1. Check [Troubleshooting Guide](../deployment/TROUBLESHOOTING.md)
2. Review Proxmox logs: `journalctl -xe`
3. Check CT/VM logs: `pct enter 100 && journalctl -u go-snmp-olt`
4. Open issue: [GitHub Issues](https://github.com/s4lfanet/go-api-c320/issues)

---

**Last Updated**: January 12, 2026  
**Tested On**: Proxmox VE 8.1.4  
**Status**: Production Ready ✅
