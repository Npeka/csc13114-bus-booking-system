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

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewETicketService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)

	service := NewETicketService(mockBookingRepo, mockTripClient)

	assert.NotNil(t, service)
	assert.IsType(t, &eTicketServiceImpl{}, service)
}

func TestGenerateETicket_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	service := NewETicketService(mockBookingRepo, mockTripClient)

	ctx := context.Background()
	bookingID := uuid.New()
	tripID := uuid.New()

	booking := &model.Booking{
		BaseModel:         model.BaseModel{ID: bookingID},
		BookingReference:  "BK123456",
		TripID:            tripID,
		Status:            model.BookingStatusConfirmed,
		TransactionStatus: payment.TransactionStatusPaid,
		TotalAmount:       500000,
		BookingSeats: []model.BookingSeat{
			{SeatNumber: "A1"},
			{SeatNumber: "A2"},
		},
	}

	tripData := &trip.Trip{
		ID:            tripID,
		BasePrice:     250000,
		DepartureTime: time.Now().UTC().Add(-1 * time.Hour),
		ArrivalTime:   time.Now().UTC().Add(1 * time.Hour),
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockTripClient.EXPECT().
		GetTripByID(ctx, trip.GetTripByIDRequest{}, tripID).
		Return(tripData, nil).
		Times(1)

	result, err := service.GenerateETicket(ctx, bookingID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, result.Len(), 0)
}

func TestGenerateETicket_BookingNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	service := NewETicketService(mockBookingRepo, mockTripClient)

	ctx := context.Background()
	bookingID := uuid.New()

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GenerateETicket(ctx, bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestGenerateETicket_BookingNotConfirmed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	service := NewETicketService(mockBookingRepo, mockTripClient)

	ctx := context.Background()
	bookingID := uuid.New()

	booking := &model.Booking{
		BaseModel: model.BaseModel{ID: bookingID},
		Status:    model.BookingStatusPending,
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	result, err := service.GenerateETicket(ctx, bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "only available for confirmed bookings")
}

func TestGenerateETicket_SuccessWithoutTripData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	service := NewETicketService(mockBookingRepo, mockTripClient)

	ctx := context.Background()
	bookingID := uuid.New()
	tripID := uuid.New()

	booking := &model.Booking{
		BaseModel:         model.BaseModel{ID: bookingID},
		BookingReference:  "BK123456",
		TripID:            tripID,
		Status:            model.BookingStatusConfirmed,
		TransactionStatus: payment.TransactionStatusPaid,
		TotalAmount:       500000,
		BookingSeats: []model.BookingSeat{
			{SeatNumber: "A1"},
		},
	}

	mockBookingRepo.EXPECT().
		GetBookingByID(ctx, bookingID).
		Return(booking, nil).
		Times(1)

	mockTripClient.EXPECT().
		GetTripByID(ctx, trip.GetTripByIDRequest{}, tripID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GenerateETicket(ctx, bookingID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, result.Len(), 0)
}

func TestGetPaymentStatusText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := repo_mocks.NewMockBookingRepository(ctrl)
	mockTripClient := mocks.NewMockTripClient(ctrl)
	service := NewETicketService(mockBookingRepo, mockTripClient).(*eTicketServiceImpl)

	tests := []struct {
		status   payment.TransactionStatus
		expected string
	}{
		{payment.TransactionStatusPending, "Cho thanh toan"},
		{payment.TransactionStatusPaid, "Da thanh toan"},
		{payment.TransactionStatusFailed, "That bai"},
		{payment.TransactionStatusCancelled, "CANCELLED"}, // Default case returns string(status)
	}

	for _, tt := range tests {
		result := service.getPaymentStatusText(tt.status)
		assert.Equal(t, tt.expected, result)
	}
}
