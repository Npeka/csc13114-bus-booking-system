package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingSeat struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID uuid.UUID `gorm:"type:uuid;not null" json:"booking_id" validate:"required"`
	SeatID    uuid.UUID `gorm:"type:uuid;not null" json:"seat_id" validate:"required"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,min=0"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"booking,omitempty"`
}

func (BookingSeat) TableName() string { return "booking_seats" }

func (bs *BookingSeat) BeforeCreate(tx *gorm.DB) error {
	if bs.ID == uuid.Nil {
		bs.ID = uuid.New()
	}
	return nil
}
