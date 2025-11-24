package main

import (
	"bus-booking/shared/db"
	"bus-booking/shared/logger"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/model"

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
		&model.User{},
	}

	if err := mm.RunMigrations(models...); err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}

	log.Info().Msg("User-service migration completed successfully!")
}
