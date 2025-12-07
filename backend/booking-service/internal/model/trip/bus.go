package trip

import "github.com/google/uuid"

type Bus struct {
	ID           uuid.UUID `json:"id"`
	PlateNumber  string    `json:"plate_number"`
	Model        string    `json:"model"`
	BusType      string    `json:"bus_type"` // Raw string: standard, vip, sleeper, double_decker
	SeatCapacity int       `json:"seat_capacity"`
	Amenities    []string  `json:"amenities,omitempty"` // Raw strings: wifi, ac, toilet, etc.
	IsActive     bool      `json:"is_active"`
	Seats        []Seat    `json:"seats,omitempty"`
}
