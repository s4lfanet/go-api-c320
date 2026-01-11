package handler

import (
	"net/http"

	"github.com/s4lfanet/go-api-c320/internal/middleware"
	"github.com/s4lfanet/go-api-c320/internal/usecase"
	"github.com/s4lfanet/go-api-c320/internal/utils"
	"github.com/rs/zerolog/log"
)

// PonHandlerInterface defines the interface for PON port handlers
type PonHandlerInterface interface {
	GetPonPortInfo(w http.ResponseWriter, r *http.Request)
}

// PonHandler handles PON port related requests
type PonHandler struct {
	ponUsecase usecase.PonUseCaseInterface
}

// NewPonHandler creates a new PON handler instance
func NewPonHandler(ponUsecase usecase.PonUseCaseInterface) *PonHandler {
	return &PonHandler{ponUsecase: ponUsecase}
}

// GetPonPortInfo retrieves PON port information
// Example: GET /api/v1/pon/board/1/pon/1/info
func (h *PonHandler) GetPonPortInfo(w http.ResponseWriter, r *http.Request) {
	// Get pre-validated values from context
	boardIDInt, _ := middleware.GetBoardID(r.Context())
	ponIDInt, _ := middleware.GetPonID(r.Context())

	log.Info().
		Int("board_id", boardIDInt).
		Int("pon_id", ponIDInt).
		Msg("Getting PON port info")

	// Call usecase to get PON port info
	ponInfo, err := h.ponUsecase.GetPonPortInfo(r.Context(), boardIDInt, ponIDInt)
	if err != nil {
		log.Error().
			Err(err).
			Int("board_id", boardIDInt).
			Int("pon_id", ponIDInt).
			Msg("Failed to get PON port info")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Int("board_id", boardIDInt).
		Int("pon_id", ponIDInt).
		Msg("Successfully retrieved PON port info")

	// Create web response object
	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   ponInfo,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}
