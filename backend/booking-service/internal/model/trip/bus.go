package trip

import "github.com/google/uuid"

type Bus struct {
	ID           uuid.UUID `json:"id"`
	PlateNumber  string    `json:"plate_number"`
	Model        string    `json:"model"`
	BusType      BusType   `json:"bus_type"` // Raw string: standard, vip, sleeper, double_decker
	SeatCapacity int       `json:"seat_capacity"`
	Amenities    []Amenity `json:"amenities,omitempty"` // Raw strings: wifi, ac, toilet, etc.
	IsActive     bool      `json:"is_active"`
	Seats        []Seat    `json:"seats,omitempty"`
}

type BusType string

const (
	BusTypeStandard     BusType = "standard"
	BusTypeVIP          BusType = "vip"
	BusTypeSleeper      BusType = "sleeper"
	BusTypeDoubleDecker BusType = "double_decker"
)

type Amenity string

const (
	AmenityWiFi     Amenity = "wifi"
	AmenityAC       Amenity = "ac"
	AmenityToilet   Amenity = "toilet"
	AmenityTV       Amenity = "tv"
	AmenityCharging Amenity = "charging"
	AmenityBlanket  Amenity = "blanket"
	AmenityWater    Amenity = "water"
	AmenitySnacks   Amenity = "snacks"
)
