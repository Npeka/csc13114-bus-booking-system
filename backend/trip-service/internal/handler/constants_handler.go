package handler

import (
	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/service"

	"github.com/rs/zerolog/log"
)

type ConstantsHandler interface {
	GetConstants(r *ginext.Request) (*ginext.Response, error)
}

type ConstantsHandlerImpl struct {
	constantsService service.ConstantsService
}

func NewConstantsHandler(constantsService service.ConstantsService) ConstantsHandler {
	return &ConstantsHandlerImpl{
		constantsService: constantsService,
	}
}

// GetConstants godoc
// @Summary Get constants
// @Description Get constants by type (bus, route, trip). Returns all types if type parameter is not specified.
// @Tags constants
// @Accept json
// @Produce json
// @Param type query string false "Constant type" Enums(bus, route, trip)
// @Success 200 {object} ginext.Response "Constants"
// @Failure 400 {object} ginext.Response "Invalid type parameter"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/constants [get]
func (h *ConstantsHandlerImpl) GetConstants(r *ginext.Request) (*ginext.Response, error) {
	constType := r.GinCtx.Query("type")

	switch constType {
	case "bus":
		busConstants, err := h.constantsService.GetBusConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get bus constants")
			return nil, err
		}
		return ginext.NewSuccessResponse(busConstants), nil

	case "route":
		routeConstants, err := h.constantsService.GetRouteConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get route constants")
			return nil, err
		}
		return ginext.NewSuccessResponse(routeConstants), nil

	case "trip":
		tripConstants, err := h.constantsService.GetTripConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get trip constants")
			return nil, err
		}
		return ginext.NewSuccessResponse(tripConstants), nil

	case "":
		// Return all constants if type is not specified
		allConstants, err := h.constantsService.GetAllConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get all constants")
			return nil, err
		}
		return ginext.NewSuccessResponse(allConstants), nil

	default:
		return nil, ginext.NewBadRequestError("invalid type parameter. Valid values: bus, route, trip")
	}
}
