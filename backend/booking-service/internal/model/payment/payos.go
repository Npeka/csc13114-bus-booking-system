package payment

import "github.com/google/uuid"

type CreatePaymentLinkRequest struct {
	BookingID     uuid.UUID     `json:"booking_id"`
	Amount        int           `json:"amount"`
	Currency      Currency      `json:"currency"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Description   string        `json:"description"`
}
