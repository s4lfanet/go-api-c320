package config

import (
	"fmt"
	"os"
	"strconv"
)

// FirmwareVersion represents the ZTE C320 firmware version
type FirmwareVersion string

const (
	FirmwareV21 FirmwareVersion = "v2.1" // Firmware V2.1.0
	FirmwareV22 FirmwareVersion = "v2.2" // Firmware V2.2.x and newer
)

// OIDProfile contains all OID configurations for a specific firmware version
type OIDProfile struct {
	Name                         string
	BaseOID                      string
	OnuIDNamePrefix              string
	OnuTypePrefix                string
	OnuSerialNumberPrefix        string
	OnuRxPowerPrefix             string
	OnuTxPowerPrefix             string
	OnuStatusIDPrefix            string
	OnuIPAddressPrefix           string
	OnuDescriptionPrefix         string
	OnuLastOnlineTimePrefix      string
	OnuLastOfflineTimePrefix     string
	OnuLastOfflineReasonPrefix   string
	OnuGponOpticalDistancePrefix string
	Board1OnuIDBase              int
	Board1OnuTypeBase            int
	Board2OnuIDBase              int
	Board2OnuTypeBase            int
	OnuIDIncrement               int
	OnuTypeIncrement             int
}

// OID Profiles for different firmware versions
// V2.1.0 uses OID base .1.3.6.1.4.1.3902.1012 (NOT 1082!)
// V2.2+ uses .1.3.6.1.4.1.3902.1082 base
var OIDProfiles = map[FirmwareVersion]*OIDProfile{
	FirmwareV21: {
		Name:    "ZTE C320 V2.1.0",
		BaseOID: ".1.3.6.1.4.1.3902.1012", // V2.1 uses 1012, NOT 1082!
		// GPON ONU Management OIDs for V2.1.0
		// Based on actual SNMP walk results from ZTE C320 V2.1.0
		// ONU Table: .3.13.3.1.{column}.{pon_index}.{onu_id}
		// ONU Statistics: .3.31.4.1.{column}.{pon_index}.{onu_id}
		OnuIDNamePrefix:              ".3.13.3.1.5",  // ONU Device SN (STRING, e.g., "GD824CDF3")
		OnuTypePrefix:                ".3.13.3.1.10", // ONU Model (STRING, e.g., "F672YV9.1")
		OnuSerialNumberPrefix:        ".3.13.3.1.2",  // ONU Serial Number (Hex-STRING)
		OnuRxPowerPrefix:             ".3.31.4.1.100", // ONU Status - no RxPower available in V2.1
		OnuTxPowerPrefix:             ".3.31.4.1.100", // ONU Status - no TxPower available in V2.1
		OnuStatusIDPrefix:            ".3.31.4.1.100", // ONU Online Status (INTEGER: 1=online)
		OnuIPAddressPrefix:           ".3.13.3.1.3",   // ONU Password (no IP in V2.1)
		OnuDescriptionPrefix:         ".3.13.3.1.11",  // ONU Firmware Version
		OnuLastOnlineTimePrefix:      ".3.31.4.1.2",   // Timestamp
		OnuLastOfflineTimePrefix:     ".3.31.4.1.2",   // Timestamp
		OnuLastOfflineReasonPrefix:   ".3.13.3.1.4",   // Status
		OnuGponOpticalDistancePrefix: ".3.13.1.1.20",  // PON settings
		// PON Index for V2.1: pon_index = 268500992 + (pon * 256)
		// Board1 PON1: 268501248, PON2: 268501504, etc
		// Formula: base + (ponID * increment) = pon_index
		// So for PON1: 268500992 + (1 * 256) = 268501248 ✓
		Board1OnuIDBase:              268500992, // Base for Board 1 (268501248 - 256)
		Board1OnuTypeBase:            268500992, // Same as OnuID for V2.1
		Board2OnuIDBase:              268509184, // Board 2 = 268500992 + 8192
		Board2OnuTypeBase:            268509184, // Same as OnuID for V2.1
		OnuIDIncrement:               256, // V2.1 increments by 256 per PON
		OnuTypeIncrement:             256, // Same increment
	},
	FirmwareV22: {
		Name:    "ZTE C320 V2.2+",
		BaseOID: ".1.3.6.1.4.1.3902.1082",
		// Original OIDs for V2.2+ firmware
		OnuIDNamePrefix:              ".500.10.2.3.3.1.2",
		OnuTypePrefix:                ".3.50.11.2.1.17",
		OnuSerialNumberPrefix:        ".500.10.2.3.3.1.18",
		OnuRxPowerPrefix:             ".500.20.2.2.2.1.10",
		OnuTxPowerPrefix:             ".3.50.12.1.1.14",
		OnuStatusIDPrefix:            ".500.10.2.3.8.1.4",
		OnuIPAddressPrefix:           ".3.50.16.1.1.10",
		OnuDescriptionPrefix:         ".500.10.2.3.3.1.3",
		OnuLastOnlineTimePrefix:      ".500.10.2.3.8.1.5",
		OnuLastOfflineTimePrefix:     ".500.10.2.3.8.1.6",
		OnuLastOfflineReasonPrefix:   ".500.10.2.3.8.1.7",
		OnuGponOpticalDistancePrefix: ".500.10.2.3.10.1.2",
		Board1OnuIDBase:              285278464,
		Board1OnuTypeBase:            268500992,
		Board2OnuIDBase:              285278720,
		Board2OnuTypeBase:            268566528,
		OnuIDIncrement:               1,
		OnuTypeIncrement:             256,
	},
}

