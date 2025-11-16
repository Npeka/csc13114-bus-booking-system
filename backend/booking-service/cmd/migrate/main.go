package main

import (
	"log"

	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/model"
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
		&model.Booking{},
		&model.BookingSeat{},
		&model.SeatStatus{},
		&model.PaymentMethod{},
		&model.Feedback{},
	}

	if err := migrationManager.RunMigrations(models...); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("User-service migration completed successfully!")
}
