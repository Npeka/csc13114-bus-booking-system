package handler

import (
	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BusHandler interface {
	Get(r *ginext.Request) (*ginext.Response, error)
	GetList(r *ginext.Request) (*ginext.Response, error)

	Create(r *ginext.Request) (*ginext.Response, error)
	Update(r *ginext.Request) (*ginext.Response, error)
	Delete(r *ginext.Request) (*ginext.Response, error)
}

type BusHandlerImpl struct {
	service service.BusService
}

func NewBusHandler(service service.BusService) BusHandler {
	return &BusHandlerImpl{
		service: service,
	}
}

// Get godoc
// @Summary Get bus by ID
// @Description Get detailed information about a specific bus including seats
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.BusResponse} "Bus details"
// @Failure 400 {object} ginext.Response "Invalid bus ID"
// @Failure 404 {object} ginext.Response "Bus not found"
// @Router /api/v1/buses/{id} [get]
func (h *BusHandlerImpl) Get(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid bus ID")
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	bus, err := h.service.GetBusByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to get bus")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToBusResponse(bus)), nil
}

// GetList godoc
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
func (h *BusHandlerImpl) GetList(r *ginext.Request) (*ginext.Response, error) {
	var req model.ListBusesRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("Query binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	buses, total, err := h.service.ListBuses(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list buses")
		return nil, err
	}

	return ginext.NewPaginatedResponse(model.ToBusResponseList(buses), req.Page, req.PageSize, total), nil
}

// Create godoc
// @Summary Create a new bus
// @Description Create a new bus with model, plate number, seat configuration, and amenities (admin only)
// @Tags buses
// @Accept json
// @Produce json
// @Param request body model.CreateBusRequest true "Bus creation data"
// @Success 201 {object} ginext.Response{data=model.BusResponse} "Created bus"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses [post]
func (h *BusHandlerImpl) Create(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateBusRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	bus, err := h.service.CreateBus(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create bus")
		return nil, err
	}

	return ginext.NewCreatedResponse(model.ToBusResponse(bus)), nil
}

// Update godoc
// @Summary Update bus information
// @Description Update bus details such as model, plate number, bus type, amenities, or active status (admin only)
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Param request body model.UpdateBusRequest true "Bus update data"
// @Success 200 {object} ginext.Response{data=model.BusResponse} "Updated bus"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 404 {object} ginext.Response "Bus not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/{id} [put]
func (h *BusHandlerImpl) Update(r *ginext.Request) (*ginext.Response, error) {
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

	bus, err := h.service.UpdateBus(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to update bus")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToBusResponse(bus)), nil
}

// Delete godoc
// @Summary Delete a bus
// @Description Soft delete a bus by ID (admin only)
// @Tags buses
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {object} ginext.Response{data=string} "Success message"
// @Failure 400 {object} ginext.Response "Invalid bus ID"
// @Failure 404 {object} ginext.Response "Bus not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/{id} [delete]
func (h *BusHandlerImpl) Delete(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid bus ID")
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	if err = h.service.DeleteBus(r.Context(), id); err != nil {
		log.Error().Err(err).Str("bus_id", idStr).Msg("Failed to delete bus")
		return nil, err
	}

	return ginext.NewSuccessResponse("Bus deleted successfully"), nil
}
