package middleware

import (
	"context"
	"net/http"
	"strconv"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/utils"
	"github.com/go-chi/chi/v5"
)

// ValidateBoardPonParams validates board_id and pon_id URL parameters,
// ensuring they are valid integers within the expected range.
func ValidateBoardPonParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boardID := chi.URLParam(r, "board_id") // Get board_id from URL
		ponID := chi.URLParam(r, "pon_id")     // Get pon_id from URL

		// Validate board_id conversion to integer
		boardIDInt, err := strconv.Atoi(boardID)
		// Check if conversion failed or if board_id not 1 or 2
		if err != nil || (boardIDInt != 1 && boardIDInt != 2) {
			appErr := apperrors.NewValidationError(
				"board_id must be 1 or 2",
				map[string]interface{}{"received": boardID},
			) // Create validation error
			utils.HandleError(w, appErr) // Return error response
			return
		}

		// Validate pon_id conversion to integer
		ponIDInt, err := strconv.Atoi(ponID)
		// Check if conversion failed or if pon_id is out of range (1-16)
		if err != nil || ponIDInt < 1 || ponIDInt > 16 {
			appErr := apperrors.NewValidationError(
				"pon_id must be between 1 and 16",
				map[string]interface{}{"received": ponID},
			) // Create validation error
			utils.HandleError(w, appErr) // Return error response
			return
		}

		// Store validated values into request context for easier access in handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "boardID", boardIDInt)
		ctx = context.WithValue(ctx, "ponID", ponIDInt)

		next.ServeHTTP(w, r.WithContext(ctx)) // Proceed with the updated context
	})
}

// ValidateOnuIDParam validates onu_id URL parameter,
// ensuring it is a valid integer within the expected range (1-128).
func ValidateOnuIDParam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		onuID := chi.URLParam(r, "onu_id") // Get onu_id from URL

		// Validate onu_id conversion to integer
		onuIDInt, err := strconv.Atoi(onuID)
		// Check if conversion failed or if onu_id is out of range (1-128)
		if err != nil || onuIDInt < 1 || onuIDInt > 128 {
			appErr := apperrors.NewValidationError(
				"onu_id must be between 1 and 128",
				map[string]interface{}{"received": onuID},
			) // Create validation error
			utils.HandleError(w, appErr) // Return error response
			return
		}

		// Store validated value into context
		ctx := context.WithValue(r.Context(), "onuID", onuIDInt)
		next.ServeHTTP(w, r.WithContext(ctx)) // Proceed with the updated context
	})
}
