package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"bus-booking/shared/constants"
	sharedcontext "bus-booking/shared/context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// HTTPClient interface for making HTTP requests to other microservices
type HTTPClient interface {
	Get(ctx context.Context, url string, params map[string][]string, headers map[string]string) (*HTTPResponse, error)
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

// StandardResponse is the standard API response format {data: T, meta: M}
type StandardResponse[T any, M any] struct {
	Data T `json:"data"`
	Meta M `json:"meta,omitempty"`
}

// DataResponse is a simple response with only data field
type DataResponse[T any] struct {
	Data T `json:"data"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Config for HTTP client
type Config struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
	ServiceName    string
	BaseURL        string
	DefaultHeaders map[string]string
}

// Client implements HTTPClient
type Client struct {
	httpClient *http.Client
	config     *Config
	baseURL    string
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
		config:  config,
		baseURL: config.BaseURL,
	}
}

// Get makes a GET request
func (c *Client) Get(ctx context.Context, path string, params map[string][]string, headers map[string]string) (*HTTPResponse, error) {
	fullURL := c.buildURL(path)

	// Add query parameters if provided
	if len(params) > 0 {
		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}

		q := parsedURL.Query()
		for key, values := range params {
			for _, value := range values {
				q.Add(key, value)
			}
		}
		parsedURL.RawQuery = q.Encode()
		fullURL = parsedURL.String()
	}

	return c.doRequest(ctx, http.MethodGet, fullURL, nil, headers)
}

// Post makes a POST request
func (c *Client) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPost, c.buildURL(url), body, headers)
}

// Put makes a PUT request
func (c *Client) Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPut, c.buildURL(url), body, headers)
}

// Delete makes a DELETE request
func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodDelete, c.buildURL(url), nil, headers)
}

// Patch makes a PATCH request
func (c *Client) Patch(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPatch, c.buildURL(url), body, headers)
}

// buildURL builds the full URL by combining baseURL with the given path
// If path is already a full URL (starts with http:// or https://), returns it as-is
// Otherwise, prepends the baseURL to the path
func (c *Client) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}
	// If path is already a full URL, return as-is
	if len(path) >= 7 && (path[:7] == "http://" || path[:8] == "https://") {
		return path
	}
	// Combine baseURL with path
	return c.baseURL + path
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
		Str("request_id", req.Header.Get(constants.XRequestID)).
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
		Str("request_id", req.Header.Get(constants.XRequestID)).
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
		req.Header.Set(constants.XRequestID, reqCtx.RequestID)
	} else {
		req.Header.Set(constants.XRequestID, sharedcontext.GenerateRequestID())
	}

	// Set user context headers if available
	if reqCtx.UserID != uuid.Nil {
		req.Header.Set(constants.XUserID, reqCtx.UserID.String())
	}
	if reqCtx.UserRole != 0 {
		req.Header.Set(constants.XUserRole, fmt.Sprintf("%d", reqCtx.UserRole))
	}
	if reqCtx.UserEmail != "" {
		req.Header.Set(constants.XUserEmail, reqCtx.UserEmail)
	}

	// Set service name
	req.Header.Set(constants.XServiceName, c.config.ServiceName)
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

// ParseData parses response with {data: T} format
func ParseData[T any](resp *HTTPResponse) (*T, error) {
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var result DataResponse[T]
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result.Data, nil
}

// ParseListData parses response with {data: {items_field: []T}} format and extracts the list
// This is useful when the API returns {data: {seats: [...]}} and you want to get the slice directly
func ParseListData[T any](resp *HTTPResponse) ([]T, error) {
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var result DataResponse[[]T]
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}

// ParseDataWithMeta parses response with {data: T, meta: M} format
func ParseDataWithMeta[T any, M any](resp *HTTPResponse) (*T, *M, error) {
	if !resp.IsSuccess() {
		return nil, nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var result StandardResponse[T, M]
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result.Data, &result.Meta, nil
}

// ParseDataWithPagination parses response with {data: T, meta: PaginationMeta} format
func ParseDataWithPagination[T any](resp *HTTPResponse) (*T, *PaginationMeta, error) {
	return ParseDataWithMeta[T, PaginationMeta](resp)
}
