package handler

import (
	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RouteHandler interface {
	GetRoute(r *ginext.Request) (*ginext.Response, error)
	ListRoutes(r *ginext.Request) (*ginext.Response, error)
	SearchRoutes(r *ginext.Request) (*ginext.Response, error)

	CreateRoute(r *ginext.Request) (*ginext.Response, error)
	UpdateRoute(r *ginext.Request) (*ginext.Response, error)
	DeleteRoute(r *ginext.Request) (*ginext.Response, error)
}

type RouteHandlerImpl struct {
	service service.RouteService
}

func NewRouteHandler(service service.RouteService) RouteHandler {
	return &RouteHandlerImpl{
		service: service,
	}
}

// GetRoute godoc
// @Summary Get route by ID
// @Description Get detailed information about a specific route
// @Tags routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.Route} "Route details"
// @Failure 400 {object} ginext.Response "Invalid route ID"
// @Failure 404 {object} ginext.Response "Route not found"
// @Router /api/v1/routes/{id} [get]
func (h *RouteHandlerImpl) GetRoute(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid route ID")
	}

	route, err := h.service.GetRouteByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("route_id", idStr).Msg("Failed to get route")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToRouteResponse(route)), nil
}

// ListRoutes godoc
// @Summary List routes
// @Description Get a paginated list of routes, optionally filtered by operator
// @Tags routes
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param operator_id query string false "Filter by operator ID" format(uuid)
// @Success 200 {object} ginext.Response "Paginated route list"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes [get]
func (h *RouteHandlerImpl) ListRoutes(r *ginext.Request) (*ginext.Response, error) {
	var req model.ListRoutesRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	routes, total, err := h.service.ListRoutes(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list routes")
		return nil, err
	}

	return ginext.NewPaginatedResponse(model.ToRouteResponseList(routes), req.Page, req.PageSize, total), nil
}

// SearchRoutes godoc
// @Summary Search routes
// @Description Search routes by origin and destination
// @Tags routes
// @Accept json
// @Produce json
// @Param origin query string true "Origin city"
// @Param destination query string true "Destination city"
// @Success 200 {object} ginext.Response{data=[]model.Route} "List of matching routes"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/search [get]
func (h *RouteHandlerImpl) SearchRoutes(r *ginext.Request) (*ginext.Response, error) {
	var req model.SearchRoutesQueryRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	routes, err := h.service.GetRoutesByOriginDestination(r.Context(), req.Origin, req.Destination)
	if err != nil {
		log.Error().Err(err).Str("origin", req.Origin).Str("destination", req.Destination).Msg("Failed to search routes")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToRouteResponseList(routes)), nil
}

// CreateRoute godoc
// @Summary Create a new route
// @Description Create a new route with origin, destination, and distance information
// @Tags routes
// @Accept json
// @Produce json
// @Param request body model.CreateRouteRequest true "Route creation data"
// @Success 201 {object} ginext.Response{data=model.Route} "Created route"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes [post]
func (h *RouteHandlerImpl) CreateRoute(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateRouteRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	route, err := h.service.CreateRoute(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create route")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToRouteResponse(route)), nil
}

// UpdateRoute godoc
// @Summary Update route
// @Description Update route information such as origin, destination, or distance
// @Tags routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID" format(uuid)
// @Param request body model.UpdateRouteRequest true "Route update data"
// @Success 200 {object} ginext.Response{data=model.Route} "Updated route"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/{id} [put]
func (h *RouteHandlerImpl) UpdateRoute(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid route ID")
	}

	var req model.UpdateRouteRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	route, err := h.service.UpdateRoute(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("route_id", idStr).Msg("Failed to update route")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToRouteResponse(route)), nil
}

// DeleteRoute godoc
// @Summary Delete route
// @Description Delete a route by ID
// @Tags routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID" format(uuid)
// @Success 200 {object} ginext.Response "Success message"
// @Failure 400 {object} ginext.Response "Invalid route ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/{id} [delete]
func (h *RouteHandlerImpl) DeleteRoute(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid route ID")
	}

	err = h.service.DeleteRoute(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("route_id", idStr).Msg("Failed to delete route")
		return nil, err
	}

	return ginext.NewSuccessResponse("Route deleted successfully"), nil
}
