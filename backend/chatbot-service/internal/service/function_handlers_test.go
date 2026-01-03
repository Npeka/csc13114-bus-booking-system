package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bus-booking/chatbot-service/internal/model"
	"bus-booking/chatbot-service/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingServiceClient_GetBookingByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		bookingID := uuid.New()
		expectedBooking := model.BookingResponse{
			ID:         bookingID,
			Reference:  "ABC123XYZ",
			Status:     "confirmed",
			TotalPrice: 350000,
			CreatedAt:  time.Now(),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/v1/bookings/"+bookingID.String(), r.URL.Path)
			assert.Equal(t, "GET", r.Method)

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
		assert.Equal(t, bookingID, result.ID)
		assert.Equal(t, "ABC123XYZ", result.Reference)
	})

	t.Run("Booking Not Found", func(t *testing.T) {
		bookingID := uuid.New()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewBookingServiceClient(server.URL)
		result, err := client.GetBookingByID(context.Background(), bookingID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "404")
	})
}

func TestHandleSearchTrips(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTripService := mocks.NewMockTripServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			tripService: mockTripService,
		}

		args := map[string]any{
			"origin":         "Hà Nội",
			"destination":    "Đà Nẵng",
			"departure_date": "2026-01-15",
			"passengers":     float64(2),
		}

		expectedTrips := map[string]any{
			"data": []map[string]any{
				{
					"id":             uuid.New().String(),
					"origin":         "Hà Nội",
					"destination":    "Đà Nẵng",
					"departure_time": "2026-01-15T08:00:00Z",
					"price":          300000,
				},
			},
		}

		mockTripService.EXPECT().
			SearchTrips(gomock.Any(), gomock.Any()).
			Return(expectedTrips, nil)

		result := service.handleSearchTrips(context.Background(), args)

		assert.NotNil(t, result)
		assert.Nil(t, result["error"])
		assert.NotNil(t, result["trips"])
	})

	t.Run("Service Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTripService := mocks.NewMockTripServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			tripService: mockTripService,
		}

		args := map[string]any{
			"origin":      "Hà Nội",
			"destination": "Đà Nẵng",
		}

		mockTripService.EXPECT().
			SearchTrips(gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		result := service.handleSearchTrips(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"].(string), "Unable to search trips")
	})
}

func TestHandleGetTripDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripService := mocks.NewMockTripServiceClient(ctrl)
	service := &ChatbotServiceImpl{
		tripService: mockTripService,
	}

	tripID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		args := map[string]any{
			"trip_id": tripID.String(),
		}

		expectedTrip := &model.TripDetailResponse{
			ID:             tripID,
			DepartureTime:  time.Now(),
			AvailableSeats: 20,
			TotalSeats:     45,
			BasePrice:      250000,
			Route: &model.RouteDetail{
				Origin:      "Sài Gòn",
				Destination: "Đà Lạt",
			},
			Bus: &model.BusDetail{
				ID:           uuid.New(),
				LicensePlate: "51A-12345",
				TotalSeats:   45,
			},
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(expectedTrip, nil)

		result := service.handleGetTripDetails(context.Background(), args)

		assert.NotNil(t, result)
		assert.Nil(t, result["error"])
		assert.NotNil(t, result["trip"])
		assert.Contains(t, result["message"], "45 total seats")
		assert.Contains(t, result["message"], "20 available")
	})

	t.Run("Missing Trip ID", func(t *testing.T) {
		args := map[string]any{}

		result := service.handleGetTripDetails(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "trip_id is required")
	})

	t.Run("Invalid Trip ID Type", func(t *testing.T) {
		args := map[string]any{
			"trip_id": 123, // Should be string
		}

		result := service.handleGetTripDetails(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "trip_id is required")
	})

	t.Run("Service Error", func(t *testing.T) {
		args := map[string]any{
			"trip_id": tripID.String(),
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(nil, assert.AnError)

		result := service.handleGetTripDetails(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"].(string), "Unable to get trip details")
	})
}

func TestHandleGetAvailableSeats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripService := mocks.NewMockTripServiceClient(ctrl)
	service := &ChatbotServiceImpl{
		tripService: mockTripService,
	}

	tripID := uuid.New()

	t.Run("Success With Available Seats", func(t *testing.T) {
		args := map[string]any{
			"trip_id": tripID.String(),
		}

		seat1ID := uuid.New()
		seat2ID := uuid.New()
		bookedSeatID := uuid.New()

		expectedTrip := &model.TripDetailResponse{
			ID: tripID,
			Bus: &model.BusDetail{
				Seats: []model.SeatDetail{
					{
						ID:          seat1ID,
						SeatNumber:  "A1",
						Floor:       1,
						SeatType:    "normal",
						IsAvailable: true,
						Status:      &model.SeatStatus{IsBooked: false, IsLocked: false},
					},
					{
						ID:          seat2ID,
						SeatNumber:  "A2",
						Floor:       1,
						SeatType:    "vip",
						IsAvailable: true,
						Status:      &model.SeatStatus{IsBooked: false, IsLocked: false},
					},
					{
						ID:          bookedSeatID,
						SeatNumber:  "A3",
						Floor:       1,
						SeatType:    "normal",
						IsAvailable: false,
						Status:      &model.SeatStatus{IsBooked: true, IsLocked: false},
					},
				},
			},
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(expectedTrip, nil)

		result := service.handleGetAvailableSeats(context.Background(), args)

		assert.NotNil(t, result)
		assert.Nil(t, result["error"])
		assert.Equal(t, 2, result["total_available"])

		availableSeats, ok := result["available_seats"].([]map[string]any)
		require.True(t, ok)
		assert.Len(t, availableSeats, 2)
		assert.Equal(t, "A1", availableSeats[0]["seat_number"])
		assert.Equal(t, "A2", availableSeats[1]["seat_number"])
	})

	t.Run("No Available Seats", func(t *testing.T) {
		args := map[string]any{
			"trip_id": tripID.String(),
		}

		expectedTrip := &model.TripDetailResponse{
			ID: tripID,
			Bus: &model.BusDetail{
				Seats: []model.SeatDetail{
					{
						ID:          uuid.New(),
						SeatNumber:  "A1",
						IsAvailable: false,
						Status:      &model.SeatStatus{IsBooked: true, IsLocked: false},
					},
				},
			},
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(expectedTrip, nil)

		result := service.handleGetAvailableSeats(context.Background(), args)

		assert.NotNil(t, result)
		assert.Nil(t, result["error"])
		assert.Equal(t, 0, result["total_available"])
	})

	t.Run("Missing Trip ID", func(t *testing.T) {
		args := map[string]any{}

		result := service.handleGetAvailableSeats(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "trip_id is required")
	})

	t.Run("Service Error", func(t *testing.T) {
		args := map[string]any{
			"trip_id": tripID.String(),
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(nil, assert.AnError)

		result := service.handleGetAvailableSeats(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
	})
}

func TestHandleCreateGuestBooking(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTripService := mocks.NewMockTripServiceClient(ctrl)
		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			tripService:    mockTripService,
			bookingService: mockBookingService,
		}

		tripID := uuid.New()
		seatID1 := uuid.New()
		seatID2 := uuid.New()

		args := map[string]any{
			"trip_id":      tripID.String(),
			"seat_numbers": []any{"A1", "A2"},
			"full_name":    "Nguyễn Văn A",
			"email":        "test@example.com",
			"phone":        "0901234567",
			"passengers": []any{
				map[string]any{
					"name":        "Nguyễn Văn A",
					"phone":       "0901234567",
					"email":       "test@example.com",
					"seat_number": "A1",
				},
				map[string]any{
					"name":        "Nguyễn Văn B",
					"phone":       "0901234568",
					"email":       "test2@example.com",
					"seat_number": "A2",
				},
			},
		}

		tripDetails := &model.TripDetailResponse{
			ID: tripID,
			Bus: &model.BusDetail{
				Seats: []model.SeatDetail{
					{ID: seatID1, SeatNumber: "A1"},
					{ID: seatID2, SeatNumber: "A2"},
				},
			},
		}

		bookingID := uuid.New()
		bookingResponse := &model.BookingResponse{
			ID:         bookingID,
			Reference:  "ABC123XYZ",
			TotalPrice: 500000,
			Status:     "pending",
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(tripDetails, nil)

		mockBookingService.EXPECT().
			CreateGuestBooking(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, req *model.CreateGuestBookingRequest) (*model.BookingResponse, error) {
				assert.Equal(t, tripID, req.TripID)
				assert.Len(t, req.SeatIDs, 2)
				assert.Equal(t, "Nguyễn Văn A", req.FullName)
				assert.Equal(t, "test@example.com", req.Email)
				assert.Equal(t, "0901234567", req.Phone)
				assert.Len(t, req.Passengers, 2)
				return bookingResponse, nil
			})

		result := service.handleCreateGuestBooking(context.Background(), args, nil)

		assert.NotNil(t, result)
		assert.Nil(t, result["error"])
		assert.True(t, result["success"].(bool))

		booking, ok := result["booking"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "ABC123XYZ", booking["reference"])
		assert.Equal(t, float64(500000), booking["total_price"])
	})

	t.Run("Seat Not Found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTripService := mocks.NewMockTripServiceClient(ctrl)
		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			tripService:    mockTripService,
			bookingService: mockBookingService,
		}

		tripID := uuid.New()
		seatID1 := uuid.New()

		args := map[string]any{
			"trip_id":      tripID.String(),
			"seat_numbers": []any{"A1", "Z99"}, // Z99 doesn't exist
			"full_name":    "Nguyễn Văn A",
			"email":        "test@example.com",
			"phone":        "0901234567",
			"passengers": []any{
				map[string]any{
					"name":        "Nguyễn Văn A",
					"phone":       "0901234567",
					"email":       "test@example.com",
					"seat_number": "A1",
				},
			},
		}

		tripDetails := &model.TripDetailResponse{
			ID: tripID,
			Bus: &model.BusDetail{
				Seats: []model.SeatDetail{
					{ID: seatID1, SeatNumber: "A1"},
				},
			},
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(tripDetails, nil)

		result := service.handleCreateGuestBooking(context.Background(), args, nil)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"].(string), "Seat Z99 not found")
	})

	t.Run("Trip Service Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTripService := mocks.NewMockTripServiceClient(ctrl)
		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			tripService:    mockTripService,
			bookingService: mockBookingService,
		}

		tripID := uuid.New()

		args := map[string]any{
			"trip_id":      tripID.String(),
			"seat_numbers": []any{"A1"},
			"full_name":    "Test",
			"email":        "test@example.com",
			"phone":        "0901234567",
			"passengers":   []any{},
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(nil, assert.AnError)

		result := service.handleCreateGuestBooking(context.Background(), args, nil)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"].(string), "Trip not found")
	})

	t.Run("Booking Service Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTripService := mocks.NewMockTripServiceClient(ctrl)
		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			tripService:    mockTripService,
			bookingService: mockBookingService,
		}

		tripID := uuid.New()
		seatID1 := uuid.New()

		args := map[string]any{
			"trip_id":      tripID.String(),
			"seat_numbers": []any{"A1"},
			"full_name":    "Nguyễn Văn A",
			"email":        "test@example.com",
			"phone":        "0901234567",
			"passengers": []any{
				map[string]any{
					"name":        "Nguyễn Văn A",
					"phone":       "0901234567",
					"email":       "test@example.com",
					"seat_number": "A1",
				},
			},
		}

		tripDetails := &model.TripDetailResponse{
			ID: tripID,
			Bus: &model.BusDetail{
				Seats: []model.SeatDetail{
					{ID: seatID1, SeatNumber: "A1"},
				},
			},
		}

		mockTripService.EXPECT().
			GetTripByID(gomock.Any(), tripID.String()).
			Return(tripDetails, nil)

		mockBookingService.EXPECT().
			CreateGuestBooking(gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		result := service.handleCreateGuestBooking(context.Background(), args, nil)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"].(string), "Booking failed")
	})
}

