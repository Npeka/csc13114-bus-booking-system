package router

import (
	"github.com/gin-gonic/gin"

	"bus-booking/shared/middleware"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/handler"
)

type RouterConfig struct {
	ServiceName  string
	Config       *config.Config
	TripHandler  *handler.TripHandler
	RouteHandler *handler.RouteHandler
	BusHandler   *handler.BusHandler
}

func SetupRoutes(router *gin.Engine, config *RouterConfig) {
	router.Use(middleware.RequestContextMiddleware(config.ServiceName))
	router.Use(middleware.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": config.ServiceName,
		})
	})

	v1 := router.Group("/api/v1")
	{
		trips := v1.Group("/trips")
		{
			trips.POST("/search", config.TripHandler.SearchTrips)
			trips.GET("/:id", config.TripHandler.GetTrip)
			trips.POST("", config.TripHandler.CreateTrip)
			trips.PUT("/:id", config.TripHandler.UpdateTrip)
			trips.DELETE("/:id", config.TripHandler.DeleteTrip)
			trips.GET("/route/:route_id", config.TripHandler.ListTripsByRoute)
		}

		routes := v1.Group("/routes")
		{
			routes.POST("", config.RouteHandler.CreateRoute)
			routes.GET("/:id", config.RouteHandler.GetRoute)
			routes.PUT("/:id", config.RouteHandler.UpdateRoute)
			routes.DELETE("/:id", config.RouteHandler.DeleteRoute)
			routes.GET("", config.RouteHandler.ListRoutes)
			routes.GET("/search", config.RouteHandler.SearchRoutes)
		}

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
