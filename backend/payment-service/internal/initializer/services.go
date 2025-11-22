package initializer

import (
	"bus-booking/payment-service/config"
	"bus-booking/payment-service/internal/handler"
	"bus-booking/payment-service/internal/repository"
	"bus-booking/payment-service/internal/service"
	sharedDB "bus-booking/shared/db"

	"github.com/rs/zerolog/log"
)

type ServiceDependencies struct {
	TransactionHandler handler.TransactionHandler
	Repositories       *repository.Repositories
}

func InitServices(cfg *config.Config, database *sharedDB.DatabaseManager) *ServiceDependencies {
	repositories := repository.NewRepositories(database.DB)

	transactionService := service.NewTransactionService(repositories)

	transactionHandler := handler.NewTransactionHandler(transactionService)

	log.Info().Msg("All services and handlers initialized successfully")

	return &ServiceDependencies{
		TransactionHandler: transactionHandler,
		Repositories:       repositories,
	}
}
