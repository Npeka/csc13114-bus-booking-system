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

func TestTripServiceClient_SearchTrips(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tripID := uuid.New()
		expectedResponse := map[string]any{
			"data": []map[string]any{
				{
					"id":             tripID.String(),
					"origin":         "Hà Nội",
					"destination":    "Đà Nẵng",
					"departure_time": "2026-01-15T08:00:00Z",
					"price":          300000,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/trips/search", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			// Verify query parameters
			assert.Contains(t, r.URL.RawQuery, "origin=")
			assert.Contains(t, r.URL.RawQuery, "destination=")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expectedResponse)
		}))
		defer server.Close()

		client := NewTripServiceClient(server.URL)
		params := &model.TripSearchParams{
			Origin:        "Hà Nội",
			Destination:   "Đà Nẵng",
			DepartureDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		}

		result, err := client.SearchTrips(context.Background(), params)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("HTTP Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
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
		assert.Contains(t, err.Error(), "status 500")
	})

	t.Run("Invalid JSON Response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
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
	})

	t.Run("Error Response With Body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("validation error: invalid parameters"))
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
		assert.Contains(t, err.Error(), "400")
	})
}

func TestTripServiceClient_GetTripByID(t *testing.T) {
	tripID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		expectedTrip := model.TripDetailResponse{
			ID:             tripID,
			AvailableSeats: 20,
			TotalSeats:     45,
			Route: &model.RouteDetail{
				Origin:      "Sài Gòn",
				Destination: "Đà Lạt",
			},
			Bus: &model.BusDetail{
				ID:           uuid.New(),
				LicensePlate: "51A-12345",
				TotalSeats:   45,
				Seats: []model.SeatDetail{
					{
						ID:         uuid.New(),
						SeatNumber: "A1",
						SeatType:   "normal",
					},
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/trips/"+tripID.String(), r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			// Verify preload query params
			assert.Contains(t, r.URL.RawQuery, "preload_bus=true")
			assert.Contains(t, r.URL.RawQuery, "preload_seat=true")
			assert.Contains(t, r.URL.RawQuery, "seat_booking_status=true")

			response := model.APIResponse[model.TripDetailResponse]{
				Data: expectedTrip,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewTripServiceClient(server.URL)
		result, err := client.GetTripByID(context.Background(), tripID.String())

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, tripID, result.ID)
		assert.Equal(t, 20, result.AvailableSeats)
		assert.NotNil(t, result.Bus)
		assert.Len(t, result.Bus.Seats, 1)
	})

	t.Run("Trip Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewTripServiceClient(server.URL)
		result, err := client.GetTripByID(context.Background(), tripID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("Invalid JSON Response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{invalid json}"))
		}))
		defer server.Close()

		client := NewTripServiceClient(server.URL)
		result, err := client.GetTripByID(context.Background(), tripID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "decode")
	})
}

func TestBookingServiceClient_CreateGuestBooking(t *testing.T) {
	tripID := uuid.New()
	seatID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		bookingID := uuid.New()
		expectedBooking := model.BookingResponse{
			ID:         bookingID,
			Reference:  "ABC123XYZ",
			TripID:     tripID,
			TotalPrice: 500000,
			Status:     "pending",
			CreatedAt:  time.Now(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/bookings/guest", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Decode and verify request body
			var req model.CreateGuestBookingRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			assert.Equal(t, tripID, req.TripID)
			assert.Len(t, req.SeatIDs, 1)

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
			FullName: "Test User",
			Email:    "test@example.com",
			Phone:    "0901234567",
			Passengers: []model.PassengerData{
				{
					Name:   "Test User",
					Email:  "test@example.com",
					Phone:  "0901234567",
					SeatID: seatID,
				},
			},
		}

		result, err := client.CreateGuestBooking(context.Background(), request)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ABC123XYZ", result.Reference)
		assert.Equal(t, float64(500000), result.TotalPrice)
	})

	t.Run("Validation Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorResponse := model.APIResponse[interface{}]{
				Error: &model.APIError{
					Message: "Invalid seat ID",
					Code:    "VALIDATION_ERROR",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
		}))
		defer server.Close()

		client := NewBookingServiceClient(server.URL)
		request := &model.CreateGuestBookingRequest{
			TripID:   tripID,
			SeatIDs:  []uuid.UUID{uuid.Nil}, // Invalid seat
			FullName: "Test",
			Email:    "test@example.com",
			Phone:    "0901234567",
		}

		result, err := client.CreateGuestBooking(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Invalid seat ID")
	})

	t.Run("JSON Decode Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("{malformed json}"))
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
	})
}

