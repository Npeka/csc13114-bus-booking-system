package trip

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID            uuid.UUID `json:"id"`
	RouteID       uuid.UUID `json:"route_id"`
	BusID         uuid.UUID `json:"bus_id"`
	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	BasePrice     float64   `json:"base_price"`
	Status        string    `json:"status"` // Raw string: scheduled, in_progress, completed, cancelled, delayed
	IsActive      bool      `json:"is_active"`
	Route         *Route    `json:"route,omitempty"`
	Bus           *Bus      `json:"bus,omitempty"`
}

// IsBookable returns true if the trip is active and scheduled
func (t *Trip) IsBookable() bool {
	return t.IsActive && t.Status == "scheduled"
}
