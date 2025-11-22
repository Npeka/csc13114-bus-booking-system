package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	CreatedAt int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	BaseModel
	BookingID     uuid.UUID `gorm:"type:uuid;not null;index;" json:"booking_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Currency      string    `gorm:"type:varchar(10);not null" json:"currency"`
	PaymentMethod string    `gorm:"type:varchar(50);not null" json:"payment_method"`
}

func (Transaction) TableName() string { return "transactions" }

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type CreateTransactionRequest struct {
	BookingID     uuid.UUID `json:"booking_id" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	Currency      string    `json:"currency" binding:"required"`
	PaymentMethod string    `json:"payment_method" binding:"required"`
}
