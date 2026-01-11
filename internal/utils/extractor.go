package utils

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// ExtractONUID function is used to extract ONU ID from OID string
// It assumes the ONU ID is the last numeric component of the dot-separated OID string.
func ExtractONUID(oid string) string {
	// Split the OID name and take the last component
	parts := strings.Split(oid, ".")
	if len(parts) > 0 {
		// Check if the last component is a valid number
		lastComponent := parts[len(parts)-1]
		if _, err := strconv.Atoi(lastComponent); err == nil {
			return lastComponent // Return if valid integer string
		}
	}
	return "" // Return an empty string if the OID is invalid or empty (default value)
}

// ExtractIDOnuID function is used to extract ONU ID from OID interface{}
// Validates the type is string and performs extraction.
func ExtractIDOnuID(oid interface{}) int {
	if oid == nil {
		return 0 // Return 0 if input is nil
	}

	switch v := oid.(type) {
	case string:
		parts := strings.Split(v, ".") // Split string by dot
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			id, err := strconv.Atoi(lastPart) // Convert the last part to int
			if err == nil {
				return id // Return ID
			}
		}
		return 0
	default:
		return 0 // Return 0 for other types
	}
}

// ExtractName function is used to extract name from OID value
// Handling both string and byte slice types.
func ExtractName(oidValue interface{}) string {
	switch v := oidValue.(type) {
	case string:
		// Data is string, return it
		return v
	case []byte:
		// Data is a byte slice, convert to string
		return string(v)
	default:
		// Data type is not recognized
		return "Unknown" // Return "Unknown" default
	}
}

// ExtractSerialNumber function is used to extract serial number from OID value
// Handles GPON ONU serial number format: 4 bytes Vendor ID (ASCII) + 4 bytes Serial (Hex)
// Example: Hex-STRING "5A 54 45 47 D8 24 CD F3" -> "ZTEGD824CDF3"
func ExtractSerialNumber(oidValue interface{}) string {
	var data []byte

	switch v := oidValue.(type) {
	case string:
		// Remove "1," prefix if exists
		v = strings.TrimPrefix(v, "1,")
		// Check if it's already a readable serial number (ASCII printable)
		if len(v) >= 8 && isASCIIPrintable(v) {
			return v
		}
		data = []byte(v)
	case []byte:
		data = v
	default:
		return ""
	}

	// GPON Serial Number format: 4 bytes Vendor ID + 4 bytes Serial Number
	// Total 8 bytes
	if len(data) < 8 {
		// If less than 8 bytes, try to return as string if printable
		if isASCIIPrintable(string(data)) {
			return string(data)
		}
		return hex.EncodeToString(data)
	}

	// First 4 bytes are Vendor ID (ASCII characters like "ZTEG", "HWTC")
	vendorID := string(data[0:4])

	// Last 4 bytes are the serial number (to be converted to uppercase hex)
	serialHex := strings.ToUpper(hex.EncodeToString(data[4:8]))

	return vendorID + serialHex
}

// isASCIIPrintable checks if a string contains only printable ASCII characters
func isASCIIPrintable(s string) bool {
	for _, c := range s {
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}

// ConvertAndMultiply function is used to convert the PDU value to string after multiplying by 0.002 and subtracting 30
// Typically used for converting Optical Power values.
func ConvertAndMultiply(pduValue interface{}) (string, error) {
	// Type asserts pduValue to an integer type
	intValue, ok := pduValue.(int)
	if !ok {
		return "", fmt.Errorf("value is not an integer") // Error if not int
	}

	// Multiply the integer by 0.002 (scale factor)
	result := float64(intValue) * 0.002

	// Subtract 30 (offset)
	result -= 30.0

	// Convert the result to a string with two decimal places
	resultStr := strconv.FormatFloat(result, 'f', 2, 64)

	return resultStr, nil
}

// ExtractAndGetStatus function is used to extract and get status from OID value
// Maps integer status codes to human-readable strings.
func ExtractAndGetStatus(oidValue interface{}) string {
	// Check if oidValue is an integer
	intValue, ok := oidValue.(int)
	if !ok {
		return "Unknown"
	}

	switch intValue {
	case 1:
		return "Logging"
	case 2:
		return "LOS" // Loss of Signal
	case 3:
		return "Synchronization"
	case 4:
		return "Online"
	case 5:
		return "Dying Gasp" // Power failure indication
	case 6:
		return "Auth Failed"
	case 7:
		return "Offline"
	default:
		return "Unknown"
	}
}

// ExtractLastOfflineReason function is used to extract the last offline reason from OID value
// Maps integer reason codes to human-readable strings.
func ExtractLastOfflineReason(oidValue interface{}) string {
	// Check if oidValue is an integer
	intValue, ok := oidValue.(int)
	if !ok {
		return "Unknown"
	}

	switch intValue {
	case 1:
		return "Unknown"
	case 2:
		return "LOS" // Loss of Signal
	case 3:
		return "LOSi"
	case 4:
		return "LOFi" // Loss of Frame
	case 5:
		return "sfi"
	case 6:
		return "loai"
	case 7:
		return "loami"
	case 8:
		return "AuthFail"
	case 9:
		return "PowerOff"
	case 10:
		return "deactiveSucc"
	case 11:
		return "deactiveFail"
	case 12:
		return "Reboot"
	case 13:
		return "Shutdown"
	default:
		return "Unknown"
	}
}

// ExtractGponOpticalDistance function is used to extract GPON optical distance from OID value
func ExtractGponOpticalDistance(oidValue interface{}) string {
	// Check if oidValue is an integer
	intValue, ok := oidValue.(int)
	if !ok {
		return "Unknown"
	}

	return strconv.Itoa(intValue) // Convert integer to string
}
