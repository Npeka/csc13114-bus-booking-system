package service

import (
	"context"
	"time"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
)

type StatisticsService interface {
	GetBookingStats(ctx context.Context, startDate, endDate time.Time) (*model.BookingStatsResponse, error)
	GetPopularTrips(ctx context.Context, limit, days int) ([]*model.TripStatsResponse, error)
}

type StatisticsServiceImpl struct {
	repositories *repository.Repositories
}

func NewStatisticsService(repositories *repository.Repositories) StatisticsService {
	return &StatisticsServiceImpl{
		repositories: repositories,
	}
}

// GetBookingStats retrieves booking statistics for a date range
func (s *StatisticsServiceImpl) GetBookingStats(ctx context.Context, startDate, endDate time.Time) (*model.BookingStatsResponse, error) {
	stats, err := s.repositories.BookingStats.GetBookingStatsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &model.BookingStatsResponse{
		TotalBookings:     stats.TotalBookings,
		TotalRevenue:      stats.TotalRevenue,
		CancelledBookings: stats.CancelledBookings,
		CompletedBookings: stats.CompletedBookings,
		AverageRating:     stats.AverageRating,
		StartDate:         startDate,
		EndDate:           endDate,
	}, nil
}

// GetPopularTrips retrieves popular trips based on booking statistics
func (s *StatisticsServiceImpl) GetPopularTrips(ctx context.Context, limit, days int) ([]*model.TripStatsResponse, error) {
	stats, err := s.repositories.BookingStats.GetPopularTrips(ctx, limit, days)
	if err != nil {
		return nil, err
	}

	var responses []*model.TripStatsResponse
	for _, stat := range stats {
		response := &model.TripStatsResponse{
			TripID:        stat.TripID,
			TotalBookings: stat.TotalBookings,
			TotalRevenue:  stat.TotalRevenue,
			AverageRating: stat.AverageRating,
		}
		responses = append(responses, response)
	}

	return responses, nil
}
