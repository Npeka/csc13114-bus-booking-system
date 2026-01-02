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

	// Transaction type and refund fields
	TransactionType TransactionType `gorm:"type:varchar(10);not null;default:'IN';index" json:"transaction_type"`
	RefundStatus    *RefundStatus   `gorm:"type:varchar(20);index" json:"refund_status,omitempty"`
	RefundAmount    *int            `json:"refund_amount,omitempty"`
}

type Currency string
type PaymentMethod string
type TransactionStatus string
type TransactionType string

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

	TransactionTypeIn  TransactionType = "IN"
	TransactionTypeOut TransactionType = "OUT"
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

type CreateTransactionRequest struct {
	ID            uuid.UUID     `json:"id"`
	BookingID     uuid.UUID     `json:"booking_id" binding:"required"`
	Amount        int           `json:"amount" binding:"required,gt=0"`
	Currency      Currency      `json:"currency" binding:"required"`
	PaymentMethod PaymentMethod `json:"payment_method" binding:"required"`
	Description   string        `json:"description"`
	ExpiresAt     time.Time     `json:"expires_at"`
}

type TransactionResponse struct {
	ID              uuid.UUID         `json:"id"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	BookingID       uuid.UUID         `json:"booking_id"`
	UserID          uuid.UUID         `json:"user_id"`
	Amount          int               `json:"amount"`
	Currency        Currency          `json:"currency"`
	PaymentMethod   PaymentMethod     `json:"payment_method"`
	OrderCode       int64             `json:"order_code,omitempty"`
	Status          TransactionStatus `json:"status"`
	CheckoutURL     string            `json:"checkout_url,omitempty"`
	QRCode          string            `json:"qr_code,omitempty"`
	TransactionType TransactionType   `json:"transaction_type"`
	RefundStatus    *RefundStatus     `json:"refund_status,omitempty"`
	RefundAmount    *int              `json:"refund_amount,omitempty"`
}

// TransactionListQuery represents query parameters for listing transactions
type TransactionListQuery struct {
	PaginationRequest
	TransactionType *TransactionType   `form:"transaction_type"`
	Status          *TransactionStatus `form:"status"`
	RefundStatus    *RefundStatus      `form:"refund_status"`
	StartDate       *time.Time         `form:"start_date"`
	EndDate         *time.Time         `form:"end_date"`
}
