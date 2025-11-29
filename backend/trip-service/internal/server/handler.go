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
	tripRepo := repository.NewTripRepository(s.db.DB)
	routeRepo := repository.NewRouteRepository(s.db.DB)
	routeStopRepo := repository.NewRouteStopRepository(s.db.DB)
	busRepo := repository.NewBusRepository(s.db.DB)
	seatRepo := repository.NewSeatRepository(s.db.DB)

	tripService := service.NewTripService(tripRepo, routeRepo, routeStopRepo, busRepo, seatRepo)
	routeService := service.NewRouteService(routeRepo)
	busService := service.NewBusService(busRepo, seatRepo)
	routeStopService := service.NewRouteStopService(routeStopRepo, routeRepo)
	seatService := service.NewSeatService(seatRepo, busRepo)

	tripHandler := handler.NewTripHandler(tripService)
	routeHandler := handler.NewRouteHandler(routeService)
	busHandler := handler.NewBusHandler(busService)
	routeStopHandler := handler.NewRouteStopHandler(routeStopService)
	seatHandler := handler.NewSeatHandler(seatService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		TripHandler:      tripHandler,
		RouteHandler:     routeHandler,
		BusHandler:       busHandler,
		RouteStopHandler: routeStopHandler,
		SeatHandler:      seatHandler,
	})
	return engine
}