func TestHandleCreatePaymentLink(t *testing.T) {
	t.Run("Missing Booking ID", func(t *testing.T) {
		service := &ChatbotServiceImpl{}
		args := map[string]any{}

		result := service.handleCreatePaymentLink(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "booking_id is required")
	})

	t.Run("Invalid Booking ID Format", func(t *testing.T) {
		service := &ChatbotServiceImpl{}
		args := map[string]any{
			"booking_id": "invalid-uuid",
		}

		result := service.handleCreatePaymentLink(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "Invalid booking ID format")
	})
}

func TestHandleCheckBookingStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
	service := &ChatbotServiceImpl{
		bookingService: mockBookingService,
	}

	t.Run("Success - Pending Booking", func(t *testing.T) {
		args := map[string]any{
			"reference": "ABC123XYZ",
			"email":     "test@example.com",
		}

		bookingResponse := &model.BookingResponse{
			Reference:  "ABC123XYZ",
			Status:     "pending",
			TotalPrice: 300000,
			CreatedAt:  time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
		}

		mockBookingService.EXPECT().
			GetBookingByReference(gomock.Any(), "ABC123XYZ", "test@example.com").
			Return(bookingResponse, nil)

		result := service.handleCheckBookingStatus(context.Background(), args)

		assert.NotNil(t, result)
		assert.Nil(t, result["error"])
		assert.Equal(t, "ABC123XYZ", result["reference"])
		assert.Equal(t, "pending", result["status"])
		assert.Equal(t, float64(300000), result["total_price"])
		assert.Contains(t, result["message"], "pending payment")
	})

	t.Run("Success - Confirmed Booking", func(t *testing.T) {
		args := map[string]any{
			"reference": "XYZ789ABC",
			"email":     "test@example.com",
		}

		bookingResponse := &model.BookingResponse{
			Reference:  "XYZ789ABC",
			Status:     "confirmed",
			TotalPrice: 450000,
			CreatedAt:  time.Now(),
			Transaction: &model.TransactionInfo{
				Status: "completed",
			},
		}

		mockBookingService.EXPECT().
			GetBookingByReference(gomock.Any(), "XYZ789ABC", "test@example.com").
			Return(bookingResponse, nil)

		result := service.handleCheckBookingStatus(context.Background(), args)

		assert.NotNil(t, result)
		assert.Nil(t, result["error"])
		assert.Equal(t, "confirmed", result["status"])
		assert.Equal(t, "completed", result["payment_status"])
		assert.Contains(t, result["message"], "confirmed")
	})

	t.Run("Success - Cancelled Booking", func(t *testing.T) {
		args := map[string]any{
			"reference": "CANCEL123",
			"email":     "test@example.com",
		}

		bookingResponse := &model.BookingResponse{
			Reference:  "CANCEL123",
			Status:     "cancelled",
			TotalPrice: 200000,
			CreatedAt:  time.Now(),
		}

		mockBookingService.EXPECT().
			GetBookingByReference(gomock.Any(), "CANCEL123", "test@example.com").
			Return(bookingResponse, nil)

		result := service.handleCheckBookingStatus(context.Background(), args)

		assert.NotNil(t, result)
		assert.Equal(t, "cancelled", result["status"])
		assert.Contains(t, result["message"], "cancelled")
	})

	t.Run("Missing Reference", func(t *testing.T) {
		args := map[string]any{
			"email": "test@example.com",
		}

		result := service.handleCheckBookingStatus(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "reference is required")
	})

	t.Run("Missing Email", func(t *testing.T) {
		args := map[string]any{
			"reference": "ABC123",
		}

		result := service.handleCheckBookingStatus(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "email is required")
	})

	t.Run("Service Error", func(t *testing.T) {
		args := map[string]any{
			"reference": "ABC123XYZ",
			"email":     "test@example.com",
		}

		mockBookingService.EXPECT().
			GetBookingByReference(gomock.Any(), "ABC123XYZ", "test@example.com").
			Return(nil, assert.AnError)

		result := service.handleCheckBookingStatus(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"].(string), "Unable to find booking")
	})
}

