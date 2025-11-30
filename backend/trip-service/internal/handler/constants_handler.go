package handler

import (
	"time"

	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/constants"
	"bus-booking/trip-service/internal/service"

	"github.com/rs/zerolog/log"
)

type ConstantsHandler interface {
	GetConstants(r *ginext.Request) (*ginext.Response, error)
}

type ConstantsHandlerImpl struct {
	constantsService service.ConstantsService
	cacheService     service.CacheService
}

func NewConstantsHandler(constantsService service.ConstantsService, cacheService service.CacheService) ConstantsHandler {
	return &ConstantsHandlerImpl{
		constantsService: constantsService,
		cacheService:     cacheService,
	}
}

// GetConstants godoc
// @Summary Get constants
// @Description Get constants by type (bus, route, trip, search_filters, cities). Returns all types if type parameter is not specified.
// @Tags constants
// @Accept json
// @Produce json
// @Param type query string false "Constant type" Enums(bus, route, trip, search_filters, cities)
// @Success 200 {object} ginext.Response "Constants"
// @Failure 400 {object} ginext.Response "Invalid type parameter"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/constants [get]
func (h *ConstantsHandlerImpl) GetConstants(r *ginext.Request) (*ginext.Response, error) {
	constTypeStr := r.GinCtx.Query("type")
	constType := constants.ConstantType(constTypeStr)

	// Validate constant type if provided
	if constTypeStr != "" && !constType.IsValid() {
		return nil, ginext.NewBadRequestError("invalid type parameter. Valid values: bus, route, trip, search_filters, cities")
	}

	// Determine cache key
	cacheKey := "all"
	if constTypeStr != "" {
		cacheKey = string(constType)
	}

	// Try to get from cache first
	cachedData, err := h.cacheService.GetConstants(r.Context(), cacheKey)
	if err == nil && cachedData != nil {
		log.Debug().Str("cache_key", cacheKey).Msg("Cache hit for constants")
		return ginext.NewSuccessResponse(cachedData), nil
	}

	// Cache miss, fetch from service
	var data interface{}
	switch constType {
	case constants.ConstantTypeBus:
		data, err = h.constantsService.GetBusConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get bus constants")
			return nil, err
		}

	case constants.ConstantTypeRoute:
		data, err = h.constantsService.GetRouteConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get route constants")
			return nil, err
		}

	case constants.ConstantTypeTrip:
		data, err = h.constantsService.GetTripConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get trip constants")
			return nil, err
		}

	case constants.ConstantTypeSearchFilters:
		data, err = h.constantsService.GetSearchFilterConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get search filter constants")
			return nil, err
		}

	case constants.ConstantTypeCities:
		data, err = h.constantsService.GetCities(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get cities")
			return nil, err
		}

	default:
		// Return all constants if type is not specified
		data, err = h.constantsService.GetAllConstants(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get all constants")
			return nil, err
		}
	}

	// Cache the result with 24-hour TTL
	if err := h.cacheService.SetConstants(r.Context(), cacheKey, data, 24*time.Hour); err != nil {
		log.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to cache constants")
		// Continue even if caching fails
	}

	return ginext.NewSuccessResponse(data), nil
}
