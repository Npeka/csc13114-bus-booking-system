package service

import (
	"context"
	"testing"
	"time"

	"bus-booking/trip-service/internal/client/mocks"
	"bus-booking/trip-service/internal/constants"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/model/booking"
	repo_mocks "bus-booking/trip-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTripService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(
		mockTripRepo,
		mockRouteRepo,
		mockRouteStopRepo,
		mockBusRepo,
		mockSeatRepo,
		mockBookingClient,
		nil,
	)

	assert.NotNil(t, service)
	assert.IsType(t, &TripServiceImpl{}, service)
}

func TestSearchTrips_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	origin := "Ha Noi"
	destination := "Da Nang"
	req := &model.TripSearchRequest{Origin: &origin, Destination: &destination}

	expectedTrips := []model.TripDetail{
		{ID: uuid.New()},
	}

	mockTripRepo.EXPECT().
		SearchTrips(ctx, req).
		Return(expectedTrips, int64(1), nil).
		Times(1)

	trips, total, err := service.SearchTrips(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, trips, 1)
	assert.Equal(t, int64(1), total)
}

func TestSearchTrips_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	req := &model.TripSearchRequest{}

	mockTripRepo.EXPECT().
		SearchTrips(ctx, req).
		Return(nil, int64(0), assert.AnError).
		Times(1)

	trips, total, err := service.SearchTrips(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, trips)
	assert.Equal(t, int64(0), total)
}

func TestGetTripByID_Simple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()
	req := &model.GetTripByIDRequest{}

	expectedTrip := &model.Trip{
		BaseModel: model.BaseModel{ID: tripID},
		BasePrice: 100000,
	}

	mockTripRepo.EXPECT().
		GetTripByID(ctx, req, tripID).
		Return(expectedTrip, nil).
		Times(1)

	result, err := service.GetTripByID(ctx, req, tripID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tripID, result.ID)
}

func TestListTrips_ByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripIDs := []uuid.UUID{uuid.New(), uuid.New()}
	req := &model.ListTripsRequest{IDs: tripIDs}

	expectedTrips := []model.Trip{
		{BaseModel: model.BaseModel{ID: tripIDs[0]}},
		{BaseModel: model.BaseModel{ID: tripIDs[1]}},
	}

	mockTripRepo.EXPECT().
		GetTripsByIDs(ctx, tripIDs).
		Return(expectedTrips, nil).
		Times(1)

	trips, total, err := service.ListTrips(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, trips, 2)
	assert.Equal(t, int64(2), total)
}

func TestListTrips_Pagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	req := &model.ListTripsRequest{
		PaginationRequest: model.PaginationRequest{Page: 1, PageSize: 10},
	}

	expectedTrips := []model.Trip{{BaseModel: model.BaseModel{ID: uuid.New()}}}

	mockTripRepo.EXPECT().
		ListTrips(ctx, 1, 10).
		Return(expectedTrips, int64(1), nil).
		Times(1)

	trips, total, err := service.ListTrips(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, trips, 1)
	assert.Equal(t, int64(1), total)
}

func TestGetSeatAvailability_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()
	busID := uuid.New()

	trip := &model.Trip{
		BaseModel: model.BaseModel{ID: tripID},
		BusID:     busID,
		BasePrice: 100000,
	}

	seats := []model.Seat{
		{
			BaseModel:       model.BaseModel{ID: uuid.New()},
			SeatNumber:      "A1",
			PriceMultiplier: 1.0,
			IsAvailable:     true,
		},
		{
			BaseModel:       model.BaseModel{ID: uuid.New()},
			SeatNumber:      "A2",
			PriceMultiplier: 1.5,
			IsAvailable:     false,
		},
	}

	mockTripRepo.EXPECT().
		GetTripByID(ctx, gomock.Any(), tripID).
		Return(trip, nil).
		Times(1)

	mockSeatRepo.EXPECT().
		GetListByBusID(ctx, busID).
		Return(seats, nil).
		Times(1)

	result, err := service.GetSeatAvailability(ctx, tripID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tripID, result.TripID)
	assert.Equal(t, 1, result.AvailableSeats)
	assert.Equal(t, 2, result.TotalSeats)
	assert.Len(t, result.SeatMap, 2)
}

