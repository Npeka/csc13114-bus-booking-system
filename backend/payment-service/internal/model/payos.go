package model

import "time"

type CreatePaymentLinkRequest struct {
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// PaymentItem represents an item in payment
type PaymentItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

// PayOSResponse represents generic PayOS API response
type PayOSResponse struct {
	Code      string      `json:"code"`
	Desc      string      `json:"desc"`
	Data      interface{} `json:"data"`
	Signature string      `json:"signature,omitempty"`
}

// CreatePaymentLinkResponse represents response from create payment link API
type CreatePaymentLinkResponse struct {
	Code      string          `json:"code"`
	Desc      string          `json:"desc"`
	Data      PaymentLinkData `json:"data"`
	Signature string          `json:"signature"`
}

// PaymentLinkData contains payment link information
type PaymentLinkData struct {
	Bin           string `json:"bin"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	Currency      string `json:"currency"`
	PaymentLinkID string `json:"paymentLinkId"`
	Amount        int    `json:"amount"`
	Description   string `json:"description"`
	OrderCode     int64  `json:"orderCode"`
	ExpiredAt     int64  `json:"expiredAt"`
	Status        string `json:"status"`
	CheckoutURL   string `json:"checkoutUrl"`
	QRCode        string `json:"qrCode"`
}

// PaymentWebhookData represents webhook payload from PayOS
type PaymentWebhookData struct {
	Code      string                `json:"code"`
	Desc      string                `json:"desc"`
	Success   bool                  `json:"success"`
	Data      PaymentWebhookDetails `json:"data"`
	Signature string                `json:"signature"`
}

// PaymentWebhookDetails contains payment transaction details
type PaymentWebhookDetails struct {
	OrderCode              int    `json:"orderCode"`
	Amount                 int    `json:"amount"`
	Description            string `json:"description"`
	AccountNumber          string `json:"accountNumber"`
	Reference              string `json:"reference"`
	TransactionDateTime    string `json:"transactionDateTime"`
	Currency               string `json:"currency"`
	PaymentLinkID          string `json:"paymentLinkId"`
	Code                   string `json:"code"`
	Desc                   string `json:"desc"`
	CounterAccountBankID   string `json:"counterAccountBankId,omitempty"`
	CounterAccountBankName string `json:"counterAccountBankName,omitempty"`
	CounterAccountName     string `json:"counterAccountName,omitempty"`
	CounterAccountNumber   string `json:"counterAccountNumber,omitempty"`
	VirtualAccountName     string `json:"virtualAccountName,omitempty"`
	VirtualAccountNumber   string `json:"virtualAccountNumber,omitempty"`
}

// GetPaymentInfoResponse represents response from get payment info API
type GetPaymentInfoResponse struct {
	Code      string          `json:"code"`
	Desc      string          `json:"desc"`
	Data      PaymentInfoData `json:"data"`
	Signature string          `json:"signature"`
}

// PaymentInfoData contains detailed payment information
type PaymentInfoData struct {
	ID                 string               `json:"id"`
	OrderCode          int64                `json:"orderCode"`
	Amount             int                  `json:"amount"`
	AmountPaid         int                  `json:"amountPaid"`
	AmountRemaining    int                  `json:"amountRemaining"`
	Status             string               `json:"status"`
	CreatedAt          time.Time            `json:"createdAt"`
	Transactions       []PaymentTransaction `json:"transactions"`
	CanceledAt         *time.Time           `json:"canceledAt,omitempty"`
	CancellationReason string               `json:"cancellationReason,omitempty"`
}

// PaymentTransaction represents a payment transaction
type PaymentTransaction struct {
	Amount                 int       `json:"amount"`
	Description            string    `json:"description"`
	AccountNumber          string    `json:"accountNumber"`
	Reference              string    `json:"reference"`
	TransactionDateTime    time.Time `json:"transactionDateTime"`
	CounterAccountBankID   string    `json:"counterAccountBankId,omitempty"`
	CounterAccountBankName string    `json:"counterAccountBankName,omitempty"`
	CounterAccountName     string    `json:"counterAccountName,omitempty"`
	CounterAccountNumber   string    `json:"counterAccountNumber,omitempty"`
	VirtualAccountName     string    `json:"virtualAccountName,omitempty"`
	VirtualAccountNumber   string    `json:"virtualAccountNumber,omitempty"`
}

// CancelPaymentRequest represents request to cancel payment
type CancelPaymentRequest struct {
	CancellationReason string `json:"cancellationReason,omitempty"`
}

// Payment status constants
const (
	PaymentStatusPending   = "PENDING"
	PaymentStatusPaid      = "PAID"
	PaymentStatusCancelled = "CANCELLED"
	PaymentStatusExpired   = "EXPIRED"
)

// PayOS response codes
const (
	PayOSCodeSuccess = "00"
	PayOSCodeFailed  = "01"
)
