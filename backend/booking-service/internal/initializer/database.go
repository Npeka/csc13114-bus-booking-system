package initializer

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/db"
)

// InitDatabase initializes the database connection and runs migrations
func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	log.Println("Initializing database connection...")

	// Create database connection
	database, err := db.NewDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// Run seeders
	if err := database.Seed(); err != nil {
		return nil, fmt.Errorf("failed to run database seeders: %w", err)
	}

	log.Println("Database initialized successfully")
	return database.DB, nil
}
