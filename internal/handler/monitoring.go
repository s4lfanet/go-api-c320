package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/internal/usecase"
	"github.com/s4lfanet/go-api-c320/internal/utils"
)

// MonitoringHandler handles monitoring-related HTTP requests
type MonitoringHandler struct {
	monitoringUsecase *usecase.MonitoringUsecase
}

// NewMonitoringHandler creates a new MonitoringHandler instance
func NewMonitoringHandler(monitoringUsecase *usecase.MonitoringUsecase) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringUsecase: monitoringUsecase,
	}
}

// GetONUMonitoring godoc
// @Summary Get real-time ONU monitoring data
// @Description Retrieves current monitoring status, statistics, and information for a specific ONU
// @Tags Monitoring
// @Accept json
// @Produce json
// @Param pon path int true "PON Port Number (1-16)"
// @Param onuId path int true "ONU ID (1-128)"
// @Success 200 {object} webresponse.WebResponse{data=model.ONUMonitoringInfo} "Successfully retrieved ONU monitoring data"
// @Failure 404 {object} webresponse.WebResponse "ONU not found or PON not configured"
// @Failure 500 {object} webresponse.WebResponse "Internal server error"
// @Router /api/v1/monitoring/onu/{pon}/{onuId} [get]
func (h *MonitoringHandler) GetONUMonitoring(w http.ResponseWriter, r *http.Request) {
	ponPort := chi.URLParam(r, "pon")
	onuIDStr := chi.URLParam(r, "onuId")
	onuID, _ := strconv.Atoi(onuIDStr)

	log.Info().Str("pon", ponPort).Int("onu_id", onuID).Msg("Getting ONU monitoring data")

	monitoring, err := h.monitoringUsecase.GetONUMonitoring(r.Context(), ponPort, onuID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get ONU monitoring")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   monitoring,
	})
}

// GetPONMonitoring godoc
// @Summary Get PON port monitoring data
// @Description Retrieves aggregated monitoring data for a PON port including all ONUs
// @Tags Monitoring
// @Accept json
// @Produce json
// @Param pon path int true "PON Port Number (1-16)"
// @Success 200 {object} webresponse.WebResponse{data=model.PONMonitoringInfo} "Successfully retrieved PON monitoring data"
// @Failure 404 {object} webresponse.WebResponse "PON port not configured"
// @Failure 500 {object} webresponse.WebResponse "Internal server error"
// @Router /api/v1/monitoring/pon/{pon} [get]
func (h *MonitoringHandler) GetPONMonitoring(w http.ResponseWriter, r *http.Request) {
	ponPort := chi.URLParam(r, "pon")

	log.Info().Str("pon", ponPort).Msg("Getting PON monitoring data")

	monitoring, err := h.monitoringUsecase.GetPONMonitoring(r.Context(), ponPort)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get PON monitoring")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   monitoring,
	})
}

// GetOLTMonitoring godoc
// @Summary Get OLT monitoring summary
// @Description Retrieves overall OLT monitoring summary including all PON ports and ONUs
// @Tags Monitoring
// @Accept json
// @Produce json
// @Success 200 {object} webresponse.WebResponse{data=model.OLTMonitoringSummary} "Successfully retrieved OLT monitoring summary"
// @Failure 500 {object} webresponse.WebResponse "Internal server error"
// @Router /api/v1/monitoring/olt [get]
func (h *MonitoringHandler) GetOLTMonitoring(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Getting OLT monitoring summary")

	monitoring, err := h.monitoringUsecase.GetOLTMonitoring(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get OLT monitoring")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   monitoring,
	})
}
