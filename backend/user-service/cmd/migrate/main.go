package main

import (
	"log"

	sharedDB "bus-booking/shared/db"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/model"
)

func main() {
	// Load service-specific config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create migration manager
	migrationManager, err := sharedDB.NewMigrationManager(&cfg.Database, cfg.Server.Environment)
	if err != nil {
		log.Fatal("Failed to create migration manager:", err)
	}
	defer migrationManager.Close()

	// Define models for user-service
	models := []interface{}{
		&model.User{},
		// Add more user-service models here
	}

	// Run migrations
	if err := migrationManager.RunMigrations(models...); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("User-service migration completed successfully!")
}
