package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/s4lfanet/go-api-c320/config"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"github.com/rs/zerolog/log"
)

// ProvisionUseCaseInterface defines the interface for provisioning operations
type ProvisionUseCaseInterface interface {
	// ONU Discovery
	GetUnconfiguredONUs(ctx context.Context, ponPort string) ([]model.UnconfiguredONU, error)
	GetAllUnconfiguredONUs(ctx context.Context) ([]model.UnconfiguredONU, error)

	// ONU Registration
	RegisterONU(ctx context.Context, req model.ONURegistrationRequest) (*model.ONURegistrationResponse, error)
	DeleteONU(ctx context.Context, ponPort string, onuID int) error

	// ONU Configuration
	ConfigureTCONT(ctx context.Context, ponPort string, onuID int, tcontID int, profileName string) error
	ConfigureGEMPort(ctx context.Context, ponPort string, onuID int, gemportID int, tcontID int) error
	ConfigureServicePort(ctx context.Context, ponPort string, onuID int, gemportID int, vlan int, userVlan string) error
}

// ProvisionUsecase implements provisioning business logic
type ProvisionUsecase struct {
	sessionManager *repository.TelnetSessionManager
	config         *config.Config
}

// NewProvisionUsecase creates a new provision usecase instance
func NewProvisionUsecase(sessionManager *repository.TelnetSessionManager, cfg *config.Config) ProvisionUseCaseInterface {
	return &ProvisionUsecase{
		sessionManager: sessionManager,
		config:         cfg,
	}
}

// GetUnconfiguredONUs retrieves unconfigured ONUs for a specific PON port
func (u *ProvisionUsecase) GetUnconfiguredONUs(ctx context.Context, ponPort string) ([]model.UnconfiguredONU, error) {
	log.Info().Str("pon_port", ponPort).Msg("Getting unconfigured ONUs")

	// Execute show command
	command := fmt.Sprintf("show gpon onu uncfg gpon-olt_%s", ponPort)
	resp, err := u.sessionManager.ExecuteCommand(ctx, command)
	if err != nil {
		log.Error().Err(err).Str("command", command).Msg("Failed to get unconfigured ONUs")
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("command failed: %s", resp.Error)
	}

	// Parse the output
	onus := u.parseUnconfiguredONUs(resp.Output, ponPort)

	log.Info().
		Str("pon_port", ponPort).
		Int("count", len(onus)).
		Msg("Retrieved unconfigured ONUs")

	return onus, nil
}

// GetAllUnconfiguredONUs retrieves all unconfigured ONUs across all PON ports
func (u *ProvisionUsecase) GetAllUnconfiguredONUs(ctx context.Context) ([]model.UnconfiguredONU, error) {
	log.Info().Msg("Getting all unconfigured ONUs")

	// Execute show command without PON filter
	command := "show gpon onu uncfg"
	resp, err := u.sessionManager.ExecuteCommand(ctx, command)
	if err != nil {
		log.Error().Err(err).Str("command", command).Msg("Failed to get unconfigured ONUs")
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("command failed: %s", resp.Error)
	}

	// Parse the output
	onus := u.parseUnconfiguredONUs(resp.Output, "")

	log.Info().Int("count", len(onus)).Msg("Retrieved all unconfigured ONUs")

	return onus, nil
}

