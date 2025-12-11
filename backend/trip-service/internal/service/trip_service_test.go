package service

import (
	"context"
	"testing"
	"time"

	"bus-booking/trip-service/internal/constants"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTripService_CreateTrip_Success(t *testing.T) {
	// Arrange
	mockTripRepo := new(mocks.MockTripRepository)
	mockRouteRepo := new(mocks.MockRouteRepository)
	mockRouteStopRepo := new(mocks.MockRouteStopRepository)
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, nil)
	ctx := context.Background()

	routeID := uuid.New()
	busID := uuid.New()
	departureTime := time.Now().Add(24 * time.Hour)
	arrivalTime := departureTime.Add(2 * time.Hour)

	route := &model.Route{Origin: "Hà Nội", Destination: "Hải Phòng"}
	route.ID = routeID

	bus := &model.Bus{PlateNumber: "ABC-123", SeatCapacity: 45, IsActive: true}
	bus.ID = busID

	req := &model.CreateTripRequest{
		RouteID:       routeID,
		BusID:         busID,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		BasePrice:     200000,
	}

	createdTrip := &model.Trip{
		RouteID:       routeID,
		BusID:         busID,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		BasePrice:     200000,
		Status:        "scheduled",
	}
	createdTrip.ID = uuid.New()

	mockRouteRepo.On("GetRouteByID", ctx, routeID).Return(route, nil)
	mockBusRepo.On("GetBusByID", ctx, busID).Return(bus, nil)
	mockTripRepo.On("GetTripsByBusAndDateRange", ctx, busID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return([]model.Trip{}, nil)
	mockTripRepo.On("CreateTrip", ctx, mock.AnythingOfType("*model.Trip")).Return(nil)
	mockTripRepo.On("GetTripByID", ctx, mock.AnythingOfType("*model.GetTripByIDRequuest"), mock.AnythingOfType("uuid.UUID")).Return(createdTrip, nil)

	// Act
	result, err := service.CreateTrip(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockTripRepo.AssertExpectations(t)
	mockRouteRepo.AssertExpectations(t)
	mockBusRepo.AssertExpectations(t)
}

func TestTripService_GetTripByID_Success(t *testing.T) {
	// Arrange
	mockTripRepo := new(mocks.MockTripRepository)
	mockRouteRepo := new(mocks.MockRouteRepository)
	mockRouteStopRepo := new(mocks.MockRouteStopRepository)
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, nil)
	ctx := context.Background()

	tripID := uuid.New()
	expectedTrip := &model.Trip{
		RouteID:   uuid.New(),
		BusID:     uuid.New(),
		BasePrice: 200000,
		Status:    "scheduled",
	}
	expectedTrip.ID = tripID

	mockTripRepo.On("GetTripByID", ctx, mock.AnythingOfType("*model.GetTripByIDRequuest"), tripID).Return(expectedTrip, nil)

	// Act
	result, err := service.GetTripByID(ctx, &model.GetTripByIDRequest{}, tripID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tripID, result.ID)
	mockTripRepo.AssertExpectations(t)
}

func TestTripService_ListTrips_Success(t *testing.T) {
	// Arrange
	mockTripRepo := new(mocks.MockTripRepository)
	mockRouteRepo := new(mocks.MockRouteRepository)
	mockRouteStopRepo := new(mocks.MockRouteStopRepository)
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, nil)
	ctx := context.Background()

	trip1 := model.Trip{RouteID: uuid.New(), BusID: uuid.New(), Status: "scheduled"}
	trip1.ID = uuid.New()
	trip2 := model.Trip{RouteID: uuid.New(), BusID: uuid.New(), Status: "completed"}
	trip2.ID = uuid.New()
	trips := []model.Trip{trip1, trip2}

	mockTripRepo.On("ListTrips", ctx, 1, 20).Return(trips, int64(2), nil)

	// Act
	result, total, err := service.ListTrips(ctx, &model.ListTripsRequest{
		PaginationRequest: model.PaginationRequest{
			Page:     1,
			PageSize: 20,
		},
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockTripRepo.AssertExpectations(t)
}

func TestTripService_UpdateTrip_Success(t *testing.T) {
	// Arrange
	mockTripRepo := new(mocks.MockTripRepository)
	mockRouteRepo := new(mocks.MockRouteRepository)
	mockRouteStopRepo := new(mocks.MockRouteStopRepository)
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, nil)
	ctx := context.Background()

	tripID := uuid.New()
	existingTrip := &model.Trip{
		RouteID:   uuid.New(),
		BusID:     uuid.New(),
		BasePrice: 200000,
		Status:    "scheduled",
	}
	existingTrip.ID = tripID

	newPrice := float64(250000)
	newStatus := constants.TripStatusCompleted
	req := &model.UpdateTripRequest{
		BasePrice: &newPrice,
		Status:    &newStatus,
	}

	mockTripRepo.On("GetTripByID", ctx, mock.AnythingOfType("*model.GetTripByIDRequuest"), tripID).Return(existingTrip, nil)
	mockTripRepo.On("UpdateTrip", ctx, existingTrip).Return(nil)

	// Act
	result, err := service.UpdateTrip(ctx, tripID, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, newPrice, result.BasePrice)
	assert.Equal(t, newStatus, result.Status)
	mockTripRepo.AssertExpectations(t)
}

func TestTripService_DeleteTrip_Success(t *testing.T) {
	// Arrange
	mockTripRepo := new(mocks.MockTripRepository)
	mockRouteRepo := new(mocks.MockRouteRepository)
	mockRouteStopRepo := new(mocks.MockRouteStopRepository)
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, nil)
	ctx := context.Background()

	tripID := uuid.New()
	futureTime := time.Now().Add(48 * time.Hour)
	trip := &model.Trip{Status: "scheduled", DepartureTime: futureTime}
	trip.ID = tripID

	mockTripRepo.On("GetTripByID", ctx, mock.AnythingOfType("*model.GetTripByIDRequuest"), tripID).Return(trip, nil)
	mockTripRepo.On("DeleteTrip", ctx, tripID).Return(nil)

	// Act
	err := service.DeleteTrip(ctx, tripID)

	// Assert
	assert.NoError(t, err)
	mockTripRepo.AssertExpectations(t)
}

