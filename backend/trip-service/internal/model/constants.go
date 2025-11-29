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
	Bus   BusConstants   `json:"bus"`
	Route RouteConstants `json:"route"`
	Trip  TripConstants  `json:"trip"`
}
