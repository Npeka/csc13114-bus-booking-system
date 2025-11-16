package handler

import (
	"net/http"
	"time"

	"bus-booking/trip-service/internal/model"
	"bus-booking/trip-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TripHandler struct {
	tripService service.TripService
}

func NewTripHandler(tripService service.TripService) *TripHandler {
	return &TripHandler{
		tripService: tripService,
	}
}

// SearchTrips handles trip search requests
func (h *TripHandler) SearchTrips(c *gin.Context) {
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

// GetTrip handles get trip by ID requests
func (h *TripHandler) GetTrip(c *gin.Context) {
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

// CreateTrip handles trip creation requests
func (h *TripHandler) CreateTrip(c *gin.Context) {
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

// UpdateTrip handles trip update requests
func (h *TripHandler) UpdateTrip(c *gin.Context) {
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

// DeleteTrip handles trip deletion requests
func (h *TripHandler) DeleteTrip(c *gin.Context) {
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

// ListTripsByRoute handles listing trips by route
func (h *TripHandler) ListTripsByRoute(c *gin.Context) {
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
