package ginext

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bus-booking/shared/validator"
)

type Handler func(r *Request) (*Response, error)

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

func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *Error:
		c.JSON(e.Code, &GeneralBody{
			Error: &ErrorBody{Message: e.Message},
		})
	default:
		log.Error().Err(err).Msg("Unhandled error in handler")
		c.JSON(http.StatusInternalServerError, &GeneralBody{
			Error: &ErrorBody{Message: "Internal server error"},
		})
	}
}

func ValidateRequest(req interface{}) error {
	if validator.ValidatorInstance == nil {
		return NewInternalServerError("Validator not initialized")
	}

	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(req); len(validationErrors) > 0 {
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
