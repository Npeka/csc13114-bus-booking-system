package router

import (
	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/handler"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	BookingHandler    handler.BookingHandler
	PaymentHandler    handler.PaymentHandler
	FeedbackHandler   handler.FeedbackHandler
	StatisticsHandler handler.StatisticsHandler
	SeatHandler       handler.SeatHandler
	SeatLockHandler   handler.SeatLockHandler
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *Handlers) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.RequestContextMiddleware(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		bookings := v1.Group("/bookings")
		{
			bookings.POST("", ginext.WrapHandler(h.BookingHandler.CreateBooking))
			bookings.GET("/:id", ginext.WrapHandler(h.BookingHandler.GetBooking))
			bookings.POST("/:id/cancel", ginext.WrapHandler(h.BookingHandler.CancelBooking))
			bookings.PUT("/:id/status", ginext.WrapHandler(h.BookingHandler.UpdateBookingStatus))
			bookings.GET("/user/:user_id", ginext.WrapHandler(h.BookingHandler.GetUserBookings))
			bookings.GET("/trip/:trip_id", ginext.WrapHandler(h.BookingHandler.GetTripBookings))
		}

		// Trips group for trip-related endpoints
		trips := v1.Group("/trips")
		{
			trips.GET("/:trip_id/seats", ginext.WrapHandler(h.SeatHandler.GetSeatAvailability))
			trips.GET("/:trip_id/locked-seats", ginext.WrapHandler(h.SeatLockHandler.GetLockedSeats))
		}

		// Payment endpoints (using ginext)
		payment := v1.Group("/payment")
		{
			payment.GET("/methods", ginext.WrapHandler(h.PaymentHandler.GetPaymentMethods))
			payment.POST("/process", ginext.WrapHandler(h.PaymentHandler.ProcessPayment))
		}

		// Feedback endpoints (using ginext)
		feedback := v1.Group("/feedback")
		{
			feedback.POST("", ginext.WrapHandler(h.FeedbackHandler.CreateFeedback))
			feedback.GET("/booking/:booking_id", ginext.WrapHandler(h.FeedbackHandler.GetBookingFeedback))
			feedback.GET("/trip/:trip_id", ginext.WrapHandler(h.FeedbackHandler.GetTripFeedbacks))
		}

		// Statistics endpoints (using ginext)
		statistics := v1.Group("/statistics")
		{
			statistics.GET("/bookings", ginext.WrapHandler(h.StatisticsHandler.GetBookingStats))
			statistics.GET("/popular-trips", ginext.WrapHandler(h.StatisticsHandler.GetPopularTrips))
		}

		// Seat operations (using ginext)
		seats := v1.Group("/seats")
		{
			seats.POST("/reserve", ginext.WrapHandler(h.SeatHandler.ReserveSeat))
			seats.POST("/release", ginext.WrapHandler(h.SeatHandler.ReleaseSeat))
		}

		// Seat locks (using ginext pattern)
		seatLocks := v1.Group("/seat-locks")
		{
			seatLocks.POST("", ginext.WrapHandler(h.SeatLockHandler.LockSeats))
			seatLocks.DELETE("", ginext.WrapHandler(h.SeatLockHandler.UnlockSeats))
		}
	}
}
