package handler

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatHandler interface {
	GetSeatAvailability(r *ginext.Request) (*ginext.Response, error)
	ReserveSeat(r *ginext.Request) (*ginext.Response, error)
	ReleaseSeat(r *ginext.Request) (*ginext.Response, error)
}

type SeatHandlerImpl struct {
	seatService service.SeatService
}

func NewSeatHandler(seatService service.SeatService) SeatHandler {
	return &SeatHandlerImpl{
		seatService: seatService,
	}
}

// GetSeatAvailability godoc
// @Summary Get seat availability for a trip
// @Description Get available, reserved, and booked seats for a specific trip
// @Tags seats
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.SeatAvailabilityResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/trips/{trip_id}/seats [get]
func (h *SeatHandlerImpl) GetSeatAvailability(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	availability, err := h.seatService.GetSeatAvailability(r.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("Failed to get seat availability")
		return nil, err
	}

	return ginext.NewSuccessResponse(availability), nil
}

// ReserveSeat godoc
// @Summary Reserve a seat
// @Description Reserve a seat for a user temporarily
// @Tags seats
// @Accept json
// @Produce json
// @Param request body model.ReserveSeatRequest true "Seat reservation request"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/seats/reserve [post]
func (h *SeatHandlerImpl) ReserveSeat(r *ginext.Request) (*ginext.Response, error) {
	var req model.ReserveSeatRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.seatService.ReserveSeat(r.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Failed to reserve seat")
		return nil, err
	}

	return ginext.NewSuccessResponse("Seat reserved successfully"), nil
}

// ReleaseSeat godoc
// @Summary Release a reserved seat
// @Description Release a reserved seat back to available status
// @Tags seats
// @Accept json
// @Produce json
// @Param request body model.ReleaseSeatRequest true "Seat release request"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/seats/release [post]
func (h *SeatHandlerImpl) ReleaseSeat(r *ginext.Request) (*ginext.Response, error) {
	var req model.ReleaseSeatRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.seatService.ReleaseSeat(r.Context(), req.TripID, req.SeatID); err != nil {
		log.Error().Err(err).Msg("Failed to release seat")
		return nil, err
	}

	return ginext.NewSuccessResponse("Seat released successfully"), nil
}
