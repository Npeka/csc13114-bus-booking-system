package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"bus-booking/booking-service/internal/client"
	"bus-booking/booking-service/internal/constants"
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/model/payment"
	"bus-booking/booking-service/internal/model/trip"
	"bus-booking/booking-service/internal/model/user"
	"bus-booking/booking-service/internal/repository"
	"bus-booking/shared/ginext"
	"bus-booking/shared/queue"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type BookingService interface {
	CreateBooking(ctx context.Context, req *model.CreateBookingRequest, userID uuid.UUID) (*model.BookingResponse, error)
	CreateGuestBooking(ctx context.Context, req *model.CreateGuestBookingRequest) (*model.BookingResponse, error)
	UpdateBookingStatus(ctx context.Context, req *model.UpdateBookingStatusRequest, bookingID uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error)
	GetByReference(ctx context.Context, reference string, email string) (*model.BookingResponse, error)
	GetUserBookings(ctx context.Context, req model.GetUserBookingsRequest, userID uuid.UUID) ([]*model.BookingResponse, int64, error)
	GetTripBookings(ctx context.Context, req model.PaginationRequest, tripID uuid.UUID) ([]*model.BookingResponse, int64, error)

	CancelBooking(ctx context.Context, id uuid.UUID, reason string) error
	RetryPayment(ctx context.Context, bookingID uuid.UUID) (*model.BookingResponse, error)

	GetSeatStatus(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) ([]model.SeatStatusItem, error)
	ExpireBooking(ctx context.Context, bookingID uuid.UUID) error
}

type bookingServiceImpl struct {
	bookingRepo        repository.BookingRepository
	paymentClient      client.PaymentClient
	tripClient         client.TripClient
	userClient         client.UserClient
	delayedQueue       queue.DelayedQueueManager
	notificationClient client.NotificationClient
	seatLockService    SeatLockService
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	paymentClient client.PaymentClient,
	tripClient client.TripClient,
	userClient client.UserClient,
	notificationClient client.NotificationClient,
	delayedQueue queue.DelayedQueueManager,
	seatLockService SeatLockService,
) BookingService {
	return &bookingServiceImpl{
		bookingRepo:        bookingRepo,
		paymentClient:      paymentClient,
		tripClient:         tripClient,
		userClient:         userClient,
		notificationClient: notificationClient,
		delayedQueue:       delayedQueue,
		seatLockService:    seatLockService,
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
		tripData, err = s.tripClient.GetTripByID(gCtx, trip.GetTripByIDRequest{}, req.TripID)
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
	totalAmount := s.calculateTotalPrice(tripData.BasePrice, seats)

	// 5. Create booking
	expiresAt := time.Now().UTC().Add(constants.BookingPaymentTimeout)
	booking := &model.Booking{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		BookingReference:  s.generateBookingReference(),
		TripID:            req.TripID,
		UserID:            userID,
		TotalAmount:       totalAmount,
		Status:            model.BookingStatusPending,
		TransactionStatus: payment.TransactionStatusPending,
		TransactionID:     uuid.New(),
		Notes:             req.Notes,
		ExpiresAt:         &expiresAt,
	}

	// 6. Create booking seats
	for _, seat := range seats {
		booking.BookingSeats = append(booking.BookingSeats, model.BookingSeat{
			SeatID:          seat.ID,
			SeatNumber:      seat.SeatNumber,
			SeatType:        seat.SeatType,
			Floor:           seat.Floor,
			Price:           seat.CalculateSeatPrice(tripData.BasePrice),
			PriceMultiplier: seat.PriceMultiplier,
		})
	}

	// 7. Save to database
	if err := s.bookingRepo.CreateBooking(ctx, booking); err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to create booking: %v", err))
	}

	// 8. Create payment link
	transaction, err := s.paymentClient.CreateTransaction(ctx, &payment.CreateTransactionRequest{
		ID:            booking.TransactionID,
		BookingID:     booking.ID,
		Amount:        totalAmount,
		Currency:      payment.CurrencyVND,
		PaymentMethod: payment.PaymentMethodPayOS,
		Description:   fmt.Sprintf("Don hang %s", booking.BookingReference),
		ExpiresAt:     expiresAt,
	})
	if err != nil {
		// Payment creation failed - update booking status to FAILED
		log.Error().Err(err).
			Str("booking_id", booking.ID.String()).
			Str("transaction_id", booking.TransactionID.String()).
			Msg("Payment link creation failed")

		booking.Status = model.BookingStatusFailed
		booking.TransactionStatus = payment.TransactionStatusFailed

		if updateErr := s.bookingRepo.UpdateBooking(ctx, booking); updateErr != nil {
			log.Error().Err(updateErr).
				Str("booking_id", booking.ID.String()).
				Msg("Failed to update booking status after payment failure")
		}

		// Return booking with error info - user can retry payment
		resp := s.toBookingResponse(booking)
		resp.Transaction = &payment.TransactionResponse{
			ID:     booking.TransactionID,
			Status: payment.TransactionStatusFailed,
		}
		return resp, nil
	}

	// Payment link created successfully - send pending email and schedule expiration
	// 9. Send pending email and schedule expiration job
	if transaction != nil && transaction.CheckoutURL != "" {
		go func() {
			// Create a detached context with timeout for background task
			bgCtx, cancel := context.WithTimeout(context.Background(), constants.BackgroundTaskTimeout)
			defer cancel()

			// Send pending email
			s.sendBookingPendingEmail(bgCtx, booking, tripData, transaction.CheckoutURL)

			// Schedule expiration in delayed queue
			item := &queue.DelayedItem{
				Payload: booking.ID,
			}
			if err := s.delayedQueue.Schedule(bgCtx, constants.QueueNameBookingExpiry, item, *booking.ExpiresAt); err != nil {
				log.Error().Err(err).
					Str("booking_id", booking.ID.String()).
					Time("expires_at", *booking.ExpiresAt).
					Msg("Failed to schedule booking expiration")
			} else {
				log.Info().
					Str("booking_id", booking.ID.String()).
					Time("expires_at", *booking.ExpiresAt).
					Msg("Successfully scheduled booking expiration")
			}
		}()
	}

	resp := s.toBookingResponse(booking)
	resp.Transaction = transaction
	return resp, nil
}

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