// GetCurrentFirmwareVersion returns the firmware version from environment variable
// Default is V2.1 if not specified
func GetCurrentFirmwareVersion() FirmwareVersion {
	version := os.Getenv("ZTE_FIRMWARE_VERSION")
	switch version {
	case "v2.2", "V2.2", "2.2":
		return FirmwareV22
	case "v2.1", "V2.1", "2.1", "":
		return FirmwareV21
	default:
		return FirmwareV21
	}
}

// GetOIDProfile returns the OID profile for the current firmware version
func GetOIDProfile() *OIDProfile {
	return OIDProfiles[GetCurrentFirmwareVersion()]
}

// GetOIDProfileForVersion returns the OID profile for a specific firmware version
func GetOIDProfileForVersion(version FirmwareVersion) *OIDProfile {
	if profile, ok := OIDProfiles[version]; ok {
		return profile
	}
	return OIDProfiles[FirmwareV21]
}

// Helper to get environment variable with custom OID override
func getOIDEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getOIDEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// Dynamic OID variables that can be overridden via environment variables
var (
	// Base OID - can be overridden with OLT_BASE_OID environment variable
	// V2.1 uses 1012, V2.2+ uses 1082
	BaseOID1 = getOIDEnv("OLT_BASE_OID", GetOIDProfile().BaseOID)
	BaseOID2 = ".1.3.6.1.4.1.3902.1012" // V2.1 base

	// OID Prefixes - can be overridden individually
	OnuIDNamePrefix              = getOIDEnv("ONU_ID_NAME_PREFIX", GetOIDProfile().OnuIDNamePrefix)
	OnuTypePrefix                = getOIDEnv("ONU_TYPE_PREFIX", GetOIDProfile().OnuTypePrefix)
	OnuSerialNumberPrefix        = getOIDEnv("ONU_SERIAL_NUMBER_PREFIX", GetOIDProfile().OnuSerialNumberPrefix)
	OnuRxPowerPrefix             = getOIDEnv("ONU_RX_POWER_PREFIX", GetOIDProfile().OnuRxPowerPrefix)
	OnuTxPowerPrefix             = getOIDEnv("ONU_TX_POWER_PREFIX", GetOIDProfile().OnuTxPowerPrefix)
	OnuStatusIDPrefix            = getOIDEnv("ONU_STATUS_ID_PREFIX", GetOIDProfile().OnuStatusIDPrefix)
	OnuIPAddressPrefix           = getOIDEnv("ONU_IP_ADDRESS_PREFIX", GetOIDProfile().OnuIPAddressPrefix)
	OnuDescriptionPrefix         = getOIDEnv("ONU_DESCRIPTION_PREFIX", GetOIDProfile().OnuDescriptionPrefix)
	OnuLastOnlineTimePrefix      = getOIDEnv("ONU_LAST_ONLINE_PREFIX", GetOIDProfile().OnuLastOnlineTimePrefix)
	OnuLastOfflineTimePrefix     = getOIDEnv("ONU_LAST_OFFLINE_PREFIX", GetOIDProfile().OnuLastOfflineTimePrefix)
	OnuLastOfflineReasonPrefix   = getOIDEnv("ONU_LAST_OFFLINE_REASON_PREFIX", GetOIDProfile().OnuLastOfflineReasonPrefix)
	OnuGponOpticalDistancePrefix = getOIDEnv("ONU_GPON_OPTICAL_DISTANCE_PREFIX", GetOIDProfile().OnuGponOpticalDistancePrefix)

	// Board-PON ID Constants
	Board1OnuIDBase   = getOIDEnvAsInt("BOARD1_ONU_ID_BASE", GetOIDProfile().Board1OnuIDBase)
	Board1OnuTypeBase = getOIDEnvAsInt("BOARD1_ONU_TYPE_BASE", GetOIDProfile().Board1OnuTypeBase)
	Board2OnuIDBase   = getOIDEnvAsInt("BOARD2_ONU_ID_BASE", GetOIDProfile().Board2OnuIDBase)
	Board2OnuTypeBase = getOIDEnvAsInt("BOARD2_ONU_TYPE_BASE", GetOIDProfile().Board2OnuTypeBase)

	// Increment values
	OnuIDIncrement   = getOIDEnvAsInt("ONU_ID_INCREMENT", GetOIDProfile().OnuIDIncrement)
	OnuTypeIncrement = getOIDEnvAsInt("ONU_TYPE_INCREMENT", GetOIDProfile().OnuTypeIncrement)
)

