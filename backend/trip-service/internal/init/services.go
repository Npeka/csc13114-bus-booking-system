package initializer

import (
	sharedDB "bus-booking/shared/db"
	"bus-booking/trip-service/config"
	"bus-booking/trip-service/internal/handler"
	"bus-booking/trip-service/internal/repository"
	"bus-booking/trip-service/internal/service"

	"github.com/rs/zerolog/log"
)

type ServiceDependencies struct {
	TripHandler  *handler.TripHandler
	RouteHandler *handler.RouteHandler
	BusHandler   *handler.BusHandler
	Repositories *repository.Repositories
}

func InitServices(cfg *config.Config, database *sharedDB.DatabaseManager) *ServiceDependencies {
	repositories := repository.NewRepositories(database.DB)

	tripService := service.NewTripService(repositories)
	routeService := service.NewRouteService(repositories)
	busService := service.NewBusService(repositories)

	tripHandler := handler.NewTripHandler(tripService)
	routeHandler := handler.NewRouteHandler(routeService)
	busHandler := handler.NewBusHandler(busService)

	log.Info().Msg("All services and handlers initialized successfully")

	return &ServiceDependencies{
		TripHandler:  tripHandler,
		RouteHandler: routeHandler,
		BusHandler:   busHandler,
		Repositories: repositories,
	}
}
