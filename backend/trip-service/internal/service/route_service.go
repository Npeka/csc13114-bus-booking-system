package service

import (
	"context"
	"errors"
	"fmt"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RouteService interface {
	GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error)
	ListRoutes(ctx context.Context, page, limit int) ([]model.RouteSummary, int64, error)
	GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error)

	CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (*model.Route, error)
	UpdateRoute(ctx context.Context, id uuid.UUID, req *model.UpdateRouteRequest) (*model.Route, error)
	DeleteRoute(ctx context.Context, id uuid.UUID) error
}

type RouteServiceImpl struct {
	routeRepo repository.RouteRepository
}

func NewRouteService(routeRepo repository.RouteRepository) RouteService {
	return &RouteServiceImpl{
		routeRepo: routeRepo,
	}
}

func (s *RouteServiceImpl) CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (*model.Route, error) {
	log.Info().Msg("Creating new route")

	// Validate request
	if req.Origin == "" || req.Destination == "" {
		return nil, errors.New("origin and destination are required")
	}

	if req.DistanceKm <= 0 || req.EstimatedMinutes <= 0 {
		return nil, errors.New("distance and estimated time must be positive")
	}

	route := &model.Route{
		Origin:           req.Origin,
		Destination:      req.Destination,
		DistanceKm:       req.DistanceKm,
		EstimatedMinutes: req.EstimatedMinutes,
		IsActive:         true,
	}

	if err := s.routeRepo.CreateRoute(ctx, route); err != nil {
		log.Error().Err(err).Msg("Failed to create route")
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	log.Info().Str("route_id", route.ID.String()).Msg("Route created successfully")
	return route, nil
}

func (s *RouteServiceImpl) GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	route, err := s.routeRepo.GetRouteByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("route_id", id.String()).Msg("Failed to get route")
		return nil, fmt.Errorf("failed to get route: %w", err)
	}

	return route, nil
}

func (s *RouteServiceImpl) UpdateRoute(ctx context.Context, id uuid.UUID, req *model.UpdateRouteRequest) (*model.Route, error) {
	log.Info().Str("route_id", id.String()).Msg("Updating route")

	// Get existing route
	route, err := s.routeRepo.GetRouteByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	// Update fields
	if req.Origin != nil {
		route.Origin = *req.Origin
	}
	if req.Destination != nil {
		route.Destination = *req.Destination
	}
	if req.DistanceKm != nil {
		if *req.DistanceKm <= 0 {
			return nil, errors.New("distance must be positive")
		}
		route.DistanceKm = *req.DistanceKm
	}
	if req.EstimatedMinutes != nil {
		if *req.EstimatedMinutes <= 0 {
			return nil, errors.New("estimated time must be positive")
		}
		route.EstimatedMinutes = *req.EstimatedMinutes
	}
	if req.IsActive != nil {
		route.IsActive = *req.IsActive
	}

	if err := s.routeRepo.UpdateRoute(ctx, route); err != nil {
		log.Error().Err(err).Str("route_id", id.String()).Msg("Failed to update route")
		return nil, fmt.Errorf("failed to update route: %w", err)
	}

	log.Info().Str("route_id", id.String()).Msg("Route updated successfully")
	return route, nil
}

func (s *RouteServiceImpl) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	log.Info().Str("route_id", id.String()).Msg("Deleting route")

	if err := s.routeRepo.DeleteRoute(ctx, id); err != nil {
		log.Error().Err(err).Str("route_id", id.String()).Msg("Failed to delete route")
		return fmt.Errorf("failed to delete route: %w", err)
	}

	log.Info().Str("route_id", id.String()).Msg("Route deleted successfully")
	return nil
}

func (s *RouteServiceImpl) ListRoutes(ctx context.Context, page, pageSize int) ([]model.RouteSummary, int64, error) {
	routes, total, err := s.routeRepo.ListRoutes(ctx, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list routes")
		return nil, 0, fmt.Errorf("failed to list routes: %w", err)
	}

	return routes, total, nil
}

func (s *RouteServiceImpl) GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error) {
	if origin == "" || destination == "" {
		return nil, errors.New("origin and destination are required")
	}

	routes, err := s.routeRepo.GetRoutesByOriginDestination(ctx, origin, destination)
	if err != nil {
		log.Error().Err(err).Str("origin", origin).Str("destination", destination).Msg("Failed to get routes by origin/destination")
		return nil, fmt.Errorf("failed to get routes: %w", err)
	}

	return routes, nil
}
