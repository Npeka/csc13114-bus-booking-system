package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"bus-booking/booking-service/internal/model"
)

type BookingStatsRepository interface {
	GetBookingStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*model.BookingStats, error)
	GetPopularTrips(ctx context.Context, limit int, days int) ([]*model.TripBookingStats, error)
}

type bookingStatsRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingStatsRepository(db *gorm.DB) BookingStatsRepository {
	return &bookingStatsRepositoryImpl{db: db}
}

func (r *bookingStatsRepositoryImpl) GetBookingStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*model.BookingStats, error) {
	var stats model.BookingStats

	// Total bookings
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalBookings).Error; err != nil {
		return nil, fmt.Errorf("failed to count total bookings: %w", err)
	}

	// Total revenue (only confirmed bookings)
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, model.BookingStatusConfirmed).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total revenue: %w", err)
	}

	// Cancelled bookings
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, model.BookingStatusCancelled).
		Count(&stats.CancelledBookings).Error; err != nil {
		return nil, fmt.Errorf("failed to count cancelled bookings: %w", err)
	}

	// Completed bookings (confirmed status)
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, model.BookingStatusConfirmed).
		Count(&stats.CompletedBookings).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed bookings: %w", err)
	}

	// Average rating
	if err := r.db.WithContext(ctx).
		Table("reviews f").
		Joins("JOIN bookings b ON f.booking_id = b.id").
		Where("b.created_at BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(AVG(f.rating), 0)").
		Scan(&stats.AverageRating).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average rating: %w", err)
	}

	return &stats, nil
}

func (r *bookingStatsRepositoryImpl) GetPopularTrips(ctx context.Context, limit int, days int) ([]*model.TripBookingStats, error) {
	var stats []*model.TripBookingStats

	startDate := time.Now().UTC().AddDate(0, 0, -days)

	err := r.db.WithContext(ctx).
		Table("bookings b").
		Select(`
			b.trip_id,
			COUNT(*) as total_bookings,
			COALESCE(SUM(CASE WHEN b.status = 'CONFIRMED' THEN b.total_amount ELSE 0 END), 0) as total_revenue,
			COALESCE(AVG(f.rating), 0) as average_rating
		`).
		Joins("LEFT JOIN reviews f ON f.booking_id = b.id").
		Where("b.created_at >= ?", startDate).
		Group("b.trip_id").
		Order("total_bookings DESC").
		Limit(limit).
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get popular trips: %w", err)
	}

	return stats, nil
}
