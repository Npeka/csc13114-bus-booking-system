package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"bus-booking/gateway-service/internal/config"
	"bus-booking/gateway-service/internal/proxy"
	"bus-booking/shared/middleware"
)

func main() {
	// Load configuration
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
	router.Use(setupCORS(&cfg.CORS))
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	gateway.SetupRoutes(router)

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Gateway server starting on %s", addr)
	log.Printf("Loaded %d routes", len(routes.Routes))

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupCORS(cfg *config.CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		ExposeHeaders:    cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
	}

	// If allow origins contains "*", use AllowAllOrigins
	for _, origin := range cfg.AllowOrigins {
		if origin == "*" {
			corsConfig.AllowAllOrigins = true
			corsConfig.AllowOrigins = nil
			break
		}
	}

	return cors.New(corsConfig)
}
