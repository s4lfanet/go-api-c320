package usecase

import (
	"context"
	"fmt"

	"github.com/s4lfanet/go-api-c320/config"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"github.com/rs/zerolog/log"
)

// TrafficUsecaseInterface defines the interface for traffic profile usecase
type TrafficUsecaseInterface interface {
	// DBA Profile operations
	GetDBAProfile(ctx context.Context, name string) (*model.DBAProfileInfo, error)
	GetAllDBAProfiles(ctx context.Context) ([]model.DBAProfileInfo, error)
	CreateDBAProfile(ctx context.Context, req model.DBAProfileRequest) (*model.DBAProfileResponse, error)
	ModifyDBAProfile(ctx context.Context, req model.DBAProfileRequest) (*model.DBAProfileResponse, error)
	DeleteDBAProfile(ctx context.Context, name string) error

	// TCONT operations
	GetONUTCONT(ctx context.Context, ponPort string, onuID int, tcontID int) (*model.TCONTInfo, error)
	ConfigureTCONT(ctx context.Context, req model.TCONTConfigRequest) (*model.TCONTConfigResponse, error)
	DeleteTCONT(ctx context.Context, ponPort string, onuID int, tcontID int) error

	// GEMPort operations
	ConfigureGEMPort(ctx context.Context, req model.GEMPortConfigRequest) (*model.GEMPortConfigResponse, error)
	DeleteGEMPort(ctx context.Context, ponPort string, onuID int, gemportID int) error
}

// TrafficUsecase implements the traffic profile usecase interface
type TrafficUsecase struct {
	telnetSessionManager *repository.TelnetSessionManager
	config               *config.Config
}

// NewTrafficUsecase creates a new traffic usecase instance
func NewTrafficUsecase(telnetSessionManager *repository.TelnetSessionManager, cfg *config.Config) TrafficUsecaseInterface {
	return &TrafficUsecase{
		telnetSessionManager: telnetSessionManager,
		config:               cfg,
	}
}

// GetDBAProfile retrieves DBA profile information
func (u *TrafficUsecase) GetDBAProfile(ctx context.Context, name string) (*model.DBAProfileInfo, error) {
	log.Info().
		Str("profile_name", name).
		Msg("Getting DBA profile")

	// Validate profile name
	if err := validateProfileName(name); err != nil {
		return nil, fmt.Errorf("invalid profile name: %w", err)
	}

	profile, err := u.telnetSessionManager.GetDBAProfile(ctx, name)
	if err != nil {
		log.Error().
			Err(err).
			Str("profile_name", name).
			Msg("Failed to get DBA profile")
		return nil, err
	}

	log.Info().
		Str("profile_name", name).
		Int("type", profile.Type).
		Msg("Retrieved DBA profile")

	return profile, nil
}

// GetAllDBAProfiles retrieves all DBA profiles
func (u *TrafficUsecase) GetAllDBAProfiles(ctx context.Context) ([]model.DBAProfileInfo, error) {
	log.Info().Msg("Getting all DBA profiles")

	profiles, err := u.telnetSessionManager.GetAllDBAProfiles(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get DBA profiles")
		return nil, err
	}

	log.Info().
		Int("count", len(profiles)).
		Msg("Retrieved DBA profiles")

	return profiles, nil
}

// CreateDBAProfile creates a new DBA profile
func (u *TrafficUsecase) CreateDBAProfile(ctx context.Context, req model.DBAProfileRequest) (*model.DBAProfileResponse, error) {
	log.Info().
		Str("profile_name", req.Name).
		Int("type", req.Type).
		Int("assured_bandwidth", req.AssuredBandwidth).
		Int("max_bandwidth", req.MaxBandwidth).
		Msg("Creating DBA profile")

	// Validate request
	if err := validateDBAProfileRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	response, err := u.telnetSessionManager.CreateDBAProfile(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("profile_name", req.Name).
			Msg("Failed to create DBA profile")
		return response, err
	}

	log.Info().
		Str("profile_name", req.Name).
		Bool("success", response.Success).
		Msg("DBA profile creation completed")

	return response, nil
}

