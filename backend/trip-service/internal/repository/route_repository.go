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
	ListRoutes(ctx context.Context, req *model.ListRoutesRequest) ([]model.Route, int64, error)
	GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error)

	Create(ctx context.Context, route *model.Route) error
	Update(ctx context.Context, route *model.Route) error
	Delete(ctx context.Context, id uuid.UUID) error
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
		return db.Order("stop_order ASC")
	}).First(&route, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepositoryImpl) ListRoutes(ctx context.Context, req *model.ListRoutesRequest) ([]model.Route, int64, error) {
	var routes []model.Route
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Route{})

	// Apply filters
	if req.Origin != nil && *req.Origin != "" {
		query = query.Where("origin = ?", *req.Origin)
	}
	if req.Destination != nil && *req.Destination != "" {
		query = query.Where("destination = ?", *req.Destination)
	}
	if req.MinDistance != nil {
		query = query.Where("distance_km >= ?", *req.MinDistance)
	}
	if req.MaxDistance != nil {
		query = query.Where("distance_km <= ?", *req.MaxDistance)
	}
	if req.MinDuration != nil {
		query = query.Where("estimated_minutes >= ?", *req.MinDuration)
	}
	if req.MaxDuration != nil {
		query = query.Where("estimated_minutes <= ?", *req.MaxDuration)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Count total
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderBy := "created_at DESC" // default
	if req.SortBy != nil && *req.SortBy != "" {
		sortOrder := "ASC"
		if req.SortOrder != nil && *req.SortOrder == "desc" {
			sortOrder = "DESC"
		}

		switch *req.SortBy {
		case "distance":
			orderBy = "distance_km " + sortOrder
		case "duration":
			orderBy = "estimated_minutes " + sortOrder
		case "origin":
			orderBy = "origin " + sortOrder
		case "destination":
			orderBy = "destination " + sortOrder
		default:
			orderBy = "created_at DESC"
		}
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Order(orderBy).Find(&routes).Error

	return routes, total, err
}

func (r *RouteRepositoryImpl) GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error) {
	var routes []model.Route
	err := r.db.WithContext(ctx).
		Where("origin = ? AND destination = ? AND is_active = ?", origin, destination, true).
		Find(&routes).Error
	return routes, err
}

func (r *RouteRepositoryImpl) Create(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *RouteRepositoryImpl) Update(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Model(route).Updates(route).Error
}

func (r *RouteRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
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
