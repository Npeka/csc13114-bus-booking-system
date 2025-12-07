package trip

import (
	"github.com/google/uuid"
)

type Seat struct {
	ID              uuid.UUID `json:"id"`
	BusID           uuid.UUID `json:"bus_id"`
	SeatNumber      string    `json:"seat_number"`
	Row             int       `json:"row"`
	Column          int       `json:"column"`
	SeatType        string    `json:"seat_type"`
	PriceMultiplier float64   `json:"price_multiplier"`
	Floor           int       `json:"floor"`
}

func (s *Seat) CalculateSeatPrice(basePrice float64) float64 {
	return basePrice * s.PriceMultiplier
}
