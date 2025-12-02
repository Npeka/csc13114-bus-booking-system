package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/handler"
)

type Handlers struct {
	TripHandler      handler.TripHandler
	RouteHandler     handler.RouteHandler
	RouteStopHandler handler.RouteStopHandler
	BusHandler       handler.BusHandler
	SeatHandler      handler.SeatHandler
	ConstantsHandler handler.ConstantsHandler
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *Handlers) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.RequestContextMiddleware(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/constants", ginext.WrapHandler(h.ConstantsHandler.GetConstants))

		trips := v1.Group("/trips")
		{
			trips.GET("/search", ginext.WrapHandler(h.TripHandler.SearchTrips))
			trips.GET("/:id", ginext.WrapHandler(h.TripHandler.GetTrip))
			// trips.GET("/route/:route_id", ginext.WrapHandler(h.TripHandler.ListTripsByRoute))
		}

		// buses := v1.Group("/buses")
		// {
		// 	buses.GET("/:id/seats", ginext.WrapHandler(h.BusHandler.GetBusSeats))
		// 	buses.GET("/:id/seat-map", ginext.WrapHandler(h.SeatHandler.GetSeatMap))
		// }

		// seats := v1.Group("/buses/seats")
		// {
		// 	seats.POST("/bulk", ginext.WrapHandler(h.SeatHandler.CreateSeatsFromTemplate))
		// }

		// routes := v1.Group("/routes")
		// {
		// 	routes.GET("/search", ginext.WrapHandler(h.RouteHandler.SearchRoutes))
		// 	routes.GET("/:id/stops", ginext.WrapHandler(h.RouteStopHandler.ListRouteStops))
		// 	routes.POST("/:id/stops/reorder", ginext.WrapHandler(h.RouteStopHandler.ReorderStops))
		// }
	}

	adminV1 := router.Group("/api/v1")
	adminV1.Use(middleware.RequireAuthMiddleware())
	adminV1.Use(middleware.RequireRoleMiddleware(constants.RoleAdmin))
	{
		trips := adminV1.Group("/trips")
		{
			trips.GET("", ginext.WrapHandler(h.TripHandler.ListTrips))
			trips.POST("", ginext.WrapHandler(h.TripHandler.CreateTrip))
			trips.PUT("/:id", ginext.WrapHandler(h.TripHandler.UpdateTrip))
			trips.DELETE("/:id", ginext.WrapHandler(h.TripHandler.DeleteTrip))
		}

		buses := adminV1.Group("/buses")
		{
			buses.GET("", ginext.WrapHandler(h.BusHandler.ListBuses))
			buses.GET("/:id", ginext.WrapHandler(h.BusHandler.GetBus))
			buses.POST("", ginext.WrapHandler(h.BusHandler.CreateBus))
			buses.PUT("/:id", ginext.WrapHandler(h.BusHandler.UpdateBus))
			buses.DELETE("/:id", ginext.WrapHandler(h.BusHandler.DeleteBus))
		}

		seats := v1.Group("/buses/seats")
		{
			seats.POST("", ginext.WrapHandler(h.SeatHandler.CreateSeat))
			seats.PUT("/:id", ginext.WrapHandler(h.SeatHandler.UpdateSeat))
			seats.DELETE("/:id", ginext.WrapHandler(h.SeatHandler.DeleteSeat))
		}

		routes := adminV1.Group("/routes")
		{
			routes.GET("", ginext.WrapHandler(h.RouteHandler.ListRoutes))
			routes.GET("/:id", ginext.WrapHandler(h.RouteHandler.GetRoute))
			routes.POST("", ginext.WrapHandler(h.RouteHandler.CreateRoute))
			routes.PUT("/:id", ginext.WrapHandler(h.RouteHandler.UpdateRoute))
			routes.DELETE("/:id", ginext.WrapHandler(h.RouteHandler.DeleteRoute))
		}

		routeStops := adminV1.Group("/routes/stops")
		{
			routeStops.POST("", ginext.WrapHandler(h.RouteStopHandler.CreateRouteStop))
			routeStops.PUT("/:id", ginext.WrapHandler(h.RouteStopHandler.UpdateRouteStop))
			routeStops.DELETE("/:id", ginext.WrapHandler(h.RouteStopHandler.DeleteRouteStop))
		}
	}
}
