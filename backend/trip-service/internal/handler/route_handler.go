package handler

import (
	"net/http"
	"strconv"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RouteHandler struct {
	routeService service.RouteService
}

func NewRouteHandler(routeService service.RouteService) *RouteHandler {
	return &RouteHandler{
		routeService: routeService,
	}
}

// CreateRoute handles route creation requests
func (h *RouteHandler) CreateRoute(c *gin.Context) {
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

// GetRoute handles get route by ID requests
func (h *RouteHandler) GetRoute(c *gin.Context) {
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

// UpdateRoute handles route update requests
func (h *RouteHandler) UpdateRoute(c *gin.Context) {
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

// DeleteRoute handles route deletion requests
func (h *RouteHandler) DeleteRoute(c *gin.Context) {
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

// ListRoutes handles listing routes with pagination
func (h *RouteHandler) ListRoutes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

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

// SearchRoutes handles route search by origin and destination
func (h *RouteHandler) SearchRoutes(c *gin.Context) {
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
