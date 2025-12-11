package service

import (
	"context"
	"fmt"

	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RouteStopService interface {
	CreateRouteStop(ctx context.Context, req *model.CreateRouteStopRequest) (*model.RouteStop, error)
	UpdateRouteStop(ctx context.Context, id uuid.UUID, req *model.UpdateRouteStopRequest) (*model.RouteStop, error)
	MoveRouteStop(ctx context.Context, id uuid.UUID, req *model.MoveRouteStopRequest) (*model.RouteStop, error)
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
	// Verify route exists and get stops sorted by order ASC
	route, err := s.routeRepo.GetRoutesWithRouteStops(ctx, req.RouteID)
	if err != nil {
		log.Error().Err(err).Str("route_id", req.RouteID.String()).Msg("Route not found")
		return nil, ginext.NewBadRequestError("route not found")
	}

	existingStops := route.RouteStops
	newStopOrder := req.StopOrder

	// Handle default and append cases
	if len(existingStops) == 0 {
		// First stop defaults to 100
		if newStopOrder <= 0 {
			newStopOrder = 100
		}
	} else {
		// Stops are sorted ASC, last stop has max order
		maxOrder := existingStops[len(existingStops)-1].StopOrder

		if newStopOrder <= 0 || newStopOrder > maxOrder {
			// Append at end with spacing of 100
			newStopOrder = maxOrder + 100
			log.Info().Int("requested_order", req.StopOrder).Int("adjusted_order", newStopOrder).Msg("Adjusted stop order to end of list")
		}
	}

	// Shift stops with order >= newStopOrder by +1
	orderMap := make(map[uuid.UUID]int)
	for _, existingStop := range existingStops {
		if existingStop.StopOrder >= newStopOrder {
			orderMap[existingStop.ID] = existingStop.StopOrder + 1
		}
	}

	// Apply shift if needed
	if len(orderMap) > 0 {
		if err := s.stopRepo.ReorderStops(ctx, req.RouteID, orderMap); err != nil {
			log.Error().Err(err).Msg("Failed to shift existing stops")
			return nil, ginext.NewInternalServerError("failed to shift existing stops")
		}
		log.Info().Int("shifted_count", len(orderMap)).Int("new_order", newStopOrder).Msg("Shifted existing stops")
	}

	// Create the new stop
	stop := &model.RouteStop{
		RouteID:       req.RouteID,
		StopOrder:     newStopOrder,
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
		return nil, ginext.NewInternalServerError("failed to create route stop")
	}

	log.Info().
		Str("route_id", req.RouteID.String()).
		Int("stop_order", newStopOrder).
		Str("location", req.Location).
		Msg("Route stop created successfully")

	return stop, nil
}

func (s *RouteStopServiceImpl) UpdateRouteStop(ctx context.Context, id uuid.UUID, req *model.UpdateRouteStopRequest) (*model.RouteStop, error) {
	stop, err := s.stopRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("stop_id", id.String()).Msg("Route stop not found")
		return nil, ginext.NewBadRequestError("route stop not found")
	}

	// Update fields (stop_order is ignored - use MoveRouteStop instead)
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
		return nil, ginext.NewInternalServerError("failed to update route stop")
	}

	log.Info().
		Str("stop_id", id.String()).
		Int("stop_order", stop.StopOrder).
		Msg("Route stop updated successfully")

	return stop, nil
}

func (s *RouteStopServiceImpl) MoveRouteStop(ctx context.Context, id uuid.UUID, req *model.MoveRouteStopRequest) (*model.RouteStop, error) {
	// Get the stop to move
	stop, err := s.stopRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("stop_id", id.String()).Msg("Route stop not found")
		return nil, ginext.NewBadRequestError("route stop not found")
	}

	// Get all stops for the route sorted by order
	route, err := s.routeRepo.GetRoutesWithRouteStops(ctx, stop.RouteID)
	if err != nil {
		log.Error().Err(err).Str("route_id", stop.RouteID.String()).Msg("Route not found")
		return nil, ginext.NewBadRequestError("route not found")
	}

	allStops := route.RouteStops
	if len(allStops) <= 1 {
		// Only one stop, nothing to reorder
		return stop, nil
	}

	var newOrder int

	switch req.Position {
	case "first":
		// Move to first position (before current first stop)
		newOrder = allStops[0].StopOrder - 1

	case "last":
		// Move to last position (after current last stop)
		newOrder = allStops[len(allStops)-1].StopOrder + 1

	case "before":
		if req.ReferenceStopID == nil {
			return nil, ginext.NewBadRequestError("reference_stop_id required for 'before' position")
		}
		// Find reference stop
		var refIndex int = -1
		for i, s := range allStops {
			if s.ID == *req.ReferenceStopID {
				refIndex = i
				break
			}
		}
		if refIndex == -1 {
			return nil, ginext.NewBadRequestError("reference stop not found")
		}
		// Place before reference (order - 1)
		newOrder = allStops[refIndex].StopOrder - 1

	case "after":
		if req.ReferenceStopID == nil {
			return nil, ginext.NewBadRequestError("reference_stop_id required for 'after' position")
		}
		// Find reference stop
		var refIndex int = -1
		for i, s := range allStops {
			if s.ID == *req.ReferenceStopID {
				refIndex = i
				break
			}
		}
		if refIndex == -1 {
			return nil, ginext.NewBadRequestError("reference stop not found")
		}
		// Place after reference (order + 1)
		newOrder = allStops[refIndex].StopOrder + 1

	default:
		return nil, ginext.NewBadRequestError("invalid position")
	}

	// Update only the moved stop's order
	orderMap := map[uuid.UUID]int{
		id: newOrder,
	}

	if err := s.stopRepo.ReorderStops(ctx, stop.RouteID, orderMap); err != nil {
		log.Error().Err(err).Msg("Failed to move route stop")
		return nil, ginext.NewInternalServerError("failed to move route stop")
	}

	stop.StopOrder = newOrder

	log.Info().
		Str("stop_id", id.String()).
		Int("new_order", newOrder).
		Str("position", req.Position).
		Msg("Route stop moved successfully")

	return stop, nil
}

func (s *RouteStopServiceImpl) DeleteRouteStop(ctx context.Context, id uuid.UUID) error {
	if err := s.stopRepo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Str("stop_id", id.String()).Msg("Failed to delete route stop")
		return ginext.NewInternalServerError("failed to delete route stop")
	}
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