func (s *bookingServiceImpl) calculateTotalPrice(basePrice float64, seats []trip.Seat) int {
	total := 0.0
	for _, seat := range seats {
		total += seat.CalculateSeatPrice(basePrice)
	}
	return int(total)
}

// example: BK251208AB123
func (s *bookingServiceImpl) generateBookingReference() string {
	now := time.Now().UTC()
	dateStr := now.Format(constants.DateFormatBookingReference)

	randomPart := make([]byte, constants.BookingReferenceRandomLength)
	for i := range randomPart {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(constants.BookingReferenceCharset))))
		if err != nil {
			n = big.NewInt(0)
		}
		randomPart[i] = constants.BookingReferenceCharset[n.Int64()]
	}

	return constants.BookingReferencePrefix + dateStr + string(randomPart)
}

func (s *bookingServiceImpl) UpdateBookingStatus(ctx context.Context, req *model.UpdateBookingStatusRequest, bookingID uuid.UUID) error {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return ginext.NewNotFoundError("booking not found")
	}

	// Update payment status
	booking.TransactionStatus = req.TransactionStatus
	switch booking.TransactionStatus {
	case payment.TransactionStatusPending:
		booking.Status = model.BookingStatusPending

	case payment.TransactionStatusPaid:
		booking.Status = model.BookingStatusConfirmed

		// Send Confirmation Email
		go func() {
			bgCtx, cancel := context.WithTimeout(context.Background(), constants.BackgroundTaskTimeout)
			defer cancel()
			s.sendBookingConfirmationEmail(bgCtx, bookingID)
		}()

		// Schedule Trip Reminder (2 hours before departure)
		// Fetch trip to get departure time
		trip, err := s.tripClient.GetTripByID(ctx, trip.GetTripByIDRequest{
			PreLoadRoute: true,
		}, booking.TripID)
		if err != nil {
			// Log error but don't fail booking update? Or fail?
			// Failing here stops status update which is bad for payment flow.
			// Just log error.
			fmt.Printf("Failed to fetch trip for reminder scheduling: %v\n", err)
		} else {
			// Schedule for 2 hours before departure
			// If already past or less than 2 hours?
			// Calculate executeAt
			executeAt := trip.DepartureTime.Add(-constants.TripReminderBeforeDeparture)

			// If executeAt is in past, schedule for now? Or skip?
			// If now > departure, skip.
			// If now > executeAt, send immediately (schedule for now).
			if time.Now().Before(trip.DepartureTime) {
				reminderPayload := &queue.DelayedItem{
					Type:    "trip_reminder",
					Payload: booking.ID,
				}
				if err := s.delayedQueue.Schedule(ctx, constants.QueueNameTripReminder, reminderPayload, executeAt); err != nil {
					fmt.Printf("Failed to schedule trip reminder: %v\n", err)
				}
			}
		}

	case payment.TransactionStatusCancelled:
		booking.Status = model.BookingStatusCancelled
		now := time.Now()
		booking.CancelledAt = &now

	case payment.TransactionStatusExpired:
		booking.Status = model.BookingStatusExpired

	case payment.TransactionStatusFailed:
		booking.Status = model.BookingStatusFailed
		// Send Failure Email
		go func() {
			bgCtx, cancel := context.WithTimeout(context.Background(), constants.BackgroundTaskTimeout)
			defer cancel()
			s.sendBookingFailureEmail(bgCtx, bookingID, "Thanh toán thất bại")
		}()

	default:
		booking.Status = model.BookingStatusFailed
		// Send Failure Email
		go func() {
			bgCtx, cancel := context.WithTimeout(context.Background(), constants.BackgroundTaskTimeout)
			defer cancel()
			s.sendBookingFailureEmail(bgCtx, bookingID, "Lỗi không xác định")
		}()
	}

	return s.bookingRepo.UpdateBooking(ctx, booking)
}

