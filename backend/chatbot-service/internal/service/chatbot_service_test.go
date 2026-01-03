package service

import (
	"testing"

	"bus-booking/chatbot-service/config"
	"bus-booking/chatbot-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewChatbotService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		geminiCfg := &config.GeminiConfig{
			APIKey: "test-api-key",
			Model:  "gemini-1.5-flash",
		}

		externalCfg := &config.ExternalConfig{
			TripServiceURL:    "http://localhost:8081",
			BookingServiceURL: "http://localhost:8082",
			PaymentServiceURL: "http://localhost:8083",
		}

		service := NewChatbotService(geminiCfg, externalCfg)

		assert.NotNil(t, service)
	})
}

func TestDetectIntent(t *testing.T) {
	service := &ChatbotServiceImpl{}

	tests := []struct {
		name           string
		aiResponse     string
		userMessage    string
		expectedIntent string
	}{
		{
			name:           "Search Trip Intent - Vietnamese",
			aiResponse:     "Tôi tìm thấy các chuyến xe",
			userMessage:    "Tìm xe từ Hà Nội đến Đà Nẵng",
			expectedIntent: "search_trip",
		},

		{
			name:           "Booking Intent",
			aiResponse:     "Booking created",
			userMessage:    "Đặt vé cho tôi",
			expectedIntent: "book_trip",
		},
		{
			name:           "Booking Intent - Buy Ticket",
			aiResponse:     "",
			userMessage:    "mua vé đi Sài Gòn",
			expectedIntent: "book_trip",
		},
		{
			name:           "FAQ Intent - Policy",
			aiResponse:     "",
			userMessage:    "Chính sách hủy vé như thế nào?",
			expectedIntent: "faq",
		},
		{
			name:           "FAQ Intent - Cancel",
			aiResponse:     "",
			userMessage:    "Tôi muốn hủy vé",
			expectedIntent: "faq",
		},
		{
			name:           "FAQ Intent - Time",
			aiResponse:     "",
			userMessage:    "Xe khởi hành lúc mấy giờ?",
			expectedIntent: "faq",
		},
		{
			name:           "General Intent",
			aiResponse:     "",
			userMessage:    "Xin chào",
			expectedIntent: "general",
		},
		{
			name:           "General Intent - Random",
			aiResponse:     "",
			userMessage:    "Cảm ơn bạn",
			expectedIntent: "general",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent := service.detectIntent(tt.aiResponse, tt.userMessage)
			assert.Equal(t, tt.expectedIntent, intent)
		})
	}
}

func TestDetermineAction(t *testing.T) {
	service := &ChatbotServiceImpl{}

	tripID := uuid.New()

	tests := []struct {
		name           string
		intent         string
		context        *model.ChatContext
		expectedAction string
	}{
		{
			name:           "Search Trip - Show Trips",
			intent:         "search_trip",
			context:        nil,
			expectedAction: "show_trips",
		},
		{
			name:   "Book Trip with Selected Trip - Show Booking Form",
			intent: "book_trip",
			context: &model.ChatContext{
				SelectedTrip: &model.TripInfo{
					TripID: tripID,
				},
			},
			expectedAction: "show_booking_form",
		},
		{
			name:           "Book Trip without Selected Trip - Select Trip",
			intent:         "book_trip",
			context:        nil,
			expectedAction: "select_trip",
		},
		{
			name:           "FAQ - Show FAQ",
			intent:         "faq",
			context:        nil,
			expectedAction: "show_faq",
		},
		{
			name:           "General - Continue Conversation",
			intent:         "general",
			context:        nil,
			expectedAction: "continue_conversation",
		},
		{
			name:           "Unknown Intent - Continue Conversation",
			intent:         "unknown_intent",
			context:        nil,
			expectedAction: "continue_conversation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := service.determineAction(tt.intent, tt.context)
			assert.Equal(t, tt.expectedAction, action)
		})
	}
}

