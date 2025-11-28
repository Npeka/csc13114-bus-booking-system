package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
)

// BookingService defines the business logic for booking operations
type BookingService interface {
	// Booking operations
	CreateBooking(ctx context.Context, req *model.CreateBookingRequest) (*model.BookingResponse, error)
	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID, page, limit int) (*model.PaginatedBookingResponse, error)
	GetTripBookings(ctx context.Context, tripID uuid.UUID, page, limit int) (*model.PaginatedBookingResponse, error)
	CancelBooking(ctx context.Context, id uuid.UUID, userID uuid.UUID, reason string) error
	UpdateBookingStatus(ctx context.Context, id uuid.UUID, status string) error
}

// bookingServiceImpl implements BookingService
type bookingServiceImpl struct {
	repositories *repository.Repositories
}

// NewBookingService creates a new booking service
func NewBookingService(repositories *repository.Repositories) BookingService {
	return &bookingServiceImpl{
		repositories: repositories,
	}
}

// CreateBooking creates a new booking with seat reservations
func (s *bookingServiceImpl) CreateBooking(ctx context.Context, req *model.CreateBookingRequest) (*model.BookingResponse, error) {
	// TODO: This method uses old model structure and needs complete refactoring
	// The new model uses Passengers instead of BookingSeats
	// For now, return error to prevent usage
	return nil, fmt.Errorf("CreateBooking not implemented - needs refactoring for new model structure")

	/* OLD CODE - COMMENTED OUT
	log.Info().
		Str("user_id", req.UserID.String()).
		Str("trip_id", req.TripID.String()).
		Int("seat_count", len(req.SeatIDs)).
		Msg("Creating new booking")

	// Validate request
	if err := s.validateBookingRequest(ctx, req); err != nil {
		return nil, fmt.Errorf("invalid booking request: %w", err)
	}

	// Get payment method
	paymentMethod, err := s.repositories.PaymentMethod.GetPaymentMethodByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment method: %w", err)
	}

	// Create booking model
	booking := &model.Booking{
		TripID:        req.TripID,
		UserID:        &req.UserID,
		TotalAmount:   req.TotalAmount,
		Status:        model.BookingStatus("pending"),
		PaymentStatus: model.PaymentStatus("pending"),
	}

	// Create booking in database
	if err := s.repositories.Booking.CreateBooking(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Convert to response
	response := s.toBookingResponse(booking, paymentMethod)

	log.Info().
		Str("booking_id", booking.ID.String()).
		Msg("Booking created successfully")

	return response, nil
	*/
}

// GetBookingByID retrieves a booking by ID
func (s *bookingServiceImpl) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error) {
	booking, err := s.repositories.Booking.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toBookingResponse(booking, booking.PaymentMethod), nil
}

// GetUserBookings retrieves bookings for a user with pagination
func (s *bookingServiceImpl) GetUserBookings(ctx context.Context, userID uuid.UUID, page, limit int) (*model.PaginatedBookingResponse, error) {
	offset := (page - 1) * limit
	bookings, total, err := s.repositories.Booking.GetBookingsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var bookingResponses []*model.BookingResponse
	for _, booking := range bookings {
		response := s.toBookingResponse(booking, booking.PaymentMethod)
		bookingResponses = append(bookingResponses, response)
	}

	return &model.PaginatedBookingResponse{
		Data:       bookingResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}, nil
}

// GetTripBookings retrieves bookings for a trip with pagination
func (s *bookingServiceImpl) GetTripBookings(ctx context.Context, tripID uuid.UUID, page, limit int) (*model.PaginatedBookingResponse, error) {
	offset := (page - 1) * limit
	bookings, total, err := s.repositories.Booking.GetBookingsByTripID(ctx, tripID, limit, offset)
	if err != nil {
		return nil, err
	}

	var bookingResponses []*model.BookingResponse
	for _, booking := range bookings {
		response := s.toBookingResponse(booking, booking.PaymentMethod)
		bookingResponses = append(bookingResponses, response)
	}

	return &model.PaginatedBookingResponse{
		Data:       bookingResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}, nil
}

// CancelBooking cancels a booking and releases seats
func (s *bookingServiceImpl) CancelBooking(ctx context.Context, id uuid.UUID, userID uuid.UUID, reason string) error {
	// Get booking to verify ownership
	booking, err := s.repositories.Booking.GetBookingByID(ctx, id)
	if err != nil {
		return err
	}

	if booking.UserID == nil || *booking.UserID != userID {
		return fmt.Errorf("booking does not belong to user")
	}

	if booking.Status == "cancelled" {
		return fmt.Errorf("booking is already cancelled")
	}

	if booking.Status == "completed" {
		return fmt.Errorf("cannot cancel completed booking")
	}

	return s.repositories.Booking.CancelBooking(ctx, id, reason)
}

