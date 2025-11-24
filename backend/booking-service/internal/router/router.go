package router

import (
	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/handler"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	BookingHandler *handler.BookingHandler
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
			bookings.POST("", h.BookingHandler.CreateBooking)
			bookings.GET("/:id", h.BookingHandler.GetBooking)
			bookings.POST("/:id/cancel", h.BookingHandler.CancelBooking)
			bookings.PUT("/:id/status", h.BookingHandler.UpdateBookingStatus)
			bookings.GET("/user/:user_id", h.BookingHandler.GetUserBookings)
			bookings.GET("/trip/:trip_id", h.BookingHandler.GetTripBookings)
		}

		trips := v1.Group("/trips")
		{
			trips.GET("/:trip_id/seats", h.BookingHandler.GetSeatAvailability)
		}

		seats := v1.Group("/seats")
		{
			seats.POST("/reserve", h.BookingHandler.ReserveSeat)
			seats.POST("/release", h.BookingHandler.ReleaseSeat)
		}

		payment := v1.Group("/payment")
		{
			payment.GET("/methods", h.BookingHandler.GetPaymentMethods)
			payment.POST("/process", h.BookingHandler.ProcessPayment)
		}

		feedback := v1.Group("/feedback")
		{
			feedback.POST("", h.BookingHandler.CreateFeedback)
			feedback.GET("/booking/:booking_id", h.BookingHandler.GetBookingFeedback)
			feedback.GET("/trip/:trip_id", h.BookingHandler.GetTripFeedbacks)
		}

		statistics := v1.Group("/statistics")
		{
			statistics.GET("/bookings", h.BookingHandler.GetBookingStats)
			statistics.GET("/popular-trips", h.BookingHandler.GetPopularTrips)
		}
	}
}
