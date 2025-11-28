package model

import (
	"github.com/google/uuid"
)

type Passenger struct {
	BaseModel
	BookingID   uuid.UUID `json:"booking_id" gorm:"type:uuid;not null"`
	SeatID      uuid.UUID `json:"seat_id" gorm:"type:uuid;not null"`
	FullName    string    `json:"full_name" gorm:"type:varchar(255);not null"`
	IDNumber    string    `json:"id_number,omitempty" gorm:"type:varchar(50)"`
	PhoneNumber string    `json:"phone_number,omitempty" gorm:"type:varchar(20)"`
	SeatNumber  string    `json:"seat_number" gorm:"type:varchar(10);not null"`
	SeatType    string    `json:"seat_type" gorm:"type:varchar(20);not null"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
}

func (Passenger) TableName() string { return "passengers" }
