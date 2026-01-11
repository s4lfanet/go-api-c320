package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/s4lfanet/go-api-c320/config"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"github.com/rs/zerolog/log"
)

// ONUManagementUsecaseInterface defines business logic for ONU lifecycle management
type ONUManagementUsecaseInterface interface {
	RebootONU(ctx context.Context, req *model.ONURebootRequest) (*model.ONURebootResponse, error)
	BlockONU(ctx context.Context, req *model.ONUBlockRequest) (*model.ONUBlockResponse, error)
	UnblockONU(ctx context.Context, req *model.ONUBlockRequest) (*model.ONUBlockResponse, error)
	UpdateDescription(ctx context.Context, req *model.ONUDescriptionRequest) (*model.ONUDescriptionResponse, error)
	DeleteONU(ctx context.Context, req *model.ONUDeleteRequest) (*model.ONUDeleteResponse, error)
}

// ONUManagementUsecase implements ONU management business logic
type ONUManagementUsecase struct {
	telnetSessionManager *repository.TelnetSessionManager
	config               *config.Config
}

// NewONUManagementUsecase creates a new ONU management usecase
func NewONUManagementUsecase(
	telnetSessionManager *repository.TelnetSessionManager,
	cfg *config.Config,
) ONUManagementUsecaseInterface {
	return &ONUManagementUsecase{
		telnetSessionManager: telnetSessionManager,
		config:               cfg,
	}
}

