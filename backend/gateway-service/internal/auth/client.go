package auth

import (
	"bus-booking/gateway-service/config"
	"bus-booking/shared/constants"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	config     *config.AuthConfig
	httpClient *http.Client
}

type VerifyTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type VerifyTokenResponse struct {
	UserID string             `json:"user_id,omitempty"`
	Email  string             `json:"email,omitempty"`
	Role   constants.UserRole `json:"role,omitempty"`
	Name   string             `json:"name,omitempty"`
}

type UserContext struct {
	UserID      string             `json:"user_id"`
	Email       string             `json:"email"`
	Role        constants.UserRole `json:"role"`
	Name        string             `json:"name"`
	AccessToken string             `json:"access_token"`
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
func (c *Client) VerifyToken(ctx context.Context, accessToken string) (*UserContext, error) {
	// Clean token (remove Bearer prefix if present)
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	accessToken = strings.TrimSpace(accessToken)

	if accessToken == "" {
		return nil, fmt.Errorf("empty token")
	}

	// Prepare request
	reqBody := VerifyTokenRequest{AccessToken: accessToken}
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unauthorized")
	}

	// Parse response
	type responseBody struct {
		Data VerifyTokenResponse `json:"data"`
	}
	var verifyResp responseBody
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &UserContext{
		UserID:      verifyResp.Data.UserID,
		Email:       verifyResp.Data.Email,
		Role:        verifyResp.Data.Role,
		Name:        verifyResp.Data.Name,
		AccessToken: accessToken,
	}, nil
}

func (uc *UserContext) HasRole(role constants.UserRole) bool {
	return uc.Role == role
}

func (uc *UserContext) HasAnyRole(roles []constants.UserRole) bool {
	for _, role := range roles {
		if uc.Role == role {
			return true
		}
	}
	return false
}

// HasAnyRoleString checks if user has any of the specified roles (from string slice)
func (uc *UserContext) HasAnyRoleString(roleStrings []string) bool {
	for _, roleStr := range roleStrings {
		if uc.Role == constants.FromString(roleStr) {
			return true
		}
	}
	return false
}

func (uc *UserContext) ToHeaders() map[string]string {
	return map[string]string{
		"X-User-ID":      uc.UserID,
		"X-User-Email":   uc.Email,
		"X-User-Role":    uc.Role.String(),
		"X-User-Name":    uc.Name,
		"X-Access-Token": uc.AccessToken,
	}
}
