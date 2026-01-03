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
	GetBookedSeatIDs(ctx context.Context, tripID uuid.UUID) ([]uuid.UUID, error)

	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error)
	GetBookingByReference(ctx context.Context, reference string) (*model.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID uuid.UUID, statuses []model.BookingStatus, limit, offset int) ([]*model.Booking, int64, error)
	GetBookingsByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error)
	GetTripBookings(ctx context.Context, tripID uuid.UUID, page, limit int) ([]*model.Booking, int64, error)
	ListBookings(ctx context.Context, req model.ListBookingsRequest) ([]*model.Booking, int64, error)
	UpdateBooking(ctx context.Context, booking *model.Booking) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.BookingStatus) error
	CancelBooking(ctx context.Context, id uuid.UUID, reason string) error
	GetAllActiveBookingsByTripID(ctx context.Context, tripID uuid.UUID) ([]*model.Booking, error)
}

type bookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{db: db}
}

func (r *bookingRepositoryImpl) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.WithContext(ctx).
		Preload("BookingSeats").
		First(&booking, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

func (r *bookingRepositoryImpl) GetBookingByReference(ctx context.Context, reference string) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.WithContext(ctx).
		Preload("BookingSeats").
		First(&booking, "booking_reference = ?", reference).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking by reference: %w", err)
	}

	return &booking, nil
}

func (r *bookingRepositoryImpl) GetBookedSeatIDs(ctx context.Context, tripID uuid.UUID) ([]uuid.UUID, error) {
	var seatIDs []uuid.UUID
	if err := r.db.WithContext(ctx).
		Model(&model.BookingSeat{}).
		Select("booking_seats.seat_id").
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Where("bookings.trip_id = ? AND bookings.status IN ?", tripID, []model.BookingStatus{
			model.BookingStatusPending,
			model.BookingStatusConfirmed}).
		Scan(&seatIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get booked seat IDs: %w", err)
	}
	return seatIDs, nil
}

func (r *bookingRepositoryImpl) GetBookingsByUserID(ctx context.Context, userID uuid.UUID, statuses []model.BookingStatus, limit, offset int) ([]*model.Booking, int64, error) {
	var bookings []*model.Booking
	var total int64

	// Build query
	query := r.db.WithContext(ctx).Model(&model.Booking{}).Where("user_id = ?", userID)
	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Get bookings
	query = r.db.WithContext(ctx).Where("user_id = ?", userID)
	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}
	err := query.
		Preload("BookingSeats").
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

func (r *bookingRepositoryImpl) GetTripBookings(ctx context.Context, tripID uuid.UUID, page, limit int) ([]*model.Booking, int64, error) {
	offset := (page - 1) * limit
	return r.GetBookingsByTripID(ctx, tripID, limit, offset)
}

func (r *bookingRepositoryImpl) ListBookings(ctx context.Context, req model.ListBookingsRequest) ([]*model.Booking, int64, error) {
	var bookings []*model.Booking
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Booking{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.StartDate != "" {
		// Assuming ISO date string YYYY-MM-DD
		query = query.Where("created_at >= ?", req.StartDate)
	}

	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	order := "created_at DESC"
	if req.SortBy != "" {
		orderStr := "ASC"
		if req.Order == "desc" {
			orderStr = "DESC"
		}
		order = fmt.Sprintf("%s %s", req.SortBy, orderStr)
	}

	offset := (req.Page - 1) * req.PageSize
	err := query.
		Preload("BookingSeats").
		Order(order).
		Limit(req.PageSize).
		Offset(offset).
		Find(&bookings).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list bookings: %w", err)
	}

	return bookings, total, nil
}

func (r *bookingRepositoryImpl) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	if err := r.db.WithContext(ctx).Save(booking).Error; err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	return nil
}

func (r *bookingRepositoryImpl) CreateBooking(ctx context.Context, booking *model.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status model.BookingStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *bookingRepositoryImpl) CancelBooking(ctx context.Context, id uuid.UUID, reason string) error {
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

func (r *bookingRepositoryImpl) GetAllActiveBookingsByTripID(ctx context.Context, tripID uuid.UUID) ([]*model.Booking, error) {
	var bookings []*model.Booking
	if err := r.db.WithContext(ctx).
		Preload("BookingSeats").
		Where("trip_id = ? AND status IN ?", tripID, []model.BookingStatus{
			model.BookingStatusPending,
			model.BookingStatusConfirmed,
		}).
		Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("failed to get active bookings: %w", err)
	}
	return bookings, nil
}
