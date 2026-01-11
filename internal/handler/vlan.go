package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/usecase"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// VLANHandlerInterface defines the interface for VLAN handlers
type VLANHandlerInterface interface {
	GetONUVLAN(w http.ResponseWriter, r *http.Request)
	GetAllServicePorts(w http.ResponseWriter, r *http.Request)
	ConfigureVLAN(w http.ResponseWriter, r *http.Request)
	ModifyVLAN(w http.ResponseWriter, r *http.Request)
	DeleteVLAN(w http.ResponseWriter, r *http.Request)
}

// VLANHandler implements the VLAN handler interface
type VLANHandler struct {
	vlanUsecase usecase.VLANUsecaseInterface
}

// NewVLANHandler creates a new VLAN handler instance
func NewVLANHandler(vlanUsecase usecase.VLANUsecaseInterface) VLANHandlerInterface {
	return &VLANHandler{
		vlanUsecase: vlanUsecase,
	}
}

// GetONUVLAN retrieves VLAN configuration for a specific ONU
// @Summary Get ONU VLAN configuration
// @Description Retrieve VLAN configuration details for a specific ONU
// @Tags VLAN
// @Accept json
// @Produce json
// @Param pon path string true "PON Port (e.g., 1/1/1)"
// @Param onu_id path int true "ONU ID"
// @Success 200 {object} utils.WebResponse{data=model.ONUVLANInfo}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/vlan/onu/{pon}/{onu_id} [get]
func (h *VLANHandler) GetONUVLAN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ponPort := chi.URLParam(r, "pon")
	onuIDStr := chi.URLParam(r, "onu_id")

	if ponPort == "" || onuIDStr == "" {
		appErr := apperrors.NewValidationError("PON port and ONU ID are required", nil)
		utils.HandleError(w, appErr)
		return
	}

	onuID, err := strconv.Atoi(onuIDStr)
	if err != nil {
		appErr := apperrors.NewValidationError("Invalid ONU ID", map[string]interface{}{"onu_id": onuIDStr})
		utils.HandleError(w, appErr)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("Getting ONU VLAN configuration")

	vlanInfo, err := h.vlanUsecase.GetONUVLAN(ctx, ponPort, onuID)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Msg("Failed to get ONU VLAN")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("svlan", vlanInfo.SVLAN).
		Msg("Retrieved ONU VLAN configuration")

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   vlanInfo,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// GetAllServicePorts retrieves all service-port configurations
// @Summary Get all service-port configurations
// @Description Retrieve all service-port (VLAN) configurations from the OLT
// @Tags VLAN
// @Accept json
// @Produce json
// @Success 200 {object} utils.WebResponse{data=[]model.ONUVLANInfo}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/vlan/service-ports [get]
func (h *VLANHandler) GetAllServicePorts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Info().Msg("Getting all service-port configurations")

	servicePorts, err := h.vlanUsecase.GetAllServicePorts(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get service-ports")
		utils.HandleError(w, err)
		return
	}

	log.Info().Int("count", len(servicePorts)).Msg("Retrieved service-port configurations")

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   servicePorts,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// ConfigureVLAN configures VLAN for an ONU
// @Summary Configure ONU VLAN
// @Description Configure VLAN (service-port) for an ONU
// @Tags VLAN
// @Accept json
// @Produce json
// @Param request body model.VLANConfigRequest true "VLAN Configuration Request"
// @Success 201 {object} utils.WebResponse{data=model.VLANConfigResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/vlan/onu [post]
func (h *VLANHandler) ConfigureVLAN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.VLANConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		appErr := apperrors.NewValidationError("Invalid request body", map[string]interface{}{"error": err.Error()})
		utils.HandleError(w, appErr)
		return
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("svlan", req.SVLAN).
		Int("cvlan", req.CVLAN).
		Str("vlan_mode", req.VLANMode).
		Msg("Configuring ONU VLAN")

	response, err := h.vlanUsecase.ConfigureVLAN(ctx, req)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Msg("Failed to configure VLAN")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("service_port_id", response.ServicePortID).
		Msg("VLAN configured successfully")

	webResponse := utils.WebResponse{
		Code:   http.StatusCreated,
		Status: "Created",
		Data:   response,
	}
	utils.SendJSONResponse(w, http.StatusCreated, webResponse)
}

// ModifyVLAN modifies existing VLAN configuration for an ONU
// @Summary Modify ONU VLAN
// @Description Modify existing VLAN (service-port) configuration for an ONU
// @Tags VLAN
// @Accept json
// @Produce json
// @Param request body model.VLANConfigRequest true "VLAN Configuration Request"
// @Success 200 {object} utils.WebResponse{data=model.VLANConfigResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/vlan/onu [put]
func (h *VLANHandler) ModifyVLAN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.VLANConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		appErr := apperrors.NewValidationError("Invalid request body", map[string]interface{}{"error": err.Error()})
		utils.HandleError(w, appErr)
		return
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("svlan", req.SVLAN).
		Msg("Modifying ONU VLAN")

	response, err := h.vlanUsecase.ModifyVLAN(ctx, req)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Msg("Failed to modify VLAN")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("service_port_id", response.ServicePortID).
		Msg("VLAN modified successfully")

	webResponse := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.SendJSONResponse(w, http.StatusOK, webResponse)
}

// DeleteVLAN removes VLAN configuration for an ONU
// @Summary Delete ONU VLAN
// @Description Remove VLAN (service-port) configuration for an ONU
// @Tags VLAN
// @Accept json
// @Produce json
// @Param pon path string true "PON Port (e.g., 1/1/1)"
// @Param onu_id path int true "ONU ID"
// @Success 200 {object} utils.WebResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/vlan/onu/{pon}/{onu_id} [delete]
func (h *VLANHandler) DeleteVLAN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ponPort := chi.URLParam(r, "pon")
	onuIDStr := chi.URLParam(r, "onu_id")

	if ponPort == "" || onuIDStr == "" {
		appErr := apperrors.NewValidationError("PON port and ONU ID are required", nil)
		utils.HandleError(w, appErr)
		return
	}

	onuID, err := strconv.Atoi(onuIDStr)
	if err != nil {
		appErr := apperrors.NewValidationError("Invalid ONU ID", map[string]interface{}{"onu_id": onuIDStr})
		utils.HandleError(w, appErr)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("Deleting ONU VLAN")

	err = h.vlanUsecase.DeleteVLAN(ctx, ponPort, onuID)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Msg("Failed to delete VLAN")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("VLAN deleted successfully")

	responseData := map[string]interface{}{
		"pon_port":   ponPort,
		"onu_id":     onuID,
		"message":    "VLAN deleted successfully",
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   responseData,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}
