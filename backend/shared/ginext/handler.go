package ginext

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/shared/validator"
)

// Handler defines the signature for wrapped handlers
type Handler func(r *Request) (*Response, error)

// WrapHandler wraps a handler function with error handling and response formatting
func WrapHandler(handler Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err  error
			resp *Response
		)

		defer func() {
			if r := recover(); r != nil {
				// Handle panic - convert to error
				switch v := r.(type) {
				case *Error:
					err = v
				case error:
					err = NewInternalServerError(v.Error())
				default:
					err = NewInternalServerError(fmt.Sprintf("panic: %v", v))
				}
			}

			if err != nil {
				handleError(c, err)
				return
			}

			if resp == nil {
				return
			}

			// Set headers if any
			for k, v := range resp.Header {
				for _, v_ := range v {
					c.Header(k, v_)
				}
			}

			// Send response
			if resp.GeneralBody != nil && (resp.Data != nil || resp.Error != nil || resp.Message != "") {
				c.JSON(resp.Code, resp.GeneralBody)
			} else {
				c.Status(resp.Code)
			}
		}()

		req := NewRequest(c)
		resp, err = handler(req)
	}
}

// handleError handles different types of errors and sends appropriate response
func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *Error:
		// Structured error with known status code
		c.JSON(e.Code, &GeneralBody{
			Success:      false,
			ErrorMessage: e.Message,
			ErrorCode:    fmt.Sprintf("HTTP_%d", e.Code),
			Error:        e.Details,
		})
	default:
		// Unknown error - treat as internal server error
		log.Error().Err(err).Msg("Unhandled error in handler")
		c.JSON(http.StatusInternalServerError, &GeneralBody{
			Success:      false,
			ErrorMessage: "Internal server error",
			ErrorCode:    "INTERNAL_SERVER_ERROR",
		})
	}
}

// ValidateRequest validates a request struct using the global validator
func ValidateRequest(req interface{}) error {
	if validator.ValidatorInstance == nil {
		return NewInternalServerError("Validator not initialized")
	}

	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(req); len(validationErrors) > 0 {
		// Format validation errors into a readable message
		var errorMsg string
		for i, ve := range validationErrors {
			if i > 0 {
				errorMsg += ", "
			}
			errorMsg += fmt.Sprintf("%s: %s", ve.Field, ve.Message)
		}
		return NewValidationError("Validation failed: " + errorMsg)
	}

	return nil
}
