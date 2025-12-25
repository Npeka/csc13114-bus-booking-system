package service

import (
	"context"
	"testing"
	"time"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewStatisticsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingStatsRepository(ctrl)
	service := NewStatisticsService(mockRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &StatisticsServiceImpl{}, service)
}

func TestGetBookingStats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingStatsRepository(ctrl)
	service := NewStatisticsService(mockRepo)

	ctx := context.Background()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	expectedStats := &model.BookingStats{
		TotalBookings:     100,
		TotalRevenue:      5000000,
		CancelledBookings: 10,
		CompletedBookings: 85,
		AverageRating:     4.5,
	}

	mockRepo.EXPECT().
		GetBookingStatsByDateRange(ctx, startDate, endDate).
		Return(expectedStats, nil).
		Times(1)

	result, err := service.GetBookingStats(ctx, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStats.TotalBookings, result.TotalBookings)
	assert.Equal(t, expectedStats.TotalRevenue, result.TotalRevenue)
	assert.Equal(t, expectedStats.CancelledBookings, result.CancelledBookings)
	assert.Equal(t, expectedStats.CompletedBookings, result.CompletedBookings)
	assert.Equal(t, expectedStats.AverageRating, result.AverageRating)
	assert.Equal(t, startDate, result.StartDate)
	assert.Equal(t, endDate, result.EndDate)
}

func TestGetBookingStats_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingStatsRepository(ctrl)
	service := NewStatisticsService(mockRepo)

	ctx := context.Background()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	expectedErr := assert.AnError

	mockRepo.EXPECT().
		GetBookingStatsByDateRange(ctx, startDate, endDate).
		Return(nil, expectedErr).
		Times(1)

	result, err := service.GetBookingStats(ctx, startDate, endDate)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
}

func TestGetPopularTrips_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingStatsRepository(ctrl)
	service := NewStatisticsService(mockRepo)

	ctx := context.Background()
	limit := 5
	days := 30

	tripID1 := uuid.New()
	tripID2 := uuid.New()

	expectedStats := []*model.TripBookingStats{
		{
			TripID:        tripID1,
			TotalBookings: 50,
			TotalRevenue:  2500000,
			AverageRating: 4.8,
		},
		{
			TripID:        tripID2,
			TotalBookings: 45,
			TotalRevenue:  2250000,
			AverageRating: 4.6,
		},
	}

	mockRepo.EXPECT().
		GetPopularTrips(ctx, limit, days).
		Return(expectedStats, nil).
		Times(1)

	result, err := service.GetPopularTrips(ctx, limit, days)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, tripID1, result[0].TripID)
	assert.Equal(t, int64(50), result[0].TotalBookings)
	assert.Equal(t, float64(2500000), result[0].TotalRevenue)
	assert.Equal(t, 4.8, result[0].AverageRating)
}

func TestGetPopularTrips_EmptyResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingStatsRepository(ctrl)
	service := NewStatisticsService(mockRepo)

	ctx := context.Background()
	limit := 10
	days := 7

	mockRepo.EXPECT().
		GetPopularTrips(ctx, limit, days).
		Return([]*model.TripBookingStats{}, nil).
		Times(1)

	result, err := service.GetPopularTrips(ctx, limit, days)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetPopularTrips_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingStatsRepository(ctrl)
	service := NewStatisticsService(mockRepo)

	ctx := context.Background()
	limit := 5
	days := 30

	expectedErr := assert.AnError

	mockRepo.EXPECT().
		GetPopularTrips(ctx, limit, days).
		Return(nil, expectedErr).
		Times(1)

	result, err := service.GetPopularTrips(ctx, limit, days)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
}

func TestGetPopularTrips_NilStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingStatsRepository(ctrl)
	service := NewStatisticsService(mockRepo)

	ctx := context.Background()
	limit := 5
	days := 30

	mockRepo.EXPECT().
		GetPopularTrips(ctx, limit, days).
		Return(nil, nil).
		Times(1)

	result, err := service.GetPopularTrips(ctx, limit, days)

	assert.NoError(t, err)
	assert.Empty(t, result)
}
