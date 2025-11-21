package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/bus-booking/internal/infrastructure/chatbot"
)

type ChatbotUsecase struct {
	chatbot chatbot.Chatbot
}

func NewChatbotUsecase(bot chatbot.Chatbot) *ChatbotUsecase {
	return &ChatbotUsecase{
		chatbot: bot,
	}
}

// SendMessage processes a user message and returns AI response
func (uc *ChatbotUsecase) SendMessage(ctx context.Context, message string, conversationID string, language string) (*chatbot.ChatResponse, error) {
	// Generate conversation ID if not provided
	if conversationID == "" {
		conversationID = uuid.New().String()
	}

	// Get response from chatbot
	response, err := uc.chatbot.SendMessage(ctx, message, conversationID, language)
	if err != nil {
		return nil, err
	}

	// TODO: Save conversation to database for history
	// conversationRepo.SaveMessage(ctx, conversationID, "user", message)
	// conversationRepo.SaveMessage(ctx, conversationID, "assistant", response.Message)

	return response, nil
}

// GetHistory retrieves conversation history
func (uc *ChatbotUsecase) GetHistory(ctx context.Context, conversationID string) ([]chatbot.Message, error) {
	return uc.chatbot.GetConversationHistory(ctx, conversationID)
}
