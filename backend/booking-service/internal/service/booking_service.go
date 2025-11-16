package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

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

	// Seat operations
	GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error)
	ReserveSeat(ctx context.Context, req *model.ReserveSeatRequest) error
	ReleaseSeat(ctx context.Context, tripID, seatID uuid.UUID) error
	CheckReservationExpiry(ctx context.Context) error

	// Payment operations
	GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethodResponse, error)
	ProcessPayment(ctx context.Context, req *model.ProcessPaymentRequest) (*model.PaymentResponse, error)

	// Feedback operations
	CreateFeedback(ctx context.Context, req *model.CreateFeedbackRequest) (*model.FeedbackResponse, error)
	GetBookingFeedback(ctx context.Context, bookingID uuid.UUID) (*model.FeedbackResponse, error)
	GetTripFeedbacks(ctx context.Context, tripID uuid.UUID, page, limit int) (*model.PaginatedFeedbackResponse, error)

	// Statistics
	GetBookingStats(ctx context.Context, startDate, endDate time.Time) (*model.BookingStatsResponse, error)
	GetPopularTrips(ctx context.Context, limit, days int) ([]*model.TripStatsResponse, error)
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
		ID:              uuid.New(),
		UserID:          req.UserID,
		TripID:          req.TripID,
		PaymentMethodID: req.PaymentMethodID,
		TotalAmount:     req.TotalAmount,
		Status:          "pending",
		PassengerName:   req.PassengerName,
		PassengerPhone:  req.PassengerPhone,
		PassengerEmail:  req.PassengerEmail,
		SpecialRequests: req.SpecialRequests,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	// Create booking seats
	for _, seatID := range req.SeatIDs {
		bookingSeat := model.BookingSeat{
			ID:        uuid.New(),
			BookingID: booking.ID,
			SeatID:    seatID,
			Price:     req.SeatPrice,
			CreatedAt: time.Now().UTC(),
		}
		booking.BookingSeats = append(booking.BookingSeats, bookingSeat)
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
}

