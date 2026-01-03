package payment

import "github.com/google/uuid"

type RefundRequest struct {
	BookingID    uuid.UUID `json:"booking_id"`
	Reason       string    `json:"reason"`
	RefundAmount int       `json:"refund_amount"`
}

type RefundResponse struct {
	ID uuid.UUID `json:"id"`
	RefundStatus string `json:"refund_status"`
}
