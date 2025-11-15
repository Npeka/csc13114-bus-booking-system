package ginext

import (
	"net/http"
)

// Response represents HTTP response structure
type Response struct {
	Code   int
	Header http.Header
	*GeneralBody
}

// GeneralBody represents the response body structure
type GeneralBody struct {
	Data         interface{} `json:"data,omitempty"`
	Message      string      `json:"message,omitempty"`
	Success      bool        `json:"success"`
	ErrorCode    string      `json:"error_code,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
	Error        interface{} `json:"error,omitempty"`
}

// WithResponseCode sets response code
func (g *GeneralBody) WithResponseCode(code string) *GeneralBody {
	g.ErrorCode = code
	return g
}

// ResponseOption function type for response options
type ResponseOption func(response *Response)

// NewResponse makes a new response with empty body
func NewResponse(code int, opts ...ResponseOption) *Response {
	r := &Response{
		Code:        code,
		GeneralBody: &GeneralBody{Success: code < 400},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// WithResponseCodeOption sets response code option
func WithResponseCodeOption(code string) ResponseOption {
	return func(response *Response) {
		response.GeneralBody.WithResponseCode(code)
	}
}

// NewResponseData makes a new response with body data
func NewResponseData(code int, data interface{}, message string, opts ...ResponseOption) *Response {
	r := &Response{
		Code: code,
		GeneralBody: &GeneralBody{
			Data:    data,
			Message: message,
			Success: code < 400,
		},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewBody creates a new general body
func NewBody(data interface{}, err interface{}) *GeneralBody {
	return &GeneralBody{
		Data:    data,
		Error:   err,
		Success: err == nil,
	}
}

// Success response helpers
func NewSuccessResponse(data interface{}, message string) *Response {
	return NewResponseData(http.StatusOK, data, message)
}

func NewCreatedResponse(data interface{}, message string) *Response {
	return NewResponseData(http.StatusCreated, data, message)
}

func NewNoContentResponse() *Response {
	return NewResponse(http.StatusNoContent)
}
