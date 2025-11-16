package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bus struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OperatorID    uuid.UUID      `gorm:"type:uuid;not null" json:"operator_id" validate:"required"`
	PlateNumber   string         `gorm:"type:varchar(20);unique;not null" json:"plate_number" validate:"required"`
	Model         string         `gorm:"type:varchar(255);not null" json:"model" validate:"required"`
	SeatCapacity  int            `gorm:"type:integer;not null" json:"seat_capacity" validate:"required,min=1,max=100"`
	AmenitiesJSON string         `gorm:"type:text" json:"-"`
	Amenities     []string       `gorm:"-" json:"amenities"`
	IsActive      bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Operator *Operator `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"operator,omitempty"`
	Seats    []Seat    `gorm:"foreignKey:BusID" json:"seats,omitempty"`
	Trips    []Trip    `gorm:"foreignKey:BusID" json:"trips,omitempty"`
}

func (Bus) TableName() string { return "buses" }

func (b *Bus) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b *Bus) BeforeSave(tx *gorm.DB) error {
	if len(b.Amenities) > 0 {
		amenitiesJSON, err := json.Marshal(b.Amenities)
		if err != nil {
			return err
		}
		b.AmenitiesJSON = string(amenitiesJSON)
	}
	return nil
}

func (b *Bus) AfterFind(tx *gorm.DB) error {
	if b.AmenitiesJSON != "" {
		return json.Unmarshal([]byte(b.AmenitiesJSON), &b.Amenities)
	}
	return nil
}
