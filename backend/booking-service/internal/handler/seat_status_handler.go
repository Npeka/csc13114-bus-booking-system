package handler

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatStatusHandler interface {
	InitSeatsForTrip(r *ginext.Request) (*ginext.Response, error)
	GetSeatAvailability(r *ginext.Request) (*ginext.Response, error)
}

type SeatStatusHandlerImpl struct {
	seatStatusService service.SeatStatusService
}

func NewSeatStatusHandler(seatStatusService service.SeatStatusService) SeatStatusHandler {
	return &SeatStatusHandlerImpl{
		seatStatusService: seatStatusService,
	}
}

// InitSeatsForTrip godoc
// @Summary Initialize seat statuses for a trip
// @Description Called by Trip Service when a new trip is created
// @Tags seat-status
// @Accept json
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Param request body model.InitSeatsRequest true "Seat initialization data"
// @Success 200 {object} ginext.Response
// @Failure 400 {object} ginext.Response
// @Router /api/v1/trips/{trip_id}/seats/init [post]
func (h *SeatStatusHandlerImpl) InitSeatsForTrip(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("invalid trip id")
		return nil, ginext.NewBadRequestError("invalid trip id")
	}

	var req model.InitSeatsRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.seatStatusService.InitSeatsForTrip(r.Context(), tripID, req.Seats); err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("failed to initialize seats")
		return nil, err
	}

	return ginext.NewSuccessResponse("seats initialized successfully"), nil
}

// GetSeatAvailability godoc
// @Summary Get seat availability for a trip
// @Description Get all seats and their statuses for a specific trip
// @Tags seat-status
// @Accept json
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Success 200 {object} ginext.Response{data=[]model.SeatStatus}
// @Failure 400 {object} ginext.Response
// @Router /api/v1/trips/{trip_id}/seats [get]
func (h *SeatStatusHandlerImpl) GetSeatAvailability(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("invalid trip id")
		return nil, ginext.NewBadRequestError("invalid trip id")
	}

	seats, err := h.seatStatusService.GetSeatAvailability(r.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("failed to get seat availability")
		return nil, err
	}

	return ginext.NewSuccessResponse(seats), nil
}
