package model

import (
	"time"

	"github.com/google/uuid"
)

// TripSearchRequest represents trip search parameters
type TripSearchRequest struct {
	Origin        string    `json:"origin" validate:"required"`
	Destination   string    `json:"destination" validate:"required"`
	DepartureDate time.Time `json:"departure_date" validate:"required"`
	Passengers    int       `json:"passengers" validate:"required,min=1,max=10"`
	SeatType      string    `json:"seat_type,omitempty" validate:"omitempty,oneof=standard premium vip"`
	PriceMin      *float64  `json:"price_min,omitempty" validate:"omitempty,min=0"`
	PriceMax      *float64  `json:"price_max,omitempty" validate:"omitempty,min=0"`
	OperatorID    *uuid.UUID `json:"operator_id,omitempty"`
	SortBy        string    `json:"sort_by,omitempty" validate:"omitempty,oneof=price departure_time arrival_time"`
	SortOrder     string    `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	Page          int       `json:"page,omitempty" validate:"omitempty,min=1"`
	Limit         int       `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
}

// TripSearchResponse represents trip search results
type TripSearchResponse struct {
	Trips      []TripDetail `json:"trips"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int          `json:"total_pages"`
}

// TripDetail represents detailed trip information with availability
type TripDetail struct {
	ID                uuid.UUID `json:"id"`
	RouteID           uuid.UUID `json:"route_id"`
	BusID             uuid.UUID `json:"bus_id"`
	DepartureTime     time.Time `json:"departure_time"`
	ArrivalTime       time.Time `json:"arrival_time"`
	BasePrice         float64   `json:"base_price"`
	Status            string    `json:"status"`
	AvailableSeats    int       `json:"available_seats"`
	TotalSeats        int       `json:"total_seats"`
	Duration          string    `json:"duration"`
	
	// Route information
	Origin            string    `json:"origin"`
	Destination       string    `json:"destination"`
	DistanceKm        int       `json:"distance_km"`
	
	// Bus information
	BusModel          string    `json:"bus_model"`
	BusPlateNumber    string    `json:"bus_plate_number"`
	BusAmenities      []string  `json:"bus_amenities"`
	
	// Operator information
	OperatorID        uuid.UUID `json:"operator_id"`
	OperatorName      string    `json:"operator_name"`
}

// SeatAvailabilityRequest represents seat availability check request
type SeatAvailabilityRequest struct {
	TripID uuid.UUID `json:"trip_id" validate:"required"`
}

// SeatAvailabilityResponse represents seat availability with detailed seat info
type SeatAvailabilityResponse struct {
	TripID         uuid.UUID    `json:"trip_id"`
	AvailableSeats int          `json:"available_seats"`
	TotalSeats     int          `json:"total_seats"`
	Seats          []SeatDetail `json:"seats"`
}

// SeatDetail represents detailed seat information with availability
type SeatDetail struct {
	ID        uuid.UUID `json:"id"`
	SeatCode  string    `json:"seat_code"`
	SeatType  string    `json:"seat_type"`
	IsBooked  bool      `json:"is_booked"`
	IsLocked  bool      `json:"is_locked"`
	Price     float64   `json:"price"`
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