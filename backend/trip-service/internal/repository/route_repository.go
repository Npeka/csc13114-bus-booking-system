package repository

import (
	"context"

	"bus-booking/trip-service/internal/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RouteRepository interface {
	CreateRoute(ctx context.Context, route *model.Route) error
	GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error)
	UpdateRoute(ctx context.Context, route *model.Route) error
	DeleteRoute(ctx context.Context, id uuid.UUID) error
	ListRoutes(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.RouteSummary, int64, error)
	GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error)
}

type RouteRepositoryImpl struct {
	db *gorm.DB
}

func NewRouteRepository(db *gorm.DB) RouteRepository {
	return &RouteRepositoryImpl{db: db}
}

func (r *RouteRepositoryImpl) CreateRoute(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *RouteRepositoryImpl) GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	var route model.Route
	err := r.db.WithContext(ctx).
		Preload("Operator").
		First(&route, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepositoryImpl) UpdateRoute(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Model(route).Updates(route).Error
}

func (r *RouteRepositoryImpl) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Route{}, "id = ?", id).Error
}

func (r *RouteRepositoryImpl) ListRoutes(ctx context.Context, operatorID *uuid.UUID, page, limit int) ([]model.RouteSummary, int64, error) {
	var results []model.RouteSummary
	var total int64

	query := r.db.WithContext(ctx).Table("routes r").
		Select(`
			r.id, r.origin, r.destination, r.distance_km, r.estimated_minutes,
			r.is_active, r.created_at, o.name as operator_name,
			COUNT(t.id) as active_trips
		`).
		Joins("JOIN operators o ON r.operator_id = o.id").
		Joins("LEFT JOIN trips t ON r.id = t.route_id AND t.is_active = true AND t.departure_time > NOW()").
		Group("r.id, o.name")

	if operatorID != nil {
		query = query.Where("r.operator_id = ?", *operatorID)
	}

	// Count total
	countQuery := r.db.WithContext(ctx).Model(&model.Route{})
	if operatorID != nil {
		countQuery = countQuery.Where("operator_id = ?", *operatorID)
	}
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
	query = query.Offset(offset).Limit(limit).Order("r.created_at DESC")

	// Execute query
	rows, err := query.Rows()
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close rows")
		}
	}()

	for rows.Next() {
		var result model.RouteSummary
		err := rows.Scan(
			&result.ID, &result.Origin, &result.Destination, &result.DistanceKm, &result.EstimatedMinutes,
			&result.IsActive, &result.CreatedAt, &result.OperatorName, &result.ActiveTrips,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}

	return results, total, nil
}

func (r *RouteRepositoryImpl) GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error) {
	var routes []model.Route
	err := r.db.WithContext(ctx).
		Where("origin = ? AND destination = ? AND is_active = ?", origin, destination, true).
		Preload("Operator").
		Find(&routes).Error
	return routes, err
}
