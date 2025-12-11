package handler

import (
	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RouteStopHandler interface {
	CreateRouteStop(r *ginext.Request) (*ginext.Response, error)
	UpdateRouteStop(r *ginext.Request) (*ginext.Response, error)
	MoveRouteStop(r *ginext.Request) (*ginext.Response, error)
	DeleteRouteStop(r *ginext.Request) (*ginext.Response, error)
	ListRouteStops(r *ginext.Request) (*ginext.Response, error)
	ReorderStops(r *ginext.Request) (*ginext.Response, error)
}

type RouteStopHandlerImpl struct {
	stopService service.RouteStopService
}

func NewRouteStopHandler(stopService service.RouteStopService) RouteStopHandler {
	return &RouteStopHandlerImpl{stopService: stopService}
}

// CreateRouteStop godoc
// @Summary Create route stop
// @Description Add a new pickup/dropoff point to a route
// @Tags route-stops
// @Accept json
// @Produce json
// @Param request body model.CreateRouteStopRequest true "Route stop data"
// @Success 201 {object} ginext.Response{data=model.RouteStop} "Created route stop"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/stops [post]
func (h *RouteStopHandlerImpl) CreateRouteStop(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateRouteStopRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	stop, err := h.stopService.CreateRouteStop(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create route stop")
		return nil, err
	}

	return ginext.NewCreatedResponse(model.ToRouteStopResponse(stop)), nil
}

// UpdateRouteStop godoc
// @Summary Update route stop
// @Description Update an existing route stop (does NOT reorder - use move endpoint)
// @Tags route-stops
// @Accept json
// @Produce json
// @Param id path string true "Stop ID" format(uuid)
// @Param request body model.UpdateRouteStopRequest true "Update data"
// @Success 200 {object} ginext.Response{data=model.RouteStop} "Updated route stop"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 404 {object} ginext.Response "Stop not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/stops/{id} [put]
func (h *RouteStopHandlerImpl) UpdateRouteStop(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid stop ID")
	}

	var req model.UpdateRouteStopRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	// Remove stop_order from update - use move endpoint instead
	req.StopOrder = nil

	stop, err := h.stopService.UpdateRouteStop(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("stop_id", idStr).Msg("Failed to update route stop")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToRouteStopResponse(stop)), nil
}

// MoveRouteStop godoc
// @Summary Move route stop
// @Description Move a stop to before/after another stop or to first/last position
// @Tags route-stops
// @Accept json
// @Produce json
// @Param id path string true "Stop ID to move" format(uuid)
// @Param request body model.MoveRouteStopRequest true "Move position data"
// @Success 200 {object} ginext.Response{data=model.RouteStop} "Moved route stop"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 404 {object} ginext.Response "Stop not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/stops/{id}/move [post]
func (h *RouteStopHandlerImpl) MoveRouteStop(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid stop ID")
	}

	var req model.MoveRouteStopRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	// Validate reference_stop_id for before/after positions
	if (req.Position == "before" || req.Position == "after") && req.ReferenceStopID == nil {
		return nil, ginext.NewBadRequestError("reference_stop_id is required for before/after positions")
	}

	stop, err := h.stopService.MoveRouteStop(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("stop_id", idStr).Msg("Failed to move route stop")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToRouteStopResponse(stop)), nil
}

// DeleteRouteStop godoc
// @Summary Delete route stop
// @Description Remove a route stop
// @Tags route-stops
// @Accept json
// @Produce json
// @Param id path string true "Stop ID" format(uuid)
// @Success 200 {object} ginext.Response "Success message"
// @Failure 400 {object} ginext.Response "Invalid stop ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/stops/{id} [delete]
func (h *RouteStopHandlerImpl) DeleteRouteStop(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid stop ID")
	}

	if err := h.stopService.DeleteRouteStop(r.Context(), id); err != nil {
		log.Error().Err(err).Msg("Failed to delete route stop")
		return nil, err
	}

	return ginext.NewSuccessResponse("Route stop deleted successfully"), nil
}

// ListRouteStops godoc
// @Summary List route stops
// @Description Get all stops for a specific route
// @Tags route-stops
// @Accept json
// @Produce json
// @Param id path string true "Route ID" format(uuid)
// @Success 200 {object} ginext.Response{data=[]model.RouteStop} "List of route stops"
// @Failure 400 {object} ginext.Response "Invalid route ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/{id}/stops [get]
func (h *RouteStopHandlerImpl) ListRouteStops(r *ginext.Request) (*ginext.Response, error) {
	routeIDStr := r.GinCtx.Param("id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid route ID")
	}

	stops, err := h.stopService.ListRouteStops(r.Context(), routeID)
	if err != nil {
		log.Error().Err(err).Str("route_id", routeIDStr).Msg("Failed to list route stops")
		return nil, err
	}

	// Convert to response with mapped constants
	responses := make([]model.RouteStopResponse, len(stops))
	for i, stop := range stops {
		responses[i] = *model.ToRouteStopResponse(&stop)
	}

	return ginext.NewSuccessResponse(responses), nil
}

// ReorderStops godoc
// @Summary Reorder route stops
// @Description Update the order of stops for a route
// @Tags route-stops
// @Accept json
// @Produce json
// @Param id path string true "Route ID" format(uuid)
// @Param request body []service.StopOrder true "Stop order data"
// @Success 200 {object} ginext.Response "Success message"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/{id}/stops/reorder [post]
func (h *RouteStopHandlerImpl) ReorderStops(r *ginext.Request) (*ginext.Response, error) {
	routeIDStr := r.GinCtx.Param("id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid route ID")
	}

	var stopOrders []service.StopOrder
	if err := r.GinCtx.ShouldBindJSON(&stopOrders); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.stopService.ReorderStops(r.Context(), routeID, stopOrders); err != nil {
		log.Error().Err(err).Msg("Failed to reorder stops")
		return nil, err
	}

	return ginext.NewNoContentResponse(), nil
}
