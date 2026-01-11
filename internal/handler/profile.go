package handler

import (
	"net/http"
	"strconv"

	"github.com/s4lfanet/go-api-c320/internal/usecase"
	"github.com/s4lfanet/go-api-c320/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// ProfileHandlerInterface defines the interface for traffic profile handlers
type ProfileHandlerInterface interface {
	GetAllTrafficProfiles(w http.ResponseWriter, r *http.Request)
	GetTrafficProfile(w http.ResponseWriter, r *http.Request)
	GetAllVlanProfiles(w http.ResponseWriter, r *http.Request)
}

// ProfileHandler handles traffic profile related requests
type ProfileHandler struct {
	profileUsecase usecase.ProfileUseCaseInterface
}

// NewProfileHandler creates a new profile handler instance
func NewProfileHandler(profileUsecase usecase.ProfileUseCaseInterface) *ProfileHandler {
	return &ProfileHandler{profileUsecase: profileUsecase}
}

// GetAllTrafficProfiles retrieves all traffic profiles
// Example: GET /api/v1/profiles/traffic
func (h *ProfileHandler) GetAllTrafficProfiles(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Getting all traffic profiles")

	// Call usecase to get all traffic profiles
	profiles, err := h.profileUsecase.GetAllTrafficProfiles(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get traffic profiles")
		utils.HandleError(w, err)
		return
	}

	log.Info().Int("count", len(profiles)).Msg("Successfully retrieved traffic profiles")

	// Create web response object
	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   profiles,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// GetTrafficProfile retrieves a specific traffic profile by ID
// Example: GET /api/v1/profiles/traffic/1
func (h *ProfileHandler) GetTrafficProfile(w http.ResponseWriter, r *http.Request) {
	// Get profile_id from URL parameter
	profileIDStr := chi.URLParam(r, "profile_id")
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		log.Warn().Str("profile_id", profileIDStr).Msg("Invalid profile ID")
		utils.HandleError(w, err)
		return
	}

	log.Info().Int("profile_id", profileID).Msg("Getting traffic profile")

	// Call usecase to get traffic profile
	profile, err := h.profileUsecase.GetTrafficProfile(r.Context(), profileID)
	if err != nil {
		log.Error().Err(err).Int("profile_id", profileID).Msg("Failed to get traffic profile")
		utils.HandleError(w, err)
		return
	}

	log.Info().Int("profile_id", profileID).Str("name", profile.Name).Msg("Successfully retrieved traffic profile")

	// Create web response object
	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   profile,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// GetAllVlanProfiles retrieves all VLAN profiles
// Example: GET /api/v1/profiles/vlan
func (h *ProfileHandler) GetAllVlanProfiles(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Getting all VLAN profiles")

	// Call usecase to get all VLAN profiles
	profiles, err := h.profileUsecase.GetAllVlanProfiles(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get VLAN profiles")
		utils.HandleError(w, err)
		return
	}

	log.Info().Int("count", len(profiles)).Msg("Successfully retrieved VLAN profiles")

	// Create web response object
	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   profiles,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

