package repository

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/s4lfanet/go-api-c320/internal/model"
)

// GetDBAProfile retrieves DBA profile information
func (m *TelnetSessionManager) GetDBAProfile(ctx context.Context, name string) (*model.DBAProfileInfo, error) {
	cmd := fmt.Sprintf("show gpon-onu-profile dba-profile %s", name)

	resp, err := m.ExecuteCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get DBA profile: %w", err)
	}

	profile := parseDBAProfileOutput(resp.Output, name)
	if profile == nil {
		return nil, fmt.Errorf("DBA profile not found: %s", name)
	}

	return profile, nil
}

// GetAllDBAProfiles retrieves all DBA profiles
func (m *TelnetSessionManager) GetAllDBAProfiles(ctx context.Context) ([]model.DBAProfileInfo, error) {
	// V2.1.0: Try different command variations for DBA profiles
	commands := []string{
		"show gpon profile tcont",           // Try ZTE standard
		"show gpon-onu-profile dba-profile", // Try V2.2+ format
		"show tcont",                         // Try short form
	}
	
	var lastErr error
	for _, cmd := range commands {
		resp, err := m.ExecuteCommand(ctx, cmd)
		if err != nil {
			lastErr = err
			continue
		}
		
		// If command succeeded, try to parse
		profiles := parseAllDBAProfiles(resp.Output)
		if len(profiles) > 0 {
			return profiles, nil
		}
	}
	
	// If all commands failed or returned empty, return empty list instead of error
	// This is normal for unconfigured OLT
	return []model.DBAProfileInfo{}, nil
}

// CreateDBAProfile creates a new DBA profile
func (m *TelnetSessionManager) CreateDBAProfile(ctx context.Context, req model.DBAProfileRequest) (*model.DBAProfileResponse, error) {
	response := &model.DBAProfileResponse{
		Name:    req.Name,
		Type:    req.Type,
		Success: false,
	}

	// Build DBA profile configuration commands
	commands := []string{
		fmt.Sprintf("gpon-onu-profile dba-profile %s", req.Name),
	}

	// Add type-specific bandwidth configuration
	switch req.Type {
	case 1: // Fixed bandwidth
		if req.FixedBandwidth <= 0 {
			response.Message = "Fixed bandwidth is required for Type 1"
			return response, fmt.Errorf("invalid fixed bandwidth")
		}
		commands = append(commands, fmt.Sprintf("type 1 fix %d", req.FixedBandwidth))

	case 2: // Assured bandwidth
		if req.AssuredBandwidth <= 0 {
			response.Message = "Assured bandwidth is required for Type 2"
			return response, fmt.Errorf("invalid assured bandwidth")
		}
		commands = append(commands, fmt.Sprintf("type 2 assure %d", req.AssuredBandwidth))

	case 3: // Assured + Maximum bandwidth
		if req.AssuredBandwidth <= 0 || req.MaxBandwidth <= 0 {
			response.Message = "Both assured and max bandwidth are required for Type 3"
			return response, fmt.Errorf("invalid bandwidth values")
		}
		commands = append(commands, fmt.Sprintf("type 3 assure %d max %d", req.AssuredBandwidth, req.MaxBandwidth))

	case 4: // Maximum bandwidth
		if req.MaxBandwidth <= 0 {
			response.Message = "Maximum bandwidth is required for Type 4"
			return response, fmt.Errorf("invalid max bandwidth")
		}
		commands = append(commands, fmt.Sprintf("type 4 max %d", req.MaxBandwidth))

	case 5: // Assured + Maximum with priority (same as Type 3)
		if req.AssuredBandwidth <= 0 || req.MaxBandwidth <= 0 {
			response.Message = "Both assured and max bandwidth are required for Type 5"
			return response, fmt.Errorf("invalid bandwidth values")
		}
		commands = append(commands, fmt.Sprintf("type 5 assure %d max %d", req.AssuredBandwidth, req.MaxBandwidth))

	default:
		response.Message = "Invalid DBA profile type (must be 1-5)"
		return response, fmt.Errorf("invalid type")
	}

	commands = append(commands, "exit")

	// Execute in config mode
	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		response.Message = fmt.Sprintf("Failed to create DBA profile: %v", err)
		return response, err
	}

	if result.Success {
		response.Success = true
		response.Message = "DBA profile created successfully"
	} else {
		response.Message = "DBA profile creation failed"
	}

	return response, nil
}

