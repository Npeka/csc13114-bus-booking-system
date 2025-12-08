package model

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	BaseModel
	BookingReference string    `json:"booking_reference" gorm:"type:varchar(20);unique;not null;index"`
	TripID           uuid.UUID `json:"trip_id" gorm:"type:uuid;not null;index"`
	UserID           uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`

	// Pricing
	TotalAmount int `json:"total_amount" gorm:"type:decimal(10,2);not null"`

	// Status
	Status        BookingStatus     `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`
	PaymentStatus TransactionStatus `json:"payment_status" gorm:"type:varchar(20);not null;default:'pending';index"`

	// Payment info - Payment Service handles PayOS integration
	PaymentOrderID string `json:"payment_order_id,omitempty" gorm:"type:varchar(255);index"`

	// Timestamps
	ExpiresAt   *time.Time `json:"expires_at,omitempty" gorm:"type:timestamptz;index"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty" gorm:"type:timestamptz"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty" gorm:"type:timestamptz"`

	// Optional
	CancellationReason string `json:"cancellation_reason,omitempty" gorm:"type:text"`
	Notes              string `json:"notes,omitempty" gorm:"type:text"`

	// Relations
	BookingSeats []BookingSeat `json:"booking_seats,omitempty" gorm:"foreignKey:BookingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusExpired   BookingStatus = "EXPIRED"
	BookingStatusFailed    BookingStatus = "FAILED"
)

func (s BookingStatus) IsValid() bool {
	switch s {
	case BookingStatusPending,
		BookingStatusConfirmed,
		BookingStatusCancelled,
		BookingStatusExpired:
		return true
	default:
		return false
	}
}

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "PENDING"
	TransactionStatusCancelled  TransactionStatus = "CANCELLED"
	TransactionStatusUnderpaid  TransactionStatus = "UNDERPAID"
	TransactionStatusPaid       TransactionStatus = "PAID"
	TransactionStatusExpired    TransactionStatus = "EXPIRED"
	TransactionStatusProcessing TransactionStatus = "PROCESSING"
	TransactionStatusFailed     TransactionStatus = "FAILED"
)

func (s TransactionStatus) IsValid() bool {
	switch s {
	case TransactionStatusPending,
		TransactionStatusCancelled,
		TransactionStatusUnderpaid,
		TransactionStatusPaid,
		TransactionStatusExpired,
		TransactionStatusProcessing,
		TransactionStatusFailed:
		return true
	default:
		return false
	}
}

func (Booking) TableName() string {
	return "bookings"
}

// BuyerInfo contains buyer information for payment
type BuyerInfo struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

// CreatePaymentRequest is the request to create payment for a booking
type CreatePaymentRequest struct {
	BuyerInfo BuyerInfo `json:"buyer_info" binding:"required"`
}

// UpdatePaymentStatusRequest updates booking payment status (internal use)
type UpdatePaymentStatusRequest struct {
	PaymentStatus  TransactionStatus `json:"payment_status" binding:"required"`
	BookingStatus  BookingStatus     `json:"booking_status" binding:"required"`
	PaymentOrderID string            `json:"payment_order_id" binding:"required"`
}
