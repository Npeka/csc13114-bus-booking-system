package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	BaseModel
	BookingID       uuid.UUID         `gorm:"type:uuid;not null;index;" json:"booking_id"`
	UserID          uuid.UUID         `gorm:"type:uuid;not null;index;" json:"user_id"`
	Amount          int               `gorm:"not null" json:"amount"`
	Currency        Currency          `gorm:"type:varchar(10);not null" json:"currency"`
	PaymentMethod   PaymentMethod     `gorm:"type:varchar(50);not null" json:"payment_method"`
	OrderCode       int64             `gorm:"index;unique" json:"order_code,omitempty"`
	PaymentLinkID   string            `gorm:"type:varchar(255)" json:"payment_link_id,omitempty"`
	Status          TransactionStatus `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	CheckoutURL     string            `gorm:"type:text" json:"checkout_url,omitempty"`
	QRCode          string            `gorm:"type:text" json:"qr_code,omitempty"`
	Reference       string            `gorm:"type:varchar(255)" json:"reference,omitempty"`
	TransactionTime *int64            `json:"transaction_time,omitempty"`
}

type Currency string
type PaymentMethod string
type TransactionStatus string

const (
	CurrencyVND Currency = "VND"

	PaymentMethodPayOS PaymentMethod = "PAYOS"

	TransactionStatusPending    TransactionStatus = "PENDING"
	TransactionStatusCancelled  TransactionStatus = "CANCELLED"
	TransactionStatusUnderpaid  TransactionStatus = "UNDERPAID"
	TransactionStatusPaid       TransactionStatus = "PAID"
	TransactionStatusExpired    TransactionStatus = "EXPIRED"
	TransactionStatusProcessing TransactionStatus = "PROCESSING"
	TransactionStatusFailed     TransactionStatus = "FAILED"
)

func (Transaction) TableName() string {
	return "transactions"
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type CreatePaymentLinkRequest struct {
	BookingID     uuid.UUID     `json:"booking_id" binding:"required"`
	Amount        int           `json:"amount" binding:"required,gt=0"`
	Currency      Currency      `json:"currency" binding:"required"`
	PaymentMethod PaymentMethod `json:"payment_method" binding:"required"`
	Description   string        `json:"description"`
}

type TransactionResponse struct {
	ID            uuid.UUID         `json:"id"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	BookingID     uuid.UUID         `json:"booking_id"`
	UserID        uuid.UUID         `json:"user_id"`
	Amount        int               `json:"amount"`
	Currency      Currency          `json:"currency"`
	PaymentMethod PaymentMethod     `json:"payment_method"`
	OrderCode     int64             `json:"order_code,omitempty"`
	Status        TransactionStatus `json:"status"`
	CheckoutURL   string            `json:"checkout_url,omitempty"`
	QRCode        string            `json:"qr_code,omitempty"`
}