// ModifyDBAProfile modifies an existing DBA profile
func (u *TrafficUsecase) ModifyDBAProfile(ctx context.Context, req model.DBAProfileRequest) (*model.DBAProfileResponse, error) {
	log.Info().
		Str("profile_name", req.Name).
		Int("type", req.Type).
		Msg("Modifying DBA profile")

	// Validate request
	if err := validateDBAProfileRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if profile exists
	_, err := u.telnetSessionManager.GetDBAProfile(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("DBA profile not found: %s", req.Name)
	}

	response, err := u.telnetSessionManager.ModifyDBAProfile(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("profile_name", req.Name).
			Msg("Failed to modify DBA profile")
		return response, err
	}

	log.Info().
		Str("profile_name", req.Name).
		Bool("success", response.Success).
		Msg("DBA profile modification completed")

	return response, nil
}

// DeleteDBAProfile deletes a DBA profile
func (u *TrafficUsecase) DeleteDBAProfile(ctx context.Context, name string) error {
	log.Info().
		Str("profile_name", name).
		Msg("Deleting DBA profile")

	// Validate profile name
	if err := validateProfileName(name); err != nil {
		return fmt.Errorf("invalid profile name: %w", err)
	}

	err := u.telnetSessionManager.DeleteDBAProfile(ctx, name)
	if err != nil {
		log.Error().
			Err(err).
			Str("profile_name", name).
			Msg("Failed to delete DBA profile")
		return err
	}

	log.Info().
		Str("profile_name", name).
		Msg("DBA profile deleted successfully")

	return nil
}

// GetONUTCONT retrieves T-CONT configuration for an ONU
func (u *TrafficUsecase) GetONUTCONT(ctx context.Context, ponPort string, onuID int, tcontID int) (*model.TCONTInfo, error) {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("tcont_id", tcontID).
		Msg("Getting TCONT configuration")

	// Validate inputs
	if err := validatePONPort(ponPort); err != nil {
		return nil, fmt.Errorf("invalid PON port: %w", err)
	}

	if err := validateONUID(onuID); err != nil {
		return nil, fmt.Errorf("invalid ONU ID: %w", err)
	}

	if err := validateTCONTID(tcontID); err != nil {
		return nil, fmt.Errorf("invalid TCONT ID: %w", err)
	}

	tcont, err := u.telnetSessionManager.GetONUTCONT(ctx, ponPort, onuID, tcontID)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Int("tcont_id", tcontID).
			Msg("Failed to get TCONT")
		return nil, err
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("tcont_id", tcontID).
		Str("profile", tcont.Profile).
		Msg("Retrieved TCONT configuration")

	return tcont, nil
}

