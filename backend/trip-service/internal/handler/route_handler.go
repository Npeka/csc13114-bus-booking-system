package handler

import (
	"net/http"
	"strconv"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RouteHandler interface {
	CreateRoute(c *gin.Context)
	GetRoute(c *gin.Context)
	UpdateRoute(c *gin.Context)
	DeleteRoute(c *gin.Context)
	ListRoutes(c *gin.Context)
	SearchRoutes(c *gin.Context)
}

type RouteHandlerImpl struct {
	routeService service.RouteService
}

func NewRouteHandler(routeService service.RouteService) RouteHandler {
	return &RouteHandlerImpl{
		routeService: routeService,
	}
}

// CreateRoute godoc
// @Summary Create a new route
// @Description Create a new route with origin, destination, and distance information
// @Tags routes
// @Accept json
// @Produce json
// @Param request body model.CreateRouteRequest true "Route creation data"
// @Success 201 {object} model.Route "Created route"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/routes [post]
func (h *RouteHandlerImpl) CreateRoute(c *gin.Context) {
	var req model.CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.routeService.CreateRoute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, route)
}

// GetRoute godoc
// @Summary Get route by ID
// @Description Get detailed information about a specific route
// @Tags routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID" format(uuid)
// @Success 200 {object} model.Route "Route details"
// @Failure 400 {object} map[string]string "Invalid route ID"
// @Failure 404 {object} map[string]string "Route not found"
// @Router /api/v1/routes/{id} [get]
func (h *RouteHandlerImpl) GetRoute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	route, err := h.routeService.GetRouteByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}

	c.JSON(http.StatusOK, route)
}

// UpdateRoute godoc
// @Summary Update route
// @Description Update route information such as origin, destination, or distance
// @Tags routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID" format(uuid)
// @Param request body model.UpdateRouteRequest true "Route update data"
// @Success 200 {object} model.Route "Updated route"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/routes/{id} [put]
func (h *RouteHandlerImpl) UpdateRoute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	var req model.UpdateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.routeService.UpdateRoute(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

// DeleteRoute godoc
// @Summary Delete route
// @Description Delete a route by ID
// @Tags routes
// @Accept json
// @Produce json
// @Param id path string true "Route ID" format(uuid)
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid route ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/routes/{id} [delete]
func (h *RouteHandlerImpl) DeleteRoute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	err = h.routeService.DeleteRoute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "route deleted successfully"})
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
// @Success 200 {object} map[string]interface{} "Paginated route list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/routes [get]
func (h *RouteHandlerImpl) ListRoutes(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit number"})
		return
	}

	var operatorID *uuid.UUID
	if operatorIDStr := c.Query("operator_id"); operatorIDStr != "" {
		if id, err := uuid.Parse(operatorIDStr); err == nil {
			operatorID = &id
		}
	}

	routes, total, err := h.routeService.ListRoutes(c.Request.Context(), operatorID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := gin.H{
		"routes":      routes,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// SearchRoutes godoc
// @Summary Search routes
// @Description Search routes by origin and destination
// @Tags routes
// @Accept json
// @Produce json
// @Param origin query string true "Origin city"
// @Param destination query string true "Destination city"
// @Success 200 {array} model.Route "List of matching routes"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/routes/search [get]
func (h *RouteHandlerImpl) SearchRoutes(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")

	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "origin and destination are required"})
		return
	}

	routes, err := h.routeService.GetRoutesByOriginDestination(c.Request.Context(), origin, destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routes)
}
