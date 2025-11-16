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
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Printf("Warning: Failed to load config file: %v. Using defaults.", err)
		cfg, _ = config.LoadConfig("")
	}

	// Load routes
	routes, err := config.LoadRoutes("routes")
	if err != nil {
		log.Fatalf("Failed to load routes: %v", err)
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
