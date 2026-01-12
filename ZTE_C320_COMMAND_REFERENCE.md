# ZTE C320 V2.1.0 - Command Reference

## Working Commands (Verified on Production OLT 136.1.1.100)

### ONU Management

#### List Unconfigured ONUs
```bash
show gpon onu uncfg
```
**Output:**
```
OnuIndex                 Sn                  State
---------------------------------------------------------------------
gpon-onu_1/1/1:1         HWTC1F14CAAD        unknown
gpon-onu_1/1/1:2         ZTEGD824CDF3        unknown
gpon-onu_1/1/1:3         ZTEGDA5918AC        unknown
```

#### Show Interface
```bash
show interface gpon-olt_1/1/1
```
**Shows:** Port status, registered ONUs count, statistics

#### Show VLAN
```bash
show vlan 100
```
**Shows:** VLAN configuration (name, description, tagged/untagged ports)

### Commands That DON'T WORK (Invalid Syntax)

❌ `show service-port` - Error 20203: Incomplete command  
❌ `show service-port all` - Error 20201: Invalid keyword  
❌ `show dba-profile` - Error 20200: Invalid command  
❌ `show dba` - Error 20200: Invalid command  
❌ `show tcont` - Error 20206: Unrecognized command  
❌ `show gemport` - Error 20206: Unrecognized command  
❌ `conf t` - Error 20200: Invalid command (use `configure terminal`)

### Important Findings

1. **ONUs NOT Registered**
   - `show interface gpon-olt_1/1/1` returns: "the number of registered onus is 0"
   - ONUs are detected (3 ONUs) but not configured
   - Status: "unknown" (not registered to OLT)

2. **No Service Ports**
   - Command syntax unknown or service-ports don't exist
   - Need to find correct command or verify feature availability

3. **VLAN Exists**
   - VLAN 100 is configured
   - But no service-ports binding ONUs to VLANs

## Correct Command Patterns

### Show Commands
```bash
show gpon onu <subcommand>
show gpon remote-onu <subcommand> <interface>
show interface <interface_name>
show vlan <vlan_id>
```

### Configuration Mode
```bash
configure terminal  # NOT 'conf t'
```

## Next Steps Required

1. **Find Service-Port Command**
   - Try: `show gpon remote-onu service-port gpon-onu_1/1/1:1`
   - Try: `show pon onu-info gpon-olt_1/1/1`
   - Try different variations

2. **Register ONUs First**
   - Before configuring VLANs, ONUs must be registered
   - Use: `onu <onu_id> type <type> sn <serial>`

3. **Find DBA/T-CONT Commands**
   - Might be under different name
   - Check profile commands
   - Verify firmware feature availability
