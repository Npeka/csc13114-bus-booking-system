package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"bus-booking/shared/middleware"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/handler"
)

// RouterConfig holds router dependencies
type RouterConfig struct {
	TripHandler  *handler.TripHandler
	RouteHandler *handler.RouteHandler
	BusHandler   *handler.BusHandler
	ServiceName  string
	Config       *config.Config
}

// SetupRoutes configures all routes for the trip service
func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	// Apply global middleware
	router.Use(middleware.RequestContextMiddleware(config.ServiceName))
	router.Use(setupLocalCORS(&config.Config.CORS))
	router.Use(middleware.Logger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": config.ServiceName,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Trip routes
		trips := v1.Group("/trips")
		{
			trips.POST("/search", config.TripHandler.SearchTrips)
			trips.GET("/:id", config.TripHandler.GetTrip)
			trips.POST("", config.TripHandler.CreateTrip)
			trips.PUT("/:id", config.TripHandler.UpdateTrip)
			trips.DELETE("/:id", config.TripHandler.DeleteTrip)
			trips.GET("/route/:route_id", config.TripHandler.ListTripsByRoute)
		}

		// Route routes
		routes := v1.Group("/routes")
		{
			routes.POST("", config.RouteHandler.CreateRoute)
			routes.GET("/:id", config.RouteHandler.GetRoute)
			routes.PUT("/:id", config.RouteHandler.UpdateRoute)
			routes.DELETE("/:id", config.RouteHandler.DeleteRoute)
			routes.GET("", config.RouteHandler.ListRoutes)
			routes.GET("/search", config.RouteHandler.SearchRoutes)
		}

		// Bus routes
		buses := v1.Group("/buses")
		{
			buses.POST("", config.BusHandler.CreateBus)
			buses.GET("/:id", config.BusHandler.GetBus)
			buses.PUT("/:id", config.BusHandler.UpdateBus)
			buses.DELETE("/:id", config.BusHandler.DeleteBus)
			buses.GET("", config.BusHandler.ListBuses)
			buses.GET("/:id/seats", config.BusHandler.GetBusSeats)
		}
	}
}

// setupLocalCORS creates CORS middleware from local config
func setupLocalCORS(cfg *config.CORSConfig) gin.HandlerFunc {
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
