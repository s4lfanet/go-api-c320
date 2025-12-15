package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCorsMiddleware_DefaultConfig(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	os.Unsetenv("CORS_ALLOWED_METHODS")
	os.Unsetenv("CORS_ALLOWED_HEADERS")
	os.Unsetenv("CORS_ALLOW_CREDENTIALS")
	os.Unsetenv("CORS_MAX_AGE")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	middleware := CorsMiddleware()
	wrappedHandler := middleware(handler)

	// Test CORS preflight request
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	// Should allow CORS
	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Errorf("Expected status OK or NoContent for OPTIONS, got %d", rr.Code)
	}
}

func TestCorsMiddleware_CustomOrigins(t *testing.T) {
	os.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com,https://test.com")
	defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorsMiddleware()
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

func TestCorsMiddleware_CustomMethods(t *testing.T) {
	os.Setenv("CORS_ALLOWED_METHODS", "GET,POST,PUT")
	defer os.Unsetenv("CORS_ALLOWED_METHODS")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorsMiddleware()
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Errorf("Expected status OK or NoContent, got %d", rr.Code)
	}
}

func TestCorsMiddleware_CustomHeaders(t *testing.T) {
	os.Setenv("CORS_ALLOWED_HEADERS", "Authorization,Content-Type,X-Custom-Header")
	defer os.Unsetenv("CORS_ALLOWED_HEADERS")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorsMiddleware()
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

func TestCorsMiddleware_AllowCredentials(t *testing.T) {
	os.Setenv("CORS_ALLOW_CREDENTIALS", "true")
	defer os.Unsetenv("CORS_ALLOW_CREDENTIALS")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorsMiddleware()
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

func TestCorsMiddleware_MaxAge(t *testing.T) {
	os.Setenv("CORS_MAX_AGE", "600")
	defer os.Unsetenv("CORS_MAX_AGE")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CorsMiddleware()
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Errorf("Expected status OK or NoContent, got %d", rr.Code)
	}
}

func TestGetEnvAsSlice(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue []string
		expected     []string
	}{
		{
			name:         "Valid comma-separated values",
			key:          "TEST_SLICE",
			envValue:     "value1,value2,value3",
			defaultValue: []string{"default"},
			expected:     []string{"value1", "value2", "value3"},
		},
		{
			name:         "Empty environment variable - use default",
			key:          "TEST_SLICE",
			envValue:     "",
			defaultValue: []string{"default1", "default2"},
			expected:     []string{"default1", "default2"},
		},
		{
			name:         "Values with spaces - should trim",
			key:          "TEST_SLICE",
			envValue:     " value1 , value2 , value3 ",
			defaultValue: []string{"default"},
			expected:     []string{"value1", "value2", "value3"},
		},
		{
			name:         "Single value",
			key:          "TEST_SLICE",
			envValue:     "single",
			defaultValue: []string{"default"},
			expected:     []string{"single"},
		},
		{
			name:         "Values with empty entries - should skip",
			key:          "TEST_SLICE",
			envValue:     "value1,,value2,  ,value3",
			defaultValue: []string{"default"},
			expected:     []string{"value1", "value2", "value3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			result := getEnvAsSlice(tt.key, tt.defaultValue)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}

			for i, val := range result {
				if val != tt.expected[i] {
					t.Errorf("At index %d: expected %s, got %s", i, tt.expected[i], val)
				}
			}
		})
	}
}

func TestGetEnvAsBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "True value",
			key:          "TEST_BOOL",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "False value",
			key:          "TEST_BOOL",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "1 value - parsed as true",
			key:          "TEST_BOOL",
			envValue:     "1",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "0 value - parsed as false",
			key:          "TEST_BOOL",
			envValue:     "0",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "Empty value - use default",
			key:          "TEST_BOOL",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "Invalid value - use default",
			key:          "TEST_BOOL",
			envValue:     "invalid",
			defaultValue: false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			result := getEnvAsBool(tt.key, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "Valid integer",
			key:          "TEST_INT",
			envValue:     "42",
			defaultValue: 10,
			expected:     42,
		},
		{
			name:         "Empty value - use default",
			key:          "TEST_INT",
			envValue:     "",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "Invalid value - use default",
			key:          "TEST_INT",
			envValue:     "invalid",
			defaultValue: 5,
			expected:     5,
		},
		{
			name:         "Negative integer",
			key:          "TEST_INT",
			envValue:     "-10",
			defaultValue: 0,
			expected:     -10,
		},
		{
			name:         "Zero",
			key:          "TEST_INT",
			envValue:     "0",
			defaultValue: 10,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			result := getEnvAsInt(tt.key, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCorsMiddleware_ActualRequest(t *testing.T) {
	os.Unsetenv("CORS_ALLOWED_ORIGINS")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := CorsMiddleware()
	handler := middleware(testHandler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "https://example.com")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	if body := strings.TrimSpace(rr.Body.String()); body != "success" {
		t.Errorf("Expected body 'success', got '%s'", body)
	}
}