// RebootONU reboots an ONU
func (u *ONUManagementUsecase) RebootONU(ctx context.Context, req *model.ONURebootRequest) (*model.ONURebootResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("Processing ONU reboot request")

	// Validate request
	if err := u.validateRebootRequest(req); err != nil {
		log.Error().Err(err).Msg("Reboot request validation failed")
		return &model.ONURebootResponse{
			PONPort: req.PONPort,
			ONUID:   req.ONUID,
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Execute reboot
	if err := u.telnetSessionManager.RebootONU(ctx, req); err != nil {
		log.Error().Err(err).Msg("Failed to reboot ONU")
		return &model.ONURebootResponse{
			PONPort: req.PONPort,
			ONUID:   req.ONUID,
			Success: false,
			Message: fmt.Sprintf("Failed to reboot ONU: %v", err),
		}, err
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU rebooted successfully")

	return &model.ONURebootResponse{
		PONPort: req.PONPort,
		ONUID:   req.ONUID,
		Success: true,
		Message: "ONU reboot command executed successfully. ONU will restart shortly.",
	}, nil
}

// BlockONU blocks (disables) an ONU
func (u *ONUManagementUsecase) BlockONU(ctx context.Context, req *model.ONUBlockRequest) (*model.ONUBlockResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("Processing ONU block request")

	// Validate request
	if err := u.validateBlockRequest(req); err != nil {
		log.Error().Err(err).Msg("Block request validation failed")
		return &model.ONUBlockResponse{
			PONPort: req.PONPort,
			ONUID:   req.ONUID,
			Blocked: false,
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Execute block
	if err := u.telnetSessionManager.BlockONU(ctx, req); err != nil {
		log.Error().Err(err).Msg("Failed to block ONU")
		return &model.ONUBlockResponse{
			PONPort: req.PONPort,
			ONUID:   req.ONUID,
			Blocked: false,
			Success: false,
			Message: fmt.Sprintf("Failed to block ONU: %v", err),
		}, err
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU blocked successfully")

	return &model.ONUBlockResponse{
		PONPort: req.PONPort,
		ONUID:   req.ONUID,
		Blocked: true,
		Success: true,
		Message: "ONU blocked successfully. ONU is now disabled.",
	}, nil
}

// UnblockONU unblocks (enables) an ONU
func (u *ONUManagementUsecase) UnblockONU(ctx context.Context, req *model.ONUBlockRequest) (*model.ONUBlockResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("Processing ONU unblock request")

	// Validate request
	if err := u.validateBlockRequest(req); err != nil {
		log.Error().Err(err).Msg("Unblock request validation failed")
		return &model.ONUBlockResponse{
			PONPort: req.PONPort,
			ONUID:   req.ONUID,
			Blocked: true,
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Execute unblock
	if err := u.telnetSessionManager.UnblockONU(ctx, req); err != nil {
		log.Error().Err(err).Msg("Failed to unblock ONU")
		return &model.ONUBlockResponse{
			PONPort: req.PONPort,
			ONUID:   req.ONUID,
			Blocked: true,
			Success: false,
			Message: fmt.Sprintf("Failed to unblock ONU: %v", err),
		}, err
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU unblocked successfully")

	return &model.ONUBlockResponse{
		PONPort: req.PONPort,
		ONUID:   req.ONUID,
		Blocked: false,
		Success: true,
		Message: "ONU unblocked successfully. ONU is now enabled.",
	}, nil
}

// UpdateDescription updates the name/description of an ONU
func (u *ONUManagementUsecase) UpdateDescription(ctx context.Context, req *model.ONUDescriptionRequest) (*model.ONUDescriptionResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Str("description", req.Description).
		Msg("Processing ONU description update request")

	// Validate request
	if err := u.validateDescriptionRequest(req); err != nil {
		log.Error().Err(err).Msg("Description update request validation failed")
		return &model.ONUDescriptionResponse{
			PONPort:     req.PONPort,
			ONUID:       req.ONUID,
			Description: req.Description,
			Success:     false,
			Message:     err.Error(),
		}, err
	}

	// Execute description update
	if err := u.telnetSessionManager.UpdateDescription(ctx, req); err != nil {
		log.Error().Err(err).Msg("Failed to update ONU description")
		return &model.ONUDescriptionResponse{
			PONPort:     req.PONPort,
			ONUID:       req.ONUID,
			Description: req.Description,
			Success:     false,
			Message:     fmt.Sprintf("Failed to update ONU description: %v", err),
		}, err
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Str("description", req.Description).
		Msg("ONU description updated successfully")

	return &model.ONUDescriptionResponse{
		PONPort:     req.PONPort,
		ONUID:       req.ONUID,
		Description: req.Description,
		Success:     true,
		Message:     "ONU description updated successfully",
	}, nil
}

// DeleteONU deletes ONU configuration
func (u *ONUManagementUsecase) DeleteONU(ctx context.Context, req *model.ONUDeleteRequest) (*model.ONUDeleteResponse, error) {
	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("Processing ONU deletion request")

	// Validate request
	if err := u.validateDeleteRequest(req); err != nil {
		log.Error().Err(err).Msg("Delete request validation failed")
		return &model.ONUDeleteResponse{
			PONPort: req.PONPort,
			ONUID:   req.ONUID,
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Execute deletion
	if err := u.telnetSessionManager.DeleteONU(ctx, req); err != nil {
		log.Error().Err(err).Msg("Failed to delete ONU")
		return &model.ONUDeleteResponse{
			PONPort: req.PONPort,
			ONUID:   req.ONUID,
			Success: false,
			Message: fmt.Sprintf("Failed to delete ONU: %v", err),
		}, err
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU deleted successfully")

	return &model.ONUDeleteResponse{
		PONPort: req.PONPort,
		ONUID:   req.ONUID,
		Success: true,
		Message: "ONU configuration deleted successfully",
	}, nil
}

// Validation functions

func (u *ONUManagementUsecase) validateRebootRequest(req *model.ONURebootRequest) error {
	if err := u.validatePONPort(req.PONPort); err != nil {
		return err
	}
	if err := u.validateONUID(req.ONUID); err != nil {
		return err
	}
	return nil
}

func (u *ONUManagementUsecase) validateBlockRequest(req *model.ONUBlockRequest) error {
	if err := u.validatePONPort(req.PONPort); err != nil {
		return err
	}
	if err := u.validateONUID(req.ONUID); err != nil {
		return err
	}
	return nil
}

func (u *ONUManagementUsecase) validateDescriptionRequest(req *model.ONUDescriptionRequest) error {
	if err := u.validatePONPort(req.PONPort); err != nil {
		return err
	}
	if err := u.validateONUID(req.ONUID); err != nil {
		return err
	}

	// Validate description
	req.Description = strings.TrimSpace(req.Description)
	if req.Description == "" {
		return fmt.Errorf("description cannot be empty")
	}
	if len(req.Description) > 64 {
		return fmt.Errorf("description too long (max 64 characters)")
	}

	// Validate description format (alphanumeric, spaces, dashes, underscores)
	validDesc := regexp.MustCompile(`^[a-zA-Z0-9\s\-_\.]+$`)
	if !validDesc.MatchString(req.Description) {
		return fmt.Errorf("description contains invalid characters (only alphanumeric, spaces, dashes, underscores, and dots allowed)")
	}

	return nil
}

func (u *ONUManagementUsecase) validateDeleteRequest(req *model.ONUDeleteRequest) error {
	if err := u.validatePONPort(req.PONPort); err != nil {
		return err
	}
	if err := u.validateONUID(req.ONUID); err != nil {
		return err
	}
	return nil
}

func (u *ONUManagementUsecase) validatePONPort(ponPort string) error {
	ponPort = strings.TrimSpace(ponPort)
	if ponPort == "" {
		return fmt.Errorf("PON port cannot be empty")
	}

	// Validate PON port format: 1/1/1 to 1/16/16 (rack/shelf/port)
	validPON := regexp.MustCompile(`^(\d+)/(\d+)/(\d+)$`)
	matches := validPON.FindStringSubmatch(ponPort)
	if matches == nil {
		return fmt.Errorf("invalid PON port format, expected format: rack/shelf/port (e.g., 1/1/1)")
	}

	return nil
}

func (u *ONUManagementUsecase) validateONUID(onuID int) error {
	if onuID < 1 || onuID > 128 {
		return fmt.Errorf("ONU ID must be between 1 and 128")
	}
	return nil
}
