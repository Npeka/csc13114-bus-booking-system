package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
)

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr returns a pointer to an int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// TimePtr returns a pointer to a time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}

// GenerateUUID generates a new UUID v4
func GenerateUUID() string {
	return uuid.New().String()
}

// TrimSpaces trims leading and trailing spaces from a string
func TrimSpaces(s string) string {
	return strings.TrimSpace(s)
}

// IsEmptyOrWhitespace checks if a string is empty or contains only whitespace
func IsEmptyOrWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(str string) string {
	words := strings.Split(str, "_")
	result := strings.ToLower(words[0])

	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result += strings.ToUpper(string(words[i][0])) + strings.ToLower(words[i][1:])
		}
	}

	return result
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(str string) string {
	words := strings.Split(str, "_")
	var result string

	for _, word := range words {
		if len(word) > 0 {
			result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return result
}

// Capitalize capitalizes the first letter of a string
func Capitalize(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
}

// Contains checks if a slice contains a specific item
func Contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Remove removes an item from a slice
func Remove[T comparable](slice []T, item T) []T {
	result := make([]T, 0, len(slice))
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// Unique returns unique items from a slice
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Map applies a function to each element of a slice
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

// Filter filters a slice based on a predicate function
func Filter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// ParseInt safely parses a string to int with default value
func ParseInt(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

// ParseInt64 safely parses a string to int64 with default value
func ParseInt64(s string, defaultValue int64) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	return defaultValue
}

// ParseFloat safely parses a string to float64 with default value
func ParseFloat(s string, defaultValue float64) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return defaultValue
}

// ParseBool safely parses a string to bool with default value
func ParseBool(s string, defaultValue bool) bool {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return defaultValue
}

// FormatTime formats time using carbon
func FormatTime(t time.Time, format string) string {
	return carbon.CreateFromStdTime(t).Format(format)
}

// ParseTime parses time string using carbon
func ParseTime(timeStr, format string) (time.Time, error) {
	c := carbon.ParseByFormat(timeStr, format)
	if c.Error != nil {
		return time.Time{}, c.Error
	}
	return c.StdTime(), nil
}

// TimeToVietnameseString converts time to Vietnamese format
func TimeToVietnameseString(t time.Time) string {
	return carbon.CreateFromStdTime(t).SetLocale("vi").ToDateTimeString()
}

// NormalizePhoneNumber normalizes Vietnamese phone numbers
func NormalizePhoneNumber(phone string) string {
	// Remove spaces, dashes, and parentheses
	phone = regexp.MustCompile(`[\s\-\(\)]`).ReplaceAllString(phone, "")

	// Remove country code prefixes
	if strings.HasPrefix(phone, "+84") {
		phone = "0" + phone[3:]
	} else if strings.HasPrefix(phone, "84") && len(phone) >= 10 {
		phone = "0" + phone[2:]
	}

	return phone
}

// ValidateVietnamesePhoneNumber validates Vietnamese phone number format
func ValidateVietnamesePhoneNumber(phone string) bool {
	normalized := NormalizePhoneNumber(phone)

	// Mobile phone patterns
	mobilePatterns := []string{
		`^0(3[2-9]|5[6|8|9]|7[0|6-9]|8[1-5]|9[0-9])\d{7}$`,
	}

	// Landline patterns
	landlinePatterns := []string{
		`^0(1[2689]|2[0-9])\d{8}$`,
	}

	patterns := append(mobilePatterns, landlinePatterns...)

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, normalized); matched {
			return true
		}
	}

	return false
}

// MaskEmail masks an email address for privacy
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		return email
	}

	masked := username[:1] + strings.Repeat("*", len(username)-2) + username[len(username)-1:]
	return masked + "@" + domain
}

// MaskPhoneNumber masks a phone number for privacy
func MaskPhoneNumber(phone string) string {
	if len(phone) < 4 {
		return phone
	}

	return phone[:3] + strings.Repeat("*", len(phone)-6) + phone[len(phone)-3:]
}

// IsValidEmail validates email format using regex
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// TruncateString truncates a string to a maximum length
func TruncateString(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}
	return str[:maxLength] + "..."
}

// PadLeft pads a string with spaces on the left to reach the target length
func PadLeft(str string, length int, padChar rune) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(string(padChar), length-len(str)) + str
}

// PadRight pads a string with spaces on the right to reach the target length
func PadRight(str string, length int, padChar rune) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(string(padChar), length-len(str))
}

// SlugifyString creates a URL-friendly slug from a string
func SlugifyString(str string) string {
	// Convert to lowercase
	str = strings.ToLower(str)

	// Replace Vietnamese characters
	replacements := map[string]string{
		"à|á|ạ|ả|ã|â|ầ|ấ|ậ|ẩ|ẫ|ă|ằ|ắ|ặ|ẳ|ẵ": "a",
		"è|é|ẹ|ẻ|ẽ|ê|ề|ế|ệ|ể|ễ":             "e",
		"ì|í|ị|ỉ|ĩ": "i",
		"ò|ó|ọ|ỏ|õ|ô|ồ|ố|ộ|ổ|ỗ|ơ|ờ|ớ|ợ|ở|ỡ": "o",
		"ù|ú|ụ|ủ|ũ|ư|ừ|ứ|ự|ử|ữ":             "u",
		"ỳ|ý|ỵ|ỷ|ỹ": "y",
		"đ":         "d",
	}

	for pattern, replacement := range replacements {
		str = regexp.MustCompile(pattern).ReplaceAllString(str, replacement)
	}

	// Replace non-alphanumeric characters with hyphens
	str = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(str, "")
	str = regexp.MustCompile(`[-\s]+`).ReplaceAllString(str, "-")

	// Trim hyphens from start and end
	str = strings.Trim(str, "-")

	return str
}
