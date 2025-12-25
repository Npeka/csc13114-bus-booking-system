package service

import (
	"context"
	"testing"
	"time"

	"bus-booking/booking-service/internal/client/mocks"
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/model/payment"
	"bus-booking/booking-service/internal/model/trip"
	repo_mocks "bus-booking/booking-service/internal/repository/mocks"
	service_mocks "bus-booking/booking-service/internal/service/mocks"
	queue_mocks "bus-booking/shared/queue/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBookingService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	assert.NotNil(t, service)
	assert.IsType(t, &bookingServiceImpl{}, service)
}

func TestGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetByID(ctx, bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()
	transactionID := uuid.New()

	booking := &model.Booking{
		BaseModel:     model.BaseModel{ID: bookingID},
		TransactionID: transactionID,
		Status:        model.BookingStatusConfirmed,
	}

	transaction := &payment.TransactionResponse{
		ID:     transactionID,
		Status: payment.TransactionStatusPaid,
		Amount: 500000,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockPaymentClient.EXPECT().
		GetTransactionByID(ctx, transactionID).
		Return(transaction, nil).
		Times(1)

	result, err := service.GetByID(ctx, bookingID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, bookingID, result.ID)
	assert.NotNil(t, result.Transaction)
	assert.Equal(t, payment.TransactionStatusPaid, result.Transaction.Status)
	assert.Equal(t, 500000, result.Transaction.Amount)
}

func TestGetByID_TransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()
	transactionID := uuid.New()

	booking := &model.Booking{
		BaseModel:     model.BaseModel{ID: bookingID},
		TransactionID: transactionID,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	// Transaction fetch fails
	mockPaymentClient.EXPECT().
		GetTransactionByID(ctx, transactionID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetByID(ctx, bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "giao dịch")
}

func TestGetByReference_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	reference := "BK123456"
	email := "test@example.com"

	booking := &model.Booking{
		BaseModel:        model.BaseModel{ID: uuid.New()},
		BookingReference: reference,
	}

	mockBookingRepo.EXPECT().
		GetBookingByReference(ctx, reference).
		Return(booking, nil).
		Times(1)

	result, err := service.GetByReference(ctx, reference, email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, reference, result.BookingReference)
}

func TestGetByReference_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	reference := "BK999999"
	email := "test@example.com"

	mockBookingRepo.EXPECT().
		GetBookingByReference(ctx, reference).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetByReference(ctx, reference, email)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTripBookings_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	tripID := uuid.New()
	req := model.PaginationRequest{Page: 1, PageSize: 10}

	bookings := []*model.Booking{
		{BaseModel: model.BaseModel{ID: uuid.New()}, TripID: tripID},
	}

	mockBookingRepo.EXPECT().
		GetTripBookings(ctx, tripID, req.PageSize, 0). // pageSize=10, offset=0
		Return(bookings, int64(1), nil).
		Times(1)

	result, total, err := service.GetTripBookings(ctx, req, tripID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
}

func TestCalculateTotalPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	).(*bookingServiceImpl)

	seats := []trip.Seat{
		{PriceMultiplier: 1.0},
		{PriceMultiplier: 1.5}, // VIP seat
	}

	total := service.calculateTotalPrice(100000, seats)

	assert.Equal(t, 250000, total) // 100k + 150k
}

func TestGenerateBookingReference(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	).(*bookingServiceImpl)

	ref := service.generateBookingReference()

	assert.NotEmpty(t, ref)
	assert.True(t, len(ref) >= 10) // At least BK + date + some chars
	assert.True(t, ref[:2] == "BK")
}

func TestGetSeatNumbersResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	).(*bookingServiceImpl)

	seats := []model.BookingSeat{
		{SeatNumber: "A1"},
		{SeatNumber: "A2"},
		{SeatNumber: "B1"},
	}

	result := service.getSeatNumbersResults(seats)

	assert.Contains(t, result, "A1")
	assert.Contains(t, result, "A2")
	assert.Contains(t, result, "B1")
}

func TestToBookingResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	).(*bookingServiceImpl)

	bookingID := uuid.New()
	userID := uuid.New()
	tripID := uuid.New()

	booking := &model.Booking{
		BaseModel:         model.BaseModel{ID: bookingID},
		BookingReference:  "BK123456",
		TripID:            tripID,
		UserID:            userID,
		TotalAmount:       100000,
		Status:            model.BookingStatusConfirmed,
		TransactionStatus: payment.TransactionStatusPaid,
		BookingSeats: []model.BookingSeat{
			{SeatNumber: "A1", Price: 50000},
			{SeatNumber: "A2", Price: 50000},
		},
	}

	result := service.toBookingResponse(booking)

	assert.NotNil(t, result)
	assert.Equal(t, bookingID, result.ID)
	assert.Equal(t, "BK123456", result.BookingReference)
	assert.Equal(t, tripID, result.TripID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, 100000, result.TotalAmount)
	assert.Equal(t, model.BookingStatusConfirmed, result.Status)
	assert.Len(t, result.Seats, 2)
}

func TestCheckSeatAvailability_AllAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	).(*bookingServiceImpl)

	ctx := context.Background()
	tripID := uuid.New()
	seatIDs := []uuid.UUID{uuid.New(), uuid.New()}

	// No booked seats
	mockBookingRepo.EXPECT().
		GetBookedSeatIDs(ctx, tripID).
		Return([]uuid.UUID{}, nil).
		Times(1)

	available, err := service.checkSeatAvailability(ctx, tripID, seatIDs)

	assert.NoError(t, err)
	assert.True(t, available)
}

func TestCheckSeatAvailability_SeatTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	).(*bookingServiceImpl)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	seatIDs := []uuid.UUID{seatID}

	// Seat already booked
	mockBookingRepo.EXPECT().
		GetBookedSeatIDs(ctx, tripID).
		Return([]uuid.UUID{seatID}, nil).
		Times(1)

	available, err := service.checkSeatAvailability(ctx, tripID, seatIDs)

	assert.NoError(t, err)
	assert.False(t, available)
}

func TestCheckSeatAvailability_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	).(*bookingServiceImpl)

	ctx := context.Background()
	tripID := uuid.New()
	seatIDs := []uuid.UUID{uuid.New()}

	mockBookingRepo.EXPECT().
		GetBookedSeatIDs(ctx, tripID).
		Return(nil, assert.AnError).
		Times(1)

	available, err := service.checkSeatAvailability(ctx, tripID, seatIDs)

	assert.Error(t, err)
	assert.False(t, available)
}

func TestGetSeatStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	tripID := uuid.New()
	seat1 := uuid.New()
	seat2 := uuid.New()
	seat3 := uuid.New()
	seatIDs := []uuid.UUID{seat1, seat2, seat3}

	// Mock booked seats
	mockBookingRepo.EXPECT().
		GetBookedSeatIDs(ctx, tripID).
		Return([]uuid.UUID{seat1}, nil).
		Times(1)

	// Mock locked seats
	mockSeatLockService.EXPECT().
		GetLockedSeats(ctx, tripID).
		Return([]uuid.UUID{seat2}, nil).
		Times(1)

	result, err := service.GetSeatStatus(ctx, tripID, seatIDs)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.True(t, result[0].IsBooked)
	assert.False(t, result[0].IsLocked)
	assert.False(t, result[1].IsBooked)
	assert.True(t, result[1].IsLocked)
	assert.False(t, result[2].IsBooked)
	assert.False(t, result[2].IsLocked)
}

