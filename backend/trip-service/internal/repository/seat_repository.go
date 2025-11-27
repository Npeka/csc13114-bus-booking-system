package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatRepository interface {
	Create(ctx context.Context, seat *model.Seat) error
	CreateBulk(ctx context.Context, seats []model.Seat) error
	Update(ctx context.Context, seat *model.Seat) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Seat, error)
	ListByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error)
	GetSeatMap(ctx context.Context, busID uuid.UUID) ([]model.Seat, error)
	CountByBusID(ctx context.Context, busID uuid.UUID) (int64, error)
}

type SeatRepositoryImpl struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &SeatRepositoryImpl{db: db}
}

func (r *SeatRepositoryImpl) Create(ctx context.Context, seat *model.Seat) error {
	return r.db.WithContext(ctx).Create(seat).Error
}

func (r *SeatRepositoryImpl) CreateBulk(ctx context.Context, seats []model.Seat) error {
	return r.db.WithContext(ctx).Create(&seats).Error
}

func (r *SeatRepositoryImpl) Update(ctx context.Context, seat *model.Seat) error {
	return r.db.WithContext(ctx).Save(seat).Error
}

func (r *SeatRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Seat{}, id).Error
}

func (r *SeatRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Seat, error) {
	var seat model.Seat
	err := r.db.WithContext(ctx).
		Preload("Bus").
		First(&seat, id).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepositoryImpl) ListByBusID(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
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

func (r *SeatRepositoryImpl) GetSeatMap(ctx context.Context, busID uuid.UUID) ([]model.Seat, error) {
	var seats []model.Seat
	err := r.db.WithContext(ctx).
		Where("bus_id = ? AND is_available = ?", busID, true).
		Order("floor ASC, row ASC, \"column\" ASC").
		Find(&seats).Error
	if err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *SeatRepositoryImpl) CountByBusID(ctx context.Context, busID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Seat{}).
		Where("bus_id = ?", busID).
		Count(&count).Error
	return count, err
}
