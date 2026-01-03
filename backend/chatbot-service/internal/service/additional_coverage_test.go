package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bus-booking/chatbot-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Additional comprehensive tests for HTTP clients to reach 70% coverage

func TestSearchTrips_ErrorBodyReadFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100") // Lie about content length
		w.WriteHeader(http.StatusBadRequest)
		// Write less than promised to cause read issues
	}))
	defer server.Close()

	client := NewTripServiceClient(server.URL)
	params := &model.TripSearchParams{
		Origin:      "Hà Nội",
		Destination: "Đà Nẵng",
	}

	result, err := client.SearchTrips(context.Background(), params)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSearchTrips_EmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"data": []any{},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewTripServiceClient(server.URL)
	params := &model.TripSearchParams{
		Origin:      "Hà Nội",
		Destination: "Unknown City",
	}

	result, err := client.SearchTrips(context.Background(), params)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetTripByID_ErrorBodyReadFailure(t *testing.T) {
	tripID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewTripServiceClient(server.URL)
	result, err := client.GetTripByID(context.Background(), tripID.String())

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTripByID_EmptyResponse(t *testing.T) {
	tripID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	client := NewTripServiceClient(server.URL)
	result, err := client.GetTripByID(context.Background(), tripID.String())

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCreateGuestBooking_MinimalFields(t *testing.T) {
	tripID := uuid.New()
	seatID := uuid.New()
	bookingID := uuid.New()

	expectedBooking := model.BookingResponse{
		ID:         bookingID,
		Reference:  "MIN123",
		TripID:     tripID,
		TotalPrice: 100000,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := model.APIResponse[model.BookingResponse]{
			Data: expectedBooking,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewBookingServiceClient(server.URL)
	request := &model.CreateGuestBookingRequest{
		TripID:   tripID,
		SeatIDs:  []uuid.UUID{seatID},
		FullName: "A",
		Email:    "a@b.c",
		Phone:    "0",
	}

	result, err := client.CreateGuestBooking(context.Background(), request)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "MIN123", result.Reference)
}

func TestCreateGuestBooking_ErrorBodyReadFailure(t *testing.T) {
	tripID := uuid.New()
	seatID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewBookingServiceClient(server.URL)
	request := &model.CreateGuestBookingRequest{
		TripID:   tripID,
		SeatIDs:  []uuid.UUID{seatID},
		FullName: "Test",
		Email:    "test@example.com",
		Phone:    "0901234567",
	}

	result, err := client.CreateGuestBooking(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetBookingByID_ErrorBodyReadFailure(t *testing.T) {
	bookingID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewBookingServiceClient(server.URL)
	result, err := client.GetBookingByID(context.Background(), bookingID.String())

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetBookingByID_AllFieldsPopulated(t *testing.T) {
	bookingID := uuid.New()
	tripID := uuid.New()
	expiresAt := time.Now().Add(15 * time.Minute)

	expectedBooking := model.BookingResponse{
		ID:         bookingID,
		Reference:  "FULL123XYZ",
		Status:     "confirmed",
		TotalPrice: 500000,
		TripID:     tripID,
		CreatedAt:  time.Now(),
		ExpiresAt:  &expiresAt,
		Transaction: &model.TransactionInfo{
			ID:     uuid.New(),
			Status: "completed",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := model.APIResponse[model.BookingResponse]{
			Data: expectedBooking,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewBookingServiceClient(server.URL)
	result, err := client.GetBookingByID(context.Background(), bookingID.String())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "FULL123XYZ", result.Reference)
	assert.NotNil(t, result.Transaction)
}

func TestCreateTransaction_ErrorBodyReadFailure(t *testing.T) {
	bookingID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	client := NewPaymentServiceClient(server.URL)
	request := &model.CreateTransactionRequest{
		BookingID:     bookingID,
		Amount:        500000,
		Currency:      "VND",
		PaymentMethod: "payos",
	}

	result, err := client.CreateTransaction(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetBookingByReference_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResponse := model.APIResponse[interface{}]{
			Error: &model.APIError{
				Message: "Booking not found with provided credentials",
				Code:    "NOT_FOUND",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}))
	defer server.Close()

	client := NewBookingServiceClient(server.URL)
	result, err := client.GetBookingByReference(context.Background(), "NOTFOUND", "wrong@email.com")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "404")
}

// Tests for GetFAQAnswer to increase coverage
func TestGetFAQAnswer_DirectKeywordMatch(t *testing.T) {
	service := &ChatbotServiceImpl{
		faqKnowledge: []model.FAQ{
			{
				Question: "Chính sách hủy vé như thế nào?",
				Answer:   "Hoàn 70% trước 24 giờ",
				Keywords: []string{"hủy", "cancel", "hoàn tiền"},
			},
			{
				Question: "Hành lý được bao nhiêu kg?",
				Answer:   "Tối đa 20kg miễn phí",
				Keywords: []string{"hành lý", "luggage", "kg"},
			},
		},
	}

	tests := []struct {
		name           string
		question       string
		expectedAnswer string
	}{
		{
			name:           "Match cancel keyword",
			question:       "Tôi muốn hủy vé",
			expectedAnswer: "Hoàn 70% trước 24 giờ",
		},
		{
			name:           "Match luggage keyword",
			question:       "Hành lý được mang bao nhiêu?",
			expectedAnswer: "Tối đa 20kg miễn phí",
		},
		{
			name:           "Match with uppercase",
			question:       "CHÍNH SÁCH HỦY",
			expectedAnswer: "Hoàn 70% trước 24 giờ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			answer, err := service.GetFAQAnswer(context.Background(), tt.question)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedAnswer, answer)
		})
	}
}

// Tests for normalizeCityName with actual cityAliases map
func TestNormalizeCityName_ActualMappings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Sài Gòn to TP.HCM", "Sài Gòn", "TP. Hồ Chí Minh"},
		{"Saigon case-insensitive", "saigon", "TP. Hồ Chí Minh"},
		{"SG to TP.HCM", "SG", "TP. Hồ Chí Minh"},
		{"Hanoi to Hà Nội", "Hanoi", "Hà Nội"},
		{"Ha Noi to Hà Nội", "Ha Noi", "Hà Nội"},
		{"Dalat to Đà Lạt", "Dalat", "Đà Lạt"},
		{"Da Nang to Đà Nẵng", "Da Nang", "Đà Nẵng"},
		{"Unknown city returns same", "Unknown City", "Unknown City"},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeCityName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
