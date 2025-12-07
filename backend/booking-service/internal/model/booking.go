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
	TotalAmount float64 `json:"total_amount" gorm:"type:decimal(10,2);not null"`

	// Status
	Status        BookingStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"type:varchar(20);not null;default:'pending';index"`

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
	BookingStatusPending   BookingStatus = "pending"   // Chờ thanh toán
	BookingStatusConfirmed BookingStatus = "confirmed" // Đã xác nhận và thanh toán
	BookingStatusCancelled BookingStatus = "cancelled" // Đã hủy
	BookingStatusExpired   BookingStatus = "expired"   // Hết hạn (quá 15 phút chưa thanh toán)
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

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"  // Chờ thanh toán
	PaymentStatusPaid     PaymentStatus = "paid"     // Đã thanh toán
	PaymentStatusRefunded PaymentStatus = "refunded" // Đã hoàn tiền
	PaymentStatusFailed   PaymentStatus = "failed"   // Thanh toán thất bại
)

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusPending,
		PaymentStatusPaid,
		PaymentStatusRefunded,
		PaymentStatusFailed:
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
	PaymentStatus  string `json:"payment_status" binding:"required"`
	BookingStatus  string `json:"booking_status" binding:"required"`
	PaymentOrderID string `json:"payment_order_id" binding:"required"`
}
