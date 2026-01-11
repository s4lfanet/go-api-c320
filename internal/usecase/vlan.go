package usecase

import (
	"context"
	"fmt"
	"regexp"

	"github.com/s4lfanet/go-api-c320/config"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"github.com/rs/zerolog/log"
)

// VLANUsecaseInterface defines the interface for VLAN usecase
type VLANUsecaseInterface interface {
	GetONUVLAN(ctx context.Context, ponPort string, onuID int) (*model.ONUVLANInfo, error)
	GetAllServicePorts(ctx context.Context) ([]model.ONUVLANInfo, error)
	ConfigureVLAN(ctx context.Context, req model.VLANConfigRequest) (*model.VLANConfigResponse, error)
	ModifyVLAN(ctx context.Context, req model.VLANConfigRequest) (*model.VLANConfigResponse, error)
	DeleteVLAN(ctx context.Context, ponPort string, onuID int) error
}

// VLANUsecase implements the VLAN usecase interface
type VLANUsecase struct {
	telnetSessionManager *repository.TelnetSessionManager
	config               *config.Config
}

// NewVLANUsecase creates a new VLAN usecase instance
func NewVLANUsecase(telnetSessionManager *repository.TelnetSessionManager, cfg *config.Config) VLANUsecaseInterface {
	return &VLANUsecase{
		telnetSessionManager: telnetSessionManager,
		config:               cfg,
	}
}

// GetONUVLAN retrieves VLAN configuration for a specific ONU
func (u *VLANUsecase) GetONUVLAN(ctx context.Context, ponPort string, onuID int) (*model.ONUVLANInfo, error) {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("Getting ONU VLAN configuration")

	// Validate inputs
	if err := validatePONPort(ponPort); err != nil {
		return nil, fmt.Errorf("invalid PON port: %w", err)
	}

	if err := validateONUID(onuID); err != nil {
		return nil, fmt.Errorf("invalid ONU ID: %w", err)
	}

	vlanInfo, err := u.telnetSessionManager.GetONUVLAN(ctx, ponPort, onuID)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Msg("Failed to get ONU VLAN")
		return nil, err
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("svlan", vlanInfo.SVLAN).
		Int("cvlan", vlanInfo.CVLAN).
		Msg("Retrieved ONU VLAN configuration")

	return vlanInfo, nil
}

// GetAllServicePorts retrieves all service-port configurations
func (u *VLANUsecase) GetAllServicePorts(ctx context.Context) ([]model.ONUVLANInfo, error) {
	log.Info().Msg("Getting all service-port configurations")

	servicePorts, err := u.telnetSessionManager.GetAllServicePorts(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get service-ports")
		return nil, err
	}

	log.Info().
		Int("count", len(servicePorts)).
		Msg("Retrieved service-port configurations")

	return servicePorts, nil
}

// ConfigureVLAN configures VLAN for an ONU
func (u *VLANUsecase) ConfigureVLAN(ctx context.Context, req model.VLANConfigRequest) (*model.VLANConfigResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("svlan", req.SVLAN).
		Int("cvlan", req.CVLAN).
		Str("vlan_mode", req.VLANMode).
		Msg("Configuring ONU VLAN")

	// Validate request
	if err := u.validateVLANRequest(&req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Configure VLAN via telnet
	response, err := u.telnetSessionManager.ConfigureONUVLAN(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Msg("Failed to configure VLAN")
		return nil, err
	}

	if !response.Success {
		log.Warn().
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Str("message", response.Message).
			Msg("VLAN configuration returned unsuccessful")
		return response, fmt.Errorf("configuration failed: %s", response.Message)
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("service_port_id", response.ServicePortID).
		Msg("VLAN configured successfully")

	return response, nil
}

// ModifyVLAN modifies existing VLAN configuration for an ONU
func (u *VLANUsecase) ModifyVLAN(ctx context.Context, req model.VLANConfigRequest) (*model.VLANConfigResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("svlan", req.SVLAN).
		Msg("Modifying ONU VLAN")

	// Validate request
	if err := u.validateVLANRequest(&req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if VLAN exists
	existingVLAN, err := u.telnetSessionManager.GetONUVLAN(ctx, req.PONPort, req.ONUID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing VLAN: %w", err)
	}

	if existingVLAN.ServicePortID == 0 {
		return nil, fmt.Errorf("no existing VLAN configuration found for ONU %s:%d", req.PONPort, req.ONUID)
	}

	// Update service port ID in request
	req.ServicePortID = existingVLAN.ServicePortID

	// Configure (update) VLAN via telnet
	response, err := u.telnetSessionManager.ConfigureONUVLAN(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Msg("Failed to modify VLAN")
		return nil, err
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("service_port_id", response.ServicePortID).
		Msg("VLAN modified successfully")

	return response, nil
}

// DeleteVLAN removes VLAN configuration for an ONU
func (u *VLANUsecase) DeleteVLAN(ctx context.Context, ponPort string, onuID int) error {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("Deleting ONU VLAN")

	// Validate inputs
	if err := validatePONPort(ponPort); err != nil {
		return fmt.Errorf("invalid PON port: %w", err)
	}

	if err := validateONUID(onuID); err != nil {
		return fmt.Errorf("invalid ONU ID: %w", err)
	}

	// Delete VLAN via telnet
	err := u.telnetSessionManager.DeleteONUVLAN(ctx, ponPort, onuID)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Msg("Failed to delete VLAN")
		return err
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("VLAN deleted successfully")

	return nil
}

// validateVLANRequest validates VLAN configuration request
func (u *VLANUsecase) validateVLANRequest(req *model.VLANConfigRequest) error {
	// Validate PON port
	if err := validatePONPort(req.PONPort); err != nil {
		return fmt.Errorf("invalid PON port: %w", err)
	}

	// Validate ONU ID
	if err := validateONUID(req.ONUID); err != nil {
		return fmt.Errorf("invalid ONU ID: %w", err)
	}

	// Validate SVLAN
	if req.SVLAN < 1 || req.SVLAN > 4094 {
		return fmt.Errorf("SVLAN must be between 1 and 4094")
	}

	// Validate CVLAN if provided
	if req.CVLAN > 0 && (req.CVLAN < 1 || req.CVLAN > 4094) {
		return fmt.Errorf("CVLAN must be between 1 and 4094")
	}

	// Validate VLAN mode
	validModes := map[string]bool{
		"tag":         true,
		"translation": true,
		"transparent": true,
	}
	if !validModes[req.VLANMode] {
		return fmt.Errorf("invalid VLAN mode: must be 'tag', 'translation', or 'transparent'")
	}

	// Validate priority
	if req.Priority < 0 || req.Priority > 7 {
		return fmt.Errorf("priority must be between 0 and 7")
	}

	// CVLAN is required for translation mode
	if req.VLANMode == "translation" && req.CVLAN == 0 {
		return fmt.Errorf("CVLAN is required for translation mode")
	}

	return nil
}

// validatePONPort validates PON port format (e.g., 1/1/1)
func validatePONPort(ponPort string) error {
	// PON port format: rack/shelf/slot (e.g., 1/1/1)
	pattern := regexp.MustCompile(`^\d+/\d+/\d+$`)
	if !pattern.MatchString(ponPort) {
		return fmt.Errorf("invalid PON port format, expected: rack/shelf/slot (e.g., 1/1/1)")
	}
	return nil
}

// validateONUID validates ONU ID range
func validateONUID(onuID int) error {
	if onuID < 1 || onuID > 128 {
		return fmt.Errorf("ONU ID must be between 1 and 128")
	}
	return nil
}
