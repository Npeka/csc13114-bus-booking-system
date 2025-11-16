package repository

import (
	"context"
	"time"

	"bus-booking/booking-service/internal/model"

	"github.com/google/uuid"
)

// BookingRepository defines the interface for booking data operations
type BookingRepository interface {
	// Booking operations
	CreateBooking(ctx context.Context, booking *model.Booking) error
	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error)
	GetBookingsByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Booking, int64, error)
	UpdateBooking(ctx context.Context, booking *model.Booking) error
	CancelBooking(ctx context.Context, id uuid.UUID, reason string) error

	// Seat operations
	GetSeatStatusByTripID(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error)
	UpdateSeatStatus(ctx context.Context, seatStatus *model.SeatStatus) error
	BulkUpdateSeatStatus(ctx context.Context, seatStatuses []*model.SeatStatus) error
	GetAvailableSeats(ctx context.Context, tripID uuid.UUID) ([]*model.SeatStatus, error)
	ReserveSeat(ctx context.Context, tripID, seatID uuid.UUID, userID uuid.UUID, reservationTime time.Duration) error
	ReleaseSeat(ctx context.Context, tripID, seatID uuid.UUID) error

	// Payment methods
	GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethod, error)
	GetPaymentMethodByID(ctx context.Context, id uuid.UUID) (*model.PaymentMethod, error)
	GetPaymentMethodByCode(ctx context.Context, code string) (*model.PaymentMethod, error)

	// Feedback operations
	CreateFeedback(ctx context.Context, feedback *model.Feedback) error
	GetFeedbackByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Feedback, error)
	GetFeedbacksByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Feedback, int64, error)
	UpdateFeedback(ctx context.Context, feedback *model.Feedback) error
	DeleteFeedback(ctx context.Context, id uuid.UUID) error

	// Statistics
	GetBookingStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*model.BookingStats, error)
	GetPopularTrips(ctx context.Context, limit int, days int) ([]*model.TripBookingStats, error)
}

// BookingStats represents booking statistics
type BookingStats struct {
	TotalBookings     int64   `json:"total_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	CancelledBookings int64   `json:"cancelled_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	AverageRating     float64 `json:"average_rating"`
}

// TripBookingStats represents trip booking statistics
type TripBookingStats struct {
	TripID        uuid.UUID `json:"trip_id"`
	TotalBookings int64     `json:"total_bookings"`
	TotalRevenue  float64   `json:"total_revenue"`
	AverageRating float64   `json:"average_rating"`
}
