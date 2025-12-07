package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatStatus struct {
	BaseModel
	TripID        uuid.UUID  `gorm:"type:uuid;not null" json:"trip_id" validate:"required"`
	SeatID        uuid.UUID  `gorm:"type:uuid;not null" json:"seat_id" validate:"required"`
	SeatNumber    string     `gorm:"type:varchar(10);not null" json:"seat_number" validate:"required"`
	Status        string     `gorm:"type:varchar(50);not null;default:'available'" json:"status"`
	UserID        *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	BookingID     *uuid.UUID `gorm:"type:uuid" json:"booking_id,omitempty"`
	ReservedUntil *time.Time `gorm:"type:timestamptz" json:"reserved_until,omitempty"`

	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"booking,omitempty"`
}

func (SeatStatus) TableName() string { return "seat_statuses" }

func (ss *SeatStatus) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	return nil
}
