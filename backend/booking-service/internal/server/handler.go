package server

import (
	"bus-booking/booking-service/internal/handler"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/booking-service/internal/router"
	"bus-booking/booking-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() http.Handler {
	repositories := repository.NewRepositories(s.db.DB)

	bookingService := service.NewBookingService(repositories)
	seatLockService := service.NewSeatLockService(repositories.SeatLock)
	paymentService := service.NewPaymentService(repositories.PaymentMethod, repositories.Booking)
	feedbackService := service.NewFeedbackService(repositories)
	statisticsService := service.NewStatisticsService(repositories)
	seatService := service.NewSeatService(repositories)

	bookingHandler := handler.NewBookingHandler(bookingService)
	seatLockHandler := handler.NewSeatLockHandler(seatLockService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)
	seatHandler := handler.NewSeatHandler(seatService)

	if s.cfg.Server.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	router.SetupRoutes(engine, s.cfg, &router.Handlers{
		BookingHandler:    bookingHandler,
		SeatLockHandler:   seatLockHandler,
		PaymentHandler:    paymentHandler,
		FeedbackHandler:   feedbackHandler,
		StatisticsHandler: statisticsHandler,
		SeatHandler:       seatHandler,
	})
	return engine
}
