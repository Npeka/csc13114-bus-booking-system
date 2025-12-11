package mocks

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
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
