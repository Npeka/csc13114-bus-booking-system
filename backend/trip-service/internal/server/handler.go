package server

import (
	"bus-booking/trip-service/internal/handler"
	"bus-booking/trip-service/internal/repository"
	"bus-booking/trip-service/internal/router"
	"bus-booking/trip-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	repositories := repository.NewRepositories(s.db.DB)

	tripService := service.NewTripService(repositories)
	routeService := service.NewRouteService(repositories)
	busService := service.NewBusService(repositories)

	tripHandler := handler.NewTripHandler(tripService)
	routeHandler := handler.NewRouteHandler(routeService)
	busHandler := handler.NewBusHandler(busService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		TripHandler:  tripHandler,
		RouteHandler: routeHandler,
		BusHandler:   busHandler,
	})
	return engine
}
