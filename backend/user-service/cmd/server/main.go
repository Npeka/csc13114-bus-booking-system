package main

// @title User Service API
// @version 1.0
// @description API for user management and authentication in the bus booking system
// @description This service handles user registration, authentication, profile management, and authorization.
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
	"bus-booking/user-service/config"
	_ "bus-booking/user-service/docs"
	"bus-booking/user-service/internal/initializer"
	"bus-booking/user-service/internal/server"
)

func main() {
	cfg := config.MustLoadConfig()
	logger.MustSetupLogger(&cfg.Log)
	validator.MustSetupValidator()

	pg := db.MustNewPostgresConnection(&cfg.Database)
	rd := db.MustNewRedisConnection(&cfg.Redis)
	fa := initializer.MustNewFirebase(&cfg.Firebase)

	sv := server.NewServer(cfg, pg, rd, fa)
	defer sv.Close()

	sv.Run()
}
