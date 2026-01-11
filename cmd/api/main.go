package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/app"
)

func main() {
	// Load .env file if exists (for local development and production)
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using environment variables")
	}

	// Initialize application
	server := app.New()                                     // Create a new instance of the application
	ctx, cancel := context.WithCancel(context.Background()) // Create a new context with a cancel function
	defer cancel()                                          // Cancel context when the main function is finished

	// Start an application server in a goroutine
	go func() {
		err := server.Start(ctx) // Start the application server
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start server") // Log error message
			cancel()                                           // Cancel context if an error occurred
		}
	}()

	// Create a channel to wait for a signal to stop the application
	stopSignal := make(chan struct{})

	// You can replace the select statement with a simple channel receiver
	<-stopSignal

	// Log that the application is stopping
	log.Info().Msg("Application is stopping")
}
