package appinit

import (
	sharedDB "bus-booking/shared/db"
	"bus-booking/trip-service/internal/handler"
	"bus-booking/trip-service/internal/repository"
	"bus-booking/trip-service/internal/service"

	"github.com/rs/zerolog/log"
)

type ServiceHandlers struct {
	TripHandler  *handler.TripHandler
	RouteHandler *handler.RouteHandler
	BusHandler   *handler.BusHandler
}

func InitServices(database *sharedDB.DatabaseManager) *ServiceHandlers {
	// Initialize repository
	tripRepo := repository.NewTripRepository(database.DB)

	// Initialize services
	tripService := service.NewTripService(tripRepo)
	routeService := service.NewRouteService(tripRepo)
	busService := service.NewBusService(tripRepo)

	// Initialize handlers
	tripHandler := handler.NewTripHandler(tripService)
	routeHandler := handler.NewRouteHandler(routeService)
	busHandler := handler.NewBusHandler(busService)

	log.Info().Msg("All services and handlers initialized successfully")

	return &ServiceHandlers{
		TripHandler:  tripHandler,
		RouteHandler: routeHandler,
		BusHandler:   busHandler,
	}
}
