package repository

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/s4lfanet/go-api-c320/internal/model"
)

// GetONUVLAN retrieves VLAN configuration for a specific ONU
func (m *TelnetSessionManager) GetONUVLAN(ctx context.Context, ponPort string, onuID int) (*model.ONUVLANInfo, error) {
	// Command to show service-port configuration for ONU
	cmd := fmt.Sprintf("show gpon remote-onu interface gpon-olt_%s id %d", ponPort, onuID)

	resp, err := m.ExecuteCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get ONU VLAN info: %w", err)
	}

	// Parse the output to extract VLAN information
	vlanInfo := &model.ONUVLANInfo{
		PONPort: ponPort,
		ONUID:   onuID,
	}

	// Parse service-port information
	// Example output:
	// gpon-olt_1/1/1:1  service-port 123 vlan 100 ...
	servicePortRegex := regexp.MustCompile(`service-port\s+(\d+)`)
	vlanRegex := regexp.MustCompile(`vlan\s+(\d+)`)

	if matches := servicePortRegex.FindStringSubmatch(resp.Output); len(matches) > 1 {
		if portID, err := strconv.Atoi(matches[1]); err == nil {
			vlanInfo.ServicePortID = portID
		}
	}

	if matches := vlanRegex.FindStringSubmatch(resp.Output); len(matches) > 1 {
		if vlan, err := strconv.Atoi(matches[1]); err == nil {
			vlanInfo.SVLAN = vlan
		}
	}

	// Get service-port details for more VLAN info
	if vlanInfo.ServicePortID > 0 {
		detailCmd := fmt.Sprintf("show service-port id %d", vlanInfo.ServicePortID)
		detailResp, err := m.ExecuteCommand(ctx, detailCmd)
		if err == nil {
			vlanInfo = parseServicePortDetails(detailResp.Output, vlanInfo)
		}
	}

	return vlanInfo, nil
}

// ConfigureONUVLAN configures VLAN for an ONU (creates or modifies service-port)
func (m *TelnetSessionManager) ConfigureONUVLAN(ctx context.Context, req model.VLANConfigRequest) (*model.VLANConfigResponse, error) {
	response := &model.VLANConfigResponse{
		PONPort:  req.PONPort,
		ONUID:    req.ONUID,
		SVLAN:    req.SVLAN,
		CVLAN:    req.CVLAN,
		VLANMode: req.VLANMode,
		Success:  false,
	}

	// Build service-port command based on VLAN mode
	var cmd string

	// Check if service-port already exists
	existingVLAN, err := m.GetONUVLAN(ctx, req.PONPort, req.ONUID)
	if err == nil && existingVLAN.ServicePortID > 0 {
		// Update existing service-port
		cmd = m.buildUpdateServicePortCommand(req, existingVLAN.ServicePortID)
		response.ServicePortID = existingVLAN.ServicePortID
	} else {
		// Create new service-port
		cmd = m.buildCreateServicePortCommand(req)
	}

	// Execute in config mode
	commands := []string{cmd}
	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		response.Message = fmt.Sprintf("Failed to configure VLAN: %v", err)
		return response, err
	}

	// Check if successful
	if result.Success {
		response.Success = true
		response.Message = "VLAN configured successfully"

		// Parse service-port ID from output if it's a new creation
		if response.ServicePortID == 0 {
			portIDRegex := regexp.MustCompile(`service-port\s+(\d+)`)
			if matches := portIDRegex.FindStringSubmatch(result.Responses[0].Output); len(matches) > 1 {
				if portID, err := strconv.Atoi(matches[1]); err == nil {
					response.ServicePortID = portID
				}
			}
		}
	} else {
		response.Message = "VLAN configuration failed"
		return response, fmt.Errorf("configuration failed: %s", result.Responses[0].Output)
	}

	return response, nil
}

// DeleteONUVLAN removes VLAN configuration for an ONU (deletes service-port)
func (m *TelnetSessionManager) DeleteONUVLAN(ctx context.Context, ponPort string, onuID int) error {
	// First, get the service-port ID
	vlanInfo, err := m.GetONUVLAN(ctx, ponPort, onuID)
	if err != nil {
		return fmt.Errorf("failed to get VLAN info: %w", err)
	}

	if vlanInfo.ServicePortID == 0 {
		return fmt.Errorf("no service-port found for ONU %s:%d", ponPort, onuID)
	}

	// Delete the service-port
	cmd := fmt.Sprintf("no service-port %d", vlanInfo.ServicePortID)
	commands := []string{cmd}

	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		return fmt.Errorf("failed to delete service-port: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("failed to delete service-port: %s", result.Responses[0].Output)
	}

	return nil
}

