package handler

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	sharedcontext "bus-booking/shared/context"

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
	CreatePayment(r *ginext.Request) (*ginext.Response, error)
	UpdatePaymentStatus(r *ginext.Request) (*ginext.Response, error)
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
	userID := sharedcontext.GetUserID(r.GinCtx)

	var req model.CreateBookingRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind and validate request")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	booking, err := h.bookingService.CreateBooking(r.Context(), &req, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create booking")
		return nil, err
	}

	return ginext.NewSuccessResponse(booking), nil
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
		log.Error().Err(err).Str("id", idStr).Msg("invalid booking id")
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	booking, err := h.bookingService.GetBookingByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to get booking")
		return nil, err
	}

	return ginext.NewSuccessResponse(booking), nil
}

// GetUserBookings godoc
// @Summary Get user bookings
// @Description Get all bookings for a specific user with pagination
// @Tags bookings
// @Produce json
// @Param user_id path string true "User ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Success 200 {object} ginext.Response{data=model.PaginatedBookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/user/{user_id} [get]
func (h *BookingHandlerImpl) GetUserBookings(r *ginext.Request) (*ginext.Response, error) {
	userIDStr := r.GinCtx.Param("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Error().Err(err).Str("user_id", userIDStr).Msg("invalid user id")
		return nil, ginext.NewBadRequestError("invalid user id")
	}

	var req model.PaginationRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind query parameters")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	req.Normalize()

	bookings, total, err := h.bookingService.GetUserBookings(r.Context(), req, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userIDStr).Msg("failed to get user bookings")
		return nil, err
	}

	return ginext.NewPaginatedResponse(bookings, req.Page, req.PageSize, total), nil
}

// GetTripBookings godoc
// @Summary Get trip bookings
// @Description Get all bookings for a specific trip with pagination
// @Tags bookings
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Success 200 {object} ginext.Response{data=model.PaginatedBookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/trip/{trip_id} [get]
func (h *BookingHandlerImpl) GetTripBookings(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("invalid trip id")
		return nil, ginext.NewBadRequestError("invalid trip id")
	}

	var req model.PaginationRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind query parameters")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	req.Normalize()

	bookings, total, err := h.bookingService.GetTripBookings(r.Context(), req, tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("failed to get trip bookings")
		return nil, err
	}

	return ginext.NewPaginatedResponse(bookings, req.Page, req.PageSize, total), nil
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
		log.Error().Err(err).Str("id", idStr).Msg("invalid booking id")
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	var req model.CancelBookingRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.bookingService.CancelBooking(r.Context(), id, req.UserID, req.Reason); err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to cancel booking")
		return nil, err
	}

	return ginext.NewSuccessResponse("booking cancelled successfully"), nil
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
		log.Error().Err(err).Str("id", idStr).Msg("invalid booking id")
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	var req model.UpdateBookingStatusRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.bookingService.UpdateBookingStatus(r.Context(), id, req.Status); err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to update booking status")
		return nil, err
	}

	return ginext.NewSuccessResponse("booking status updated successfully"), nil
}

// CreatePayment godoc
// @Summary Create payment link for booking
// @Description Create a PayOS payment link for a booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID" format(uuid)
// @Param request body model.CreatePaymentRequest true "Buyer information"
// @Success 200 {object} ginext.Response{data=client.PaymentLinkResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/bookings/{id}/payment [post]
func (h *BookingHandlerImpl) CreatePayment(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("invalid booking id")
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	var req model.CreatePaymentRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	paymentResp, err := h.bookingService.CreatePayment(r.Context(), id, &req.BuyerInfo)
	if err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to create payment")
		return nil, err
	}

	return ginext.NewSuccessResponse(paymentResp), nil
}

// UpdatePaymentStatus godoc
// @Summary Update booking payment status
// @Description Update booking payment status (internal use by payment service)
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID" format(uuid)
// @Param request body model.UpdatePaymentStatusRequest true "Payment status update"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/bookings/{id}/payment-status [put]
func (h *BookingHandlerImpl) UpdatePaymentStatus(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("invalid booking id")
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	var req model.UpdatePaymentStatusRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.bookingService.UpdatePaymentStatus(r.Context(), id, req.PaymentStatus, req.BookingStatus, req.PaymentOrderID); err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to update payment status")
		return nil, err
	}

	return ginext.NewSuccessResponse("payment status updated successfully"), nil
}
