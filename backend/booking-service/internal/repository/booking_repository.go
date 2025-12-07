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
	CheckSeatAvailability(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) (map[uuid.UUID]bool, error)
	GetBookedSeatsForTrip(ctx context.Context, tripID uuid.UUID) ([]uuid.UUID, error)

	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error)
	GetBookingsByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error)
	GetTripBookings(ctx context.Context, tripID uuid.UUID, page, limit int) ([]*model.Booking, int64, error)
	UpdateBooking(ctx context.Context, booking *model.Booking) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.BookingStatus) error
	CancelBooking(ctx context.Context, id uuid.UUID, reason string) error
}

type bookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{db: db}
}

func (r *bookingRepositoryImpl) CreateBooking(ctx context.Context, booking *model.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepositoryImpl) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.WithContext(ctx).
		Preload("BookingSeats").
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

	err := r.db.WithContext(ctx).
		Preload("BookingSeats").
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

// GetTripBookings with pagination
func (r *bookingRepositoryImpl) GetTripBookings(ctx context.Context, tripID uuid.UUID, page, limit int) ([]*model.Booking, int64, error) {
	offset := (page - 1) * limit
	return r.GetBookingsByTripID(ctx, tripID, limit, offset)
}

func (r *bookingRepositoryImpl) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	if err := r.db.WithContext(ctx).Save(booking).Error; err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of a booking
func (r *bookingRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status model.BookingStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *bookingRepositoryImpl) CancelBooking(ctx context.Context, id uuid.UUID, reason string) error {
	// Simply update booking status - no need to manage seats
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&model.Booking{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":              model.BookingStatusCancelled,
			"cancellation_reason": reason,
			"cancelled_at":        &now,
			"updated_at":          now,
		}).Error
}

// GetBookedSeatsForTrip returns all seat IDs that are booked for a trip with valid status
func (r *bookingRepositoryImpl) GetBookedSeatsForTrip(ctx context.Context, tripID uuid.UUID) ([]uuid.UUID, error) {
	var seatIDs []uuid.UUID

	// Get all booking_seats for bookings that are confirmed or pending (not cancelled/expired)
	err := r.db.WithContext(ctx).
		Model(&model.BookingSeat{}).
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Where("bookings.trip_id = ?", tripID).
		Where("bookings.status IN ?", []model.BookingStatus{
			model.BookingStatusPending,
			model.BookingStatusConfirmed,
		}).
		Pluck("booking_seats.seat_id", &seatIDs).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to get booked seats: %w", err)
	}

	return seatIDs, nil
}

// CheckSeatAvailability checks if seats are available (not booked) for a trip
// Returns a map of seatID -> isBooked (true if already booked, false if available)
func (r *bookingRepositoryImpl) CheckSeatAvailability(ctx context.Context, tripID uuid.UUID, seatIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	bookedSeats, err := r.GetBookedSeatsForTrip(ctx, tripID)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	bookedMap := make(map[uuid.UUID]bool)
	for _, seatID := range bookedSeats {
		bookedMap[seatID] = true
	}

	// Check each requested seat
	result := make(map[uuid.UUID]bool)
	for _, seatID := range seatIDs {
		result[seatID] = bookedMap[seatID]
	}

	return result, nil
}
