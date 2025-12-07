package model

import (
	"github.com/google/uuid"
)

type BookingSeat struct {
	BaseModel
	BookingID uuid.UUID `json:"booking_id" gorm:"type:uuid;not null;index"`
	SeatID    uuid.UUID `json:"seat_id" gorm:"type:uuid;not null"`

	// snapshot of seat info at booking time
	SeatNumber      string  `json:"seat_number" gorm:"type:varchar(10);not null"`
	SeatType        string  `json:"seat_type" gorm:"type:varchar(50);not null"`
	Floor           int     `json:"floor" gorm:"type:int;not null;default:1"`
	Price           float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	PriceMultiplier float64 `json:"price_multiplier" gorm:"type:decimal(3,2);not null;default:1.0"`
}

func (BookingSeat) TableName() string {
	return "booking_seats"
}
