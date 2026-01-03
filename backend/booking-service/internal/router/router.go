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
	StatisticsHandler handler.StatisticsHandler
	SeatLockHandler   handler.SeatLockHandler
	ReviewHandler     handler.ReviewHandler
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
			bookings.POST("/guest", ginext.WrapHandler(h.BookingHandler.CreateGuestBooking))
			bookings.GET("/lookup", ginext.WrapHandler(h.BookingHandler.GetByReference))
			bookings.GET("/:id", ginext.WrapHandler(h.BookingHandler.GetByID))

			bookings.GET("/:id/eticket", func(c *gin.Context) {
				r := ginext.NewRequest(c)
				if err := h.BookingHandler.DownloadETicket(r); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": gin.H{"message": err.Error()},
					})
				}
			})
		}

		seatLocks := v1.Group("/seat-locks")
		{
			seatLocks.POST("", ginext.WrapHandler(h.SeatLockHandler.LockSeats))
			seatLocks.DELETE("", ginext.WrapHandler(h.SeatLockHandler.UnlockSeats))
		}

		trips := v1.Group("/trips")
		{
			trips.GET("/:trip_id/locked-seats", ginext.WrapHandler(h.SeatLockHandler.GetLockedSeats))
			trips.GET("/:trip_id/reviews", ginext.WrapHandler(h.ReviewHandler.GetTripReviews))
			trips.GET("/:trip_id/reviews/summary", ginext.WrapHandler(h.ReviewHandler.GetTripReviewSummary))
		}
	}

	userV1 := router.Group("/api/v1")
	userV1.Use(middleware.RequireAuth())
	{
		bookings := userV1.Group("/bookings")
		{
			bookings.POST("", ginext.WrapHandler(h.BookingHandler.CreateBooking))
			bookings.POST("/:id/cancel", ginext.WrapHandler(h.BookingHandler.CancelBooking))
			bookings.POST("/:id/retry-payment", ginext.WrapHandler(h.BookingHandler.RetryPayment))
			bookings.GET("/user/:user_id", ginext.WrapHandler(h.BookingHandler.GetUserBookings))
			bookings.POST("/:id/review", ginext.WrapHandler(h.ReviewHandler.CreateReview))
			bookings.GET("/:id/review", ginext.WrapHandler(h.ReviewHandler.GetReviewByBooking))
		}

		reviews := userV1.Group("/reviews")
		{
			reviews.PUT("/:id", ginext.WrapHandler(h.ReviewHandler.UpdateReview))
			reviews.DELETE("/:id", ginext.WrapHandler(h.ReviewHandler.DeleteReview))
		}

		users := userV1.Group("/users")
		{
			users.GET("/:user_id/reviews", ginext.WrapHandler(h.ReviewHandler.GetUserReviews))
		}
	}

	adminV1 := router.Group("/api/v1")
	adminV1.Use(middleware.RequireAuth())
	adminV1.Use(middleware.RequireRole(constants.RoleAdmin))
	{
		bookings := adminV1.Group("/bookings")
		{
			bookings.GET("/trip/:trip_id", ginext.WrapHandler(h.BookingHandler.GetTripBookings))
			bookings.GET("/trip/:trip_id/passengers", ginext.WrapHandler(h.BookingHandler.GetTripPassengers))
			bookings.GET("", ginext.WrapHandler(h.BookingHandler.ListBookings))
			bookings.POST("/:id/check-in", ginext.WrapHandler(h.BookingHandler.CheckInPassenger))
		}

		statistics := adminV1.Group("/statistics")
		{
			statistics.GET("/bookings", ginext.WrapHandler(h.StatisticsHandler.GetBookingStats))
			statistics.GET("/popular-trips", ginext.WrapHandler(h.StatisticsHandler.GetPopularTrips))
		}

		reviews := adminV1.Group("/reviews")
		{
			reviews.PUT("/:id/moderate", ginext.WrapHandler(h.ReviewHandler.ModerateReview))
		}
	}

	internalV1 := router.Group("/api/v1")
	{
		bookings := internalV1.Group("/bookings")
		{
			bookings.PUT("/:id/status", ginext.WrapHandler(h.BookingHandler.UpdateBookingStatus))
			bookings.GET("/trips/:trip_id/seats/status", ginext.WrapHandler(h.BookingHandler.GetSeatStatus))
		}
	}
}
