package server

import (
	"bus-booking/payment-service/internal/handler"
	"bus-booking/payment-service/internal/repository"
	"bus-booking/payment-service/internal/router"
	"bus-booking/payment-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	repositories := repository.NewRepositories(s.db.DB)

	// Initialize PayOS client
	payosClient := service.NewPayOSClient(
		s.cfg.PayOS.ClientID,
		s.cfg.PayOS.APIKey,
		s.cfg.PayOS.ChecksumKey,
	)

	// Initialize services with PayOS client
	transactionService := service.NewTransactionService(
		repositories,
		payosClient,
		s.cfg.PayOS.ReturnURL,
		s.cfg.PayOS.CancelURL,
	)

	transactionHandler := handler.NewTransactionHandler(transactionService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		TransactionHandler: transactionHandler,
	})
	return engine
}
