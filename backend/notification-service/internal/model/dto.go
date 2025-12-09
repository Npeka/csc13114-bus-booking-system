package model

// SendOTPEmailRequest represents the request to send an OTP email
type SendOTPEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

// SendOTPEmailResponse represents the response after sending an OTP email
type SendOTPEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TripReminderRequest represents the request to send a trip reminder email
type TripReminderRequest struct {
	Email             string `json:"email" binding:"required,email"`
	PassengerName     string `json:"passenger_name" binding:"required"`
	BookingReference  string `json:"booking_reference" binding:"required"`
	DepartureLocation string `json:"departure_location" binding:"required"`
	Destination       string `json:"destination" binding:"required"`
	DepartureTime     string `json:"departure_time" binding:"required"`
	SeatNumbers       string `json:"seat_numbers" binding:"required"`
	BusPlate          string `json:"bus_plate" binding:"required"`
	PickupPoint       string `json:"pickup_point" binding:"required"`
	TicketLink        string `json:"ticket_link" binding:"required"`
}

// BookingConfirmationRequest represents the request to send booking confirmation email
type BookingConfirmationRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Name             string `json:"name" binding:"required"`
	BookingReference string `json:"booking_reference" binding:"required"`
	From             string `json:"from" binding:"required"`
	To               string `json:"to" binding:"required"`
	DepartureTime    string `json:"departure_time" binding:"required"`
	SeatNumbers      string `json:"seat_numbers" binding:"required"`
	TotalAmount      int    `json:"total_amount" binding:"required"`
	TicketLink       string `json:"ticket_link" binding:"required"`
}

// BookingFailureRequest represents the request to send booking failure email
type BookingFailureRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Name             string `json:"name" binding:"required"`
	BookingReference string `json:"booking_reference" binding:"required"`
	Reason           string `json:"reason" binding:"required"`
	From             string `json:"from" binding:"required"`
	To               string `json:"to" binding:"required"`
	DepartureTime    string `json:"departure_time" binding:"required"`
	BookingLink      string `json:"booking_link" binding:"required"`
}

// BookingPendingRequest represents the request to send booking pending email
// BookingPendingRequest represents the request to send booking pending email
type BookingPendingRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Name             string `json:"name" binding:"required"`
	BookingReference string `json:"booking_reference" binding:"required"`
	From             string `json:"from" binding:"required"`
	To               string `json:"to" binding:"required"`
	DepartureTime    string `json:"departure_time" binding:"required"`
	TotalAmount      int    `json:"total_amount" binding:"required"`
	PaymentLink      string `json:"payment_link" binding:"required"`
}

type NotificationType string

const (
	NotificationTypeOTP                 NotificationType = "OTP"
	NotificationTypeTripReminder        NotificationType = "TRIP_REMINDER"
	NotificationTypeBookingConfirmation NotificationType = "BOOKING_CONFIRMATION"
	NotificationTypeBookingFailure      NotificationType = "BOOKING_FAILURE"
	NotificationTypeBookingPending      NotificationType = "BOOKING_PENDING"
)

// GenericNotificationRequest represents a unified request for all notifications
type GenericNotificationRequest struct {
	Type    NotificationType       `json:"type" binding:"required"`
	Payload map[string]interface{} `json:"payload" binding:"required"`
}
