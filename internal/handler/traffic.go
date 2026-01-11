package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/usecase"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// TrafficHandlerInterface defines the interface for traffic profile handlers
type TrafficHandlerInterface interface {
	// DBA Profile handlers
	GetDBAProfile(w http.ResponseWriter, r *http.Request)
	GetAllDBAProfiles(w http.ResponseWriter, r *http.Request)
	CreateDBAProfile(w http.ResponseWriter, r *http.Request)
	ModifyDBAProfile(w http.ResponseWriter, r *http.Request)
	DeleteDBAProfile(w http.ResponseWriter, r *http.Request)

	// TCONT handlers
	GetONUTCONT(w http.ResponseWriter, r *http.Request)
	ConfigureTCONT(w http.ResponseWriter, r *http.Request)
	DeleteTCONT(w http.ResponseWriter, r *http.Request)

	// GEMPort handlers
	ConfigureGEMPort(w http.ResponseWriter, r *http.Request)
	DeleteGEMPort(w http.ResponseWriter, r *http.Request)
}

// TrafficHandler handles traffic profile related requests
type TrafficHandler struct {
	trafficUsecase usecase.TrafficUsecaseInterface
}

// NewTrafficHandler creates a new traffic handler instance
func NewTrafficHandler(trafficUsecase usecase.TrafficUsecaseInterface) *TrafficHandler {
	return &TrafficHandler{
		trafficUsecase: trafficUsecase,
	}
}

