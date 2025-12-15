package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// Logger is a middleware function that logs incoming HTTP requests and their details
// using the provided zerolog.Logger instance. It captures information such as request
// time, remote address, request path, protocol, method, user agent, response status,
// bytes in/out, and elapsed time. It also handles panics and logs them as errors
func Logger(logger zerolog.Logger) func(next http.Handler) http.Handler { // Function to return Logger middleware handler
	return func(next http.Handler) http.Handler { // Return the actual middleware function
		fn := func(w http.ResponseWriter, r *http.Request) { // Define the handler function
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor) // Wrap the response writer to capture status and size
			startTime := time.Now()                                 // Record start time of request

			defer func() { // Defer logging execution until after the request is processed
				endTime := time.Now()                 // Record end time
				elapsedTime := endTime.Sub(startTime) // Calculate duration

				if r := recover(); r != nil && r != http.ErrAbortHandler { // Recover from panics
					logger.Error().Interface("recover", r).Bytes("stack", debug.Stack()).Msg("incoming_request_panic") // Log panic details
					ww.WriteHeader(http.StatusInternalServerError)                                                     // Respond with 500 Internal Server Error
				}

				// Log request details using structured logging
				logger.Info().Fields(map[string]interface{}{
					"time":         startTime.Format(time.RFC3339), // Format start time as RFC3339
					"remote_addr":  r.RemoteAddr,                   // Remote IP address
					"path":         r.URL.Path,                     // Request path
					"proto":        r.Proto,                        // Protocol version
					"method":       r.Method,                       // HTTP method
					"user_agent":   r.UserAgent(),                  // User Agent string
					"status":       http.StatusText(ww.Status()),   // Text description of HTTP status
					"status_code":  ww.Status(),                    // Numeric HTTP status code
					"bytes_in":     r.ContentLength,                // Request content length
					"bytes_out":    ww.BytesWritten(),              // Response body size
					"elapsed_time": elapsedTime.String(),           // Processing duration as a string
				}).Msg("incoming_request") // Log message
			}()

			next.ServeHTTP(ww, r) // Serve the request using the wrapped response writer
		}

		return http.HandlerFunc(fn) // Return the handler function
	}
}
