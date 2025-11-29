package handler

import (
	"time"

	"bus-booking/shared/ginext"
	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type TripHandler interface {
	SearchTrips(r *ginext.Request) (*ginext.Response, error)
	GetTrip(r *ginext.Request) (*ginext.Response, error)
	ListTripsByRoute(r *ginext.Request) (*ginext.Response, error)

	// Admin only
	CreateTrip(r *ginext.Request) (*ginext.Response, error)
	UpdateTrip(r *ginext.Request) (*ginext.Response, error)
	DeleteTrip(r *ginext.Request) (*ginext.Response, error)
}

type TripHandlerImpl struct {
	tripService service.TripService
}

func NewTripHandler(tripService service.TripService) TripHandler {
	return &TripHandlerImpl{
		tripService: tripService,
	}
}

// SearchTrips godoc
// @Summary Search trips
// @Description Search for available trips based on origin, destination, and other criteria
// @Tags trips
// @Accept json
// @Produce json
// @Param request body model.TripSearchRequest true "Trip search criteria"
// @Success 200 {object} ginext.Response{data=model.TripSearchResponse} "List of matching trips"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/trips/search [post]
func (h *TripHandlerImpl) SearchTrips(r *ginext.Request) (*ginext.Response, error) {
	var req model.TripSearchRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	trips, total, err := h.tripService.SearchTrips(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search trips")
		return nil, err
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	if req.Limit == 0 {
		totalPages = 1
	}

	response := model.TripSearchResponse{
		Trips:      trips,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}

	return ginext.NewSuccessResponse(response), nil
}

// GetTrip godoc
// @Summary Get trip by ID
// @Description Get detailed information about a specific trip
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "Trip ID" format(uuid)
// @Success 200 {object} ginext.Response{data=model.Trip} "Trip details"
// @Failure 400 {object} ginext.Response "Invalid trip ID"
// @Failure 404 {object} ginext.Response "Trip not found"
// @Router /api/v1/trips/{id} [get]
func (h *TripHandlerImpl) GetTrip(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	trip, err := h.tripService.GetTripByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("trip_id", idStr).Msg("Failed to get trip")
		return nil, err
	}

	return ginext.NewSuccessResponse(trip), nil
}

// CreateTrip godoc
// @Summary Create a new trip
// @Description Create a new trip with route, bus, and schedule information
// @Tags trips
// @Accept json
// @Produce json
// @Param request body model.CreateTripRequest true "Trip creation data"
// @Success 201 {object} ginext.Response{data=model.Trip} "Created trip"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/trips [post]
func (h *TripHandlerImpl) CreateTrip(r *ginext.Request) (*ginext.Response, error) {
	var req model.CreateTripRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	trip, err := h.tripService.CreateTrip(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create trip")
		return nil, err
	}

	return ginext.NewCreatedResponse(trip), nil
}

// UpdateTrip godoc
// @Summary Update trip
// @Description Update trip information such as schedule, price, or status
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "Trip ID" format(uuid)
// @Param request body model.UpdateTripRequest true "Trip update data"
// @Success 200 {object} ginext.Response{data=model.Trip} "Updated trip"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/trips/{id} [put]
func (h *TripHandlerImpl) UpdateTrip(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	var req model.UpdateTripRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	trip, err := h.tripService.UpdateTrip(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("trip_id", idStr).Msg("Failed to update trip")
		return nil, err
	}

	return ginext.NewSuccessResponse(trip), nil
}

// DeleteTrip godoc
// @Summary Delete trip
// @Description Delete a trip by ID
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "Trip ID" format(uuid)
// @Success 200 {object} ginext.Response "Success message"
// @Failure 400 {object} ginext.Response "Invalid trip ID"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/trips/{id} [delete]
func (h *TripHandlerImpl) DeleteTrip(r *ginext.Request) (*ginext.Response, error) {
	idStr := r.GinCtx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	err = h.tripService.DeleteTrip(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("trip_id", idStr).Msg("Failed to delete trip")
		return nil, err
	}

	return ginext.NewNoContentResponse(), nil
}

// ListTripsByRoute godoc
// @Summary List trips by route and date
// @Description Get all trips for a specific route on a given date
// @Tags trips
// @Accept json
// @Produce json
// @Param route_id path string true "Route ID" format(uuid)
// @Param date query string true "Date in YYYY-MM-DD format" example(2024-01-15)
// @Success 200 {object} ginext.Response{data=[]model.Trip} "List of trips"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/{route_id}/trips [get]
func (h *TripHandlerImpl) ListTripsByRoute(r *ginext.Request) (*ginext.Response, error) {
	routeIDStr := r.GinCtx.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid route ID")
	}

	var req model.ListTripsByRouteRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	date, err := parseDate(req.Date)
	if err != nil {
		return nil, ginext.NewBadRequestError("invalid date format")
	}

	trips, err := h.tripService.GetTripsByRouteAndDate(r.Context(), routeID, date)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list trips by route")
		return nil, err
	}

	return ginext.NewSuccessResponse(trips), nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
