package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/bus-booking/internal/usecases"
)

type ChatbotHandler struct {
	chatbotUsecase *usecases.ChatbotUsecase
}

func NewChatbotHandler(chatbotUsecase *usecases.ChatbotUsecase) *ChatbotHandler {
	return &ChatbotHandler{
		chatbotUsecase: chatbotUsecase,
	}
}

// SendMessage godoc
// @Summary Send message to AI chatbot
// @Description Send a message to the AI assistant and get response
// @Tags chatbot
// @Accept json
// @Produce json
// @Param request body ChatbotRequest true "Chat message"
// @Success 200 {object} ChatbotResponse
// @Failure 400 {object} ErrorResponse
// @Router /chatbot/message [post]
func (h *ChatbotHandler) SendMessage(c *gin.Context) {
	var req ChatbotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Default to Vietnamese if no language specified
	if req.Language == "" {
		req.Language = "vi"
	}

	response, err := h.chatbotUsecase.SendMessage(c.Request.Context(), req.Message, req.ConversationID, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetHistory godoc
// @Summary Get conversation history
// @Description Retrieve chat history for a conversation
// @Tags chatbot
// @Produce json
// @Param conversation_id query string true "Conversation ID"
// @Success 200 {object} HistoryResponse
// @Failure 400 {object} ErrorResponse
// @Router /chatbot/history [get]
func (h *ChatbotHandler) GetHistory(c *gin.Context) {
	conversationID := c.Query("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "conversation_id is required"})
		return
	}

	history, err := h.chatbotUsecase.GetHistory(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert to response format
	messages := make([]Message, len(history))
	for i, msg := range history {
		messages[i] = Message{
			Role:      msg.Role,
			Content:   msg.Content,
			Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, HistoryResponse{Messages: messages})
}

// Request/Response types
type ChatbotRequest struct {
	Message        string `json:"message" binding:"required"`
	ConversationID string `json:"conversation_id,omitempty"`
	Language       string `json:"language,omitempty"` // "vi" or "en"
}

type ChatbotResponse struct {
	Message        string            `json:"message"`
	Suggestions    []string          `json:"suggestions,omitempty"`
	QuickActions   []QuickAction     `json:"quick_actions,omitempty"`
	ConversationID string            `json:"conversation_id"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type QuickAction struct {
	Label  string `json:"label"`
	Action string `json:"action"`
	Data   string `json:"data,omitempty"`
}

type HistoryResponse struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
