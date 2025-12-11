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
	CreateGuestBooking(r *ginext.Request) (*ginext.Response, error)

	GetByID(r *ginext.Request) (*ginext.Response, error)
	GetByReference(r *ginext.Request) (*ginext.Response, error)
	GetUserBookings(r *ginext.Request) (*ginext.Response, error)
	GetTripBookings(r *ginext.Request) (*ginext.Response, error)

	CancelBooking(r *ginext.Request) (*ginext.Response, error)
	RetryPayment(r *ginext.Request) (*ginext.Response, error)

	UpdateBookingStatus(r *ginext.Request) (*ginext.Response, error)
	GetSeatStatus(r *ginext.Request) (*ginext.Response, error)

	DownloadETicket(r *ginext.Request) error
}

type BookingHandlerImpl struct {
	bookingService service.BookingService
	eTicketService service.ETicketService
}

func NewBookingHandler(bookingService service.BookingService, eTicketService service.ETicketService) BookingHandler {
	return &BookingHandlerImpl{
		bookingService: bookingService,
		eTicketService: eTicketService,
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

// CreateGuestBooking godoc
// @Summary Create a guest booking
// @Description Create a booking without authentication for guest users
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body model.CreateGuestBookingRequest true "Guest booking creation request"
// @Success 201 {object} ginext.Response{data=model.BookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/guest [post]
func (h *BookingHandlerImpl) CreateGuestBooking(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateGuestBookingRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind and validate guest booking request")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	booking, err := h.bookingService.CreateGuestBooking(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create guest booking")
		return nil, err
	}

	return ginext.NewSuccessResponse(booking), nil
}

// GetByID godoc
// @Summary Get booking by ID
// @Description Get a specific booking by its ID
// @Tags bookings
// @Produce json
// @Param id path string true "Booking ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.BookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/bookings/{id} [get]
func (h *BookingHandlerImpl) GetByID(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("invalid booking id")
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	booking, err := h.bookingService.GetByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to get booking")
		return nil, err
	}

	return ginext.NewSuccessResponse(booking), nil
}

// GetByReference godoc
// @Summary Get booking by reference number
// @Description Get booking details using booking reference number (for guest lookup)
// @Tags bookings
// @Produce json
// @Param reference query string true "Booking reference number"
// @Param email query string false "Email for verification"
// @Success 200 {object} ginext.Response{data=model.BookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Router /api/v1/bookings/lookup [get]
func (h *BookingHandlerImpl) GetByReference(r *ginext.Request) (*ginext.Response, error) {
	reference := r.GinCtx.Query("reference")
	email := r.GinCtx.Query("email")

	if reference == "" {
		return nil, ginext.NewBadRequestError("Booking reference is required")
	}

	booking, err := h.bookingService.GetByReference(r.Context(), reference, email)
	if err != nil {
		log.Error().Err(err).Str("reference", reference).Msg("failed to get booking by reference")
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

	var req model.GetUserBookingsRequest
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

	if err := h.bookingService.CancelBooking(r.Context(), id, req.Reason); err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to cancel booking")
		return nil, err
	}

	return ginext.NewSuccessResponse("booking cancelled successfully"), nil
}

// RetryPayment godoc
// @Summary Retry payment for a booking
// @Description Create a new payment link for a failed or expired booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true \"Booking ID\" format(uuid)
// @Success 200 {object} ginext.Response{data=model.BookingResponse}
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/{id}/retry-payment [post]
func (h *BookingHandlerImpl) RetryPayment(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("invalid booking id")
		return nil, ginext.NewBadRequestError("invalid booking id")
	}

	booking, err := h.bookingService.RetryPayment(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to retry payment")
		return nil, err
	}

	return ginext.NewSuccessResponse(booking), nil
}

// UpdateBookingStatus godoc
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

	if err := h.bookingService.UpdateBookingStatus(r.Context(), &req, id); err != nil {
		log.Error().Err(err).Str("booking_id", idStr).Msg("failed to update payment status")
		return nil, err
	}

	return ginext.NewSuccessResponse("payment status updated successfully"), nil
}

// GetSeatStatus godoc
// @Summary Get seat status for a trip
// @Description Get booking status of seats for a specific trip
// @Tags bookings
// @Accept json
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Param seat_ids query []string true "Seat IDs" collectionFormat(multi)
// @Success 200 {object} ginext.Response{data=[]model.SeatStatusItem}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/trips/{trip_id}/seats/status [get]
func (h *BookingHandlerImpl) GetSeatStatus(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("invalid trip id")
		return nil, ginext.NewBadRequestError("invalid trip id")
	}

	var req model.GetSeatStatusRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind query params")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	seatIDs := make([]uuid.UUID, 0, len(req.SeatIDs))
	for _, seatIDStr := range req.SeatIDs {
		seatID, err := uuid.Parse(seatIDStr)
		if err != nil {
			log.Error().Err(err).Str("seat_id", seatIDStr).Msg("invalid seat id")
			return nil, ginext.NewBadRequestError("invalid seat id: " + seatIDStr)
		}
		seatIDs = append(seatIDs, seatID)
	}

	seatStatuses, err := h.bookingService.GetSeatStatus(r.Context(), tripID, seatIDs)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("failed to get seat status")
		return nil, err
	}

	return ginext.NewSuccessResponse(seatStatuses), nil
}

// DownloadETicket godoc
// @Summary Download e-ticket PDF
// @Description Download e-ticket PDF for a confirmed booking
// @Tags bookings
// @Produce application/pdf
// @Param id path string true "Booking ID"
// @Success 200 {file} binary "PDF file"
// @Failure 400 {object} ginext.Response
// @Failure 404 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/bookings/{id}/eticket [get]
func (h *BookingHandlerImpl) DownloadETicket(r *ginext.Request) error {
	bookingIDStr := r.GinCtx.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingIDStr).Msg("invalid booking id")
		return ginext.NewBadRequestError("invalid booking id")
	}

	// Generate PDF
	pdfBuffer, err := h.eTicketService.GenerateETicket(r.Context(), bookingID)
	if err != nil {
		log.Error().Err(err).Str("booking_id", bookingIDStr).Msg("failed to generate e-ticket")
		return err
	}

	// Set headers for PDF download
	r.GinCtx.Header("Content-Type", "application/pdf")
	r.GinCtx.Header("Content-Disposition", "attachment; filename=eticket_"+bookingIDStr[:8]+".pdf")
	r.GinCtx.Data(200, "application/pdf", pdfBuffer.Bytes())

	return nil
}
