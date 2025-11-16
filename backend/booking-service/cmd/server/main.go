package main

import (
	"log"
	"os"

	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/initializer"
)

func main() {
	log.Println("Starting Booking Service...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded - Environment: %s, Port: %d", cfg.Server.Environment, cfg.Server.Port)

	// Initialize database
	db, err := initializer.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize services
	bookingHandler, err := initializer.InitServices(db)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Start HTTP server
	if err := initializer.InitServer(cfg, bookingHandler); err != nil {
		log.Printf("Server shutdown with error: %v", err)
		os.Exit(1)
	}

	log.Println("Booking Service stopped")
}
