package mocks

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

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
