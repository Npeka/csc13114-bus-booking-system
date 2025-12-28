package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BankAccount represents a user's bank account for refund purposes
type BankAccount struct {
	BaseModel
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	BankCode      string    `gorm:"type:varchar(20);not null" json:"bank_code"`
	AccountNumber string    `gorm:"type:varchar(50);not null" json:"account_number"`
	AccountHolder string    `gorm:"type:varchar(100);not null" json:"account_holder"`
	IsPrimary     bool      `gorm:"default:false;index" json:"is_primary"`
}

func (BankAccount) TableName() string {
	return "bank_accounts"
}

func (b *BankAccount) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// BankConstant represents a Vietnamese bank
type BankConstant struct {
	Code      string `json:"code"`
	ShortName string `json:"short_name"`
	Name      string `json:"name"`
	Logo      string `json:"logo,omitempty"`
}

// BankAccountRequest represents the request to create/update a bank account
type BankAccountRequest struct {
	BankCode      string `json:"bank_code" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
	AccountHolder string `json:"account_holder" binding:"required,min=1,max=100"`
}

// BankAccountResponse represents the response for bank account
type BankAccountResponse struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        uuid.UUID `json:"user_id"`
	BankCode      string    `json:"bank_code"`
	BankName      string    `json:"bank_name"` // Resolved from BankCode
	AccountNumber string    `json:"account_number"`
	AccountHolder string    `json:"account_holder"`
	IsPrimary     bool      `json:"is_primary"`
}

// ToResponse converts BankAccount to BankAccountResponse
func (b *BankAccount) ToResponse(bankName string) *BankAccountResponse {
	return &BankAccountResponse{
		ID:            b.ID,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		UserID:        b.UserID,
		BankCode:      b.BankCode,
		BankName:      bankName,
		AccountNumber: b.AccountNumber,
		AccountHolder: b.AccountHolder,
		IsPrimary:     b.IsPrimary,
	}
}
