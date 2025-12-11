package mocks

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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
