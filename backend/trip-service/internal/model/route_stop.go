package model

import (
	"time"

	"bus-booking/trip-service/internal/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteStop struct {
	ID            uuid.UUID          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	RouteID       uuid.UUID          `gorm:"type:uuid;not null;index:idx_route_stops_route" json:"route_id" validate:"required"`
	StopOrder     int                `gorm:"type:integer;not null" json:"stop_order" validate:"required,min=1"`
	StopType      constants.StopType `gorm:"type:varchar(20);not null" json:"stop_type" validate:"required"`
	Location      string             `gorm:"type:varchar(255);not null;index:idx_route_stops_location" json:"location" validate:"required"`
	Address       string             `gorm:"type:text" json:"address"`
	Latitude      *float64           `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude     *float64           `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	OffsetMinutes int                `gorm:"type:integer;not null" json:"offset_minutes" validate:"min=0"`
	IsActive      bool               `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt     time.Time          `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time          `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt     `gorm:"index" json:"-"`

	Route *Route `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"route,omitempty"`
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

// CreateRouteStopRequest is the request body for creating a route stop
type CreateRouteStopRequest struct {
	RouteID       uuid.UUID          `json:"route_id" validate:"required"`
	StopOrder     int                `json:"stop_order" validate:"required,min=1"`
	StopType      constants.StopType `json:"stop_type" validate:"required"`
	Location      string             `json:"location" validate:"required"`
	Address       string             `json:"address"`
	Latitude      *float64           `json:"latitude,omitempty"`
	Longitude     *float64           `json:"longitude,omitempty"`
	OffsetMinutes int                `json:"offset_minutes" validate:"min=0"`
}

// UpdateRouteStopRequest is the request body for updating a route stop
type UpdateRouteStopRequest struct {
	StopOrder     *int                `json:"stop_order,omitempty" validate:"omitempty,min=1"`
	StopType      *constants.StopType `json:"stop_type,omitempty"`
	Location      *string             `json:"location,omitempty"`
	Address       *string             `json:"address,omitempty"`
	Latitude      *float64            `json:"latitude,omitempty"`
	Longitude     *float64            `json:"longitude,omitempty"`
	OffsetMinutes *int                `json:"offset_minutes,omitempty" validate:"omitempty,min=0"`
	IsActive      *bool               `json:"is_active,omitempty"`
}