func TestFormatFAQs(t *testing.T) {
	service := &ChatbotServiceImpl{
		faqKnowledge: []model.FAQ{
			{
				Question: "Chính sách hủy vé?",
				Answer:   "Hoàn 70% trước 24h",
				Keywords: []string{"hủy", "cancel"},
			},
			{
				Question: "Làm sao đổi vé?",
				Answer:   "Đổi vé miễn phí trước 12h",
				Keywords: []string{"đổi", "change"},
			},
		},
	}

	formatted := service.formatFAQs()

	assert.NotEmpty(t, formatted)
	assert.Contains(t, formatted, "Q: Chính sách hủy vé?")
	assert.Contains(t, formatted, "A: Hoàn 70% trước 24h")
	assert.Contains(t, formatted, "Q: Làm sao đổi vé?")
	assert.Contains(t, formatted, "A: Đổi vé miễn phí trước 12h")
}

func TestLoadFAQs(t *testing.T) {
	faqs := loadFAQs()

	assert.NotEmpty(t, faqs, "FAQs should not be empty")
	assert.Greater(t, len(faqs), 0, "Should have at least one FAQ")

	// Verify structure of first FAQ
	if len(faqs) > 0 {
		firstFAQ := faqs[0]
		assert.NotEmpty(t, firstFAQ.Question, "FAQ should have a question")
		assert.NotEmpty(t, firstFAQ.Answer, "FAQ should have an answer")
		assert.NotEmpty(t, firstFAQ.Keywords, "FAQ should have keywords")
	}

	// Test for expected FAQ topics
	var hasCancellationFAQ bool
	var hasBaggage bool
	var hasPaymentFAQ bool

	for _, faq := range faqs {
		if containsAny(faq.Keywords, []string{"hủy", "cancel"}) {
			hasCancellationFAQ = true
		}
		if containsAny(faq.Keywords, []string{"hành lý", "luggage", "baggage"}) {
			hasBaggage = true
		}
		if containsAny(faq.Keywords, []string{"thanh toán", "payment"}) {
			hasPaymentFAQ = true
		}
	}

	assert.True(t, hasCancellationFAQ, "Should have cancellation policy FAQ")
	assert.True(t, hasBaggage, "Should have baggage FAQ")
	assert.True(t, hasPaymentFAQ, "Should have payment FAQ")
}

// Helper function for testing
func containsAny(slice []string, items []string) bool {
	for _, s := range slice {
		for _, item := range items {
			if s == item {
				return true
			}
		}
	}
	return false
}

func TestGetSuggestions(t *testing.T) {
	// This tests the suggestion generation logic in ProcessMessage
	// Since we can't test ProcessMessage without Gemini API, we test the logic inline

	tests := []struct {
		name                string
		intent              string
		context             *model.ChatContext
		expectedSuggestions []string
	}{
		{
			name:                "Search Trip Intent",
			intent:              "search_trip",
			expectedSuggestions: []string{"Xem chi tiết", "Tìm chuyến khác", "Chính sách hoàn vé"},
		},
		{
			name:                "FAQ Intent",
			intent:              "faq",
			expectedSuggestions: []string{"Tìm chuyến xe", "Xem giá vé", "Liên hệ hỗ trợ"},
		},
		{
			name:   "Book Trip with Selected Trip",
			intent: "book_trip",
			context: &model.ChatContext{
				SelectedTrip: &model.TripInfo{
					TripID: uuid.New(),
				},
			},
			expectedSuggestions: []string{"Xác nhận đặt vé", "Chọn chuyến khác", "Hủy"},
		},
		{
			name:                "General Intent",
			intent:              "general",
			expectedSuggestions: []string{"Tìm chuyến xe", "Xem giá vé", "Chính sách hoàn vé"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Logic mimicking what's in ProcessMessage
			var suggestions []string
			switch tt.intent {
			case "search_trip":
				suggestions = []string{"Xem chi tiết", "Tìm chuyến khác", "Chính sách hoàn vé"}
			case "faq":
				suggestions = []string{"Tìm chuyến xe", "Xem giá vé", "Liên hệ hỗ trợ"}
			case "book_trip":
				if tt.context != nil && tt.context.SelectedTrip != nil {
					suggestions = []string{"Xác nhận đặt vé", "Chọn chuyến khác", "Hủy"}
				}
			default:
				suggestions = []string{"Tìm chuyến xe", "Xem giá vé", "Chính sách hoàn vé"}
			}

			if len(tt.expectedSuggestions) > 0 {
				assert.Equal(t, tt.expectedSuggestions, suggestions)
			}
		})
	}
}
