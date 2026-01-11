package repository

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// OpticalInfo represents ONU optical power information
type OpticalInfo struct {
	OnuID             int     `json:"onu_id"`
	RxPower           float64 `json:"rx_power"`           // OLT receives from ONU (dBm)
	TxPower           float64 `json:"tx_power"`           // ONU transmit power (dBm)
	Temperature       float64 `json:"temperature"`        // ONU temperature (°C)
	Voltage           float64 `json:"voltage"`            // ONU voltage (V)
	BiasCurrent       float64 `json:"bias_current"`       // Bias current (mA)
	OLTRxPower        float64 `json:"olt_rx_power"`       // OLT received power (dBm)
	RxPowerStatus     string  `json:"rx_power_status"`    // normal/low/high
	TxPowerStatus     string  `json:"tx_power_status"`    // normal/low/high
	TemperatureStatus string  `json:"temperature_status"` // normal/low/high
}

// GetONUOpticalInfo retrieves optical power information for a specific ONU via Telnet
// Command: show gpon onu optical-info gpon-olt_1/{board}/{pon} {onu_id}
func (m *TelnetSessionManager) GetONUOpticalInfo(ctx context.Context, boardID, ponID, onuID int) (*OpticalInfo, error) {
	// Format: show gpon onu optical-info gpon-olt_1/1/1 1
	cmd := fmt.Sprintf("show gpon onu optical-info gpon-olt_1/%d/%d %d", boardID, ponID, onuID)

	log.Info().
		Int("board_id", boardID).
		Int("pon_id", ponID).
		Int("onu_id", onuID).
		Str("command", cmd).
		Msg("Getting ONU optical info via Telnet")

	resp, err := m.ExecuteCommand(ctx, cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute optical-info command")
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("command failed: %s", resp.Output)
	}

	// Parse the optical info response
	opticalInfo, err := parseOpticalInfoResponse(resp.Output, onuID)
	if err != nil {
		log.Error().Err(err).Str("response", resp.Output).Msg("Failed to parse optical info")
		return nil, err
	}

	log.Info().
		Int("onu_id", onuID).
		Float64("rx_power", opticalInfo.RxPower).
		Float64("tx_power", opticalInfo.TxPower).
		Float64("temperature", opticalInfo.Temperature).
		Msg("Successfully retrieved ONU optical info")

	return opticalInfo, nil
}

// GetPONOpticalInfo retrieves optical power information for all ONUs on a PON port
func (m *TelnetSessionManager) GetPONOpticalInfo(ctx context.Context, boardID, ponID int) ([]*OpticalInfo, error) {
	// First get list of ONUs on this PON
	// Format: show gpon onu optical-info gpon-olt_1/1/1
	cmd := fmt.Sprintf("show gpon onu optical-info gpon-olt_1/%d/%d", boardID, ponID)

	log.Info().
		Int("board_id", boardID).
		Int("pon_id", ponID).
		Str("command", cmd).
		Msg("Getting all ONU optical info for PON port via Telnet")

	resp, err := m.ExecuteCommand(ctx, cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute optical-info command for PON")
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("command failed: %s", resp.Output)
	}

	// Parse the optical info response for multiple ONUs
	opticalInfoList, err := parsePONOpticalInfoResponse(resp.Output)
	if err != nil {
		log.Error().Err(err).Str("response", resp.Output).Msg("Failed to parse PON optical info")
		return nil, err
	}

	log.Info().
		Int("pon_id", ponID).
		Int("onu_count", len(opticalInfoList)).
		Msg("Successfully retrieved PON optical info")

	return opticalInfoList, nil
}

