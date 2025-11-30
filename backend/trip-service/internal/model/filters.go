package model

import (
	"time"

	"bus-booking/trip-service/internal/constants"

	"github.com/google/uuid"
)

// TripSearchRequest represents advanced search criteria for trips

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
