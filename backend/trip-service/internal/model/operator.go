package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Operator struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	ContactEmail string         `gorm:"type:varchar(255);unique;not null" json:"contact_email" validate:"required,email"`
	ContactPhone string         `gorm:"type:varchar(20)" json:"contact_phone"`
	Status       string         `gorm:"type:varchar(50);not null;default:'pending'" json:"status" validate:"oneof=pending approved rejected"`
	ApprovedAt   *time.Time     `gorm:"type:timestamptz" json:"approved_at,omitempty"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Routes []Route `gorm:"foreignKey:OperatorID" json:"routes,omitempty"`
	Buses  []Bus   `gorm:"foreignKey:OperatorID" json:"buses,omitempty"`
}

func (Operator) TableName() string { return "operators" }

func (o *Operator) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
