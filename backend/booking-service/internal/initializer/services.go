package initializer

import (
	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/handler"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/booking-service/internal/service"
	sharedDB "bus-booking/shared/db"

	"github.com/rs/zerolog/log"
)

type ServiceDependencies struct {
	BookingHandler *handler.BookingHandler
	Repositories   *repository.Repositories
}

func InitServices(cfg *config.Config, database *sharedDB.DatabaseManager) *ServiceDependencies {
	repositories := repository.NewRepositories(database.DB)

	bookingService := service.NewBookingService(repositories)

	bookingHandler := handler.NewBookingHandler(bookingService)

	log.Info().Msg("All services and handlers initialized successfully")

	return &ServiceDependencies{
		BookingHandler: bookingHandler,
		Repositories:   repositories,
	}
}
