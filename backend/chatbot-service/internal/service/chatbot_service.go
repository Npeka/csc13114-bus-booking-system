package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"bus-booking/chatbot-service/config"
	"bus-booking/chatbot-service/internal/model"

	"github.com/rs/zerolog/log"
	"google.golang.org/genai"
)

type ChatbotService interface {
	ProcessMessage(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error)
	ExtractTripSearchParams(ctx context.Context, message string) (*model.TripSearchParams, error)
	GetFAQAnswer(ctx context.Context, question string) (string, error)
}

type ChatbotServiceImpl struct {
	genaiClient    *genai.Client
	config         *config.GeminiConfig
	faqKnowledge   []model.FAQ
	tripService    TripServiceClient
	bookingService BookingServiceClient
	paymentService PaymentServiceClient // NEW: Payment service client
}

func NewChatbotService(
	cfg *config.GeminiConfig,
	external *config.ExternalConfig,
) ChatbotService {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Gemini client")
	}

	return &ChatbotServiceImpl{
		genaiClient:    client,
		config:         cfg,
		faqKnowledge:   loadFAQs(),
		tripService:    NewTripServiceClient(external.TripServiceURL),
		bookingService: NewBookingServiceClient(external.BookingServiceURL),
		paymentService: NewPaymentServiceClient(external.PaymentServiceURL), // NEW
	}
}

