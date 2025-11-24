package handler

import (
	"net/http"
	"strconv"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BusHandler interface {
	CreateBus(c *gin.Context)
	GetBus(c *gin.Context)
	UpdateBus(c *gin.Context)
	DeleteBus(c *gin.Context)
	ListBuses(c *gin.Context)
	GetBusSeats(c *gin.Context)
}

type BusHandlerImpl struct {
	busService service.BusService
}

func NewBusHandler(busService service.BusService) BusHandler {
	return &BusHandlerImpl{
		busService: busService,
	}
}

// CreateBus handles bus creation requests
func (h *BusHandlerImpl) CreateBus(c *gin.Context) {
	var req model.CreateBusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bus, err := h.busService.CreateBus(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bus)
}

// GetBus handles get bus by ID requests
func (h *BusHandlerImpl) GetBus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bus ID"})
		return
	}

	bus, err := h.busService.GetBusByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bus not found"})
		return
	}

	c.JSON(http.StatusOK, bus)
}

// UpdateBus handles bus update requests
func (h *BusHandlerImpl) UpdateBus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bus ID"})
		return
	}

	var req model.UpdateBusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bus, err := h.busService.UpdateBus(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bus)
}

// DeleteBus handles bus deletion requests
func (h *BusHandlerImpl) DeleteBus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bus ID"})
		return
	}

	err = h.busService.DeleteBus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bus deleted successfully"})
}

// ListBuses handles listing buses with pagination
func (h *BusHandlerImpl) ListBuses(c *gin.Context) {
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

	buses, total, err := h.busService.ListBuses(c.Request.Context(), operatorID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := gin.H{
		"buses":       buses,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetBusSeats handles getting seats for a specific bus
func (h *BusHandlerImpl) GetBusSeats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bus ID"})
		return
	}

	seats, err := h.busService.GetBusSeats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}