func TestGetTripsByRouteAndDate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	routeID := uuid.New()
	date := time.Now()

	route := &model.Route{BaseModel: model.BaseModel{ID: routeID}}
	trips := []model.Trip{{BaseModel: model.BaseModel{ID: uuid.New()}}}

	mockRouteRepo.EXPECT().
		GetRouteByID(ctx, routeID).
		Return(route, nil).
		Times(1)

	mockTripRepo.EXPECT().
		GetTripsByRouteAndDate(ctx, routeID, date).
		Return(trips, nil).
		Times(1)

	result, err := service.GetTripsByRouteAndDate(ctx, routeID, date)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetTripsByRouteAndDate_InvalidRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	date := time.Now()

	result, err := service.GetTripsByRouteAndDate(ctx, uuid.Nil, date)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "route ID is required")
}

func TestCreateTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	now := time.Now()
	req := &model.CreateTripRequest{
		RouteID:       uuid.New(),
		BusID:         uuid.New(),
		DepartureTime: now.Add(48 * time.Hour),
		ArrivalTime:   now.Add(60 * time.Hour),
		BasePrice:     100000,
	}

	route := &model.Route{BaseModel: model.BaseModel{ID: req.RouteID}}
	bus := &model.Bus{BaseModel: model.BaseModel{ID: req.BusID}, IsActive: true}
	createdTrip := &model.Trip{BaseModel: model.BaseModel{ID: uuid.New()}}

	mockRouteRepo.EXPECT().GetRouteByID(ctx, req.RouteID).Return(route, nil).Times(1)
	mockBusRepo.EXPECT().GetBusByID(ctx, req.BusID).Return(bus, nil).Times(1)
	mockTripRepo.EXPECT().GetTripsByBusAndDateRange(ctx, req.BusID, gomock.Any(), gomock.Any()).Return([]model.Trip{}, nil).Times(1)
	mockTripRepo.EXPECT().CreateTrip(ctx, gomock.Any()).Return(nil).Times(1)
	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), gomock.Any()).Return(createdTrip, nil).Times(1)

	result, err := service.CreateTrip(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateTrip_ArrivalBeforeDeparture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	now := time.Now()
	req := &model.CreateTripRequest{
		DepartureTime: now.Add(60 * time.Hour),
		ArrivalTime:   now.Add(48 * time.Hour), // Before departure!
	}

	result, err := service.CreateTrip(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "arrival time must be after departure time")
}

func TestCreateTrip_PastDeparture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	past := time.Now().Add(-1 * time.Hour)
	req := &model.CreateTripRequest{
		DepartureTime: past,
		ArrivalTime:   past.Add(2 * time.Hour),
	}

	result, err := service.CreateTrip(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be in the past")
}

func TestUpdateTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()
	newPrice := 150000.0

	existingTrip := &model.Trip{
		BaseModel:     model.BaseModel{ID: tripID},
		DepartureTime: time.Now().Add(48 * time.Hour),
		BasePrice:     100000,
	}

	req := &model.UpdateTripRequest{
		BasePrice: &newPrice,
	}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(existingTrip, nil).Times(2)
	mockTripRepo.EXPECT().UpdateTrip(ctx, gomock.Any()).Return(nil).Times(1)

	result, err := service.UpdateTrip(ctx, tripID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateTrip_NegativePrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()
	negativePrice := -1000.0

	existingTrip := &model.Trip{BaseModel: model.BaseModel{ID: tripID}}
	req := &model.UpdateTripRequest{BasePrice: &negativePrice}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(existingTrip, nil).Times(1)

	result, err := service.UpdateTrip(ctx, tripID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "non-negative")
}

func TestDeleteTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()

	trip := &model.Trip{
		BaseModel:     model.BaseModel{ID: tripID},
		Status:        "scheduled",
		DepartureTime: time.Now().Add(48 * time.Hour),
	}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(trip, nil).Times(1)
	mockTripRepo.EXPECT().DeleteTrip(ctx, tripID).Return(nil).Times(1)

	err := service.DeleteTrip(ctx, tripID)

	assert.NoError(t, err)
}

func TestDeleteTrip_NotScheduled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()

	trip := &model.Trip{
		BaseModel: model.BaseModel{ID: tripID},
		Status:    "completed",
	}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(trip, nil).Times(1)

	err := service.DeleteTrip(ctx, tripID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only scheduled trips")
}

func TestRescheduleTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()
	newDeparture := time.Now().Add(72 * time.Hour)
	newArrival := newDeparture.Add(12 * time.Hour)

	trip := &model.Trip{BaseModel: model.BaseModel{ID: tripID}}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(trip, nil).Times(1)
	mockTripRepo.EXPECT().UpdateTrip(ctx, gomock.Any()).Return(nil).Times(1)

	err := service.RescheduleTrip(ctx, tripID, newDeparture, newArrival)

	assert.NoError(t, err)
}

func TestRescheduleTrip_InvalidTimes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()
	newDeparture := time.Now().Add(72 * time.Hour)
	newArrival := newDeparture.Add(-1 * time.Hour) // Before departure!

	trip := &model.Trip{BaseModel: model.BaseModel{ID: tripID}}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(trip, nil).Times(1)

	err := service.RescheduleTrip(ctx, tripID, newDeparture, newArrival)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "arrival time must be after departure time")
}

func TestGetCompletedTripsForReschedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	expectedTrips := []model.Trip{{BaseModel: model.BaseModel{ID: uuid.New()}}}

	mockTripRepo.EXPECT().
		GetCompletedTripsForReschedule(ctx).
		Return(expectedTrips, nil).
		Times(1)

	result, err := service.GetCompletedTripsForReschedule(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetTripByID_WithSeatStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()
	seatID1 := uuid.New()
	seatID2 := uuid.New()

	req := &model.GetTripByIDRequest{
		SeatBookingStatus: true,
		PreloadBus:        true,
		PreloadSeat:       true,
	}

	seats := []model.Seat{
		{BaseModel: model.BaseModel{ID: seatID1}, SeatNumber: "A1"},
		{BaseModel: model.BaseModel{ID: seatID2}, SeatNumber: "A2"},
	}

	trip := &model.Trip{
		BaseModel: model.BaseModel{ID: tripID},
		Bus: &model.Bus{
			BaseModel: model.BaseModel{ID: uuid.New()},
			Seats:     seats,
		},
	}

	seatStatuses := []booking.SeatStatus{
		{SeatID: seatID1, IsBooked: true, IsLocked: false},
		// seatID2 is missing (should default to available)
	}

	mockTripRepo.EXPECT().GetTripByID(ctx, req, tripID).Return(trip, nil).Times(1)
	mockBookingClient.EXPECT().GetSeatStatus(ctx, tripID, gomock.Any()).Return(seatStatuses, nil).Times(1)

	result, err := service.GetTripByID(ctx, req, tripID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Bus)
	assert.Len(t, result.Bus.Seats, 2)

	// Check mapped status
	assert.NotNil(t, result.Bus.Seats[0].Status)
	assert.True(t, result.Bus.Seats[0].Status.IsBooked)

	// Check default status
	assert.NotNil(t, result.Bus.Seats[1].Status)
	assert.False(t, result.Bus.Seats[1].Status.IsBooked)
}

func TestUpdateTrip_Validations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()

	// Setup existing trip
	now := time.Now()
	existingTrip := &model.Trip{
		BaseModel:     model.BaseModel{ID: tripID},
		DepartureTime: now.Add(24 * time.Hour),
		ArrivalTime:   now.Add(30 * time.Hour),
	}

	// Case 1: Past Departure Time
	pastTime := now.Add(-1 * time.Hour)
	req1 := &model.UpdateTripRequest{DepartureTime: &pastTime}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(existingTrip, nil).Times(1)
	_, err1 := service.UpdateTrip(ctx, tripID, req1)
	assert.Error(t, err1)
	assert.Contains(t, err1.Error(), "past")

	// Case 2: Arrival Before Departure (Departure updated)
	futureDep := now.Add(48 * time.Hour)
	futureArr := now.Add(40 * time.Hour) // Before new departure
	req2 := &model.UpdateTripRequest{
		DepartureTime: &futureDep,
		ArrivalTime:   &futureArr,
	}
	// Note: Logic inside UpdateTrip checks arrival vs existing/new departure
	// But actually UpdateTrip implementation checks req.ArrivalTime vs trip.DepartureTime independently if updated separately
	// Wait, let's look at the implementation logic again:
	/*
		if req.DepartureTime != nil { trip.DepartureTime = *req.DepartureTime }
		if req.ArrivalTime != nil {
			if req.ArrivalTime.Before(trip.DepartureTime) { return error }
			trip.ArrivalTime = *req.ArrivalTime
		}
	*/
	// So if both provided, trip.DepartureTime is updated first, then ArrivalTime check uses NEW value. Correct.

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(existingTrip, nil).Times(1)
	_, err2 := service.UpdateTrip(ctx, tripID, req2)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "after departure")

	// Case 3: Arrival Before Departure (Only Arrival updated)
	badArr := now.Add(20 * time.Hour) // Before existing departure (24h)
	req3 := &model.UpdateTripRequest{ArrivalTime: &badArr}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(existingTrip, nil).Times(1)
	_, err3 := service.UpdateTrip(ctx, tripID, req3)
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "after departure")
}

