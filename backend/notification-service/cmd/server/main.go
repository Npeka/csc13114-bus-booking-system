package main

// @title Payment Service API
// @version 1.0
// @description API for payment processing in the bus booking system
// @description This service handles payment transactions, refunds, and payment history.
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
	"bus-booking/notification-service/config"
	_ "bus-booking/notification-service/docs"
	"bus-booking/notification-service/internal/server"
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
