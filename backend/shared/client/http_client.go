package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"bus-booking/shared/constants"
	sharedcontext "bus-booking/shared/context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// HTTPClient interface for making HTTP requests to other microservices
type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
	Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error)
	Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
	Patch(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error)
}

// HTTPResponse represents HTTP response
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string][]string
}

// Config for HTTP client
type Config struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
	ServiceName    string
	DefaultHeaders map[string]string
}

// Client implements HTTPClient
type Client struct {
	httpClient *http.Client
	config     *Config
}

// NewHTTPClient creates a new HTTP client for microservice communication
func NewHTTPClient(config *Config) HTTPClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = time.Second
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}
}

// Get makes a GET request
func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodGet, url, nil, headers)
}

// Post makes a POST request
func (c *Client) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPost, url, body, headers)
}

// Put makes a PUT request
func (c *Client) Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPut, url, body, headers)
}

// Delete makes a DELETE request
func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodDelete, url, nil, headers)
}

// Patch makes a PATCH request
func (c *Client) Patch(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPatch, url, body, headers)
}

// doRequest performs the actual HTTP request with retry logic
func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.config.RetryDelay * time.Duration(attempt)):
			}
		}

		resp, err := c.makeRequest(ctx, method, url, body, headers)
		if err != nil {
			lastErr = err
			log.Warn().
				Str("method", method).
				Str("url", url).
				Int("attempt", attempt+1).
				Str("error", err.Error()).
				Msg("HTTP request failed, retrying...")
			continue
		}

		// Don't retry on client errors (4xx)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return resp, nil
		}

		// Retry on server errors (5xx)
		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			log.Warn().
				Str("method", method).
				Str("url", url).
				Int("status_code", resp.StatusCode).
				Int("attempt", attempt+1).
				Msg("Server error, retrying...")
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.config.MaxRetries+1, lastErr)
}

// makeRequest creates and executes a single HTTP request
func (c *Client) makeRequest(ctx context.Context, method, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	var reqBody io.Reader

	// Prepare request body
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	c.setHeaders(req, headers)

	// Extract request context and set microservice headers
	c.setMicroserviceHeaders(req, ctx)

	// Log request
	log.Debug().
		Str("method", method).
		Str("url", url).
		Str("request_id", req.Header.Get(constants.HeaderRequestID)).
		Msg("Making HTTP request")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response
	log.Debug().
		Str("method", method).
		Str("url", url).
		Int("status_code", resp.StatusCode).
		Str("request_id", req.Header.Get(constants.HeaderRequestID)).
		Msg("HTTP request completed")

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// setHeaders sets request headers
func (c *Client) setHeaders(req *http.Request, headers map[string]string) {
	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("microservice-client/%s", c.config.ServiceName))

	// Set configured default headers
	for key, value := range c.config.DefaultHeaders {
		req.Header.Set(key, value)
	}

	// Set provided headers (can override defaults)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// setMicroserviceHeaders sets microservice-specific headers from context
func (c *Client) setMicroserviceHeaders(req *http.Request, ctx context.Context) {
	// Extract request context from standard context
	reqCtx := sharedcontext.FromRequestContext(ctx)

	// Set request ID (generate if not present)
	if reqCtx.RequestID != "" {
		req.Header.Set(constants.HeaderRequestID, reqCtx.RequestID)
	} else {
		req.Header.Set(constants.HeaderRequestID, sharedcontext.GenerateRequestID())
	}

	// Set user context headers if available
	if reqCtx.UserID != uuid.Nil {
		req.Header.Set(constants.HeaderUserID, reqCtx.UserID.String())
	}
	if reqCtx.UserRole != 0 {
		req.Header.Set(constants.HeaderUserRole, fmt.Sprintf("%d", reqCtx.UserRole))
	}
	if reqCtx.UserEmail != "" {
		req.Header.Set(constants.HeaderUserEmail, reqCtx.UserEmail)
	}

	// Set service name
	req.Header.Set(constants.HeaderServiceName, c.config.ServiceName)
}

// UnmarshalResponse unmarshals JSON response body into target struct
func (r *HTTPResponse) UnmarshalResponse(target interface{}) error {
	return json.Unmarshal(r.Body, target)
}

// IsSuccess checks if response status code indicates success (2xx)
func (r *HTTPResponse) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError checks if response status code indicates client error (4xx)
func (r *HTTPResponse) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError checks if response status code indicates server error (5xx)
func (r *HTTPResponse) IsServerError() bool {
	return r.StatusCode >= 500
}
