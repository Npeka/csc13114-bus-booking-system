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

// CreateBus godoc
// @Summary Create a new bus
// @Description Create a new bus with operator, model, and seat capacity information
// @Tags buses
// @Accept json
// @Produce json
// @Param request body model.CreateBusRequest true "Bus creation data"
// @Success 201 {object} model.Bus "Created bus"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/buses [post]
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

// GetBus godoc
// @Summary Get bus by ID
// @Description Get detailed information about a specific bus
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {object} model.Bus "Bus details"
// @Failure 400 {object} map[string]string "Invalid bus ID"
// @Failure 404 {object} map[string]string "Bus not found"
// @Router /api/v1/buses/{id} [get]
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

// UpdateBus godoc
// @Summary Update bus
// @Description Update bus information such as model, plate number, or amenities
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Param request body model.UpdateBusRequest true "Bus update data"
// @Success 200 {object} model.Bus "Updated bus"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/buses/{id} [put]
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

// DeleteBus godoc
// @Summary Delete bus
// @Description Delete a bus by ID
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid bus ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/buses/{id} [delete]
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

// ListBuses godoc
// @Summary List buses
// @Description Get a paginated list of buses, optionally filtered by operator
// @Tags buses
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param operator_id query string false "Filter by operator ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Paginated bus list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/buses [get]
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

// GetBusSeats godoc
// @Summary Get bus seats
// @Description Get all seats for a specific bus
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {array} model.Seat "List of bus seats"
// @Failure 400 {object} map[string]string "Invalid bus ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/buses/{id}/seats [get]
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
