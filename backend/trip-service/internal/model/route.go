package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Route struct {
	BaseModel
	Origin           string  `gorm:"type:varchar(255);not null" json:"origin" validate:"required"`
	Destination      string  `gorm:"type:varchar(255);not null" json:"destination" validate:"required"`
	DistanceKm       float64 `gorm:"type:decimal(10,2);not null" json:"distance_km" validate:"required,min=1"`
	EstimatedMinutes int     `gorm:"type:integer;not null" json:"estimated_minutes" validate:"required,min=1"`
	IsActive         bool    `gorm:"type:boolean;not null;default:true" json:"is_active"`

	Trips      []Trip      `gorm:"foreignKey:RouteID" json:"trips,omitempty"`
	RouteStops []RouteStop `gorm:"foreignKey:RouteID" json:"route_stops,omitempty"`
}

func (Route) TableName() string {
	return "routes"
}

func (r *Route) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type RouteResponse struct {
	ID               uuid.UUID           `json:"id"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	Origin           string              `json:"origin"`
	Destination      string              `json:"destination"`
	DistanceKm       float64             `json:"distance_km"`
	EstimatedMinutes int                 `json:"estimated_minutes"`
	IsActive         bool                `json:"is_active"`
	RouteStops       []RouteStopResponse `json:"route_stops,omitempty"`
}

type RouteListResponse struct {
	Routes     []RouteSummary `json:"routes"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type RouteSummary struct {
	ID               uuid.UUID `json:"id"`
	Origin           string    `json:"origin"`
	Destination      string    `json:"destination"`
	DistanceKm       float64   `json:"distance_km"`
	EstimatedMinutes int       `json:"estimated_minutes"`
	ActiveTrips      int       `json:"active_trips"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateRouteRequest struct {
	Origin           string                   `json:"origin" validate:"required,min=2,max=255"`
	Destination      string                   `json:"destination" validate:"required,min=2,max=255"`
	DistanceKm       float64                  `json:"distance_km" validate:"required,min=1"`
	EstimatedMinutes int                      `json:"estimated_minutes" validate:"required,min=1"`
	RouteStops       []CreateRouteStopRequest `json:"route_stops" validate:"required,dive"`
}

type UpdateRouteRequest struct {
	Origin           *string  `json:"origin,omitempty" validate:"omitempty,min=2,max=255"`
	Destination      *string  `json:"destination,omitempty" validate:"omitempty,min=2,max=255"`
	DistanceKm       *float64 `json:"distance_km,omitempty" validate:"omitempty,min=1"`
	EstimatedMinutes *int     `json:"estimated_minutes,omitempty" validate:"omitempty,min=1"`
	IsActive         *bool    `json:"is_active,omitempty"`
}
