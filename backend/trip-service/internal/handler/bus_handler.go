package handler

import (
	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BusHandler interface {
	GetBus(r *ginext.Request) (*ginext.Response, error)
	ListBuses(r *ginext.Request) (*ginext.Response, error)
	GetBusSeats(r *ginext.Request) (*ginext.Response, error)

	CreateBus(r *ginext.Request) (*ginext.Response, error)
	UpdateBus(r *ginext.Request) (*ginext.Response, error)
	DeleteBus(r *ginext.Request) (*ginext.Response, error)
}

type BusHandlerImpl struct {
	busService service.BusService
}

func NewBusHandler(busService service.BusService) BusHandler {
	return &BusHandlerImpl{
		busService: busService,
	}
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
		log.Error().Err(err).Msg("Invalid bus ID")
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	bus, err := h.busService.GetBusByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to get bus")
		return nil, err
	}

	return ginext.NewSuccessResponse(bus), nil
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
		log.Error().Err(err).Msg("Query binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	buses, total, err := h.busService.ListBuses(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list buses")
		return nil, err
	}

	return ginext.NewPaginatedResponse(buses, req.Page, req.PageSize, total), nil
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
		log.Error().Err(err).Msg("Invalid bus ID")
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	seats, err := h.busService.GetBusSeats(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to get bus seats")
		return nil, err
	}

	return ginext.NewSuccessResponse(seats), nil
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
		log.Error().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	bus, err := h.busService.CreateBus(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create bus")
		return nil, err
	}

	return ginext.NewCreatedResponse(bus), nil
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
		log.Error().Err(err).Msg("Invalid bus ID")
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	var req model.UpdateBusRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	bus, err := h.busService.UpdateBus(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to update bus")
		return nil, err
	}

	return ginext.NewSuccessResponse(bus), nil
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
		log.Error().Err(err).Msg("Invalid bus ID")
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	if err = h.busService.DeleteBus(r.Context(), id); err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to delete bus")
		return nil, err
	}

	return ginext.NewSuccessResponse("Bus deleted successfully"), nil
}
