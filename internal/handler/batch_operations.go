package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/usecase"
	"github.com/s4lfanet/go-api-c320/internal/utils"
)

// BatchOperationsHandlerInterface defines HTTP handlers for batch ONU operations
type BatchOperationsHandlerInterface interface {
	BatchRebootONUs(w http.ResponseWriter, r *http.Request)
	BatchBlockONUs(w http.ResponseWriter, r *http.Request)
	BatchUnblockONUs(w http.ResponseWriter, r *http.Request)
	BatchDeleteONUs(w http.ResponseWriter, r *http.Request)
	BatchUpdateDescriptions(w http.ResponseWriter, r *http.Request)
}

// BatchOperationsHandler implements batch ONU operations HTTP handlers
type BatchOperationsHandler struct {
	batchUsecase usecase.BatchOperationsUsecaseInterface
}

// NewBatchOperationsHandler creates a new batch operations handler
func NewBatchOperationsHandler(batchUsecase usecase.BatchOperationsUsecaseInterface) BatchOperationsHandlerInterface {
	return &BatchOperationsHandler{
		batchUsecase: batchUsecase,
	}
}

// BatchRebootONUs godoc
// @Summary      Batch Reboot ONUs
// @Description  Reboot multiple ONUs in a single operation (max 50 ONUs)
// @Tags         Batch Operations
// @Accept       json
// @Produce      json
// @Param        request body model.BatchONURebootRequest true "Batch Reboot Request"
// @Success      200 {object} utils.WebResponse{data=model.BatchONURebootResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/batch/reboot [post]
func (h *BatchOperationsHandler) BatchRebootONUs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.BatchONURebootRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode batch reboot request")
		utils.HandleError(w, err)
		return
	}

	response, err := h.batchUsecase.BatchRebootONUs(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute batch reboot")
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

// BatchBlockONUs godoc
// @Summary      Batch Block ONUs
// @Description  Block (disable) multiple ONUs in a single operation (max 50 ONUs)
// @Tags         Batch Operations
// @Accept       json
// @Produce      json
// @Param        request body model.BatchONUBlockRequest true "Batch Block Request"
// @Success      200 {object} utils.WebResponse{data=model.BatchONUBlockResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/batch/block [post]
func (h *BatchOperationsHandler) BatchBlockONUs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.BatchONUBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode batch block request")
		utils.HandleError(w, err)
		return
	}

	req.Block = true

	response, err := h.batchUsecase.BatchBlockONUs(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute batch block")
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

// BatchUnblockONUs godoc
// @Summary      Batch Unblock ONUs
// @Description  Unblock (enable) multiple ONUs in a single operation (max 50 ONUs)
// @Tags         Batch Operations
// @Accept       json
// @Produce      json
// @Param        request body model.BatchONUBlockRequest true "Batch Unblock Request (without block field)"
// @Success      200 {object} utils.WebResponse{data=model.BatchONUBlockResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/batch/unblock [post]
func (h *BatchOperationsHandler) BatchUnblockONUs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.BatchONUBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode batch unblock request")
		utils.HandleError(w, err)
		return
	}

	req.Block = false

	response, err := h.batchUsecase.BatchUnblockONUs(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute batch unblock")
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

// BatchDeleteONUs godoc
// @Summary      Batch Delete ONUs
// @Description  Delete multiple ONU configurations in a single operation (max 50 ONUs)
// @Tags         Batch Operations
// @Accept       json
// @Produce      json
// @Param        request body model.BatchONUDeleteRequest true "Batch Delete Request"
// @Success      200 {object} utils.WebResponse{data=model.BatchONUDeleteResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/batch/delete [post]
func (h *BatchOperationsHandler) BatchDeleteONUs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.BatchONUDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode batch delete request")
		utils.HandleError(w, err)
		return
	}

	response, err := h.batchUsecase.BatchDeleteONUs(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute batch delete")
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

// BatchUpdateDescriptions godoc
// @Summary      Batch Update ONU Descriptions
// @Description  Update descriptions for multiple ONUs in a single operation (max 50 ONUs)
// @Tags         Batch Operations
// @Accept       json
// @Produce      json
// @Param        request body model.BatchONUDescriptionRequest true "Batch Description Update Request"
// @Success      200 {object} utils.WebResponse{data=model.BatchONUDescriptionResponse}
// @Failure      400 {object} utils.ErrorResponse
// @Failure      500 {object} utils.ErrorResponse
// @Router       /api/v1/batch/descriptions [put]
func (h *BatchOperationsHandler) BatchUpdateDescriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.BatchONUDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode batch description update request")
		utils.HandleError(w, err)
		return
	}

	response, err := h.batchUsecase.BatchUpdateDescriptions(ctx, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute batch description update")
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
