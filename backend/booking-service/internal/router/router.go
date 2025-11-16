package router

import (
	"time"

	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the main router
func SetupRouter(cfg *config.Config, bookingHandler *handler.BookingHandler) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Second,
	}
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", bookingHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Booking routes
		bookings := v1.Group("/bookings")
		{
			bookings.POST("", bookingHandler.CreateBooking)
			bookings.GET("/:id", bookingHandler.GetBooking)
			bookings.POST("/:id/cancel", bookingHandler.CancelBooking)
			bookings.PUT("/:id/status", bookingHandler.UpdateBookingStatus)
			bookings.GET("/user/:user_id", bookingHandler.GetUserBookings)
			bookings.GET("/trip/:trip_id", bookingHandler.GetTripBookings)
		}

		// Trip routes (for seat management)
		trips := v1.Group("/trips")
		{
			trips.GET("/:trip_id/seats", bookingHandler.GetSeatAvailability)
		}

		// Seat management routes
		seats := v1.Group("/seats")
		{
			seats.POST("/reserve", bookingHandler.ReserveSeat)
			seats.POST("/release", bookingHandler.ReleaseSeat)
		}

		// Payment routes
		payment := v1.Group("/payment")
		{
			payment.GET("/methods", bookingHandler.GetPaymentMethods)
			payment.POST("/process", bookingHandler.ProcessPayment)
		}

		// Feedback routes
		feedback := v1.Group("/feedback")
		{
			feedback.POST("", bookingHandler.CreateFeedback)
			feedback.GET("/booking/:booking_id", bookingHandler.GetBookingFeedback)
			feedback.GET("/trip/:trip_id", bookingHandler.GetTripFeedbacks)
		}

		// Statistics routes
		statistics := v1.Group("/statistics")
		{
			statistics.GET("/bookings", bookingHandler.GetBookingStats)
			statistics.GET("/popular-trips", bookingHandler.GetPopularTrips)
		}
	}

	return router
}
