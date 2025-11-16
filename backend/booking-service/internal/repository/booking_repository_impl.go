package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bus-booking/booking-service/internal/model"
)

// bookingRepositoryImpl implements BookingRepository
type bookingRepositoryImpl struct {
	db *gorm.DB
}

// NewBookingRepository creates a new booking repository
func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{db: db}
}

// CreateBooking creates a new booking with seats
func (r *bookingRepositoryImpl) CreateBooking(ctx context.Context, booking *model.Booking) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create booking
		if err := tx.Create(booking).Error; err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}

		// Update seat statuses to booked
		for _, seat := range booking.BookingSeats {
			if err := tx.Model(&model.SeatStatus{}).
				Where("trip_id = ? AND seat_id = ?", booking.TripID, seat.SeatID).
				Updates(map[string]interface{}{
					"status":     "booked",
					"user_id":    booking.UserID,
					"booking_id": booking.ID,
					"updated_at": time.Now().UTC(),
				}).Error; err != nil {
				return fmt.Errorf("failed to update seat status: %w", err)
			}
		}

		return nil
	})
}

// GetBookingByID retrieves a booking by ID with all related data
func (r *bookingRepositoryImpl) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.WithContext(ctx).
		Preload("BookingSeats").
		Preload("PaymentMethod").
		First(&booking, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

// GetBookingsByUserID retrieves bookings for a user with pagination
func (r *bookingRepositoryImpl) GetBookingsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error) {
	var bookings []*model.Booking
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Get bookings
	err := r.db.WithContext(ctx).
		Preload("BookingSeats").
		Preload("PaymentMethod").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bookings).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings: %w", err)
	}

	return bookings, total, nil
}

// GetBookingsByTripID retrieves bookings for a trip with pagination
func (r *bookingRepositoryImpl) GetBookingsByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error) {
	var bookings []*model.Booking
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("trip_id = ?", tripID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Get bookings
	err := r.db.WithContext(ctx).
		Preload("BookingSeats").
		Preload("PaymentMethod").
		Where("trip_id = ?", tripID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bookings).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings: %w", err)
	}

	return bookings, total, nil
}

// UpdateBooking updates a booking
func (r *bookingRepositoryImpl) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	if err := r.db.WithContext(ctx).Save(booking).Error; err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	return nil
}

// CancelBooking cancels a booking and releases seats
func (r *bookingRepositoryImpl) CancelBooking(ctx context.Context, id uuid.UUID, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get booking
		var booking model.Booking
		if err := tx.Preload("BookingSeats").First(&booking, "id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}

		// Update booking status
		if err := tx.Model(&booking).Updates(map[string]interface{}{
			"status":              "cancelled",
			"cancellation_reason": reason,
			"cancelled_at":        time.Now().UTC(),
			"updated_at":          time.Now().UTC(),
		}).Error; err != nil {
			return fmt.Errorf("failed to update booking status: %w", err)
		}

		// Release seats
		for _, seat := range booking.BookingSeats {
			if err := tx.Model(&model.SeatStatus{}).
				Where("trip_id = ? AND seat_id = ?", booking.TripID, seat.SeatID).
				Updates(map[string]interface{}{
					"status":     "available",
					"user_id":    nil,
					"booking_id": nil,
					"updated_at": time.Now().UTC(),
				}).Error; err != nil {
				return fmt.Errorf("failed to release seat: %w", err)
			}
		}

		return nil
	})
}

// GetSeatStatusByTripID retrieves all seat statuses for a trip
func (r *bookingRepositoryImpl) GetSeatStatusByTripID(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error) {
	var seatStatuses []*model.SeatStatus
	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("seat_number").
		Find(&seatStatuses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get seat statuses: %w", err)
	}

	return seatStatuses, nil
}

// UpdateSeatStatus updates a single seat status
func (r *bookingRepositoryImpl) UpdateSeatStatus(ctx context.Context, seatStatus *model.SeatStatus) error {
	if err := r.db.WithContext(ctx).Save(seatStatus).Error; err != nil {
		return fmt.Errorf("failed to update seat status: %w", err)
	}
	return nil
}

// BulkUpdateSeatStatus updates multiple seat statuses
func (r *bookingRepositoryImpl) BulkUpdateSeatStatus(ctx context.Context, seatStatuses []*model.SeatStatus) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, seatStatus := range seatStatuses {
			if err := tx.Save(seatStatus).Error; err != nil {
				return fmt.Errorf("failed to update seat status: %w", err)
			}
		}
		return nil
	})
}

// GetAvailableSeats retrieves available seats for a trip
func (r *bookingRepositoryImpl) GetAvailableSeats(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error) {
	var seatStatuses []*model.SeatStatus
	err := r.db.WithContext(ctx).
		Where("trip_id = ? AND status = ?", tripID, "available").
		Order("seat_number").
		Find(&seatStatuses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get available seats: %w", err)
	}

	return seatStatuses, nil
}

