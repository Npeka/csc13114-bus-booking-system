package model

import (
	"bus-booking/trip-service/internal/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteStop struct {
	BaseModel
	RouteID       uuid.UUID          `gorm:"type:uuid;not null;index" json:"route_id"`
	StopOrder     int                `gorm:"type:integer;not null" json:"stop_order"`
	StopType      constants.StopType `gorm:"type:varchar(50);not null" json:"stop_type"` // pickup, dropoff, both
	Location      string             `gorm:"type:varchar(255);not null" json:"location"`
	Address       string             `gorm:"type:text;not null" json:"address"`
	Latitude      *float64           `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude     *float64           `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	OffsetMinutes int                `gorm:"type:integer;not null" json:"offset_minutes"`
	IsActive      bool               `gorm:"type:boolean;not null;default:true" json:"is_active"`
	Route         Route              `gorm:"foreignKey:RouteID" json:"route,omitempty"`
}

func (RouteStop) TableName() string {
	return "route_stops"
}

func (rs *RouteStop) BeforeCreate(tx *gorm.DB) error {
	if rs.ID == uuid.Nil {
		rs.ID = uuid.New()
	}
	return nil
}

// Request models
type CreateRouteStopRequest struct {
	RouteID       uuid.UUID          `json:"route_id,omitempty"` // Optional - auto-assigned when creating route with stops
	StopOrder     int                `json:"stop_order" validate:"min=1"`
	StopType      constants.StopType `json:"stop_type" validate:"required,oneof=pickup dropoff both"`
	Location      string             `json:"location" validate:"required,min=2,max=255"`
	Address       string             `json:"address" validate:"required"`
	Latitude      *float64           `json:"latitude,omitempty" validate:"omitempty,min=-90,max=90"`
	Longitude     *float64           `json:"longitude,omitempty" validate:"omitempty,min=-180,max=180"`
	OffsetMinutes int                `json:"offset_minutes" validate:"min=0"`
}

type UpdateRouteStopRequest struct {
	StopOrder     *int                `json:"stop_order,omitempty" validate:"omitempty,min=1"`
	StopType      *constants.StopType `json:"stop_type,omitempty" validate:"omitempty,oneof=pickup dropoff both"`
	Location      *string             `json:"location,omitempty" validate:"omitempty,min=2,max=255"`
	Address       *string             `json:"address,omitempty"`
	Latitude      *float64            `json:"latitude,omitempty" validate:"omitempty,min=-90,max=90"`
	Longitude     *float64            `json:"longitude,omitempty" validate:"omitempty,min=-180,max=180"`
	OffsetMinutes *int                `json:"offset_minutes,omitempty" validate:"omitempty,min=0"`
	IsActive      *bool               `json:"is_active,omitempty"`
}

// MoveRouteStopRequest allows frontend to specify position relative to another stop
type MoveRouteStopRequest struct {
	Position        string     `json:"position" validate:"required,oneof=before after first last"` // before, after, first, last
	ReferenceStopID *uuid.UUID `json:"reference_stop_id,omitempty"`                                // Required for before/after
}

type RouteStopResponse struct {
	ID            uuid.UUID          `json:"id"`
	RouteID       uuid.UUID          `json:"route_id"`
	StopOrder     int                `json:"stop_order"`
	StopType      constants.StopType `json:"stop_type"` // Raw string value
	Location      string             `json:"location"`
	Address       string             `json:"address"`
	Latitude      *float64           `json:"latitude,omitempty"`
	Longitude     *float64           `json:"longitude,omitempty"`
	OffsetMinutes int                `json:"offset_minutes"`
	IsActive      bool               `json:"is_active"`
}
