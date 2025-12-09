package trip

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID            uuid.UUID  `json:"id"`
	RouteID       uuid.UUID  `json:"route_id"`
	BusID         uuid.UUID  `json:"bus_id"`
	Date          string     `json:"date"` // stored as string "2006-01-02"
	DepartureTime time.Time  `json:"departure_time"`
	ArrivalTime   time.Time  `json:"arrival_time"`
	BasePrice     float64    `json:"base_price"`
	Status        TripStatus `json:"status"`
	IsActive      bool       `json:"is_active"`

	// Expansion fields
	Route *Route `json:"route,omitempty"`
	Bus   *Bus   `json:"bus,omitempty"`
}

type TripStatus string

const (
	TripStatusScheduled  TripStatus = "scheduled"
	TripStatusInProgress TripStatus = "in_progress"
	TripStatusCompleted  TripStatus = "completed"
	TripStatusCancelled  TripStatus = "cancelled"
	TripStatusDelayed    TripStatus = "delayed"
)

func (t *Trip) IsBookable() bool {
	return t.IsActive && t.Status == "scheduled"
}
