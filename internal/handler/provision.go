package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	apperrors "github.com/s4lfanet/go-api-c320/internal/errors"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/usecase"
	"github.com/s4lfanet/go-api-c320/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// ProvisionHandlerInterface defines the interface for provisioning handlers
type ProvisionHandlerInterface interface {
	GetUnconfiguredONUs(w http.ResponseWriter, r *http.Request)
	GetUnconfiguredONUsByPON(w http.ResponseWriter, r *http.Request)
	RegisterONU(w http.ResponseWriter, r *http.Request)
	DeleteONU(w http.ResponseWriter, r *http.Request)
}

// ProvisionHandler handles ONU provisioning related requests
type ProvisionHandler struct {
	provisionUsecase usecase.ProvisionUseCaseInterface
}

// NewProvisionHandler creates a new provision handler instance
func NewProvisionHandler(provisionUsecase usecase.ProvisionUseCaseInterface) *ProvisionHandler {
	return &ProvisionHandler{
		provisionUsecase: provisionUsecase,
	}
}

// GetUnconfiguredONUs retrieves all unconfigured ONUs
// @Summary Get all unconfigured ONUs
// @Description Retrieve list of all unconfigured ONUs discovered on all PON ports
// @Tags Provisioning
// @Accept json
// @Produce json
// @Success 200 {object} utils.WebResponse{data=[]model.UnconfiguredONU}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/onu/unconfigured [get]
func (h *ProvisionHandler) GetUnconfiguredONUs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Info().Msg("Getting all unconfigured ONUs")

	onus, err := h.provisionUsecase.GetAllUnconfiguredONUs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get unconfigured ONUs")
		utils.HandleError(w, err)
		return
	}

	log.Info().Int("count", len(onus)).Msg("Retrieved unconfigured ONUs")

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   onus,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// GetUnconfiguredONUsByPON retrieves unconfigured ONUs for a specific PON port
// @Summary Get unconfigured ONUs by PON port
// @Description Retrieve list of unconfigured ONUs discovered on a specific PON port
// @Tags Provisioning
// @Accept json
// @Produce json
// @Param pon path string true "PON Port (e.g., 1/1/1)"
// @Success 200 {object} utils.WebResponse{data=[]model.UnconfiguredONU}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/onu/unconfigured/{pon} [get]
func (h *ProvisionHandler) GetUnconfiguredONUsByPON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ponPort := chi.URLParam(r, "pon")

	if ponPort == "" {
		appErr := apperrors.NewValidationError("PON port is required", nil)
		utils.HandleError(w, appErr)
		return
	}

	log.Info().Str("pon_port", ponPort).Msg("Getting unconfigured ONUs for PON port")

	onus, err := h.provisionUsecase.GetUnconfiguredONUs(ctx, ponPort)
	if err != nil {
		log.Error().Err(err).Str("pon_port", ponPort).Msg("Failed to get unconfigured ONUs")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("count", len(onus)).
		Msg("Retrieved unconfigured ONUs")

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   onus,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// RegisterONU registers a new ONU to the OLT
// @Summary Register a new ONU
// @Description Register a new ONU with TCONT, GEMPORT, and service port configuration
// @Tags Provisioning
// @Accept json
// @Produce json
// @Param request body model.ONURegistrationRequest true "ONU Registration Request"
// @Success 201 {object} utils.WebResponse{data=model.ONURegistrationResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/onu/register [post]
func (h *ProvisionHandler) RegisterONU(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.ONURegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		appErr := apperrors.NewValidationError("Invalid request body", map[string]interface{}{"error": err.Error()})
		utils.HandleError(w, appErr)
		return
	}

	// Validate request
	if err := h.validateRegistrationRequest(&req); err != nil {
		log.Error().Err(err).Msg("Invalid registration request")
		appErr := apperrors.NewValidationError(err.Error(), nil)
		utils.HandleError(w, appErr)
		return
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Str("serial", req.SerialNumber).
		Str("type", req.ONUType).
		Msg("Registering ONU")

	resp, err := h.provisionUsecase.RegisterONU(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register ONU")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Msg("ONU registered successfully")

	response := utils.WebResponse{
		Code:   http.StatusCreated,
		Status: "Created",
		Data:   resp,
	}
	utils.SendJSONResponse(w, http.StatusCreated, response)
}

// DeleteONU deletes an ONU from the OLT
// @Summary Delete an ONU
// @Description Delete an ONU from the OLT by PON port and ONU ID
// @Tags Provisioning
// @Accept json
// @Produce json
// @Param pon path string true "PON Port (e.g., 1/1/1)"
// @Param onu_id path int true "ONU ID"
// @Success 200 {object} utils.WebResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/onu/{pon}/{onu_id} [delete]
func (h *ProvisionHandler) DeleteONU(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ponPort := chi.URLParam(r, "pon")
	onuIDStr := chi.URLParam(r, "onu_id")

	if ponPort == "" || onuIDStr == "" {
		appErr := apperrors.NewValidationError("PON port and ONU ID are required", nil)
		utils.HandleError(w, appErr)
		return
	}

	// Parse ONU ID
	onuID, err := strconv.Atoi(onuIDStr)
	if err != nil {
		appErr := apperrors.NewValidationError("Invalid ONU ID", map[string]interface{}{"onu_id": onuIDStr})
		utils.HandleError(w, appErr)
		return
	}

	// Validate ONU ID range (1-128 for most OLTs)
	if onuID < 1 || onuID > 128 {
		appErr := apperrors.NewValidationError("ONU ID must be between 1 and 128", nil)
		utils.HandleError(w, appErr)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("Deleting ONU")

	err = h.provisionUsecase.DeleteONU(ctx, ponPort, onuID)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Msg("Failed to delete ONU")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("ONU deleted successfully")

	responseData := map[string]interface{}{
		"pon_port":   ponPort,
		"onu_id":     onuID,
		"message":    "ONU deleted successfully",
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   responseData,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// validateRegistrationRequest validates the ONU registration request
func (h *ProvisionHandler) validateRegistrationRequest(req *model.ONURegistrationRequest) error {
	if req.PONPort == "" {
		return fmt.Errorf("pon_port is required")
	}

	if req.ONUID < 1 || req.ONUID > 128 {
		return fmt.Errorf("onu_id must be between 1 and 128")
	}

	if req.ONUType == "" {
		return fmt.Errorf("onu_type is required")
	}

	if req.SerialNumber == "" {
		return fmt.Errorf("serial_number is required")
	}

	// Validate serial number format (typically 12 characters: 4 vendor + 8 hex)
	if len(req.SerialNumber) < 8 || len(req.SerialNumber) > 16 {
		return fmt.Errorf("invalid serial_number format")
	}

	// Validate profile if provided
	if req.Profile.DBAProfile != "" {
		if req.Profile.VLAN < 1 || req.Profile.VLAN > 4094 {
			return fmt.Errorf("vlan must be between 1 and 4094")
		}
	}

	return nil
}
