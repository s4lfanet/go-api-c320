package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID_GeneratesID(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get request ID from context
		requestID := r.Context().Value(RequestIDKey)
		if requestID == nil {
			t.Error("Expected requestID in context, got nil")
		}

		// Check that it's a non-empty string
		if requestIDStr, ok := requestID.(string); !ok || requestIDStr == "" {
			t.Errorf("Expected non-empty string requestID, got %v", requestID)
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := RequestID(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	// Don't set X-Request-ID header - should be generated
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check that the X-Request-ID header was set in response
	responseID := rr.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Error("Expected X-Request-ID header in response")
	}
}

func TestRequestID_UsesClientProvidedID(t *testing.T) {
	clientRequestID := "client-provided-id-12345"

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get request ID from context
		requestID := r.Context().Value(RequestIDKey)
		if requestID == nil {
			t.Error("Expected requestID in context, got nil")
		}

		// Check that it matches the client-provided ID
		if requestIDStr, ok := requestID.(string); ok {
			if requestIDStr != clientRequestID {
				t.Errorf("Expected requestID '%s', got '%s'", clientRequestID, requestIDStr)
			}
		} else {
			t.Error("requestID is not a string")
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := RequestID(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", clientRequestID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check that response has the same X-Request-ID
	responseID := rr.Header().Get("X-Request-ID")
	if responseID != clientRequestID {
		t.Errorf("Expected response X-Request-ID '%s', got '%s'", clientRequestID, responseID)
	}
}

func TestRequestID_ResponseHeaderSet(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestID(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check that X-Request-ID header exists in response
	responseID := rr.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Error("Expected X-Request-ID header in response, got empty string")
	}
}

func TestRequestID_UniqueIDsGenerated(t *testing.T) {
	var firstID, secondID string

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestID(testHandler)

	// First request
	req1 := httptest.NewRequest("GET", "/test", nil)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	firstID = rr1.Header().Get("X-Request-ID")

	// Second request
	req2 := httptest.NewRequest("GET", "/test", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	secondID = rr2.Header().Get("X-Request-ID")

	// IDs should be different
	if firstID == "" || secondID == "" {
		t.Error("Expected non-empty request IDs")
	}

	if firstID == secondID {
		t.Errorf("Expected unique IDs, but got same ID: %s", firstID)
	}
}

func TestRequestID_ContextPropagation(t *testing.T) {
	var contextRequestID string

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract requestID from context
		if val := r.Context().Value(RequestIDKey); val != nil {
			contextRequestID = val.(string)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestID(testHandler)

	providedID := "test-id-123"
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", providedID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check that context value matches the header
	if contextRequestID != providedID {
		t.Errorf("Expected context requestID '%s', got '%s'", providedID, contextRequestID)
	}

	// Check that response header matches
	responseID := rr.Header().Get("X-Request-ID")
	if responseID != providedID {
		t.Errorf("Expected response header '%s', got '%s'", providedID, responseID)
	}
}

func TestRequestID_EmptyHeaderGeneratesID(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(RequestIDKey)
		if requestID == nil || requestID == "" {
			t.Error("Expected generated requestID, got nil or empty")
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestID(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "") // Empty header
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Should generate new ID when header is empty
	responseID := rr.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Error("Expected generated X-Request-ID, got empty string")
	}
}

func TestRequestID_HandlerCalledCorrectly(t *testing.T) {
	handlerCalled := false

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	handler := RequestID(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	if body := rr.Body.String(); body != "success" {
		t.Errorf("Expected body 'success', got '%s'", body)
	}
}
