package main

import (
	"log"

	"bus-booking/shared/db"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/model"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	mm := db.MustNewMigrationManager(&cfg.Database)
	defer mm.Close()

	models := []interface{}{
		&model.User{},
	}

	if err := mm.RunMigrations(models...); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("User-service migration completed successfully!")
}
