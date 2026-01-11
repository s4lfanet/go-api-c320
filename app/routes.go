package app

import (
	"net/http"
	"os"
	"time"

	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/handler"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func loadRoutes(onuHandler *handler.OnuHandler, ponHandler *handler.PonHandler, profileHandler *handler.ProfileHandler, cardHandler *handler.CardHandler, provisionHandler *handler.ProvisionHandler, vlanHandler handler.VLANHandlerInterface, trafficHandler handler.TrafficHandlerInterface, onuMgmtHandler handler.ONUManagementHandlerInterface) http.Handler { // Function to configure and return the HTTP router

	// Initialize logger
	l := log.Output(zerolog.ConsoleWriter{ // Create a new logger with console writer output
		Out: os.Stdout, // Set output to standard out
	})

	// Initialize router using chi
	router := chi.NewRouter() // Create a new instance of Chi router

	// Request ID tracking (should be first for all requests)
	router.Use(middleware.RequestID) // Add unique request ID to each request

	// Security middleware
	router.Use(middleware.SecurityHeaders)                  // Apply security headers middleware
	router.Use(middleware.RequestTimeout(90 * time.Second)) // Set request timeout to the 90s (allows cold cache SNMP queries up to the 60s without a timeout)
	router.Use(middleware.RateLimiter(100, 200))            // Apply rate limiting: 100 requests per second, burst up to 200
	router.Use(middleware.MaxBodySize(1 << 20))             // Limit request body size to 1MB (1 << 20 bytes)

	// Middleware for logging requests
	router.Use(middleware.Logger(l)) // Apply logging middleware using the initialized logger

	// Middleware for CORS (now configurable via environment variables)
	router.Use(middleware.CorsMiddleware()) // Apply Cross-Origin Resource Sharing (CORS) middleware

	// Define a simple root endpoint
	router.Get("/", rootHandler) // Register the GET handler for the root path "/"

	// Create a group for /api/v1/
	apiV1Group := chi.NewRouter() // Create a new router instance for API version 1 group

	// Define routes for /api/v1/
	apiV1Group.Route("/board", func(r chi.Router) { // Create a route group starting with "/board"
		r.Route("/{board_id}/pon/{pon_id}", func(r chi.Router) { // Nested route group with board_id and pon_id parameters
			// Apply validation middleware for board_id and pon_id
			r.Use(middleware.ValidateBoardPonParams) // specific middleware to validate these parameters

			r.Get("/", onuHandler.GetByBoardIDAndPonID)             // GET /board/{board_id}/pon/{pon_id}/ - Fetch ONUs by board and PON
			r.Delete("/", onuHandler.DeleteCache)                   // DELETE /board/{board_id}/pon/{pon_id}/ - Delete cache for board/pon
			r.Get("/info", ponHandler.GetPonPortInfo)               // GET .../info - Fetch PON port information
			r.Get("/onu_id/empty", onuHandler.GetEmptyOnuID)        // GET .../onu_id/empty - Fetch empty ONU IDs
			r.Get("/onu_id_sn", onuHandler.GetOnuIDAndSerialNumber) // GET .../onu_id_sn - Fetch ONU IDs and serial numbers
			r.Get("/onu_id/update", onuHandler.UpdateEmptyOnuID)    // GET .../onu_id/update - Update empty ONU IDs (Note: GET used for update seems unusual but following existing code)

			// Routes with onu_id parameter
			r.Route("/onu/{onu_id}", func(r chi.Router) { // Nested route group for specific ONU ID
				r.Use(middleware.ValidateOnuIDParam)             // Validate onu_id parameter
				r.Get("/", onuHandler.GetByBoardIDPonIDAndOnuID) // GET .../onu/{onu_id} - Fetch specific ONU details
			})
		})
	})

	// Define routes for /api/v1/paginate
	apiV1Group.Route("/paginate", func(r chi.Router) { // Create a route group for pagination
		r.Route("/board/{board_id}/pon/{pon_id}", func(r chi.Router) { // Nested route group with board and PON IDs
			r.Use(middleware.ValidateBoardPonParams)                // Apply parameter validation
			r.Get("/", onuHandler.GetByBoardIDAndPonIDWithPaginate) // GET .../ - Fetch paginated ONU list
		})
	})

	// Define routes for /api/v1/profiles
	apiV1Group.Route("/profiles", func(r chi.Router) { // Create a route group for profiles
		r.Route("/traffic", func(r chi.Router) { // Nested route group for traffic profiles
			r.Get("/", profileHandler.GetAllTrafficProfiles)         // GET /profiles/traffic - Fetch all traffic profiles
			r.Get("/{profile_id}", profileHandler.GetTrafficProfile) // GET /profiles/traffic/{profile_id} - Fetch specific traffic profile
		})
		r.Route("/vlan", func(r chi.Router) { // Nested route group for VLAN profiles
			r.Get("/", profileHandler.GetAllVlanProfiles) // GET /profiles/vlan - Fetch all VLAN profiles
		})
	})

	// Define routes for /api/v1/system
	apiV1Group.Route("/system", func(r chi.Router) { // Create a route group for system information
		r.Route("/cards", func(r chi.Router) { // Nested route group for card/slot info
			r.Get("/", cardHandler.GetAllCards)                  // GET /system/cards - Fetch all cards
			r.Get("/{rack}/{shelf}/{slot}", cardHandler.GetCard) // GET /system/cards/{rack}/{shelf}/{slot} - Fetch specific card
		})
	})

	// Define routes for /api/v1/onu (provisioning)
	apiV1Group.Route("/onu", func(r chi.Router) {
		r.Get("/unconfigured", provisionHandler.GetUnconfiguredONUs)            // GET all unconfigured ONUs
		r.Get("/unconfigured/{pon}", provisionHandler.GetUnconfiguredONUsByPON) // GET unconfigured ONUs by PON port
		r.Post("/register", provisionHandler.RegisterONU)                       // POST register new ONU
		r.Delete("/{pon}/{onu_id}", provisionHandler.DeleteONU)                 // DELETE ONU
	})

	// Define routes for /api/v1/vlan (VLAN management)
	apiV1Group.Route("/vlan", func(r chi.Router) {
		r.Get("/onu/{pon}/{onu_id}", vlanHandler.GetONUVLAN)    // GET ONU VLAN configuration
		r.Get("/service-ports", vlanHandler.GetAllServicePorts) // GET all service-port configurations
		r.Post("/onu", vlanHandler.ConfigureVLAN)               // POST configure ONU VLAN
		r.Put("/onu", vlanHandler.ModifyVLAN)                   // PUT modify ONU VLAN
		r.Delete("/onu/{pon}/{onu_id}", vlanHandler.DeleteVLAN) // DELETE ONU VLAN
	})

	// Define routes for /api/v1/traffic (Traffic profile management)
	apiV1Group.Route("/traffic", func(r chi.Router) {
		// DBA Profile routes
		r.Get("/dba-profiles", trafficHandler.GetAllDBAProfiles)         // GET all DBA profiles
		r.Get("/dba-profile/{name}", trafficHandler.GetDBAProfile)       // GET specific DBA profile
		r.Post("/dba-profile", trafficHandler.CreateDBAProfile)          // POST create DBA profile
		r.Put("/dba-profile", trafficHandler.ModifyDBAProfile)           // PUT modify DBA profile
		r.Delete("/dba-profile/{name}", trafficHandler.DeleteDBAProfile) // DELETE DBA profile

		// TCONT routes
		r.Get("/tcont/{pon}/{onu_id}/{tcont_id}", trafficHandler.GetONUTCONT)    // GET T-CONT configuration
		r.Post("/tcont", trafficHandler.ConfigureTCONT)                          // POST configure T-CONT
		r.Delete("/tcont/{pon}/{onu_id}/{tcont_id}", trafficHandler.DeleteTCONT) // DELETE T-CONT

		// GEMPort routes
		r.Post("/gemport", trafficHandler.ConfigureGEMPort)                            // POST configure GEM port
		r.Delete("/gemport/{pon}/{onu_id}/{gemport_id}", trafficHandler.DeleteGEMPort) // DELETE GEM port
	})

	// Define routes for /api/v1/onu-management (ONU lifecycle management)
	apiV1Group.Route("/onu-management", func(r chi.Router) {
		r.Post("/reboot", onuMgmtHandler.RebootONU)             // POST reboot ONU
		r.Post("/block", onuMgmtHandler.BlockONU)               // POST block (disable) ONU
		r.Post("/unblock", onuMgmtHandler.UnblockONU)           // POST unblock (enable) ONU
		r.Put("/description", onuMgmtHandler.UpdateDescription) // PUT update ONU description
		r.Delete("/{pon}/{onu_id}", onuMgmtHandler.DeleteONU)   // DELETE ONU configuration
	})

	// Mount /api/v1/ to root router
	router.Mount("/api/v1", apiV1Group) // Mount the API v1 group to the main router under /api/v1 prefix

	return router // Return the configured router
}

// rootHandler is a simple handler for a root endpoint
func rootHandler(w http.ResponseWriter, _ *http.Request) { // Handler function for the root URL
	w.WriteHeader(http.StatusOK)                                // Set the HTTP status code to 200 OK
	_, _ = w.Write([]byte("Hello, this is the root endpoint!")) // Write a simple welcome message to the response body
}
