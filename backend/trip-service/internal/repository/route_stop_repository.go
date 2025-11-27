package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteStopRepository interface {
	Create(ctx context.Context, stop *model.RouteStop) error
	Update(ctx context.Context, stop *model.RouteStop) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.RouteStop, error)
	ListByRouteID(ctx context.Context, routeID uuid.UUID) ([]model.RouteStop, error)
	ReorderStops(ctx context.Context, routeID uuid.UUID, stopOrders map[uuid.UUID]int) error
}

type RouteStopRepositoryImpl struct {
	db *gorm.DB
}

func NewRouteStopRepository(db *gorm.DB) RouteStopRepository {
	return &RouteStopRepositoryImpl{db: db}
}

func (r *RouteStopRepositoryImpl) Create(ctx context.Context, stop *model.RouteStop) error {
	return r.db.WithContext(ctx).Create(stop).Error
}

func (r *RouteStopRepositoryImpl) Update(ctx context.Context, stop *model.RouteStop) error {
	return r.db.WithContext(ctx).Save(stop).Error
}

func (r *RouteStopRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.RouteStop{}, id).Error
}

func (r *RouteStopRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.RouteStop, error) {
	var stop model.RouteStop
	err := r.db.WithContext(ctx).
		Preload("Route").
		First(&stop, id).Error
	if err != nil {
		return nil, err
	}
	return &stop, nil
}

func (r *RouteStopRepositoryImpl) ListByRouteID(ctx context.Context, routeID uuid.UUID) ([]model.RouteStop, error) {
	var stops []model.RouteStop
	err := r.db.WithContext(ctx).
		Where("route_id = ? AND is_active = ?", routeID, true).
		Order("stop_order ASC").
		Find(&stops).Error
	if err != nil {
		return nil, err
	}
	return stops, nil
}

func (r *RouteStopRepositoryImpl) ReorderStops(ctx context.Context, routeID uuid.UUID, stopOrders map[uuid.UUID]int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for stopID, order := range stopOrders {
			if err := tx.Model(&model.RouteStop{}).
				Where("id = ? AND route_id = ?", stopID, routeID).
				Update("stop_order", order).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