func TestGetSeatStatus_EmptySeats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	tripID := uuid.New()

	result, err := service.GetSeatStatus(ctx, tripID, []uuid.UUID{})

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetSeatStatus_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	tripID := uuid.New()
	seatIDs := []uuid.UUID{uuid.New()}

	mockBookingRepo.EXPECT().
		GetBookedSeatIDs(ctx, tripID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetSeatStatus(ctx, tripID, seatIDs)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetUserBookings_WithTrips(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	userID := uuid.New()
	tripID1 := uuid.New()
	tripID2 := uuid.New()

	req := model.GetUserBookingsRequest{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	bookings := []*model.Booking{
		{
			BaseModel: model.BaseModel{ID: uuid.New()},
			TripID:    tripID1,
			UserID:    userID,
		},
		{
			BaseModel: model.BaseModel{ID: uuid.New()},
			TripID:    tripID2,
			UserID:    userID,
		},
	}

	trips := []trip.Trip{
		{
			ID:        tripID1,
			BasePrice: 100000,
			Route: &trip.Route{
				Origin:      "Ha Noi",
				Destination: "Da Nang",
			},
			Bus: &trip.Bus{
				Model: "Hyundai Universe",
			},
		},
		{
			ID:        tripID2,
			BasePrice: 150000,
			Route: &trip.Route{
				Origin:      "Ho Chi Minh",
				Destination: "Vung Tau",
			},
			Bus: &trip.Bus{
				Model: "Mercedes Benz",
			},
		},
	}

	// Mock GetBookingsByUserID
	mockBookingRepo.EXPECT().
		GetBookingsByUserID(ctx, userID, req.Status, 10, 0).
		Return(bookings, int64(2), nil).
		Times(1)

	// Mock GetTripsByIDs
	mockTripClient.EXPECT().
		GetTripsByIDs(ctx, gomock.Any(), gomock.Any()).
		Return(trips, nil).
		Times(1)

	result, total, err := service.GetUserBookings(ctx, req, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	assert.NotNil(t, result[0].Trip)
	assert.Equal(t, "Ha Noi", result[0].Trip.Origin)
	assert.NotNil(t, result[1].Trip)
	assert.Equal(t, "Ho Chi Minh", result[1].Trip.Origin)
}

func TestCreateGuestBooking_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	req := &model.CreateGuestBookingRequest{
		Email: "",
		Phone: "", // Both empty
	}

	result, err := service.CreateGuestBooking(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email hoặc số điện thoại")
}

func TestCreateGuestBooking_UserClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	req := &model.CreateGuestBookingRequest{
		Email:    "guest@example.com",
		FullName: "Guest User",
	}

	mockUserClient.EXPECT().
		CreateGuest(ctx, gomock.Any()).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.CreateGuestBooking(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "không thể tạo tài khoản khách")
}

func TestCancelBooking_AlreadyCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	booking := &model.Booking{
		BaseModel: model.BaseModel{ID: bookingID},
		Status:    model.BookingStatusCancelled,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	err := service.CancelBooking(ctx, bookingID, "test")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already cancelled")
}

func TestCancelBooking_CannotCancelConfirmed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	booking := &model.Booking{
		BaseModel: model.BaseModel{ID: bookingID},
		Status:    model.BookingStatusConfirmed,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	err := service.CancelBooking(ctx, bookingID, "test")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot cancel confirmed")
}

func TestCancelBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()
	txID := uuid.New()

	booking := &model.Booking{
		BaseModel:     model.BaseModel{ID: bookingID},
		Status:        model.BookingStatusPending,
		TransactionID: txID,
	}

	cancelledTx := &payment.TransactionResponse{
		ID:     txID,
		Status: payment.TransactionStatusCancelled,
	}

	updatedBooking := &model.Booking{
		BaseModel:         model.BaseModel{ID: bookingID},
		Status:            model.BookingStatusCancelled,
		TransactionID:     txID,
		TransactionStatus: payment.TransactionStatusPending,
	}

	// Mock sequence
	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		CancelBooking(ctx, bookingID, "test reason").
		Return(nil).
		Times(1)

	mockPaymentClient.EXPECT().
		CancelTransaction(ctx, txID).
		Return(cancelledTx, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(updatedBooking, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		UpdateBooking(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	err := service.CancelBooking(ctx, bookingID, "test reason")

	assert.NoError(t, err)
}

func TestCancelBooking_PaymentCancelFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()
	txID := uuid.New()

	booking := &model.Booking{
		BaseModel:     model.BaseModel{ID: bookingID},
		Status:        model.BookingStatusPending,
		TransactionID: txID,
	}

	// Mock sequence
	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		CancelBooking(ctx, bookingID, "test").
		Return(nil).
		Times(1)

	// Payment cancellation fails
	mockPaymentClient.EXPECT().
		CancelTransaction(ctx, txID).
		Return(nil, assert.AnError).
		Times(1)

	// Should still return nil (booking is cancelled despite payment error)
	err := service.CancelBooking(ctx, bookingID, "test")

	assert.NoError(t, err) // Function logs error but doesn't fail
}

func TestCancelBooking_GetBookingAfterCancelFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()
	txID := uuid.New()

	booking := &model.Booking{
		BaseModel:     model.BaseModel{ID: bookingID},
		Status:        model.BookingStatusPending,
		TransactionID: txID,
	}

	cancelledTx := &payment.TransactionResponse{
		ID:     txID,
		Status: payment.TransactionStatusCancelled,
	}

	// Mock sequence
	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		CancelBooking(ctx, bookingID, "test").
		Return(nil).
		Times(1)

	mockPaymentClient.EXPECT().
		CancelTransaction(ctx, txID).
		Return(cancelledTx, nil).
		Times(1)

	// Second GetBookingByID fails
	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(nil, assert.AnError).
		Times(1)

	// Should log error but still return nil (booking is already cancelled)
	err := service.CancelBooking(ctx, bookingID, "test")

	assert.NoError(t, err) // Logs error but succeeds
}

func TestGetUserBookings_TripEnrichmentFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	userID := uuid.New()
	tripID := uuid.New()

	req := model.GetUserBookingsRequest{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	bookings := []*model.Booking{
		{BaseModel: model.BaseModel{ID: uuid.New()}, TripID: tripID, UserID: userID},
	}

	mockBookingRepo.EXPECT().
		GetBookingsByUserID(ctx, userID, req.Status, 10, 0).
		Return(bookings, int64(1), nil).
		Times(1)

	// Trip enrichment fails
	mockTripClient.EXPECT().
		GetTripsByIDs(ctx, gomock.Any(), gomock.Any()).
		Return(nil, assert.AnError).
		Times(1)

	// Should still return bookings (just without trip data)
	result, total, err := service.GetUserBookings(ctx, req, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	assert.Nil(t, result[0].Trip) // No trip data
}

func TestGetTripBookings_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	tripID := uuid.New()
	req := model.PaginationRequest{Page: 1, PageSize: 10}

	mockBookingRepo.EXPECT().
		GetTripBookings(ctx, tripID, 10, 0).
		Return(nil, int64(0), assert.AnError).
		Times(1)

	result, total, err := service.GetTripBookings(ctx, req, tripID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
}

func TestCreateBooking_SeatUnavailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	userID := uuid.New()
	tripID := uuid.New()
	seatID := uuid.New()

	req := &model.CreateBookingRequest{
		TripID:  tripID,
		SeatIDs: []uuid.UUID{seatID},
	}

	// Seat is already booked
	mockBookingRepo.EXPECT().
		GetBookedSeatIDs(ctx, tripID).
		Return([]uuid.UUID{seatID}, nil).
		Times(1)

	result, err := service.CreateBooking(ctx, req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already booked")
}

func TestCreateBooking_PaymentCreationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	userID := uuid.New()
	tripID := uuid.New()
	seatID := uuid.New()

	req := &model.CreateBookingRequest{
		TripID:  tripID,
		SeatIDs: []uuid.UUID{seatID},
	}

	tripData := &trip.Trip{
		ID:        tripID,
		BasePrice: 100000,
	}

	seatData := trip.Seat{
		ID:              seatID,
		SeatNumber:      "A1",
		PriceMultiplier: 1.0,
	}

	// Mock sequence
	mockBookingRepo.EXPECT().
		GetBookedSeatIDs(ctx, tripID).
		Return([]uuid.UUID{}, nil).
		Times(1)

	mockTripClient.EXPECT().
		GetTripByID(gomock.Any(), gomock.Any(), tripID).
		Return(tripData, nil).
		Times(1)

	mockTripClient.EXPECT().
		ListSeatsByIDs(gomock.Any(), gomock.Any()).
		Return([]trip.Seat{seatData}, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		CreateBooking(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	// Payment fails
	mockPaymentClient.EXPECT().
		CreateTransaction(ctx, gomock.Any()).
		Return(nil, assert.AnError).
		Times(1)

	// Should update to FAILED
	mockBookingRepo.EXPECT().
		UpdateBooking(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	result, err := service.CreateBooking(ctx, req, userID)

	assert.NoError(t, err) // Returns booking with failed status
	assert.NotNil(t, result)
	assert.Equal(t, payment.TransactionStatusFailed, result.Transaction.Status)
}

func TestUpdateBookingStatus_Cancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	booking := &model.Booking{
		BaseModel: model.BaseModel{ID: bookingID},
		Status:    model.BookingStatusPending,
	}

	req := &model.UpdateBookingStatusRequest{
		TransactionStatus: payment.TransactionStatusCancelled,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		UpdateBooking(ctx, gomock.Any()).
		Do(func(_ context.Context, b *model.Booking) {
			assert.Equal(t, model.BookingStatusCancelled, b.Status)
			assert.NotNil(t, b.CancelledAt)
		}).
		Return(nil).
		Times(1)

	err := service.UpdateBookingStatus(ctx, req, bookingID)

	assert.NoError(t, err)
}

func TestUpdateBookingStatus_Expired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	booking := &model.Booking{
		BaseModel: model.BaseModel{ID: bookingID},
		Status:    model.BookingStatusPending,
	}

	req := &model.UpdateBookingStatusRequest{
		TransactionStatus: payment.TransactionStatusExpired,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockBookingRepo.EXPECT().
		UpdateBooking(ctx, gomock.Any()).
		Do(func(_ context.Context, b *model.Booking) {
			assert.Equal(t, model.BookingStatusExpired, b.Status)
		}).
		Return(nil).
		Times(1)

	err := service.UpdateBookingStatus(ctx, req, bookingID)

	assert.NoError(t, err)
}

func TestRetryPayment_BookingNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.RetryPayment(ctx, bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestRetryPayment_NotRetryableState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	// Booking is confirmed - not retryable
	booking := &model.Booking{
		BaseModel: model.BaseModel{ID: bookingID},
		Status:    model.BookingStatusConfirmed,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	result, err := service.RetryPayment(ctx, bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not in a retryable state")
}

func TestRetryPayment_ExpiredBeyondGracePeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	// Expired 2 hours ago
	expiredTime := time.Now().UTC().Add(-2 * time.Hour)
	booking := &model.Booking{
		BaseModel: model.BaseModel{ID: bookingID},
		Status:    model.BookingStatusExpired,
		ExpiresAt: &expiredTime,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	result, err := service.RetryPayment(ctx, bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "expired beyond retry period")
}

// ========== ExpireBooking - Validation paths ==========

func TestExpireBooking_BookingNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(nil, assert.AnError).
		Times(1)

	err := service.ExpireBooking(ctx, bookingID)

	assert.Error(t, err)
}

func TestExpireBooking_AlreadyProcessed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	mockUserClient := mocks.NewMockUserClient(ctrl)
	mockNotificationClient := mocks.NewMockNotificationClient(ctrl)
	mockDelayedQueue := queue_mocks.NewMockDelayedQueueManager(ctrl)
	mockSeatLockService := service_mocks.NewMockSeatLockService(ctrl)

	service := NewBookingService(
		mockBookingRepo,
		mockPaymentClient,
		mockTripClient,
		mockUserClient,
		mockNotificationClient,
		mockDelayedQueue,
		mockSeatLockService,
	)

	ctx := context.Background()
	bookingID := uuid.New()

	// Already confirmed - skip expiration
	booking := &model.Booking{
		BaseModel: model.BaseModel{ID: bookingID},
		Status:    model.BookingStatusConfirmed,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	err := service.ExpireBooking(ctx, bookingID)

	assert.NoError(t, err) // Returns nil - already processed
}
