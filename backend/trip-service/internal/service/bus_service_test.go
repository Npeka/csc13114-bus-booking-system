package service

import (
	"context"
	"errors"
	"testing"

	"bus-booking/trip-service/internal/constants"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBusService_CreateBus_Success(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	req := &model.CreateBusRequest{
		PlateNumber: "ABC-123",
		Model:       "Mercedes Sprinter",
		Floors: []model.FloorConfig{
			{Floor: 1, SeatCapacity: 20},
		},
		Amenities: []string{"WiFi", "AC"},
	}

	mockBusRepo.On("GetBusByPlateNumber", ctx, req.PlateNumber).Return(nil, errors.New("not found"))
	mockBusRepo.On("CreateBus", ctx, mock.AnythingOfType("*model.Bus")).Return(nil)

	// Act
	result, err := service.CreateBus(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.PlateNumber, result.PlateNumber)
	assert.Equal(t, req.Model, result.Model)
	assert.Equal(t, 20, result.SeatCapacity)
	assert.Equal(t, 20, len(result.Seats))
	assert.True(t, result.IsActive)
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_CreateBus_MultipleFloors(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	req := &model.CreateBusRequest{
		PlateNumber: "ABC-123",
		Model:       "Double Decker",
		Floors: []model.FloorConfig{
			{Floor: 1, SeatCapacity: 20},
			{Floor: 2, SeatCapacity: 24},
		},
		Amenities: []string{"WiFi"},
	}

	mockBusRepo.On("GetBusByPlateNumber", ctx, req.PlateNumber).Return(nil, errors.New("not found"))
	mockBusRepo.On("CreateBus", ctx, mock.AnythingOfType("*model.Bus")).Return(nil)

	// Act
	result, err := service.CreateBus(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 44, result.SeatCapacity)
	assert.Equal(t, 44, len(result.Seats))

	// Check that multi-floor buses have floor prefix in seat numbers
	hasFloorPrefix := false
	for _, seat := range result.Seats {
		if len(seat.SeatNumber) > 2 && seat.SeatNumber[0] == 'F' {
			hasFloorPrefix = true
			break
		}
	}
	assert.True(t, hasFloorPrefix, "Multi-floor bus should have floor prefix in seat numbers")
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_CreateBus_PlateNumberExists(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	req := &model.CreateBusRequest{
		PlateNumber: "ABC-123",
		Model:       "Mercedes",
		Floors:      []model.FloorConfig{{Floor: 1, SeatCapacity: 20}},
	}

	existingBus := &model.Bus{PlateNumber: "ABC-123"}
	existingBus.ID = uuid.New()
	mockBusRepo.On("GetBusByPlateNumber", ctx, req.PlateNumber).Return(existingBus, nil)

	// Act
	result, err := service.CreateBus(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "plate number already exists")
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_CreateBus_ExceedCapacity(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	req := &model.CreateBusRequest{
		PlateNumber: "ABC-123",
		Model:       "Large Bus",
		Floors: []model.FloorConfig{
			{Floor: 1, SeatCapacity: 60},
			{Floor: 2, SeatCapacity: 50},
		},
	}

	mockBusRepo.On("GetBusByPlateNumber", ctx, req.PlateNumber).Return(nil, errors.New("not found"))

	// Act
	result, err := service.CreateBus(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot exceed 100")
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_GetBusByID_Success(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	busID := uuid.New()
	expectedBus := &model.Bus{
		PlateNumber:  "ABC-123",
		Model:        "Mercedes",
		SeatCapacity: 45,
		IsActive:     true,
		Seats: []model.Seat{
			{SeatNumber: "A1", SeatType: constants.SeatTypeVIP},
		},
	}
	expectedBus.ID = busID

	mockBusRepo.On("GetBusWithSeatsByID", ctx, busID).Return(expectedBus, nil)

	// Act
	result, err := service.GetBusByID(ctx, busID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, busID, result.ID)
	assert.Equal(t, 1, len(result.Seats))
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_ListBuses_Success(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	bus1 := model.Bus{PlateNumber: "ABC-123"}
	bus1.ID = uuid.New()
	bus2 := model.Bus{PlateNumber: "XYZ-789"}
	bus2.ID = uuid.New()
	buses := []model.Bus{bus1, bus2}

	req := model.ListBusesRequest{PaginationRequest: model.PaginationRequest{Page: 1, PageSize: 20}}
	mockBusRepo.On("ListBuses", ctx, 1, 20).Return(buses, int64(2), nil)

	// Act
	result, total, err := service.ListBuses(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_UpdateBus_Success(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	busID := uuid.New()
	existingBus := &model.Bus{
		PlateNumber:  "ABC-123",
		Model:        "Mercedes",
		SeatCapacity: 45,
	}
	existingBus.ID = busID

	newModel := "New Model"
	newPlate := "NEW-456"
	req := &model.UpdateBusRequest{
		PlateNumber: &newPlate,
		Model:       &newModel,
	}

	mockBusRepo.On("GetBusByID", ctx, busID).Return(existingBus, nil)
	mockBusRepo.On("GetBusByPlateNumber", ctx, newPlate).Return(nil, errors.New("not found"))
	mockBusRepo.On("UpdateBus", ctx, mock.AnythingOfType("*model.Bus")).Return(nil)

	// Act
	result, err := service.UpdateBus(ctx, busID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newPlate, result.PlateNumber)
	assert.Equal(t, newModel, result.Model)
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_UpdateBus_PlateNumberConflict(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	busID := uuid.New()
	existingBus := &model.Bus{PlateNumber: "ABC-123"}
	existingBus.ID = busID
	conflictBus := &model.Bus{PlateNumber: "NEW-456"}
	conflictBus.ID = uuid.New()

	newPlate := "NEW-456"
	req := &model.UpdateBusRequest{PlateNumber: &newPlate}

	mockBusRepo.On("GetBusByID", ctx, busID).Return(existingBus, nil)
	mockBusRepo.On("GetBusByPlateNumber", ctx, newPlate).Return(conflictBus, nil)

	// Act
	result, err := service.UpdateBus(ctx, busID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "plate number already exists")
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_DeleteBus_Success(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	busID := uuid.New()
	mockBusRepo.On("DeleteBus", ctx, busID).Return(nil)

	// Act
	err := service.DeleteBus(ctx, busID)

	// Assert
	assert.NoError(t, err)
	mockBusRepo.AssertExpectations(t)
}

func TestBusService_GetBusSeats_Success(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo)
	ctx := context.Background()

	busID := uuid.New()
	seat1 := model.Seat{SeatNumber: "A1", BusID: busID}
	seat1.ID = uuid.New()
	seat2 := model.Seat{SeatNumber: "A2", BusID: busID}
	seat2.ID = uuid.New()
	seats := []model.Seat{seat1, seat2}

	mockSeatRepo.On("ListByBusID", ctx, busID).Return(seats, nil)

	// Act
	result, err := service.GetBusSeats(ctx, busID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockSeatRepo.AssertExpectations(t)
}

func TestBusService_GenerateSeats_VIPSeats(t *testing.T) {
	// Arrange
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)
	service := NewBusService(mockBusRepo, mockSeatRepo).(*BusServiceImpl)

	floors := []model.FloorConfig{{Floor: 1, SeatCapacity: 20}}

	// Act
	seats := service.generateSeatsForBus(floors)

	// Assert
	assert.Equal(t, 20, len(seats))

	// First row should be VIP
	firstRowSeats := 0
	for _, seat := range seats {
		if seat.Row == 1 {
			assert.Equal(t, constants.SeatTypeVIP, seat.SeatType)
			firstRowSeats++
		}
	}
	assert.Greater(t, firstRowSeats, 0)
}