// GetBookingByID retrieves a booking by ID
func (s *bookingServiceImpl) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error) {
	booking, err := s.repositories.Booking.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toBookingResponse(booking, &booking.PaymentMethod), nil
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
		response := s.toBookingResponse(booking, &booking.PaymentMethod)
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
		response := s.toBookingResponse(booking, &booking.PaymentMethod)
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

	if booking.UserID != userID {
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

	booking.Status = status
	booking.UpdatedAt = time.Now().UTC()

	if status == "completed" {
		completedAt := time.Now().UTC()
		booking.CompletedAt = &completedAt
	}

	return s.repositories.Booking.UpdateBooking(ctx, booking)
}

// GetSeatAvailability retrieves seat availability for a trip
func (s *bookingServiceImpl) GetSeatAvailability(ctx context.Context, tripID uuid.UUID) (*model.SeatAvailabilityResponse, error) {
	seatStatuses, err := s.repositories.SeatStatus.GetSeatStatusByTripID(ctx, tripID)
	if err != nil {
		return nil, err
	}

	var availableSeats, reservedSeats, bookedSeats []uuid.UUID
	seatMap := make(map[uuid.UUID]*model.SeatStatus)

	for _, seatStatus := range seatStatuses {
		seatMap[seatStatus.SeatID] = seatStatus

		switch seatStatus.Status {
		case "available":
			availableSeats = append(availableSeats, seatStatus.SeatID)
		case "reserved":
			// Check if reservation has expired
			if seatStatus.ReservedUntil != nil && time.Now().UTC().After(*seatStatus.ReservedUntil) {
				// Release expired reservation
				if err := s.repositories.SeatStatus.ReleaseSeat(ctx, tripID, seatStatus.SeatID); err != nil {
					log.Error().Err(err).Msg("Failed to release expired seat reservation")
				} else {
					availableSeats = append(availableSeats, seatStatus.SeatID)
				}
			} else {
				reservedSeats = append(reservedSeats, seatStatus.SeatID)
			}
		case "booked":
			bookedSeats = append(bookedSeats, seatStatus.SeatID)
		}
	}

	return &model.SeatAvailabilityResponse{
		TripID:         tripID,
		AvailableSeats: availableSeats,
		ReservedSeats:  reservedSeats,
		BookedSeats:    bookedSeats,
		SeatDetails:    seatMap,
	}, nil
}

// ReserveSeat reserves a seat for a user
func (s *bookingServiceImpl) ReserveSeat(ctx context.Context, req *model.ReserveSeatRequest) error {
	// Check if seat is available
	seatStatuses, err := s.repositories.SeatStatus.GetSeatStatusByTripID(ctx, req.TripID)
	if err != nil {
		return err
	}

	for _, seatStatus := range seatStatuses {
		if seatStatus.SeatID == req.SeatID {
			if seatStatus.Status != "available" {
				return fmt.Errorf("seat is not available")
			}
			break
		}
	}

	// Default reservation time is 15 minutes
	reservationTime := 15 * time.Minute
	if req.ReservationMinutes > 0 {
		reservationTime = time.Duration(req.ReservationMinutes) * time.Minute
	}

	return s.repositories.SeatStatus.ReserveSeat(ctx, req.TripID, req.SeatID, req.UserID, reservationTime)
}

// ReleaseSeat releases a reserved seat
func (s *bookingServiceImpl) ReleaseSeat(ctx context.Context, tripID, seatID uuid.UUID) error {
	return s.repositories.SeatStatus.ReleaseSeat(ctx, tripID, seatID)
}

// CheckReservationExpiry checks and releases expired seat reservations
func (s *bookingServiceImpl) CheckReservationExpiry(ctx context.Context) error {
	// This should be called periodically by a background job
	// For now, it's a placeholder implementation
	log.Info().Msg("Checking for expired seat reservations")

	// In a real implementation, you would:
	// 1. Query all reserved seats with expired reservations
	// 2. Release them back to available status
	// 3. Optionally notify users about expired reservations

	return nil
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
		PaymentMethodID: booking.PaymentMethodID,
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

// CreateFeedback creates feedback for a booking
func (s *bookingServiceImpl) CreateFeedback(ctx context.Context, req *model.CreateFeedbackRequest) (*model.FeedbackResponse, error) {
	// Verify booking exists and belongs to user
	booking, err := s.repositories.Booking.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		return nil, err
	}

	if booking.UserID != req.UserID {
		return nil, fmt.Errorf("booking does not belong to user")
	}

	if booking.Status != "completed" {
		return nil, fmt.Errorf("can only provide feedback for completed bookings")
	}

	// Check if feedback already exists
	if _, err := s.repositories.Feedback.GetFeedbackByBookingID(ctx, req.BookingID); err == nil {
		return nil, fmt.Errorf("feedback already exists for this booking")
	}

	feedback := &model.Feedback{
		ID:        uuid.New(),
		UserID:    req.UserID,
		BookingID: req.BookingID,
		TripID:    booking.TripID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repositories.Feedback.CreateFeedback(ctx, feedback); err != nil {
		return nil, err
	}

	return &model.FeedbackResponse{
		ID:        feedback.ID,
		UserID:    feedback.UserID,
		BookingID: feedback.BookingID,
		TripID:    feedback.TripID,
		Rating:    feedback.Rating,
		Comment:   feedback.Comment,
		CreatedAt: feedback.CreatedAt,
	}, nil
}

// GetBookingFeedback retrieves feedback for a booking
func (s *bookingServiceImpl) GetBookingFeedback(ctx context.Context, bookingID uuid.UUID) (*model.FeedbackResponse, error) {
	feedback, err := s.repositories.Feedback.GetFeedbackByBookingID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	return &model.FeedbackResponse{
		ID:        feedback.ID,
		UserID:    feedback.UserID,
		BookingID: feedback.BookingID,
		TripID:    feedback.TripID,
		Rating:    feedback.Rating,
		Comment:   feedback.Comment,
		CreatedAt: feedback.CreatedAt,
	}, nil
}

// GetTripFeedbacks retrieves feedbacks for a trip with pagination
func (s *bookingServiceImpl) GetTripFeedbacks(ctx context.Context, tripID uuid.UUID, page, limit int) (*model.PaginatedFeedbackResponse, error) {
	offset := (page - 1) * limit
	feedbacks, total, err := s.repositories.Feedback.GetFeedbacksByTripID(ctx, tripID, limit, offset)
	if err != nil {
		return nil, err
	}

	var feedbackResponses []*model.FeedbackResponse
	for _, feedback := range feedbacks {
		response := &model.FeedbackResponse{
			ID:        feedback.ID,
			UserID:    feedback.UserID,
			BookingID: feedback.BookingID,
			TripID:    feedback.TripID,
			Rating:    feedback.Rating,
			Comment:   feedback.Comment,
			CreatedAt: feedback.CreatedAt,
		}
		feedbackResponses = append(feedbackResponses, response)
	}

	return &model.PaginatedFeedbackResponse{
		Data:       feedbackResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}, nil
}

// GetBookingStats retrieves booking statistics for a date range
func (s *bookingServiceImpl) GetBookingStats(ctx context.Context, startDate, endDate time.Time) (*model.BookingStatsResponse, error) {
	stats, err := s.repositories.BookingStats.GetBookingStatsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &model.BookingStatsResponse{
		TotalBookings:     stats.TotalBookings,
		TotalRevenue:      stats.TotalRevenue,
		CancelledBookings: stats.CancelledBookings,
		CompletedBookings: stats.CompletedBookings,
		AverageRating:     stats.AverageRating,
		StartDate:         startDate,
		EndDate:           endDate,
	}, nil
}

// GetPopularTrips retrieves popular trips based on booking statistics
func (s *bookingServiceImpl) GetPopularTrips(ctx context.Context, limit, days int) ([]*model.TripStatsResponse, error) {
	stats, err := s.repositories.BookingStats.GetPopularTrips(ctx, limit, days)
	if err != nil {
		return nil, err
	}

	var responses []*model.TripStatsResponse
	for _, stat := range stats {
		response := &model.TripStatsResponse{
			TripID:        stat.TripID,
			TotalBookings: stat.TotalBookings,
			TotalRevenue:  stat.TotalRevenue,
			AverageRating: stat.AverageRating,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// Helper methods

// validateBookingRequest validates a booking request
func (s *bookingServiceImpl) validateBookingRequest(ctx context.Context, req *model.CreateBookingRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}
	if req.TripID == uuid.Nil {
		return fmt.Errorf("trip ID is required")
	}
	if len(req.SeatIDs) == 0 {
		return fmt.Errorf("at least one seat must be selected")
	}
	if req.TotalAmount <= 0 {
		return fmt.Errorf("total amount must be greater than 0")
	}
	if req.PassengerName == "" {
		return fmt.Errorf("passenger name is required")
	}
	if req.PassengerPhone == "" {
		return fmt.Errorf("passenger phone is required")
	}

	// Check if all seats are available
	for _, seatID := range req.SeatIDs {
		seatStatuses, err := s.repositories.SeatStatus.GetSeatStatusByTripID(ctx, req.TripID)
		if err != nil {
			return fmt.Errorf("failed to check seat availability: %w", err)
		}

		seatFound := false
		for _, seatStatus := range seatStatuses {
			if seatStatus.SeatID == seatID {
				seatFound = true
				if seatStatus.Status != "available" {
					return fmt.Errorf("seat %s is not available", seatID.String())
				}
				break
			}
		}

		if !seatFound {
			return fmt.Errorf("seat %s not found for trip", seatID.String())
		}
	}

	return nil
}

// toBookingResponse converts a booking model to response
func (s *bookingServiceImpl) toBookingResponse(booking *model.Booking, paymentMethod *model.PaymentMethod) *model.BookingResponse {
	response := &model.BookingResponse{
		ID:              booking.ID,
		UserID:          booking.UserID,
		TripID:          booking.TripID,
		Status:          booking.Status,
		TotalAmount:     booking.TotalAmount,
		PassengerName:   booking.PassengerName,
		PassengerPhone:  booking.PassengerPhone,
		PassengerEmail:  booking.PassengerEmail,
		SpecialRequests: booking.SpecialRequests,
		CreatedAt:       booking.CreatedAt,
		UpdatedAt:       booking.UpdatedAt,
		CompletedAt:     booking.CompletedAt,
		CancelledAt:     booking.CancelledAt,
	}

	if paymentMethod != nil {
		response.PaymentMethod = &model.PaymentMethodResponse{
			ID:          paymentMethod.ID,
			Name:        paymentMethod.Name,
			Code:        paymentMethod.Code,
			Description: paymentMethod.Description,
			IsActive:    paymentMethod.IsActive,
		}
	}

	for _, seat := range booking.BookingSeats {
		seatResponse := model.BookingSeatResponse{
			ID:     seat.ID,
			SeatID: seat.SeatID,
			Price:  seat.Price,
		}
		response.Seats = append(response.Seats, seatResponse)
	}

	return response
}
