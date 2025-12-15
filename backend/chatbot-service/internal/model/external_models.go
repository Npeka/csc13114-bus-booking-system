package model

import (
	"time"

	"github.com/google/uuid"
)

// APIResponse is the standard API response wrapper used by all services
type APIResponse[T any] struct {
	Data  T           `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
	Meta  interface{} `json:"meta,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// TripDetailResponse represents a detailed trip from trip-service
type TripDetailResponse struct {
	ID             uuid.UUID    `json:"id"`
	RouteID        uuid.UUID    `json:"route_id"`
	BusID          uuid.UUID    `json:"bus_id"`
	DepartureTime  time.Time    `json:"departure_time"`
	ArrivalTime    time.Time    `json:"arrival_time"`
	BasePrice      float64      `json:"base_price"`
	Status         string       `json:"status"`
	AvailableSeats int          `json:"available_seats"`
	TotalSeats     int          `json:"total_seats"`
	Route          *RouteDetail `json:"route,omitempty"`
	Bus            *BusDetail   `json:"bus,omitempty"`
	PriceTiers     []PriceTier  `json:"price_tiers,omitempty"`
}

type RouteDetail struct {
	ID          uuid.UUID   `json:"id"`
	Origin      string      `json:"origin"`
	Destination string      `json:"destination"`
	Distance    float64     `json:"distance"`
	Duration    int         `json:"duration"` // in minutes
	RouteStops  []RouteStop `json:"route_stops,omitempty"`
}

type RouteStop struct {
	ID         uuid.UUID `json:"id"`
	RouteID    uuid.UUID `json:"route_id"`
	StopOrder  int       `json:"stop_order"`
	Location   string    `json:"location"`
	ArrivalMin int       `json:"arrival_min"` // minutes from start
}

type BusDetail struct {
	ID           uuid.UUID    `json:"id"`
	LicensePlate string       `json:"license_plate"`
	BusType      string       `json:"bus_type"`
	TotalSeats   int          `json:"total_seats"`
	Floor        int          `json:"floor"`
	Amenities    []string     `json:"amenities"`
	Seats        []SeatDetail `json:"seats,omitempty"`
}

type SeatDetail struct {
	ID          uuid.UUID   `json:"id"`
	BusID       uuid.UUID   `json:"bus_id"`
	SeatNumber  string      `json:"seat_number"`
	Floor       int         `json:"floor"`
	Row         int         `json:"row"`
	Column      int         `json:"column"`
	SeatType    string      `json:"seat_type"`
	IsAvailable bool        `json:"is_available"`
	Status      *SeatStatus `json:"status,omitempty"` // Booking status from booking-service
}

type SeatStatus struct {
	SeatID   uuid.UUID `json:"seat_id"`
	IsBooked bool      `json:"is_booked"`
	IsLocked bool      `json:"is_locked"`
}

type PriceTier struct {
	SeatType string  `json:"seat_type"`
	Price    float64 `json:"price"`
}

// BookingResponse represents a booking from booking-service
type BookingResponse struct {
	ID          uuid.UUID        `json:"id"`
	Reference   string           `json:"reference"`
	TripID      uuid.UUID        `json:"trip_id"`
	UserID      *uuid.UUID       `json:"user_id,omitempty"`
	FullName    string           `json:"full_name"`
	Email       string           `json:"email"`
	Phone       string           `json:"phone"`
	TotalPrice  float64          `json:"total_price"`
	Status      string           `json:"status"` // pending, confirmed, cancelled
	BookedSeats []BookedSeat     `json:"booked_seats"`
	Passengers  []PassengerInfo  `json:"passengers"`
	ExpiresAt   *time.Time       `json:"expires_at,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	Transaction *TransactionInfo `json:"transaction,omitempty"`
}

type BookedSeat struct {
	SeatID     uuid.UUID `json:"seat_id"`
	SeatNumber string    `json:"seat_number"`
	Price      float64   `json:"price"`
}

type PassengerInfo struct {
	Name   string    `json:"name"`
	Phone  string    `json:"phone"`
	Email  string    `json:"email"`
	SeatID uuid.UUID `json:"seat_id"`
}

type TransactionInfo struct {
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	CheckoutURL string    `json:"checkout_url,omitempty"`
	QRCode      string    `json:"qr_code,omitempty"`
}

// CreateGuestBookingRequest for guest booking creation
type CreateGuestBookingRequest struct {
	TripID     uuid.UUID       `json:"trip_id"`
	SeatIDs    []uuid.UUID     `json:"seat_ids"`
	FullName   string          `json:"full_name"`
	Email      string          `json:"email"`
	Phone      string          `json:"phone"`
	Passengers []PassengerData `json:"passengers"`
}

type PassengerData struct {
	Name   string    `json:"name"`
	Phone  string    `json:"phone"`
	Email  string    `json:"email"`
	SeatID uuid.UUID `json:"seat_id"`
}

// TransactionResponse from payment-service
type TransactionResponse struct {
	ID            uuid.UUID `json:"id"`
	BookingID     uuid.UUID `json:"booking_id"`
	Amount        int       `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	CheckoutURL   string    `json:"checkout_url,omitempty"`
	QRCode        string    `json:"qr_code,omitempty"`
	OrderCode     int64     `json:"order_code,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateTransactionRequest for payment transaction creation
type CreateTransactionRequest struct {
	ID            uuid.UUID `json:"id,omitempty"` // Optional, payment service may generate
	BookingID     uuid.UUID `json:"booking_id"`
	Amount        int       `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Description   string    `json:"description,omitempty"`
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
}
