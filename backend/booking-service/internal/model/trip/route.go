package trip

import (
	"time"

	"github.com/google/uuid"
)

type Route struct {
	ID               uuid.UUID   `json:"id"`
	Origin           string      `json:"origin"`
	Destination      string      `json:"destination"`
	DistanceKm       float64     `json:"distance_km"`
	EstimatedMinutes int         `json:"estimated_minutes"`
	IsActive         bool        `json:"is_active"`
	RouteStops       []RouteStop `json:"route_stops,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}
