package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
)

func TestLogger_SuccessfulRequest(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	logger := zerolog.New(&buf).With().Timestamp().Logger()

	// Test handler that returns 200 OK
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap handler with logger middleware
	middleware := Logger(logger)
	handler := middleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "192.168.1.1:12345"

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	// Check that log was written
	logOutput := buf.String()
	if logOutput == "" {
		t.Error("Expected log output, got empty string")
	}

	// Parse log as JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("Failed to parse log as JSON: %v", err)
	}

	// Verify log fields
	if logEntry["level"] != "info" {
		t.Errorf("Expected level 'info', got %v", logEntry["level"])
	}

	if logEntry["message"] != "incoming_request" {
		t.Errorf("Expected message 'incoming_request', got %v", logEntry["message"])
	}

	if logEntry["method"] != "GET" {
		t.Errorf("Expected method 'GET', got %v", logEntry["method"])
	}

	if logEntry["path"] != "/api/test" {
		t.Errorf("Expected path '/api/test', got %v", logEntry["path"])
	}

	if logEntry["status_code"] != float64(200) {
		t.Errorf("Expected status_code 200, got %v", logEntry["status_code"])
	}
}

func TestLogger_ErrorResponse(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).With().Timestamp().Logger()

	// Test handler that returns 500 error
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	})

	middleware := Logger(logger)
	handler := middleware(testHandler)

	req := httptest.NewRequest("POST", "/api/error", bytes.NewBufferString("test body"))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}

	// Parse log
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("Failed to parse log as JSON: %v", err)
	}

	if logEntry["status_code"] != float64(500) {
		t.Errorf("Expected status_code 500, got %v", logEntry["status_code"])
	}

	if logEntry["method"] != "POST" {
		t.Errorf("Expected method 'POST', got %v", logEntry["method"])
	}
}

func TestLogger_PanicRecovery(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).With().Timestamp().Logger()

	// Test handler that panics
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	middleware := Logger(logger)
	handler := middleware(testHandler)

	req := httptest.NewRequest("GET", "/api/panic", nil)
	rr := httptest.NewRecorder()

	// Should not panic - middleware should recover
	handler.ServeHTTP(rr, req)

	// Check that panic was logged as error
	logOutput := buf.String()
	if logOutput == "" {
		t.Error("Expected log output for panic, got empty string")
	}

	// Should contain both error log for panic and info log for request
	// The panic recovery should log at ERROR level
	if !bytes.Contains(buf.Bytes(), []byte("error")) {
		t.Error("Expected error level log for panic")
	}

	if !bytes.Contains(buf.Bytes(), []byte("incoming_request_panic")) {
		t.Error("Expected 'incoming_request_panic' message in log")
	}
}

func TestLogger_DifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			var buf bytes.Buffer
			logger := zerolog.New(&buf).With().Timestamp().Logger()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := Logger(logger)
			handler := middleware(testHandler)

			req := httptest.NewRequest(method, "/api/test", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Parse log
			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Errorf("Failed to parse log as JSON: %v", err)
			}

			if logEntry["method"] != method {
				t.Errorf("Expected method '%s', got %v", method, logEntry["method"])
			}
		})
	}
}

func TestLogger_LogsElapsedTime(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).With().Timestamp().Logger()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Logger(logger)
	handler := middleware(testHandler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Parse log
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("Failed to parse log as JSON: %v", err)
	}

	// Check that elapsed_time field exists
	if _, ok := logEntry["elapsed_time"]; !ok {
		t.Error("Expected 'elapsed_time' field in log")
	}
}

func TestLogger_BytesInOut(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf).With().Timestamp().Logger()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response body"))
	})

	middleware := Logger(logger)
	handler := middleware(testHandler)

	requestBody := "request data"
	req := httptest.NewRequest("POST", "/api/test", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Parse log
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("Failed to parse log as JSON: %v", err)
	}

	// Check bytes_in (content length)
	if _, ok := logEntry["bytes_in"]; !ok {
		t.Error("Expected 'bytes_in' field in log")
	}

	// Check bytes_out
	if _, ok := logEntry["bytes_out"]; !ok {
		t.Error("Expected 'bytes_out' field in log")
	}
}
