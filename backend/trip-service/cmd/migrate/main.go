package main

import (
	"log"

	sharedDB "bus-booking/shared/db"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/model"
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
		&model.Operator{},
		&model.Route{},
		&model.Bus{},
		&model.Seat{},
		&model.Trip{},
	}

	if err := migrationManager.RunMigrations(models...); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("User-service migration completed successfully!")
}
