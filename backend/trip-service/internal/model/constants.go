package model

// SeatTypeConstant represents a seat type constant with details
type SeatTypeConstant struct {
	Value           string  `json:"value"`
	DisplayName     string  `json:"display_name"`
	PriceMultiplier float64 `json:"price_multiplier"`
}

// AmenityConstant represents an amenity constant with details
type AmenityConstant struct {
	Value       string `json:"value"`
	DisplayName string `json:"display_name"`
}

// BusTypeConstant represents a bus type constant
type BusTypeConstant struct {
	Value       string `json:"value"`
	DisplayName string `json:"display_name"`
}

// StopTypeConstant represents a stop type constant
type StopTypeConstant struct {
	Value       string `json:"value"`
	DisplayName string `json:"display_name"`
}

// TripStatusConstant represents a trip status constant
type TripStatusConstant struct {
	Value       string `json:"value"`
	DisplayName string `json:"display_name"`
}

// ConstantDisplay represents a constant value with its display name
type ConstantDisplay struct {
	Value       string `json:"value"`
	DisplayName string `json:"display_name"`
}

// FilterPriceRange represents a price range filter option
type FilterPriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// FilterTimeSlot represents a time slot filter option
type FilterTimeSlot struct {
	StartTime   string `json:"start_time"` // HH:MM format
	EndTime     string `json:"end_time"`   // HH:MM format
	DisplayName string `json:"display_name"`
}

// SearchFilterConstants contains all filter options for trip search
type SearchFilterConstants struct {
	SortOptions []ConstantDisplay  `json:"sort_options"`
	PriceRanges []FilterPriceRange `json:"price_ranges"`
	TimeSlots   []FilterTimeSlot   `json:"time_slots"`
	SeatTypes   []SeatTypeConstant `json:"seat_types"`
	Amenities   []AmenityConstant  `json:"amenities"`
	Cities      []string           `json:"cities"`
}

// BusConstants contains constants related to buses and seats
type BusConstants struct {
	SeatTypes []SeatTypeConstant `json:"seat_types"`
	Amenities []AmenityConstant  `json:"amenities"`
	BusTypes  []BusTypeConstant  `json:"bus_types"`
}

// RouteConstants contains constants related to routes and stops
type RouteConstants struct {
	StopTypes []StopTypeConstant `json:"stop_types"`
}

// TripConstants contains constants related to trips
type TripConstants struct {
	TripStatuses []TripStatusConstant `json:"trip_statuses"`
}

// ConstantsResponse contains all constants grouped by domain
type ConstantsResponse struct {
	Bus           BusConstants          `json:"bus"`
	Route         RouteConstants        `json:"route"`
	Trip          TripConstants         `json:"trip"`
	SearchFilters SearchFilterConstants `json:"search_filters"`
}