func TestUpdateTrip_FullUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockRouteRepo := repo_mocks.NewMockRouteRepository(ctrl)
	mockRouteStopRepo := repo_mocks.NewMockRouteStopRepository(ctrl)
	mockBusRepo := repo_mocks.NewMockBusRepository(ctrl)
	mockSeatRepo := repo_mocks.NewMockSeatRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)

	service := NewTripService(mockTripRepo, mockRouteRepo, mockRouteStopRepo, mockBusRepo, mockSeatRepo, mockBookingClient, nil)

	ctx := context.Background()
	tripID := uuid.New()
	existingTrip := &model.Trip{BaseModel: model.BaseModel{ID: tripID}}

	status := constants.TripStatusCancelled
	isActive := false

	req := &model.UpdateTripRequest{
		Status:   &status,
		IsActive: &isActive,
	}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(existingTrip, nil).Times(2)
	mockTripRepo.EXPECT().UpdateTrip(ctx, gomock.Any()).Do(func(_ context.Context, tr *model.Trip) {
		assert.Equal(t, constants.TripStatusCancelled, tr.Status)
		assert.False(t, tr.IsActive)
	}).Return(nil).Times(1)

	result, err := service.UpdateTrip(ctx, tripID, req)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestProcessTripStatusUpdates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	service := &TripServiceImpl{tripRepo: mockTripRepo}

	ctx := context.Background()

	mockTripRepo.EXPECT().UpdateTripStatuses(ctx).Return(nil).Times(1)

	err := service.ProcessTripStatusUpdates(ctx)
	assert.NoError(t, err)
}

func TestCancelTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	mockBookingClient := mocks.NewMockBookingClient(ctrl)
	mockPaymentClient := mocks.NewMockPaymentClient(ctrl)

	service := &TripServiceImpl{
		tripRepo:      mockTripRepo,
		bookingClient: mockBookingClient,
		paymentClient: mockPaymentClient,
	}

	ctx := context.Background()
	tripID := uuid.New()

	trip := &model.Trip{
		BaseModel: model.BaseModel{ID: tripID},
		Status:    constants.TripStatusScheduled,
	}

	bookings := []*booking.Booking{
		{
			ID:                uuid.New(),
			TransactionStatus: "PAID",
			TotalAmount:       100000,
		},
		{
			ID:                uuid.New(),
			TransactionStatus: "PENDING",
		},
	}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(trip, nil).Times(1)
	mockTripRepo.EXPECT().UpdateTrip(ctx, gomock.Any()).Do(func(_ context.Context, tr *model.Trip) {
		assert.Equal(t, constants.TripStatusCancelled, tr.Status)
	}).Return(nil).Times(1)

	mockBookingClient.EXPECT().GetTripBookings(ctx, tripID).Return(bookings, nil).Times(1)

	// Expect refund for PAID booking
	mockPaymentClient.EXPECT().CreateRefund(ctx, gomock.Any()).Return(nil, nil).Times(1)

	// Expect cancel for both bookings
	mockBookingClient.EXPECT().CancelBooking(ctx, gomock.Any(), gomock.Any()).Return(nil).Times(2)

	err := service.CancelTrip(ctx, tripID)
	assert.NoError(t, err)
}

func TestCancelTrip_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripRepo := repo_mocks.NewMockTripRepository(ctrl)
	service := &TripServiceImpl{tripRepo: mockTripRepo}

	ctx := context.Background()
	tripID := uuid.New()

	trip := &model.Trip{
		BaseModel: model.BaseModel{ID: tripID},
		Status:    constants.TripStatusCompleted, // Cannot cancel completed
	}

	mockTripRepo.EXPECT().GetTripByID(ctx, gomock.Any(), tripID).Return(trip, nil).Times(1)

	err := service.CancelTrip(ctx, tripID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot cancel trip")
}
