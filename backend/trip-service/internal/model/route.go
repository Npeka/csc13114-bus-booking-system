package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Route struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OperatorID       uuid.UUID      `gorm:"type:uuid;not null" json:"operator_id" validate:"required"`
	Origin           string         `gorm:"type:varchar(255);not null" json:"origin" validate:"required"`
	Destination      string         `gorm:"type:varchar(255);not null" json:"destination" validate:"required"`
	DistanceKm       int            `gorm:"type:integer;not null" json:"distance_km" validate:"required,min=1"`
	EstimatedMinutes int            `gorm:"type:integer;not null" json:"estimated_minutes" validate:"required,min=1"`
	IsActive         bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt        time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	Operator *Operator `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"operator,omitempty"`
	Trips    []Trip    `gorm:"foreignKey:RouteID" json:"trips,omitempty"`
}

func (Route) TableName() string { return "routes" }

func (r *Route) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
