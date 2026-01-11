package model

import "time"

// TelnetCommand represents a single Telnet command to execute
type TelnetCommand struct {
	Command      string        `json:"command"`
	ExpectPrompt string        `json:"expect_prompt,omitempty"`
	Timeout      time.Duration `json:"timeout,omitempty"`
}

// TelnetResponse represents the response from a Telnet command
type TelnetResponse struct {
	Command   string `json:"command"`
	Output    string `json:"output"`
	Error     string `json:"error,omitempty"`
	Success   bool   `json:"success"`
	Timestamp string `json:"timestamp"`
}

// TelnetBatchResponse represents responses from multiple Telnet commands
type TelnetBatchResponse struct {
	Commands  []TelnetCommand  `json:"commands"`
	Responses []TelnetResponse `json:"responses"`
	Success   bool             `json:"success"`
	TotalTime string           `json:"total_time"`
}

// TelnetConnectionInfo represents information about a Telnet connection
type TelnetConnectionInfo struct {
	Host       string    `json:"host"`
	Port       int       `json:"port"`
	Connected  bool      `json:"connected"`
	Mode       string    `json:"mode"` // "user", "enable", "config"
	LastActive time.Time `json:"last_active"`
	Uptime     string    `json:"uptime"`
}

// TelnetError represents a Telnet-specific error
type TelnetError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Command     string `json:"command,omitempty"`
	RawOutput   string `json:"raw_output,omitempty"`
	Recoverable bool   `json:"recoverable"`
}

func (e *TelnetError) Error() string {
	return e.Message
}

// Common Telnet error codes
const (
	ErrCodeConnectionFailed = "TELNET_CONNECTION_FAILED"
	ErrCodeAuthFailed       = "TELNET_AUTH_FAILED"
	ErrCodeTimeout          = "TELNET_TIMEOUT"
	ErrCodeSessionBusy      = "TELNET_SESSION_BUSY"
	ErrCodeCommandFailed    = "TELNET_COMMAND_FAILED"
	ErrCodeInvalidPrompt    = "TELNET_INVALID_PROMPT"
	ErrCodeDisconnected     = "TELNET_DISCONNECTED"
	ErrCodeConfigSaveFailed = "CONFIG_SAVE_FAILED"
)

// NewTelnetError creates a new Telnet error
func NewTelnetError(code, message string, recoverable bool) *TelnetError {
	return &TelnetError{
		Code:        code,
		Message:     message,
		Recoverable: recoverable,
	}
}

// UnconfiguredONU represents an unconfigured ONU discovered on PON port
type UnconfiguredONU struct {
	PONPort      string `json:"pon_port"`
	SerialNumber string `json:"serial_number"`
	Type         string `json:"type,omitempty"`
	DiscoveredAt string `json:"discovered_at"`
}

// ONURegistrationRequest represents a request to register a new ONU
type ONURegistrationRequest struct {
	PONPort      string `json:"pon_port" validate:"required"`
	ONUID        int    `json:"onu_id" validate:"required,min=1,max=128"`
	ONUType      string `json:"onu_type" validate:"required"`
	SerialNumber string `json:"serial_number" validate:"required"`
	Name         string `json:"name,omitempty"`
	Profile      struct {
		DBAProfile string `json:"dba_profile" validate:"required"`
		VLAN       int    `json:"vlan" validate:"required,min=1,max=4094"`
	} `json:"profile" validate:"required"`
}

