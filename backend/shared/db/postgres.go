package db

import (
	"context"
	"database/sql"
	"fmt"
	stdlog "log"
	"os"
	"time"

	"bus-booking/shared/config"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager handles database connections and operations
type DatabaseManager struct {
	DB     *gorm.DB
	SqlDB  *sql.DB
	Config *config.DatabaseConfig
}

// NewPostgresConnection creates a new PostgreSQL connection with GORM
func NewPostgresConnection(cfg *config.DatabaseConfig, env string) (*DatabaseManager, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode, cfg.TimeZone,
	)

	// Configure GORM logger based on environment
	var logLevel logger.LogLevel
	switch env {
	case "development":
		logLevel = logger.Info
	case "staging":
		logLevel = logger.Warn
	case "production":
		logLevel = logger.Error
	default:
		logLevel = logger.Silent
	}

	gormLogger := logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  env == "development",
			ParameterizedQueries:      env == "production", // Hide SQL parameters in production
		},
	)

	// GORM configuration
	gormConfig := &gorm.Config{
		Logger:                                   gormLogger,
		NowFunc:                                  func() time.Time { return time.Now().UTC() },
		PrepareStmt:                              true, // Enable prepared statements for better performance
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Name).
		Str("ssl_mode", cfg.SSLMode).
		Msg("Successfully connected to PostgreSQL database")

	return &DatabaseManager{
		DB:     db,
		SqlDB:  sqlDB,
		Config: cfg,
	}, nil
}

// Close closes the database connection
func (dm *DatabaseManager) Close() error {
	if dm.SqlDB != nil {
		return dm.SqlDB.Close()
	}
	return nil
}

// HealthCheck performs a health check on the database
func (dm *DatabaseManager) HealthCheck() error {
	if dm.SqlDB == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return dm.SqlDB.PingContext(ctx)
}

// GetStats returns database connection statistics
func (dm *DatabaseManager) GetStats() sql.DBStats {
	if dm.SqlDB == nil {
		return sql.DBStats{}
	}
	return dm.SqlDB.Stats()
}

// AutoMigrate runs database migrations for given models
func (dm *DatabaseManager) AutoMigrate(models ...interface{}) error {
	if err := dm.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	log.Info().Msg("Database migrations completed successfully")
	return nil
}
