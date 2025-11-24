package main

import (
	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/model"
	"bus-booking/shared/db"
	"bus-booking/shared/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.MustLoadConfig()
	logger.MustSetupLogger(&cfg.Log)

	mm := db.MustNewMigrationManager(&cfg.Database)
	defer func() {
		if err := mm.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close migrator")
		}
	}()

	models := []interface{}{
		&model.Booking{},
		&model.BookingSeat{},
		&model.SeatStatus{},
		&model.PaymentMethod{},
		&model.Feedback{},
	}

	if err := mm.RunMigrations(models...); err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}

	log.Info().Msg("Booking-service migration completed successfully!")
}
