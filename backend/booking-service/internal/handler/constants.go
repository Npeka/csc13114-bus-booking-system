package handler

// HTTP Status constants
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusInternalServerError = 500
)

// Error messages
const (
	ErrInvalidRequest     = "Invalid request format"
	ErrValidationFailed   = "Validation failed"
	ErrBookingNotFound    = "Booking not found"
	ErrUserNotFound       = "User not found"
	ErrTripNotFound       = "Trip not found"
	ErrSeatNotAvailable   = "Seat not available"
	ErrPaymentFailed      = "Payment processing failed"
	ErrFeedbackNotFound   = "Feedback not found"
	ErrUnauthorizedAccess = "Unauthorized access"
	ErrInternalServer     = "Internal server error"
)

// Validation constants
const (
	MinPageSize     = 1
	MaxPageSize     = 100
	DefaultPageSize = 10
	DefaultPage     = 1
)

// Booking status constants
const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCompleted = "completed"
	BookingStatusCancelled = "cancelled"
)

// Seat status constants
const (
	SeatStatusAvailable = "available"
	SeatStatusReserved  = "reserved"
	SeatStatusBooked    = "booked"
)

// Payment status constants
const (
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"
)

// Context keys
const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"
)

// Header constants
const (
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
	HeaderUserID        = "X-User-ID"
	HeaderRequestID     = "X-Request-ID"
)

// Content type constants
const (
	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
)

// Date format constants
const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02T15:04:05Z07:00"
)
