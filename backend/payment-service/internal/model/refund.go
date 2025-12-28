package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Refund represents a refund request for a cancelled booking
type Refund struct {
	BaseModel
	BookingID      uuid.UUID    `gorm:"type:uuid;not null;index" json:"booking_id"`
	TransactionID  uuid.UUID    `gorm:"type:uuid;not null;index" json:"transaction_id"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	RefundAmount   int          `gorm:"not null" json:"refund_amount"`
	RefundStatus   RefundStatus `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"refund_status"`
	RefundReason   string       `gorm:"type:text;not null" json:"refund_reason"`
	RejectedReason *string      `gorm:"type:text" json:"rejected_reason,omitempty"`
	ProcessedBy    *uuid.UUID   `gorm:"type:uuid" json:"processed_by,omitempty"`
	ProcessedAt    *time.Time   `json:"processed_at,omitempty"`

	// Relations (not stored in DB, loaded via preload)
	Transaction *Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
}

type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "PENDING"
	RefundStatusProcessing RefundStatus = "PROCESSING"
	RefundStatusCompleted  RefundStatus = "COMPLETED"
	RefundStatusRejected   RefundStatus = "REJECTED"
)

func (Refund) TableName() string {
	return "refunds"
}

func (r *Refund) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// RefundRequest represents a request to create a refund
type RefundRequest struct {
	BookingID    uuid.UUID `json:"booking_id" binding:"required"`
	Reason       string    `json:"reason" binding:"required,min=10,max=500"`
	RefundAmount int       `json:"refund_amount" binding:"required,gt=0"`
}

// RefundResponse represents a refund transaction with user info
type RefundResponse struct {
	ID                    uuid.UUID    `json:"id"`
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`
	BookingID             uuid.UUID    `json:"booking_id"`
	UserID                uuid.UUID    `json:"user_id"`
	RefundAmount          int          `json:"refund_amount"`
	RefundStatus          RefundStatus `json:"refund_status"`
	RefundReason          string       `json:"refund_reason"`
	OriginalTransactionID uuid.UUID    `json:"original_transaction_id"` // For compatibility
	ProcessedBy           *uuid.UUID   `json:"processed_by,omitempty"`
	ProcessedAt           *time.Time   `json:"processed_at,omitempty"`

	// User bank account info (for export)
	BankCode      string `json:"bank_code,omitempty"`
	BankName      string `json:"bank_name,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	AccountHolder string `json:"account_holder,omitempty"`
}

// RefundListQuery represents query parameters for listing refunds
type RefundListQuery struct {
	PaginationRequest
	Status    *RefundStatus `form:"status"`
	StartDate *time.Time    `form:"start_date"`
	EndDate   *time.Time    `form:"end_date"`
}

// UpdateRefundStatusRequest represents request to update refund status
type UpdateRefundStatusRequest struct {
	Status RefundStatus `json:"status" binding:"required,oneof=PROCESSING COMPLETED REJECTED"`
}

// RefundExportItem represents an item in the refund Excel export
type RefundExportItem struct {
	BookingReference string
	UserName         string
	BankCode         string
	BankName         string
	AccountNumber    string
	AccountHolder    string
	RefundAmount     int
	Reason           string
	CreatedDate      time.Time
}

// ExportRefundsRequest represents request to export refunds
type ExportRefundsRequest struct {
	RefundIDs []uuid.UUID `json:"refund_ids" binding:"required,min=1"`
}

// TransactionStats represents transaction statistics (including refunds)
type TransactionStats struct {
	TotalTransactions  int `json:"total_transactions"`
	TotalIn            int `json:"total_in"`  // Total payment received
	TotalOut           int `json:"total_out"` // Total refunds paid
	PendingRefunds     int `json:"pending_refunds"`
	PendingRefundCount int `json:"pending_refund_count"`
}
