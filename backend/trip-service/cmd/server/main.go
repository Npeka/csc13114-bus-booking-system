package main

// @title Trip Service API
// @version 1.0
// @description API for managing trips, routes, buses, and operators in the bus booking system
//
// @contact.name API Support
// @contact.email support@busbooking.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"bus-booking/shared/db"
	"bus-booking/shared/logger"
	"bus-booking/shared/validator"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/server"
)

func main() {
	cfg := config.MustLoadConfig()
	logger.MustSetupLogger(&cfg.Log)
	validator.MustSetupValidator()

	pg := db.MustNewPostgresConnection(&cfg.Database)
	sv := server.NewServer(cfg, pg)
	defer sv.Close()

	sv.Run()
}
