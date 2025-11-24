package validator

import (
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"bus-booking/shared/response"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// ValidatorInstance holds the global validator instance
var ValidatorInstance *CustomValidator

func MustSetupValidator() {
	InitValidator()
}

// InitValidator initializes the custom validator with custom validation rules and tag name func
func InitValidator() {
	v := validator.New()

	// Register custom tag name function to use json tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validation rules
	registerCustomValidations(v)

	ValidatorInstance = &CustomValidator{validator: v}

	// Set the validator for Gin binding
	binding.Validator = ValidatorInstance
}

// Engine returns the underlying validator engine
func (cv *CustomValidator) Engine() interface{} {
	return cv.validator
}

// ValidateStruct validates a struct and returns error (for gin binding)
func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	return cv.validator.Struct(obj)
}

// ValidateStructDetailed validates a struct and returns formatted validation errors
func (cv *CustomValidator) ValidateStructDetailed(obj interface{}) []response.ValidationError {
	var validationErrors []response.ValidationError

	if err := cv.validator.Struct(obj); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, response.ValidationError{
				Field:   toSnakeCase(err.Field()),
				Message: getErrorMessage(err),
				Value:   err.Value().(string),
			})
		}
	}

	return validationErrors
}

// registerCustomValidations registers custom validation rules
func registerCustomValidations(v *validator.Validate) {
	// Register phone number validation
	v.RegisterValidation("phone", validatePhone)

	// Register password strength validation
	v.RegisterValidation("password", validatePassword)

	// Register Vietnamese phone number validation
	v.RegisterValidation("vnphone", validateVNPhone)

	// Register ID card validation (Vietnamese)
	v.RegisterValidation("vnid", validateVNID)

	// Register license plate validation (Vietnamese)
	v.RegisterValidation("vnplate", validateVNPlate)

	// Register no special chars validation
	v.RegisterValidation("nospecial", validateNoSpecialChars)

	// Register alphanumeric with spaces validation
	v.RegisterValidation("alphanumspc", validateAlphaNumSpace)

	// Register seat code validation (A1, B2, etc.)
	v.RegisterValidation("seatcode", validateSeatCode)

	// Register time format validation (HH:MM)
	v.RegisterValidation("timeformat", validateTimeFormat)

	// Register date format validation (YYYY-MM-DD)
	v.RegisterValidation("dateformat", validateDateFormat)

	// Register enum validation
	v.RegisterValidation("enum", validateEnum)
}

// validatePhone validates phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}

	// Remove common prefixes and check if remaining is numeric
	phone = strings.TrimPrefix(phone, "+")
	phone = strings.TrimPrefix(phone, "84")
	phone = strings.TrimPrefix(phone, "0")

	for _, r := range phone {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateVNPhone validates Vietnamese phone number
func validateVNPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Remove spaces and dashes
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Check for valid Vietnamese phone patterns
	patterns := []string{
		"^(\\+84|84|0)(3[2-9]|5[6|8|9]|7[0|6-9]|8[1-5]|9[0-9])[0-9]{7}$",
		"^(\\+84|84|0)(1[2689]|2[0-9])[0-9]{8}$", // Landline
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, phone); matched {
			return true
		}
	}

	return false
}

// validateVNID validates Vietnamese ID card number
func validateVNID(fl validator.FieldLevel) bool {
	id := fl.Field().String()

	// Remove spaces
	id = strings.ReplaceAll(id, " ", "")

	// Old format: 9 digits or New format: 12 digits
	if len(id) != 9 && len(id) != 12 {
		return false
	}

	// Check if all characters are digits
	for _, r := range id {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// validateVNPlate validates Vietnamese license plate
func validateVNPlate(fl validator.FieldLevel) bool {
	plate := fl.Field().String()

	// Remove spaces and dashes
	plate = strings.ReplaceAll(plate, " ", "")
	plate = strings.ReplaceAll(plate, "-", "")

	// Vietnamese license plate patterns
	patterns := []string{
		"^[0-9]{2}[A-Z]{1,2}[0-9]{4,5}$",          // Standard format
		"^[0-9]{2}[A-Z]{1,2}[0-9]{3}\\.[0-9]{2}$", // With dot separator
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, strings.ToUpper(plate)); matched {
			return true
		}
	}

	return false
}

// validateNoSpecialChars validates that field contains no special characters
func validateNoSpecialChars(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' {
			return false
		}
	}

	return true
}

// validateAlphaNumSpace validates alphanumeric characters with spaces (allows multiple consecutive spaces)
func validateAlphaNumSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Allow empty string
	if value == "" {
		return true
	}

	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' {
			return false
		}
	}

	return true
}

// validateSeatCode validates seat code format (A1, B2, etc.)
func validateSeatCode(fl validator.FieldLevel) bool {
	seatCode := fl.Field().String()

	pattern := "^[A-Z][0-9]{1,2}$"
	matched, _ := regexp.MatchString(pattern, strings.ToUpper(seatCode))

	return matched
}

// validateTimeFormat validates time format (HH:MM)
func validateTimeFormat(fl validator.FieldLevel) bool {
	timeStr := fl.Field().String()

	pattern := "^([01][0-9]|2[0-3]):[0-5][0-9]$"
	matched, _ := regexp.MatchString(pattern, timeStr)

	return matched
}

// validateDateFormat validates date format (YYYY-MM-DD)
func validateDateFormat(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	pattern := "^[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$"
	matched, _ := regexp.MatchString(pattern, dateStr)

	return matched
}

// validateEnum validates enum values
func validateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := fl.Param()

	allowedValues := strings.Split(param, " ")

	for _, allowed := range allowedValues {
		if value == allowed {
			return true
		}
	}

	return false
}

// getErrorMessage returns a human-readable error message for validation errors
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "phone":
		return "Invalid phone number format"
	case "password":
		return "Password must be at least 8 characters with uppercase, lowercase, number and special character"
	case "vnphone":
		return "Invalid Vietnamese phone number"
	case "vnid":
		return "Invalid Vietnamese ID card number"
	case "vnplate":
		return "Invalid Vietnamese license plate"
	case "min":
		return "Value must be at least " + fe.Param()
	case "max":
		return "Value must be at most " + fe.Param()
	case "len":
		return "Value must be exactly " + fe.Param() + " characters"
	case "oneof":
		return "Value must be one of: " + fe.Param()
	case "url":
		return "Invalid URL format"
	case "uuid":
		return "Invalid UUID format"
	case "numeric":
		return "Value must be numeric"
	case "alpha":
		return "Value must contain only letters"
	case "alphanum":
		return "Value must contain only letters and numbers"
	case "nospecial":
		return "Value must not contain special characters"
	case "alphanumspc":
		return "Value must contain only letters, numbers and spaces"
	case "seatcode":
		return "Invalid seat code format (e.g., A1, B12)"
	case "timeformat":
		return "Invalid time format (HH:MM)"
	case "dateformat":
		return "Invalid date format (YYYY-MM-DD)"
	case "enum":
		return "Invalid value. Allowed values: " + fe.Param()
	case "gt":
		return "Value must be greater than " + fe.Param()
	case "gte":
		return "Value must be greater than or equal to " + fe.Param()
	case "lt":
		return "Value must be less than " + fe.Param()
	case "lte":
		return "Value must be less than or equal to " + fe.Param()
	default:
		return "Invalid value"
	}
}

// toSnakeCase converts camelCase to snake_case
func toSnakeCase(str string) string {
	var result []rune

	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}
