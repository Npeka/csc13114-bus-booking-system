package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"bus-booking/payment-service/internal/model"

	"github.com/rs/zerolog/log"
)

// PayOSClient handles PayOS API integration
type PayOSClient struct {
	clientID    string
	apiKey      string
	checksumKey string
	baseURL     string
	httpClient  *http.Client
}

// NewPayOSClient creates a new PayOS client
func NewPayOSClient(clientID, apiKey, checksumKey string) *PayOSClient {
	return &PayOSClient{
		clientID:    clientID,
		apiKey:      apiKey,
		checksumKey: checksumKey,
		baseURL:     "https://api-merchant.payos.vn",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreatePaymentLink creates a payment link via PayOS API
func (c *PayOSClient) CreatePaymentLink(req *model.CreatePaymentLinkRequest) (*model.CreatePaymentLinkResponse, error) {
	// Generate signature
	req.Signature = c.generateSignature(req)

	// Marshal request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.baseURL+"/v2/payment-requests", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-client-id", c.clientID)
	httpReq.Header.Set("x-api-key", c.apiKey)

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(body)).
			Msg("PayOS API returned error")
		return nil, fmt.Errorf("PayOS API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Unmarshal response
	var paymentResp model.CreatePaymentLinkResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Verify signature
	if !c.verifyResponseSignature(&paymentResp) {
		return nil, fmt.Errorf("invalid response signature")
	}

	return &paymentResp, nil
}

// GetPaymentInfo retrieves payment information
func (c *PayOSClient) GetPaymentInfo(orderCode int64) (*model.GetPaymentInfoResponse, error) {
	url := fmt.Sprintf("%s/v2/payment-requests/%d", c.baseURL, orderCode)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("x-client-id", c.clientID)
	httpReq.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PayOS API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var paymentInfo model.GetPaymentInfoResponse
	if err := json.Unmarshal(body, &paymentInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &paymentInfo, nil
}

// CancelPayment cancels a payment
func (c *PayOSClient) CancelPayment(orderCode int64, reason string) (*model.GetPaymentInfoResponse, error) {
	url := fmt.Sprintf("%s/v2/payment-requests/%d/cancel", c.baseURL, orderCode)

	cancelReq := model.CancelPaymentRequest{
		CancellationReason: reason,
	}

	jsonData, err := json.Marshal(cancelReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-client-id", c.clientID)
	httpReq.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PayOS API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var paymentInfo model.GetPaymentInfoResponse
	if err := json.Unmarshal(body, &paymentInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &paymentInfo, nil
}

// VerifyWebhookSignature verifies webhook signature from PayOS
func (c *PayOSClient) VerifyWebhookSignature(webhookData *model.PaymentWebhookData) bool {
	// Extract data for signature verification
	data := webhookData.Data

	// Sort fields alphabetically and create signature string
	signatureStr := fmt.Sprintf(
		"amount=%d&code=%s&desc=%s&orderCode=%d&paymentLinkId=%s&reference=%s&transactionDateTime=%s",
		data.Amount,
		data.Code,
		data.Desc,
		data.OrderCode,
		data.PaymentLinkID,
		data.Reference,
		data.TransactionDateTime,
	)

	// Calculate HMAC SHA256
	h := hmac.New(sha256.New, []byte(c.checksumKey))
	h.Write([]byte(signatureStr))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return expectedSignature == webhookData.Signature
}

// generateSignature generates signature for payment link creation
func (c *PayOSClient) generateSignature(req *model.CreatePaymentLinkRequest) string {
	// Sort parameters alphabetically and create signature string
	signatureStr := fmt.Sprintf(
		"amount=%d&cancelUrl=%s&description=%s&orderCode=%d&returnUrl=%s",
		req.Amount,
		req.CancelURL,
		req.Description,
		req.OrderCode,
		req.ReturnURL,
	)

	// Calculate HMAC SHA256
	h := hmac.New(sha256.New, []byte(c.checksumKey))
	h.Write([]byte(signatureStr))
	return hex.EncodeToString(h.Sum(nil))
}

// verifyResponseSignature verifies signature in PayOS response
func (c *PayOSClient) verifyResponseSignature(resp *model.CreatePaymentLinkResponse) bool {
	if resp.Signature == "" {
		return false
	}

	// Extract fields for verification (sorted alphabetically)
	data := resp.Data
	var fields []string

	fields = append(fields, fmt.Sprintf("accountName=%s", data.AccountName))
	fields = append(fields, fmt.Sprintf("accountNumber=%s", data.AccountNumber))
	fields = append(fields, fmt.Sprintf("amount=%d", data.Amount))
	fields = append(fields, fmt.Sprintf("bin=%s", data.Bin))
	fields = append(fields, fmt.Sprintf("checkoutUrl=%s", data.CheckoutURL))
	fields = append(fields, fmt.Sprintf("currency=%s", data.Currency))
	fields = append(fields, fmt.Sprintf("description=%s", data.Description))
	fields = append(fields, fmt.Sprintf("orderCode=%d", data.OrderCode))
	fields = append(fields, fmt.Sprintf("paymentLinkId=%s", data.PaymentLinkID))
	fields = append(fields, fmt.Sprintf("qrCode=%s", data.QRCode))
	fields = append(fields, fmt.Sprintf("status=%s", data.Status))

	sort.Strings(fields)
	signatureStr := strings.Join(fields, "&")

	// Calculate HMAC SHA256
	h := hmac.New(sha256.New, []byte(c.checksumKey))
	h.Write([]byte(signatureStr))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return expectedSignature == resp.Signature
}
