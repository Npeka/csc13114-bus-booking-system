package service

import (
	"context"
	"sort"

	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RouteService interface {
	GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error)
	ListRoutes(ctx context.Context, page, pageSize int) ([]model.Route, int64, error)
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

func (s *RouteServiceImpl) GetRouteByID(ctx context.Context, id uuid.UUID) (*model.Route, error) {
	route, err := s.routeRepo.GetRoutesWithRouteStops(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("route_id", id.String()).Msg("Failed to get route")
		return nil, ginext.NewInternalServerError("failed to get route")
	}
	return route, nil
}

func (s *RouteServiceImpl) ListRoutes(ctx context.Context, page, pageSize int) ([]model.Route, int64, error) {
	routes, total, err := s.routeRepo.ListRoutes(ctx, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list routes")
		return nil, 0, ginext.NewInternalServerError("failed to list routes")
	}
	return routes, total, nil
}

func (s *RouteServiceImpl) GetRoutesByOriginDestination(ctx context.Context, origin, destination string) ([]model.Route, error) {
	routes, err := s.routeRepo.GetRoutesByOriginDestination(ctx, origin, destination)
	if err != nil {
		log.Error().Err(err).Str("origin", origin).Str("destination", destination).Msg("Failed to get routes by origin/destination")
		return nil, ginext.NewInternalServerError("failed to get routes by origin/destination")
	}
	return routes, nil
}

func (s *RouteServiceImpl) CreateRoute(ctx context.Context, req *model.CreateRouteRequest) (*model.Route, error) {
	// Create route
	route := &model.Route{
		Origin:           req.Origin,
		Destination:      req.Destination,
		DistanceKm:       req.DistanceKm,
		EstimatedMinutes: req.EstimatedMinutes,
		IsActive:         true,
	}

	// Sort route stops by stop_order from frontend
	sort.Slice(req.RouteStops, func(i, j int) bool {
		return req.RouteStops[i].StopOrder < req.RouteStops[j].StopOrder
	})

	// Generate route stops with normalized order (100, 200, 300...)
	routeStops := make([]model.RouteStop, len(req.RouteStops))
	for i, stopReq := range req.RouteStops {
		routeStops[i] = model.RouteStop{
			StopOrder:     (i + 1) * 100, // Normalize to multiples of 100
			StopType:      stopReq.StopType,
			Location:      stopReq.Location,
			Address:       stopReq.Address,
			Latitude:      stopReq.Latitude,
			Longitude:     stopReq.Longitude,
			OffsetMinutes: stopReq.OffsetMinutes,
			IsActive:      true,
		}
	}
	route.RouteStops = routeStops

	// Create route with stops in single transaction
	if err := s.routeRepo.CreateRoute(ctx, route); err != nil {
		log.Error().Err(err).Msg("Failed to create route")
		return nil, ginext.NewInternalServerError("failed to create route")
	}

	return route, nil
}

func (s *RouteServiceImpl) UpdateRoute(ctx context.Context, id uuid.UUID, req *model.UpdateRouteRequest) (*model.Route, error) {
	route, err := s.routeRepo.GetRouteByID(ctx, id)
	if err != nil {
		return nil, ginext.NewInternalServerError("failed to get route")
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
			return nil, ginext.NewBadRequestError("distance must be positive")
		}
		route.DistanceKm = *req.DistanceKm
	}

	if req.EstimatedMinutes != nil {
		if *req.EstimatedMinutes <= 0 {
			return nil, ginext.NewBadRequestError("estimated time must be positive")
		}
		route.EstimatedMinutes = *req.EstimatedMinutes
	}
	if req.IsActive != nil {
		route.IsActive = *req.IsActive
	}

	if err := s.routeRepo.UpdateRoute(ctx, route); err != nil {
		log.Error().Err(err).Str("route_id", id.String()).Msg("Failed to update route")
		return nil, ginext.NewInternalServerError("failed to update route")
	}

	return route, nil
}

func (s *RouteServiceImpl) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	if err := s.routeRepo.DeleteRoute(ctx, id); err != nil {
		log.Error().Err(err).Str("route_id", id.String()).Msg("Failed to delete route")
		return ginext.NewInternalServerError("failed to delete route")
	}
	return nil
}