func (s *bookingServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.BookingResponse, error) {
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

// GetByReference retrieves booking by reference number for guest lookup
func (s *bookingServiceImpl) GetByReference(ctx context.Context, reference string, email string) (*model.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByReference(ctx, reference)
	if err != nil {
		return nil, ginext.NewNotFoundError("Booking not found with this reference number")
	}

	// For guest bookings, we trust the reference number is unique enough
	// In production, you might want to verify email matches the user's email
	return s.toBookingResponse(booking), nil
}

// GetUserBookings retrieves bookings for a user with pagination
func (s *bookingServiceImpl) GetUserBookings(ctx context.Context, req model.GetUserBookingsRequest, userID uuid.UUID) ([]*model.BookingResponse, int64, error) {
	offset := (req.Page - 1) * req.PageSize
	bookings, total, err := s.bookingRepo.GetBookingsByUserID(ctx, userID, req.Status, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*model.BookingResponse, len(bookings))

	// Collect unique trip IDs
	uniqueTripIDs := make(map[uuid.UUID]bool)
	for _, booking := range bookings {
		uniqueTripIDs[booking.TripID] = true
	}

	// Convert to slice
	tripIDs := make([]uuid.UUID, 0, len(uniqueTripIDs))
	for tripID := range uniqueTripIDs {
		tripIDs = append(tripIDs, tripID)
	}

	// Batch fetch all trips with preloads in a single API call
	var tripDataMap map[uuid.UUID]*trip.Trip
	if len(tripIDs) > 0 {
		trips, err := s.tripClient.GetTripsByIDs(ctx, trip.GetTripByIDRequest{
			PreLoadRoute: true,
			PreloadBus:   true,
		}, tripIDs)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to batch fetch trip data for bookings")
		} else {
			// Build map for quick lookup
			tripDataMap = make(map[uuid.UUID]*trip.Trip, len(trips))
			for i := range trips {
				tripDataMap[trips[i].ID] = &trips[i]
			}
		}
	}

	// Build responses with trip data
	for i, booking := range bookings {
		resp := s.toBookingResponse(booking)

		// Add trip info if available
		if tripData, ok := tripDataMap[booking.TripID]; ok {
			resp.Trip = &model.TripBasicInfo{
				Origin:        tripData.Route.Origin,
				Destination:   tripData.Route.Destination,
				DepartureTime: tripData.DepartureTime,
				BusName:       tripData.Bus.Model,
			}
		}

		responses[i] = resp
	}

	return responses, total, nil
}

