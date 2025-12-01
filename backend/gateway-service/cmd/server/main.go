package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"bus-booking/gateway-service/config"
	"bus-booking/gateway-service/internal/proxy"
	"bus-booking/shared/logger"
	"bus-booking/shared/middleware"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.MustLoadConfig()
	routes := config.MustLoadRoutes()
	logger.MustSetupLogger(&cfg.Log)

	// Debug: Print loaded services configuration
	log.Info().Msg("Loaded services configuration:")
	for serviceName, serviceConfig := range cfg.ServicesMap {
		log.Info().Msgf("Service %s: URL=%s, Timeout=%d, Retries=%d", serviceName, serviceConfig.URL, serviceConfig.Timeout, serviceConfig.Retries)
	}

	// Debug: Print loaded routes
	log.Info().Msgf("Loaded %d routes:", len(routes.Routes))
	for i, route := range routes.Routes {
		log.Info().Msgf("Route %d: %s %v -> %s", i+1, route.Path, route.Methods, route.Service)
	}

	// Create gateway
	gateway := proxy.NewGateway(cfg, routes)

	// Setup Gin
	if gin.Mode() == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add global middleware
	router.Use(middleware.RequestContextMiddleware("gateway-service"))
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	gateway.SetupRoutes(router)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Gateway server starting on %s", addr)
	log.Printf("Loaded %d routes", len(routes.Routes))

	if err := router.Run(addr); err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}
}
