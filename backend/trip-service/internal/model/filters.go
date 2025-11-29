package model

import (
	"time"

	"bus-booking/trip-service/internal/constants"

	"github.com/google/uuid"
)

// TripSearchRequest represents advanced search criteria for trips
type TripSearchRequest struct {
	// Basic search
	Origin      string    `json:"origin" validate:"required"`
	Destination string    `json:"destination" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`

	// Advanced filters
	DepartureTimeStart *time.Time           `json:"departure_time_start,omitempty"`
	DepartureTimeEnd   *time.Time           `json:"departure_time_end,omitempty"`
	MinPrice           *float64             `json:"min_price,omitempty" validate:"omitempty,min=0"`
	MaxPrice           *float64             `json:"max_price,omitempty" validate:"omitempty,min=0"`
	SeatTypes          []constants.SeatType `json:"seat_types,omitempty"`
	Amenities          []constants.Amenity  `json:"amenities,omitempty"`

	// Sorting
	SortBy    string `json:"sort_by" validate:"omitempty,oneof=price departure_time duration"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`

	// Pagination
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// TripSearchResponse represents the search results with metadata
type TripSearchResponse struct {
	Trips      []TripDetail   `json:"trips"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
	Filters    *SearchFilters `json:"filters,omitempty"`
}

// SearchFilters provides aggregated filter options
type SearchFilters struct {
	PriceRange     PriceRange     `json:"price_range"`
	AvailableTimes []TimeSlot     `json:"available_times"`
	SeatTypes      []SeatTypeInfo `json:"seat_types"`
}

// PriceRange represents min and max prices in search results
type PriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// TimeSlot represents a time range with trip count
type TimeSlot struct {
	Label string    `json:"label"` // "Morning", "Afternoon", "Evening", "Night"
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Count int       `json:"count"`
}

// SeatTypeInfo represents seat type availability
type SeatTypeInfo struct {
	Type  constants.SeatType `json:"type"`
	Count int                `json:"count"`
}

// TripDetail represents detailed trip information in search results
type TripDetail struct {
	ID             uuid.UUID            `json:"id"`
	RouteID        uuid.UUID            `json:"route_id"`
	BusID          uuid.UUID            `json:"bus_id"`
	DepartureTime  time.Time            `json:"departure_time"`
	ArrivalTime    time.Time            `json:"arrival_time"`
	BasePrice      float64              `json:"base_price"`
	Status         constants.TripStatus `json:"status"`
	AvailableSeats int                  `json:"available_seats"`
	TotalSeats     int                  `json:"total_seats"`

	// Route info
	Origin          string `json:"origin"`
	Destination     string `json:"destination"`
	DistanceKm      int    `json:"distance_km"`
	DurationMinutes int    `json:"duration_minutes"`

	// Bus info
	BusModel  string              `json:"bus_model"`
	BusType   constants.BusType   `json:"bus_type,omitempty"`
	Amenities []constants.Amenity `json:"amenities,omitempty"`

	// Pricing tiers
	PriceTiers []PriceTier `json:"price_tiers,omitempty"`
}

// PriceTier represents pricing for different seat types
type PriceTier struct {
	SeatType        constants.SeatType `json:"seat_type"`
	BasePrice       float64            `json:"base_price"`
	PriceMultiplier float64            `json:"price_multiplier"`
	FinalPrice      float64            `json:"final_price"`
	AvailableCount  int                `json:"available_count"`
}

// SeatAvailabilityResponse represents detailed seat availability
type SeatAvailabilityResponse struct {
	TripID         uuid.UUID          `json:"trip_id"`
	TotalSeats     int                `json:"total_seats"`
	AvailableSeats int                `json:"available_seats"`
	BookedSeats    int                `json:"booked_seats"`
	SeatMap        []SeatAvailability `json:"seat_map"`
	PriceTiers     []PriceTier        `json:"price_tiers"`
}

// SeatAvailability represents individual seat availability
type SeatAvailability struct {
	SeatID      uuid.UUID          `json:"seat_id"`
	SeatNumber  string             `json:"seat_number"`
	SeatType    constants.SeatType `json:"seat_type"`
	Price       float64            `json:"price"`
	IsAvailable bool               `json:"is_available"`
	Row         int                `json:"row"`
	Column      int                `json:"column"`
	Floor       int                `json:"floor"`
}
