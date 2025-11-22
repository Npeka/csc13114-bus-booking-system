package main

import (
	"log"

	"bus-booking/payment-service/config"
	"bus-booking/payment-service/internal/model"
	sharedDB "bus-booking/shared/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	migrationManager, err := sharedDB.NewMigrationManager(&cfg.Database, cfg.Server.Environment)
	if err != nil {
		log.Fatal("Failed to create migration manager:", err)
	}
	defer migrationManager.Close()

	models := []interface{}{
		&model.Transaction{},
	}

	if err := migrationManager.RunMigrations(models...); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("User-service migration completed successfully!")
}
