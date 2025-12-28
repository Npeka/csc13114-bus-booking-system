package service

import (
	"context"
	"testing"
	"time"

	"bus-booking/payment-service/config"
	"bus-booking/payment-service/internal/model"

	"github.com/payOSHQ/payos-lib-golang/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewPayOSService(t *testing.T) {
	cfg := config.PayOSConfig{
		ClientID:    "test-client-id",
		APIKey:      "test-api-key",
		ChecksumKey: "test-checksum-key",
		ReturnURL:   "https://example.com/return",
		CancelURL:   "https://example.com/cancel",
	}

	service := NewPayOSService(cfg)

	assert.NotNil(t, service)
	assert.IsType(t, &PayOSServiceImpl{}, service)

	impl := service.(*PayOSServiceImpl)
	assert.Equal(t, "https://example.com/return", impl.ReturnURL)
	assert.Equal(t, "https://example.com/cancel", impl.CancelURL)
}

func TestToTransactionStatus(t *testing.T) {
	cfg := config.PayOSConfig{
		ClientID:    "test-client-id",
		APIKey:      "test-api-key",
		ChecksumKey: "test-checksum-key",
		ReturnURL:   "https://example.com/return",
		CancelURL:   "https://example.com/cancel",
	}

	service := NewPayOSService(cfg)

	tests := []struct {
		name           string
		payOSStatus    payos.PaymentLinkStatus
		expectedStatus model.TransactionStatus
	}{
		{
			name:           "Pending status",
			payOSStatus:    payos.PaymentLinkStatusPending,
			expectedStatus: model.TransactionStatusPending,
		},
		{
			name:           "Paid status",
			payOSStatus:    payos.PaymentLinkStatusPaid,
			expectedStatus: model.TransactionStatusPaid,
		},
		{
			name:           "Cancelled status",
			payOSStatus:    payos.PaymentLinkStatusCancelled,
			expectedStatus: model.TransactionStatusCancelled,
		},
		{
			name:           "Expired status",
			payOSStatus:    payos.PaymentLinkStatusExpired,
			expectedStatus: model.TransactionStatusExpired,
		},
		{
			name:           "Processing status",
			payOSStatus:    payos.PaymentLinkStatusProcessing,
			expectedStatus: model.TransactionStatusProcessing,
		},
		{
			name:           "Failed status",
			payOSStatus:    payos.PaymentLinkStatusFailed,
			expectedStatus: model.TransactionStatusFailed,
		},
		{
			name:           "Underpaid status",
			payOSStatus:    payos.PaymentLinkStatusUnderpaid,
			expectedStatus: model.TransactionStatusUnderpaid,
		},
		{
			name:           "Unknown status defaults to Pending",
			payOSStatus:    payos.PaymentLinkStatus("UNKNOWN"),
			expectedStatus: model.TransactionStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ToTransactionStatus(tt.payOSStatus)
			assert.Equal(t, tt.expectedStatus, result)
		})
	}
}

func TestCreatePaymentLink_ExpirationValidation(t *testing.T) {
	cfg := config.PayOSConfig{
		ClientID:    "test-client-id",
		APIKey:      "test-api-key",
		ChecksumKey: "test-checksum-key",
		ReturnURL:   "https://example.com/return",
		CancelURL:   "https://example.com/cancel",
	}

	service := NewPayOSService(cfg)
	ctx := context.Background()

	t.Run("Expired time in the past", func(t *testing.T) {
		req := &model.CreatePaymentLinkRequest{
			Amount:      100000,
			Description: "Test payment",
			ExpiresAt:   time.Now().Add(-1 * time.Hour), // Past time
		}

		result, err := service.CreatePaymentLink(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "must be in the future")
	})

	t.Run("Expired time beyond 2038", func(t *testing.T) {
		req := &model.CreatePaymentLinkRequest{
			Amount:      100000,
			Description: "Test payment",
			ExpiresAt:   time.Date(2039, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		result, err := service.CreatePaymentLink(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "beyond Year 2038")
	})
}

// Note: Integration tests with actual PayOS API are skipped in unit tests
// These would require valid credentials and should be tested separately
func TestCreatePaymentLink_Integration_Skipped(t *testing.T) {
	t.Skip("Skipping integration test - requires actual PayOS credentials")

	// This is an example of how to structure an integration test
	// cfg := config.PayOSConfig{
	// 	ClientID:    os.Getenv("PAYOS_CLIENT_ID"),
	// 	APIKey:      os.Getenv("PAYOS_API_KEY"),
	// 	ChecksumKey: os.Getenv("PAYOS_CHECKSUM_KEY"),
	// 	ReturnURL:   "https://example.com/return",
	// 	CancelURL:   "https://example.com/cancel",
	// }
	//
	// service := NewPayOSService(cfg)
	// ctx := context.Background()
	//
	// req := &model.CreatePaymentLinkRequest{
	// 	Amount:      100000,
	// 	Description: "Test payment",
	// 	ExpiresAt:   time.Now().Add(15 * time.Minute),
	// }
	//
	// result, err := service.CreatePaymentLink(ctx, req)
	//
	// assert.NoError(t, err)
	// assert.NotNil(t, result)
	// assert.NotEmpty(t, result.CheckoutUrl)
}

func TestGetPaymentLink_Integration_Skipped(t *testing.T) {
	t.Skip("Skipping integration test - requires actual PayOS credentials")
}

func TestCancelPaymentLink_Integration_Skipped(t *testing.T) {
	t.Skip("Skipping integration test - requires actual PayOS credentials")
}

func TestVerifyWebhook_Integration_Skipped(t *testing.T) {
	t.Skip("Skipping integration test - requires actual PayOS credentials")
}
