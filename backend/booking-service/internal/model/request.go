package model

import (
	"bus-booking/booking-service/internal/model/payment"
	"time"

	"github.com/google/uuid"
)

// CreateBookingRequest represents simplified booking creation request
// Backend will calculate price from Trip Service
type CreateBookingRequest struct {
	TripID  uuid.UUID   `json:"trip_id" binding:"required"`
	SeatIDs []uuid.UUID `json:"seat_ids" binding:"required,min=1,max=10,dive"`
	Notes   string      `json:"notes,omitempty"`
}

// CreateGuestBookingRequest represents guest booking creation (without authentication)
type CreateGuestBookingRequest struct {
	CreateBookingRequest
	FullName string `json:"full_name" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=15"`
}

// CancelBookingRequest represents booking cancellation request
type CancelBookingRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Reason string    `json:"reason" binding:"required"`
}

// UpdateBookingStatusRequest represents booking status update request
type UpdateBookingStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// InitSeatsRequest represents request to initialize seats for a trip
type InitSeatsRequest struct {
	Seats []SeatInitData `json:"seats" binding:"required,min=1,dive"`
}

// SeatInitData represents seat initialization data from trip service
type SeatInitData struct {
	SeatID     uuid.UUID `json:"seat_id" binding:"required"`
	SeatNumber string    `json:"seat_number" binding:"required"`
}

// ReserveSeatRequest represents seat reservation request
type ReserveSeatRequest struct {
	TripID             uuid.UUID `json:"trip_id" validate:"required"`
	SeatID             uuid.UUID `json:"seat_id" validate:"required"`
	UserID             uuid.UUID `json:"user_id" validate:"required"`
	ReservationMinutes int       `json:"reservation_minutes,omitempty"`
}

// ReleaseSeatRequest represents seat release request
type ReleaseSeatRequest struct {
	TripID uuid.UUID `json:"trip_id" validate:"required"`
	SeatID uuid.UUID `json:"seat_id" validate:"required"`
}

// BookingResponse represents simplified booking response
type BookingResponse struct {
	ID                uuid.UUID                 `json:"id"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	BookingReference  string                    `json:"booking_reference"`
	TripID            uuid.UUID                 `json:"trip_id"`
	UserID            uuid.UUID                 `json:"user_id"`
	TotalAmount       int                       `json:"total_amount"`
	Status            BookingStatus             `json:"status"`
	TransactionStatus payment.TransactionStatus `json:"transaction_status"`
	TransactionID     uuid.UUID                 `json:"transaction_id,omitempty"`
	Notes             string                    `json:"notes,omitempty"`
	ExpiresAt         *time.Time                `json:"expires_at,omitempty"`
	ConfirmedAt       *time.Time                `json:"confirmed_at,omitempty"`
	CancelledAt       *time.Time                `json:"cancelled_at,omitempty"`

	// Seats info
	Seats       []BookingSeatResponse        `json:"seats"`
	Transaction *payment.TransactionResponse `json:"transaction,omitempty"`
}

// BookingSeatResponse represents booking seat in response
type BookingSeatResponse struct {
	ID              uuid.UUID `json:"id"`
	SeatID          uuid.UUID `json:"seat_id"`
	SeatNumber      string    `json:"seat_number"`
	SeatType        string    `json:"seat_type"`
	Floor           int       `json:"floor"`
	Price           float64   `json:"price"`
	PriceMultiplier float64   `json:"price_multiplier"`
}

// PaymentResponse represents payment response
type PaymentResponse struct {
	BookingID       uuid.UUID `json:"booking_id"`
	Amount          float64   `json:"amount"`
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
	Status          string    `json:"status"`
	TransactionID   string    `json:"transaction_id"`
	ProcessedAt     time.Time `json:"processed_at"`
}

// FeedbackResponse represents feedback response
type FeedbackResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	BookingID uuid.UUID `json:"booking_id"`
	TripID    uuid.UUID `json:"trip_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// SeatAvailabilityResponse represents seat availability response
type SeatAvailabilityResponse struct {
	TripID         uuid.UUID                 `json:"trip_id"`
	AvailableSeats []uuid.UUID               `json:"available_seats"`
	ReservedSeats  []uuid.UUID               `json:"reserved_seats"`
	BookedSeats    []uuid.UUID               `json:"booked_seats"`
	SeatDetails    map[uuid.UUID]*SeatStatus `json:"seat_details"`
}

// Paginated responses
type PaginatedBookingResponse struct {
	Data       []*BookingResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int64              `json:"total_pages"`
}

type PaginatedFeedbackResponse struct {
	Data       []*FeedbackResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int64               `json:"total_pages"`
}

// Standard API responses
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// Seat lock requests
type LockSeatsRequest struct {
	TripID    uuid.UUID   `json:"trip_id" binding:"required"`
	SeatIDs   []uuid.UUID `json:"seat_ids" binding:"required,min=1,max=10,dive"`
	SessionID string      `json:"session_id" binding:"required"`
}

type UnlockSeatsRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

// GetSeatStatusRequest represents request to check seat status for a trip
type GetSeatStatusRequest struct {
	SeatIDs []string `form:"seat_ids"`
}

// SeatStatusItem represents booking status of a single seat
type SeatStatusItem struct {
	SeatID   uuid.UUID `json:"seat_id"`
	IsBooked bool      `json:"is_booked"`
	IsLocked bool      `json:"is_locked"`
}

// SeatBookingStatus represents booking status (for backward compatibility)
type SeatBookingStatus struct {
	IsBooked bool `json:"is_booked"`
	IsLocked bool `json:"is_locked"`
}
