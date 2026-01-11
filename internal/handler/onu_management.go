package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/usecase"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// ONUManagementHandlerInterface defines HTTP handlers for ONU lifecycle management
type ONUManagementHandlerInterface interface {
	RebootONU(w http.ResponseWriter, r *http.Request)
	BlockONU(w http.ResponseWriter, r *http.Request)
	UnblockONU(w http.ResponseWriter, r *http.Request)
	UpdateDescription(w http.ResponseWriter, r *http.Request)
	DeleteONU(w http.ResponseWriter, r *http.Request)
}

// ONUManagementHandler implements ONU management HTTP handlers
type ONUManagementHandler struct {
	onuMgmtUsecase usecase.ONUManagementUsecaseInterface
}

// NewONUManagementHandler creates a new ONU management handler
func NewONUManagementHandler(onuMgmtUsecase usecase.ONUManagementUsecaseInterface) ONUManagementHandlerInterface {
	return &ONUManagementHandler{
		onuMgmtUsecase: onuMgmtUsecase,
	}
}

// RebootONU godoc
// @Summary      Reboot ONU
// @Description  Reboot/reset an ONU on the specified PON port
// @Tags         ONU Management
// @Accept       json
// @Produce      json
// @Param        request body model.ONURebootRequest true "ONU Reboot Request"
// @Success      200 {object} utils.WebResponse{data=model.ONURebootResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/onu-management/reboot [post]
func (h *ONUManagementHandler) RebootONU(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.ONURebootRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode reboot request")
		appErr := apperrors.NewValidationError("Invalid request body", nil)
		utils.HandleError(w, appErr)
		return
	}

	response, err := h.onuMgmtUsecase.RebootONU(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reboot ONU")
		utils.HandleError(w, err)
		return
	}

	webResp := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.SendJSONResponse(w, http.StatusOK, webResp)
}

// BlockONU godoc
// @Summary      Block ONU
// @Description  Block (disable) an ONU on the specified PON port
// @Tags         ONU Management
// @Accept       json
// @Produce      json
// @Param        request body model.ONUBlockRequest true "ONU Block Request"
// @Success      200 {object} utils.WebResponse{data=model.ONUBlockResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/onu-management/block [post]
func (h *ONUManagementHandler) BlockONU(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.ONUBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode block request")
		appErr := apperrors.NewValidationError("Invalid request body", nil)
		utils.HandleError(w, appErr)
		return
	}

	req.Block = true

	response, err := h.onuMgmtUsecase.BlockONU(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to block ONU")
		utils.HandleError(w, err)
		return
	}

	webResp := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.SendJSONResponse(w, http.StatusOK, webResp)
}

// UnblockONU godoc
// @Summary      Unblock ONU
// @Description  Unblock (enable) an ONU on the specified PON port
// @Tags         ONU Management
// @Accept       json
// @Produce      json
// @Param        request body model.ONUBlockRequest true "ONU Unblock Request"
// @Success      200 {object} utils.WebResponse{data=model.ONUBlockResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/onu-management/unblock [post]
func (h *ONUManagementHandler) UnblockONU(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.ONUBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode unblock request")
		appErr := apperrors.NewValidationError("Invalid request body", nil)
		utils.HandleError(w, appErr)
		return
	}

	req.Block = false

	response, err := h.onuMgmtUsecase.UnblockONU(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unblock ONU")
		utils.HandleError(w, err)
		return
	}

	webResp := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.SendJSONResponse(w, http.StatusOK, webResp)
}

// UpdateDescription godoc
// @Summary      Update ONU Description
// @Description  Update the name/description of an ONU
// @Tags         ONU Management
// @Accept       json
// @Produce      json
// @Param        request body model.ONUDescriptionRequest true "ONU Description Update Request"
// @Success      200 {object} utils.WebResponse{data=model.ONUDescriptionResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/onu-management/description [put]
func (h *ONUManagementHandler) UpdateDescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.ONUDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode description update request")
		appErr := apperrors.NewValidationError("Invalid request body", nil)
		utils.HandleError(w, appErr)
		return
	}

	response, err := h.onuMgmtUsecase.UpdateDescription(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update ONU description")
		utils.HandleError(w, err)
		return
	}

	webResp := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.SendJSONResponse(w, http.StatusOK, webResp)
}

// DeleteONU godoc
// @Summary      Delete ONU
// @Description  Delete ONU configuration from the OLT (removes all settings)
// @Tags         ONU Management
// @Accept       json
// @Produce      json
// @Param        pon path string true "PON Port (e.g., 1-1-1)"
// @Param        onu_id path int true "ONU ID (1-128)"
// @Success      200 {object} utils.WebResponse{data=model.ONUDeleteResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/onu-management/{pon}/{onu_id} [delete]
func (h *ONUManagementHandler) DeleteONU(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ponPort := chi.URLParam(r, "pon")
	if ponPort == "" {
		appErr := apperrors.NewValidationError("PON port is required", nil)
		utils.HandleError(w, appErr)
		return
	}

	onuIDStr := chi.URLParam(r, "onu_id")
	if onuIDStr == "" {
		appErr := apperrors.NewValidationError("ONU ID is required", nil)
		utils.HandleError(w, appErr)
		return
	}

	onuID, err := strconv.Atoi(onuIDStr)
	if err != nil {
		appErr := apperrors.NewValidationError("Invalid ONU ID format", map[string]interface{}{"onu_id": onuIDStr})
		utils.HandleError(w, appErr)
		return
	}

	req := &model.ONUDeleteRequest{
		PONPort: ponPort,
		ONUID:   onuID,
	}

	response, err := h.onuMgmtUsecase.DeleteONU(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete ONU")
		utils.HandleError(w, err)
		return
	}

	webResp := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.SendJSONResponse(w, http.StatusOK, webResp)
}
