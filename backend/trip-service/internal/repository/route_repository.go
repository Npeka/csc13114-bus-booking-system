package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteRepository interface {
	GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error)
	GetRoutesWithRouteStops(ctx context.Context, id uuid.UUID) (*model.Route, error)
	ListRoutes(ctx context.Context, page, limit int) ([]model.Route, int64, error)
	GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error)

	CreateRoute(ctx context.Context, route *model.Route) error
	UpdateRoute(ctx context.Context, route *model.Route) error
	DeleteRoute(ctx context.Context, id uuid.UUID) error
}

type RouteRepositoryImpl struct {
	db *gorm.DB
}

func NewRouteRepository(db *gorm.DB) RouteRepository {
	return &RouteRepositoryImpl{db: db}
}

func (r *RouteRepositoryImpl) GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	var route model.Route
	if err := r.db.WithContext(ctx).First(&route, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepositoryImpl) GetRoutesWithRouteStops(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	var route model.Route
	if err := r.db.WithContext(ctx).Preload("RouteStops", func(db *gorm.DB) *gorm.DB {
		return db.Order("stop_number ASC")
	}).First(&route, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepositoryImpl) ListRoutes(ctx context.Context, page, pageSize int) ([]model.Route, int64, error) {
	var routes []model.Route
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Route{})

	// Count total
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&routes).Error

	return routes, total, err
}

func (r *RouteRepositoryImpl) GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error) {
	var routes []model.Route
	err := r.db.WithContext(ctx).
		Where("origin = ? AND destination = ? AND is_active = ?", origin, destination, true).
		Find(&routes).Error
	return routes, err
}

func (r *RouteRepositoryImpl) CreateRoute(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *RouteRepositoryImpl) UpdateRoute(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Model(route).Updates(route).Error
}

func (r *RouteRepositoryImpl) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Route{}, "id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.RouteStop{}, "route_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
