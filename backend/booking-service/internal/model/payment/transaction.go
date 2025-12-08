package payment

import (
	"time"

	"github.com/google/uuid"
)

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

type TransactionResponse struct {
	ID            uuid.UUID         `json:"id"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	BookingID     uuid.UUID         `json:"booking_id"`
	Amount        int               `json:"amount"`
	Currency      Currency          `json:"currency"`
	PaymentMethod PaymentMethod     `json:"payment_method"`
	OrderCode     int64             `json:"order_code,omitempty"`
	Status        TransactionStatus `json:"status"`
	CheckoutURL   string            `json:"checkout_url,omitempty"`
	QRCode        string            `json:"qr_code,omitempty"`
}
