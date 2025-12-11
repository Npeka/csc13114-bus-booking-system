package mocks

import (
	"context"
	"time"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockTripRepository is a mock implementation of TripRepository
type MockTripRepository struct {
	mock.Mock
}

func (m *MockTripRepository) Create(ctx context.Context, trip *model.Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockTripRepository) GetTripByID(ctx context.Context, req *model.GetTripByIDRequest, id uuid.UUID) (*model.Trip, error) {
	args := m.Called(ctx, req, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Trip), args.Error(1)
}

func (m *MockTripRepository) GetTripsByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Trip, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]model.Trip), args.Error(1)
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
