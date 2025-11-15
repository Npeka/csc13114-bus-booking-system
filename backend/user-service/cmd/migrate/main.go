package main

import (
	"log"

	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/db"
	"bus-booking/user-service/internal/model"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	database, err := db.NewPostgresConnection(&cfg.Database, cfg.Server.Environment)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	log.Println("Starting database migration...")

	// Enable UUID extension
	if err := database.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatal("Failed to create UUID extension:", err)
	}

	// Auto migrate the schema
	if err := database.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully!")
}
