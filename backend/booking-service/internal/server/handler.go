package server

import (
	"bus-booking/booking-service/internal/client"
	"bus-booking/booking-service/internal/handler"
	"bus-booking/booking-service/internal/jobs"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/booking-service/internal/router"
	"bus-booking/booking-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) buildHandler() (http.Handler, *jobs.BookingExpirationJob, *jobs.TripReminderJob) {
	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(s.db.DB)
	feedbackRepo := repository.NewFeedbackRepository(s.db.DB)
	bookingStatsRepo := repository.NewBookingStatsRepository(s.db.DB)

	// Initialize HTTP clients for other services
	tripClient := client.NewTripClient(s.cfg.ServiceName, s.cfg.External.TripServiceURL)
	paymentClient := client.NewPaymentClient(s.cfg.ServiceName, s.cfg.External.PaymentServiceURL)
	userClient := client.NewUserClient(s.cfg.ServiceName, s.cfg.External.UserServiceURL)
	notificationClient := client.NewNotificationClient(s.cfg.ServiceName, s.cfg.External.NotificationServiceURL)

	// Initialize services
	seatLockRepo := repository.NewSeatLockRepository(s.db.DB)

	bookingService := service.NewBookingService(bookingRepo, paymentClient, tripClient, userClient, notificationClient, s.delayedQueue)
	feedbackService := service.NewFeedbackService(bookingRepo, feedbackRepo)
	statisticsService := service.NewStatisticsService(bookingStatsRepo)
	eTicketService := service.NewETicketService(bookingRepo, tripClient)

	// Initialize Jobs
	bookingExpirationJob := jobs.NewBookingExpirationJob(bookingService, seatLockRepo, s.delayedQueue)
	tripReminderJob := jobs.NewTripReminderJob(bookingRepo, s.delayedQueue, notificationClient, tripClient, userClient)

	// Initialize handlers
	bookingHandler := handler.NewBookingHandler(bookingService, eTicketService)
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
	return engine, bookingExpirationJob, tripReminderJob
}
