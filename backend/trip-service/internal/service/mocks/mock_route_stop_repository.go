package mocks

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

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
