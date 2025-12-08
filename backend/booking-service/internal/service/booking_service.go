package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"bus-booking/booking-service/internal/client"
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/model/payment"
	"bus-booking/booking-service/internal/model/trip"
	"bus-booking/booking-service/internal/model/user"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type BookingService interface {
	CreateBooking(ctx context.Context, req *model.CreateBookingRequest, userID uuid.UUID) (*model.BookingResponse, error)
	CreateGuestBooking(ctx context.Context, req *model.CreateGuestBookingRequest) (*model.BookingResponse, error)

	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error)
	GetBookingByReference(ctx context.Context, reference string, email string) (*model.BookingResponse, error)
	GetUserBookings(ctx context.Context, req model.PaginationRequest, userID uuid.UUID) ([]*model.BookingResponse, int64, error)
	GetTripBookings(ctx context.Context, req model.PaginationRequest, tripID uuid.UUID) ([]*model.BookingResponse, int64, error)

	CancelBooking(ctx context.Context, id uuid.UUID, userID uuid.UUID, reason string) error
	UpdateBookingStatus(ctx context.Context, id uuid.UUID, status string) error

	UpdatePaymentStatus(ctx context.Context, req *model.UpdatePaymentStatusRequest, bookingID uuid.UUID) error
	GetSeatStatus(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) ([]model.SeatStatusItem, error)
}

type bookingServiceImpl struct {
	bookingRepo   repository.BookingRepository
	paymentClient client.PaymentClient
	tripClient    client.TripClient
	userClient    client.UserClient
	// notificationQueue queue.NotificationQueue // TODO: Uncomment when notification service is ready
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	paymentClient client.PaymentClient,
	tripClient client.TripClient,
	userClient client.UserClient,
	// notificationQueue queue.NotificationQueue, // TODO: Uncomment when notification service is ready
) BookingService {
	return &bookingServiceImpl{
		bookingRepo:   bookingRepo,
		paymentClient: paymentClient,
		tripClient:    tripClient,
		userClient:    userClient,
		// notificationQueue: notificationQueue, // TODO: Uncomment when notification service is ready
	}
}

