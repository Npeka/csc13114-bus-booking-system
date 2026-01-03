package model

import (
	"bus-booking/trip-service/internal/constants"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Bus struct {
	BaseModel
	PlateNumber  string            `gorm:"type:varchar(20);unique;not null" json:"plate_number" validate:"required"`
	Model        string            `gorm:"type:varchar(255);not null" json:"model" validate:"required"`
	BusType      constants.BusType `gorm:"type:varchar(20);not null;default:'standard'" json:"bus_type" validate:"required"` // constants.BusType values
	SeatCapacity int               `gorm:"type:integer;not null" json:"seat_capacity" validate:"required,min=1,max=100"`
	Amenities    pq.StringArray    `gorm:"type:text[]" json:"amenities"` // constants.Amenity values
	ImageURLs    pq.StringArray    `gorm:"type:text[];column:image_urls" json:"image_urls"`
	IsActive     bool              `gorm:"type:boolean;not null;default:true" json:"is_active"`

	Seats []Seat `gorm:"foreignKey:BusID" json:"seats"`
}

func (Bus) TableName() string {
	return "buses"
}

func (b *Bus) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BusResponse struct {
	ID           uuid.UUID         `json:"id"`
	PlateNumber  string            `json:"plate_number"`
	Model        string            `json:"model"`
	BusType      constants.BusType `json:"bus_type"` // Raw string value: constants.BusType
	SeatCapacity int               `json:"seat_capacity"`
	Amenities    []string          `json:"amenities"` // Raw string values: constants.Amenity
	ImageURLs    []string          `json:"image_urls"`
	IsActive     bool              `json:"is_active"`

	Seats []SeatResponse `json:"seats,omitempty"`
}

type CreateBusRequest struct {
	PlateNumber string            `json:"plate_number" validate:"required,min=3,max=20"`
	Model       string            `json:"model" validate:"required,min=2,max=255"`
	BusType     constants.BusType `json:"bus_type" validate:"required,oneof=standard vip sleeper double_decker"`
	Floors      []FloorConfig     `json:"floors" validate:"required,min=1,max=2,dive"`
	Amenities   []string          `json:"amenities"`
	IsActive    bool              `json:"is_active"`
}

// FloorConfig defines the seat layout for one floor
type FloorConfig struct {
	Floor   int          `json:"floor" validate:"required,min=1,max=2"`
	Rows    int          `json:"rows" validate:"required,min=1,max=20"`   // Number of seat rows
	Columns int          `json:"columns" validate:"required,min=1,max=5"` // Number of seats per row
	Seats   []SeatConfig `json:"seats" validate:"required,dive"`          // Individual seat configurations
}

// SeatConfig defines configuration for a single seat
type SeatConfig struct {
	Row             int                `json:"row" validate:"required,min=1"`
	Column          int                `json:"column" validate:"required,min=1"`
	SeatType        constants.SeatType `json:"seat_type" validate:"required,oneof=standard vip sleeper"`
	PriceMultiplier *float64           `json:"price_multiplier,omitempty" validate:"omitempty,min=0.5,max=5.0"` // Optional override
}

type UpdateBusRequest struct {
	PlateNumber *string            `json:"plate_number" validate:"omitempty,min=3,max=20"`
	Model       *string            `json:"model" validate:"omitempty,min=2,max=255"`
	BusType     *constants.BusType `json:"bus_type" validate:"omitempty,oneof=standard vip sleeper double_decker"`
	Amenities   *[]string          `json:"amenities"`
	IsActive    *bool              `json:"is_active"`
}
