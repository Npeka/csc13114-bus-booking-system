package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatRepository interface {
	CreateSeats(ctx context.Context, seats []model.Seat) error
	GetSeatsByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error)
	GetSeatByID(ctx context.Context, id uuid.UUID) (*model.Seat, error)
	UpdateSeat(ctx context.Context, seat *model.Seat) error
	DeleteSeat(ctx context.Context, id uuid.UUID) error
}

type SeatRepositoryImpl struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &SeatRepositoryImpl{db: db}
}

func (r *SeatRepositoryImpl) CreateSeats(ctx context.Context, seats []model.Seat) error {
	return r.db.WithContext(ctx).Create(&seats).Error
}

func (r *SeatRepositoryImpl) GetSeatsByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
	var seats []model.Seat
	err := r.db.WithContext(ctx).
		Where("bus_id = ? AND is_active = ?", busID, true).
		Order("seat_code ASC").
		Find(&seats).Error
	return seats, err
}

func (r *SeatRepositoryImpl) GetSeatByID(ctx context.Context, id uuid.UUID) (*model.Seat, error) {
	var seat model.Seat
	err := r.db.WithContext(ctx).First(&seat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepositoryImpl) UpdateSeat(ctx context.Context, seat *model.Seat) error {
	return r.db.WithContext(ctx).Model(seat).Updates(seat).Error
}

func (r *SeatRepositoryImpl) DeleteSeat(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Seat{}, "id = ?", id).Error
}