// ProcessMessage handles incoming chat messages with Gemini Function Calling
func (s *ChatbotServiceImpl) ProcessMessage(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error) {
	// Define function declarations
	searchTripsFunc := &genai.FunctionDeclaration{
		Name:        "searchTrips",
		Description: "Search for bus trips between cities. Use this to find available trips when the user asks about routes, schedules, or availability.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"origin": {
					Type:        genai.TypeString,
					Description: "Origin city name (e.g., 'Sài Gòn', 'Hà Nội', 'Đà Nẵng')",
				},
				"destination": {
					Type:        genai.TypeString,
					Description: "Destination city name (e.g., 'Sài Gòn', 'Hà Nội', 'Đà Lạt')",
				},
				"departure_date": {
					Type:        genai.TypeString,
					Description: "Departure date in YYYY-MM-DD format. Leave empty if not specified.",
				},
				"passengers": {
					Type:        genai.TypeInteger,
					Description: "Number of passengers. Default is 1.",
				},
			},
			Required: []string{"origin", "destination"},
		},
	}

	getTripDetailsFunc := &genai.FunctionDeclaration{
		Name:        "getTripDetails",
		Description: "Get detailed information about a specific trip including seat map and availability. Use when user wants to see trip details or select seats.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"trip_id": {
					Type:        genai.TypeString,
					Description: "The UUID of the trip to get details for",
				},
			},
			Required: []string{"trip_id"},
		},
	}

	createGuestBookingFunc := &genai.FunctionDeclaration{
		Name:        "createGuestBooking",
		Description: "Create a booking for a guest user with passenger details and seat selection. Use ONLY when you have trip_id, seat numbers, and complete passenger information (name, phone, email).",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"trip_id": {
					Type:        genai.TypeString,
					Description: "The UUID of the trip",
				},
				"seat_numbers": {
					Type:        genai.TypeArray,
					Description: "Array of seat numbers (e.g., ['A1', 'A2'])",
					Items:       &genai.Schema{Type: genai.TypeString},
				},
				"full_name": {
					Type:        genai.TypeString,
					Description: "Full name of the primary passenger",
				},
				"email": {
					Type:        genai.TypeString,
					Description: "Email address for booking confirmation",
				},
				"phone": {
					Type:        genai.TypeString,
					Description: "Phone number for contact",
				},
				"passengers": {
					Type:        genai.TypeArray,
					Description: "Array of passenger details matching the number of seats",
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"name":        {Type: genai.TypeString, Description: "Passenger name"},
							"phone":       {Type: genai.TypeString, Description: "Passenger phone"},
							"email":       {Type: genai.TypeString, Description: "Passenger email"},
							"seat_number": {Type: genai.TypeString, Description: "Assigned seat number"},
						},
					},
				},
			},
			Required: []string{"trip_id", "seat_numbers", "full_name", "email", "phone", "passengers"},
		},
	}

	getAvailableSeatsFunc := &genai.FunctionDeclaration{
		Name:        "getAvailableSeats",
		Description: "Get only the available (not booked or locked) seats for a trip. Use when user wants to see which seats they can choose.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"trip_id": {
					Type:        genai.TypeString,
					Description: "The UUID of the trip to check seat availability",
				},
			},
			Required: []string{"trip_id"},
		},
	}

	createPaymentLinkFunc := &genai.FunctionDeclaration{
		Name:        "createPaymentLink",
		Description: "Generate a payment link for a booking. Use when user wants to pay or after booking is created. Requires booking_id from a previously created booking.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"booking_id": {
					Type:        genai.TypeString,
					Description: "The UUID of the booking to create payment for",
				},
			},
			Required: []string{"booking_id"},
		},
	}

	checkBookingStatusFunc := &genai.FunctionDeclaration{
		Name:        "checkBookingStatus",
		Description: "Check the status of a booking using the booking reference code and email. Use when user asks to check their booking or payment status.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"reference": {
					Type:        genai.TypeString,
					Description: "The booking reference code (e.g., 'ABC123XYZ')",
				},
				"email": {
					Type:        genai.TypeString,
					Description: "Email address used for the booking",
				},
			},
			Required: []string{"reference", "email"},
		},
	}

	// Create system instruction with booking flow rules
	systemInstruction := &genai.Content{
		Parts: []*genai.Part{
			{Text: `You are a helpful assistant for BusTicket.vn, a bus booking platform in Vietnam. You can help users search trips AND complete bookings.

CRITICAL RULES:
1. ALWAYS use the searchTrips function to retrieve real trip data when users ask about routes, schedules, or availability
2. NEVER make up or hallucinate trip information (times, prices, routes, seat numbers) - always call functions to get real data
3. When user selects a trip, call getTripDetails to show seat map and availability
4. When user provides seat selection and passenger info, call createGuestBooking to complete the booking
5. For general questions (policies, FAQs, contact info), you can answer directly using your knowledge

BOOKING FLOW:
1. User searches trips → call searchTrips
2. User selects a trip → call getTripDetails to show seat map
3. User chooses seats and provides passenger info → call createGuestBooking
4. After booking created → inform user of reference code and total price

Your responsibilities:
- Search for bus trips (searchTrips)
- Show trip details with seat maps (getTripDetails)
- Collect passenger information step-by-step if incomplete
- Create bookings (createGuestBooking)
- Answer questions about policies, cancellation, refunds

DATA COLLECTION RULES:
- For createGuestBooking, you MUST have: trip_id, seat_numbers, full_name, email, phone, and passengers array
- If user hasn't provided all info, ASK for missing details one by one
- Passengers array must match the number of seats selected
- Each passenger needs: name, phone, email, seat_number

Always respond in Vietnamese when the user speaks Vietnamese. Be friendly, clear, and helpful.`},
		},
	}

	// Build conversation history
	history := []*genai.Content{}
	for _, msg := range req.History {
		role := msg.Role
		if role == "assistant" {
			role = "model" // Gemini uses "model" instead of "assistant"
		}
		history = append(history, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: msg.Content}},
		})
	}

	// Add user message
	history = append(history, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: req.Message}},
	})

	// Configure generation with tools (now includes 6 functions)
	maxTokens := int32(s.config.MaxTokens)
	if s.config.MaxTokens > 2147483647 {
		maxTokens = 2147483647 // Max int32 value
	}
	genConfig := &genai.GenerateContentConfig{
		Temperature:       &s.config.Temperature,
		MaxOutputTokens:   maxTokens,
		SystemInstruction: systemInstruction,
		Tools:             []*genai.Tool{{FunctionDeclarations: []*genai.FunctionDeclaration{searchTripsFunc, getTripDetailsFunc, getAvailableSeatsFunc, createGuestBookingFunc, createPaymentLinkFunc, checkBookingStatusFunc}}},
	}

	// Call Gemini API
	resp, err := s.genaiClient.Models.GenerateContent(ctx, s.config.Model, history, genConfig)
	if err != nil {
		log.Error().Err(err).Msg("Gemini API call failed")
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	candidate := resp.Candidates[0]

	// Check if Gemini wants to call a function
	for _, part := range candidate.Content.Parts {
		if part.FunctionCall != nil {
			fc := part.FunctionCall
			log.Info().Str("function", fc.Name).Msg("Gemini requested function call")

			// Execute function call and get response
			var funcResp map[string]any
			switch fc.Name {
			case "searchTrips":
				funcResp = s.handleSearchTrips(ctx, fc.Args)
			case "getTripDetails":
				funcResp = s.handleGetTripDetails(ctx, fc.Args)
			case "getAvailableSeats":
				funcResp = s.handleGetAvailableSeats(ctx, fc.Args)
			case "createGuestBooking":
				funcResp = s.handleCreateGuestBooking(ctx, fc.Args, req.Context)
			case "createPaymentLink":
				funcResp = s.handleCreatePaymentLink(ctx, fc.Args)
			case "checkBookingStatus":
				funcResp = s.handleCheckBookingStatus(ctx, fc.Args)
			default:
				log.Warn().Str("function", fc.Name).Msg("Unknown function call")
				funcResp = map[string]any{"error": "Unknown function"}
			}

			// Add function response to history and call again
			history = append(history, candidate.Content)
			history = append(history, &genai.Content{
				Role: "function",
				Parts: []*genai.Part{
					{
						FunctionResponse: &genai.FunctionResponse{
							Name:     fc.Name,
							Response: funcResp,
						},
					},
				},
			})

			// Call Gemini again with function response
			resp, err = s.genaiClient.Models.GenerateContent(ctx, s.config.Model, history, genConfig)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get final response from Gemini")
				return nil, fmt.Errorf("failed to get final AI response: %w", err)
			}

			if len(resp.Candidates) > 0 {
				candidate = resp.Candidates[0]
			}
		}
	}

	// Extract text response
	var aiMessage string
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			aiMessage += part.Text
		}
	}

	if aiMessage == "" {
		aiMessage = "Xin lỗi, tôi không thể xử lý yêu cầu của bạn lúc này."
	}

	// Determine intent and action for frontend
	intent := s.detectIntent(aiMessage, req.Message)
	action := s.determineAction(intent, req.Context)

	// Prepare suggestions based on intent
	var suggestions []string
	switch intent {
	case "search_trip":
		suggestions = []string{"Xem chi tiết", "Tìm chuyến khác", "Chính sách hoàn vé"}
	case "faq":
		suggestions = []string{"Tìm chuyến xe", "Xem giá vé", "Liên hệ hỗ trợ"}
	case "book_trip":
		if req.Context != nil && req.Context.SelectedTrip != nil {
			suggestions = []string{"Xác nhận đặt vé", "Chọn chuyến khác", "Hủy"}
		}
	default:
		suggestions = []string{"Tìm chuyến xe", "Xem giá vé", "Chính sách hoàn vé"}
	}

	return &model.ChatResponse{
		Message:     aiMessage,
		Intent:      intent,
		Action:      action,
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

	temp := float32(0.3)
	config := &genai.GenerateContentConfig{
		Temperature:     &temp,
		MaxOutputTokens: int32(200),
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
	}

	resp, err := s.genaiClient.Models.GenerateContent(ctx, s.config.Model, []*genai.Content{
		{Role: "user", Parts: []*genai.Part{{Text: message}}},
	}, config)
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	// Extract text from response
	var content string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			content += part.Text
		}
	}

	content = strings.TrimSpace(content)

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// Parse JSON response
	var params model.TripSearchParams
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

	// If no exact match, use Gemini with FAQ context
	systemPrompt := fmt.Sprintf(`You are a customer service assistant for a bus booking system.
Answer the user's question based on this FAQ knowledge:

%s

If the question is not covered in the FAQ, provide a helpful general answer.`, s.formatFAQs())

	temp := float32(0.7)
	config := &genai.GenerateContentConfig{
		Temperature:     &temp,
		MaxOutputTokens: int32(300),
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
	}

	resp, err := s.genaiClient.Models.GenerateContent(ctx, s.config.Model, []*genai.Content{
		{Role: "user", Parts: []*genai.Part{{Text: question}}},
	}, config)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Extract text from response
	var answer string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			answer += part.Text
		}
	}

	return answer, nil
}

// Helper methods

func (s *ChatbotServiceImpl) detectIntent(aiResponse, userMessage string) string {
	messageLower := strings.ToLower(userMessage)

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
		if strings.Contains(messageLower, keyword) {
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