func (s *bookingServiceImpl) CreateBooking(ctx context.Context, req *model.CreateBookingRequest, userID uuid.UUID) (*model.BookingResponse, error) {
	// 1. Validate seat IDs
	seatAvailability, err := s.checkSeatAvailability(ctx, req.TripID, req.SeatIDs)
	if err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to check seat availability: %v", err))
	}
	if !seatAvailability {
		return nil, ginext.NewBadRequestError("one or more selected seats are already booked")
	}

	// 2. Fetch trip and seat details concurrently
	var (
		tripData *trip.Trip
		seats    []trip.Seat
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		tripData, err = s.tripClient.GetTripByID(gCtx, req.TripID)
		if err != nil {
			return fmt.Errorf("failed to get trip data: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		seats, err = s.tripClient.ListSeatsByIDs(gCtx, req.SeatIDs)
		if err != nil {
			return fmt.Errorf("failed to list seats: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// 4. Calculate total amount
	totalAmount := s.CalculateTotalPrice(tripData.BasePrice, seats)

	// 5. Create booking
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	booking := &model.Booking{
		BookingReference:  s.generateBookingReference(),
		TripID:            req.TripID,
		UserID:            userID,
		TotalAmount:       totalAmount,
		Status:            model.BookingStatusPending,
		TransactionStatus: payment.TransactionStatusPending,
		Notes:             req.Notes,
		ExpiresAt:         &expiresAt,
	}

	// 6. Create booking seats
	for _, seat := range seats {
		bookingSeat := model.BookingSeat{
			SeatID:          seat.ID,
			SeatNumber:      seat.SeatNumber,
			SeatType:        seat.SeatType,
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

	// 8. Create payment link
	transaction, err := s.paymentClient.CreatePaymentLink(ctx, &payment.CreatePaymentLinkRequest{
		BookingID:     booking.ID,
		Amount:        totalAmount,
		Currency:      payment.CurrencyVND,
		PaymentMethod: payment.PaymentMethodPayOS,
		Description:   fmt.Sprintf("Don hang %s", booking.BookingReference),
	})
	if err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to create payment link: %v", err))
	}

	booking.TransactionID = transaction.ID
	booking.TransactionStatus = transaction.Status

	go func() {
		if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
			log.Printf("failed to update booking: %v", err)
		}
	}()

	// 9. Push notification to queue (khi có notification service)
	// TODO: Uncomment khi notification service đã ready
	/*
		if s.notificationQueue != nil {
			emailNotif := &queue.EmailNotification{
				To:           "user@example.com", // TODO: Get from user service
				Subject:      "Xác nhận đặt vé",
				TemplateName: "booking_confirmation",
				TemplateData: map[string]interface{}{
					"booking_reference": booking.BookingReference,
					"total_amount":      booking.TotalAmount,
					"trip_id":           booking.TripID.String(),
					"seat_numbers":      s.getSeatNumbers(booking.BookingSeats),
				},
				Priority: 1, // High priority
			}

			// Push to queue (non-blocking, errors logged only)
			if err := s.notificationQueue.PushEmailNotification(ctx, emailNotif); err != nil {
				log.Error().Err(err).Msg("Failed to push email notification to queue")
			}
		}
	*/

	// 9. Return response
	return s.toBookingResponse(booking), nil
}

// CreateGuestBooking creates a booking for guest users (without authentication)
func (s *bookingServiceImpl) CreateGuestBooking(ctx context.Context, req *model.CreateGuestBookingRequest) (*model.BookingResponse, error) {
	// 1. Validate contact information
	if req.Email == "" && req.Phone == "" {
		return nil, ginext.NewBadRequestError("phải cung cấp email hoặc số điện thoại")
	}

	// 2. Create or get guest account from user service
	guest, err := s.userClient.CreateGuest(ctx, &user.CreateGuestRequest{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
	})
	if err != nil {
		return nil, ginext.NewInternalServerError("không thể tạo tài khoản khách")
	}

	// 3. Use existing CreateBooking logic with guest user ID
	return s.CreateBooking(ctx, &model.CreateBookingRequest{
		TripID:  req.TripID,
		SeatIDs: req.SeatIDs,
		Notes:   req.Notes,
	}, guest.ID)
}

func (s *bookingServiceImpl) checkSeatAvailability(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) (bool, error) {
	bookedSeatIDs, err := s.bookingRepo.GetBookedSeatIDs(ctx, tripID)
	if err != nil {
		return false, err
	}

	bookedMap := make(map[uuid.UUID]bool)
	for _, bookedID := range bookedSeatIDs {
		bookedMap[bookedID] = true
	}

	for _, seatID := range seatIDs {
		if bookedMap[seatID] {
			return false, nil
		}
	}

	return true, nil
}

func (s *bookingServiceImpl) CalculateTotalPrice(basePrice float64, seats []trip.Seat) int {
	total := 0.0
	for _, seat := range seats {
		total += seat.CalculateSeatPrice(basePrice)
	}
	return int(total)
}

func (s *bookingServiceImpl) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get payment status from payment service
	transaction, err := s.paymentClient.GetTransactionByID(ctx, booking.TransactionID)
	if err != nil {
		return nil, ginext.NewInternalServerError("không thể lấy thông tin giao dịch")
	}

	resp := s.toBookingResponse(booking)
	resp.Transaction = transaction
	return resp, nil
}

// GetBookingByReference retrieves booking by reference number for guest lookup
func (s *bookingServiceImpl) GetBookingByReference(ctx context.Context, reference string, email string) (*model.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByReference(ctx, reference)
	if err != nil {
		return nil, ginext.NewNotFoundError("Booking not found with this reference number")
	}

	// For guest bookings, we trust the reference number is unique enough
	// In production, you might want to verify email matches the user's email
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

// UpdatePaymentStatus updates booking payment status (called by payment service)
func (s *bookingServiceImpl) UpdatePaymentStatus(ctx context.Context, req *model.UpdatePaymentStatusRequest, bookingID uuid.UUID) error {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return ginext.NewNotFoundError("booking not found")
	}

	// Update payment status
	booking.TransactionStatus = req.PaymentStatus
	booking.Status = req.BookingStatus

	// If payment is successful, set confirmed time
	if req.PaymentStatus == payment.TransactionStatusPaid && req.BookingStatus == model.BookingStatusConfirmed {
		now := time.Now()
		booking.ConfirmedAt = &now
	}

	return s.bookingRepo.UpdateBooking(ctx, booking)
}

// Helper methods

func (s *bookingServiceImpl) toBookingResponse(booking *model.Booking) *model.BookingResponse {
	resp := &model.BookingResponse{
		ID:                booking.ID,
		CreatedAt:         booking.CreatedAt,
		UpdatedAt:         booking.UpdatedAt,
		BookingReference:  booking.BookingReference,
		TripID:            booking.TripID,
		UserID:            booking.UserID,
		TotalAmount:       booking.TotalAmount,
		Status:            booking.Status,
		TransactionStatus: booking.TransactionStatus,
		TransactionID:     booking.TransactionID,
		Notes:             booking.Notes,
		ExpiresAt:         booking.ExpiresAt,
		ConfirmedAt:       booking.ConfirmedAt,
		CancelledAt:       booking.CancelledAt,
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

// example: BK251208AB123
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

func (s *bookingServiceImpl) GetSeatStatus(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) ([]model.SeatStatusItem, error) {
	if len(seatIDs) == 0 {
		return []model.SeatStatusItem{}, nil
	}

	bookedSeatIDs, err := s.bookingRepo.GetBookedSeatIDs(ctx, tripID)
	if err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to get booked seats: %v", err))
	}

	bookedMap := make(map[uuid.UUID]bool)
	for _, bookedID := range bookedSeatIDs {
		bookedMap[bookedID] = true
	}

	result := make([]model.SeatStatusItem, len(seatIDs))
	for i, seatID := range seatIDs {
		result[i] = model.SeatStatusItem{
			SeatID:   seatID,
			IsBooked: bookedMap[seatID],
			IsLocked: false,
		}
	}

	return result, nil
}
