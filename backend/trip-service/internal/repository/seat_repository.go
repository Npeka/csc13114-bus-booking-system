package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Seat, error)
	GetListByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Seat, error)
	GetListByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error)
	Update(ctx context.Context, seat *model.Seat) error
}

type SeatRepositoryImpl struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &SeatRepositoryImpl{db: db}
}

func (r *SeatRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Seat, error) {
	var seat model.Seat
	if err := r.db.WithContext(ctx).First(&seat, id).Error; err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepositoryImpl) GetListByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Seat, error) {
	var seats []model.Seat
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *SeatRepositoryImpl) GetListByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
	var seats []model.Seat
	err := r.db.WithContext(ctx).
		Where("bus_id = ?", busID).
		Order("floor ASC, row ASC, \"column\" ASC").
		Find(&seats).Error
	if err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *SeatRepositoryImpl) Update(ctx context.Context, seat *model.Seat) error {
	return r.db.WithContext(ctx).Save(seat).Error
}
