package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"bus-booking/gateway-service/config"
	"bus-booking/gateway-service/internal/proxy"
	"bus-booking/shared/middleware"
)

func main() {
	cfg := config.MustLoadConfig("config/config.yaml")
	routes := config.MustLoadRoutes("routes")

	// Debug: Print loaded services configuration
	log.Printf("Loaded services configuration:")
	for serviceName, serviceConfig := range cfg.ServicesMap {
		log.Printf("Service %s: URL=%s, Timeout=%d, Retries=%d", serviceName, serviceConfig.URL, serviceConfig.Timeout, serviceConfig.Retries)
	}

	// Debug: Print loaded routes
	log.Printf("Loaded %d routes:", len(routes.Routes))
	for i, route := range routes.Routes {
		log.Printf("Route %d: %s %v -> %s", i+1, route.Path, route.Methods, route.Service)
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
		log.Fatalf("Failed to start server: %v", err)
	}
}
