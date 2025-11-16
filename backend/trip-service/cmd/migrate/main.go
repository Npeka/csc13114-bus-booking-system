package main

import (
	"log"

	sharedConfig "bus-booking/shared/config"
	sharedDB "bus-booking/shared/db"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/model"
)

func main() {
	log.Println("Starting database migration...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	dbConfig := &sharedConfig.DatabaseConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Name:            cfg.Database.Name,
		Username:        cfg.Database.Username,
		Password:        cfg.Database.Password,
		SSLMode:         cfg.Database.SSLMode,
		TimeZone:        cfg.Database.TimeZone,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	}

	database, err := sharedDB.NewPostgresConnection(dbConfig, "development")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable UUID extension
	if err := database.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Warning: Failed to create uuid-ossp extension: %v", err)
	}

	// Auto-migrate models
	err = database.DB.AutoMigrate(
		&model.Operator{},
		&model.Route{},
		&model.Bus{},
		&model.Seat{},
		&model.Trip{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migration completed successfully!")
}
