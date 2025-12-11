package constants

import "time"

// Booking timeout and expiration durations
const (
	// BookingDefaultTimeout is the default timeout duration for a booking (5 minutes)
	BookingDefaultTimeout = 5 * time.Minute

	// BookingPaymentTimeout is the payment timeout duration for a booking (5 minutes)
	BookingPaymentTimeout = 5 * time.Minute

	// BookingRetryGracePeriod is the grace period after expiration during which payment can be retried (60 minutes)
	BookingRetryGracePeriod = 60 * time.Minute

	// BookingExpirationGracePeriod is the grace period before actually expiring a booking (1 minute)
	BookingExpirationGracePeriod = 1 * time.Minute
)

// Notification and background task timeouts
const (
	// BackgroundTaskTimeout is the default timeout for background tasks like sending emails (1 minute)
	BackgroundTaskTimeout = 1 * time.Minute
)

// Trip reminder scheduling
const (
	// TripReminderBeforeDeparture is how long before departure to send reminder (2 hours)
	TripReminderBeforeDeparture = 2 * time.Hour
)

// Booking reference generation
const (
	// BookingReferencePrefix is the prefix for booking reference numbers
	BookingReferencePrefix = "BK"

	// BookingReferenceRandomLength is the length of random characters in booking reference
	BookingReferenceRandomLength = 4

	// BookingReferenceCharset is the charset used for random part of booking reference
	BookingReferenceCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Queue names
const (
	// QueueNameBookingExpiry is the queue name for booking expiration jobs
	QueueNameBookingExpiry = "booking_expiry"

	// QueueNameTripReminder is the queue name for trip reminder jobs
	QueueNameTripReminder = "trip_reminder"
)

// Date/Time formats
const (
	// DateFormatBookingReference is the date format used in booking references (YYMMDD)
	DateFormatBookingReference = "060102"

	// DateTimeFormatDisplay is the display format for date/time in emails and UI (15:04 02/01/2006)
	DateTimeFormatDisplay = "15:04 02/01/2006"
)

// URLs (should be configurable via environment variables in production)
const (
	// DefaultFrontendURL is the default frontend URL
	DefaultFrontendURL = "http://localhost:3000"
)