// ONURegistrationResponse represents the response after ONU registration
type ONURegistrationResponse struct {
	PONPort       string `json:"pon_port"`
	ONUID         int    `json:"onu_id"`
	SerialNumber  string `json:"serial_number"`
	ServicePortID int    `json:"service_port_id,omitempty"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
}

// TrafficProfileRequest represents a request to create traffic profile
type TrafficProfileRequest struct {
	Name             string `json:"name" validate:"required,max=32"`
	Type             int    `json:"type" validate:"required,min=1,max=5"`
	AssuredBandwidth int    `json:"assured_bandwidth,omitempty"`              // in Kbps
	MaxBandwidth     int    `json:"max_bandwidth" validate:"required,min=64"` // in Kbps
}

// TrafficProfileResponse represents the response after traffic profile creation
type TrafficProfileResponse struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ServicePortRequest represents a request to create service port
type ServicePortRequest struct {
	ServicePortID int    `json:"service_port_id,omitempty"` // Auto-assign if 0
	PONPort       string `json:"pon_port" validate:"required"`
	ONUID         int    `json:"onu_id" validate:"required"`
	GEMPort       int    `json:"gemport" validate:"required"`
	CoS           int    `json:"cos" validate:"min=0,max=7"`
	UserVLAN      string `json:"user_vlan"` // "untagged" or VLAN ID
	VLAN          int    `json:"vlan" validate:"required,min=1,max=4094"`
	RxCTTR        int    `json:"rx_cttr,omitempty"` // RX traffic profile
	TxCTTR        int    `json:"tx_cttr,omitempty"` // TX traffic profile
}

// ServicePortResponse represents the response after service port creation
type ServicePortResponse struct {
	ServicePortID int    `json:"service_port_id"`
	PONPort       string `json:"pon_port"`
	ONUID         int    `json:"onu_id"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
}

// ONURebootRequest represents a request to reboot an ONU
type ONURebootRequest struct {
	PONPort string `json:"pon_port" validate:"required"`
	ONUID   int    `json:"onu_id" validate:"required,min=1,max=128"`
}

