package model

import (
	"time"

	"github.com/google/uuid"
)

// PaginationRequest represents common pagination query parameters
type PaginationRequest struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=20"`
}

// ListBusesRequest represents query parameters for listing buses
type ListBusesRequest struct {
	PaginationRequest
	OperatorID string `form:"operator_id"`
}

// ListRoutesRequest represents query parameters for listing routes
type ListRoutesRequest struct {
	PaginationRequest
	OperatorID string `form:"operator_id"`
}

// SearchRoutesQueryRequest represents query parameters for searching routes
type SearchRoutesQueryRequest struct {
	Origin      string `form:"origin" binding:"required"`
	Destination string `form:"destination" binding:"required"`
}

// ListTripsByRouteRequest represents query parameters for listing trips by route
type ListTripsByRouteRequest struct {
	Date string `form:"date" binding:"required"`
}

// SeatAvailabilityRequest represents seat availability check request
type SeatAvailabilityRequest struct {
	TripID uuid.UUID `json:"trip_id" validate:"required"`
}

// OperatorListResponse represents operator listing
type OperatorListResponse struct {
	Operators  []OperatorSummary `json:"operators"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// OperatorSummary represents summary operator information
type OperatorSummary struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	Status       string    `json:"status"`
	ActiveRoutes int       `json:"active_routes"`
	ActiveBuses  int       `json:"active_buses"`
	CreatedAt    time.Time `json:"created_at"`
}

// RouteListResponse represents route listing
type RouteListResponse struct {
	Routes     []RouteSummary `json:"routes"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// RouteSummary represents summary route information
type RouteSummary struct {
	ID               uuid.UUID `json:"id"`
	Origin           string    `json:"origin"`
	Destination      string    `json:"destination"`
	DistanceKm       int       `json:"distance_km"`
	EstimatedMinutes int       `json:"estimated_minutes"`
	OperatorName     string    `json:"operator_name"`
	ActiveTrips      int       `json:"active_trips"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

// CreateOperatorRequest represents operator creation request
type CreateOperatorRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=255"`
	ContactEmail string `json:"contact_email" validate:"required,email"`
	ContactPhone string `json:"contact_phone,omitempty" validate:"omitempty,min=10,max=20"`
}

// UpdateOperatorRequest represents operator update request
type UpdateOperatorRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	ContactEmail *string `json:"contact_email,omitempty" validate:"omitempty,email"`
	ContactPhone *string `json:"contact_phone,omitempty" validate:"omitempty,min=10,max=20"`
}

// CreateRouteRequest represents route creation request
type CreateRouteRequest struct {
	OperatorID       uuid.UUID `json:"operator_id" validate:"required"`
	Origin           string    `json:"origin" validate:"required,min=2,max=255"`
	Destination      string    `json:"destination" validate:"required,min=2,max=255"`
	DistanceKm       int       `json:"distance_km" validate:"required,min=1"`
	EstimatedMinutes int       `json:"estimated_minutes" validate:"required,min=1"`
}

// UpdateRouteRequest represents route update request
type UpdateRouteRequest struct {
	Origin           *string `json:"origin,omitempty" validate:"omitempty,min=2,max=255"`
	Destination      *string `json:"destination,omitempty" validate:"omitempty,min=2,max=255"`
	DistanceKm       *int    `json:"distance_km,omitempty" validate:"omitempty,min=1"`
	EstimatedMinutes *int    `json:"estimated_minutes,omitempty" validate:"omitempty,min=1"`
	IsActive         *bool   `json:"is_active,omitempty"`
}

// CreateBusRequest represents bus creation request
type CreateBusRequest struct {
	OperatorID   uuid.UUID `json:"operator_id" validate:"required"`
	PlateNumber  string    `json:"plate_number" validate:"required,min=3,max=20"`
	Model        string    `json:"model" validate:"required,min=2,max=255"`
	SeatCapacity int       `json:"seat_capacity" validate:"required,min=1,max=100"`
	Amenities    []string  `json:"amenities,omitempty"`
}

// UpdateBusRequest represents bus update request
type UpdateBusRequest struct {
	PlateNumber  *string   `json:"plate_number,omitempty" validate:"omitempty,min=3,max=20"`
	Model        *string   `json:"model,omitempty" validate:"omitempty,min=2,max=255"`
	SeatCapacity *int      `json:"seat_capacity,omitempty" validate:"omitempty,min=1,max=100"`
	Amenities    *[]string `json:"amenities,omitempty"`
	IsActive     *bool     `json:"is_active,omitempty"`
}

// CreateTripRequest represents trip creation request
type CreateTripRequest struct {
	RouteID       uuid.UUID `json:"route_id" validate:"required"`
	BusID         uuid.UUID `json:"bus_id" validate:"required"`
	DepartureTime time.Time `json:"departure_time" validate:"required"`
	ArrivalTime   time.Time `json:"arrival_time" validate:"required"`
	BasePrice     float64   `json:"base_price" validate:"required,min=0"`
}

// UpdateTripRequest represents trip update request
type UpdateTripRequest struct {
	DepartureTime *time.Time `json:"departure_time,omitempty" validate:"omitempty"`
	ArrivalTime   *time.Time `json:"arrival_time,omitempty" validate:"omitempty"`
	BasePrice     *float64   `json:"base_price,omitempty" validate:"omitempty,min=0"`
	Status        *string    `json:"status,omitempty" validate:"omitempty,oneof=scheduled in_progress completed cancelled"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

// Booking requests
type CreateBookingRequest struct {
	TripID     uuid.UUID          `json:"trip_id" validate:"required"`
	SeatIDs    []uuid.UUID        `json:"seat_ids" validate:"required,min=1,max=10"`
	Passengers []PassengerRequest `json:"passengers" validate:"required,dive"`
	IsGuest    bool               `json:"is_guest"`
	GuestInfo  *GuestInfo         `json:"guest_info,omitempty"`
}

type PassengerRequest struct {
	SeatID      uuid.UUID `json:"seat_id" validate:"required"`
	FullName    string    `json:"full_name" validate:"required,min=2,max=255"`
	IDNumber    string    `json:"id_number,omitempty" validate:"omitempty,min=5,max=50"`
	PhoneNumber string    `json:"phone_number,omitempty" validate:"omitempty,min=10,max=20"`
}

type GuestInfo struct {
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,min=10,max=20"`
	Name  string `json:"name" validate:"required,min=2,max=255"`
}

type LockSeatsRequest struct {
	TripID    uuid.UUID   `json:"trip_id" validate:"required"`
	SeatIDs   []uuid.UUID `json:"seat_ids" validate:"required,min=1,max=10"`
	SessionID string      `json:"session_id" validate:"required"`
}

type UnlockSeatsRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}

type CancelBookingRequest struct {
	Reason string `json:"reason,omitempty"`
}

// Booking query parameters
type BookingLookupRequest struct {
	Reference string `form:"reference" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
}

type ListBookingsRequest struct {
	PaginationRequest
	Status string `form:"status"`
}
