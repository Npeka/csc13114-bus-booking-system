// Generic migration template for any service
// Copy this file to your service's cmd/migrate/main.go and modify the models section

package main

import (
	"log"

	sharedConfig "bus-booking/shared/config"
	sharedDB "bus-booking/shared/db"
	// Import your service models here
	// "bus-booking/your-service/internal/model"
)

func main() {
	// Load base config (works with any service that embeds BaseConfig)
	cfg, err := sharedConfig.LoadConfig[sharedConfig.BaseConfig]()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create migration manager
	migrationManager, err := sharedDB.NewMigrationManager(&cfg.Database, cfg.Server.Environment)
	if err != nil {
		log.Fatal("Failed to create migration manager:", err)
	}
	defer migrationManager.Close()

	// Define models for your service
	models := []interface{}{
		// Add your service models here
		// &model.YourModel{},
	}

	if len(models) == 0 {
		log.Println("No models defined for migration")
		return
	}

	// Run migrations
	if err := migrationManager.RunMigrations(models...); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed successfully!")
}
