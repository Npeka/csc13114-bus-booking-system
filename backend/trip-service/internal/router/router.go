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
	router.Use(middleware.RequestContext(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/constants", ginext.WrapHandler(h.ConstantsHandler.GetList))

		trips := v1.Group("/trips")
		{
			trips.GET("/search", ginext.WrapHandler(h.TripHandler.SearchTrips))
			trips.GET("/:id", ginext.WrapHandler(h.TripHandler.GetByID))
		}

		buses := v1.Group("/buses")
		{
			buses.GET("/:id", ginext.WrapHandler(h.BusHandler.Get))
		}
	}

	adminV1 := router.Group("/api/v1")
	adminV1.Use(middleware.RequireAuth())
	adminV1.Use(middleware.RequireRole(constants.RoleAdmin))
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
			buses.GET("", ginext.WrapHandler(h.BusHandler.GetList))
			buses.POST("", ginext.WrapHandler(h.BusHandler.Create))
			buses.PUT("/:id", ginext.WrapHandler(h.BusHandler.Update))
			buses.DELETE("/:id", ginext.WrapHandler(h.BusHandler.Delete))
			buses.POST("/:id/images", ginext.WrapHandler(h.BusHandler.UploadImages))
			buses.DELETE("/:id/images", ginext.WrapHandler(h.BusHandler.DeleteImage))
		}

		seats := adminV1.Group("/buses/seats")
		{
			seats.PUT("/:id", ginext.WrapHandler(h.SeatHandler.Update))
		}

		routes := adminV1.Group("/routes")
		{
			routes.GET("", ginext.WrapHandler(h.RouteHandler.GetList))
			routes.GET("/:id", ginext.WrapHandler(h.RouteHandler.GetByID))
			routes.POST("", ginext.WrapHandler(h.RouteHandler.Create))
			routes.PUT("/:id", ginext.WrapHandler(h.RouteHandler.Update))
			routes.DELETE("/:id", ginext.WrapHandler(h.RouteHandler.Delete))
		}

		routeStops := adminV1.Group("/routes/stops")
		{
			routeStops.POST("", ginext.WrapHandler(h.RouteStopHandler.CreateRouteStop))
			routeStops.PUT("/:id", ginext.WrapHandler(h.RouteStopHandler.UpdateRouteStop))
			routeStops.POST("/:id/move", ginext.WrapHandler(h.RouteStopHandler.MoveRouteStop))
			routeStops.DELETE("/:id", ginext.WrapHandler(h.RouteStopHandler.DeleteRouteStop))
		}

	}

	internalV1 := router.Group("/api/v1")
	{
		seats := internalV1.Group("/buses/seats")
		{
			seats.GET("/ids", ginext.WrapHandler(h.SeatHandler.GetListByIDs))
		}
	}
}
