package service

import (
	"context"
	"testing"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSeatService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	service := NewSeatService(mockSeatRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &SeatServiceImpl{}, service)
}

func TestGetListByIDs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	service := NewSeatService(mockSeatRepo)

	ctx := context.Background()
	seatIDs := []uuid.UUID{uuid.New(), uuid.New()}

	expectedSeats := []model.Seat{
		{
			BaseModel:       model.BaseModel{ID: seatIDs[0]},
			SeatNumber:      "A1",
			PriceMultiplier: 1.0,
		},
		{
			BaseModel:       model.BaseModel{ID: seatIDs[1]},
			SeatNumber:      "A2",
			PriceMultiplier: 1.0,
		},
	}

	mockSeatRepo.EXPECT().
		GetListByIDs(ctx, seatIDs).
		Return(expectedSeats, nil).
		Times(1)

	result, err := service.GetListByIDs(ctx, seatIDs)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "A1", result[0].SeatNumber)
	assert.Equal(t, "A2", result[1].SeatNumber)
}

func TestGetListByIDs_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	service := NewSeatService(mockSeatRepo)

	ctx := context.Background()
	seatIDs := []uuid.UUID{uuid.New()}

	mockSeatRepo.EXPECT().
		GetListByIDs(ctx, seatIDs).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.GetListByIDs(ctx, seatIDs)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list seats")
}

func TestUpdateSeat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	service := NewSeatService(mockSeatRepo)

	ctx := context.Background()
	seatID := uuid.New()

	newSeatNumber := "B1"
	newPriceMultiplier := 1.2

	existingSeat := &model.Seat{
		BaseModel:       model.BaseModel{ID: seatID},
		SeatNumber:      "A1",
		PriceMultiplier: 1.0,
	}

	req := &model.UpdateSeatRequest{
		SeatNumber:      &newSeatNumber,
		PriceMultiplier: &newPriceMultiplier,
	}

	mockSeatRepo.EXPECT().
		GetByID(ctx, seatID).
		Return(existingSeat, nil).
		Times(1)

	mockSeatRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Do(func(_ context.Context, seat *model.Seat) {
			assert.Equal(t, "B1", seat.SeatNumber)
			assert.Equal(t, 1.2, seat.PriceMultiplier)
		}).
		Return(nil).
		Times(1)

	result, err := service.Update(ctx, req, seatID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "B1", result.SeatNumber)
	assert.Equal(t, 1.2, result.PriceMultiplier)
}

func TestUpdate_SeatNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	service := NewSeatService(mockSeatRepo)

	ctx := context.Background()
	seatID := uuid.New()

	req := &model.UpdateSeatRequest{}

	mockSeatRepo.EXPECT().
		GetByID(ctx, seatID).
		Return(nil, assert.AnError).
		Times(1)

	result, err := service.Update(ctx, req, seatID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "seat not found")
}

func TestUpdate_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeatRepo := mocks.NewMockSeatRepository(ctrl)
	service := NewSeatService(mockSeatRepo)

	ctx := context.Background()
	seatID := uuid.New()

	existingSeat := &model.Seat{
		BaseModel:  model.BaseModel{ID: seatID},
		SeatNumber: "A1",
	}

	req := &model.UpdateSeatRequest{}

	mockSeatRepo.EXPECT().
		GetByID(ctx, seatID).
		Return(existingSeat, nil).
		Times(1)

	mockSeatRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(assert.AnError).
		Times(1)

	result, err := service.Update(ctx, req, seatID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update seat")
}
