package handler

import (
	"bus-booking/payment-service/internal/service"
	"bus-booking/shared/ginext"

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
// @Description Get constants by type (banks). Returns all types if type parameter is not specified.
// @Tags constants
// @Accept json
// @Produce json
// @Param type query string false "Constant type" Enums(banks)
// @Success 200 {object} ginext.Response "Constants"
// @Failure 400 {object} ginext.Response "Invalid type parameter"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/constants [get]
func (h *ConstantsHandlerImpl) GetConstants(r *ginext.Request) (*ginext.Response, error) {
	constType := r.GinCtx.Query("type")

	// Validate constant type if provided
	if constType != "" && constType != "banks" {
		return nil, ginext.NewBadRequestError("invalid type parameter. Valid values: banks")
	}

	// Fetch banks
	banks, err := h.constantsService.GetBanks(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get banks")
		return nil, ginext.NewInternalServerError("failed to get banks")
	}

	// Return banks array directly
	return ginext.NewSuccessResponse(banks), nil
}
