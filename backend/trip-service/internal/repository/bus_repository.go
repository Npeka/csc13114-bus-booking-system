package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusRepository interface {
	GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error)
	GetBusWithSeatsByID(ctx context.Context, id uuid.UUID) (*model.Bus, error)
	ListBuses(ctx context.Context, page, pageSize int) ([]model.Bus, int64, error)
	GetBusByPlateNumber(ctx context.Context, plateNumber string) (*model.Bus, error)
	CreateBus(ctx context.Context, bus *model.Bus) error
	UpdateBus(ctx context.Context, bus *model.Bus) error
	DeleteBus(ctx context.Context, id uuid.UUID) error
}

type BusRepositoryImpl struct {
	db *gorm.DB
}

func NewBusRepository(db *gorm.DB) BusRepository {
	return &BusRepositoryImpl{db: db}
}

func (r *BusRepositoryImpl) GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	var bus model.Bus
	if err := r.db.WithContext(ctx).First(&bus, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *BusRepositoryImpl) GetBusWithSeatsByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	var bus model.Bus
	if err := r.db.WithContext(ctx).Model(&model.Bus{}).
		Preload("Seats").
		First(&bus, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *BusRepositoryImpl) ListBuses(ctx context.Context, page, pageSize int) ([]model.Bus, int64, error) {
	var buses []model.Bus
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Bus{})

	// Count total
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&buses).Error

	return buses, total, err
}

func (r *BusRepositoryImpl) GetBusByPlateNumber(ctx context.Context, plateNumber string) (*model.Bus, error) {
	var bus model.Bus
	if err := r.db.WithContext(ctx).
		Where("plate_number = ?", plateNumber).
		First(&bus).Error; err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *BusRepositoryImpl) CreateBus(ctx context.Context, bus *model.Bus) error {
	return r.db.WithContext(ctx).Create(bus).Error
}

func (r *BusRepositoryImpl) UpdateBus(ctx context.Context, bus *model.Bus) error {
	return r.db.WithContext(ctx).Model(bus).Updates(bus).Error
}

func (r *BusRepositoryImpl) DeleteBus(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Seat{}, "bus_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.Bus{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
