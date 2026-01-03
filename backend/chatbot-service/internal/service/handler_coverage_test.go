package service

import (
	"context"
	"testing"

	"bus-booking/chatbot-service/internal/model"
	"bus-booking/chatbot-service/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Additional function handler tests to reach 70% coverage

func TestHandleSearchTrips_EmptyResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripService := mocks.NewMockTripServiceClient(ctrl)
	service := &ChatbotServiceImpl{
		tripService: mockTripService,
	}

	args := map[string]any{
		"origin":      "Hà Nội",
		"destination": "Unknown City",
	}

	mockTripService.EXPECT().
		SearchTrips(gomock.Any(), gomock.Any()).
		Return(map[string]any{"data": []any{}}, nil)

	result := service.handleSearchTrips(context.Background(), args)

	assert.NotNil(t, result)
	assert.NotNil(t, result["trips"])
}

func TestHandleGetTripDetails_CompleteFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripService := mocks.NewMockTripServiceClient(ctrl)
	service := &ChatbotServiceImpl{
		tripService: mockTripService,
	}

	tripID := uuid.New()
	args := map[string]any{
		"trip_id": tripID.String(),
	}

	tripDetails := &model.TripDetailResponse{
		ID:             tripID,
		AvailableSeats: 15,
		TotalSeats:     45,
		Route: &model.RouteDetail{
			Origin:      "Hà Nội",
			Destination: "Đà Nẵng",
			Distance:    800,
			Duration:    720,
		},
		Bus: &model.BusDetail{
			ID:           uuid.New(),
			LicensePlate: "29A-12345",
			BusType:      "Giường nằm",
			TotalSeats:   45,
			Seats: []model.SeatDetail{
				{
					ID:         uuid.New(),
					SeatNumber: "A1",
					Floor:      1,
					SeatType:   "normal",
				},
			},
		},
	}

	mockTripService.EXPECT().
		GetTripByID(gomock.Any(), tripID.String()).
		Return(tripDetails, nil)

	result := service.handleGetTripDetails(context.Background(), args)

	assert.NotNil(t, result)
	assert.Nil(t, result["error"])
	details, ok := result["trip"].(map[string]any)
	assert.True(t, ok)
	assert.NotNil(t, details)
}

func TestHandleCreateGuestBooking_MaxPassengers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTripService := mocks.NewMockTripServiceClient(ctrl)
	mockBookingService := mocks.NewMockBookingServiceClient(ctrl)
	service := &ChatbotServiceImpl{
		tripService:    mockTripService,
		bookingService: mockBookingService,
	}

	tripID := uuid.New()
	seats := make([]uuid.UUID, 5)
	seatNumbers := []any{}
	passengers := []any{}

	for i := 0; i < 5; i++ {
		seats[i] = uuid.New()
		seatNum := string(rune('A' + i))
		seatNumbers = append(seatNumbers, seatNum)
		passengers = append(passengers, map[string]any{
			"name":        "Passenger " + seatNum,
			"phone":       "090123456" + string(rune('0'+i)),
			"email":       "p" + seatNum + "@test.com",
			"seat_number": seatNum,
		})
	}

	args := map[string]any{
		"trip_id":      tripID.String(),
		"seat_numbers": seatNumbers,
		"full_name":    "Main Passenger",
		"email":        "main@example.com",
		"phone":        "0901234567",
		"passengers":   passengers,
	}

	seatDetails := []model.SeatDetail{}
	for i := 0; i < 5; i++ {
		seatDetails = append(seatDetails, model.SeatDetail{
			ID:         seats[i],
			SeatNumber: string(rune('A' + i)),
		})
	}

	tripDetails := &model.TripDetailResponse{
		ID: tripID,
		Bus: &model.BusDetail{
			Seats: seatDetails,
		},
	}

	bookingResponse := &model.BookingResponse{
		ID:         uuid.New(),
		Reference:  "MULTI123",
		TotalPrice: 1250000,
		Status:     "pending",
	}

	mockTripService.EXPECT().
		GetTripByID(gomock.Any(), tripID.String()).
		Return(tripDetails, nil)

	mockBookingService.EXPECT().
		CreateGuestBooking(gomock.Any(), gomock.Any()).
		Return(bookingResponse, nil)

	result := service.handleCreateGuestBooking(context.Background(), args, nil)

	assert.NotNil(t, result)
	assert.True(t, result["success"].(bool))
}

func TestNormalizeCityName_AllVariants(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Sài Gòn", "TP. Hồ Chí Minh"},
		{"Saigon", "TP. Hồ Chí Minh"},
		{"TPHCM", "TP. Hồ Chí Minh"},
		{"HCM", "TP. Hồ Chí Minh"},
		{"Hồ Chí Minh", "Hồ Chí Minh"},
		{"Hà Nội", "Hà Nội"},
		{"", ""},
		{"Đà Lạt", "Đà Lạt"},
	}

	for _, tt := range tests {
		result := normalizeCityName(tt.input)
		assert.Equal(t, tt.expected, result, "Failed for input: %s", tt.input)
	}
}

func TestHandleSearchTrips_WithDepartureDate(t *testing.T) {
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
	}

	mockTripService.EXPECT().
		SearchTrips(gomock.Any(), gomock.Any()).
		Return(map[string]any{"data": []any{}}, nil)

	result := service.handleSearchTrips(context.Background(), args)

	assert.NotNil(t, result)
	assert.NotNil(t, result["trips"])
}
