package service

import (
	"context"
	"testing"
	"time"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSeatService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &SeatServiceImpl{}, service)
}

func TestGetSeatAvailability_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID1 := uuid.New()
	seatID2 := uuid.New()
	seatID3 := uuid.New()

	seatStatuses := []*model.SeatStatus{
		{SeatID: seatID1, Status: "available"},
		{SeatID: seatID2, Status: "reserved"},
		{SeatID: seatID3, Status: "booked"},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seatStatuses, nil).
		Times(1)

	result, err := service.GetSeatAvailability(ctx, tripID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tripID, result.TripID)
	assert.Len(t, result.AvailableSeats, 1)
	assert.Len(t, result.ReservedSeats, 1)
	assert.Len(t, result.BookedSeats, 1)
	assert.Contains(t, result.AvailableSeats, seatID1)
	assert.Contains(t, result.ReservedSeats, seatID2)
	assert.Contains(t, result.BookedSeats, seatID3)
}

func TestGetSeatAvailability_WithExpiredReservation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	pastTime := time.Now().UTC().Add(-1 * time.Hour)

	seatStatuses := []*model.SeatStatus{
		{
			SeatID:        seatID,
			Status:        "reserved",
			ReservedUntil: &pastTime,
		},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seatStatuses, nil).
		Times(1)

	mockRepo.EXPECT().
		ReleaseSeat(ctx, tripID, seatID).
		Return(nil).
		Times(1)

	result, err := service.GetSeatAvailability(ctx, tripID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.AvailableSeats, 1)
	assert.Empty(t, result.ReservedSeats)
}

func TestGetSeatAvailability_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetSeatAvailability(ctx, tripID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReserveSeat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	userID := uuid.New()

	req := &model.ReserveSeatRequest{
		TripID:             tripID,
		SeatID:             seatID,
		UserID:             userID,
		ReservationMinutes: 15,
	}

	seatStatuses := []*model.SeatStatus{
		{SeatID: seatID, Status: "available"},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seatStatuses, nil).
		Times(1)

	mockRepo.EXPECT().
		ReserveSeat(ctx, tripID, seatID, userID, 15*time.Minute).
		Return(nil).
		Times(1)

	err := service.ReserveSeat(ctx, req)

	assert.NoError(t, err)
}

func TestReserveSeat_SeatNotAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	userID := uuid.New()

	req := &model.ReserveSeatRequest{
		TripID: tripID,
		SeatID: seatID,
		UserID: userID,
	}

	seatStatuses := []*model.SeatStatus{
		{SeatID: seatID, Status: "booked"},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seatStatuses, nil).
		Times(1)

	err := service.ReserveSeat(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not available")
}

func TestReserveSeat_DefaultReservationTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	userID := uuid.New()

	req := &model.ReserveSeatRequest{
		TripID:             tripID,
		SeatID:             seatID,
		UserID:             userID,
		ReservationMinutes: 0, // Should use default 15 minutes
	}

	seatStatuses := []*model.SeatStatus{
		{SeatID: seatID, Status: "available"},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seatStatuses, nil).
		Times(1)

	mockRepo.EXPECT().
		ReserveSeat(ctx, tripID, seatID, userID, 15*time.Minute).
		Return(nil).
		Times(1)

	err := service.ReserveSeat(ctx, req)

	assert.NoError(t, err)
}

func TestReleaseSeat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()

	mockRepo.EXPECT().
		ReleaseSeat(ctx, tripID, seatID).
		Return(nil).
		Times(1)

	err := service.ReleaseSeat(ctx, tripID, seatID)

	assert.NoError(t, err)
}

func TestReleaseSeat_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()

	mockRepo.EXPECT().
		ReleaseSeat(ctx, tripID, seatID).
		Return(assert.AnError).
		Times(1)

	err := service.ReleaseSeat(ctx, tripID, seatID)

	assert.Error(t, err)
}

func TestCheckReservationExpiry_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatService(mockRepo)

	ctx := context.Background()

	err := service.CheckReservationExpiry(ctx)

	assert.NoError(t, err)
}
