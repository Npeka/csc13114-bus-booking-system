package handler

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SeatLockHandler interface {
	LockSeats(r *ginext.Request) (*ginext.Response, error)
	UnlockSeats(r *ginext.Request) (*ginext.Response, error)
	GetLockedSeats(r *ginext.Request) (*ginext.Response, error)
}

type SeatLockHandlerImpl struct {
	lockService service.SeatLockService
}

func NewSeatLockHandler(lockService service.SeatLockService) SeatLockHandler {
	return &SeatLockHandlerImpl{lockService: lockService}
}

// LockSeats godoc
// @Summary Lock seats temporarily
// @Description Lock selected seats for 15 minutes during booking process
// @Tags seat-locks
// @Accept json
// @Produce json
// @Param request body model.LockSeatsRequest true "Seat lock data"
// @Success 200 {object} ginext.Response "Seats locked successfully"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 409 {object} ginext.Response "Seats already locked"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/seat-locks [post]
func (h *SeatLockHandlerImpl) LockSeats(r *ginext.Request) (*ginext.Response, error) {
	var req model.LockSeatsRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.lockService.LockSeats(r.Context(), req.TripID, req.SeatIDs, req.SessionID); err != nil {
		log.Error().Err(err).Msg("Failed to lock seats")
		return nil, err
	}

	return ginext.NewSuccessResponse("Seats locked successfully"), nil
}

// UnlockSeats godoc
// @Summary Unlock seats
// @Description Release locked seats for a session
// @Tags seat-locks
// @Accept json
// @Produce json
// @Param request body model.UnlockSeatsRequest true "Unlock data"
// @Success 200 {object} ginext.Response "Seats unlocked successfully"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/seat-locks [delete]
func (h *SeatLockHandlerImpl) UnlockSeats(r *ginext.Request) (*ginext.Response, error) {
	var req model.UnlockSeatsRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if err := h.lockService.UnlockSeats(r.Context(), req.SessionID); err != nil {
		log.Error().Err(err).Msg("Failed to unlock seats")
		return nil, err
	}

	return ginext.NewSuccessResponse("Seats unlocked successfully"), nil
}

// GetLockedSeats godoc
// @Summary Get locked seats
// @Description Get list of locked seats for a trip
// @Tags seat-locks
// @Accept json
// @Produce json
// @Param trip_id path string true "Trip ID" format(uuid)
// @Success 200 {object} ginext.Response{data=[]uuid.UUID} "Locked seat IDs"
// @Failure 400 {object} ginext.Response "Invalid trip ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/trips/{trip_id}/locked-seats [get]
func (h *SeatLockHandlerImpl) GetLockedSeats(r *ginext.Request) (*ginext.Response, error) {
	tripIDStr := r.GinCtx.Param("trip_id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	seatIDs, err := h.lockService.GetLockedSeats(r.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripIDStr).Msg("Failed to get locked seats")
		return nil, err
	}

	return ginext.NewSuccessResponse(seatIDs), nil
}
