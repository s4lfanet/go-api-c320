package model

// CardInfo represents a card/slot information in the OLT chassis
type CardInfo struct {
	Rack         int    `json:"rack"`          // Rack number
	Shelf        int    `json:"shelf"`         // Shelf number
	Slot         int    `json:"slot"`          // Slot number
	CardType     string `json:"card_type"`     // Card type/model
	Status       string `json:"status"`        // Card status (active/inactive)
	SerialNumber string `json:"serial_number"` // Card serial number
	HardwareVer  string `json:"hardware_ver"`  // Hardware version
	SoftwareVer  string `json:"software_ver"`  // Software version
	Description  string `json:"description"`   // Card description
}
