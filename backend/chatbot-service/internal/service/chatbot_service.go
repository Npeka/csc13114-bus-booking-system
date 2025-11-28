package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"bus-booking/chatbot-service/config"
	"bus-booking/chatbot-service/internal/model"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type ChatbotService interface {
	ProcessMessage(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error)
	ExtractTripSearchParams(ctx context.Context, message string) (*model.TripSearchParams, error)
	GetFAQAnswer(ctx context.Context, question string) (string, error)
}

type ChatbotServiceImpl struct {
	openaiClient   *openai.Client
	config         *config.OpenAIConfig
	faqKnowledge   []model.FAQ
	tripService    TripServiceClient
	bookingService BookingServiceClient
}

func NewChatbotService(
	cfg *config.OpenAIConfig,
	externalCfg *config.ExternalConfig,
) ChatbotService {
	client := openai.NewClient(cfg.APIKey)

	return &ChatbotServiceImpl{
		openaiClient:   client,
		config:         cfg,
		faqKnowledge:   loadFAQs(),
		tripService:    NewTripServiceClient(externalCfg.TripServiceURL),
		bookingService: NewBookingServiceClient(externalCfg.BookingServiceURL),
	}
}

// ProcessMessage handles incoming chat messages
func (s *ChatbotServiceImpl) ProcessMessage(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error) {
	// Build conversation history
	messages := s.buildMessages(req)

	// Call OpenAI API
	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       s.config.Model,
			Messages:    messages,
			Temperature: s.config.Temperature,
			MaxTokens:   s.config.MaxTokens,
		},
	)

	if err != nil {
		log.Error().Err(err).Msg("OpenAI API call failed")
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	aiMessage := resp.Choices[0].Message.Content

	// Determine intent and action
	intent := s.detectIntent(aiMessage, req.Message)
	action := s.determineAction(intent, req.Context)

	// Handle different intents
	var responseData interface{}
	var suggestions []string

	switch intent {
	case "search_trip":
		params, err := s.ExtractTripSearchParams(ctx, req.Message)
		if err == nil {
			trips, err := s.tripService.SearchTrips(ctx, params)
			if err == nil {
				responseData = trips
				suggestions = []string{"Xem chi tiết", "Tìm chuyến khác", "Đặt vé"}
			}
		}
	case "faq":
		// Already handled by AI, no additional data needed
		suggestions = []string{"Tìm chuyến xe", "Xem chính sách", "Liên hệ hỗ trợ"}
	case "book_trip":
		if req.Context != nil && req.Context.SelectedTrip != nil {
			suggestions = []string{"Xác nhận đặt vé", "Chọn chuyến khác", "Hủy"}
		}
	}

	return &model.ChatResponse{
		Message:     aiMessage,
		Intent:      intent,
		Action:      action,
		Data:        responseData,
		Context:     req.Context,
		Suggestions: suggestions,
	}, nil
}

// ExtractTripSearchParams extracts search parameters from natural language
func (s *ChatbotServiceImpl) ExtractTripSearchParams(ctx context.Context, message string) (*model.TripSearchParams, error) {
	systemPrompt := `You are a travel assistant. Extract trip search parameters from the user's message.
Return ONLY a JSON object with these fields:
{
  "origin": "city name",
  "destination": "city name",
  "departure_date": "YYYY-MM-DD or empty",
  "passengers": number or 1
}

Examples:
- "Tìm xe từ Hà Nội đi Đà Nẵng ngày 30/11" -> {"origin": "Hà Nội", "destination": "Đà Nẵng", "departure_date": "2025-11-30", "passengers": 1}
- "Tôi muốn đi Sài Gòn từ Huế" -> {"origin": "Huế", "destination": "Sài Gòn", "departure_date": "", "passengers": 1}
`

	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: s.config.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
			Temperature: 0.3,
			MaxTokens:   200,
		},
	)

	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// Parse JSON response
	var params model.TripSearchParams
	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &params); err != nil {
		log.Error().Err(err).Str("content", content).Msg("Failed to parse search params")
		return nil, fmt.Errorf("failed to parse search parameters: %w", err)
	}

	return &params, nil
}

