package mocks

import (
	"context"
	"time"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockBusRepository is a mock implementation of BusRepository
type MockBusRepository struct {
	mock.Mock
}

func (m *MockBusRepository) CreateBus(ctx context.Context, bus *model.Bus) error {
	args := m.Called(ctx, bus)
	return args.Error(0)
}

func (m *MockBusRepository) GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bus), args.Error(1)
}

func (m *MockBusRepository) GetBusWithSeatsByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bus), args.Error(1)
}

func (m *MockBusRepository) GetBusByPlateNumber(ctx context.Context, plateNumber string) (*model.Bus, error) {
	args := m.Called(ctx, plateNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Bus), args.Error(1)
}

func (m *MockBusRepository) ListBuses(ctx context.Context, page, pageSize int) ([]model.Bus, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.Bus), args.Get(1).(int64), args.Error(2)
}

func (m *MockBusRepository) UpdateBus(ctx context.Context, bus *model.Bus) error {
	args := m.Called(ctx, bus)
	return args.Error(0)
}

func (m *MockBusRepository) DeleteBus(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSeatRepository is a mock implementation of SeatRepository
type MockSeatRepository struct {
	mock.Mock
}

func (m *MockSeatRepository) Create(ctx context.Context, seat *model.Seat) error {
	args := m.Called(ctx, seat)
	return args.Error(0)
}

func (m *MockSeatRepository) CreateBulk(ctx context.Context, seats []model.Seat) error {
	args := m.Called(ctx, seats)
	return args.Error(0)
}

func (m *MockSeatRepository) CreateWithTx(ctx context.Context, seat *model.Seat, tx *gorm.DB) error {
	args := m.Called(ctx, seat, tx)
	return args.Error(0)
}

func (m *MockSeatRepository) CreateBulkWithTx(ctx context.Context, seats []model.Seat, tx *gorm.DB) error {
	args := m.Called(ctx, seats, tx)
	return args.Error(0)
}

func (m *MockSeatRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Seat, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Seat), args.Error(1)
}

func (m *MockSeatRepository) ListByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
	args := m.Called(ctx, busID)
	return args.Get(0).([]model.Seat), args.Error(1)
}

func (m *MockSeatRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Seat, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]model.Seat), args.Error(1)
}

func (m *MockSeatRepository) GetSeatMap(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
	args := m.Called(ctx, busID)
	return args.Get(0).([]model.Seat), args.Error(1)
}

func (m *MockSeatRepository) CountByBusID(ctx context.Context, busID uuid.UUID) (int64, error) {
	args := m.Called(ctx, busID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSeatRepository) Update(ctx context.Context, seat *model.Seat) error {
	args := m.Called(ctx, seat)
	return args.Error(0)
}

func (m *MockSeatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSeatRepository) BulkUpdateAvailability(ctx context.Context, seatIDs []uuid.UUID, isAvailable bool) error {
	args := m.Called(ctx, seatIDs, isAvailable)
	return args.Error(0)
}

// MockRouteRepository is a mock implementation of RouteRepository
type MockRouteRepository struct {
	mock.Mock
}

func (m *MockRouteRepository) CreateRoute(ctx context.Context, route *model.Route) error {
	args := m.Called(ctx, route)
	return args.Error(0)
}

func (m *MockRouteRepository) GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Route), args.Error(1)
}

func (m *MockRouteRepository) GetRoutesWithRouteStops(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Route), args.Error(1)
}

func (m *MockRouteRepository) ListRoutes(ctx context.Context, req *model.ListRoutesRequest) ([]model.Route, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]model.Route), args.Get(1).(int64), args.Error(2)
}

func (m *MockRouteRepository) GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error) {
	args := m.Called(ctx, origin, destination)
	return args.Get(0).([]model.Route), args.Error(1)
}

func (m *MockRouteRepository) UpdateRoute(ctx context.Context, route *model.Route) error {
	args := m.Called(ctx, route)
	return args.Error(0)
}

func (m *MockRouteRepository) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockRouteStopRepository is a mock implementation of RouteStopRepository
type MockRouteStopRepository struct {
	mock.Mock
}

func (m *MockRouteStopRepository) Create(ctx context.Context, routeStop *model.RouteStop) error {
	args := m.Called(ctx, routeStop)
	return args.Error(0)
}

func (m *MockRouteStopRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.RouteStop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RouteStop), args.Error(1)
}

func (m *MockRouteStopRepository) ListByRouteID(ctx context.Context, routeID uuid.UUID) ([]model.RouteStop, error) {
	args := m.Called(ctx, routeID)
	return args.Get(0).([]model.RouteStop), args.Error(1)
}

func (m *MockRouteStopRepository) Update(ctx context.Context, routeStop *model.RouteStop) error {
	args := m.Called(ctx, routeStop)
	return args.Error(0)
}

func (m *MockRouteStopRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRouteStopRepository) ReorderStops(ctx context.Context, routeID uuid.UUID, stopOrders map[uuid.UUID]int) error {
	args := m.Called(ctx, routeID, stopOrders)
	return args.Error(0)
}

// MockTripRepository is a mock implementation of TripRepository
type MockTripRepository struct {
	mock.Mock
}

func (m *MockTripRepository) Create(ctx context.Context, trip *model.Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockTripRepository) GetTripByID(ctx context.Context, req *model.GetTripByIDRequuest, id uuid.UUID) (*model.Trip, error) {
	args := m.Called(ctx, req, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Trip), args.Error(1)
}

func (m *MockTripRepository) CreateTrip(ctx context.Context, trip *model.Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockTripRepository) ListTrips(ctx context.Context, page, pageSize int) ([]model.Trip, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.Trip), args.Get(1).(int64), args.Error(2)
}

func (m *MockTripRepository) SearchTrips(ctx context.Context, req *model.TripSearchRequest) ([]model.TripDetail, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]model.TripDetail), args.Get(1).(int64), args.Error(2)
}

func (m *MockTripRepository) GetTripsByRouteAndDate(ctx context.Context, routeID uuid.UUID, date time.Time) ([]model.Trip, error) {
	args := m.Called(ctx, routeID, date)
	return args.Get(0).([]model.Trip), args.Error(1)
}

func (m *MockTripRepository) GetTripsByBusAndDateRange(ctx context.Context, busID uuid.UUID, startDate, endDate time.Time) ([]model.Trip, error) {
	args := m.Called(ctx, busID, startDate, endDate)
	return args.Get(0).([]model.Trip), args.Error(1)
}

func (m *MockTripRepository) UpdateTrip(ctx context.Context, trip *model.Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockTripRepository) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
