package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Bus struct {
	BaseModel
	PlateNumber  string         `gorm:"type:varchar(20);unique;not null" json:"plate_number" validate:"required"`
	Model        string         `gorm:"type:varchar(255);not null" json:"model" validate:"required"`
	BusType      string         `gorm:"type:varchar(20);not null;default:'standard'" json:"bus_type" validate:"required"` // constants.BusType values
	SeatCapacity int            `gorm:"type:integer;not null" json:"seat_capacity" validate:"required,min=1,max=100"`
	Amenities    pq.StringArray `gorm:"type:text[]" json:"amenities"` // constants.Amenity values
	IsActive     bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`

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
	ID           uuid.UUID `json:"id"`
	PlateNumber  string    `json:"plate_number"`
	Model        string    `json:"model"`
	BusType      string    `json:"bus_type"` // Raw string value: constants.BusType
	SeatCapacity int       `json:"seat_capacity"`
	Amenities    []string  `json:"amenities"` // Raw string values: constants.Amenity
	IsActive     bool      `json:"is_active"`

	Seats []SeatResponse `json:"seats,omitempty"`
}

type FloorConfig struct {
	Floor        int `json:"floor" validate:"required,min=1"`
	SeatCapacity int `json:"seat_capacity" validate:"required,min=1"`
}

type CreateBusRequest struct {
	PlateNumber string        `json:"plate_number" validate:"required,min=3,max=20"`
	Model       string        `json:"model" validate:"required,min=2,max=255"`
	BusType     string        `json:"bus_type" validate:"required,oneof=standard vip sleeper double_decker"` // constants.BusType
	Floors      []FloorConfig `json:"floors" validate:"required,min=1,dive"`
	Amenities   []string      `json:"amenities"` // constants.Amenity values
	IsActive    bool          `json:"is_active"`
}

type UpdateBusRequest struct {
	PlateNumber  *string   `json:"plate_number" validate:"omitempty,min=3,max=20"`
	Model        *string   `json:"model" validate:"omitempty,min=2,max=255"`
	BusType      *string   `json:"bus_type" validate:"omitempty,oneof=standard vip sleeper double_decker"`
	SeatCapacity *int      `json:"seat_capacity" validate:"omitempty,min=1,max=100"`
	Amenities    *[]string `json:"amenities"`
	IsActive     *bool     `json:"is_active"`
}
