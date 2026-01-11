package app

import (
	"context"
	"fmt"
	"net/http"

	rds "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	"github.com/s4lfanet/go-api-c320/internal/handler"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"github.com/s4lfanet/go-api-c320/internal/usecase"
	"github.com/s4lfanet/go-api-c320/pkg/graceful"
	"github.com/s4lfanet/go-api-c320/pkg/redis"
	"github.com/s4lfanet/go-api-c320/pkg/snmp"
)

// App represents the main application structure that holds the HTTP router
// and manages the application lifecycle, including dependency initialization
// and server startup.
type App struct { // Define the App struct
	router http.Handler // HTTP handler for routing requests
}

// New creates and returns a new instance of the App with initialized dependencies.
// It prepares the application for startup but does not start the server.
func New() *App { // Factory function to create a new App instance
	return &App{} // Return a pointer to a new App struct
}

// Start initializes the application components, sets up connections to external services
// (Redis and SNMP), and starts the HTTP server. It handles graceful shutdown on context
// cancellation and ensures proper cleanup of resources.
//
// Parameters:
//   - ctx: context.Context for cancellation and timeout propagation
//
// Returns:
//   - error: returns any error that occurs during application startup or shutdown
func (a *App) Start(ctx context.Context) error { // Method to start the application

	// Load configuration from environment variables (no config file needed)
	// Board/PON OID mappings are generated dynamically using mathematical formulas
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load config")
		return err
	}

	// Initialize Redis client
	redisClient := redis.NewRedisClient(cfg) // Create a new Redis client using the configuration

	// Check Redis connection
	err = redisClient.Ping(ctx).Err() // Ping Redis to verify connection
	if err != nil {                   // Check if ping failed
		log.Error().Err(err).Msg("Failed to ping Redis server") // Log the error
	} else { // If ping succeeded
		log.Info().Msg("Redis server successfully connected") // Log success message
	}

	// Close Redis client
	defer func(redisClient *rds.Client) { // Defer closure of a Redis client until Start function exits
		err := redisClient.Close() // Close the Redis connection
		if err != nil {            // Check if closing failed
			log.Error().Err(err).Msg("Failed to close Redis client") // Log the error
		}
	}(redisClient) // Pass redisClient to the deferred function

	// Initialize SNMP connection
	snmpConn, err := snmp.SetupSnmpConnection(cfg) // Setup SNMP connection using configuration
	if err != nil {                                // Check if setup failed
		log.Error().Err(err).Msg("Failed to setup SNMP connection") // Log the error
		return fmt.Errorf("SNMP setup failed: %w", err)
	}

	// Check SNMP connection
	/*
		if SNMP Connection with wrong credentials in SNMP v3, return error is nil
		if SNMP Connection with the wrong Port in SNMP v2 v2c, return error is nil
		if SNMP Connection with wrong community v2 v2c, return error is nil

		Connect creates and opens a socket. Because UDP is a connectionless protocol,
		you won't know if the remote host is responding until you send packets.
		Neither will you know if the host is regularly disappearing and reappearing.
	*/

	if snmpConn.Connect() != nil { // Attempt to connect to SNMP agent (UDP socket creation)
		log.Error().Err(err).Msg("Failed to connect to SNMP server") // Log connection failure
	} else { // If connection setup (socket creation) succeeded
		log.Info().Msg("SNMP server successfully connected") // Log success message
	}

	// Close SNMP connection after application shutdown
	defer func() { // Defer closure of SNMP connection
		if err := snmpConn.Conn.Close(); err != nil { // Close the SNMP connection and check for error
			log.Error().Err(err).Msg("Failed to close SNMP connection") // Log the error
		}
	}()

	// Initialize repository
	snmpRepo := repository.NewPonRepository(snmpConn.Target, snmpConn.Community, snmpConn.Port) // Create a new PON repository with SNMP details
	redisRepo := repository.NewOnuRedisRepo(redisClient)                                        // Create new ONU Redis repository
	onuRepo := repository.NewOnuRepository(snmpConn, cfg)                                       // Create new ONU repository for monitoring

	// Initialize Telnet session manager
	telnetCfg := config.LoadTelnetConfig()                                // Load telnet configuration
	telnetSessionManager := repository.GetGlobalSessionManager(telnetCfg) // Get global telnet session manager

	// Initialize usecase
	onuUsecase := usecase.NewOnuUsecase(snmpRepo, redisRepo, cfg)                                // Create new ONU usecase with repositories and config
	ponUsecase := usecase.NewPonUsecase(snmpRepo, redisRepo, cfg)                                // Create new PON usecase with repositories and config
	profileUsecase := usecase.NewProfileUsecase(snmpRepo, redisRepo, cfg)                        // Create new Profile usecase with repositories and config
	cardUsecase := usecase.NewCardUsecase(snmpRepo, redisRepo, cfg)                              // Create new Card usecase with repositories and config
	provisionUsecase := usecase.NewProvisionUsecase(telnetSessionManager, cfg)                   // Create new Provision usecase with telnet manager
	vlanUsecase := usecase.NewVLANUsecase(telnetSessionManager, cfg)                             // Create new VLAN usecase with telnet manager
	trafficUsecase := usecase.NewTrafficUsecase(telnetSessionManager, cfg)                       // Create new Traffic usecase with telnet manager
	onuMgmtUsecase := usecase.NewONUManagementUsecase(telnetSessionManager, cfg)                 // Create new ONU Management usecase with telnet manager
	batchUsecase := usecase.NewBatchOperationsUsecase(telnetSessionManager, onuMgmtUsecase, cfg) // Create new Batch Operations usecase
	monitoringUsecase := usecase.NewMonitoringUsecase(snmpConn, cfg, onuRepo)                    // Create new Monitoring usecase (Phase 7.1)

	// Initialize handler
	onuHandler := handler.NewOnuHandler(onuUsecase)                      // Create new ONU handler with usecase
	ponHandler := handler.NewPonHandler(ponUsecase)                      // Create new PON handler with usecase
	profileHandler := handler.NewProfileHandler(profileUsecase)          // Create new Profile handler with usecase
	cardHandler := handler.NewCardHandler(cardUsecase)                   // Create new Card handler with usecase
	provisionHandler := handler.NewProvisionHandler(provisionUsecase)    // Create new Provision handler with usecase
	vlanHandler := handler.NewVLANHandler(vlanUsecase)                   // Create new VLAN handler with usecase
	trafficHandler := handler.NewTrafficHandler(trafficUsecase)          // Create new Traffic handler with usecase
	onuMgmtHandler := handler.NewONUManagementHandler(onuMgmtUsecase)    // Create new ONU Management handler with usecase
	batchHandler := handler.NewBatchOperationsHandler(batchUsecase)      // Create new Batch Operations handler with usecase
	monitoringHandler := handler.NewMonitoringHandler(monitoringUsecase) // Create new Monitoring handler (Phase 7.1)

	// Initialize Config Backup handler (Phase 6.2)
	configBackupUsecase := usecase.NewConfigBackupUsecase(cfg, onuMgmtUsecase, vlanUsecase, trafficUsecase, provisionUsecase) // Create config backup usecase
	configBackupHandler := handler.NewConfigBackupHandler(configBackupUsecase)                                                // Create new Config Backup handler

	// Initialize router
	a.router = loadRoutes(onuHandler, ponHandler, profileHandler, cardHandler, provisionHandler, vlanHandler, trafficHandler, onuMgmtHandler, batchHandler, configBackupHandler, monitoringHandler) // Load all routes and middleware, assigning to app router

	// Start server
	addr := "8081"          // Define the server address/port
	server := &http.Server{ // Create a new HTTP server struct
		Addr:    ":" + addr, // Set the address
		Handler: a.router,   // Set the handler (router)
	}

	// Start server at given address
	log.Info().Msgf("Application started at %s", addr) // Log startup message with address

	// Graceful shutdown
	return graceful.Shutdown(ctx, server) // Start a server with graceful shutdown handling
}
