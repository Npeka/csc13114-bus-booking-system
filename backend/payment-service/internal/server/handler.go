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
	bankAccountRepo := repository.NewBankAccountRepository(s.db.DB)
	refundRepo := repository.NewRefundRepository(s.db.DB) // NEW

	// Initialize PayOS client
	payosClient := service.NewPayOSService(s.cfg.PayOS)
	bookingClient := client.NewBookingClient(s.cfg.ServiceName, s.cfg.External.BookingServiceURL)

	// Initialize constants and Excel services
	constantsService := service.NewConstantsService()
	excelService := service.NewExcelService()

	// Initialize services with PayOS client
	transactionService := service.NewTransactionService(
		transactionRepo,
		bookingClient,
		payosClient,
	)

	bankAccountService := service.NewBankAccountService(
		bankAccountRepo,
		constantsService,
	)

	refundService := service.NewRefundService(
		refundRepo, // NEW - first param
		transactionRepo,
		bankAccountRepo,
		constantsService,
		excelService,
	)

	transactionHandler := handler.NewTransactionHandler(transactionService)
	bankAccountHandler := handler.NewBankAccountHandler(bankAccountService)
	constantsHandler := handler.NewConstantsHandler(constantsService)
	refundHandler := handler.NewRefundHandler(refundService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		TransactionHandler: transactionHandler,
		BankAccountHandler: bankAccountHandler,
		ConstantsHandler:   constantsHandler,
		RefundHandler:      refundHandler,
	})
	return engine
}
