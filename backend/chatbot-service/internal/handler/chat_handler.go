package handler

import (
	"bus-booking/chatbot-service/internal/model"
	"bus-booking/chatbot-service/internal/service"
	"bus-booking/shared/ginext"

	"github.com/rs/zerolog/log"
)

type ChatHandler interface {
	Chat(r *ginext.Request) (*ginext.Response, error)
	ExtractSearchParams(r *ginext.Request) (*ginext.Response, error)
	GetFAQ(r *ginext.Request) (*ginext.Response, error)
}

type ChatHandlerImpl struct {
	chatbotService service.ChatbotService
}

func NewChatHandler(chatbotService service.ChatbotService) ChatHandler {
	return &ChatHandlerImpl{
		chatbotService: chatbotService,
	}
}

// Chat godoc
// @Summary Chat with AI assistant
// @Description Send a message to the chatbot and get an AI-powered response
// @Tags chatbot
// @Accept json
// @Produce json
// @Param request body model.ChatRequest true "Chat request"
// @Success 200 {object} ginext.Response{data=model.ChatResponse}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/chat [post]
func (h *ChatHandlerImpl) Chat(r *ginext.Request) (*ginext.Response, error) {
	var req model.ChatRequest
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("Invalid request body")
		return nil, ginext.NewBadRequestError(err.Error())
	}

	response, err := h.chatbotService.ProcessMessage(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process chat message")
		return nil, err
	}

	return ginext.NewSuccessResponse(response, "Message processed successfully"), nil
}

// ExtractSearchParams godoc
// @Summary Extract trip search parameters
// @Description Extract trip search parameters from natural language
// @Tags chatbot
// @Accept json
// @Produce json
// @Param message query string true "Natural language message"
// @Success 200 {object} ginext.Response{data=model.TripSearchParams}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/chat/extract-search [get]
func (h *ChatHandlerImpl) ExtractSearchParams(r *ginext.Request) (*ginext.Response, error) {
	message := r.GinCtx.Query("message")
	if message == "" {
		return nil, ginext.NewBadRequestError("message is required")
	}

	params, err := h.chatbotService.ExtractTripSearchParams(r.Context(), message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to extract search parameters")
		return nil, err
	}

	return ginext.NewSuccessResponse(params, "Search parameters extracted successfully"), nil
}

// GetFAQ godoc
// @Summary Get FAQ answer
// @Description Get answer to frequently asked questions
// @Tags chatbot
// @Accept json
// @Produce json
// @Param question query string true "Question"
// @Success 200 {object} ginext.Response{data=string}
// @Failure 400 {object} ginext.Response
// @Failure 500 {object} ginext.Response
// @Router /api/v1/chat/faq [get]
func (h *ChatHandlerImpl) GetFAQ(r *ginext.Request) (*ginext.Response, error) {
	question := r.GinCtx.Query("question")
	if question == "" {
		return nil, ginext.NewBadRequestError("question is required")
	}

	answer, err := h.chatbotService.GetFAQAnswer(r.Context(), question)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get FAQ answer")
		return nil, err
	}

	return ginext.NewSuccessResponse(answer, "FAQ answer retrieved successfully"), nil
}
