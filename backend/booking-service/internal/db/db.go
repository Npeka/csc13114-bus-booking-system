package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/model"
)

// DB wraps the GORM database connection
type DB struct {
	*gorm.DB
}

// NewDB creates a new database connection
func NewDB(cfg *config.Config) (*DB, error) {
	dbLogger := logger.Default
	if cfg.IsDevelopment() {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger:          dbLogger,
		NowFunc:         func() time.Time { return time.Now().UTC() },
		CreateBatchSize: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	log.Println("Running database migrations...")

	// Create UUID extension if not exists
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	// Auto migrate all models
	if err := db.AutoMigrate(
		&model.Booking{},
		&model.BookingSeat{},
		&model.SeatStatus{},
		&model.PaymentMethod{},
		&model.Feedback{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Seed runs database seeders
func (db *DB) Seed() error {
	log.Println("Running database seeders...")

	// Seed default payment methods
	if err := db.seedPaymentMethods(); err != nil {
		return fmt.Errorf("failed to seed payment methods: %w", err)
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// seedPaymentMethods seeds default payment methods
func (db *DB) seedPaymentMethods() error {
	paymentMethods := []model.PaymentMethod{
		{
			Name:        "Cash",
			Code:        "CASH",
			Description: "Pay with cash at pickup location",
			IsActive:    true,
		},
		{
			Name:        "Credit Card",
			Code:        "CREDIT_CARD",
			Description: "Pay with credit/debit card",
			IsActive:    true,
		},
		{
			Name:        "Bank Transfer",
			Code:        "BANK_TRANSFER",
			Description: "Pay via bank transfer",
			IsActive:    true,
		},
		{
			Name:        "Digital Wallet",
			Code:        "DIGITAL_WALLET",
			Description: "Pay with digital wallet (Momo, ZaloPay, etc.)",
			IsActive:    true,
		},
	}

	for _, method := range paymentMethods {
		var existing model.PaymentMethod
		if err := db.Where("code = ?", method.Code).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&method).Error; err != nil {
					return fmt.Errorf("failed to create payment method %s: %w", method.Code, err)
				}
				log.Printf("Created payment method: %s", method.Name)
			} else {
				return fmt.Errorf("failed to check payment method %s: %w", method.Code, err)
			}
		}
	}

	return nil
}

// IsHealthy checks if the database connection is healthy
func (db *DB) IsHealthy() bool {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
}
