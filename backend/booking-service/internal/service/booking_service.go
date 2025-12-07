package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"bus-booking/booking-service/internal/client"
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
)

type BookingService interface {
	CreateBooking(ctx context.Context, req *model.CreateBookingRequest, userID uuid.UUID) (*model.BookingResponse, error)
	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error)
	GetUserBookings(ctx context.Context, req model.PaginationRequest, userID uuid.UUID) ([]*model.BookingResponse, int64, error)
	GetTripBookings(ctx context.Context, req model.PaginationRequest, tripID uuid.UUID) ([]*model.BookingResponse, int64, error)
	CancelBooking(ctx context.Context, id uuid.UUID, userID uuid.UUID, reason string) error
	UpdateBookingStatus(ctx context.Context, id uuid.UUID, status string) error
	CreatePayment(ctx context.Context, bookingID uuid.UUID, buyerInfo *model.BuyerInfo) (*client.PaymentLinkResponse, error)
	UpdatePaymentStatus(ctx context.Context, bookingID uuid.UUID, paymentStatus, bookingStatus, paymentOrderID string) error
}

type bookingServiceImpl struct {
	bookingRepo   repository.BookingRepository
	paymentClient client.PaymentClient
	tripClient    client.TripClient
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	paymentClient client.PaymentClient,
	tripClient client.TripClient,
) BookingService {
	return &bookingServiceImpl{
		bookingRepo:   bookingRepo,
		paymentClient: paymentClient,
		tripClient:    tripClient,
	}
}

func (s *bookingServiceImpl) CreateBooking(ctx context.Context, req *model.CreateBookingRequest, userID uuid.UUID) (*model.BookingResponse, error) {
	tripData, err := s.tripClient.GetTrip(ctx, req.TripID)
	if err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to get trip details: %v", err))
	}

	// if !tripData.IsBookable() {
	// 	return nil, ginext.NewBadRequestError("trip is not available for booking")
	// }

	seatAvailability, err := s.bookingRepo.CheckSeatAvailability(ctx, req.TripID, req.SeatIDs)
	if err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to check seat availability: %v", err))
	}

	for seatID, isBooked := range seatAvailability {
		if isBooked {
			return nil, ginext.NewBadRequestError(fmt.Sprintf("seat %s is already booked", seatID))
		}
	}

	// 3. Get seat metadata from trip service (for pricing and display info)
	validatedSeats, err := s.tripClient.GetSeatsMetadata(ctx, req.TripID, req.SeatIDs)
	if err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	// 4. Calculate total amount
	totalAmount := s.tripClient.CalculateTotalPrice(tripData.BasePrice, validatedSeats)

	// 5. Create booking
	booking := &model.Booking{
		BookingReference: s.generateBookingReference(),
		TripID:           req.TripID,
		UserID:           userID,
		TotalAmount:      totalAmount,
		Status:           model.BookingStatusPending,
		PaymentStatus:    model.PaymentStatusPending,
		Notes:            req.Notes,
	}

	// Set expiration (15 minutes for pending bookings)
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	booking.ExpiresAt = &expiresAt

	// 6. Create booking seats
	for seatID, seat := range validatedSeats {
		bookingSeat := model.BookingSeat{
			SeatID:          seatID,
			SeatNumber:      seat.SeatNumber,
			SeatType:        seat.SeatType, // Raw string now (standard, vip, sleeper)
			Floor:           seat.Floor,
			Price:           seat.CalculateSeatPrice(tripData.BasePrice),
			PriceMultiplier: seat.PriceMultiplier,
		}
		booking.BookingSeats = append(booking.BookingSeats, bookingSeat)
	}

	// 7. Save to database
	if err := s.bookingRepo.CreateBooking(ctx, booking); err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to create booking: %v", err))
	}

	// 8. Return response
	return s.toBookingResponse(booking), nil
}

// GetBookingByID retrieves a booking by ID
func (s *bookingServiceImpl) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toBookingResponse(booking), nil
}

// GetUserBookings retrieves bookings for a user with pagination
func (s *bookingServiceImpl) GetUserBookings(ctx context.Context, req model.PaginationRequest, userID uuid.UUID) ([]*model.BookingResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	bookings, total, err := s.bookingRepo.GetBookingsByUserID(ctx, userID, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	var bookingResponses []*model.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, s.toBookingResponse(booking))
	}

	return bookingResponses, total, nil
}

// CancelBooking cancels a booking
func (s *bookingServiceImpl) CancelBooking(ctx context.Context, id uuid.UUID, userID uuid.UUID, reason string) error {
	booking, err := s.bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		return err
	}

	if booking.UserID != userID {
		return ginext.NewForbiddenError("you don't have permission to cancel this booking")
	}

	if booking.Status == model.BookingStatusCancelled {
		return ginext.NewBadRequestError("booking is already cancelled")
	}

	if booking.Status == model.BookingStatusConfirmed {
		return ginext.NewBadRequestError("cannot cancel confirmed booking")
	}

	return s.bookingRepo.CancelBooking(ctx, id, reason)
}

