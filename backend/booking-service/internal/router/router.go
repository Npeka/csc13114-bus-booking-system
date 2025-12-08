package router

import (
	"bus-booking/booking-service/config"
	"bus-booking/booking-service/internal/handler"
	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/shared/health"
	"bus-booking/shared/middleware"
	"bus-booking/shared/swagger"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	BookingHandler    handler.BookingHandler
	FeedbackHandler   handler.FeedbackHandler
	StatisticsHandler handler.StatisticsHandler
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *Handlers) {
	router.Use(middleware.Logger())
	router.Use(middleware.SetupCORS(&cfg.CORS))
	router.Use(middleware.RequestContext(cfg.ServiceName))
	router.GET(health.Path, health.Handler(cfg.ServiceName))
	router.GET(swagger.Path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		bookings := v1.Group("/bookings")
		{
			// Public routes
			bookings.POST("/guest", ginext.WrapHandler(h.BookingHandler.CreateGuestBooking))
			bookings.GET("/lookup", ginext.WrapHandler(h.BookingHandler.GetBookingByReference))

			// E-ticket download - public route (accessible by booking ID)
			bookings.GET("/:id/eticket", func(c *gin.Context) {
				r := ginext.NewRequest(c)
				if err := h.BookingHandler.DownloadETicket(r); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": gin.H{"message": err.Error()},
					})
				}
			})

			// Authenticated routes
			bookings.Use(middleware.RequireAuth())
			bookings.POST("", ginext.WrapHandler(h.BookingHandler.CreateBooking))
			bookings.GET("/:id", ginext.WrapHandler(h.BookingHandler.GetBooking))
			bookings.POST("/:id/cancel", ginext.WrapHandler(h.BookingHandler.CancelBooking))
			bookings.GET("/user/:user_id", ginext.WrapHandler(h.BookingHandler.GetUserBookings))
		}

		internal := v1.Group("/bookings")
		{
			internal.PUT("/:id/payment-status", ginext.WrapHandler(h.BookingHandler.UpdatePaymentStatus))
			internal.GET("/trips/:trip_id/seats/status", ginext.WrapHandler(h.BookingHandler.GetSeatStatus))
		}

		feedback := v1.Group("/feedback")
		feedback.Use(middleware.RequireAuth())
		{
			feedback.POST("", ginext.WrapHandler(h.FeedbackHandler.CreateFeedback))
			feedback.GET("/booking/:booking_id", ginext.WrapHandler(h.FeedbackHandler.GetBookingFeedback))
			feedback.GET("/trip/:trip_id", ginext.WrapHandler(h.FeedbackHandler.GetTripFeedbacks))
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.RequireAuth())
		admin.Use(middleware.RequireRole(constants.RoleAdmin))
		{
			admin.PUT("/bookings/:id/status", ginext.WrapHandler(h.BookingHandler.UpdateBookingStatus))
			admin.GET("/bookings/trip/:trip_id", ginext.WrapHandler(h.BookingHandler.GetTripBookings))

			admin.GET("/statistics/bookings", ginext.WrapHandler(h.StatisticsHandler.GetBookingStats))
			admin.GET("/statistics/popular-trips", ginext.WrapHandler(h.StatisticsHandler.GetPopularTrips))
		}
	}
}
