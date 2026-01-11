package handler

import (
	"net/http"
	"strconv"

	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/usecase"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// CardHandlerInterface defines the interface for card/slot handlers
type CardHandlerInterface interface {
	GetAllCards(w http.ResponseWriter, r *http.Request)
	GetCard(w http.ResponseWriter, r *http.Request)
}

// CardHandler handles card/slot related requests
type CardHandler struct {
	cardUsecase usecase.CardUseCaseInterface
}

// NewCardHandler creates a new card handler instance
func NewCardHandler(cardUsecase usecase.CardUseCaseInterface) *CardHandler {
	return &CardHandler{cardUsecase: cardUsecase}
}

// GetAllCards retrieves all card/slot information
// Example: GET /api/v1/system/cards
func (h *CardHandler) GetAllCards(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Getting all card information")

	// Call usecase to get all cards
	cards, err := h.cardUsecase.GetAllCards(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get card information")
		utils.HandleError(w, err)
		return
	}

	log.Info().Int("count", len(cards)).Msg("Successfully retrieved card information")

	// Create web response object
	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   cards,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

// GetCard retrieves specific card information
// Example: GET /api/v1/system/cards/1/1/1
func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	// Get rack, shelf, slot from URL parameters
	rackStr := chi.URLParam(r, "rack")
	shelfStr := chi.URLParam(r, "shelf")
	slotStr := chi.URLParam(r, "slot")

	rack, err := strconv.Atoi(rackStr)
	if err != nil {
		log.Warn().Str("rack", rackStr).Msg("Invalid rack number")
		utils.HandleError(w, err)
		return
	}

	shelf, err := strconv.Atoi(shelfStr)
	if err != nil {
		log.Warn().Str("shelf", shelfStr).Msg("Invalid shelf number")
		utils.HandleError(w, err)
		return
	}

	slot, err := strconv.Atoi(slotStr)
	if err != nil {
		log.Warn().Str("slot", slotStr).Msg("Invalid slot number")
		utils.HandleError(w, err)
		return
	}

	log.Info().Int("rack", rack).Int("shelf", shelf).Int("slot", slot).Msg("Getting card information")

	// Call usecase to get card information
	card, err := h.cardUsecase.GetCard(r.Context(), rack, shelf, slot)
	if err != nil {
		log.Error().
			Err(err).
			Int("rack", rack).
			Int("shelf", shelf).
			Int("slot", slot).
			Msg("Failed to get card information")
		utils.HandleError(w, err)
		return
	}

	log.Info().
		Int("rack", rack).
		Int("shelf", shelf).
		Int("slot", slot).
		Str("card_type", card.CardType).
		Msg("Successfully retrieved card information")

	// Create web response object
	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   card,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}
