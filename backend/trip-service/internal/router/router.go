package router

import (
	"github.com/gin-gonic/gin"

	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/handler"
)

type RouterConfig struct {
	Config       *config.Config
	TripHandler  *handler.TripHandler
	RouteHandler *handler.RouteHandler
	BusHandler   *handler.BusHandler
}

func SetupRoutes(router *gin.Engine, cfg *RouterConfig) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.Config.CORS))
	router.Use(middleware.RequestContextMiddleware(cfg.Config.ServiceName))
	router.GET(health.Path, health.Handler(cfg.Config.ServiceName))

	v1 := router.Group("/api/v1")
	{
		trips := v1.Group("/trips")
		{
			trips.POST("/search", cfg.TripHandler.SearchTrips)
			trips.GET("/:id", cfg.TripHandler.GetTrip)
			trips.POST("", cfg.TripHandler.CreateTrip)
			trips.PUT("/:id", cfg.TripHandler.UpdateTrip)
			trips.DELETE("/:id", cfg.TripHandler.DeleteTrip)
			trips.GET("/route/:route_id", cfg.TripHandler.ListTripsByRoute)
		}

		routes := v1.Group("/routes")
		{
			routes.POST("", cfg.RouteHandler.CreateRoute)
			routes.GET("/:id", cfg.RouteHandler.GetRoute)
			routes.PUT("/:id", cfg.RouteHandler.UpdateRoute)
			routes.DELETE("/:id", cfg.RouteHandler.DeleteRoute)
			routes.GET("", cfg.RouteHandler.ListRoutes)
			routes.GET("/search", cfg.RouteHandler.SearchRoutes)
		}

		buses := v1.Group("/buses")
		{
			buses.POST("", cfg.BusHandler.CreateBus)
			buses.GET("/:id", cfg.BusHandler.GetBus)
			buses.PUT("/:id", cfg.BusHandler.UpdateBus)
			buses.DELETE("/:id", cfg.BusHandler.DeleteBus)
			buses.GET("", cfg.BusHandler.ListBuses)
			buses.GET("/:id/seats", cfg.BusHandler.GetBusSeats)
		}
	}
}
