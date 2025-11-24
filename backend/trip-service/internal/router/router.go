package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/handler"
)

type Handlers struct {
	TripHandler  handler.TripHandler
	RouteHandler handler.RouteHandler
	BusHandler   handler.BusHandler
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *Handlers) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.RequestContextMiddleware(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		trips := v1.Group("/trips")
		{
			trips.POST("/search", h.TripHandler.SearchTrips)
			trips.GET("/:id", h.TripHandler.GetTrip)
			trips.POST("", h.TripHandler.CreateTrip)
			trips.PUT("/:id", h.TripHandler.UpdateTrip)
			trips.DELETE("/:id", h.TripHandler.DeleteTrip)
			trips.GET("/route/:route_id", h.TripHandler.ListTripsByRoute)
		}

		routes := v1.Group("/routes")
		{
			routes.POST("", h.RouteHandler.CreateRoute)
			routes.GET("/:id", h.RouteHandler.GetRoute)
			routes.PUT("/:id", h.RouteHandler.UpdateRoute)
			routes.DELETE("/:id", h.RouteHandler.DeleteRoute)
			routes.GET("", h.RouteHandler.ListRoutes)
			routes.GET("/search", h.RouteHandler.SearchRoutes)
		}

		buses := v1.Group("/buses")
		{
			buses.POST("", h.BusHandler.CreateBus)
			buses.GET("/:id", h.BusHandler.GetBus)
			buses.PUT("/:id", h.BusHandler.UpdateBus)
			buses.DELETE("/:id", h.BusHandler.DeleteBus)
			buses.GET("", h.BusHandler.ListBuses)
			buses.GET("/:id/seats", h.BusHandler.GetBusSeats)
		}
	}
}
