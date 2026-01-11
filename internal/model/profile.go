package model

// TrafficProfile represents a traffic profile configuration
type TrafficProfile struct {
	ProfileID int    `json:"profile_id"` // Traffic profile ID
	Name      string `json:"name"`       // Profile name
	CIR       int    `json:"cir"`        // Committed Information Rate (kbps)
	PIR       int    `json:"pir"`        // Peak Information Rate (kbps)
	MaxBW     int    `json:"max_bw"`     // Maximum Bandwidth (kbps)
}

// VlanProfile represents a VLAN profile configuration
type VlanProfile struct {
	Name        string `json:"name"`         // VLAN profile name
	VlanID      int    `json:"vlan_id"`      // VLAN ID
	Priority    int    `json:"priority"`     // VLAN priority
	Mode        string `json:"mode"`         // VLAN mode (tag/untag)
	Description string `json:"description"`  // Profile description
}
