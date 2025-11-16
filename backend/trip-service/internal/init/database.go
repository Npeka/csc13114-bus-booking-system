package appinit

import (
	"fmt"

	sharedCfg "bus-booking/shared/config"
	sharedDB "bus-booking/shared/db"
	"bus-booking/trip-service/config"

	"github.com/rs/zerolog/log"
)

func InitDatabase(cfg *config.Config) (*sharedDB.DatabaseManager, error) {
	log.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Str("ssl_mode", cfg.Database.SSLMode).
		Msg("Connecting to PostgreSQL database")

	// Convert local config to shared config
	dbConfig := &sharedCfg.DatabaseConfig{
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
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	log.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Str("ssl_mode", cfg.Database.SSLMode).
		Msg("Successfully connected to PostgreSQL database")

	return database, nil
}
