package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// UserClient interface for user service communication
type UserClient interface {
	CreateGuestAccount(ctx context.Context, req *CreateGuestAccountRequest) (*GuestAccountResponse, error)
}

type userClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient(baseURL string) UserClient {
	return &userClientImpl{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// Request/Response types
type CreateGuestAccountRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type GuestAccountResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email,omitempty"`
	Phone    string    `json:"phone,omitempty"`
	Role     int       `json:"role"`
}

type UserServiceResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message,omitempty"`
	Data    *GuestAccountResponse `json:"data,omitempty"`
}

// CreateGuestAccount creates a guest user account
func (c *userClientImpl) CreateGuestAccount(ctx context.Context, req *CreateGuestAccountRequest) (*GuestAccountResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/guest", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call user service: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("body", string(body)).
			Msg("User service returned error")
		return nil, fmt.Errorf("user service error: %s", string(body))
	}

	var userResp UserServiceResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !userResp.Success || userResp.Data == nil {
		return nil, fmt.Errorf("invalid response from user service")
	}

	return userResp.Data, nil
}
