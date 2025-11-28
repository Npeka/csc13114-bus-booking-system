package main

// @title Chatbot Service API
// @version 1.0
// @description AI-powered chatbot service for bus ticket booking assistance
// @description This service provides natural language processing for trip search, FAQ handling, and booking assistance.
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
	"bus-booking/chatbot-service/config"
	_ "bus-booking/chatbot-service/docs"
	"bus-booking/chatbot-service/internal/server"
	"bus-booking/shared/logger"
	"bus-booking/shared/validator"
)

func main() {
	cfg := config.MustLoadConfig()
	logger.MustSetupLogger(&cfg.Log)
	validator.MustSetupValidator()

	sv := server.NewServer(cfg)
	defer sv.Close()

	sv.Run()
}
