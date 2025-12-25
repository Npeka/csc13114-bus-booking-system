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

func TestNewSeatStatusService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &SeatStatusServiceImpl{}, service)
}

func TestInitSeatsForTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID1 := uuid.New()
	seatID2 := uuid.New()

	seats := []model.SeatInitData{
		{SeatID: seatID1, SeatNumber: "A1"},
		{SeatID: seatID2, SeatNumber: "A2"},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return([]*model.SeatStatus{}, nil).
		Times(1)

	mockRepo.EXPECT().
		BulkUpdateSeatStatus(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, statuses []*model.SeatStatus) error {
			assert.Len(t, statuses, 2)
			assert.Equal(t, tripID, statuses[0].TripID)
			assert.Equal(t, "available", statuses[0].Status)
			return nil
		}).
		Times(1)

	err := service.InitSeatsForTrip(ctx, tripID, seats)

	assert.NoError(t, err)
}

func TestInitSeatsForTrip_AlreadyInitialized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()

	existingSeats := []*model.SeatStatus{
		{TripID: tripID, SeatID: uuid.New(), Status: "available"},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(existingSeats, nil).
		Times(1)

	err := service.InitSeatsForTrip(ctx, tripID, []model.SeatInitData{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already initialized")
}

func TestInitSeatsForTrip_CheckExistingError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(nil, assert.AnError).
		Times(1)

	err := service.InitSeatsForTrip(ctx, tripID, []model.SeatInitData{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check existing seats")
}

func TestSeatStatusGetSeatAvailability_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()

	seats := []*model.SeatStatus{
		{SeatID: seatID, Status: "available"},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seats, nil).
		Times(1)

	result, err := service.GetSeatAvailability(ctx, tripID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "available", result[0].Status)
}

func TestSeatStatusGetSeatAvailability_WithExpiredReservation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	userID := uuid.New()
	pastTime := time.Now().UTC().Add(-1 * time.Hour)

	seats := []*model.SeatStatus{
		{
			SeatID:        seatID,
			Status:        "reserved",
			ReservedUntil: &pastTime,
			UserID:        &userID,
		},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seats, nil).
		Times(1)

	mockRepo.EXPECT().
		ReleaseSeat(ctx, tripID, seatID).
		Return(nil).
		Times(1)

	result, err := service.GetSeatAvailability(ctx, tripID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "available", result[0].Status)
	assert.Nil(t, result[0].ReservedUntil)
}

func TestCheckAndReleaseSeat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	pastTime := time.Now().UTC().Add(-1 * time.Hour)

	seats := []*model.SeatStatus{
		{
			SeatID:        seatID,
			Status:        "reserved",
			ReservedUntil: &pastTime,
		},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seats, nil).
		Times(1)

	mockRepo.EXPECT().
		ReleaseSeat(ctx, tripID, seatID).
		Return(nil).
		Times(1)

	err := service.CheckAndReleaseSeat(ctx, tripID, seatID)

	assert.NoError(t, err)
}

func TestCheckAndReleaseSeat_NotExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	futureTime := time.Now().UTC().Add(1 * time.Hour)

	seats := []*model.SeatStatus{
		{
			SeatID:        seatID,
			Status:        "reserved",
			ReservedUntil: &futureTime,
		},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seats, nil).
		Times(1)

	err := service.CheckAndReleaseSeat(ctx, tripID, seatID)

	assert.NoError(t, err)
}

func TestCheckAndReleaseSeat_SeatNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()
	differentSeatID := uuid.New()

	seats := []*model.SeatStatus{
		{SeatID: differentSeatID, Status: "available"},
	}

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(seats, nil).
		Times(1)

	err := service.CheckAndReleaseSeat(ctx, tripID, seatID)

	assert.NoError(t, err)
}

func TestCheckAndReleaseSeat_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSeatStatusRepository(ctrl)
	service := NewSeatStatusService(mockRepo)

	ctx := context.Background()
	tripID := uuid.New()
	seatID := uuid.New()

	mockRepo.EXPECT().
		GetSeatStatusByTripID(ctx, tripID).
		Return(nil, assert.AnError).
		Times(1)

	err := service.CheckAndReleaseSeat(ctx, tripID, seatID)

	assert.Error(t, err)
}
