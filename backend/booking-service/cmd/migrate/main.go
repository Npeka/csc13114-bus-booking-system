package main

import (
	"log"
	"os"

	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/initializer"
)

func main() {
	log.Println("Starting database migration...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database and run migrations
	_, err = initializer.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
		os.Exit(1)
	}

	log.Println("Database migration completed successfully")
}
