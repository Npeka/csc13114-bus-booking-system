package server

import (
	"bus-booking/booking-service/internal/client"
	"bus-booking/booking-service/internal/handler"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/booking-service/internal/router"
	"bus-booking/booking-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(s.db.DB)
	feedbackRepo := repository.NewFeedbackRepository(s.db.DB)
	bookingStatsRepo := repository.NewBookingStatsRepository(s.db.DB)

	// Initialize HTTP clients for other services
	tripClient := client.NewTripClient(s.cfg.ServiceName, s.cfg.External.TripServiceURL)
	paymentClient := client.NewPaymentClient(s.cfg.ServiceName, s.cfg.External.PaymentServiceURL)

	// Initialize services
	bookingService := service.NewBookingService(bookingRepo, paymentClient, tripClient)
	feedbackService := service.NewFeedbackService(bookingRepo, feedbackRepo)
	statisticsService := service.NewStatisticsService(bookingStatsRepo)

	// Initialize handlers
	bookingHandler := handler.NewBookingHandler(bookingService)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		BookingHandler:    bookingHandler,
		FeedbackHandler:   feedbackHandler,
		StatisticsHandler: statisticsHandler,
	})
	return engine
}
