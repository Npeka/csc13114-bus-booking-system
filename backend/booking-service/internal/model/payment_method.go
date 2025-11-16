package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Code        string         `gorm:"type:varchar(50);not null;unique" json:"code" validate:"required"`
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PaymentMethod) TableName() string { return "payment_methods" }

func (pm *PaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}
