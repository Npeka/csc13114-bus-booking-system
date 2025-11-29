package model

import (
	"bus-booking/trip-service/internal/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Seat struct {
	BaseModel
	BusID           uuid.UUID          `gorm:"type:uuid;not null;index:idx_seats_bus" json:"bus_id" validate:"required"`
	SeatNumber      string             `gorm:"type:varchar(10);not null" json:"seat_number" validate:"required"`
	Row             int                `gorm:"type:integer;not null" json:"row" validate:"required,min=1"`
	Column          int                `gorm:"type:integer;not null" json:"column" validate:"required,min=1"`
	SeatType        constants.SeatType `gorm:"type:varchar(20);not null" json:"seat_type" validate:"required"`
	PriceMultiplier float64            `gorm:"type:decimal(3,2);not null;default:1.0" json:"price_multiplier" validate:"min=0.5,max=5.0"`
	IsAvailable     bool               `gorm:"type:boolean;not null;default:true" json:"is_available"`
	Floor           int                `gorm:"type:integer;not null;default:1" json:"floor" validate:"min=1,max=2"`

	Bus *Bus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bus,omitempty"`
}

func (Seat) TableName() string {
	return "seats"
}

func (s *Seat) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	// Set default price multiplier based on seat type if not set
	if s.PriceMultiplier == 0 {
		s.PriceMultiplier = s.SeatType.GetPriceMultiplier()
	}
	return nil
}

type CreateSeatRequest struct {
	BusID           uuid.UUID          `json:"bus_id" validate:"required"`
	SeatNumber      string             `json:"seat_number" validate:"required"`
	Row             int                `json:"row" validate:"required,min=1"`
	Column          int                `json:"column" validate:"required,min=1"`
	SeatType        constants.SeatType `json:"seat_type" validate:"required"`
	PriceMultiplier *float64           `json:"price_multiplier,omitempty" validate:"omitempty,min=0.5,max=5.0"`
	Floor           int                `json:"floor" validate:"min=1,max=2"`
}

type UpdateSeatRequest struct {
	SeatNumber      *string             `json:"seat_number,omitempty"`
	Row             *int                `json:"row,omitempty" validate:"omitempty,min=1"`
	Column          *int                `json:"column,omitempty" validate:"omitempty,min=1"`
	SeatType        *constants.SeatType `json:"seat_type,omitempty"`
	PriceMultiplier *float64            `json:"price_multiplier,omitempty" validate:"omitempty,min=0.5,max=5.0"`
	IsAvailable     *bool               `json:"is_available,omitempty"`
	Floor           *int                `json:"floor,omitempty" validate:"omitempty,min=1,max=2"`
}

type BulkCreateSeatsRequest struct {
	BusID uuid.UUID           `json:"bus_id" validate:"required"`
	Seats []CreateSeatRequest `json:"seats" validate:"required,min=1,dive"`
}

type SeatMapResponse struct {
	BusID      uuid.UUID      `json:"bus_id"`
	TotalSeats int            `json:"total_seats"`
	Seats      []SeatDetail   `json:"seats"`
	Layout     SeatLayoutInfo `json:"layout"`
}

type SeatDetail struct {
	ID              uuid.UUID          `json:"id"`
	SeatNumber      string             `json:"seat_number"`
	Row             int                `json:"row"`
	Column          int                `json:"column"`
	SeatType        constants.SeatType `json:"seat_type"`
	PriceMultiplier float64            `json:"price_multiplier"`
	IsAvailable     bool               `json:"is_available"`
	Floor           int                `json:"floor"`
}

type SeatLayoutInfo struct {
	MaxRows    int `json:"max_rows"`
	MaxColumns int `json:"max_columns"`
	Floors     int `json:"floors"`
}

type LockSeatsRequest struct {
	TripID    uuid.UUID   `json:"trip_id" validate:"required"`
	SeatIDs   []uuid.UUID `json:"seat_ids" validate:"required,min=1,max=10"`
	SessionID string      `json:"session_id" validate:"required"`
}
