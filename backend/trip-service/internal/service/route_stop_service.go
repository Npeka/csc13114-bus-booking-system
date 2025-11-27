package service

import (
	"context"
	"fmt"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RouteStopService interface {
	CreateRouteStop(ctx context.Context, req *model.CreateRouteStopRequest) (*model.RouteStop, error)
	UpdateRouteStop(ctx context.Context, id uuid.UUID, req *model.UpdateRouteStopRequest) (*model.RouteStop, error)
	DeleteRouteStop(ctx context.Context, id uuid.UUID) error
	ListRouteStops(ctx context.Context, routeID uuid.UUID) ([]model.RouteStop, error)
	ReorderStops(ctx context.Context, routeID uuid.UUID, stopOrders []StopOrder) error
}

type StopOrder struct {
	StopID uuid.UUID `json:"stop_id" validate:"required"`
	Order  int       `json:"order" validate:"required,min=1"`
}

type RouteStopServiceImpl struct {
	stopRepo  repository.RouteStopRepository
	routeRepo repository.RouteRepository
}

func NewRouteStopService(stopRepo repository.RouteStopRepository, routeRepo repository.RouteRepository) RouteStopService {
	return &RouteStopServiceImpl{
		stopRepo:  stopRepo,
		routeRepo: routeRepo,
	}
}

func (s *RouteStopServiceImpl) CreateRouteStop(ctx context.Context, req *model.CreateRouteStopRequest) (*model.RouteStop, error) {
	// Verify route exists
	_, err := s.routeRepo.GetRouteByID(ctx, req.RouteID)
	if err != nil {
		log.Error().Err(err).Str("route_id", req.RouteID.String()).Msg("Route not found")
		return nil, fmt.Errorf("route not found: %w", err)
	}

	stop := &model.RouteStop{
		RouteID:       req.RouteID,
		StopOrder:     req.StopOrder,
		StopType:      req.StopType,
		Location:      req.Location,
		Address:       req.Address,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		OffsetMinutes: req.OffsetMinutes,
		IsActive:      true,
	}
	if err := s.stopRepo.Create(ctx, stop); err != nil {
		log.Error().Err(err).Msg("Failed to create route stop")
		return nil, fmt.Errorf("failed to create route stop: %w", err)
	}

	log.Info().Str("stop_id", stop.ID.String()).Msg("Route stop created successfully")
	return stop, nil
}

func (s *RouteStopServiceImpl) UpdateRouteStop(ctx context.Context, id uuid.UUID, req *model.UpdateRouteStopRequest) (*model.RouteStop, error) {
	stop, err := s.stopRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("stop_id", id.String()).Msg("Route stop not found")
		return nil, fmt.Errorf("route stop not found: %w", err)
	}

	// Update fields if provided
	if req.StopOrder != nil {
		stop.StopOrder = *req.StopOrder
	}
	if req.StopType != nil {
		stop.StopType = *req.StopType
	}
	if req.Location != nil {
		stop.Location = *req.Location
	}
	if req.Address != nil {
		stop.Address = *req.Address
	}
	if req.Latitude != nil {
		stop.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		stop.Longitude = req.Longitude
	}
	if req.OffsetMinutes != nil {
		stop.OffsetMinutes = *req.OffsetMinutes
	}
	if req.IsActive != nil {
		stop.IsActive = *req.IsActive
	}

	if err := s.stopRepo.Update(ctx, stop); err != nil {
		log.Error().Err(err).Msg("Failed to update route stop")
		return nil, fmt.Errorf("failed to update route stop: %w", err)
	}

	log.Info().Str("stop_id", stop.ID.String()).Msg("Route stop updated successfully")
	return stop, nil
}

func (s *RouteStopServiceImpl) DeleteRouteStop(ctx context.Context, id uuid.UUID) error {
	if err := s.stopRepo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Str("stop_id", id.String()).Msg("Failed to delete route stop")
		return fmt.Errorf("failed to delete route stop: %w", err)
	}

	log.Info().Str("stop_id", id.String()).Msg("Route stop deleted successfully")
	return nil
}

func (s *RouteStopServiceImpl) ListRouteStops(ctx context.Context, routeID uuid.UUID) ([]model.RouteStop, error) {
	stops, err := s.stopRepo.ListByRouteID(ctx, routeID)
	if err != nil {
		log.Error().Err(err).Str("route_id", routeID.String()).Msg("Failed to list route stops")
		return nil, fmt.Errorf("failed to list route stops: %w", err)
	}

	return stops, nil
}

func (s *RouteStopServiceImpl) ReorderStops(ctx context.Context, routeID uuid.UUID, stopOrders []StopOrder) error {
	// Convert to map for repository
	orderMap := make(map[uuid.UUID]int)
	for _, so := range stopOrders {
		orderMap[so.StopID] = so.Order
	}

	if err := s.stopRepo.ReorderStops(ctx, routeID, orderMap); err != nil {
		log.Error().Err(err).Str("route_id", routeID.String()).Msg("Failed to reorder stops")
		return fmt.Errorf("failed to reorder stops: %w", err)
	}

	log.Info().Str("route_id", routeID.String()).Msg("Route stops reordered successfully")
	return nil
}