// GetTripBookings retrieves all bookings for a trip with pagination
func (s *bookingServiceImpl) GetTripBookings(ctx context.Context, req model.PaginationRequest, tripID uuid.UUID) ([]*model.BookingResponse, int64, error) {
	bookings, total, err := s.bookingRepo.GetTripBookings(ctx, tripID, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*model.BookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		responses = append(responses, s.toBookingResponse(booking))
	}

	return responses, total, nil
}

// UpdateBookingStatus updates the status of a booking (admin only)
func (s *bookingServiceImpl) UpdateBookingStatus(ctx context.Context, id uuid.UUID, status string) error {
	// Validate status
	bookingStatus := model.BookingStatus(status)
	if !bookingStatus.IsValid() {
		return ginext.NewBadRequestError("invalid booking status")
	}

	return s.bookingRepo.UpdateStatus(ctx, id, bookingStatus)
}

// CreatePayment creates a payment link for the booking
func (s *bookingServiceImpl) CreatePayment(ctx context.Context, bookingID uuid.UUID, buyerInfo *model.BuyerInfo) (*client.PaymentLinkResponse, error) {
	// Get booking details
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, ginext.NewNotFoundError("booking not found")
	}

	// Validate booking status
	if booking.Status != model.BookingStatusPending {
		return nil, ginext.NewBadRequestError("booking is not in pending status")
	}

	// Check if booking has expired
	if booking.ExpiresAt != nil && time.Now().After(*booking.ExpiresAt) {
		return nil, ginext.NewBadRequestError("booking has expired")
	}

	// Create payment link request
	paymentReq := &client.CreatePaymentLinkRequest{
		BookingID:     bookingID,
		Amount:        booking.TotalAmount,
		Currency:      "VND",
		PaymentMethod: "PAYOS",
		Description:   fmt.Sprintf("Thanh toán vé %s", booking.BookingReference),
		BuyerName:     buyerInfo.Name,
		BuyerEmail:    buyerInfo.Email,
		BuyerPhone:    buyerInfo.Phone,
	}

	// Call payment service
	paymentResp, err := s.paymentClient.CreatePaymentLink(ctx, paymentReq)
	if err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to create payment link: %v", err))
	}

	// Update booking with payment order ID
	booking.PaymentOrderID = fmt.Sprintf("%d", paymentResp.OrderCode)
	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return nil, ginext.NewInternalServerError("failed to update booking with payment info")
	}

	return paymentResp, nil
}

// UpdatePaymentStatus updates booking payment status (called by payment service)
func (s *bookingServiceImpl) UpdatePaymentStatus(ctx context.Context, bookingID uuid.UUID, paymentStatus, bookingStatus, paymentOrderID string) error {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return ginext.NewNotFoundError("booking not found")
	}

	// Update payment status
	booking.PaymentStatus = model.PaymentStatus(paymentStatus)
	booking.Status = model.BookingStatus(bookingStatus)
	booking.PaymentOrderID = paymentOrderID

	// If payment is successful, set confirmed time
	if paymentStatus == string(model.PaymentStatusPaid) && bookingStatus == string(model.BookingStatusConfirmed) {
		now := time.Now()
		booking.ConfirmedAt = &now
	}

	return s.bookingRepo.UpdateBooking(ctx, booking)
}

// Helper methods

func (s *bookingServiceImpl) toBookingResponse(booking *model.Booking) *model.BookingResponse {
	resp := &model.BookingResponse{
		ID:               booking.ID,
		BookingReference: booking.BookingReference,
		TripID:           booking.TripID,
		UserID:           booking.UserID,
		TotalAmount:      booking.TotalAmount,
		Status:           string(booking.Status),
		PaymentStatus:    string(booking.PaymentStatus),
		PaymentOrderID:   booking.PaymentOrderID,
		Notes:            booking.Notes,
		ExpiresAt:        booking.ExpiresAt,
		ConfirmedAt:      booking.ConfirmedAt,
		CancelledAt:      booking.CancelledAt,
		CreatedAt:        booking.CreatedAt,
		UpdatedAt:        booking.UpdatedAt,
	}

	// Map seats
	for _, seat := range booking.BookingSeats {
		resp.Seats = append(resp.Seats, model.BookingSeatResponse{
			ID:              seat.ID,
			SeatID:          seat.SeatID,
			SeatNumber:      seat.SeatNumber,
			SeatType:        seat.SeatType,
			Floor:           seat.Floor,
			Price:           seat.Price,
			PriceMultiplier: seat.PriceMultiplier,
		})
	}

	return resp
}

func (s *bookingServiceImpl) generateBookingReference() string {
	now := time.Now().UTC()
	dateStr := now.Format("060102") // YYMMDD

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomPart := make([]byte, 4)
	for i := range randomPart {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			n = big.NewInt(0)
		}
		randomPart[i] = charset[n.Int64()]
	}

	return "BK" + dateStr + string(randomPart)
}
