package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bus-booking/gateway-service/internal/config"
)

type Client struct {
	config     *config.AuthConfig
	httpClient *http.Client
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}

type VerifyTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
	Name   string `json:"name,omitempty"`
}

type UserContext struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Name   string `json:"name"`
}

func NewClient(config *config.AuthConfig) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// VerifyToken calls user-service to verify JWT token
func (c *Client) VerifyToken(ctx context.Context, token string) (*UserContext, error) {
	// Clean token (remove Bearer prefix if present)
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, fmt.Errorf("empty token")
	}

	// Prepare request
	reqBody := VerifyTokenRequest{Token: token}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := c.config.UserServiceURL + c.config.VerifyEndpoint
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token verification failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var verifyResp VerifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !verifyResp.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &UserContext{
		UserID: verifyResp.UserID,
		Email:  verifyResp.Email,
		Role:   verifyResp.Role,
		Name:   verifyResp.Name,
	}, nil
}

// HasRole checks if user has required role
func (uc *UserContext) HasRole(role string) bool {
	return uc.Role == role
}

// HasAnyRole checks if user has any of the required roles
func (uc *UserContext) HasAnyRole(roles []string) bool {
	for _, role := range roles {
		if uc.Role == role {
			return true
		}
	}
	return false
}

// ToHeaders converts user context to headers for downstream services
func (uc *UserContext) ToHeaders() map[string]string {
	return map[string]string{
		"X-User-ID":    uc.UserID,
		"X-User-Email": uc.Email,
		"X-User-Role":  uc.Role,
		"X-User-Name":  uc.Name,
	}
}
