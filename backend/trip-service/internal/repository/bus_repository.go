package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusRepository interface {
	CreateBus(ctx context.Context, bus *model.Bus) error
	GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error)
	UpdateBus(ctx context.Context, bus *model.Bus) error
	DeleteBus(ctx context.Context, id uuid.UUID) error
	ListBuses(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.Bus, int64, error)
	GetBusByPlateNumber(ctx context.Context, plateNumber string) (*model.Bus, error)
}

type BusRepositoryImpl struct {
	db *gorm.DB
}

func NewBusRepository(db *gorm.DB) BusRepository {
	return &BusRepositoryImpl{db: db}
}

func (r *BusRepositoryImpl) CreateBus(ctx context.Context, bus *model.Bus) error {
	return r.db.WithContext(ctx).Create(bus).Error
}

func (r *BusRepositoryImpl) GetBusByID(ctx context.Context, id uuid.UUID) (*model.Bus, error) {
	var bus model.Bus
	err := r.db.WithContext(ctx).Preload("Operator").First(&bus, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *BusRepositoryImpl) UpdateBus(ctx context.Context, bus *model.Bus) error {
	return r.db.WithContext(ctx).Model(bus).Updates(bus).Error
}

func (r *BusRepositoryImpl) DeleteBus(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Bus{}, "id = ?", id).Error
}

func (r *BusRepositoryImpl) ListBuses(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.Bus, int64, error) {
	var buses []model.Bus
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Bus{}).Preload("Operator")
	if operatorID != nil {
		query = query.Where("operator_id = ?", *operatorID)
	}

	// Count total
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&buses).Error

	return buses, total, err
}

func (r *BusRepositoryImpl) GetBusByPlateNumber(ctx context.Context, plateNumber string) (*model.Bus, error) {
	var bus model.Bus
	err := r.db.WithContext(ctx).
		Preload("Operator").
		Where("plate_number = ?", plateNumber).
		First(&bus).Error
	if err != nil {
		return nil, err
	}
	return &bus, nil
}
