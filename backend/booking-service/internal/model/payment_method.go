package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Code        string `gorm:"type:varchar(50);not null;unique" json:"code" validate:"required"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"type:boolean;not null;default:true" json:"is_active"`
}

func (PaymentMethod) TableName() string { return "payment_methods" }

func (pm *PaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}
