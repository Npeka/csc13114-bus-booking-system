package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	TripID          uuid.UUID   `json:"trip_id" validate:"required"`
	UserID          uuid.UUID   `json:"user_id" validate:"required"`
	SeatIDs         []uuid.UUID `json:"seat_ids" validate:"required,min=1,max=10"`
	PaymentMethodID uuid.UUID   `json:"payment_method_id" validate:"required"`
	TotalAmount     float64     `json:"total_amount" validate:"required,min=0"`
	SeatPrice       float64     `json:"seat_price" validate:"required,min=0"`
	PassengerName   string      `json:"passenger_name" validate:"required"`
	PassengerPhone  string      `json:"passenger_phone" validate:"required"`
	PassengerEmail  string      `json:"passenger_email,omitempty" validate:"omitempty,email"`
	SpecialRequests string      `json:"special_requests,omitempty"`
}

// CancelBookingRequest represents booking cancellation request
type CancelBookingRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Reason string    `json:"reason" validate:"required"`
}

// UpdateBookingStatusRequest represents booking status update request
type UpdateBookingStatusRequest struct {
	Status string `json:"status" validate:"required"`
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

// ProcessPaymentRequest represents payment processing request
type ProcessPaymentRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	Token     string    `json:"token,omitempty"`
}

// CreateFeedbackRequest represents feedback creation request
type CreateFeedbackRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Comment   string    `json:"comment,omitempty"`
}

// Response types

// BookingResponse represents booking response
type BookingResponse struct {
	ID                 uuid.UUID              `json:"id"`
	UserID             uuid.UUID              `json:"user_id"`
	TripID             uuid.UUID              `json:"trip_id"`
	Status             string                 `json:"status"`
	TotalAmount        float64                `json:"total_amount"`
	PassengerName      string                 `json:"passenger_name"`
	PassengerPhone     string                 `json:"passenger_phone"`
	PassengerEmail     string                 `json:"passenger_email,omitempty"`
	SpecialRequests    string                 `json:"special_requests,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	CompletedAt        *time.Time             `json:"completed_at,omitempty"`
	CancelledAt        *time.Time             `json:"cancelled_at,omitempty"`
	CancellationReason string                 `json:"cancellation_reason,omitempty"`
	Seats              []BookingSeatResponse  `json:"seats"`
	PaymentMethod      *PaymentMethodResponse `json:"payment_method,omitempty"`
}

// BookingSeatResponse represents booking seat response
type BookingSeatResponse struct {
	ID     uuid.UUID `json:"id"`
	SeatID uuid.UUID `json:"seat_id"`
	Price  float64   `json:"price"`
}

// PaymentMethodResponse represents payment method response
type PaymentMethodResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
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
	Limit      int                `json:"limit"`
	TotalPages int64              `json:"total_pages"`
}

type PaginatedFeedbackResponse struct {
	Data       []*FeedbackResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int64               `json:"total_pages"`
}

// Statistics responses
type BookingStatsResponse struct {
	TotalBookings     int64     `json:"total_bookings"`
	TotalRevenue      float64   `json:"total_revenue"`
	CancelledBookings int64     `json:"cancelled_bookings"`
	CompletedBookings int64     `json:"completed_bookings"`
	AverageRating     float64   `json:"average_rating"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
}

type TripStatsResponse struct {
	TripID        uuid.UUID `json:"trip_id"`
	TotalBookings int64     `json:"total_bookings"`
	TotalRevenue  float64   `json:"total_revenue"`
	AverageRating float64   `json:"average_rating"`
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

// Validation method
func (req *CreateBookingRequest) Validate() error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}
	if req.TripID == uuid.Nil {
		return fmt.Errorf("trip ID is required")
	}
	if len(req.SeatIDs) == 0 {
		return fmt.Errorf("at least one seat must be selected")
	}
	if req.TotalAmount <= 0 {
		return fmt.Errorf("total amount must be greater than 0")
	}
	if req.PassengerName == "" {
		return fmt.Errorf("passenger name is required")
	}
	if req.PassengerPhone == "" {
		return fmt.Errorf("passenger phone is required")
	}
	return nil
}
