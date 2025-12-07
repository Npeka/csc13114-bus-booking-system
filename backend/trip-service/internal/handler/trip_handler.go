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
	ListTrips(r *ginext.Request) (*ginext.Response, error)
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
// @Description Search for available trips based on origin, destination, and other criteria. All filters are optional.
// @Tags trips
// @Accept json
// @Produce json
// @Param origin query string false "Origin city (partial match)"
// @Param destination query string false "Destination city (partial match)"
// @Param departure_time_start query string false "Departure time start (ISO8601 or HH:MM)" example(2025-12-01T06:00:00Z)
// @Param departure_time_end query string false "Departure time end (ISO8601 or HH:MM)" example(2025-12-01T22:00:00Z)
// @Param arrival_time_start query string false "Arrival time start (ISO8601 or HH:MM)"
// @Param arrival_time_end query string false "Arrival time end (ISO8601 or HH:MM)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param seat_types query []string false "Seat types" collectionFormat(multi)
// @Param amenities query []string false "Amenities" collectionFormat(multi)
// @Param status query string false "Trip status (for admin)" Enums(scheduled, in_progress, completed, cancelled)
// @Param sort_by query string false "Sort by field" Enums(price, departure_time, duration)
// @Param sort_order query string false "Sort order" Enums(asc, desc)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Success 200 {object} ginext.Response "List of matching trips with pagination"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/trips/search [get]
func (h *TripHandlerImpl) SearchTrips(r *ginext.Request) (*ginext.Response, error) {
	var req model.TripSearchRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("Query binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	trips, total, err := h.tripService.SearchTrips(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search trips")
		return nil, err
	}

	return ginext.NewPaginatedResponse(trips, req.Page, req.PageSize, total), nil
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
		log.Error().Err(err).Str("trip_id", idStr).Msg("Invalid trip ID")
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	trip, err := h.tripService.GetTripByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("trip_id", idStr).Msg("Failed to get trip")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToTripResponse(trip)), nil
}

// ListTrips godoc
// @Summary List trips
// @Description Get a paginated list of trips
// @Tags trips
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Success 200 {object} ginext.Response "Paginated trip list"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/trips [get]
func (h *TripHandlerImpl) ListTrips(r *ginext.Request) (*ginext.Response, error) {
	var req model.ListTripsRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("Query binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	trips, total, err := h.tripService.ListTrips(r.Context(), req.Page, req.PageSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list trips")
		return nil, err
	}

	return ginext.NewPaginatedResponse(model.ToTripResponseList(trips), req.Page, req.PageSize, total), nil
}

// ListTripsByRoute godoc
// @Summary List trips by route and date
// @Description Get all trips for a specific route on a given date
// @Tags trips
// @Accept json
// @Produce json
// @Param route_id path string true "Route ID" format(uuid)
// @Param date query string true "Date in DD/MM/YYYY format" example(15/01/2024)
// @Success 200 {object} ginext.Response{data=[]model.Trip} "List of trips"
// @Failure 400 {object} ginext.Response "Invalid request"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /api/v1/routes/{route_id}/trips [get]
func (h *TripHandlerImpl) ListTripsByRoute(r *ginext.Request) (*ginext.Response, error) {
	routeIDStr := r.GinCtx.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		log.Error().Err(err).Str("route_id", routeIDStr).Msg("Invalid route ID")
		return nil, ginext.NewBadRequestError("invalid route ID")
	}

	var req model.ListTripsByRouteRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		return nil, ginext.NewBadRequestError(err.Error())
	}

	date, err := parseDate(req.Date)
	if err != nil {
		log.Error().Err(err).Str("date", req.Date).Msg("Invalid date format")
		return nil, ginext.NewBadRequestError("invalid date format")
	}

	trips, err := h.tripService.GetTripsByRouteAndDate(r.Context(), routeID, date)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list trips by route")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToTripResponseList(trips)), nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("02/01/2006", dateStr)
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
		log.Error().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	trip, err := h.tripService.CreateTrip(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create trip")
		return nil, err
	}

	return ginext.NewCreatedResponse(model.ToTripResponse(trip)), nil
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
		log.Error().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	trip, err := h.tripService.UpdateTrip(r.Context(), id, &req)
	if err != nil {
		log.Error().Err(err).Str("trip_id", idStr).Msg("Failed to update trip")
		return nil, err
	}

	return ginext.NewSuccessResponse(model.ToTripResponse(trip)), nil
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
		log.Error().Err(err).Str("trip_id", idStr).Msg("Invalid trip ID")
		return nil, ginext.NewBadRequestError("invalid trip ID")
	}

	if err = h.tripService.DeleteTrip(r.Context(), id); err != nil {
		log.Error().Err(err).Str("trip_id", idStr).Msg("Failed to delete trip")
		return nil, err
	}

	return ginext.NewSuccessResponse("Trip deleted successfully"), nil
}
