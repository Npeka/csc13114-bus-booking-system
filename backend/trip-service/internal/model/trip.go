package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Trip struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	RouteID       uuid.UUID      `gorm:"type:uuid;not null" json:"route_id" validate:"required"`
	BusID         uuid.UUID      `gorm:"type:uuid;not null" json:"bus_id" validate:"required"`
	DepartureTime time.Time      `gorm:"type:timestamptz;not null" json:"departure_time" validate:"required"`
	ArrivalTime   time.Time      `gorm:"type:timestamptz;not null" json:"arrival_time" validate:"required"`
	BasePrice     float64        `gorm:"type:decimal(10,2);not null" json:"base_price" validate:"required,min=0"`
	Status        string         `gorm:"type:varchar(50);not null;default:'scheduled'" json:"status" validate:"oneof=scheduled in_progress completed cancelled"`
	IsActive      bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Route *Route `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"route,omitempty"`
	Bus   *Bus   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bus,omitempty"`
}

func (Trip) TableName() string { return "trips" }

func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
