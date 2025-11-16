package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
)

// BookingHandler handles HTTP requests for booking operations
type BookingHandler struct {
	bookingService service.BookingService
}

// NewBookingHandler creates a new booking handler
func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking with seat selection
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body model.CreateBookingRequest true "Booking creation request"
// @Success 201 {object} model.BookingResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req model.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
		return
	}

	// Create booking
	booking, err := h.bookingService.CreateBooking(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create booking")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Failed to create booking",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// GetBooking godoc
// @Summary Get booking by ID
// @Description Get a specific booking by its ID
// @Tags bookings
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} model.BookingResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid booking ID",
			Message: "Booking ID must be a valid UUID",
		})
		return
	}

	booking, err := h.bookingService.GetBookingByID(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("booking_id", id.String()).Msg("Failed to get booking")
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Error:   "Booking not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// GetUserBookings godoc
// @Summary Get user bookings
// @Description Get all bookings for a specific user with pagination
// @Tags bookings
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} model.PaginatedBookingResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /bookings/user/{user_id} [get]
func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid UUID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	bookings, err := h.bookingService.GetUserBookings(c.Request.Context(), userID, page, limit)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to get user bookings")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Failed to get bookings",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// GetTripBookings godoc
// @Summary Get trip bookings
// @Description Get all bookings for a specific trip with pagination
// @Tags bookings
// @Produce json
// @Param trip_id path string true "Trip ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} model.PaginatedBookingResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /bookings/trip/{trip_id} [get]
func (h *BookingHandler) GetTripBookings(c *gin.Context) {
	tripIDStr := c.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid trip ID",
			Message: "Trip ID must be a valid UUID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	bookings, err := h.bookingService.GetTripBookings(c.Request.Context(), tripID, page, limit)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID.String()).Msg("Failed to get trip bookings")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Failed to get bookings",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking and release seats
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Param request body model.CancelBookingRequest true "Cancellation request"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /bookings/{id}/cancel [post]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid booking ID",
			Message: "Booking ID must be a valid UUID",
		})
		return
	}

	var req model.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	if err := h.bookingService.CancelBooking(c.Request.Context(), id, req.UserID, req.Reason); err != nil {
		log.Error().Err(err).Str("booking_id", id.String()).Msg("Failed to cancel booking")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Failed to cancel booking",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Message: "Booking cancelled successfully",
	})
}

// UpdateBookingStatus godoc
// @Summary Update booking status
// @Description Update the status of a booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Param request body model.UpdateBookingStatusRequest true "Status update request"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /bookings/{id}/status [put]
func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid booking ID",
			Message: "Booking ID must be a valid UUID",
		})
		return
	}

	var req model.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	if err := h.bookingService.UpdateBookingStatus(c.Request.Context(), id, req.Status); err != nil {
		log.Error().Err(err).Str("booking_id", id.String()).Msg("Failed to update booking status")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Failed to update booking status",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Message: "Booking status updated successfully",
	})
}

// GetSeatAvailability godoc
// @Summary Get seat availability for a trip
// @Description Get available, reserved, and booked seats for a specific trip
// @Tags seats
// @Produce json
// @Param trip_id path string true "Trip ID"
// @Success 200 {object} model.SeatAvailabilityResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /trips/{trip_id}/seats [get]
func (h *BookingHandler) GetSeatAvailability(c *gin.Context) {
	tripIDStr := c.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid trip ID",
			Message: "Trip ID must be a valid UUID",
		})
		return
	}

	availability, err := h.bookingService.GetSeatAvailability(c.Request.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID.String()).Msg("Failed to get seat availability")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Failed to get seat availability",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// ReserveSeat godoc
// @Summary Reserve a seat
// @Description Reserve a seat for a user temporarily
// @Tags seats
// @Accept json
// @Produce json
// @Param request body model.ReserveSeatRequest true "Seat reservation request"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /seats/reserve [post]
func (h *BookingHandler) ReserveSeat(c *gin.Context) {
	var req model.ReserveSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	if err := h.bookingService.ReserveSeat(c.Request.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Failed to reserve seat")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Failed to reserve seat",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Message: "Seat reserved successfully",
	})
}

