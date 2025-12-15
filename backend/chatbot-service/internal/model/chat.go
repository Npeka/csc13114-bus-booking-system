package model

import (
	"time"

	"github.com/google/uuid"
)

// ChatMessage represents a single message in a conversation
type ChatMessage struct {
	Role    string `json:"role"`    // "user", "assistant", "system"
	Content string `json:"content"` // Message content
}

// ChatRequest represents a chatbot request from the user
type ChatRequest struct {
	SessionID string        `json:"session_id,omitempty"` // Optional session tracking
	Message   string        `json:"message" binding:"required"`
	Context   *ChatContext  `json:"context,omitempty"` // Optional context
	History   []ChatMessage `json:"history,omitempty"` // Conversation history
}

// ChatContext holds conversation context
type ChatContext struct {
	UserID           *uuid.UUID        `json:"user_id,omitempty"`
	CurrentStep      string            `json:"current_step,omitempty"` // "search", "select_trip", "select_seats", "collect_info", "booking", "payment"
	TripSearchID     *uuid.UUID        `json:"trip_search_id,omitempty"`
	SelectedTrip     *TripInfo         `json:"selected_trip,omitempty"`
	SelectedSeats    []SelectedSeat    `json:"selected_seats,omitempty"`    // NEW: Seats user wants to book
	PassengerDetails []PassengerDetail `json:"passenger_details,omitempty"` // NEW: Collected passenger info
	BookingID        *uuid.UUID        `json:"booking_id,omitempty"`        // NEW: Created booking ID
	BookingReference string            `json:"booking_reference,omitempty"` // NEW: Booking reference code
}

// SelectedSeat represents a seat selected by the user
type SelectedSeat struct {
	SeatNumber string    `json:"seat_number"` // e.g., "A1"
	SeatID     uuid.UUID `json:"seat_id"`     // UUID from trip service
}

// PassengerDetail represents collected passenger information
type PassengerDetail struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	SeatNumber string `json:"seat_number"` // Which seat this passenger is assigned to
}

// TripInfo holds basic trip information for context
type TripInfo struct {
	TripID        uuid.UUID `json:"trip_id"`
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	DepartureTime time.Time `json:"departure_time"`
	Price         float64   `json:"price"`
}

// ChatResponse represents the chatbot's response
type ChatResponse struct {
	Message     string       `json:"message"`
	Intent      string       `json:"intent,omitempty"`      // "search_trip", "book_trip", "faq", "other"
	Action      string       `json:"action,omitempty"`      // "show_trips", "show_booking_form", etc.
	Data        interface{}  `json:"data,omitempty"`        // Additional data (trip results, booking info, etc.)
	Context     *ChatContext `json:"context,omitempty"`     // Updated context
	Suggestions []string     `json:"suggestions,omitempty"` // Quick reply suggestions
}

// TripSearchParams represents extracted search parameters
type TripSearchParams struct {
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	DepartureDate time.Time `json:"departure_date,omitempty"`
	Passengers    int       `json:"passengers,omitempty"`
}

// FAQ represents frequently asked questions
type FAQ struct {
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Keywords []string `json:"keywords"` // For matching
}
