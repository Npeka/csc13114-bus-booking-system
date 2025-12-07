package server

import (
	"bus-booking/payment-service/internal/client"
	"bus-booking/payment-service/internal/handler"
	"bus-booking/payment-service/internal/repository"
	"bus-booking/payment-service/internal/router"
	"bus-booking/payment-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	transactionRepo := repository.NewTransactionRepository(s.db.DB)

	// Initialize PayOS client
	payosClient := service.NewPayOSClient(s.cfg.PayOS)
	bookingClient := client.NewBookingClient(s.cfg.ServiceName, s.cfg.External.BookingServiceURL)

	// Initialize services with PayOS client
	transactionService := service.NewTransactionService(
		transactionRepo,
		bookingClient,
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
