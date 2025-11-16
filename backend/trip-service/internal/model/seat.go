package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Seat struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BusID     uuid.UUID      `gorm:"type:uuid;not null" json:"bus_id" validate:"required"`
	SeatCode  string         `gorm:"type:varchar(10);not null" json:"seat_code" validate:"required"`
	SeatType  string         `gorm:"type:varchar(50);not null;default:'standard'" json:"seat_type" validate:"oneof=standard premium vip"`
	IsActive  bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Bus *Bus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bus,omitempty"`
}

func (Seat) TableName() string { return "seats" }

func (s *Seat) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
