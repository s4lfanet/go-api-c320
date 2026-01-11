package model

// PonPortConfig represents PON port configuration
type PonPortConfig struct {
	Board       int    `json:"board"`        // Board ID
	PON         int    `json:"pon"`          // PON ID
	AdminStatus string `json:"admin_status"` // Administrative status (up/down)
	// Add more fields as discovered from OID .3.11.3.1.{col}
}

// PonPortStats represents PON port statistics
type PonPortStats struct {
	Board            int    `json:"board"`              // Board ID
	PON              int    `json:"pon"`                // PON ID
	OpticalPower     string `json:"optical_power"`      // Optical power level
	Distance         int    `json:"distance"`           // Distance or range
	OperStatus       string `json:"oper_status"`        // Operational status
	OnuCount         int    `json:"onu_count"`          // Number of ONUs registered
	OnuOnlineCount   int    `json:"onu_online_count"`   // Number of ONUs online
	OnuOfflineCount  int    `json:"onu_offline_count"`  // Number of ONUs offline
	TotalRxBytes     uint64 `json:"total_rx_bytes"`     // Total received bytes
	TotalTxBytes     uint64 `json:"total_tx_bytes"`     // Total transmitted bytes
	TotalRxPackets   uint64 `json:"total_rx_packets"`   // Total received packets
	TotalTxPackets   uint64 `json:"total_tx_packets"`   // Total transmitted packets
}

// PonPortInfo combines config and stats
type PonPortInfo struct {
	Board       int    `json:"board"`        // Board ID
	PON         int    `json:"pon"`          // PON ID
	AdminStatus string `json:"admin_status"` // Administrative status
	OperStatus  string `json:"oper_status"`  // Operational status
	OnuCount    int    `json:"onu_count"`    // Number of ONUs
	// Statistics from .3.11.5.1
	Distance int `json:"distance"` // Distance setting
}
