package db

import (
	"fmt"

	"bus-booking/shared/config"

	"github.com/rs/zerolog/log"
)

// MigrationManager handles database migrations for all services
type MigrationManager struct {
	db *DatabaseManager
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(cfg *config.DatabaseConfig, env string) (*MigrationManager, error) {
	dbManager, err := NewPostgresConnection(cfg, env)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database for migration: %w", err)
	}

	return &MigrationManager{
		db: dbManager,
	}, nil
}

// EnableExtensions enables common PostgreSQL extensions
func (mm *MigrationManager) EnableExtensions() error {
	extensions := []string{
		"uuid-ossp", // For UUID generation
		"citext",    // Case-insensitive text
		"pg_trgm",   // For trigram matching (text search)
	}

	for _, ext := range extensions {
		if err := mm.db.DB.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\"", ext)).Error; err != nil {
			log.Warn().Str("extension", ext).Err(err).Msg("Failed to create extension")
		} else {
			log.Info().Str("extension", ext).Msg("Extension enabled")
		}
	}

	return nil
}

// RunMigrations runs database migrations for the provided models
func (mm *MigrationManager) RunMigrations(models ...interface{}) error {
	if len(models) == 0 {
		log.Warn().Msg("No models provided for migration")
		return nil
	}

	log.Info().Int("model_count", len(models)).Msg("Starting database migration...")

	// Enable extensions first
	if err := mm.EnableExtensions(); err != nil {
		return fmt.Errorf("failed to enable extensions: %w", err)
	}

	// Run auto migration
	if err := mm.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	log.Info().Msg("Database migration completed successfully!")
	return nil
}

// Close closes the database connection
func (mm *MigrationManager) Close() error {
	return mm.db.Close()
}

// GetDB returns the underlying database manager (for advanced usage)
func (mm *MigrationManager) GetDB() *DatabaseManager {
	return mm.db
}
