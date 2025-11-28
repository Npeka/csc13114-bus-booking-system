package handler

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BookingHandler interface {
	CreateBooking(r *ginext.Request) (*ginext.Response, error)
	GetBooking(r *ginext.Request) (*ginext.Response, error)
	GetUserBookings(r *ginext.Request) (*ginext.Response, error)
	GetTripBookings(r *ginext.Request) (*ginext.Response, error)
	CancelBooking(r *ginext.Request) (*ginext.Response, error)
	UpdateBookingStatus(r *ginext.Request) (*ginext.Response, error)
}

type BookingHandlerImpl struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) BookingHandler {
	return &BookingHandlerImpl{
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
// @Success 201 {object} ginext.Response{data=model.BookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings [post]
func (h *BookingHandlerImpl) CreateBooking(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateBookingRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	// Create booking
	booking, err := h.bookingService.CreateBooking(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create booking")
		return nil, err
	}

	return ginext.NewSuccessResponse(booking, "Booking created successfully"), nil
}

// GetBooking godoc
// @Summary Get booking by ID
// @Description Get a specific booking by its ID
// @Tags bookings
// @Produce json
// @Param id path string true "Booking ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.BookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/bookings/{id} [get]
func (h *BookingHandlerImpl) GetBooking(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid booking ID")
	}

	booking, err := h.bookingService.GetBookingByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("Failed to get booking")
		return nil, err
	}

	return ginext.NewSuccessResponse(booking, "Booking retrieved successfully"), nil
}

// GetUserBookings godoc
// @Summary Get user bookings
// @Description Get all bookings for a specific user with pagination
// @Tags bookings
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} ginext.Response{data=model.PaginatedBookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/user/{user_id} [get]
func (h *BookingHandlerImpl) GetUserBookings(r *ginext.Request) (*ginext.Response, error) {
	userIDStr := r.GinCtx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid user ID")
	}

	var query struct {
		Page  int `form:"page,default=1"`
		Limit int `form:"limit,default=10"`
	}

	if err := r.GinCtx.ShouldBindQuery(&query); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	// Set defaults if not provided
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 10
	}

	bookings, err := h.bookingService.GetUserBookings(r.Context(), userID, query.Page, query.Limit)
	if err != nil {
		log.Error().Err(err).Str("user_id", userIDStr).Msg("Failed to get user bookings")
		return nil, err
	}

	return ginext.NewSuccessResponse(bookings, "User bookings retrieved successfully"), nil
}

// GetTripBookings godoc
// @Summary Get trip bookings
// @Description Get all bookings for a specific trip with pagination
// @Tags bookings
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} ginext.Response{data=model.PaginatedBookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/trip/{trip_id} [get]
func (h *BookingHandlerImpl) GetTripBookings(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	var query struct {
		Page  int `form:"page,default=1"`
		Limit int `form:"limit,default=10"`
	}

	if err := r.GinCtx.ShouldBindQuery(&query); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	// Set defaults if not provided
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 10
	}

	bookings, err := h.bookingService.GetTripBookings(r.Context(), tripID, query.Page, query.Limit)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("Failed to get trip bookings")
		return nil, err
	}

	return ginext.NewSuccessResponse(bookings, "Trip bookings retrieved successfully"), nil
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking and release seats
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID" format(uuid)
// @Param request body model.CancelBookingRequest true "Cancellation request"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/bookings/{id}/cancel [post]
func (h *BookingHandlerImpl) CancelBooking(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid booking ID")
	}

	var req model.CancelBookingRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.bookingService.CancelBooking(r.Context(), id, req.UserID, req.Reason); err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("Failed to cancel booking")
		return nil, err
	}

	return ginext.NewSuccessResponse(nil, "Booking cancelled successfully"), nil
}

// UpdateBookingStatus godoc
// @Summary Update booking status
// @Description Update the status of a booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID" format(uuid)
// @Param request body model.UpdateBookingStatusRequest true "Status update request"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/bookings/{id}/status [put]
func (h *BookingHandlerImpl) UpdateBookingStatus(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid booking ID")
	}

	var req model.UpdateBookingStatusRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.bookingService.UpdateBookingStatus(r.Context(), id, req.Status); err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("Failed to update booking status")
		return nil, err
	}

	return ginext.NewSuccessResponse(nil, "Booking status updated successfully"), nil
}
