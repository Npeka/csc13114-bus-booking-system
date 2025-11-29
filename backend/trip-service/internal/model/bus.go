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
	SeatCapacity int            `gorm:"type:integer;not null" json:"seat_capacity" validate:"required,min=1,max=100"`
	Amenities    pq.StringArray `gorm:"type:text[]" json:"amenities"`
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
	SeatCapacity int       `json:"seat_capacity"`
	Amenities    []string  `json:"amenities"`
	IsActive     bool      `json:"is_active"`

	Seats []Seat `gorm:"foreignKey:BusID" json:"seats"`
	Trips []Trip `gorm:"foreignKey:BusID" json:"trips"`
}

type CreateBusRequest struct {
	PlateNumber  string   `json:"plate_number" validate:"required,min=3,max=20"`
	Model        string   `json:"model" validate:"required,min=2,max=255"`
	SeatCapacity int      `json:"seat_capacity" validate:"required,min=1,max=100"`
	Amenities    []string `json:"amenities"`
	IsActive     bool     `json:"is_active"`
}

type UpdateBusRequest struct {
	PlateNumber  *string   `json:"plate_number" validate:"omitempty,min=3,max=20"`
	Model        *string   `json:"model" validate:"omitempty,min=2,max=255"`
	SeatCapacity *int      `json:"seat_capacity" validate:"omitempty,min=1,max=100"`
	Amenities    *[]string `json:"amenities"`
	IsActive     *bool     `json:"is_active"`
}