// buildCreateServicePortCommand builds command to create new service-port
func (m *TelnetSessionManager) buildCreateServicePortCommand(req model.VLANConfigRequest) string {
	// Base command: service-port <index> vlan <vlan-id> gpon <gpon-olt_id> gemport <gemport-id> queue <queue-id>
	// For ZTE C320, typical format:
	// service-port vlan 100 gpon 1/1/1 gemport 1 multi-service user-vlan 100

	cmd := fmt.Sprintf("service-port vlan %d gpon %s gemport %d multi-service",
		req.SVLAN, req.PONPort, req.ONUID)

	// Add user-vlan (CVLAN) if specified
	if req.CVLAN > 0 {
		cmd += fmt.Sprintf(" user-vlan %d", req.CVLAN)
	} else {
		cmd += fmt.Sprintf(" user-vlan %d", req.SVLAN)
	}

	// Add VLAN mode
	switch req.VLANMode {
	case "tag":
		cmd += " tag-transform default"
	case "translation":
		if req.CVLAN > 0 {
			cmd += fmt.Sprintf(" vlan-translation %d", req.CVLAN)
		}
	case "transparent":
		cmd += " vlan-transparent"
	}

	// Add priority if specified
	if req.Priority > 0 {
		cmd += fmt.Sprintf(" rx-cttr %d tx-cttr %d", req.Priority, req.Priority)
	}

	return cmd
}

// buildUpdateServicePortCommand builds command to update existing service-port
func (m *TelnetSessionManager) buildUpdateServicePortCommand(req model.VLANConfigRequest, servicePortID int) string {
	// For updating, we typically delete and recreate
	// But for ZTE, we can use modify command:
	// service-port <id> vlan <new-vlan>

	cmd := fmt.Sprintf("service-port %d vlan %d", servicePortID, req.SVLAN)

	if req.CVLAN > 0 {
		cmd += fmt.Sprintf(" user-vlan %d", req.CVLAN)
	}

	return cmd
}

// parseServicePortDetails parses detailed service-port output
func parseServicePortDetails(output string, vlanInfo *model.ONUVLANInfo) *model.ONUVLANInfo {
	// Parse output to extract VLAN details
	// Example:
	// Index  VLAN  User-VLAN  Mode          ...
	// 123    100   100        tag-transform ...

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Index") || strings.HasPrefix(line, "---") {
			continue
		}

		// Parse line with regex
		// Format: <index> <vlan> <user-vlan> <mode> ...
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			// SVLAN
			if vlan, err := strconv.Atoi(fields[1]); err == nil {
				vlanInfo.SVLAN = vlan
			}
			// CVLAN (user-vlan)
			if cvlan, err := strconv.Atoi(fields[2]); err == nil {
				vlanInfo.CVLAN = cvlan
			}
			// Mode
			vlanInfo.VLANMode = fields[3]
		}
	}

	return vlanInfo
}

// GetAllServicePorts retrieves all service-port configurations
func (m *TelnetSessionManager) GetAllServicePorts(ctx context.Context) ([]model.ONUVLANInfo, error) {
	cmd := "show service-port"

	resp, err := m.ExecuteCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get service-ports: %w", err)
	}

	return parseAllServicePorts(resp.Output), nil
}

// parseAllServicePorts parses output of "show service-port"
func parseAllServicePorts(output string) []model.ONUVLANInfo {
	var servicePorts []model.ONUVLANInfo

	// Example output:
	// Index  VLAN  Gpon-Port  Gem-Port  User-VLAN  Mode
	// 1      100   1/1/1:1    1         100        tag-transform

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Index") || strings.HasPrefix(line, "---") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 6 {
			sp := model.ONUVLANInfo{}

			// Parse index
			if idx, err := strconv.Atoi(fields[0]); err == nil {
				sp.ServicePortID = idx
			}

			// Parse VLAN
			if vlan, err := strconv.Atoi(fields[1]); err == nil {
				sp.SVLAN = vlan
			}

			// Parse Gpon-Port (format: 1/1/1:1)
			gponPort := fields[2]
			if parts := strings.Split(gponPort, ":"); len(parts) == 2 {
				sp.PONPort = parts[0]
				if onuID, err := strconv.Atoi(parts[1]); err == nil {
					sp.ONUID = onuID
				}
			}

			// Parse User-VLAN
			if len(fields) > 4 {
				if cvlan, err := strconv.Atoi(fields[4]); err == nil {
					sp.CVLAN = cvlan
				}
			}

			// Parse Mode
			if len(fields) > 5 {
				sp.VLANMode = fields[5]
			}

			servicePorts = append(servicePorts, sp)
		}
	}

	return servicePorts
}
