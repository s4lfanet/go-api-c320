package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/cors"
)

// CorsMiddleware is a middleware function that sets up CORS (Cross-Origin Resource Sharing)
// Configuration is loaded from environment variables for easy maintenance
func CorsMiddleware() func(next http.Handler) http.Handler {
	// Get CORS configuration from environment variables
	allowedOrigins := getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"https://*", "http://*"})
	allowedMethods := getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-API-Key", "X-Request-ID"})
	allowCredentials := getEnvAsBool("CORS_ALLOW_CREDENTIALS", false)
	maxAge := getEnvAsInt("CORS_MAX_AGE", 300)

	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,                   // Configurable allowed origins
		AllowedMethods:   allowedMethods,                   // Configurable HTTP methods
		AllowedHeaders:   allowedHeaders,                   // Configurable allowed headers
		ExposedHeaders:   []string{"Link", "X-Request-ID"}, // Expose headers including Request ID
		AllowCredentials: allowCredentials,                 // Configurable credential support
		MaxAge:           maxAge,                           // Configurable preflight cache duration
	})
}

// getEnvAsSlice retrieves an environment variable as a comma-separated slice
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}

// getEnvAsBool retrieves an environment variable as boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
