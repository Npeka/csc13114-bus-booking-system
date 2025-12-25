package handler

import (
	"bus-booking/shared/ginext"
	"bus-booking/shared/utils"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatHandler interface {
	GetListByIDs(r *ginext.Request) (*ginext.Response, error)
	Update(r *ginext.Request) (*ginext.Response, error)
}

type SeatHandlerImpl struct {
	service service.SeatService
}

func NewSeatHandler(service service.SeatService) SeatHandler {
	return &SeatHandlerImpl{service: service}
}

// GetListByIDs godoc
// @Summary Get seats by IDs
// @Description Get multiple seats by their IDs (internal use)
// @Tags seats
// @Accept json
// @Produce json
// @Param seat_ids query []string true "Seat IDs" collectionFormat(multi)
// @Success 200 {object} ginext.Response{data=[]model.Seat} "List of seats"
// @Failure 400 {object} ginext.Response "Invalid seat IDs"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/seats/ids [get]
func (h *SeatHandlerImpl) GetListByIDs(r *ginext.Request) (*ginext.Response, error) {
	var req model.ListSeatsByIDsRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid query parameters")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	seatIDs, err := utils.ParseUUIDs(req.SeatIDs)
	if err != nil {
		log.Debug().Err(err).Msg("Invalid seat IDs")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	seats, err := h.service.GetListByIDs(r.Context(), seatIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list seats")
		return nil, err
	}

	return ginext.NewSuccessResponse(seats), nil
}

// Update godoc
// @Summary Update seat availability
// @Description Update seat availability status (admin only)
// @Tags seats
// @Accept json
// @Produce json
// @Param id path string true "Seat ID" format(uuid)
// @Param request body model.UpdateSeatRequest true "Update data"
// @Success 200 {object} ginext.Response{data=model.SeatResponse} "Updated seat"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 404 {object} ginext.Response "Seat not found"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/buses/seats/{id} [put]
func (h *SeatHandlerImpl) Update(r *ginext.Request) (*ginext.Response, error) {
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

	seat, err := h.service.Update(r.Context(), &req, id)
	if err != nil {
		log.Error().Err(err).Str("seat_id", idStr).Msg("Failed to update seat")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToSeatResponse(seat)), nil
}