func TestBookingServiceClient_GetBookingByReference(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		bookingID := uuid.New()
		expectedBooking := model.BookingResponse{
			ID:         bookingID,
			Reference:  "ABC123XYZ",
			Status:     "confirmed",
			TotalPrice: 300000,
			CreatedAt:  time.Now(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/bookings/lookup", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			// Verify query parameters
			assert.Equal(t, "ABC123XYZ", r.URL.Query().Get("reference"))
			assert.Equal(t, "test@example.com", r.URL.Query().Get("email"))

			response := model.APIResponse[model.BookingResponse]{
				Data: expectedBooking,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewBookingServiceClient(server.URL)
		result, err := client.GetBookingByReference(context.Background(), "ABC123XYZ", "test@example.com")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ABC123XYZ", result.Reference)
		assert.Equal(t, "confirmed", result.Status)
	})

	t.Run("Booking Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewBookingServiceClient(server.URL)
		result, err := client.GetBookingByReference(context.Background(), "INVALID", "test@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("JSON Decode Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{broken json"))
		}))
		defer server.Close()

		client := NewBookingServiceClient(server.URL)
		result, err := client.GetBookingByReference(context.Background(), "ABC123", "test@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "decode")
	})
}

func TestPaymentServiceClient_CreateTransaction(t *testing.T) {
	bookingID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		transactionID := uuid.New()
		expectedTransaction := model.TransactionResponse{
			ID:            transactionID,
			BookingID:     bookingID,
			Amount:        500000,
			Currency:      "VND",
			PaymentMethod: "payos",
			Status:        "pending",
			CheckoutURL:   "https://payment.example.com/checkout",
			CreatedAt:     time.Now(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/transactions", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Decode and verify request body
			var req model.CreateTransactionRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			assert.Equal(t, bookingID, req.BookingID)
			assert.Equal(t, 500000, req.Amount)

			response := model.APIResponse[model.TransactionResponse]{
				Data: expectedTransaction,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response)
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

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, bookingID, result.BookingID)
		assert.Equal(t, "pending", result.Status)
		assert.NotEmpty(t, result.CheckoutURL)
	})

	t.Run("Payment Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorResponse := model.APIResponse[interface{}]{
				Error: &model.APIError{
					Message: "Payment provider unavailable",
					Code:    "PAYMENT_ERROR",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(errorResponse)
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
		assert.Contains(t, err.Error(), "Payment provider unavailable")
	})

	t.Run("Invalid Amount", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorResponse := model.APIResponse[interface{}]{
				Error: &model.APIError{
					Message: "Amount must be positive",
					Code:    "INVALID_AMOUNT",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
		}))
		defer server.Close()

		client := NewPaymentServiceClient(server.URL)
		request := &model.CreateTransactionRequest{
			BookingID:     bookingID,
			Amount:        -100,
			Currency:      "VND",
			PaymentMethod: "payos",
		}

		result, err := client.CreateTransaction(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Amount must be positive")
	})

	t.Run("JSON Decode Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("{{invalid}}"))
		}))
		defer server.Close()

		client := NewPaymentServiceClient(server.URL)
		request := &model.CreateTransactionRequest{
			BookingID:     bookingID,
			Amount:        100000,
			Currency:      "VND",
			PaymentMethod: "payos",
		}

		result, err := client.CreateTransaction(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
