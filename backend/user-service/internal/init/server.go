package appinit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/handler"
	"bus-booking/user-service/internal/router"
)

// InitHTTPServer configures the HTTP server and routes
func InitHTTPServer(cfg *config.Config, userHandler *handler.UserHandler, authHandler *handler.AuthHandler, serviceName string) *http.Server {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	ginRouter := gin.New()

	// Setup routes using the router package
	routerConfig := &router.RouterConfig{
		UserHandler: userHandler,
		AuthHandler: authHandler,
		ServiceName: serviceName,
	}
	router.SetupRoutes(ginRouter, routerConfig)

	// Create HTTP server
	server := &http.Server{
		Addr:           cfg.GetServerAddr(),
		Handler:        ginRouter,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	log.Info().Str("address", cfg.GetServerAddr()).Msg("HTTP server configured")
	return server
}