// CancelBooking cancels a booking
func (s *bookingServiceImpl) CancelBooking(ctx context.Context, id uuid.UUID, reason string) error {
	booking, err := s.bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		return err
	}

	if booking.Status == model.BookingStatusCancelled {
		return ginext.NewBadRequestError("booking is already cancelled")
	}

	if booking.Status == model.BookingStatusConfirmed {
		return ginext.NewBadRequestError("cannot cancel confirmed booking")
	}

	// Cancel the booking first
	if err := s.bookingRepo.CancelBooking(ctx, id, reason); err != nil {
		return err
	}

	// Try to cancel payment if transaction exists
	transaction, err := s.paymentClient.CancelTransaction(ctx, booking.TransactionID)
	if err != nil {
		// Log error but don't fail booking cancellation
		log.Error().
			Err(err).
			Str("booking_id", id.String()).
			Str("transaction_id", booking.TransactionID.String()).
			Msg("Failed to cancel payment, but booking is cancelled")
		// Continue - booking is already cancelled
	} else {
		// Fetch fresh booking data to avoid overwriting the CANCELLED status
		updatedBooking, err := s.bookingRepo.GetBookingByID(ctx, id)
		if err != nil {
			log.Error().
				Err(err).
				Str("booking_id", id.String()).
				Msg("Failed to fetch booking after payment cancellation")
			// Continue - booking is cancelled, just couldn't update transaction status
			return nil
		}

		// Update transaction status only
		updatedBooking.TransactionStatus = transaction.Status
		if err := s.bookingRepo.UpdateBooking(ctx, updatedBooking); err != nil {
			log.Error().
				Err(err).
				Str("booking_id", id.String()).
				Msg("Failed to update booking transaction status after payment cancellation")
			// Continue - payment is cancelled, just status update failed
		} else {
			log.Info().
				Str("booking_id", id.String()).
				Str("transaction_id", booking.TransactionID.String()).
				Str("new_transaction_status", string(transaction.Status)).
				Msg("Successfully cancelled payment and updated booking transaction status")
		}
	}

	return nil
}

// RetryPayment creates a new payment link for a failed or expired booking
func (s *bookingServiceImpl) RetryPayment(ctx context.Context, bookingID uuid.UUID) (*model.BookingResponse, error) {
	// 1. Get booking
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, ginext.NewNotFoundError("booking not found")
	}

	// 2. Validate booking is in retryable state
	if booking.Status != model.BookingStatusFailed && booking.Status != model.BookingStatusExpired {
		return nil, ginext.NewBadRequestError("booking is not in a retryable state")
	}

	// 3. Check that booking hasn't expired beyond grace period (60 minutes)
	if booking.ExpiresAt != nil {
		gracePeriod := booking.ExpiresAt.Add(constants.BookingRetryGracePeriod)
		if time.Now().UTC().After(gracePeriod) {
			return nil, ginext.NewBadRequestError("booking has expired beyond retry period")
		}
	}

	// 4. Get trip data for amount validation
	tripData, err := s.tripClient.GetTripByID(ctx, trip.GetTripByIDRequest{
		PreLoadRoute: true,
	}, booking.TripID)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to fetch trip data")
	}

	// 5. Create new transaction ID and expiration
	newTransactionID := uuid.New()
	expiresAt := time.Now().UTC().Add(constants.BookingPaymentTimeout)

	// 6. Create new payment link
	transaction, err := s.paymentClient.CreateTransaction(ctx, &payment.CreateTransactionRequest{
		ID:            newTransactionID,
		BookingID:     booking.ID,
		Amount:        booking.TotalAmount,
		Currency:      payment.CurrencyVND,
		PaymentMethod: payment.PaymentMethodPayOS,
		Description:   fmt.Sprintf("Don hang %s (Thu lai)", booking.BookingReference),
		ExpiresAt:     expiresAt,
	})
	if err != nil {
		log.Error().Err(err).
			Str("booking_id", booking.ID.String()).
			Msg("Failed to create retry payment link")
		return nil, ginext.NewInternalServerError("failed to create payment link")
	}

	// 7. Update booking with new transaction and expiry
	booking.TransactionID = newTransactionID
	booking.TransactionStatus = payment.TransactionStatusPending
	booking.Status = model.BookingStatusPending
	booking.ExpiresAt = &expiresAt

	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		log.Error().Err(err).
			Str("booking_id", booking.ID.String()).
			Msg("Failed to update booking with new transaction")
		return nil, ginext.NewInternalServerError("failed to update booking")
	}

	// 8. Schedule expiration in delayed queue
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		item := &queue.DelayedItem{
			Payload: booking.ID,
		}
		if err := s.delayedQueue.Schedule(bgCtx, constants.QueueNameBookingExpiry, item, expiresAt); err != nil {
			log.Error().Err(err).
				Str("booking_id", booking.ID.String()).
				Time("expires_at", expiresAt).
				Msg("Failed to schedule booking expiration on retry")
		} else {
			log.Info().
				Str("booking_id", booking.ID.String()).
				Time("expires_at", expiresAt).
				Msg("Successfully scheduled booking expiration on retry")
		}
	}()

	// 9. Send pending email with new payment link
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), constants.BackgroundTaskTimeout)
		defer cancel()
		s.sendBookingPendingEmail(bgCtx, booking, tripData, transaction.CheckoutURL)
	}()

	log.Info().
		Str("booking_id", booking.ID.String()).
		Str("new_transaction_id", newTransactionID.String()).
		Msg("Payment retry successful")

	resp := s.toBookingResponse(booking)
	resp.Transaction = transaction
	return resp, nil
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

