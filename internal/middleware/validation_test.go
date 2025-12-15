package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/utils"
	"github.com/go-chi/chi/v5"
)

func TestValidateBoardPonParams(t *testing.T) {
	tests := []struct {
		name           string
		boardID        string
		ponID          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid board 1 pon 1",
			boardID:        "1",
			ponID:          "1",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid board 2 pon 16",
			boardID:        "2",
			ponID:          "16",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid board 0",
			boardID:        "0",
			ponID:          "1",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid board 3",
			boardID:        "3",
			ponID:          "1",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid board non-numeric",
			boardID:        "abc",
			ponID:          "1",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid pon 0",
			boardID:        "1",
			ponID:          "0",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid pon 17",
			boardID:        "1",
			ponID:          "17",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid pon non-numeric",
			boardID:        "1",
			ponID:          "xyz",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Valid board 1 pon 8",
			boardID:        "1",
			ponID:          "8",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid board 2 pon 10",
			boardID:        "2",
			ponID:          "10",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that checks context values
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If we get here, validation passed
				boardIDInt := r.Context().Value("boardID")
				ponIDInt := r.Context().Value("ponID")

				if boardIDInt == nil || ponIDInt == nil {
					t.Error("Expected boardID and ponID in context")
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Wrap with validation middleware
			handler := ValidateBoardPonParams(testHandler)

			// Create a chi router to set URL params
			r := chi.NewRouter()
			r.With(ValidateBoardPonParams).Get("/board/{board_id}/pon/{pon_id}", testHandler)

			// Create request
			req := httptest.NewRequest("GET", "/board/"+tt.boardID+"/pon/"+tt.ponID, nil)

			// Set chi URL params manually
			rctx := chi.NewRouteContext()

			// Set chi URL params manually
			rctx.URLParams.Add("board_id", tt.boardID)

			// Set chi URL params manually
			rctx.URLParams.Add("pon_id", tt.ponID)

			// Set chi context
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create recorder
			rr := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %v, got %v", tt.expectedStatus, status)
			}

			// If error expected, check the response
			if tt.expectError {
				var response utils.ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}

				if response.Code != http.StatusBadRequest {
					t.Errorf("Expected error code %v, got %v", http.StatusBadRequest, response.Code)
				}
			}
		})
	}
}

func TestValidateOnuIDParam(t *testing.T) {
	tests := []struct {
		name           string
		onuID          string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid ONU ID 1",
			onuID:          "1",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid ONU ID 64",
			onuID:          "64",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Valid ONU ID 128",
			onuID:          "128",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid ONU ID 0",
			onuID:          "0",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid ONU ID 129",
			onuID:          "129",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid ONU ID non-numeric",
			onuID:          "abc",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid ONU ID negative",
			onuID:          "-5",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// If we get here, validation passed
				onuIDInt := r.Context().Value("onuID")

				if onuIDInt == nil {
					t.Error("Expected onuID in context")
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Wrap with validation middleware
			handler := ValidateOnuIDParam(testHandler)

			// Create a request with chi context
			req := httptest.NewRequest("GET", "/onu/"+tt.onuID, nil)

			// Set chi URL params manually
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("onu_id", tt.onuID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Expected status %v, got %v", tt.expectedStatus, status)
			}

			// If error expected, check the response
			if tt.expectError {
				var response utils.ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}

				if response.Code != http.StatusBadRequest {
					t.Errorf("Expected error code %v, got %v", http.StatusBadRequest, response.Code)
				}
			}
		})
	}
}