// ConfigureTCONT configures T-CONT for an ONU
func (u *TrafficUsecase) ConfigureTCONT(ctx context.Context, req model.TCONTConfigRequest) (*model.TCONTConfigResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("tcont_id", req.TCONTID).
		Str("profile", req.Profile).
		Msg("Configuring TCONT")

	// Validate request
	if err := validateTCONTRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	response, err := u.telnetSessionManager.ConfigureTCONT(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Int("tcont_id", req.TCONTID).
			Msg("Failed to configure TCONT")
		return response, err
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("tcont_id", req.TCONTID).
		Bool("success", response.Success).
		Msg("TCONT configuration completed")

	return response, nil
}

// DeleteTCONT deletes T-CONT from an ONU
func (u *TrafficUsecase) DeleteTCONT(ctx context.Context, ponPort string, onuID int, tcontID int) error {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("tcont_id", tcontID).
		Msg("Deleting TCONT")

	// Validate inputs
	if err := validatePONPort(ponPort); err != nil {
		return fmt.Errorf("invalid PON port: %w", err)
	}

	if err := validateONUID(onuID); err != nil {
		return fmt.Errorf("invalid ONU ID: %w", err)
	}

	if err := validateTCONTID(tcontID); err != nil {
		return fmt.Errorf("invalid TCONT ID: %w", err)
	}

	err := u.telnetSessionManager.DeleteTCONT(ctx, ponPort, onuID, tcontID)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Int("tcont_id", tcontID).
			Msg("Failed to delete TCONT")
		return err
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("tcont_id", tcontID).
		Msg("TCONT deleted successfully")

	return nil
}

// ConfigureGEMPort configures GEM port for an ONU
func (u *TrafficUsecase) ConfigureGEMPort(ctx context.Context, req model.GEMPortConfigRequest) (*model.GEMPortConfigResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("gemport_id", req.GEMPortID).
		Int("tcont_id", req.TCONTID).
		Msg("Configuring GEM port")

	// Validate request
	if err := validateGEMPortRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	response, err := u.telnetSessionManager.ConfigureGEMPort(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Int("gemport_id", req.GEMPortID).
			Msg("Failed to configure GEM port")
		return response, err
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("gemport_id", req.GEMPortID).
		Bool("success", response.Success).
		Msg("GEM port configuration completed")

	return response, nil
}

// DeleteGEMPort deletes GEM port from an ONU
func (u *TrafficUsecase) DeleteGEMPort(ctx context.Context, ponPort string, onuID int, gemportID int) error {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("gemport_id", gemportID).
		Msg("Deleting GEM port")

	// Validate inputs
	if err := validatePONPort(ponPort); err != nil {
		return fmt.Errorf("invalid PON port: %w", err)
	}

	if err := validateONUID(onuID); err != nil {
		return fmt.Errorf("invalid ONU ID: %w", err)
	}

	if err := validateGEMPortID(gemportID); err != nil {
		return fmt.Errorf("invalid GEM port ID: %w", err)
	}

	err := u.telnetSessionManager.DeleteGEMPort(ctx, ponPort, onuID, gemportID)
	if err != nil {
		log.Error().
			Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Int("gemport_id", gemportID).
			Msg("Failed to delete GEM port")
		return err
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("gemport_id", gemportID).
		Msg("GEM port deleted successfully")

	return nil
}

// Validation functions

func validateProfileName(name string) error {
	if name == "" {
		return fmt.Errorf("profile name is required")
	}
	if len(name) > 32 {
		return fmt.Errorf("profile name too long (max 32 characters)")
	}
	return nil
}

func validateDBAProfileRequest(req model.DBAProfileRequest) error {
	if err := validateProfileName(req.Name); err != nil {
		return err
	}

	if req.Type < 1 || req.Type > 5 {
		return fmt.Errorf("invalid DBA profile type (must be 1-5)")
	}

	// Validate bandwidth values based on type
	switch req.Type {
	case 1: // Fixed
		if req.FixedBandwidth < 64 {
			return fmt.Errorf("fixed bandwidth must be at least 64 Kbps")
		}
	case 2: // Assured
		if req.AssuredBandwidth < 64 {
			return fmt.Errorf("assured bandwidth must be at least 64 Kbps")
		}
	case 3, 5: // Assured + Max
		if req.AssuredBandwidth < 64 {
			return fmt.Errorf("assured bandwidth must be at least 64 Kbps")
		}
		if req.MaxBandwidth < 64 {
			return fmt.Errorf("max bandwidth must be at least 64 Kbps")
		}
		if req.AssuredBandwidth > req.MaxBandwidth {
			return fmt.Errorf("assured bandwidth cannot exceed max bandwidth")
		}
	case 4: // Max only
		if req.MaxBandwidth < 64 {
			return fmt.Errorf("max bandwidth must be at least 64 Kbps")
		}
	}

	return nil
}

func validateTCONTID(tcontID int) error {
	if tcontID < 1 || tcontID > 8 {
		return fmt.Errorf("TCONT ID must be between 1 and 8")
	}
	return nil
}

func validateTCONTRequest(req model.TCONTConfigRequest) error {
	if err := validatePONPort(req.PONPort); err != nil {
		return err
	}

	if err := validateONUID(req.ONUID); err != nil {
		return err
	}

	if err := validateTCONTID(req.TCONTID); err != nil {
		return err
	}

	if err := validateProfileName(req.Profile); err != nil {
		return fmt.Errorf("invalid DBA profile: %w", err)
	}

	if req.Name != "" && len(req.Name) > 32 {
		return fmt.Errorf("TCONT name too long (max 32 characters)")
	}

	return nil
}

func validateGEMPortID(gemportID int) error {
	if gemportID < 1 || gemportID > 128 {
		return fmt.Errorf("GEM port ID must be between 1 and 128")
	}
	return nil
}

func validateGEMPortRequest(req model.GEMPortConfigRequest) error {
	if err := validatePONPort(req.PONPort); err != nil {
		return err
	}

	if err := validateONUID(req.ONUID); err != nil {
		return err
	}

	if err := validateGEMPortID(req.GEMPortID); err != nil {
		return err
	}

	if err := validateTCONTID(req.TCONTID); err != nil {
		return err
	}

	if req.Name != "" && len(req.Name) > 32 {
		return fmt.Errorf("GEM port name too long (max 32 characters)")
	}

	if req.Queue < 0 || req.Queue > 8 {
		return fmt.Errorf("queue must be between 0 and 8")
	}

	return nil
}