// UpdateBookingStatus updates the status of a booking
func (s *bookingServiceImpl) UpdateBookingStatus(ctx context.Context, id uuid.UUID, status string) error {
	booking, err := s.repositories.Booking.GetBookingByID(ctx, id)
	if err != nil {
		return err
	}

	booking.Status = model.BookingStatus(status)
	booking.UpdatedAt = time.Now().UTC()

	if status == "completed" {
		confirmedAt := time.Now().UTC()
		booking.ConfirmedAt = &confirmedAt
	}

	return s.repositories.Booking.UpdateBooking(ctx, booking)
}

// GetPaymentMethods retrieves all available payment methods
func (s *bookingServiceImpl) GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethodResponse, error) {
	paymentMethods, err := s.repositories.PaymentMethod.GetPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*model.PaymentMethodResponse
	for _, pm := range paymentMethods {
		response := &model.PaymentMethodResponse{
			ID:          pm.ID,
			Name:        pm.Name,
			Code:        pm.Code,
			Description: pm.Description,
			IsActive:    pm.IsActive,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// ProcessPayment processes payment for a booking
func (s *bookingServiceImpl) ProcessPayment(ctx context.Context, req *model.ProcessPaymentRequest) (*model.PaymentResponse, error) {
	// Get booking
	booking, err := s.repositories.Booking.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		return nil, err
	}

	if booking.Status != "pending" {
		return nil, fmt.Errorf("booking is not in pending status")
	}

	// In a real implementation, you would integrate with payment gateway here
	// For now, we'll simulate payment processing

	response := &model.PaymentResponse{
		BookingID:       req.BookingID,
		Amount:          booking.TotalAmount,
		PaymentMethodID: uuid.Nil,
		Status:          "completed",
		TransactionID:   uuid.New().String(),
		ProcessedAt:     time.Now().UTC(),
	}

	// Update booking status to confirmed
	if err := s.UpdateBookingStatus(ctx, req.BookingID, "confirmed"); err != nil {
		return nil, fmt.Errorf("failed to update booking status: %w", err)
	}

	return response, nil
}

// Helper methods

// // validateBookingRequest validates a booking request
// func (s *bookingServiceImpl) validateBookingRequest(ctx context.Context, req *model.CreateBookingRequest) error {
// 	if req.UserID == uuid.Nil {
// 		return fmt.Errorf("user ID is required")
// 	}
// 	if req.TripID == uuid.Nil {
// 		return fmt.Errorf("trip ID is required")
// 	}
// 	if len(req.SeatIDs) == 0 {
// 		return fmt.Errorf("at least one seat must be selected")
// 	}
// 	if req.TotalAmount <= 0 {
// 		return fmt.Errorf("total amount must be greater than 0")
// 	}
// 	if req.PassengerName == "" {
// 		return fmt.Errorf("passenger name is required")
// 	}
// 	if req.PassengerPhone == "" {
// 		return fmt.Errorf("passenger phone is required")
// 	}

// 	// Check if all seats are available
// 	for _, seatID := range req.SeatIDs {
// 		seatStatuses, err := s.repositories.SeatStatus.GetSeatStatusByTripID(ctx, req.TripID)
// 		if err != nil {
// 			return fmt.Errorf("failed to check seat availability: %w", err)
// 		}

// 		seatFound := false
// 		for _, seatStatus := range seatStatuses {
// 			if seatStatus.SeatID == seatID {
// 				seatFound = true
// 				if seatStatus.Status != "available" {
// 					return fmt.Errorf("seat %s is not available", seatID.String())
// 				}
// 				break
// 			}
// 		}

// 		if !seatFound {
// 			return fmt.Errorf("seat %s not found for trip", seatID.String())
// 		}
// 	}

// 	return nil
// }

// toBookingResponse converts a booking model to response
func (s *bookingServiceImpl) toBookingResponse(booking *model.Booking, paymentMethod string) *model.BookingResponse {
	var userID uuid.UUID
	if booking.UserID != nil {
		userID = *booking.UserID
	}

	response := &model.BookingResponse{
		ID:                 booking.ID,
		UserID:             userID,
		TripID:             booking.TripID,
		Status:             string(booking.Status),
		TotalAmount:        booking.TotalAmount,
		PassengerName:      booking.GuestName,
		PassengerPhone:     booking.GuestPhone,
		PassengerEmail:     booking.GuestEmail,
		SpecialRequests:    "",
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
		CompletedAt:        booking.ConfirmedAt,
		CancelledAt:        booking.CancelledAt,
		CancellationReason: booking.CancellationReason,
	}

	// Map payment method if available
	if paymentMethod != "" {
		// Note: PaymentMethod is now a string in the new model, not a relation
		// This needs to be refactored if we want to return full payment method details
		response.PaymentMethod = nil
	}

	// Map passengers to seats response
	for _, passenger := range booking.Passengers {
		seatResponse := model.BookingSeatResponse{
			ID:     passenger.ID,
			SeatID: passenger.SeatID,
			Price:  passenger.Price,
		}
		response.Seats = append(response.Seats, seatResponse)
	}

	return response
}