func (s *bookingServiceImpl) GetSeatStatus(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) ([]model.SeatStatusItem, error) {
	if len(seatIDs) == 0 {
		return []model.SeatStatusItem{}, nil
	}

	// Get booked seats
	bookedSeatIDs, err := s.bookingRepo.GetBookedSeatIDs(ctx, tripID)
	if err != nil {
		return nil, ginext.NewInternalServerError(fmt.Sprintf("failed to get booked seats: %v", err))
	}

	bookedMap := make(map[uuid.UUID]bool)
	for _, bookedID := range bookedSeatIDs {
		bookedMap[bookedID] = true
	}

	// Get locked seats
	lockedSeatIDs, err := s.seatLockService.GetLockedSeats(ctx, tripID)
	if err != nil {
		// Log error but don't fail - locked seats check is not critical
		log.Warn().Err(err).Str("trip_id", tripID.String()).Msg("Failed to get locked seats, continuing without lock info")
		lockedSeatIDs = []uuid.UUID{}
	}

	lockedMap := make(map[uuid.UUID]bool)
	for _, lockedID := range lockedSeatIDs {
		lockedMap[lockedID] = true
	}

	result := make([]model.SeatStatusItem, len(seatIDs))
	for i, seatID := range seatIDs {
		result[i] = model.SeatStatusItem{
			SeatID:   seatID,
			IsBooked: bookedMap[seatID],
			IsLocked: lockedMap[seatID],
		}
	}

	return result, nil
}

func (s *bookingServiceImpl) ExpireBooking(ctx context.Context, bookingID uuid.UUID) error {
	// 1. Get booking
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	// 2. Check if eligible for expiration
	if booking.Status != model.BookingStatusPending {
		return nil // Already processed (paid, cancelled, or failed)
	}

	// Double check expiry time (give 1 minute grace period)
	if booking.ExpiresAt != nil {
		graceDeadline := booking.ExpiresAt.Add(constants.BookingExpirationGracePeriod)
		now := time.Now().UTC()

		// Check if still within grace period
		if now.Before(graceDeadline) {
			// Re-schedule for after grace period ends
			log.Info().
				Str("booking_id", booking.ID.String()).
				Time("expires_at", *booking.ExpiresAt).
				Time("grace_deadline", graceDeadline).
				Msg("Booking still within grace period, re-scheduling expiration")

			// Schedule to execute right after grace period
			item := &queue.DelayedItem{
				Payload: booking.ID,
			}
			if err := s.delayedQueue.Schedule(ctx, constants.QueueNameBookingExpiry, item, graceDeadline); err != nil {
				log.Error().Err(err).
					Str("booking_id", booking.ID.String()).
					Time("grace_deadline", graceDeadline).
					Msg("Failed to re-schedule booking expiration after grace period")
				return err
			}

			return nil
		}
	}

	// 3. Grace period has passed - Update status to Expired
	booking.Status = model.BookingStatusExpired
	booking.TransactionStatus = payment.TransactionStatusExpired
	now := time.Now().UTC()
	booking.UpdatedAt = now

	// 4. Try to cancel payment transaction
	if booking.TransactionID != uuid.Nil {
		_, err := s.paymentClient.CancelTransaction(ctx, booking.TransactionID)
		if err != nil {
			// Log error but don't fail expiration - transaction will expire on PayOS side
			log.Warn().Err(err).
				Str("booking_id", booking.ID.String()).
				Str("transaction_id", booking.TransactionID.String()).
				Msg("Failed to cancel transaction, but booking will still expire")
		} else {
			log.Info().
				Str("booking_id", booking.ID.String()).
				Str("transaction_id", booking.TransactionID.String()).
				Msg("Successfully cancelled transaction")
		}
	}

	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return fmt.Errorf("failed to expire booking: %w", err)
	}

	log.Info().
		Str("booking_id", booking.ID.String()).
		Time("expires_at", *booking.ExpiresAt).
		Msg("Successfully expired booking after grace period")

	// Send Failure Email (Expired)
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), constants.BackgroundTaskTimeout)
		defer cancel()
		s.sendBookingFailureEmail(bgCtx, bookingID, "Hết hạn thanh toán")
	}()

	return nil
}

