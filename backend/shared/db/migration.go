package db

import (
	"fmt"
	"strings"

	"bus-booking/shared/config"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// MigrationManager handles database migrations for all services
type MigrationManager struct {
	db *DatabaseManager
}

func MustNewMigrationManager(cfg *config.DatabaseConfig) *MigrationManager {
	mm, err := NewMigrationManager(cfg)
	if err != nil {
		panic(err)
	}
	return mm
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(cfg *config.DatabaseConfig) (*MigrationManager, error) {
	dbManager, err := NewPostgresConnection(cfg)
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

	// Handle table schema synchronization for each model
	for _, model := range models {
		if err := mm.syncTableSchema(model); err != nil {
			return fmt.Errorf("failed to sync schema for model %T: %w", model, err)
		}
	}

	// Run auto migration to ensure all new fields are created
	if err := mm.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	log.Info().Msg("Database migration completed successfully!")
	return nil
}

// syncTableSchema synchronizes table schema with model definition
func (mm *MigrationManager) syncTableSchema(model interface{}) error {
	stmt := &gorm.Statement{DB: mm.db.DB}
	if err := stmt.Parse(model); err != nil {
		return fmt.Errorf("failed to parse model: %w", err)
	}

	tableName := stmt.Schema.Table
	log.Info().Str("table", tableName).Msg("Synchronizing table schema...")

	// Check if table exists
	var tableExists bool
	err := mm.db.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = ?)", tableName).Scan(&tableExists).Error
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if !tableExists {
		log.Info().Str("table", tableName).Msg("Table doesn't exist, will be created by AutoMigrate")
		return nil
	}

	// Get current columns from database
	currentColumns, err := mm.getCurrentColumns(tableName)
	if err != nil {
		return fmt.Errorf("failed to get current columns: %w", err)
	}

	// Get expected columns from model
	expectedColumns := mm.getExpectedColumns(stmt.Schema)

	// Drop columns that no longer exist in model
	for colName := range currentColumns {
		if _, exists := expectedColumns[colName]; !exists {
			// Skip system columns that shouldn't be dropped
			if mm.isSystemColumn(colName) {
				continue
			}

			log.Info().Str("table", tableName).Str("column", colName).Msg("Dropping obsolete column")
			if err := mm.dropColumn(tableName, colName); err != nil {
				log.Warn().Str("table", tableName).Str("column", colName).Err(err).Msg("Failed to drop column")
			}
		}
	}

	return nil
}

// getCurrentColumns gets current columns from database table
func (mm *MigrationManager) getCurrentColumns(tableName string) (map[string]string, error) {
	var columns []struct {
		ColumnName string `json:"column_name"`
		DataType   string `json:"data_type"`
	}

	err := mm.db.DB.Raw(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = ? AND table_schema = current_schema()
	`, tableName).Scan(&columns).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, col := range columns {
		result[col.ColumnName] = col.DataType
	}

	return result, nil
}

// getExpectedColumns gets expected columns from GORM schema
func (mm *MigrationManager) getExpectedColumns(gormSchema *schema.Schema) map[string]string {
	result := make(map[string]string)

	for _, field := range gormSchema.Fields {
		if field.DBName != "" {
			result[field.DBName] = string(field.DataType)
		}
	}

	return result
}

// isSystemColumn checks if column is a system column that shouldn't be dropped
func (mm *MigrationManager) isSystemColumn(columnName string) bool {
	systemColumns := []string{
		"id", "created_at", "updated_at", "deleted_at",
	}

	columnName = strings.ToLower(columnName)
	for _, sysCol := range systemColumns {
		if columnName == sysCol {
			return true
		}
	}

	return false
}

// dropColumn drops a column from table
func (mm *MigrationManager) dropColumn(tableName, columnName string) error {
	// First, try to drop any constraints on the column
	mm.dropColumnConstraints(tableName, columnName)

	// Then drop the column
	sql := fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s", tableName, columnName)
	return mm.db.DB.Exec(sql).Error
}

// dropColumnConstraints drops constraints associated with a column
func (mm *MigrationManager) dropColumnConstraints(tableName, columnName string) {
	// Get constraint names for this column
	var constraints []string
	mm.db.DB.Raw(`
		SELECT constraint_name
		FROM information_schema.constraint_column_usage
		WHERE table_name = ? AND column_name = ? AND table_schema = current_schema()
	`, tableName, columnName).Scan(&constraints)

	// Drop each constraint
	for _, constraint := range constraints {
		sql := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s", tableName, constraint)
		if err := mm.db.DB.Exec(sql).Error; err != nil {
			log.Warn().Str("constraint", constraint).Err(err).Msg("Failed to drop constraint")
		}
	}
}

// Close closes the database connection
func (mm *MigrationManager) Close() error {
	return mm.db.Close()
}

// GetDB returns the underlying database manager (for advanced usage)
func (mm *MigrationManager) GetDB() *DatabaseManager {
	return mm.db
}
