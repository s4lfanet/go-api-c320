package model

import "time"

// ConfigBackup represents a configuration backup for an ONU or entire OLT
type ConfigBackup struct {
	ID          string         `json:"id"`                    // UUID for backup
	Type        string         `json:"type"`                  // "onu" or "olt"
	Timestamp   time.Time      `json:"timestamp"`             // When backup was created
	Description string         `json:"description,omitempty"` // User-provided description
	Metadata    BackupMetadata `json:"metadata"`              // Additional metadata
	Config      interface{}    `json:"config"`                // ONU config or OLT config
}

// BackupMetadata contains additional information about the backup
type BackupMetadata struct {
	CreatedBy    string            `json:"created_by,omitempty"`    // User/system that created backup
	Source       string            `json:"source"`                  // Source OLT IP or identifier
	Version      string            `json:"version"`                 // Firmware version at backup time
	TotalONUs    int               `json:"total_onus,omitempty"`    // For OLT-level backups
	Tags         []string          `json:"tags,omitempty"`          // Custom tags
	CustomFields map[string]string `json:"custom_fields,omitempty"` // Additional custom metadata
}

// ONUConfigBackup represents configuration for a single ONU
type ONUConfigBackup struct {
	PONPort      string `json:"pon_port"`              // e.g., "1/1/1"
	ONUID        int    `json:"onu_id"`                // ONU ID on the PON
	SerialNumber string `json:"serial_number"`         // ONU serial number
	Type         string `json:"type"`                  // ONU type/model
	Name         string `json:"name,omitempty"`        // ONU description/name
	AuthMethod   string `json:"auth_method,omitempty"` // "sn", "password", "loid"

	// VLAN Configuration
	VLANs []ONUVLANConfig `json:"vlans,omitempty"` // VLAN configurations

	// Traffic/QoS Configuration
	TCONTs   []ONUTCONTConfig   `json:"tconts,omitempty"`   // T-CONT configurations
	GEMPorts []ONUGEMPortConfig `json:"gemports,omitempty"` // GEM port configurations

	// Service Ports
	ServicePorts []ONUServicePortConfig `json:"service_ports,omitempty"` // Service port configurations

	// Additional settings
	AdminState string `json:"admin_state,omitempty"` // "enabled" or "disabled"
	OperState  string `json:"oper_state,omitempty"`  // "online", "offline", etc.

	// Custom settings
	CustomConfig map[string]interface{} `json:"custom_config,omitempty"` // For extensibility
}

// ONUVLANConfig represents VLAN configuration for an ONU
type ONUVLANConfig struct {
	UserVLAN    int    `json:"user_vlan"`          // User-side VLAN
	ServiceVLAN int    `json:"service_vlan"`       // Service-side VLAN
	Mode        string `json:"mode"`               // "tag", "untag", "translation"
	Priority    int    `json:"priority,omitempty"` // 802.1p priority
	TPID        string `json:"tpid,omitempty"`     // Tag Protocol ID
}

// ONUTCONTConfig represents T-CONT configuration
type ONUTCONTConfig struct {
	TCONTID     int    `json:"tcont_id"`               // T-CONT ID (1-8)
	Name        string `json:"name,omitempty"`         // T-CONT name
	ProfileName string `json:"profile_name,omitempty"` // DBA profile name
	Type        int    `json:"type,omitempty"`         // T-CONT type
}

// ONUGEMPortConfig represents GEM port configuration
type ONUGEMPortConfig struct {
	GEMPortID  int    `json:"gemport_id"`           // GEM port ID
	Name       string `json:"name,omitempty"`       // GEM port name
	TCONTID    int    `json:"tcont_id"`             // Associated T-CONT
	Direction  string `json:"direction,omitempty"`  // "upstream", "downstream", "both"
	Encryption bool   `json:"encryption,omitempty"` // AES encryption enabled
}

// ONUServicePortConfig represents service port configuration
type ONUServicePortConfig struct {
	PortID      int `json:"port_id"`              // Service port ID
	VPort       int `json:"vport"`                // Virtual port
	UserVLAN    int `json:"user_vlan"`            // User VLAN
	ServiceVLAN int `json:"service_vlan"`         // Service VLAN
	GEMPortID   int `json:"gemport_id,omitempty"` // Associated GEM port
}