// Email Helper Methods

func (s *bookingServiceImpl) sendBookingPendingEmail(ctx context.Context, booking *model.Booking, trip *trip.Trip, paymentLink string) {
	userData, err := s.userClient.GetUserByID(ctx, booking.UserID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", booking.UserID.String()).
			Msg("Failed to fetch user for booking pending email")
		return
	}

	// Check if trip and route are populated
	if trip == nil || trip.Route == nil {
		log.Error().
			Str("booking_id", booking.ID.String()).
			Msg("Trip or Route is nil, cannot send booking pending email")
		return
	}

	// Route might not be populated if just checking tripData unless fetched with preloads?
	// tripClient.GetTripByID should return full trip details including Route.
	// Assuming Trip struct has Route populated.

	req := &client.BookingPendingRequest{
		Email:            userData.Email,
		Name:             userData.FullName,
		BookingReference: booking.BookingReference,
		From:             trip.Route.Origin,
		To:               trip.Route.Destination,
		DepartureTime:    trip.DepartureTime.Format(constants.DateTimeFormatDisplay),
		TotalAmount:      booking.TotalAmount,
		PaymentLink:      paymentLink,
	}

	if err := s.notificationClient.SendBookingPending(ctx, req); err != nil {
		log.Error().
			Err(err).
			Str("booking_id", booking.ID.String()).
			Msg("Failed to send booking pending email")
	}
}

func (s *bookingServiceImpl) sendBookingConfirmationEmail(ctx context.Context, bookingID uuid.UUID) {
	// Re-fetch everything to ensure fresh data
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		fmt.Printf("Failed to get booking for confirmation email: %v\n", err)
		return
	}

	user, err := s.userClient.GetUserByID(ctx, booking.UserID)
	if err != nil {
		fmt.Printf("Failed to get user for confirmation email: %v\n", err)
		return
	}

	trip, err := s.tripClient.GetTripByID(ctx, trip.GetTripByIDRequest{
		PreLoadRoute: true,
	}, booking.TripID)
	if err != nil {
		fmt.Printf("Failed to get trip for confirmation email: %v\n", err)
		return
	}

	req := &client.BookingConfirmationRequest{
		Email:            user.Email,
		Name:             user.FullName,
		BookingReference: booking.BookingReference,
		From:             trip.Route.Origin,
		To:               trip.Route.Destination,
		DepartureTime:    trip.DepartureTime.Format(constants.DateTimeFormatDisplay),
		SeatNumbers:      s.getSeatNumbersResults(booking.BookingSeats),
		TotalAmount:      booking.TotalAmount,
		TicketLink:       fmt.Sprintf("%s/booking/ticket/%s", constants.DefaultFrontendURL, booking.BookingReference),
	}

	if err := s.notificationClient.SendBookingConfirmation(ctx, req); err != nil {
		fmt.Printf("Failed to send booking confirmation email: %v\n", err)
	}
}

func (s *bookingServiceImpl) sendBookingFailureEmail(ctx context.Context, bookingID uuid.UUID, reason string) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		fmt.Printf("Failed to get booking for failure email: %v\n", err)
		return
	}

	user, err := s.userClient.GetUserByID(ctx, booking.UserID)
	if err != nil {
		fmt.Printf("Failed to get user for failure email: %v\n", err)
		return
	}

	trip, err := s.tripClient.GetTripByID(ctx, trip.GetTripByIDRequest{
		PreLoadRoute: true,
	}, booking.TripID)
	if err != nil {
		fmt.Printf("Failed to get trip for failure email: %v\n", err)
		return
	}

	req := &client.BookingFailureRequest{
		Email:            user.Email,
		Name:             user.FullName,
		BookingReference: booking.BookingReference,
		Reason:           reason,
		From:             trip.Route.Origin,
		To:               trip.Route.Destination,
		DepartureTime:    trip.DepartureTime.Format(constants.DateTimeFormatDisplay),
		BookingLink:      constants.DefaultFrontendURL,
	}

	if err := s.notificationClient.SendBookingFailure(ctx, req); err != nil {
		fmt.Printf("Failed to send booking failure email: %v\n", err)
	}
}

func (s *bookingServiceImpl) getSeatNumbersResults(seats []model.BookingSeat) string {
	var numbers string
	for i, seat := range seats {
		if i > 0 {
			numbers += ", "
		}
		numbers += seat.SeatNumber
	}
	return numbers
}
