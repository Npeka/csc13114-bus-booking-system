package handler

import (
	"time"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/rs/zerolog/log"
)

type StatisticsHandler interface {
	GetBookingStats(r *ginext.Request) (*ginext.Response, error)
	GetPopularTrips(r *ginext.Request) (*ginext.Response, error)
}

type StatisticsHandlerImpl struct {
	service service.StatisticsService
}

func NewStatisticsHandler(
	service service.StatisticsService,
) StatisticsHandler {
	return &StatisticsHandlerImpl{
		service: service,
	}
}

// GetBookingStats godoc
// @Summary Get booking statistics
// @Description Get booking statistics for a date range
// @Tags statistics
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} ginext.Response{data=model.BookingStatsResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/statistics/bookings [get]
func (h *StatisticsHandlerImpl) GetBookingStats(r *ginext.Request) (*ginext.Response, error) {
	var req model.BookingStatsRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("Invalid query parameters")
		return nil, ginext.NewBadRequestError("start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		log.Error().Err(err).Msg("Invalid start_date format")
		return nil, ginext.NewBadRequestError("invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		log.Error().Err(err).Msg("Invalid end_date format")
		return nil, ginext.NewBadRequestError("invalid end_date format, use YYYY-MM-DD")
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	stats, err := h.service.GetBookingStats(r.Context(), startDate, endDate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get booking statistics")
		return nil, err
	}

	return ginext.NewSuccessResponse(stats, "Statistics retrieved successfully"), nil
}

// GetPopularTrips godoc
// @Summary Get popular trips
// @Description Get popular trips based on booking statistics
// @Tags statistics
// @Produce json
// @Param limit query int false "Number of trips to return" default(10)
// @Param days query int false "Number of days to look back" default(30)
// @Success 200 {object} ginext.Response{data=[]model.TripStatsResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/statistics/popular-trips [get]
func (h *StatisticsHandlerImpl) GetPopularTrips(r *ginext.Request) (*ginext.Response, error) {
	var req model.PopularTripsRequest
	if err := r.GinCtx.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("Invalid query parameters")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}
	if req.Days < 1 || req.Days > 365 {
		req.Days = 30
	}

	trips, err := h.service.GetPopularTrips(r.Context(), req.Limit, req.Days)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get popular trips")
		return nil, err
	}

	return ginext.NewSuccessResponse(trips, "Popular trips retrieved successfully"), nil
}