// OLTConfigBackup represents configuration for entire OLT
type OLTConfigBackup struct {
	OLTIP           string `json:"olt_ip"`             // OLT management IP
	Hostname        string `json:"hostname,omitempty"` // OLT hostname
	Model           string `json:"model"`              // OLT model (e.g., "C320")
	FirmwareVersion string `json:"firmware_version"`   // Firmware version

	// ONU configurations
	ONUs []ONUConfigBackup `json:"onus"` // All ONU configurations

	// Global configurations
	GlobalVLANs []GlobalVLANConfig `json:"global_vlans,omitempty"` // Global VLAN definitions
	DBAProfiles []DBAProfileConfig `json:"dba_profiles,omitempty"` // DBA profiles

	// PON port configurations
	PONPorts []PONPortConfig `json:"pon_ports,omitempty"` // PON port settings
}

// GlobalVLANConfig represents global VLAN definition
type GlobalVLANConfig struct {
	VLANID      int    `json:"vlan_id"`               // VLAN ID
	Name        string `json:"name,omitempty"`        // VLAN name
	Description string `json:"description,omitempty"` // VLAN description
}

// DBAProfileConfig represents DBA profile configuration
type DBAProfileConfig struct {
	Name    string `json:"name"`              // Profile name
	Type    string `json:"type"`              // Profile type
	Fixed   int    `json:"fixed,omitempty"`   // Fixed bandwidth (kbps)
	Assured int    `json:"assured,omitempty"` // Assured bandwidth (kbps)
	Maximum int    `json:"maximum,omitempty"` // Maximum bandwidth (kbps)
}

// PONPortConfig represents PON port configuration
type PONPortConfig struct {
	PONPort    string `json:"pon_port"`              // PON port identifier (e.g., "1/1/1")
	AdminState string `json:"admin_state,omitempty"` // "enabled" or "disabled"
	OperState  string `json:"oper_state,omitempty"`  // Operational state
	ActiveONUs int    `json:"active_onus,omitempty"` // Number of active ONUs
}

// BackupListItem represents a backup entry in the list
type BackupListItem struct {
	ID          string    `json:"id"`                    // Backup ID
	Type        string    `json:"type"`                  // "onu" or "olt"
	Timestamp   time.Time `json:"timestamp"`             // When backup was created
	Description string    `json:"description,omitempty"` // User description
	Size        int64     `json:"size"`                  // Backup file size in bytes
	ONUCount    int       `json:"onu_count,omitempty"`   // Number of ONUs in backup
	Source      string    `json:"source,omitempty"`      // Source identifier
	Tags        []string  `json:"tags,omitempty"`        // Tags
}

// RestoreRequest represents a request to restore configuration
type RestoreRequest struct {
	BackupID     string   `json:"backup_id"`               // ID of backup to restore
	TargetType   string   `json:"target_type,omitempty"`   // "same", "different" - where to restore
	TargetPON    string   `json:"target_pon,omitempty"`    // Target PON (if different from backup)
	TargetONUID  int      `json:"target_onu_id,omitempty"` // Target ONU ID (if different)
	Overwrite    bool     `json:"overwrite,omitempty"`     // Overwrite existing config
	DryRun       bool     `json:"dry_run,omitempty"`       // Simulate restore without applying
	RestoreItems []string `json:"restore_items,omitempty"` // Specific items to restore: "vlan", "tcont", "gemport", "service_port"
}

// RestoreResult represents the result of a restore operation
type RestoreResult struct {
	BackupID     string              `json:"backup_id"`         // ID of backup used
	Success      bool                `json:"success"`           // Overall success
	Message      string              `json:"message"`           // Summary message
	RestoredONUs int                 `json:"restored_onus"`     // Number of ONUs restored
	FailedONUs   int                 `json:"failed_onus"`       // Number of failed ONUs
	Details      []RestoreItemResult `json:"details,omitempty"` // Detailed results
	DryRun       bool                `json:"dry_run,omitempty"` // Was this a dry run
}

// RestoreItemResult represents result for individual restore item
type RestoreItemResult struct {
	PONPort  string `json:"pon_port,omitempty"` // PON port
	ONUID    int    `json:"onu_id,omitempty"`   // ONU ID
	ItemType string `json:"item_type"`          // "onu", "vlan", "tcont", etc.
	Success  bool   `json:"success"`            // Item restore success
	Message  string `json:"message,omitempty"`  // Result message
	Error    string `json:"error,omitempty"`    // Error message if failed
}

// BackupCreateRequest represents request to create a backup
type BackupCreateRequest struct {
	Type         string   `json:"type"`                    // "onu" or "olt"
	PONPort      string   `json:"pon_port,omitempty"`      // For ONU backup
	ONUID        int      `json:"onu_id,omitempty"`        // For ONU backup
	Description  string   `json:"description,omitempty"`   // User description
	Tags         []string `json:"tags,omitempty"`          // Custom tags
	IncludeItems []string `json:"include_items,omitempty"` // What to include: "vlan", "tcont", "gemport", "service_port"
}