func TestNormalizeCityName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Sài Gòn variant",
			input:    "Sài Gòn",
			expected: "TP. Hồ Chí Minh",
		},
		{
			name:     "Saigon variant",
			input:    "Saigon",
			expected: "TP. Hồ Chí Minh",
		},
		{
			name:     "TPHCM variant",
			input:    "TPHCM",
			expected: "TP. Hồ Chí Minh",
		},
		{
			name:     "HCM variant",
			input:    "HCM",
			expected: "TP. Hồ Chí Minh",
		},
		{
			name:     "Hồ Chí Minh - no change",
			input:    "Hồ Chí Minh",
			expected: "Hồ Chí Minh",
		},
		{
			name:     "Hà Nội - no change",
			input:    "Hà Nội",
			expected: "Hà Nội",
		},
		{
			name:     "Đà Nẵng - no change",
			input:    "Đà Nẵng",
			expected: "Đà Nẵng",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Unknown city",
			input:    "Nha Trang",
			expected: "Nha Trang",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeCityName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleCreatePaymentLink_Comprehensive(t *testing.T) {
	t.Run("Success - Create Payment Link", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		mockPaymentService := mocks.NewMockPaymentServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			bookingService: mockBookingService,
			paymentService: mockPaymentService,
		}

		bookingID := uuid.New()
		args := map[string]any{
			"booking_id": bookingID.String(),
		}

		expiresAt := time.Now().Add(15 * time.Minute)
		bookingResponse := &model.BookingResponse{
			ID:         bookingID,
			Reference:  "ABC123XYZ",
			Status:     "pending",
			TotalPrice: 300000,
			ExpiresAt:  &expiresAt,
		}

		transactionResponse := &model.TransactionResponse{
			ID:          uuid.New(),
			BookingID:   bookingID,
			Amount:      300000,
			Currency:    "VND",
			Status:      "pending",
			CheckoutURL: "https://payment.example.com/checkout",
			QRCode:      "data:image/png;base64,abc123",
		}

		mockBookingService.EXPECT().
			GetBookingByID(gomock.Any(), bookingID.String()).
			Return(bookingResponse, nil)

		mockPaymentService.EXPECT().
			CreateTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, req *model.CreateTransactionRequest) (*model.TransactionResponse, error) {
				assert.Equal(t, bookingID, req.BookingID)
				assert.Equal(t, 300000, req.Amount)
				assert.Equal(t, "VND", req.Currency)
				return transactionResponse, nil
			})

		result := service.handleCreatePaymentLink(context.Background(), args)

		assert.NotNil(t, result)
		assert.Nil(t, result["error"])
		assert.True(t, result["success"].(bool))
		assert.Contains(t, result["checkout_url"], "payment.example.com")
	})

	t.Run("Booking Already Paid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			bookingService: mockBookingService,
		}

		bookingID := uuid.New()
		args := map[string]any{
			"booking_id": bookingID.String(),
		}

		bookingResponse := &model.BookingResponse{
			ID:         bookingID,
			Reference:  "ABC123XYZ",
			Status:     "confirmed", // Already paid
			TotalPrice: 300000,
		}

		mockBookingService.EXPECT().
			GetBookingByID(gomock.Any(), bookingID.String()).
			Return(bookingResponse, nil)

		result := service.handleCreatePaymentLink(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "already paid")
	})

	t.Run("Booking Cancelled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			bookingService: mockBookingService,
		}

		bookingID := uuid.New()
		args := map[string]any{
			"booking_id": bookingID.String(),
		}

		bookingResponse := &model.BookingResponse{
			ID:         bookingID,
			Reference:  "ABC123XYZ",
			Status:     "cancelled",
			TotalPrice: 300000,
		}

		mockBookingService.EXPECT().
			GetBookingByID(gomock.Any(), bookingID.String()).
			Return(bookingResponse, nil)

		result := service.handleCreatePaymentLink(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"], "cannot be paid")
	})

	t.Run("Booking Service Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			bookingService: mockBookingService,
		}

		bookingID := uuid.New()
		args := map[string]any{
			"booking_id": bookingID.String(),
		}

		mockBookingService.EXPECT().
			GetBookingByID(gomock.Any(), bookingID.String()).
			Return(nil, assert.AnError)

		result := service.handleCreatePaymentLink(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"].(string), "Cannot find booking")
	})

	t.Run("Payment Service Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
		mockPaymentService := mocks.NewMockPaymentServiceClient(ctrl)
		service := &ChatbotServiceImpl{
			bookingService: mockBookingService,
			paymentService: mockPaymentService,
		}

		bookingID := uuid.New()
		args := map[string]any{
			"booking_id": bookingID.String(),
		}

		expiresAt := time.Now().Add(15 * time.Minute)
		bookingResponse := &model.BookingResponse{
			ID:         bookingID,
			Reference:  "ABC123XYZ",
			Status:     "pending",
			TotalPrice: 300000,
			ExpiresAt:  &expiresAt,
		}

		mockBookingService.EXPECT().
			GetBookingByID(gomock.Any(), bookingID.String()).
			Return(bookingResponse, nil)

		mockPaymentService.EXPECT().
			CreateTransaction(gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		result := service.handleCreatePaymentLink(context.Background(), args)

		assert.NotNil(t, result)
		assert.NotNil(t, result["error"])
		assert.Contains(t, result["error"].(string), "Không thể tạo link thanh toán")
	})
}