// ModifyDBAProfile modifies an existing DBA profile
func (m *TelnetSessionManager) ModifyDBAProfile(ctx context.Context, req model.DBAProfileRequest) (*model.DBAProfileResponse, error) {
	response := &model.DBAProfileResponse{
		Name:    req.Name,
		Type:    req.Type,
		Success: false,
	}

	// Build modification commands (delete old type, add new type)
	commands := []string{
		fmt.Sprintf("gpon-onu-profile dba-profile %s", req.Name),
		"no type", // Remove existing type configuration
	}

	// Add new type configuration
	switch req.Type {
	case 1:
		commands = append(commands, fmt.Sprintf("type 1 fix %d", req.FixedBandwidth))
	case 2:
		commands = append(commands, fmt.Sprintf("type 2 assure %d", req.AssuredBandwidth))
	case 3:
		commands = append(commands, fmt.Sprintf("type 3 assure %d max %d", req.AssuredBandwidth, req.MaxBandwidth))
	case 4:
		commands = append(commands, fmt.Sprintf("type 4 max %d", req.MaxBandwidth))
	case 5:
		commands = append(commands, fmt.Sprintf("type 5 assure %d max %d", req.AssuredBandwidth, req.MaxBandwidth))
	}

	commands = append(commands, "exit")

	// Execute in config mode
	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		response.Message = fmt.Sprintf("Failed to modify DBA profile: %v", err)
		return response, err
	}

	if result.Success {
		response.Success = true
		response.Message = "DBA profile modified successfully"
	} else {
		response.Message = "DBA profile modification failed"
	}

	return response, nil
}

// DeleteDBAProfile deletes a DBA profile
func (m *TelnetSessionManager) DeleteDBAProfile(ctx context.Context, name string) error {
	commands := []string{
		fmt.Sprintf("no gpon-onu-profile dba-profile %s", name),
	}

	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		return fmt.Errorf("failed to delete DBA profile: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("DBA profile deletion failed")
	}

	return nil
}

// GetONUTCONT retrieves T-CONT configuration for an ONU
func (m *TelnetSessionManager) GetONUTCONT(ctx context.Context, ponPort string, onuID int, tcontID int) (*model.TCONTInfo, error) {
	cmd := fmt.Sprintf("show gpon remote-onu interface gpon-onu_%s:%d", ponPort, onuID)

	resp, err := m.ExecuteCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get TCONT info: %w", err)
	}

	tcont := parseTCONTOutput(resp.Output, ponPort, onuID, tcontID)
	if tcont == nil {
		return nil, fmt.Errorf("TCONT %d not found for ONU %s:%d", tcontID, ponPort, onuID)
	}

	return tcont, nil
}

// ConfigureTCONT configures T-CONT for an ONU
func (m *TelnetSessionManager) ConfigureTCONT(ctx context.Context, req model.TCONTConfigRequest) (*model.TCONTConfigResponse, error) {
	response := &model.TCONTConfigResponse{
		PONPort: req.PONPort,
		ONUID:   req.ONUID,
		TCONTID: req.TCONTID,
		Profile: req.Profile,
		Success: false,
	}

	// Enter ONU configuration mode
	commands := []string{
		fmt.Sprintf("interface gpon-onu_%s:%d", req.PONPort, req.ONUID),
	}

	// Build TCONT command
	if req.Name != "" {
		commands = append(commands, fmt.Sprintf("tcont %d name %s profile %s", req.TCONTID, req.Name, req.Profile))
	} else {
		commands = append(commands, fmt.Sprintf("tcont %d profile %s", req.TCONTID, req.Profile))
	}

	commands = append(commands, "exit")

	// Execute in config mode
	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		response.Message = fmt.Sprintf("Failed to configure TCONT: %v", err)
		return response, err
	}

	if result.Success {
		response.Success = true
		response.Message = "TCONT configured successfully"
	} else {
		response.Message = "TCONT configuration failed"
	}

	return response, nil
}

// DeleteTCONT deletes T-CONT from an ONU
func (m *TelnetSessionManager) DeleteTCONT(ctx context.Context, ponPort string, onuID int, tcontID int) error {
	commands := []string{
		fmt.Sprintf("interface gpon-onu_%s:%d", ponPort, onuID),
		fmt.Sprintf("no tcont %d", tcontID),
		"exit",
	}

	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		return fmt.Errorf("failed to delete TCONT: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("TCONT deletion failed")
	}

	return nil
}

