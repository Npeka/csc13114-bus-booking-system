package handler

import (
	"net/http"
	"time"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TripHandler interface {
	SearchTrips(c *gin.Context)
	GetTrip(c *gin.Context)
	CreateTrip(c *gin.Context)
	UpdateTrip(c *gin.Context)
	DeleteTrip(c *gin.Context)
	ListTripsByRoute(c *gin.Context)
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
// @Success 200 {object} model.TripSearchResponse "List of matching trips"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/trips/search [post]
func (h *TripHandlerImpl) SearchTrips(c *gin.Context) {
	var req model.TripSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trips, total, err := h.tripService.SearchTrips(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	if req.Limit == 0 {
		totalPages = 1
	}

	// Convert Trip to TripDetail
	var tripDetails []model.TripDetail
	for _, trip := range trips {
		tripDetails = append(tripDetails, model.TripDetail{
			ID:            trip.ID,
			RouteID:       trip.RouteID,
			BusID:         trip.BusID,
			DepartureTime: trip.DepartureTime,
			ArrivalTime:   trip.ArrivalTime,
			BasePrice:     trip.BasePrice,
			Status:        trip.Status,
		})
	}

	response := model.TripSearchResponse{
		Trips:      tripDetails,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetTrip godoc
// @Summary Get trip by ID
// @Description Get detailed information about a specific trip
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "Trip ID" format(uuid)
// @Success 200 {object} model.Trip "Trip details"
// @Failure 400 {object} map[string]string "Invalid trip ID"
// @Failure 404 {object} map[string]string "Trip not found"
// @Router /api/v1/trips/{id} [get]
func (h *TripHandlerImpl) GetTrip(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID"})
		return
	}

	trip, err := h.tripService.GetTripByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	c.JSON(http.StatusOK, trip)
}

// CreateTrip godoc
// @Summary Create a new trip
// @Description Create a new trip with route, bus, and schedule information
// @Tags trips
// @Accept json
// @Produce json
// @Param request body model.CreateTripRequest true "Trip creation data"
// @Success 201 {object} model.Trip "Created trip"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/trips [post]
func (h *TripHandlerImpl) CreateTrip(c *gin.Context) {
	var req model.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip, err := h.tripService.CreateTrip(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trip)
}

// UpdateTrip godoc
// @Summary Update trip
// @Description Update trip information such as schedule, price, or status
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "Trip ID" format(uuid)
// @Param request body model.UpdateTripRequest true "Trip update data"
// @Success 200 {object} model.Trip "Updated trip"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/trips/{id} [put]
func (h *TripHandlerImpl) UpdateTrip(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID"})
		return
	}

	var req model.UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip, err := h.tripService.UpdateTrip(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trip)
}

// DeleteTrip godoc
// @Summary Delete trip
// @Description Delete a trip by ID
// @Tags trips
// @Accept json
// @Produce json
// @Param id path string true "Trip ID" format(uuid)
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid trip ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/trips/{id} [delete]
func (h *TripHandlerImpl) DeleteTrip(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID"})
		return
	}

	err = h.tripService.DeleteTrip(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trip deleted successfully"})
}

// ListTripsByRoute godoc
// @Summary List trips by route and date
// @Description Get all trips for a specific route on a given date
// @Tags trips
// @Accept json
// @Produce json
// @Param route_id path string true "Route ID" format(uuid)
// @Param date query string true "Date in YYYY-MM-DD format" example(2024-01-15)
// @Success 200 {array} model.Trip "List of trips"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/routes/{route_id}/trips [get]
func (h *TripHandlerImpl) ListTripsByRoute(c *gin.Context) {
	routeIDStr := c.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid route ID"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	date, err := parseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	trips, err := h.tripService.GetTripsByRouteAndDate(c.Request.Context(), routeID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trips)
}

// parseDate parses date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
