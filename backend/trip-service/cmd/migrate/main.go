package main

import (
	"bus-booking/shared/db"
	"bus-booking/shared/logger"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/model"

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
		&model.Operator{},
		&model.Route{},
		&model.Bus{},
		&model.Seat{},
		&model.Trip{},
	}

	if err := mm.RunMigrations(models...); err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}

	log.Info().Msg("Trip-service migration completed successfully!")
}
