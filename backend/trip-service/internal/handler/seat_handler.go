package handler

import (
	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatHandler interface {
	CreateSeat(r *ginext.Request) (*ginext.Response, error)
	CreateSeatsFromTemplate(r *ginext.Request) (*ginext.Response, error)
	UpdateSeat(r *ginext.Request) (*ginext.Response, error)
	DeleteSeat(r *ginext.Request) (*ginext.Response, error)
	GetSeatMap(r *ginext.Request) (*ginext.Response, error)
}

type SeatHandlerImpl struct {
	seatService service.SeatService
}

func NewSeatHandler(seatService service.SeatService) SeatHandler {
	return &SeatHandlerImpl{seatService: seatService}
}

// CreateSeat godoc
// @Summary Create seat
// @Description Add a new seat to a bus
// @Tags seats
// @Accept json
// @Produce json
// @Param request body model.CreateSeatRequest true "Seat data"
// @Success 201 {object} ginext.Response{data=model.Seat} "Created seat"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/seats [post]
func (h *SeatHandlerImpl) CreateSeat(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateSeatRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	seat, err := h.seatService.CreateSeat(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create seat")
		return nil, err
	}

	return ginext.NewCreatedResponse(seat), nil
}

// CreateSeatsFromTemplate godoc
// @Summary Bulk create seats
// @Description Create multiple seats for a bus from a template
// @Tags seats
// @Accept json
// @Produce json
// @Param request body model.BulkCreateSeatsRequest true "Bulk seat data"
// @Success 201 {object} ginext.Response{data=[]model.Seat} "Created seats"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/seats/bulk [post]
func (h *SeatHandlerImpl) CreateSeatsFromTemplate(r *ginext.Request) (*ginext.Response, error) {
	var req model.BulkCreateSeatsRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	seats, err := h.seatService.CreateSeatsFromTemplate(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create seats")
		return nil, err
	}

	return ginext.NewCreatedResponse(seats), nil
}

// UpdateSeat godoc
// @Summary Update seat
// @Description Update an existing seat
// @Tags seats
// @Accept json
// @Produce json
// @Param id path string true "Seat ID" format(uuid)
// @Param request body model.UpdateSeatRequest true "Update data"
// @Success 200 {object} ginext.Response{data=model.Seat} "Updated seat"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 404 {object} ginext.Response "Seat not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/seats/{id} [put]
func (h *SeatHandlerImpl) UpdateSeat(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid seat ID")
	}

	var req model.UpdateSeatRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	seat, err := h.seatService.UpdateSeat(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("seat_id", idStr).Msg("Failed to update seat")
		return nil, err
	}

	return ginext.NewSuccessResponse(seat), nil
}

// DeleteSeat godoc
// @Summary Delete seat
// @Description Remove a seat from a bus
// @Tags seats
// @Accept json
// @Produce json
// @Param id path string true "Seat ID" format(uuid)
// @Success 200 {object} ginext.Response "Success message"
// @Failure 400 {object} ginext.Response "Invalid seat ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/seats/{id} [delete]
func (h *SeatHandlerImpl) DeleteSeat(r *ginext.Request) (*ginext.Response, error) {
	id, err := uuid.Parse(r.GinCtx.Param("id"))
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid seat ID")
	}

	if err := h.seatService.DeleteSeat(r.Context(), id); err != nil {
		log.Error().Err(err).Msg("Failed to delete seat")
		return nil, err
	}

	return ginext.NewSuccessResponse("Seat deleted successfully"), nil
}

// GetSeatMap godoc
// @Summary Get seat map
// @Description Get the complete seat map for a bus
// @Tags seats
// @Accept json
// @Produce json
// @Param id path string true "Bus ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.SeatMapResponse} "Seat map"
// @Failure 400 {object} ginext.Response "Invalid bus ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/{id}/seat-map [get]
func (h *SeatHandlerImpl) GetSeatMap(r *ginext.Request) (*ginext.Response, error) {
	busID, err := uuid.Parse(r.GinCtx.Param("id"))
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid bus ID")
	}

	seatMap, err := h.seatService.GetSeatMap(r.Context(), busID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get seat map")
		return nil, err
	}

	return ginext.NewSuccessResponse(seatMap), nil
}
