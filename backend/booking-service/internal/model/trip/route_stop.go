package trip

import "github.com/google/uuid"

type RouteStop struct {
	ID            uuid.UUID `json:"id"`
	RouteID       uuid.UUID `json:"route_id"`
	StopOrder     int       `json:"stop_order"`
	StopType      StopType  `json:"stop_type"`
	Location      string    `json:"location"`
	Address       string    `json:"address"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	OffsetMinutes int       `json:"offset_minutes"`
	IsActive      bool      `json:"is_active"`
}

type StopType struct {
	Value       string `json:"value"`
	DisplayName string `json:"display_name"`
}
