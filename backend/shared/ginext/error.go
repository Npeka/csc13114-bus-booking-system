package ginext

import (
	"fmt"
	"net/http"
)

// Error represents a structured error with HTTP status code
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// NewError creates a new structured error
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithDetails creates a new structured error with details
func NewErrorWithDetails(code int, message, details string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error constructors
func NewBadRequestError(message string) *Error {
	return NewError(http.StatusBadRequest, message)
}

func NewUnauthorizedError(message string) *Error {
	return NewError(http.StatusUnauthorized, message)
}

func NewForbiddenError(message string) *Error {
	return NewError(http.StatusForbidden, message)
}

func NewNotFoundError(message string) *Error {
	return NewError(http.StatusNotFound, message)
}

func NewConflictError(message string) *Error {
	return NewError(http.StatusConflict, message)
}

func NewValidationError(message string) *Error {
	return NewError(http.StatusUnprocessableEntity, message)
}

func NewInternalServerError(message string) *Error {
	return NewError(http.StatusInternalServerError, message)
}
