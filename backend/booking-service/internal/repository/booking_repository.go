package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bus-booking/booking-service/internal/model"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *model.Booking) error
	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error)
	GetBookingsByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error)
	UpdateBooking(ctx context.Context, booking *model.Booking) error
	CancelBooking(ctx context.Context, id uuid.UUID, reason string) error
}

type bookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{db: db}
}

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

func (r *bookingRepositoryImpl) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	if err := r.db.WithContext(ctx).Save(booking).Error; err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	return nil
}

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