// ConfigureGEMPort configures GEM port for an ONU
func (m *TelnetSessionManager) ConfigureGEMPort(ctx context.Context, req model.GEMPortConfigRequest) (*model.GEMPortConfigResponse, error) {
	response := &model.GEMPortConfigResponse{
		PONPort:   req.PONPort,
		ONUID:     req.ONUID,
		GEMPortID: req.GEMPortID,
		TCONTID:   req.TCONTID,
		Success:   false,
	}

	// Enter ONU configuration mode
	commands := []string{
		fmt.Sprintf("interface gpon-onu_%s:%d", req.PONPort, req.ONUID),
	}

	// Build GEMPort command
	if req.Name != "" && req.Queue > 0 {
		commands = append(commands, fmt.Sprintf("gemport %d name %s tcont %d queue %d", req.GEMPortID, req.Name, req.TCONTID, req.Queue))
	} else if req.Name != "" {
		commands = append(commands, fmt.Sprintf("gemport %d name %s tcont %d", req.GEMPortID, req.Name, req.TCONTID))
	} else {
		commands = append(commands, fmt.Sprintf("gemport %d tcont %d", req.GEMPortID, req.TCONTID))
	}

	commands = append(commands, "exit")

	// Execute in config mode
	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		response.Message = fmt.Sprintf("Failed to configure GEM port: %v", err)
		return response, err
	}

	if result.Success {
		response.Success = true
		response.Message = "GEM port configured successfully"
	} else {
		response.Message = "GEM port configuration failed"
	}

	return response, nil
}

// DeleteGEMPort deletes GEM port from an ONU
func (m *TelnetSessionManager) DeleteGEMPort(ctx context.Context, ponPort string, onuID int, gemportID int) error {
	commands := []string{
		fmt.Sprintf("interface gpon-onu_%s:%d", ponPort, onuID),
		fmt.Sprintf("no gemport %d", gemportID),
		"exit",
	}

	result, err := m.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		return fmt.Errorf("failed to delete GEM port: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("GEM port deletion failed")
	}

	return nil
}

// parseDBAProfileOutput parses the output of "show gpon-onu-profile dba-profile" command
func parseDBAProfileOutput(output, name string) *model.DBAProfileInfo {
	profile := &model.DBAProfileInfo{
		Name: name,
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse type line: "type 4 assure 10240 max 10240"
		typeRegex := regexp.MustCompile(`type\s+(\d+)\s*(.*)`)
		if matches := typeRegex.FindStringSubmatch(line); len(matches) > 1 {
			if profileType, err := strconv.Atoi(matches[1]); err == nil {
				profile.Type = profileType

				// Parse bandwidth values based on type
				bandwidthPart := matches[2]

				// Parse "fix" value for Type 1
				fixRegex := regexp.MustCompile(`fix\s+(\d+)`)
				if fixMatches := fixRegex.FindStringSubmatch(bandwidthPart); len(fixMatches) > 1 {
					if bw, err := strconv.Atoi(fixMatches[1]); err == nil {
						profile.FixedBandwidth = bw
					}
				}

				// Parse "assure" value
				assureRegex := regexp.MustCompile(`assure\s+(\d+)`)
				if assureMatches := assureRegex.FindStringSubmatch(bandwidthPart); len(assureMatches) > 1 {
					if bw, err := strconv.Atoi(assureMatches[1]); err == nil {
						profile.AssuredBandwidth = bw
					}
				}

				// Parse "max" value
				maxRegex := regexp.MustCompile(`max\s+(\d+)`)
				if maxMatches := maxRegex.FindStringSubmatch(bandwidthPart); len(maxMatches) > 1 {
					if bw, err := strconv.Atoi(maxMatches[1]); err == nil {
						profile.MaxBandwidth = bw
					}
				}
			}
		}
	}

	// Return nil if type was not found (profile doesn't exist)
	if profile.Type == 0 {
		return nil
	}

	return profile
}

// parseAllDBAProfiles parses the output to get all DBA profiles
func parseAllDBAProfiles(output string) []model.DBAProfileInfo {
	var profiles []model.DBAProfileInfo

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse profile name from header or list output
		// Example: "dba-profile UP-10M"
		nameRegex := regexp.MustCompile(`(?:dba-profile|Name)\s+([A-Za-z0-9_-]+)`)
		if matches := nameRegex.FindStringSubmatch(line); len(matches) > 1 {
			profileName := matches[1]

			// Create basic profile info (detailed info requires individual query)
			profile := model.DBAProfileInfo{
				Name: profileName,
			}

			profiles = append(profiles, profile)
		}
	}

	return profiles
}

// parseTCONTOutput parses the output to extract TCONT information
func parseTCONTOutput(output, ponPort string, onuID, tcontID int) *model.TCONTInfo {
	tcont := &model.TCONTInfo{
		PONPort: ponPort,
		ONUID:   onuID,
		TCONTID: tcontID,
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse TCONT line: "tcont 1 name TCONT_DATA profile UP-10M"
		tcontRegex := regexp.MustCompile(fmt.Sprintf(`tcont\s+%d\s+(?:name\s+(\S+)\s+)?profile\s+(\S+)`, tcontID))
		if matches := tcontRegex.FindStringSubmatch(line); len(matches) > 1 {
			if matches[1] != "" {
				tcont.Name = matches[1]
			}
			if matches[2] != "" {
				tcont.Profile = matches[2]
			}
			return tcont
		}
	}

	// Return nil if TCONT not found
	if tcont.Profile == "" {
		return nil
	}

	return tcont
}