// parseUnconfiguredONUs parses the output of "show gpon onu uncfg" command
func (u *ProvisionUsecase) parseUnconfiguredONUs(output string, filterPonPort string) []model.UnconfiguredONU {
	onus := make([]model.UnconfiguredONU, 0)

	// V2.1.0 Output format:
	// OnuIndex                 Sn                  State
	// ---------------------------------------------------------------------
	// gpon-onu_1/1/1:1         HWTC1F14CAAD        unknown
	// gpon-onu_1/1/1:2         ZTEGD824CDF3        unknown
	// gpon-onu_1/1/1:3         ZTEGDA5918AC        unknown

	lines := strings.Split(output, "\n")

	// Regex to match ONU line: gpon-onu_X/X/X:Y (SerialNumber) (State)
	onuRegex := regexp.MustCompile(`gpon-onu_(\d+/\d+/\d+):(\d+)\s+([A-Z0-9]+)\s+(\w+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "OnuIndex") || strings.HasPrefix(line, "----") {
			continue
		}

		matches := onuRegex.FindStringSubmatch(line)
		if len(matches) >= 5 {
			ponPort := matches[1]        // e.g., "1/1/1"
			onuID := matches[2]          // e.g., "1"
			serialNumber := matches[3]   // e.g., "HWTC1F14CAAD"
			state := matches[4]          // e.g., "unknown"

			// Filter by PON port if specified
			if filterPonPort != "" && ponPort != filterPonPort {
				continue
			}

			// Determine ONU type from serial number (first 4 chars are vendor)
			onuType := "Unknown"
			if len(serialNumber) >= 4 {
				vendor := serialNumber[:4]
				onuType = u.guessONUType(vendor, serialNumber)
			}

			onus = append(onus, model.UnconfiguredONU{
				PONPort:      ponPort,
				SerialNumber: serialNumber,
				Type:         onuType,
				DiscoveredAt: time.Now().Format(time.RFC3339),
				LOID:         "",  // V2.1.0 doesn't show LOID in uncfg list
			})

			log.Debug().
				Str("pon", ponPort).
				Str("onu_id", onuID).
				Str("sn", serialNumber).
				Str("state", state).
				Msg("Parsed unconfigured ONU")
		}
	}

	return onus
}

// guessONUType attempts to guess ONU type from serial number
func (u *ProvisionUsecase) guessONUType(vendor, serialNumber string) string {
	switch vendor {
	case "ZTEG":
		// ZTE ONUs - try to determine model from serial pattern
		return "ZTE (Auto-detect)"
	case "HWTC":
		return "Huawei (Auto-detect)"
	case "FIBR":
		return "FiberHome (Auto-detect)"
	case "ALCL":
		return "Nokia/Alcatel-Lucent"
	default:
		return fmt.Sprintf("Unknown (%s)", vendor)
	}
}

// RegisterONU registers a new ONU to the OLT
func (u *ProvisionUsecase) RegisterONU(ctx context.Context, req model.ONURegistrationRequest) (*model.ONURegistrationResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Str("serial", req.SerialNumber).
		Str("type", req.ONUType).
		Msg("Registering ONU")

	// Prepare commands
	commands := []string{
		fmt.Sprintf("interface gpon-olt_%s", req.PONPort),
		fmt.Sprintf("onu %d type %s sn %s", req.ONUID, req.ONUType, req.SerialNumber),
	}

	// Add name if provided
	if req.Name != "" {
		commands = append(commands, fmt.Sprintf("onu %d name \"%s\"", req.ONUID, req.Name))
	}

	commands = append(commands, "exit")

	// Execute in config mode
	result, err := u.sessionManager.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register ONU")
		return &model.ONURegistrationResponse{
			PONPort:      req.PONPort,
			ONUID:        req.ONUID,
			SerialNumber: req.SerialNumber,
			Success:      false,
			Message:      fmt.Sprintf("Failed to register ONU: %v", err),
		}, err
	}

	// Check for errors in responses
	for _, resp := range result.Responses {
		if !resp.Success || strings.Contains(strings.ToLower(resp.Output), "error") {
			log.Error().
				Str("command", resp.Command).
				Str("output", resp.Output).
				Msg("Command execution failed")
			return &model.ONURegistrationResponse{
				PONPort:      req.PONPort,
				ONUID:        req.ONUID,
				SerialNumber: req.SerialNumber,
				Success:      false,
				Message:      fmt.Sprintf("ONU registration failed: %s", resp.Output),
			}, fmt.Errorf("registration failed")
		}
	}

	// Configure TCONT and GEMPORT if profile is provided
	if req.Profile.DBAProfile != "" {
		log.Info().Msg("Configuring ONU with profile")

		// Configure TCONT
		if err := u.ConfigureTCONT(ctx, req.PONPort, req.ONUID, 1, req.Profile.DBAProfile); err != nil {
			log.Error().Err(err).Msg("Failed to configure TCONT")
			// Don't fail registration, just log the error
		}

		// Configure GEMPORT
		if err := u.ConfigureGEMPort(ctx, req.PONPort, req.ONUID, 1, 1); err != nil {
			log.Error().Err(err).Msg("Failed to configure GEMPORT")
		}

		// Configure service port on ONU
		if req.Profile.VLAN > 0 {
			if err := u.ConfigureServicePort(ctx, req.PONPort, req.ONUID, 1, req.Profile.VLAN, "untagged"); err != nil {
				log.Error().Err(err).Msg("Failed to configure service port")
			}
		}
	}

	// Save configuration
	if err := u.sessionManager.SaveConfiguration(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to save configuration")
		// Don't fail registration
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU registered successfully")

	return &model.ONURegistrationResponse{
		PONPort:      req.PONPort,
		ONUID:        req.ONUID,
		SerialNumber: req.SerialNumber,
		Success:      true,
		Message:      "ONU registered successfully",
	}, nil
}

// DeleteONU deletes an ONU from the OLT
func (u *ProvisionUsecase) DeleteONU(ctx context.Context, ponPort string, onuID int) error {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("Deleting ONU")

	commands := []string{
		fmt.Sprintf("interface gpon-olt_%s", ponPort),
		fmt.Sprintf("no onu %d", onuID),
		"exit",
	}

	result, err := u.sessionManager.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete ONU")
		return err
	}

	// Check for errors
	for _, resp := range result.Responses {
		if !resp.Success {
			log.Error().
				Str("command", resp.Command).
				Str("output", resp.Output).
				Msg("Delete command failed")
			return fmt.Errorf("delete failed: %s", resp.Output)
		}
	}

	// Save configuration
	if err := u.sessionManager.SaveConfiguration(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to save configuration")
		return err
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("ONU deleted successfully")

	return nil
}

// ConfigureTCONT configures TCONT for an ONU
func (u *ProvisionUsecase) ConfigureTCONT(ctx context.Context, ponPort string, onuID int, tcontID int, profileName string) error {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("tcont_id", tcontID).
		Str("profile", profileName).
		Msg("Configuring TCONT")

	commands := []string{
		fmt.Sprintf("interface gpon-onu_%s:%d", ponPort, onuID),
		fmt.Sprintf("tcont %d name TCONT_%d profile %s", tcontID, tcontID, profileName),
		"exit",
	}

	result, err := u.sessionManager.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		log.Error().Err(err).Msg("Failed to configure TCONT")
		return err
	}

	// Check for errors
	for _, resp := range result.Responses {
		if strings.Contains(strings.ToLower(resp.Output), "error") {
			return fmt.Errorf("TCONT configuration failed: %s", resp.Output)
		}
	}

	log.Info().Msg("TCONT configured successfully")
	return nil
}

// ConfigureGEMPort configures GEMPORT for an ONU
func (u *ProvisionUsecase) ConfigureGEMPort(ctx context.Context, ponPort string, onuID int, gemportID int, tcontID int) error {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("gemport_id", gemportID).
		Int("tcont_id", tcontID).
		Msg("Configuring GEMPORT")

	commands := []string{
		fmt.Sprintf("interface gpon-onu_%s:%d", ponPort, onuID),
		fmt.Sprintf("gemport %d name GEM_%d tcont %d", gemportID, gemportID, tcontID),
		"exit",
	}

	result, err := u.sessionManager.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		log.Error().Err(err).Msg("Failed to configure GEMPORT")
		return err
	}

	// Check for errors
	for _, resp := range result.Responses {
		if strings.Contains(strings.ToLower(resp.Output), "error") {
			return fmt.Errorf("GEMPORT configuration failed: %s", resp.Output)
		}
	}

	log.Info().Msg("GEMPORT configured successfully")
	return nil
}

// ConfigureServicePort configures service port on ONU (VLAN mapping)
func (u *ProvisionUsecase) ConfigureServicePort(ctx context.Context, ponPort string, onuID int, gemportID int, vlan int, userVlan string) error {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("gemport_id", gemportID).
		Int("vlan", vlan).
		Str("user_vlan", userVlan).
		Msg("Configuring service port")

	commands := []string{
		fmt.Sprintf("interface gpon-onu_%s:%d", ponPort, onuID),
		fmt.Sprintf("service-port 1 vport 1 user-vlan %s vlan %d", userVlan, vlan),
		"exit",
	}

	result, err := u.sessionManager.ExecuteInConfigMode(ctx, commands)
	if err != nil {
		log.Error().Err(err).Msg("Failed to configure service port")
		return err
	}

	// Check for errors
	for _, resp := range result.Responses {
		if strings.Contains(strings.ToLower(resp.Output), "error") {
			return fmt.Errorf("service port configuration failed: %s", resp.Output)
		}
	}

	log.Info().Msg("Service port configured successfully")
	return nil
}