// GetDBAProfile retrieves DBA profile information
// @Summary Get DBA profile
// @Description Retrieve DBA (Dynamic Bandwidth Allocation) profile details by name
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param name path string true "DBA Profile Name"
// @Success 200 {object} utils.WebResponse{data=model.DBAProfileInfo}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/dba-profile/{name} [get]
func (h *TrafficHandler) GetDBAProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profileName := chi.URLParam(r, "name")

	log.Info().Str("profile_name", profileName).Msg("Getting DBA profile")

	profile, err := h.trafficUsecase.GetDBAProfile(ctx, profileName)
	if err != nil {
		log.Error().Err(err).Str("profile_name", profileName).Msg("Failed to get DBA profile")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   profile,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// GetAllDBAProfiles retrieves all DBA profiles
// @Summary Get all DBA profiles
// @Description Retrieve list of all DBA profiles configured on the OLT
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Success 200 {object} utils.WebResponse{data=[]model.DBAProfileInfo}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/dba-profiles [get]
func (h *TrafficHandler) GetAllDBAProfiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Info().Msg("Getting all DBA profiles")

	profiles, err := h.trafficUsecase.GetAllDBAProfiles(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get DBA profiles")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   profiles,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// CreateDBAProfile creates a new DBA profile
// @Summary Create DBA profile
// @Description Create a new DBA (Dynamic Bandwidth Allocation) profile
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param request body model.DBAProfileRequest true "DBA Profile Configuration"
// @Success 200 {object} utils.WebResponse{data=model.DBAProfileResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/dba-profile [post]
func (h *TrafficHandler) CreateDBAProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.DBAProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("profile_name", req.Name).
		Int("type", req.Type).
		Msg("Creating DBA profile")

	result, err := h.trafficUsecase.CreateDBAProfile(ctx, req)
	if err != nil {
		log.Error().Err(err).Str("profile_name", req.Name).Msg("Failed to create DBA profile")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// ModifyDBAProfile modifies an existing DBA profile
// @Summary Modify DBA profile
// @Description Modify an existing DBA profile configuration
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param request body model.DBAProfileRequest true "DBA Profile Configuration"
// @Success 200 {object} utils.WebResponse{data=model.DBAProfileResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/dba-profile [put]
func (h *TrafficHandler) ModifyDBAProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.DBAProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("profile_name", req.Name).
		Int("type", req.Type).
		Msg("Modifying DBA profile")

	result, err := h.trafficUsecase.ModifyDBAProfile(ctx, req)
	if err != nil {
		log.Error().Err(err).Str("profile_name", req.Name).Msg("Failed to modify DBA profile")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// DeleteDBAProfile deletes a DBA profile
// @Summary Delete DBA profile
// @Description Delete a DBA profile by name
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param name path string true "DBA Profile Name"
// @Success 200 {object} utils.WebResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/dba-profile/{name} [delete]
func (h *TrafficHandler) DeleteDBAProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profileName := chi.URLParam(r, "name")

	log.Info().Str("profile_name", profileName).Msg("Deleting DBA profile")

	err := h.trafficUsecase.DeleteDBAProfile(ctx, profileName)
	if err != nil {
		log.Error().Err(err).Str("profile_name", profileName).Msg("Failed to delete DBA profile")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   map[string]string{"message": "DBA profile deleted successfully"},
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// GetONUTCONT retrieves T-CONT configuration for an ONU
// @Summary Get ONU T-CONT
// @Description Retrieve T-CONT (Transmission Container) configuration for an ONU
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param pon path string true "PON Port (e.g., 1/1/1)"
// @Param onu_id path int true "ONU ID"
// @Param tcont_id path int true "T-CONT ID"
// @Success 200 {object} utils.WebResponse{data=model.TCONTInfo}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/tcont/{pon}/{onu_id}/{tcont_id} [get]
func (h *TrafficHandler) GetONUTCONT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ponPort := chi.URLParam(r, "pon")
	onuIDStr := chi.URLParam(r, "onu_id")
	tcontIDStr := chi.URLParam(r, "tcont_id")

	onuID, err := strconv.Atoi(onuIDStr)
	if err != nil {
		log.Error().Err(err).Str("onu_id", onuIDStr).Msg("Invalid ONU ID")
		utils.HandleError(w, err)
		return
	}

	tcontID, err := strconv.Atoi(tcontIDStr)
	if err != nil {
		log.Error().Err(err).Str("tcont_id", tcontIDStr).Msg("Invalid T-CONT ID")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("tcont_id", tcontID).
		Msg("Getting T-CONT configuration")

	tcont, err := h.trafficUsecase.GetONUTCONT(ctx, ponPort, onuID, tcontID)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Int("tcont_id", tcontID).
			Msg("Failed to get T-CONT")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   tcont,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// ConfigureTCONT configures T-CONT for an ONU
// @Summary Configure T-CONT
// @Description Configure T-CONT (Transmission Container) for an ONU
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param request body model.TCONTConfigRequest true "T-CONT Configuration"
// @Success 200 {object} utils.WebResponse{data=model.TCONTConfigResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/tcont [post]
func (h *TrafficHandler) ConfigureTCONT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.TCONTConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("tcont_id", req.TCONTID).
		Str("profile", req.Profile).
		Msg("Configuring T-CONT")

	result, err := h.trafficUsecase.ConfigureTCONT(ctx, req)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Msg("Failed to configure T-CONT")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// DeleteTCONT deletes T-CONT from an ONU
// @Summary Delete T-CONT
// @Description Delete T-CONT configuration from an ONU
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param pon path string true "PON Port (e.g., 1/1/1)"
// @Param onu_id path int true "ONU ID"
// @Param tcont_id path int true "T-CONT ID"
// @Success 200 {object} utils.WebResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/tcont/{pon}/{onu_id}/{tcont_id} [delete]
func (h *TrafficHandler) DeleteTCONT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ponPort := chi.URLParam(r, "pon")
	onuIDStr := chi.URLParam(r, "onu_id")
	tcontIDStr := chi.URLParam(r, "tcont_id")

	onuID, err := strconv.Atoi(onuIDStr)
	if err != nil {
		log.Error().Err(err).Str("onu_id", onuIDStr).Msg("Invalid ONU ID")
		utils.HandleError(w, err)
		return
	}

	tcontID, err := strconv.Atoi(tcontIDStr)
	if err != nil {
		log.Error().Err(err).Str("tcont_id", tcontIDStr).Msg("Invalid T-CONT ID")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("tcont_id", tcontID).
		Msg("Deleting T-CONT")

	err = h.trafficUsecase.DeleteTCONT(ctx, ponPort, onuID, tcontID)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Int("tcont_id", tcontID).
			Msg("Failed to delete T-CONT")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   map[string]string{"message": "T-CONT deleted successfully"},
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// ConfigureGEMPort configures GEM port for an ONU
// @Summary Configure GEM port
// @Description Configure GEM (GPON Encapsulation Method) port for an ONU
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param request body model.GEMPortConfigRequest true "GEM Port Configuration"
// @Success 200 {object} utils.WebResponse{data=model.GEMPortConfigResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/gemport [post]
func (h *TrafficHandler) ConfigureGEMPort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.GEMPortConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", req.PONPort).
		Int("onu_id", req.ONUID).
		Int("gemport_id", req.GEMPortID).
		Int("tcont_id", req.TCONTID).
		Msg("Configuring GEM port")

	result, err := h.trafficUsecase.ConfigureGEMPort(ctx, req)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", req.PONPort).
			Int("onu_id", req.ONUID).
			Msg("Failed to configure GEM port")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// DeleteGEMPort deletes GEM port from an ONU
// @Summary Delete GEM port
// @Description Delete GEM port configuration from an ONU
// @Tags Traffic Management
// @Accept json
// @Produce json
// @Param pon path string true "PON Port (e.g., 1/1/1)"
// @Param onu_id path int true "ONU ID"
// @Param gemport_id path int true "GEM Port ID"
// @Success 200 {object} utils.WebResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/traffic/gemport/{pon}/{onu_id}/{gemport_id} [delete]
func (h *TrafficHandler) DeleteGEMPort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ponPort := chi.URLParam(r, "pon")
	onuIDStr := chi.URLParam(r, "onu_id")
	gemportIDStr := chi.URLParam(r, "gemport_id")

	onuID, err := strconv.Atoi(onuIDStr)
	if err != nil {
		log.Error().Err(err).Str("onu_id", onuIDStr).Msg("Invalid ONU ID")
		utils.HandleError(w, err)
		return
	}

	gemportID, err := strconv.Atoi(gemportIDStr)
	if err != nil {
		log.Error().Err(err).Str("gemport_id", gemportIDStr).Msg("Invalid GEM port ID")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Int("gemport_id", gemportID).
		Msg("Deleting GEM port")

	err = h.trafficUsecase.DeleteGEMPort(ctx, ponPort, onuID, gemportID)
	if err != nil {
		log.Error().Err(err).
			Str("pon_port", ponPort).
			Int("onu_id", onuID).
			Int("gemport_id", gemportID).
			Msg("Failed to delete GEM port")
		utils.HandleError(w, err)
		return
	}

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   map[string]string{"message": "GEM port deleted successfully"},
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}
