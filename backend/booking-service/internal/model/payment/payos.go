package payment

import (
	"time"

	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	ID            uuid.UUID     `json:"id"`
	BookingID     uuid.UUID     `json:"booking_id"`
	Amount        int           `json:"amount"`
	Currency      Currency      `json:"currency"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Description   string        `json:"description"`
	ExpiresAt     time.Time     `json:"expires_at"`
}