// ONURebootResponse represents the response after ONU reboot
type ONURebootResponse struct {
	PONPort string `json:"pon_port"`
	ONUID   int    `json:"onu_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ONUBlockRequest represents a request to block/unblock an ONU
type ONUBlockRequest struct {
	PONPort string `json:"pon_port" validate:"required"`
	ONUID   int    `json:"onu_id" validate:"required,min=1,max=128"`
	Block   bool   `json:"block"` // true = block (disable), false = unblock (enable)
}

// ONUBlockResponse represents the response after ONU block/unblock operation
type ONUBlockResponse struct {
	PONPort string `json:"pon_port"`
	ONUID   int    `json:"onu_id"`
	Blocked bool   `json:"blocked"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ONUDescriptionRequest represents a request to update ONU description/name
type ONUDescriptionRequest struct {
	PONPort     string `json:"pon_port" validate:"required"`
	ONUID       int    `json:"onu_id" validate:"required,min=1,max=128"`
	Description string `json:"description" validate:"required,max=64"`
}

// ONUDescriptionResponse represents the response after updating ONU description
type ONUDescriptionResponse struct {
	PONPort     string `json:"pon_port"`
	ONUID       int    `json:"onu_id"`
	Description string `json:"description"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}

// ONUDeleteRequest represents a request to delete an ONU configuration
type ONUDeleteRequest struct {
	PONPort string `json:"pon_port" validate:"required"`
	ONUID   int    `json:"onu_id" validate:"required,min=1,max=128"`
}

// ONUDeleteResponse represents the response after ONU deletion
type ONUDeleteResponse struct {
	PONPort string `json:"pon_port"`
	ONUID   int    `json:"onu_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ONUVLANInfo represents VLAN configuration for an ONU
type ONUVLANInfo struct {
	PONPort       string `json:"pon_port"`
	ONUID         int    `json:"onu_id"`
	SVLAN         int    `json:"svlan"`     // Service VLAN
	CVLAN         int    `json:"cvlan"`     // Customer VLAN
	VLANMode      string `json:"vlan_mode"` // "tag", "translation", "transparent"
	Priority      int    `json:"priority"`  // CoS/Priority
	ServicePortID int    `json:"service_port_id"`
}

// VLANConfigRequest represents a request to configure ONU VLAN
type VLANConfigRequest struct {
	PONPort       string `json:"pon_port" validate:"required"`
	ONUID         int    `json:"onu_id" validate:"required"`
	SVLAN         int    `json:"svlan" validate:"required,min=1,max=4094"`
	CVLAN         int    `json:"cvlan,omitempty" validate:"omitempty,min=1,max=4094"`
	VLANMode      string `json:"vlan_mode" validate:"required,oneof=tag translation transparent"`
	Priority      int    `json:"priority" validate:"min=0,max=7"`
	ServicePortID int    `json:"service_port_id,omitempty"`
}

// VLANConfigResponse represents the response after VLAN configuration
type VLANConfigResponse struct {
	PONPort       string `json:"pon_port"`
	ONUID         int    `json:"onu_id"`
	SVLAN         int    `json:"svlan"`
	CVLAN         int    `json:"cvlan"`
	VLANMode      string `json:"vlan_mode"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	ServicePortID int    `json:"service_port_id,omitempty"`
}

// DBAProfileInfo represents DBA (Dynamic Bandwidth Allocation) profile information
type DBAProfileInfo struct {
	Name             string `json:"name"`
	Type             int    `json:"type"`                        // 1=Fixed, 2=Assured, 3=Assured+Max, 4=Max, 5=Assured+Max+Priority
	FixedBandwidth   int    `json:"fixed_bandwidth,omitempty"`   // Kbps (Type 1)
	AssuredBandwidth int    `json:"assured_bandwidth,omitempty"` // Kbps (Type 2,3,5)
	MaxBandwidth     int    `json:"max_bandwidth,omitempty"`     // Kbps (Type 3,4,5)
}

// DBAProfileRequest represents a request to create/modify DBA profile
type DBAProfileRequest struct {
	Name             string `json:"name" validate:"required,max=32"`
	Type             int    `json:"type" validate:"required,min=1,max=5"`
	FixedBandwidth   int    `json:"fixed_bandwidth,omitempty" validate:"omitempty,min=64"`
	AssuredBandwidth int    `json:"assured_bandwidth,omitempty" validate:"omitempty,min=64"`
	MaxBandwidth     int    `json:"max_bandwidth,omitempty" validate:"omitempty,min=64"`
}

// DBAProfileResponse represents the response after DBA profile operation
type DBAProfileResponse struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TCONTInfo represents T-CONT (Transmission Container) information
type TCONTInfo struct {
	PONPort   string `json:"pon_port"`
	ONUID     int    `json:"onu_id"`
	TCONTID   int    `json:"tcont_id"`
	Name      string `json:"name"`
	Profile   string `json:"profile"` // DBA profile name
	GEMPorts  []int  `json:"gemports,omitempty"`
	Bandwidth int    `json:"bandwidth,omitempty"` // Current allocated bandwidth in Kbps
}

// TCONTConfigRequest represents a request to configure T-CONT
type TCONTConfigRequest struct {
	PONPort string `json:"pon_port" validate:"required"`
	ONUID   int    `json:"onu_id" validate:"required,min=1,max=128"`
	TCONTID int    `json:"tcont_id" validate:"required,min=1,max=8"`
	Name    string `json:"name,omitempty" validate:"omitempty,max=32"`
	Profile string `json:"profile" validate:"required,max=32"`
}

// TCONTConfigResponse represents the response after T-CONT configuration
type TCONTConfigResponse struct {
	PONPort string `json:"pon_port"`
	ONUID   int    `json:"onu_id"`
	TCONTID int    `json:"tcont_id"`
	Profile string `json:"profile"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GEMPortInfo represents GEM (GPON Encapsulation Method) port information
type GEMPortInfo struct {
	PONPort   string `json:"pon_port"`
	ONUID     int    `json:"onu_id"`
	GEMPortID int    `json:"gemport_id"`
	Name      string `json:"name,omitempty"`
	TCONTID   int    `json:"tcont_id"`
	Queue     int    `json:"queue,omitempty"`
}

// GEMPortConfigRequest represents a request to configure GEM port
type GEMPortConfigRequest struct {
	PONPort   string `json:"pon_port" validate:"required"`
	ONUID     int    `json:"onu_id" validate:"required,min=1,max=128"`
	GEMPortID int    `json:"gemport_id" validate:"required,min=1,max=128"`
	Name      string `json:"name,omitempty" validate:"omitempty,max=32"`
	TCONTID   int    `json:"tcont_id" validate:"required,min=1,max=8"`
	Queue     int    `json:"queue,omitempty" validate:"omitempty,min=1,max=8"`
}

// GEMPortConfigResponse represents the response after GEM port configuration
type GEMPortConfigResponse struct {
	PONPort   string `json:"pon_port"`
	ONUID     int    `json:"onu_id"`
	GEMPortID int    `json:"gemport_id"`
	TCONTID   int    `json:"tcont_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

// TrafficProfileAssignmentRequest represents a request to assign traffic profile to ONU
type TrafficProfileAssignmentRequest struct {
	PONPort    string `json:"pon_port" validate:"required"`
	ONUID      int    `json:"onu_id" validate:"required,min=1,max=128"`
	DBAProfile string `json:"dba_profile" validate:"required,max=32"`
	TCONTID    int    `json:"tcont_id,omitempty" validate:"omitempty,min=1,max=8"`
}

// TrafficProfileAssignmentResponse represents the response after traffic profile assignment
type TrafficProfileAssignmentResponse struct {
	PONPort    string `json:"pon_port"`
	ONUID      int    `json:"onu_id"`
	DBAProfile string `json:"dba_profile"`
	TCONTID    int    `json:"tcont_id"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
}

// ============================================
// Phase 6: Batch Operations Models
// ============================================

// ONUTarget represents a single ONU target for batch operations
type ONUTarget struct {
	PONPort string `json:"pon_port" validate:"required"`
	ONUID   int    `json:"onu_id" validate:"required,min=1,max=128"`
}

// BatchOperationResult represents the result of a single operation in a batch
type BatchOperationResult struct {
	PONPort string `json:"pon_port"`
	ONUID   int    `json:"onu_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// BatchONURebootRequest represents a request to reboot multiple ONUs
type BatchONURebootRequest struct {
	Targets []ONUTarget `json:"targets" validate:"required,min=1,max=50,dive"`
}

// BatchONURebootResponse represents the response after batch ONU reboot
type BatchONURebootResponse struct {
	TotalTargets    int                    `json:"total_targets"`
	SuccessCount    int                    `json:"success_count"`
	FailureCount    int                    `json:"failure_count"`
	Results         []BatchOperationResult `json:"results"`
	ExecutionTimeMs int64                  `json:"execution_time_ms"`
}

// BatchONUBlockRequest represents a request to block/unblock multiple ONUs
type BatchONUBlockRequest struct {
	Targets []ONUTarget `json:"targets" validate:"required,min=1,max=50,dive"`
	Block   bool        `json:"block"` // true=block, false=unblock
}

// BatchONUBlockResponse represents the response after batch ONU block/unblock
type BatchONUBlockResponse struct {
	Blocked         bool                   `json:"blocked"` // Operation type (blocked or unblocked)
	TotalTargets    int                    `json:"total_targets"`
	SuccessCount    int                    `json:"success_count"`
	FailureCount    int                    `json:"failure_count"`
	Results         []BatchOperationResult `json:"results"`
	ExecutionTimeMs int64                  `json:"execution_time_ms"`
}

// BatchONUDeleteRequest represents a request to delete multiple ONUs
type BatchONUDeleteRequest struct {
	Targets []ONUTarget `json:"targets" validate:"required,min=1,max=50,dive"`
}

// BatchONUDeleteResponse represents the response after batch ONU deletion
type BatchONUDeleteResponse struct {
	TotalTargets    int                    `json:"total_targets"`
	SuccessCount    int                    `json:"success_count"`
	FailureCount    int                    `json:"failure_count"`
	Results         []BatchOperationResult `json:"results"`
	ExecutionTimeMs int64                  `json:"execution_time_ms"`
}

// BatchONUDescriptionRequest represents a request to update descriptions for multiple ONUs
type BatchONUDescriptionRequest struct {
	Targets []ONUDescriptionTarget `json:"targets" validate:"required,min=1,max=50,dive"`
}

// ONUDescriptionTarget represents a single ONU with description for batch update
type ONUDescriptionTarget struct {
	PONPort     string `json:"pon_port" validate:"required"`
	ONUID       int    `json:"onu_id" validate:"required,min=1,max=128"`
	Description string `json:"description" validate:"required,max=64"`
}

// BatchONUDescriptionResponse represents the response after batch description update
type BatchONUDescriptionResponse struct {
	TotalTargets    int                    `json:"total_targets"`
	SuccessCount    int                    `json:"success_count"`
	FailureCount    int                    `json:"failure_count"`
	Results         []BatchOperationResult `json:"results"`
	ExecutionTimeMs int64                  `json:"execution_time_ms"`
}