// ReserveSeat reserves a seat for a user with expiration
func (r *bookingRepositoryImpl) ReserveSeat(ctx context.Context, tripID, seatID uuid.UUID, userID uuid.UUID, reservationTime time.Duration) error {
	expiresAt := time.Now().UTC().Add(reservationTime)

	err := r.db.WithContext(ctx).Model(&model.SeatStatus{}).
		Where("trip_id = ? AND seat_id = ? AND status = ?", tripID, seatID, "available").
		Updates(map[string]interface{}{
			"status":         "reserved",
			"user_id":        userID,
			"reserved_until": &expiresAt,
			"updated_at":     time.Now().UTC(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to reserve seat: %w", err)
	}

	return nil
}

// ReleaseSeat releases a reserved seat
func (r *bookingRepositoryImpl) ReleaseSeat(ctx context.Context, tripID, seatID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&model.SeatStatus{}).
		Where("trip_id = ? AND seat_id = ?", tripID, seatID).
		Updates(map[string]interface{}{
			"status":         "available",
			"user_id":        nil,
			"booking_id":     nil,
			"reserved_until": nil,
			"updated_at":     time.Now().UTC(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to release seat: %w", err)
	}

	return nil
}

// GetPaymentMethods retrieves all active payment methods
func (r *bookingRepositoryImpl) GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethod, error) {
	var paymentMethods []*model.PaymentMethod
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name").
		Find(&paymentMethods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	return paymentMethods, nil
}

// GetPaymentMethodByID retrieves a payment method by ID
func (r *bookingRepositoryImpl) GetPaymentMethodByID(ctx context.Context, id uuid.UUID) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := r.db.WithContext(ctx).First(&paymentMethod, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment method not found")
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	return &paymentMethod, nil
}

// GetPaymentMethodByCode retrieves a payment method by code
func (r *bookingRepositoryImpl) GetPaymentMethodByCode(ctx context.Context, code string) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := r.db.WithContext(ctx).First(&paymentMethod, "code = ?", code).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment method not found")
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	return &paymentMethod, nil
}

// CreateFeedback creates a new feedback
func (r *bookingRepositoryImpl) CreateFeedback(ctx context.Context, feedback *model.Feedback) error {
	if err := r.db.WithContext(ctx).Create(feedback).Error; err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}
	return nil
}

// GetFeedbackByBookingID retrieves feedback for a booking
func (r *bookingRepositoryImpl) GetFeedbackByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Feedback, error) {
	var feedback model.Feedback
	err := r.db.WithContext(ctx).First(&feedback, "booking_id = ?", bookingID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("feedback not found")
		}
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	return &feedback, nil
}

// GetFeedbacksByTripID retrieves feedbacks for a trip with pagination
func (r *bookingRepositoryImpl) GetFeedbacksByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Feedback, int64, error) {
	var feedbacks []*model.Feedback
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&model.Feedback{}).
		Where("trip_id = ?", tripID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count feedbacks: %w", err)
	}

	// Get feedbacks
	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&feedbacks).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get feedbacks: %w", err)
	}

	return feedbacks, total, nil
}

// UpdateFeedback updates a feedback
func (r *bookingRepositoryImpl) UpdateFeedback(ctx context.Context, feedback *model.Feedback) error {
	if err := r.db.WithContext(ctx).Save(feedback).Error; err != nil {
		return fmt.Errorf("failed to update feedback: %w", err)
	}
	return nil
}

// DeleteFeedback deletes a feedback
func (r *bookingRepositoryImpl) DeleteFeedback(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&model.Feedback{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete feedback: %w", err)
	}
	return nil
}

// GetBookingStatsByDateRange retrieves booking statistics for a date range
func (r *bookingRepositoryImpl) GetBookingStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*model.BookingStats, error) {
	var stats model.BookingStats

	// Total bookings
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalBookings).Error; err != nil {
		return nil, fmt.Errorf("failed to count total bookings: %w", err)
	}

	// Total revenue (only completed bookings)
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "completed").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total revenue: %w", err)
	}

	// Cancelled bookings
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "cancelled").
		Count(&stats.CancelledBookings).Error; err != nil {
		return nil, fmt.Errorf("failed to count cancelled bookings: %w", err)
	}

	// Completed bookings
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "completed").
		Count(&stats.CompletedBookings).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed bookings: %w", err)
	}

	// Average rating
	if err := r.db.WithContext(ctx).
		Table("feedbacks f").
		Joins("JOIN bookings b ON f.booking_id = b.id").
		Where("b.created_at BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(AVG(f.rating), 0)").
		Scan(&stats.AverageRating).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average rating: %w", err)
	}

	return &stats, nil
}

// GetPopularTrips retrieves popular trips based on booking count
func (r *bookingRepositoryImpl) GetPopularTrips(ctx context.Context, limit int, days int) ([]*model.TripBookingStats, error) {
	var stats []*model.TripBookingStats

	startDate := time.Now().UTC().AddDate(0, 0, -days)

	err := r.db.WithContext(ctx).
		Table("bookings b").
		Select(`
			b.trip_id,
			COUNT(*) as total_bookings,
			COALESCE(SUM(CASE WHEN b.status = 'completed' THEN b.total_amount ELSE 0 END), 0) as total_revenue,
			COALESCE(AVG(f.rating), 0) as average_rating
		`).
		Joins("LEFT JOIN feedbacks f ON f.booking_id = b.id").
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