// GenerateBoardPonOID generates all OID configurations for a specific Board-PON combination
// using mathematical formulas instead of hardcoded config file entries.
//
// Formula:
//   - onuIDSuffix = baseOnuID + ponID
//   - onuTypeSuffix = baseOnuType + (ponID * 256)
//
// Example:
//
//	Board 1, PON 1: 285278465, 268501248
//	Board 1, PON 2: 285278466, 268501504 (+1, +256)
//	Board 2, PON 1: 285278721, 268566784
func GenerateBoardPonOID(boardID, ponID int) (*BoardPonConfig, error) {
	// Validate inputs
	if boardID < 1 || boardID > 2 {
		return nil, fmt.Errorf("invalid boardID: %d (must be 1 or 2)", boardID)
	}
	if ponID < 1 || ponID > 16 {
		return nil, fmt.Errorf("invalid ponID: %d (must be 1-16)", ponID)
	}

	// Determine base values based on board
	var baseOnuID, baseOnuType int
	switch boardID {
	case 1:
		baseOnuID = Board1OnuIDBase
		baseOnuType = Board1OnuTypeBase
	case 2:
		baseOnuID = Board2OnuIDBase
		baseOnuType = Board2OnuTypeBase
	default:
		return nil, fmt.Errorf("unsupported boardID: %d", boardID)
	}

	// Calculate suffixes using formula
	onuIDSuffix := baseOnuID + (ponID * OnuIDIncrement)
	onuTypeSuffix := baseOnuType + (ponID * OnuTypeIncrement)

	// Generate full OIDs by concatenating prefix + suffix
	return &BoardPonConfig{
		OnuIDNameOID:              fmt.Sprintf("%s.%d", OnuIDNamePrefix, onuIDSuffix),
		OnuTypeOID:                fmt.Sprintf("%s.%d", OnuTypePrefix, onuTypeSuffix),
		OnuSerialNumberOID:        fmt.Sprintf("%s.%d", OnuSerialNumberPrefix, onuIDSuffix),
		OnuRxPowerOID:             fmt.Sprintf("%s.%d", OnuRxPowerPrefix, onuIDSuffix),
		OnuTxPowerOID:             fmt.Sprintf("%s.%d", OnuTxPowerPrefix, onuTypeSuffix),
		OnuStatusOID:              fmt.Sprintf("%s.%d", OnuStatusIDPrefix, onuIDSuffix),
		OnuIPAddressOID:           fmt.Sprintf("%s.%d", OnuIPAddressPrefix, onuTypeSuffix),
		OnuDescriptionOID:         fmt.Sprintf("%s.%d", OnuDescriptionPrefix, onuIDSuffix),
		OnuLastOnlineOID:          fmt.Sprintf("%s.%d", OnuLastOnlineTimePrefix, onuIDSuffix),
		OnuLastOfflineOID:         fmt.Sprintf("%s.%d", OnuLastOfflineTimePrefix, onuIDSuffix),
		OnuLastOfflineReasonOID:   fmt.Sprintf("%s.%d", OnuLastOfflineReasonPrefix, onuIDSuffix),
		OnuGponOpticalDistanceOID: fmt.Sprintf("%s.%d", OnuGponOpticalDistancePrefix, onuIDSuffix),
	}, nil
}

// InitializeBoardPonMap generates all 32 Board-PON configurations dynamically.
// This replaces the need for a 20KB config file with 384 lines of OID mappings.
func InitializeBoardPonMap() (map[BoardPonKey]*BoardPonConfig, error) {
	boardPonMap := make(map[BoardPonKey]*BoardPonConfig, 32) // Pre-allocate for 32 entries (2 boards * 16 PONs)

	for boardID := 1; boardID <= 2; boardID++ {
		for ponID := 1; ponID <= 16; ponID++ {
			cfg, err := GenerateBoardPonOID(boardID, ponID)
			if err != nil {
				return nil, fmt.Errorf("failed to generate OID for Board%dPon%d: %w", boardID, ponID, err)
			}
			boardPonMap[BoardPonKey{BoardID: boardID, PonID: ponID}] = cfg
		}
	}

	return boardPonMap, nil
}
