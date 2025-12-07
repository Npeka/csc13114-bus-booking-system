package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	BaseModel
	BookingID     uuid.UUID `gorm:"type:uuid;not null;index;" json:"booking_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Currency      string    `gorm:"type:varchar(10);not null" json:"currency"`
	PaymentMethod string    `gorm:"type:varchar(50);not null" json:"payment_method"`

	// PayOS fields
	OrderCode       int64  `gorm:"index;unique" json:"order_code,omitempty"`
	PaymentLinkID   string `gorm:"type:varchar(255)" json:"payment_link_id,omitempty"`
	Status          string `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	CheckoutURL     string `gorm:"type:text" json:"checkout_url,omitempty"`
	QRCode          string `gorm:"type:text" json:"qr_code,omitempty"`
	Reference       string `gorm:"type:varchar(255)" json:"reference,omitempty"`
	TransactionTime *int64 `json:"transaction_time,omitempty"`
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
	Description   string    `json:"description"`
	BuyerName     string    `json:"buyer_name"`
	BuyerEmail    string    `json:"buyer_email"`
	BuyerPhone    string    `json:"buyer_phone"`
}

type TransactionResponse struct {
	ID            uuid.UUID `json:"id"`
	BookingID     uuid.UUID `json:"booking_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	OrderCode     int64     `json:"order_code,omitempty"`
	Status        string    `json:"status"`
	CheckoutURL   string    `json:"checkout_url,omitempty"`
	QRCode        string    `json:"qr_code,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PaymentCallbackRequest struct {
	OrderCode int64  `form:"orderCode" binding:"required"`
	Status    string `form:"status"`
	Cancel    bool   `form:"cancel"`
}
