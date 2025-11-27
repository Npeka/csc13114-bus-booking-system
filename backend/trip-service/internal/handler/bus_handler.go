package handler

import (
	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BusHandler interface {
	CreateBus(r *ginext.Request) (*ginext.Response, error)
	GetBus(r *ginext.Request) (*ginext.Response, error)
	UpdateBus(r *ginext.Request) (*ginext.Response, error)
	DeleteBus(r *ginext.Request) (*ginext.Response, error)
	ListBuses(r *ginext.Request) (*ginext.Response, error)
	GetBusSeats(r *ginext.Request) (*ginext.Response, error)
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
// @Success 201 {object} ginext.Response{data=model.Bus} "Created bus"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses [post]
func (h *BusHandlerImpl) CreateBus(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateBusRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	bus, err := h.busService.CreateBus(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create bus")
		return nil, err
	}

	return ginext.NewSuccessResponse(bus, "Bus created successfully"), nil
}

// GetBus godoc
// @Summary Get bus by ID
// @Description Get detailed information about a specific bus
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.Bus} "Bus details"
// @Failure 400 {object} ginext.Response "Invalid bus ID"
// @Failure 404 {object} ginext.Response "Bus not found"
// @Router /api/v1/buses/{id} [get]
func (h *BusHandlerImpl) GetBus(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	bus, err := h.busService.GetBusByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to get bus")
		return nil, err
	}

	return ginext.NewSuccessResponse(bus, "Bus retrieved successfully"), nil
}

// UpdateBus godoc
// @Summary Update bus
// @Description Update bus information such as model, plate number, or amenities
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Param request body model.UpdateBusRequest true "Bus update data"
// @Success 200 {object} ginext.Response{data=model.Bus} "Updated bus"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/{id} [put]
func (h *BusHandlerImpl) UpdateBus(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	var req model.UpdateBusRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	bus, err := h.busService.UpdateBus(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to update bus")
		return nil, err
	}

	return ginext.NewSuccessResponse(bus, "Bus updated successfully"), nil
}

// DeleteBus godoc
// @Summary Delete bus
// @Description Delete a bus by ID
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {object} ginext.Response "Success message"
// @Failure 400 {object} ginext.Response "Invalid bus ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/{id} [delete]
func (h *BusHandlerImpl) DeleteBus(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	err = h.busService.DeleteBus(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to delete bus")
		return nil, err
	}

	return ginext.NewSuccessResponse(nil, "Bus deleted successfully"), nil
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
// @Success 200 {object} ginext.Response "Paginated bus list"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses [get]
func (h *BusHandlerImpl) ListBuses(r *ginext.Request) (*ginext.Response, error) {
	var req model.ListBusesRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	var operatorID *uuid.UUID
	if req.OperatorID != "" {
		id, err := uuid.Parse(req.OperatorID)
		if err != nil {
			return nil, ginext.NewBadRequestError("invalid operator ID")
		}
		operatorID = &id
	}

	buses, total, err := h.busService.ListBuses(r.Context(), operatorID, req.Page, req.Limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list buses")
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	response := gin.H{
		"buses":       buses,
		"total":       total,
		"page":        req.Page,
		"limit":       req.Limit,
		"total_pages": totalPages,
	}

	return ginext.NewSuccessResponse(response, "Buses retrieved successfully"), nil
}

// GetBusSeats godoc
// @Summary Get bus seats
// @Description Get all seats for a specific bus
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {object} ginext.Response{data=[]model.Seat} "List of bus seats"
// @Failure 400 {object} ginext.Response "Invalid bus ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/{id}/seats [get]
func (h *BusHandlerImpl) GetBusSeats(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	seats, err := h.busService.GetBusSeats(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to get bus seats")
		return nil, err
	}

	return ginext.NewSuccessResponse(seats, "Bus seats retrieved successfully"), nil
}