// parseOpticalInfoResponse parses the output of "show gpon onu optical-info" command for a single ONU
// Expected output format example:
// OLT-Rx Optical-Power(dBm): -18.45
// ONU-Rx Optical-Power(dBm): -19.23
// ONU-Tx Optical-Power(dBm): 2.35
// ONU Laser BIAS-Current(mA): 15.2
// ONU Temperature(C): 42.5
// ONU Voltage(V): 3.28
func parseOpticalInfoResponse(response string, onuID int) (*OpticalInfo, error) {
	info := &OpticalInfo{
		OnuID:             onuID,
		RxPowerStatus:     "unknown",
		TxPowerStatus:     "unknown",
		TemperatureStatus: "unknown",
	}

	lines := strings.Split(response, "\n")

	// Regular expressions for parsing
	oltRxRegex := regexp.MustCompile(`OLT-?Rx\s+Optical-?Power.*?:\s*([-\d.]+)`)
	onuRxRegex := regexp.MustCompile(`ONU-?Rx\s+Optical-?Power.*?:\s*([-\d.]+)`)
	onuTxRegex := regexp.MustCompile(`ONU-?Tx\s+Optical-?Power.*?:\s*([-\d.]+)`)
	biasRegex := regexp.MustCompile(`(?i)BIAS.*?Current.*?:\s*([-\d.]+)`)
	tempRegex := regexp.MustCompile(`(?i)Temperature.*?:\s*([-\d.]+)`)
	voltageRegex := regexp.MustCompile(`(?i)Voltage.*?:\s*([-\d.]+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// OLT RX Power (what OLT receives from ONU)
		if matches := oltRxRegex.FindStringSubmatch(line); len(matches) > 1 {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				info.OLTRxPower = val
				info.RxPower = val // For compatibility
			}
		}

		// ONU RX Power (what ONU receives)
		if matches := onuRxRegex.FindStringSubmatch(line); len(matches) > 1 {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				info.RxPower = val
			}
		}

		// ONU TX Power
		if matches := onuTxRegex.FindStringSubmatch(line); len(matches) > 1 {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				info.TxPower = val
			}
		}

		// Bias Current
		if matches := biasRegex.FindStringSubmatch(line); len(matches) > 1 {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				info.BiasCurrent = val
			}
		}

		// Temperature
		if matches := tempRegex.FindStringSubmatch(line); len(matches) > 1 {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				info.Temperature = val
			}
		}

		// Voltage
		if matches := voltageRegex.FindStringSubmatch(line); len(matches) > 1 {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				info.Voltage = val
			}
		}
	}

	// Determine status based on power levels
	info.RxPowerStatus = classifyPowerLevel(info.RxPower, -28, -8)
	info.TxPowerStatus = classifyPowerLevel(info.TxPower, 0, 5)
	info.TemperatureStatus = classifyTemperature(info.Temperature)

	return info, nil
}

// parsePONOpticalInfoResponse parses optical info for all ONUs on a PON port
// Expected table format:
// OnuId OLT-Rx(dBm) ONU-Rx(dBm) ONU-Tx(dBm) Temp(C) Voltage(V) Current(mA)
//
//	1     -18.45      -19.23       2.35     42.5     3.28       15.2
//	2     -20.12      -21.45       1.98     45.0     3.30       14.8
func parsePONOpticalInfoResponse(response string) ([]*OpticalInfo, error) {
	var results []*OpticalInfo
	lines := strings.Split(response, "\n")

	// Find header line to understand column positions
	dataStarted := false

	// Pattern for table row with ONU data
	// OnuId followed by power values
	rowRegex := regexp.MustCompile(`^\s*(\d+)\s+([-\d.]+)\s+([-\d.]+)\s+([-\d.]+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Look for header indicators to know when data starts
		if strings.Contains(line, "OnuId") || strings.Contains(line, "OLT-Rx") {
			dataStarted = true
			continue
		}

		// Skip separator lines
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "===") {
			continue
		}

		if !dataStarted {
			continue
		}

		// Try to parse as data row
		if matches := rowRegex.FindStringSubmatch(line); len(matches) > 4 {
			info := &OpticalInfo{
				RxPowerStatus:     "unknown",
				TxPowerStatus:     "unknown",
				TemperatureStatus: "unknown",
			}

			if id, err := strconv.Atoi(matches[1]); err == nil {
				info.OnuID = id
			}
			if val, err := strconv.ParseFloat(matches[2], 64); err == nil {
				info.OLTRxPower = val
				info.RxPower = val
			}
			if val, err := strconv.ParseFloat(matches[3], 64); err == nil {
				info.RxPower = val // ONU RX power
			}
			if val, err := strconv.ParseFloat(matches[4], 64); err == nil {
				info.TxPower = val
			}

			// Try to get additional columns if available
			parts := strings.Fields(line)
			if len(parts) >= 5 {
				if val, err := strconv.ParseFloat(parts[4], 64); err == nil {
					info.Temperature = val
				}
			}
			if len(parts) >= 6 {
				if val, err := strconv.ParseFloat(parts[5], 64); err == nil {
					info.Voltage = val
				}
			}
			if len(parts) >= 7 {
				if val, err := strconv.ParseFloat(parts[6], 64); err == nil {
					info.BiasCurrent = val
				}
			}

			// Determine status
			info.RxPowerStatus = classifyPowerLevel(info.RxPower, -28, -8)
			info.TxPowerStatus = classifyPowerLevel(info.TxPower, 0, 5)
			info.TemperatureStatus = classifyTemperature(info.Temperature)

			results = append(results, info)
		}
	}

	return results, nil
}

// classifyPowerLevel categorizes power level as normal/low/high
func classifyPowerLevel(power, minThreshold, maxThreshold float64) string {
	if power == 0 {
		return "unknown"
	}
	if power < minThreshold {
		return "low"
	}
	if power > maxThreshold {
		return "high"
	}
	return "normal"
}

// classifyTemperature categorizes temperature as normal/low/high
func classifyTemperature(temp float64) string {
	if temp == 0 {
		return "unknown"
	}
	if temp < 0 {
		return "low"
	}
	if temp > 70 {
		return "high"
	}
	return "normal"
}
