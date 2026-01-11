package model

import "time"

// ONUMonitoringInfo represents real-time monitoring data for a single ONU
type ONUMonitoringInfo struct {
	PonPort      string         `json:"pon_port"`
	OnuID        int            `json:"onu_id"`
	SerialNumber string         `json:"serial_number"`
	Model        string         `json:"model"`
	FirmwareVer  string         `json:"firmware_version"`
	OnlineStatus int            `json:"online_status"` // 1=online, 0=offline
	Statistics   *ONUStatistics `json:"statistics,omitempty"`
	LastUpdate   time.Time      `json:"last_update"`
}

// ONUStatistics represents traffic statistics for an ONU
type ONUStatistics struct {
	RxPackets uint64 `json:"rx_packets"` // Received packets
	RxBytes   uint64 `json:"rx_bytes"`   // Received bytes
	RxRate    string `json:"rx_rate"`    // Human readable rate (e.g., "1.5 Mbps")
}

// PONMonitoringInfo represents aggregated monitoring for a PON port
type PONMonitoringInfo struct {
	PonPort      string              `json:"pon_port"`
	PonIndex     int                 `json:"pon_index"`
	OnuCount     int                 `json:"onu_count"`
	OnlineCount  int                 `json:"online_count"`
	OfflineCount int                 `json:"offline_count"`
	Statistics   *PONStatistics      `json:"statistics,omitempty"`
	ONUs         []ONUMonitoringInfo `json:"onus,omitempty"`
	LastUpdate   time.Time           `json:"last_update"`
}

// PONStatistics represents aggregated traffic statistics for a PON port
type PONStatistics struct {
	RxPackets uint64 `json:"rx_packets"` // Total received packets
	RxBytes   uint64 `json:"rx_bytes"`   // Total received bytes
	RxRate    string `json:"rx_rate"`    // Human readable rate
}

// OLTMonitoringSummary represents overall OLT monitoring summary
type OLTMonitoringSummary struct {
	TotalONUs   int                 `json:"total_onus"`
	OnlineONUs  int                 `json:"online_onus"`
	OfflineONUs int                 `json:"offline_onus"`
	PONPorts    []PONMonitoringInfo `json:"pon_ports"`
	LastUpdate  time.Time           `json:"last_update"`
}