// ReleaseSeat godoc
// @Summary Release a reserved seat
// @Description Release a reserved seat back to available status
// @Tags seats
// @Accept json
// @Produce json
// @Param request body model.ReleaseSeatRequest true "Seat release request"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /seats/release [post]
func (h *BookingHandler) ReleaseSeat(c *gin.Context) {
	var req model.ReleaseSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	if err := h.bookingService.ReleaseSeat(c.Request.Context(), req.TripID, req.SeatID); err != nil {
		log.Error().Err(err).Msg("Failed to release seat")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Failed to release seat",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Message: "Seat released successfully",
	})
}

// GetPaymentMethods godoc
// @Summary Get payment methods
// @Description Get all available payment methods
// @Tags payment
// @Produce json
// @Success 200 {array} model.PaymentMethodResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /payment/methods [get]
func (h *BookingHandler) GetPaymentMethods(c *gin.Context) {
	methods, err := h.bookingService.GetPaymentMethods(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get payment methods")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Failed to get payment methods",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, methods)
}

// ProcessPayment godoc
// @Summary Process payment
// @Description Process payment for a booking
// @Tags payment
// @Accept json
// @Produce json
// @Param request body model.ProcessPaymentRequest true "Payment processing request"
// @Success 200 {object} model.PaymentResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /payment/process [post]
func (h *BookingHandler) ProcessPayment(c *gin.Context) {
	var req model.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	payment, err := h.bookingService.ProcessPayment(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process payment")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Failed to process payment",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// CreateFeedback godoc
// @Summary Create feedback
// @Description Create feedback for a completed booking
// @Tags feedback
// @Accept json
// @Produce json
// @Param request body model.CreateFeedbackRequest true "Feedback creation request"
// @Success 201 {object} model.FeedbackResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /feedback [post]
func (h *BookingHandler) CreateFeedback(c *gin.Context) {
	var req model.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	feedback, err := h.bookingService.CreateFeedback(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create feedback")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Failed to create feedback",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// GetBookingFeedback godoc
// @Summary Get booking feedback
// @Description Get feedback for a specific booking
// @Tags feedback
// @Produce json
// @Param booking_id path string true "Booking ID"
// @Success 200 {object} model.FeedbackResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /feedback/booking/{booking_id} [get]
func (h *BookingHandler) GetBookingFeedback(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid booking ID",
			Message: "Booking ID must be a valid UUID",
		})
		return
	}

	feedback, err := h.bookingService.GetBookingFeedback(c.Request.Context(), bookingID)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingID.String()).Msg("Failed to get booking feedback")
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Error:   "Feedback not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// GetTripFeedbacks godoc
// @Summary Get trip feedbacks
// @Description Get all feedbacks for a specific trip with pagination
// @Tags feedback
// @Produce json
// @Param trip_id path string true "Trip ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} model.PaginatedFeedbackResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /feedback/trip/{trip_id} [get]
func (h *BookingHandler) GetTripFeedbacks(c *gin.Context) {
	tripIDStr := c.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid trip ID",
			Message: "Trip ID must be a valid UUID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	feedbacks, err := h.bookingService.GetTripFeedbacks(c.Request.Context(), tripID, page, limit)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID.String()).Msg("Failed to get trip feedbacks")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Failed to get feedbacks",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

// GetBookingStats godoc
// @Summary Get booking statistics
// @Description Get booking statistics for a date range
// @Tags statistics
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} model.BookingStatsResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /statistics/bookings [get]
func (h *BookingHandler) GetBookingStats(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Missing date parameters",
			Message: "Both start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid start date",
			Message: "Date must be in YYYY-MM-DD format",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "Invalid end date",
			Message: "Date must be in YYYY-MM-DD format",
		})
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	stats, err := h.bookingService.GetBookingStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get booking statistics")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Failed to get statistics",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPopularTrips godoc
// @Summary Get popular trips
// @Description Get popular trips based on booking statistics
// @Tags statistics
// @Produce json
// @Param limit query int false "Number of trips to return" default(10)
// @Param days query int false "Number of days to look back" default(30)
// @Success 200 {array} model.TripStatsResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /statistics/popular-trips [get]
func (h *BookingHandler) GetPopularTrips(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	if limit < 1 || limit > 100 {
		limit = 10
	}
	if days < 1 || days > 365 {
		days = 30
	}

	trips, err := h.bookingService.GetPopularTrips(c.Request.Context(), limit, days)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get popular trips")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "Failed to get popular trips",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, trips)
}

// Health godoc
// @Summary Health check
// @Description Check service health
// @Tags health
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router /health [get]
func (h *BookingHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, model.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "booking-service",
	})
}