func TestTripService_GetSeatAvailability_Success(t *testing.T) {
	// Arrange
	mockTripRepo := new(mocks.MockTripRepository)
	mockRouteRepo := new(mocks.MockRouteRepository)
	mockRouteStopRepo := new(mocks.MockRouteStopRepository)
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, nil)
	ctx := context.Background()

	tripID := uuid.New()
	busID := uuid.New()

	trip := &model.Trip{BusID: busID}
	trip.ID = tripID

	seat1 := model.Seat{SeatNumber: "A1", BusID: busID, IsAvailable: true}
	seat1.ID = uuid.New()
	seat2 := model.Seat{SeatNumber: "A2", BusID: busID, IsAvailable: false}
	seat2.ID = uuid.New()
	seats := []model.Seat{seat1, seat2}

	mockTripRepo.On("GetTripByID", ctx, mock.AnythingOfType("*model.GetTripByIDRequuest"), tripID).Return(trip, nil)
	mockSeatRepo.On("ListByBusID", ctx, busID).Return(seats, nil)

	// Act
	result, err := service.GetSeatAvailability(ctx, tripID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalSeats)
	assert.Equal(t, 1, result.AvailableSeats)
	mockTripRepo.AssertExpectations(t)
	mockSeatRepo.AssertExpectations(t)
}

func TestTripService_GetTripsByRouteAndDate_Success(t *testing.T) {
	// Arrange
	mockTripRepo := new(mocks.MockTripRepository)
	mockRouteRepo := new(mocks.MockRouteRepository)
	mockRouteStopRepo := new(mocks.MockRouteStopRepository)
	mockBusRepo := new(mocks.MockBusRepository)
	mockSeatRepo := new(mocks.MockSeatRepository)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, nil)
	ctx := context.Background()

	routeID := uuid.New()
	date := time.Now()

	route := &model.Route{Origin: "Hà Nội", Destination: "Hải Phòng"}
	route.ID = routeID
	trip1 := model.Trip{RouteID: routeID, Status: "scheduled"}
	trip1.ID = uuid.New()
	trips := []model.Trip{trip1}

	mockRouteRepo.On("GetRouteByID", ctx, routeID).Return(route, nil)
	mockTripRepo.On("GetTripsByRouteAndDate", ctx, routeID, mock.AnythingOfType("time.Time")).Return(trips, nil)

	// Act
	result, err := service.GetTripsByRouteAndDate(ctx, routeID, date)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, routeID, result[0].RouteID)
	mockTripRepo.AssertExpectations(t)
}
