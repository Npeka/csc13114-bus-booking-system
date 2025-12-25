package server

import (
	"bus-booking/trip-service/internal/client"
	"bus-booking/trip-service/internal/cronjob"
	"bus-booking/trip-service/internal/handler"
	"bus-booking/trip-service/internal/repository"
	"bus-booking/trip-service/internal/router"
	"bus-booking/trip-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() (http.Handler, *cronjob.TripRescheduleCronJob) {
	bookingClient := client.NewBookingClient(s.cfg.ServiceName, s.cfg.External.BookingServiceURL)

	// Initialize repositories
	tripRepo := repository.NewTripRepository(s.db.DB)
	routeRepo := repository.NewRouteRepository(s.db.DB)
	routeStopRepo := repository.NewRouteStopRepository(s.db.DB)
	busRepo := repository.NewBusRepository(s.db.DB)
	seatRepo := repository.NewSeatRepository(s.db.DB)

	// Initialize services
	tripService := service.NewTripService(tripRepo, routeRepo, routeStopRepo, busRepo, seatRepo, bookingClient)
	routeService := service.NewRouteService(routeRepo)
	busService := service.NewBusService(busRepo, seatRepo)
	routeStopService := service.NewRouteStopService(routeStopRepo, routeRepo)
	seatService := service.NewSeatService(seatRepo)
	constantsService := service.NewConstantsService()

	// Initialize trip reschedule cronjob
	cronJob := cronjob.NewTripRescheduleCronJob(tripService)

	// Initialize handlers
	tripHandler := handler.NewTripHandler(tripService)
	routeHandler := handler.NewRouteHandler(routeService)
	busHandler := handler.NewBusHandler(busService)
	routeStopHandler := handler.NewRouteStopHandler(routeStopService)
	seatHandler := handler.NewSeatHandler(seatService)
	constantsHandler := handler.NewConstantsHandler(constantsService)

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
		ConstantsHandler: constantsHandler,
	})
	return engine, cronJob
}