// GetFAQAnswer searches for FAQ answers
func (s *ChatbotServiceImpl) GetFAQAnswer(ctx context.Context, question string) (string, error) {
	questionLower := strings.ToLower(question)

	// Simple keyword matching
	for _, faq := range s.faqKnowledge {
		for _, keyword := range faq.Keywords {
			if strings.Contains(questionLower, strings.ToLower(keyword)) {
				return faq.Answer, nil
			}
		}
	}

	// If no exact match, use OpenAI with FAQ context
	systemPrompt := fmt.Sprintf(`You are a customer service assistant for a bus booking system.
Answer the user's question based on this FAQ knowledge:

%s

If the question is not covered in the FAQ, provide a helpful general answer.`, s.formatFAQs())

	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: s.config.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
			Temperature: 0.7,
			MaxTokens:   300,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// Helper methods

func (s *ChatbotServiceImpl) buildMessages(req *model.ChatRequest) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a helpful travel assistant for a bus booking system in Vietnam.
Help users:
1. Search for bus trips between cities
2. Book tickets
3. Answer questions about policies, schedules, and services

Be friendly, concise, and helpful. Respond in Vietnamese when the user speaks Vietnamese.`,
		},
	}

	// Add conversation history
	for _, msg := range req.History {
		role := openai.ChatMessageRoleUser
		if msg.Role == "assistant" {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Add current message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Message,
	})

	return messages
}

func (s *ChatbotServiceImpl) detectIntent(aiResponse, userMessage string) string {
	messageLower := strings.ToLower(userMessage)
	responseLower := strings.ToLower(aiResponse)

	// Check for trip search intent
	searchKeywords := []string{"tìm", "search", "chuyến", "trip", "xe", "bus", "từ", "đến", "đi"}
	for _, keyword := range searchKeywords {
		if strings.Contains(messageLower, keyword) &&
			(strings.Contains(messageLower, "từ") || strings.Contains(messageLower, "đến")) {
			return "search_trip"
		}
	}

	// Check for booking intent
	bookingKeywords := []string{"đặt", "book", "mua vé", "ticket"}
	for _, keyword := range bookingKeywords {
		if strings.Contains(messageLower, keyword) || strings.Contains(responseLower, keyword) {
			return "book_trip"
		}
	}

	// Check for FAQ intent
	faqKeywords := []string{"chính sách", "policy", "quy định", "rule", "giờ", "time", "hủy", "cancel", "đổi", "change"}
	for _, keyword := range faqKeywords {
		if strings.Contains(messageLower, keyword) {
			return "faq"
		}
	}

	return "general"
}

func (s *ChatbotServiceImpl) determineAction(intent string, context *model.ChatContext) string {
	switch intent {
	case "search_trip":
		return "show_trips"
	case "book_trip":
		if context != nil && context.SelectedTrip != nil {
			return "show_booking_form"
		}
		return "select_trip"
	case "faq":
		return "show_faq"
	default:
		return "continue_conversation"
	}
}

func (s *ChatbotServiceImpl) formatFAQs() string {
	var sb strings.Builder
	for _, faq := range s.faqKnowledge {
		sb.WriteString(fmt.Sprintf("Q: %s\nA: %s\n\n", faq.Question, faq.Answer))
	}
	return sb.String()
}

// loadFAQs loads FAQ knowledge base
func loadFAQs() []model.FAQ {
	return []model.FAQ{
		{
			Question: "Chính sách hủy vé như thế nào?",
			Answer:   "Bạn có thể hủy vé trước 24 giờ và được hoàn 70% giá vé, trước 12 giờ hoàn 50%, và trước 6 giờ hoàn 30%. Hủy vé trong vòng 6 giờ không được hoàn tiền.",
			Keywords: []string{"hủy", "cancel", "hoàn tiền", "refund"},
		},
		{
			Question: "Làm sao để đổi vé?",
			Answer:   "Bạn có thể đổi vé sang chuyến khác cùng tuyến đường miễn phí nếu đổi trước 12 giờ. Đổi vé trong vòng 12 giờ trước giờ khởi hành phải trả thêm 10% giá vé.",
			Keywords: []string{"đổi", "change", "reschedule"},
		},
		{
			Question: "Tôi có thể mang hành lý bao nhiêu kg?",
			Answer:   "Mỗi hành khách được mang tối đa 20kg hành lý miễn phí. Hành lý quá cân sẽ tính phí 10,000 VNĐ/kg.",
			Keywords: []string{"hành lý", "luggage", "baggage", "kg"},
		},
		{
			Question: "Xe có wifi không?",
			Answer:   "Tất cả xe của chúng tôi đều có wifi miễn phí, ổ cắm sạc điện thoại, và điều hòa.",
			Keywords: []string{"wifi", "internet", "tiện ích", "amenities"},
		},
		{
			Question: "Thanh toán như thế nào?",
			Answer:   "Chúng tôi chấp nhận thanh toán qua PayOS (chuyển khoản ngân hàng, ví điện tử), thẻ tín dụng, hoặc thanh toán tại quầy.",
			Keywords: []string{"thanh toán", "payment", "pay"},
		},
		{
			Question: "Xe khởi hành ở đâu?",
			Answer:   "Điểm đón khách sẽ được hiển thị trong chi tiết chuyến xe. Bạn cần có mặt trước 15 phút so với giờ khởi hành.",
			Keywords: []string{"điểm đón", "pickup", "bến xe", "station"},
		},
	}
}
