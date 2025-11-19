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
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

// ErrorBody represents simplified error structure
type ErrorBody struct {
	Message string `json:"message"`
}

// WithError sets error message
func (g *GeneralBody) WithError(message string) *GeneralBody {
	g.Error = &ErrorBody{Message: message}
	return g
}

// ResponseOption function type for response options
type ResponseOption func(response *Response)

// NewResponse makes a new response with empty body
func NewResponse(code int, opts ...ResponseOption) *Response {
	r := &Response{
		Code:        code,
		GeneralBody: &GeneralBody{},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// WithErrorOption sets error message option
func WithErrorOption(message string) ResponseOption {
	return func(response *Response) {
		response.GeneralBody.WithError(message)
	}
}

// NewResponseData makes a new response with body data
func NewResponseData(code int, data interface{}, message string, opts ...ResponseOption) *Response {
	r := &Response{
		Code: code,
		GeneralBody: &GeneralBody{
			Data:    data,
			Message: message,
		},
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewBody creates a new general body
func NewBody(data interface{}, errMsg string) *GeneralBody {
	var errorBody *ErrorBody
	if errMsg != "" {
		errorBody = &ErrorBody{Message: errMsg}
	}
	return &GeneralBody{
		Data:  data,
		Error: errorBody,
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

// Error response helpers
func NewErrorResponse(code int, message string) *Response {
	return &Response{
		Code: code,
		GeneralBody: &GeneralBody{
			Error: &ErrorBody{Message: message},
		},
	}
}

func NewBadRequestResponse(message string) *Response {
	return NewErrorResponse(http.StatusBadRequest, message)
}

func NewUnauthorizedResponse(message string) *Response {
	return NewErrorResponse(http.StatusUnauthorized, message)
}

func NewForbiddenResponse(message string) *Response {
	return NewErrorResponse(http.StatusForbidden, message)
}

func NewNotFoundResponse(message string) *Response {
	return NewErrorResponse(http.StatusNotFound, message)
}

func NewConflictResponse(message string) *Response {
	return NewErrorResponse(http.StatusConflict, message)
}

func NewValidationErrorResponse(message string) *Response {
	return NewErrorResponse(http.StatusUnprocessableEntity, message)
}

func NewInternalServerErrorResponse(message string) *Response {
	return NewErrorResponse(http.StatusInternalServerError, message)
}
