package main

// @title Booking Service API
// @version 1.0
// @description API for booking management in the bus booking system
// @description This service handles bookings, seat reservations, and feedback.
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
	"bus-booking/booking-service/config"
	_ "bus-booking/booking-service/docs"
	"bus-booking/booking-service/internal/server"
	"bus-booking/shared/db"
	"bus-booking/shared/logger"
	"bus-booking/shared/validator"
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
