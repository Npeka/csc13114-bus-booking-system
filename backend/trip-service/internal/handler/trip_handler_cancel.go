package handler

import (
	"bus-booking/shared/ginext"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CancelTrip godoc
// @Summary Cancel trip
// @Description Cancel a trip and trigger refunds for paid bookings (Admin only)
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "Trip ID" format(uuid)
// @Success 200 {object} ginext.Response "Success message"
// @Failure 400 {object} ginext.Response "Invalid trip ID or status"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/trips/{id}/cancel [put]
func (h *TripHandlerImpl) CancelTrip(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("trip_id", idStr).Msg("Invalid trip ID")
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	if err = h.tripService.CancelTrip(r.Context(), id); err != nil {
		log.Error().Err(err).Str("trip_id", idStr).Msg("Failed to cancel trip")
		return nil, err
	}

	return ginext.NewSuccessResponse("Trip cancelled and refunds processing initiated"), nil
}
